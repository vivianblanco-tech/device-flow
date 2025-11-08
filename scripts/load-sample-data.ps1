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

# Load sample data
Write-Host "Loading sample data from scripts/sample_data.sql..." -ForegroundColor Cyan
Write-Host "Container: $containerName" -ForegroundColor Yellow
Write-Host ""

# Copy SQL file to container and execute it
$tempFile = "/tmp/sample_data.sql"
docker cp "scripts\sample_data.sql" "${containerName}:${tempFile}" 2>&1 | Out-Null

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Failed to copy SQL file to container!" -ForegroundColor Red
    exit 1
}

$result = docker exec $containerName psql -U $DB_USER -d $DB_NAME -f $tempFile 2>&1

if ($LASTEXITCODE -eq 0) {
    # Clean up temp file
    docker exec $containerName rm $tempFile 2>&1 | Out-Null
    
    Write-Host ""
    Write-Host "[SUCCESS] Sample data loaded successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Sample Users (all passwords: 'password123'):" -ForegroundColor Cyan
    Write-Host "  Logistics:        logistics@bairesdev.com" -ForegroundColor White
    Write-Host "  Client:           client1@techcorp.com" -ForegroundColor White
    Write-Host "  Warehouse:        warehouse@bairesdev.com" -ForegroundColor White
    Write-Host "  Project Manager:  pm@bairesdev.com" -ForegroundColor White
    Write-Host ""
    Write-Host "Sample Data Includes:" -ForegroundColor Cyan
    Write-Host "  - 9 users across all roles" -ForegroundColor White
    Write-Host "  - 5 client companies" -ForegroundColor White
    Write-Host "  - 10 software engineers" -ForegroundColor White
    Write-Host "  - 15 laptops (Dell, HP, Lenovo, Apple, Microsoft)" -ForegroundColor White
    Write-Host "  - 8 shipments in various statuses" -ForegroundColor White
    Write-Host "  - 5 pickup forms with detailed information" -ForegroundColor White
    Write-Host "  - 3 reception reports" -ForegroundColor White
    Write-Host "  - 2 delivery forms" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "[ERROR] Failed to load sample data!" -ForegroundColor Red
    Write-Host "Error details:" -ForegroundColor Yellow
    Write-Host $result -ForegroundColor Gray
    
    # Clean up temp file even on error
    docker exec $containerName rm $tempFile 2>&1 | Out-Null
    exit 1
}

