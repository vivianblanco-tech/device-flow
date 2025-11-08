#!/usr/bin/env pwsh
# Initialize Database with Sample Data if Empty
# This script checks if the database is empty and loads sample data if needed

Write-Host "Checking database status..." -ForegroundColor Cyan

# Wait for database to be ready
$maxAttempts = 30
$attempt = 0
$dbReady = $false

while ($attempt -lt $maxAttempts -and -not $dbReady) {
    $attempt++
    try {
        $result = docker exec laptop-tracking-db pg_isready -U postgres 2>&1
        if ($LASTEXITCODE -eq 0) {
            $dbReady = $true
            Write-Host "[OK] Database is ready" -ForegroundColor Green
        } else {
            Write-Host "Waiting for database... ($attempt/$maxAttempts)" -ForegroundColor Yellow
            Start-Sleep -Seconds 2
        }
    } catch {
        Write-Host "Waiting for database... ($attempt/$maxAttempts)" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
}

if (-not $dbReady) {
    Write-Host "[X] Database is not ready after $maxAttempts attempts" -ForegroundColor Red
    exit 1
}

# Check if users table exists and has data
try {
    $userCount = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM users;" 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[X] Failed to query database" -ForegroundColor Red
        Write-Host $userCount -ForegroundColor Red
        exit 1
    }
    
    $count = $userCount.Trim()
    Write-Host "Current user count: $count" -ForegroundColor Cyan
    
    if ($count -eq "0") {
        Write-Host "`nDatabase is empty. Loading sample data..." -ForegroundColor Yellow
        Write-Host "========================================" -ForegroundColor Yellow
        
        # Load test users
        Write-Host "`nLoading test users..." -ForegroundColor Cyan
        Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "[OK] Test users loaded successfully!" -ForegroundColor Green
        } else {
            Write-Host "[X] Failed to load test users" -ForegroundColor Red
            exit 1
        }
        
        # Ask if user wants to load test data (companies, engineers, laptops, shipments)
        Write-Host "`nDo you want to load additional test data? (companies, engineers, laptops, shipments)" -ForegroundColor Yellow
        Write-Host "This includes:" -ForegroundColor Cyan
        Write-Host "  - 5 Client Companies" -ForegroundColor Cyan
        Write-Host "  - 10 Software Engineers" -ForegroundColor Cyan
        Write-Host "  - 15 Laptops" -ForegroundColor Cyan
        Write-Host "  - 13 Shipments" -ForegroundColor Cyan
        
        $response = Read-Host "`nLoad test data? (y/N)"
        
        if ($response -eq "y" -or $response -eq "Y") {
            Write-Host "`nLoading test data..." -ForegroundColor Cyan
            Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "[OK] Test data loaded successfully!" -ForegroundColor Green
            } else {
                Write-Host "[X] Failed to load test data" -ForegroundColor Red
                exit 1
            }
        } else {
            Write-Host "Skipping test data." -ForegroundColor Yellow
        }
        
        Write-Host "`n========================================" -ForegroundColor Green
        Write-Host "[OK] Database initialization complete!" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Green
        Write-Host "`nTest Users Created:" -ForegroundColor Cyan
        Write-Host "  Email: logistics@bairesdev.com    Password: Test123!  Role: Logistics" -ForegroundColor White
        Write-Host "  Email: warehouse@bairesdev.com    Password: Test123!  Role: Warehouse" -ForegroundColor White
        Write-Host "  Email: client@bairesdev.com       Password: Test123!  Role: Client" -ForegroundColor White
        Write-Host "  Email: pm@bairesdev.com           Password: Test123!  Role: Project Manager" -ForegroundColor White
        Write-Host "`nLogin at: http://localhost:8080/login" -ForegroundColor Cyan
        
    } else {
        Write-Host "[OK] Database already contains data ($count users). Skipping initialization." -ForegroundColor Green
    }
    
} catch {
    Write-Host "[X] Error checking database: $_" -ForegroundColor Red
    exit 1
}

