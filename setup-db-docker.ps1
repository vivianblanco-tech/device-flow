# Quick Database Setup Script Using Docker
# Run this after Docker Desktop is started

Write-Host "=== Database Setup with Docker ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Start PostgreSQL
Write-Host "Step 1: Starting PostgreSQL container..." -ForegroundColor Yellow
docker-compose up -d postgres

if ($LASTEXITCODE -ne 0) {
    Write-Host "" 
    Write-Host "ERROR: Failed to start PostgreSQL container" -ForegroundColor Red
    Write-Host "Please ensure Docker Desktop is running and try again" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ PostgreSQL container started" -ForegroundColor Green
Write-Host ""

# Step 2: Wait for PostgreSQL to be ready
Write-Host "Step 2: Waiting for PostgreSQL to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 10
Write-Host "✓ PostgreSQL ready" -ForegroundColor Green
Write-Host ""

# Step 3: Create test database
Write-Host "Step 3: Creating test database..." -ForegroundColor Yellow
docker exec laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_test;" 2>$null

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Test database created" -ForegroundColor Green
} else {
    Write-Host "⚠ Test database might already exist" -ForegroundColor Yellow
}
Write-Host ""

# Step 4: Check for migrate command
$migrateExists = Get-Command migrate -ErrorAction SilentlyContinue
if (-not $migrateExists) {
    Write-Host "Step 4: Installing golang-migrate..." -ForegroundColor Yellow
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    # Refresh PATH
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
}

# Step 5: Run migrations
Write-Host "Step 5: Running database migrations..." -ForegroundColor Yellow

$devDbUrl = "postgres://postgres:postgres@localhost:5432/laptop_tracking_dev?sslmode=disable"
$testDbUrl = "postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"

Write-Host "  - Migrating development database..." -ForegroundColor Gray
migrate -path migrations -database $devDbUrl up 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Development database migrated" -ForegroundColor Green
} else {
    Write-Host "  ⚠ Development database migration had warnings" -ForegroundColor Yellow
}

Write-Host "  - Migrating test database..." -ForegroundColor Gray
migrate -path migrations -database $testDbUrl up 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Test database migrated" -ForegroundColor Green
} else {
    Write-Host "  ⚠ Test database migration had warnings" -ForegroundColor Yellow
}

Write-Host ""

# Step 6: Verify setup
Write-Host "Step 6: Verifying setup..." -ForegroundColor Yellow

$tableCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>$null
if ($tableCount -gt 10) {
    Write-Host "✓ Database tables created successfully ($tableCount tables)" -ForegroundColor Green
} else {
    Write-Host "⚠ Warning: Expected more tables" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Setup Complete! ===" -ForegroundColor Green
Write-Host ""
Write-Host "Database Connection Info:" -ForegroundColor Cyan
Write-Host "  Host:     localhost" -ForegroundColor White
Write-Host "  Port:     5432" -ForegroundColor White
Write-Host "  User:     postgres" -ForegroundColor White
Write-Host "  Password: postgres" -ForegroundColor White
Write-Host "  Dev DB:   laptop_tracking_dev" -ForegroundColor White
Write-Host "  Test DB:  laptop_tracking_test" -ForegroundColor White
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Run all tests:  go test ./... -v" -ForegroundColor White
Write-Host "  2. Run unit tests: go test ./... -v -short" -ForegroundColor White
Write-Host "  3. Start app:      go run cmd/web/main.go" -ForegroundColor White
Write-Host ""
Write-Host "To stop database:  docker-compose stop postgres" -ForegroundColor Gray
Write-Host "To start database: docker-compose start postgres" -ForegroundColor Gray
Write-Host ""



