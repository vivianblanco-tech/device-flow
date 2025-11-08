#!/usr/bin/env pwsh
# Backup PostgreSQL Database
# Creates a timestamped backup of the entire database

# Create backup directory
$backupDir = "db-backups"
if (!(Test-Path $backupDir)) {
    Write-Host "Creating backup directory: $backupDir" -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $backupDir | Out-Null
}

# Generate backup filename with timestamp
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$backupFile = "$backupDir/laptop_tracking_backup_$timestamp.sql"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Database Backup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Creating backup..." -ForegroundColor Yellow
Write-Host "File: $backupFile" -ForegroundColor Cyan

# Create backup
docker exec laptop-tracking-db pg_dump -U postgres laptop_tracking_dev > $backupFile

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n[OK] Backup created successfully!" -ForegroundColor Green
    
    # Show file size
    $fileSize = (Get-Item $backupFile).Length
    if ($fileSize -gt 1MB) {
        $sizeStr = "$([math]::Round($fileSize / 1MB, 2)) MB"
    } elseif ($fileSize -gt 1KB) {
        $sizeStr = "$([math]::Round($fileSize / 1KB, 2)) KB"
    } else {
        $sizeStr = "$fileSize bytes"
    }
    
    Write-Host "Size: $sizeStr" -ForegroundColor Cyan
    Write-Host "Location: $backupFile" -ForegroundColor Cyan
    
    # List all backups
    Write-Host "`nAll backups:" -ForegroundColor Cyan
    Get-ChildItem -Path $backupDir -Filter "*.sql" | 
        Sort-Object LastWriteTime -Descending | 
        ForEach-Object {
            $size = if ($_.Length -gt 1MB) {
                "$([math]::Round($_.Length / 1MB, 2)) MB"
            } elseif ($_.Length -gt 1KB) {
                "$([math]::Round($_.Length / 1KB, 2)) KB"
            } else {
                "$($_.Length) bytes"
            }
            Write-Host "  $($_.Name) - $size - $($_.LastWriteTime)" -ForegroundColor White
        }
    
    Write-Host "`n========================================" -ForegroundColor Green
    Write-Host "To restore this backup, run:" -ForegroundColor Cyan
    Write-Host "  .\scripts\restore-db.ps1 -BackupFile '$backupFile'" -ForegroundColor Yellow
    Write-Host "========================================" -ForegroundColor Green
    
} else {
    Write-Host "`n[X] Backup failed!" -ForegroundColor Red
    exit 1
}

