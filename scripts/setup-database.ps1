# Database Setup Script for Windows PowerShell
# This script helps set up the PostgreSQL databases for development and testing

Write-Host "=== Laptop Tracking System - Database Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if PostgreSQL is installed
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    Write-Host "ERROR: PostgreSQL 'psql' command not found." -ForegroundColor Red
    Write-Host "Please ensure PostgreSQL is installed and added to PATH." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Installation options:" -ForegroundColor Yellow
    Write-Host "1. Download from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    Write-Host "2. Or use Docker: docker-compose up -d postgres" -ForegroundColor Yellow
    exit 1
}

Write-Host "Found PostgreSQL: $($psqlPath.Source)" -ForegroundColor Green
Write-Host ""

# Get PostgreSQL password
$password = Read-Host "Enter PostgreSQL password for user 'postgres'" -AsSecureString
$BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($password)
$plainPassword = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)

# Set environment variable for this session
$env:PGPASSWORD = $plainPassword

Write-Host ""
Write-Host "Step 1: Creating development database..." -ForegroundColor Cyan

# Create development database
try {
    psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_dev;" 2>$null
    psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;"
    Write-Host "✓ Development database 'laptop_tracking_dev' created successfully" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to create development database" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Step 2: Creating test database..." -ForegroundColor Cyan

# Create test database
try {
    psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;" 2>$null
    psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
    Write-Host "✓ Test database 'laptop_tracking_test' created successfully" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to create test database" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Step 3: Creating .env file..." -ForegroundColor Cyan

# Create .env file if it doesn't exist
if (Test-Path ".env") {
    Write-Host "⚠ .env file already exists, skipping..." -ForegroundColor Yellow
} else {
    if (Test-Path ".env.example") {
        Copy-Item ".env.example" ".env"
        
        # Update password in .env
        $envContent = Get-Content ".env" -Raw
        $envContent = $envContent -replace 'DB_PASSWORD=postgres', "DB_PASSWORD=$plainPassword"
        Set-Content ".env" $envContent
        
        Write-Host "✓ .env file created from .env.example" -ForegroundColor Green
        Write-Host "  Database password has been set" -ForegroundColor Green
    } else {
        Write-Host "✗ .env.example not found, cannot create .env" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Step 4: Running database migrations..." -ForegroundColor Cyan

# Check if migrate is installed
$migratePath = Get-Command migrate -ErrorAction SilentlyContinue
if (-not $migratePath) {
    Write-Host "⚠ 'migrate' command not found, skipping migrations" -ForegroundColor Yellow
    Write-Host "  Install from: https://github.com/golang-migrate/migrate" -ForegroundColor Yellow
    Write-Host "  Or run manually: make migrate-up" -ForegroundColor Yellow
} else {
    $dbUrl = "postgres://postgres:$plainPassword@localhost:5432/laptop_tracking_dev?sslmode=disable"
    try {
        & migrate -path migrations -database $dbUrl up
        Write-Host "✓ Development database migrations completed" -ForegroundColor Green
    } catch {
        Write-Host "✗ Migration failed: $_" -ForegroundColor Red
    }
    
    $testDbUrl = "postgres://postgres:$plainPassword@localhost:5432/laptop_tracking_test?sslmode=disable"
    try {
        & migrate -path migrations -database $testDbUrl up
        Write-Host "✓ Test database migrations completed" -ForegroundColor Green
    } catch {
        Write-Host "✗ Test migration failed: $_" -ForegroundColor Red
    }
}

# Clear password from environment
$env:PGPASSWORD = ""

Write-Host ""
Write-Host "=== Database Setup Complete! ===" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Review .env file and update any necessary settings" -ForegroundColor White
Write-Host "2. Run tests: go test ./... -v" -ForegroundColor White
Write-Host "3. Start the application: go run cmd/web/main.go" -ForegroundColor White
Write-Host ""
Write-Host "To run integration tests:" -ForegroundColor Cyan
Write-Host "  go test ./... -v" -ForegroundColor White
Write-Host ""
Write-Host "To run only unit tests (no database):" -ForegroundColor Cyan
Write-Host "  go test ./... -v -short" -ForegroundColor White
Write-Host ""









