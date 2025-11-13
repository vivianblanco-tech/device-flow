#!/usr/bin/env pwsh
# Start Application with Automatic Data Loading
# This script starts the Docker containers and ensures sample data exists

param(
    [Parameter(Mandatory=$false)]
    [switch]$Build,
    
    [Parameter(Mandatory=$false)]
    [switch]$Fresh
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Laptop Tracking System - Startup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Handle fresh start
if ($Fresh) {
    Write-Host ""
    Write-Host "⚠️  FRESH START MODE" -ForegroundColor Yellow
    Write-Host "This will DELETE all existing data and volumes!" -ForegroundColor Red
    $confirm = Read-Host "Are you sure? (yes/no)"
    
    if ($confirm -ne "yes") {
        Write-Host "Cancelled." -ForegroundColor Yellow
        exit 0
    }
    
    Write-Host ""
    Write-Host "Stopping containers and removing volumes..." -ForegroundColor Yellow
    docker compose down -v
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[X] Failed to stop containers" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "[OK] Volumes removed" -ForegroundColor Green
    $Build = $true
}

# Start containers
Write-Host ""
if ($Build) {
    Write-Host "Building and starting containers..." -ForegroundColor Cyan
    docker compose up -d --build
} else {
    Write-Host "Starting containers..." -ForegroundColor Cyan
    docker compose up -d
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "[X] Failed to start containers" -ForegroundColor Red
    exit 1
}

Write-Host "[OK] Containers started" -ForegroundColor Green

# Wait a moment for services to initialize
Write-Host ""
Write-Host "Waiting for services to initialize..." -ForegroundColor Cyan
Start-Sleep -Seconds 3

# Initialize database if empty
Write-Host ""
& "$PSScriptRoot\init-db-if-empty.ps1"

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[X] Failed to initialize database" -ForegroundColor Red
    exit 1
}

# Show service status
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "[OK] Application Started Successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Services:" -ForegroundColor Cyan
Write-Host "  Web Application:  http://localhost:8080" -ForegroundColor White
Write-Host "  MailHog (Email):  http://localhost:8025" -ForegroundColor White
Write-Host "  PostgreSQL:       localhost:5432" -ForegroundColor White
Write-Host ""
Write-Host "Test Credentials (Password: Test123!):" -ForegroundColor Cyan
Write-Host "  Logistics:        logistics@bairesdev.com" -ForegroundColor White
Write-Host "  Warehouse:        warehouse@bairesdev.com" -ForegroundColor White
Write-Host "  Project Manager:  pm@bairesdev.com" -ForegroundColor White
Write-Host "  Client Users:     client@techcorp.com, admin@innovate.io" -ForegroundColor White
Write-Host ""
Write-Host "Sample Data Features:" -ForegroundColor Cyan
Write-Host "  * 15 shipments (all statuses and types)" -ForegroundColor Green
Write-Host "  * 8 client companies with contacts" -ForegroundColor Green  
Write-Host "  * 22 software engineers" -ForegroundColor Green
Write-Host "  * 35+ laptops (Dell, HP, Lenovo, Apple, ASUS, Acer)" -ForegroundColor Green
Write-Host "  * Complete pickup, reception and delivery forms" -ForegroundColor Green
Write-Host "  * Audit logs and magic links" -ForegroundColor Green
Write-Host "  * Multiple bulk shipments (2-6 laptops each)" -ForegroundColor Green
Write-Host ""
Write-Host "Data Volume:" -ForegroundColor Cyan
Write-Host "  Users: 14 | Companies: 8 | Engineers: 22" -ForegroundColor White
Write-Host "  Laptops: 35 | Shipments: 15 | Forms: 15 each" -ForegroundColor White
Write-Host ""
Write-Host "Useful Commands:" -ForegroundColor Cyan
Write-Host "  View logs:        docker compose logs -f" -ForegroundColor Yellow
Write-Host "  Stop services:    docker compose down" -ForegroundColor Yellow
Write-Host "  Restart:          docker compose restart" -ForegroundColor Yellow
Write-Host "  Backup DB:        .\scripts\backup-db.ps1" -ForegroundColor Yellow
Write-Host "  Restore DB:       .\scripts\restore-db.ps1" -ForegroundColor Yellow
Write-Host ""

