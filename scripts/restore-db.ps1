#!/usr/bin/env pwsh
# Restore PostgreSQL Database from Backup
# Restores the database from a backup file

param(
    [Parameter(Mandatory=$false)]
    [string]$BackupFile,
    
    [Parameter(Mandatory=$false)]
    [switch]$Force
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Database Restore" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# If no backup file specified, show available backups and prompt
if ([string]::IsNullOrEmpty($BackupFile)) {
    $backupDir = "db-backups"
    
    if (!(Test-Path $backupDir)) {
        Write-Host "[X] No backup directory found: $backupDir" -ForegroundColor Red
        Write-Host "`nTo create a backup, run:" -ForegroundColor Cyan
        Write-Host "  .\scripts\backup-db.ps1" -ForegroundColor Yellow
        exit 1
    }
    
    $backups = Get-ChildItem -Path $backupDir -Filter "*.sql" | Sort-Object LastWriteTime -Descending
    
    if ($backups.Count -eq 0) {
        Write-Host "[X] No backup files found in $backupDir" -ForegroundColor Red
        Write-Host "`nTo create a backup, run:" -ForegroundColor Cyan
        Write-Host "  .\scripts\backup-db.ps1" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "Available backups:" -ForegroundColor Cyan
    Write-Host ""
    for ($i = 0; $i -lt $backups.Count; $i++) {
        $backup = $backups[$i]
        $size = if ($backup.Length -gt 1MB) {
            "$([math]::Round($backup.Length / 1MB, 2)) MB"
        } elseif ($backup.Length -gt 1KB) {
            "$([math]::Round($backup.Length / 1KB, 2)) KB"
        } else {
            "$($backup.Length) bytes"
        }
        $index = $i + 1
        Write-Host "  [$index] $($backup.Name)" -ForegroundColor Yellow
        Write-Host "      Size: $size | Date: $($backup.LastWriteTime)" -ForegroundColor Gray
    }
    
    Write-Host ""
    $selection = Read-Host "Select backup number (or press Enter for latest)"
    
    if ([string]::IsNullOrEmpty($selection)) {
        $BackupFile = $backups[0].FullName
        Write-Host "Using latest backup: $($backups[0].Name)" -ForegroundColor Cyan
    } else {
        $index = [int]$selection - 1
        if ($index -ge 0 -and $index -lt $backups.Count) {
            $BackupFile = $backups[$index].FullName
            Write-Host "Using backup: $($backups[$index].Name)" -ForegroundColor Cyan
        } else {
            Write-Host "[X] Invalid selection" -ForegroundColor Red
            exit 1
        }
    }
}

# Verify backup file exists
if (!(Test-Path $BackupFile)) {
    Write-Host "[X] Backup file not found: $BackupFile" -ForegroundColor Red
    exit 1
}

# Show backup info
$backupInfo = Get-Item $BackupFile
$size = if ($backupInfo.Length -gt 1MB) {
    "$([math]::Round($backupInfo.Length / 1MB, 2)) MB"
} elseif ($backupInfo.Length -gt 1KB) {
    "$([math]::Round($backupInfo.Length / 1KB, 2)) KB"
} else {
    "$($backupInfo.Length) bytes"
}

Write-Host ""
Write-Host "Backup Details:" -ForegroundColor Cyan
Write-Host "  File: $($backupInfo.Name)" -ForegroundColor White
Write-Host "  Size: $size" -ForegroundColor White
Write-Host "  Date: $($backupInfo.LastWriteTime)" -ForegroundColor White

# Confirm restore
if (-not $Force) {
    Write-Host ""
    Write-Host "⚠️  WARNING: This will DELETE all current data and restore from backup!" -ForegroundColor Red
    $confirm = Read-Host "Are you sure you want to continue? (yes/no)"
    
    if ($confirm -ne "yes") {
        Write-Host "Restore cancelled." -ForegroundColor Yellow
        exit 0
    }
}

Write-Host ""
Write-Host "Restoring database..." -ForegroundColor Yellow

# Drop existing connections
Write-Host "Terminating existing connections..." -ForegroundColor Cyan
docker exec laptop-tracking-db psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'laptop_tracking_dev' AND pid <> pg_backend_pid();" 2>&1 | Out-Null

# Drop and recreate database
Write-Host "Dropping existing database..." -ForegroundColor Cyan
docker exec laptop-tracking-db psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_dev;" 2>&1 | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "[X] Failed to drop database" -ForegroundColor Red
    exit 1
}

Write-Host "Creating new database..." -ForegroundColor Cyan
docker exec laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;" 2>&1 | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "[X] Failed to create database" -ForegroundColor Red
    exit 1
}

# Restore from backup
Write-Host "Restoring data from backup..." -ForegroundColor Cyan
Get-Content $BackupFile | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev 2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "[OK] Database restored successfully!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    
    # Show statistics
    Write-Host ""
    Write-Host "Verifying restored data..." -ForegroundColor Cyan
    
    $userCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM users;" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Users: $($userCount.Trim())" -ForegroundColor White
    }
    
    $companyCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM client_companies;" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Client Companies: $($companyCount.Trim())" -ForegroundColor White
    }
    
    $laptopCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM laptops;" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Laptops: $($laptopCount.Trim())" -ForegroundColor White
    }
    
    $shipmentCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM shipments;" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Shipments: $($shipmentCount.Trim())" -ForegroundColor White
    }
    
    Write-Host ""
    Write-Host "You may need to restart the application:" -ForegroundColor Cyan
    Write-Host "  docker compose restart app" -ForegroundColor Yellow
    
} else {
    Write-Host ""
    Write-Host "[X] Restore failed!" -ForegroundColor Red
    Write-Host "The database may be in an inconsistent state." -ForegroundColor Yellow
    Write-Host "You may need to run migrations:" -ForegroundColor Cyan
    Write-Host "  docker compose restart app" -ForegroundColor Yellow
    exit 1
}

