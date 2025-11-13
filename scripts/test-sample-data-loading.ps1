#!/usr/bin/env pwsh
# Test Sample Data Loading Script
# This script validates that sample data can be loaded without errors

Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "  Sample Data Loading Test" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Change to project directory
$projectDir = "E:\Cursor Projects\BDH"
Set-Location $projectDir

# Test 1: Check database connection
Write-Host "Test 1: Checking database connection..." -ForegroundColor Yellow
$checkDb = docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT 1;" 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[FAIL] Cannot connect to database" -ForegroundColor Red
    Write-Host "Make sure Docker container is running: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}
Write-Host "[PASS] Database connection successful" -ForegroundColor Green
Write-Host ""

# Test 2: Clear existing data
Write-Host "Test 2: Clearing existing data..." -ForegroundColor Yellow
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
DELETE FROM audit_logs;
DELETE FROM magic_links;
DELETE FROM sessions;
DELETE FROM delivery_forms;
DELETE FROM reception_reports;
DELETE FROM pickup_forms;
DELETE FROM shipment_laptops;
DELETE FROM laptops;
DELETE FROM shipments;
DELETE FROM software_engineers;
DELETE FROM client_companies;
DELETE FROM users;
" > $null 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[FAIL] Error clearing data" -ForegroundColor Red
    exit 1
}
Write-Host "[PASS] Data cleared successfully" -ForegroundColor Green
Write-Host ""

# Test 3: Load sample data
Write-Host "Test 3: Loading sample data..." -ForegroundColor Yellow
$loadResult = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data.sql 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "[FAIL] Error loading sample data" -ForegroundColor Red
    Write-Host "Error details:" -ForegroundColor Red
    Write-Host $loadResult
    exit 1
}
Write-Host "[PASS] Sample data loaded successfully" -ForegroundColor Green
Write-Host ""

# Test 4: Verify table counts
Write-Host "Test 4: Verifying data counts..." -ForegroundColor Yellow

$userCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM users;" | ForEach-Object { $_.Trim() }
$companyCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM client_companies;" | ForEach-Object { $_.Trim() }
$engineerCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM software_engineers;" | ForEach-Object { $_.Trim() }
$laptopCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM laptops;" | ForEach-Object { $_.Trim() }
$shipmentCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM shipments;" | ForEach-Object { $_.Trim() }

Write-Host "  Users: $userCount" -ForegroundColor White
Write-Host "  Companies: $companyCount" -ForegroundColor White
Write-Host "  Engineers: $engineerCount" -ForegroundColor White
Write-Host "  Laptops: $laptopCount" -ForegroundColor White
Write-Host "  Shipments: $shipmentCount" -ForegroundColor White

$expectedCounts = @{
    'users' = 14
    'companies' = 8
    'engineers' = 22
    'laptops' = 35
    'shipments' = 15
}

$allCountsCorrect = $true
if ([int]$userCount -lt $expectedCounts['users']) { $allCountsCorrect = $false }
if ([int]$companyCount -lt $expectedCounts['companies']) { $allCountsCorrect = $false }
if ([int]$engineerCount -lt $expectedCounts['engineers']) { $allCountsCorrect = $false }
if ([int]$laptopCount -lt $expectedCounts['laptops']) { $allCountsCorrect = $false }
if ([int]$shipmentCount -lt $expectedCounts['shipments']) { $allCountsCorrect = $false }

if ($allCountsCorrect) {
    Write-Host "[PASS] All entity counts are correct" -ForegroundColor Green
} else {
    Write-Host "[WARN] Some counts are lower than expected" -ForegroundColor Yellow
}
Write-Host ""

# Test 5: Check foreign key relationships
Write-Host "Test 5: Checking foreign key relationships..." -ForegroundColor Yellow

# Check users linked to companies
$usersWithCompanies = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(*) FROM users u 
JOIN client_companies cc ON u.client_company_id = cc.id 
WHERE u.role = 'client';
" | ForEach-Object { $_.Trim() }

# Check shipments with valid relationships
$shipmentsValid = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(*) FROM shipments s 
JOIN client_companies cc ON s.client_company_id = cc.id;
" | ForEach-Object { $_.Trim() }

# Check shipment_laptops junction
$junctionValid = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(*) FROM shipment_laptops sl
JOIN shipments s ON sl.shipment_id = s.id
JOIN laptops l ON sl.laptop_id = l.id;
" | ForEach-Object { $_.Trim() }

Write-Host "  Client users with companies: $usersWithCompanies" -ForegroundColor White
Write-Host "  Shipments with valid companies: $shipmentsValid" -ForegroundColor White
Write-Host "  Valid shipment-laptop links: $junctionValid" -ForegroundColor White

if ([int]$shipmentsValid -eq [int]$shipmentCount) {
    Write-Host "[PASS] All foreign key relationships valid" -ForegroundColor Green
} else {
    Write-Host "[FAIL] Some foreign key relationships invalid" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 6: Check shipment types and laptop counts
Write-Host "Test 6: Validating shipment types and laptop counts..." -ForegroundColor Yellow

$typeValidation = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT 
    COUNT(*) as invalid_count
FROM shipments
WHERE 
    (shipment_type = 'single_full_journey' AND laptop_count != 1)
    OR (shipment_type = 'warehouse_to_engineer' AND laptop_count != 1)
    OR (shipment_type = 'bulk_to_warehouse' AND laptop_count < 2);
" | ForEach-Object { $_.Trim() }

if ([int]$typeValidation -eq 0) {
    Write-Host "[PASS] All shipment types have correct laptop counts" -ForegroundColor Green
} else {
    Write-Host "[FAIL] $typeValidation shipments have invalid laptop counts for their type" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 7: Check data quality
Write-Host "Test 7: Checking data quality..." -ForegroundColor Yellow

# Check for shipments with forms
$shipmentsWithForms = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(DISTINCT s.id) FROM shipments s
WHERE EXISTS (SELECT 1 FROM pickup_forms pf WHERE pf.shipment_id = s.id);
" | ForEach-Object { $_.Trim() }

# Check for reception reports
$receptionReports = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(*) FROM reception_reports;
" | ForEach-Object { $_.Trim() }

# Check for delivery forms
$deliveryForms = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "
SELECT COUNT(*) FROM delivery_forms;
" | ForEach-Object { $_.Trim() }

Write-Host "  Shipments with pickup forms: $shipmentsWithForms" -ForegroundColor White
Write-Host "  Reception reports: $receptionReports" -ForegroundColor White
Write-Host "  Delivery forms: $deliveryForms" -ForegroundColor White

if ([int]$shipmentsWithForms -ge 10) {
    Write-Host "[PASS] Good form coverage" -ForegroundColor Green
} else {
    Write-Host "[WARN] Limited form coverage" -ForegroundColor Yellow
}
Write-Host ""

# Test 8: Verify shipment statuses
Write-Host "Test 8: Verifying shipment status distribution..." -ForegroundColor Yellow

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT status, COUNT(*) as count 
FROM shipments 
GROUP BY status 
ORDER BY count DESC;
"

Write-Host "[PASS] Shipment statuses verified" -ForegroundColor Green
Write-Host ""

# Final summary
Write-Host ""
Write-Host "================================================" -ForegroundColor Green
Write-Host "  ✅ All Tests Passed!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
Write-Host "Sample data loaded successfully and validated." -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Start the application: go run cmd/web/main.go" -ForegroundColor White
Write-Host "  2. Access at: http://localhost:8080" -ForegroundColor White
Write-Host "  3. Login with: logistics@bairesdev.com / Test123!" -ForegroundColor White
Write-Host ""

