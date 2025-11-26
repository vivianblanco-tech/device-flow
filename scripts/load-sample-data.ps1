# PowerShell script to load sample data into the database
# Usage: .\scripts\load-sample-data.ps1

Write-Host "Loading Sample Data into Database..." -ForegroundColor Cyan
Write-Host ""

# Check if .env file exists
if (-Not (Test-Path ".env")) {
    Write-Host "Error: .env file not found!" -ForegroundColor Red
    Write-Host "Please create a .env file with your database configuration." -ForegroundColor Yellow
    exit 1
}

# Load environment variables from .env file
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        $key = $matches[1].Trim()
        $value = $matches[2].Trim()
        # Remove quotes if present
        $value = $value -replace '^["'']|["'']$', ''
        Set-Item -Path "env:$key" -Value $value
    }
}

# Build connection string
$DB_USER = $env:DB_USER
$DB_PASSWORD = $env:DB_PASSWORD
$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_NAME = $env:DB_NAME
$DB_SSLMODE = if ($env:DB_SSLMODE) { $env:DB_SSLMODE } else { "disable" }

$DB_URL = "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

Write-Host "Database: $DB_NAME" -ForegroundColor Yellow
Write-Host "Host: ${DB_HOST}:${DB_PORT}" -ForegroundColor Yellow
Write-Host ""

# Check if Docker is available
$dockerPath = Get-Command docker -ErrorAction SilentlyContinue
if (-Not $dockerPath) {
    Write-Host "Error: docker command not found!" -ForegroundColor Red
    Write-Host "Please ensure Docker is installed and running." -ForegroundColor Yellow
    exit 1
}

# Check if Docker container is running
$containerName = "laptop-tracking-db"
$containerRunning = docker ps --filter "name=$containerName" --format "{{.Names}}" 2>$null
if (-Not $containerRunning) {
    Write-Host "Error: Docker container '$containerName' is not running!" -ForegroundColor Red
    Write-Host "Please start the database container first." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Try running: docker-compose up -d" -ForegroundColor Gray
    exit 1
}

# Load comprehensive sample data (v4.0)
Write-Host "Loading comprehensive sample data (v4.0)..." -ForegroundColor Cyan
Write-Host "Container: $containerName" -ForegroundColor Yellow
Write-Host ""

# Step 1: Load base data (users, companies, engineers, laptops)
Write-Host "[1/2] Loading base data (users, companies, engineers, laptops)..." -ForegroundColor Yellow
$tempFileBase = "/tmp/comprehensive-sample-data-v4.sql"
docker cp "scripts\comprehensive-sample-data-v4.sql" "${containerName}:${tempFileBase}" 2>&1 | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Failed to copy base data SQL file to container!" -ForegroundColor Red
    exit 1
}

$resultBase = docker exec $containerName psql -U $DB_USER -d $DB_NAME -f $tempFileBase 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Failed to load base data!" -ForegroundColor Red
    Write-Host "Error details:" -ForegroundColor Yellow
    Write-Host $resultBase -ForegroundColor Gray
    docker exec $containerName rm $tempFileBase 2>&1 | Out-Null
    exit 1
}

Write-Host "[OK] Base data loaded" -ForegroundColor Green
docker exec $containerName rm $tempFileBase 2>&1 | Out-Null

# Step 2: Load shipments data (shipments, forms, reports, audit logs)
Write-Host "[2/2] Loading shipments data (shipments, forms, reports)..." -ForegroundColor Yellow
$tempFileShipments = "/tmp/comprehensive-shipments-data-v4.sql"
docker cp "scripts\comprehensive-shipments-data-v4.sql" "${containerName}:${tempFileShipments}" 2>&1 | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Failed to copy shipments data SQL file to container!" -ForegroundColor Red
    exit 1
}

$resultShipments = docker exec $containerName psql -U $DB_USER -d $DB_NAME -f $tempFileShipments 2>&1

if ($LASTEXITCODE -eq 0) {
    # Clean up temp file
    docker exec $containerName rm $tempFileShipments 2>&1 | Out-Null
    
    Write-Host "[OK] Shipments data loaded" -ForegroundColor Green
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "[SUCCESS] Comprehensive sample data loaded!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Test Credentials (Password: Test123!):" -ForegroundColor Cyan
    Write-Host "  Logistics:        logistics@bairesdev.com" -ForegroundColor White
    Write-Host "  Warehouse:        warehouse@bairesdev.com" -ForegroundColor White
    Write-Host "  Project Manager:  pm@bairesdev.com" -ForegroundColor White
    Write-Host "  Client:           client@techcorp.com" -ForegroundColor White
    Write-Host ""
    Write-Host "Sample Data Includes:" -ForegroundColor Cyan
    Write-Host "  - 30+ users across all 4 roles" -ForegroundColor White
    Write-Host "  - 15 client companies" -ForegroundColor White
    Write-Host "  - 35+ software engineers (with address confirmations)" -ForegroundColor White
    Write-Host "  - 200+ laptops (Dell, HP, Lenovo, Apple) with AUTO-GENERATED SKUs" -ForegroundColor White
    Write-Host "  - ALL laptops assigned to client companies" -ForegroundColor White
    Write-Host "  - ~100 shipments (all three types: single, bulk, warehouse-to-engineer)" -ForegroundColor White
    Write-Host "  - Average delivery time: ~2.5-2.9 days" -ForegroundColor White
    Write-Host "  - Complete pickup forms with detailed JSON data" -ForegroundColor White
    Write-Host "  - Laptop-based reception reports with approval workflow" -ForegroundColor White
    Write-Host "  - All reception reports have complete photo data" -ForegroundColor White
    Write-Host "  - All shipments have complete field data (courier, tracking, timestamps)" -ForegroundColor White
    Write-Host "  - Delivery forms with photos" -ForegroundColor White
    Write-Host "  - Audit logs tracking all activities" -ForegroundColor White
    Write-Host "  - Magic links for secure delivery confirmation" -ForegroundColor White
    Write-Host ""
    Write-Host "Shipment Types Represented:" -ForegroundColor Cyan
    Write-Host "  + Single Full Journey (Client -> Warehouse -> Engineer)" -ForegroundColor Green
    Write-Host "  + Bulk to Warehouse (2+ laptops -> Warehouse only)" -ForegroundColor Green
    Write-Host "  + Warehouse to Engineer (Inventory -> Engineer)" -ForegroundColor Green
    Write-Host ""
    Write-Host "All Shipment Statuses:" -ForegroundColor Cyan
    Write-Host "  ✓ pending_pickup_from_client" -ForegroundColor Green
    Write-Host "  ✓ pickup_from_client_scheduled" -ForegroundColor Green
    Write-Host "  ✓ picked_up_from_client" -ForegroundColor Green
    Write-Host "  ✓ in_transit_to_warehouse" -ForegroundColor Green
    Write-Host "  ✓ at_warehouse" -ForegroundColor Green
    Write-Host "  ✓ released_from_warehouse" -ForegroundColor Green
    Write-Host "  ✓ in_transit_to_engineer" -ForegroundColor Green
    Write-Host "  ✓ delivered" -ForegroundColor Green
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "[ERROR] Failed to load shipments data!" -ForegroundColor Red
    Write-Host "Error details:" -ForegroundColor Yellow
    Write-Host $resultShipments -ForegroundColor Gray
    
    # Clean up temp file even on error
    docker exec $containerName rm $tempFileShipments 2>&1 | Out-Null
    exit 1
}

