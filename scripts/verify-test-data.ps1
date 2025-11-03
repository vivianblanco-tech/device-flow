# Verify Test Data Script
# This script verifies that all test data has been properly created in the database

Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "  Laptop Tracking System - Test Data Verification" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Change to project directory
$projectDir = "E:\Cursor Projects\BDH"
Set-Location $projectDir

Write-Host "Checking database connection..." -ForegroundColor Yellow
$checkDb = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT 1;" 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Cannot connect to database. Make sure Docker container is running." -ForegroundColor Red
    Write-Host "Run: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}

Write-Host "Database connection successful" -ForegroundColor Green
Write-Host ""

# Get summary counts
Write-Host "Retrieving data summary..." -ForegroundColor Yellow
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT 'Client Companies' as entity, COUNT(*) as count FROM client_companies UNION ALL SELECT 'Software Engineers', COUNT(*) FROM software_engineers UNION ALL SELECT 'Laptops', COUNT(*) FROM laptops UNION ALL SELECT 'Shipments', COUNT(*) FROM shipments UNION ALL SELECT 'Shipment-Laptop Links', COUNT(*) FROM shipment_laptops;"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Shipment Status Breakdown" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT status, COUNT(*) as count FROM shipments GROUP BY status ORDER BY count DESC;"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Laptop Status Breakdown" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT status, COUNT(*) as count FROM laptops GROUP BY status ORDER BY count DESC;"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Laptops by Brand" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT brand, COUNT(*) as count FROM laptops GROUP BY brand ORDER BY count DESC;"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Sample Shipments with Details" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT s.id, cc.name as client, se.name as engineer, s.status, s.courier_name as courier FROM shipments s JOIN client_companies cc ON s.client_company_id = cc.id LEFT JOIN software_engineers se ON s.software_engineer_id = se.id ORDER BY s.created_at DESC LIMIT 10;"

Write-Host ""
Write-Host "================================================" -ForegroundColor Green
Write-Host "  Verification Complete!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Start the application: go run cmd/web/main.go" -ForegroundColor White
Write-Host "  2. Open browser: http://localhost:8080" -ForegroundColor White
Write-Host "  3. View test data documentation: scripts/TEST_DATA_README.md" -ForegroundColor White
Write-Host ""
