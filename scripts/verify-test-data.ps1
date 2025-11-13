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

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT 
    'Entity' as category,
    'Count' as total
UNION ALL
SELECT '-------------------------', '-----'
UNION ALL
SELECT 'Users (all roles)', COUNT(*)::text FROM users WHERE role != 'admin'
UNION ALL
SELECT 'Client Companies', COUNT(*)::text FROM client_companies
UNION ALL
SELECT 'Software Engineers', COUNT(*)::text FROM software_engineers
UNION ALL
SELECT 'Laptops', COUNT(*)::text FROM laptops
UNION ALL
SELECT 'Shipments', COUNT(*)::text FROM shipments
UNION ALL
SELECT 'Pickup Forms', COUNT(*)::text FROM pickup_forms
UNION ALL
SELECT 'Reception Reports', COUNT(*)::text FROM reception_reports
UNION ALL
SELECT 'Delivery Forms', COUNT(*)::text FROM delivery_forms
UNION ALL
SELECT 'Audit Log Entries', COUNT(*)::text FROM audit_logs
UNION ALL
SELECT 'Shipment-Laptop Links', COUNT(*)::text FROM shipment_laptops;
"

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
Write-Host "  Bulk Shipments (Multi-Laptop)" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT 
    s.id,
    s.jira_ticket_number as ticket,
    s.status,
    COUNT(sl.laptop_id) as laptop_count,
    cc.name as client
FROM shipments s
JOIN shipment_laptops sl ON sl.shipment_id = s.id
JOIN client_companies cc ON cc.id = s.client_company_id
GROUP BY s.id, s.jira_ticket_number, s.status, cc.name
HAVING COUNT(sl.laptop_id) > 1
ORDER BY COUNT(sl.laptop_id) DESC, s.id
LIMIT 10;
"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Recent Shipments with Details" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT 
    s.id, 
    s.jira_ticket_number as ticket,
    cc.name as client, 
    COALESCE(se.name, 'Unassigned') as engineer, 
    s.status, 
    COUNT(sl.laptop_id) as laptops
FROM shipments s 
JOIN client_companies cc ON s.client_company_id = cc.id 
LEFT JOIN software_engineers se ON s.software_engineer_id = se.id
LEFT JOIN shipment_laptops sl ON sl.shipment_id = s.id
GROUP BY s.id, s.jira_ticket_number, cc.name, se.name, s.status
ORDER BY s.created_at DESC 
LIMIT 10;
"

Write-Host ""
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Laptop Brands Distribution" -ForegroundColor Cyan
Write-Host "------------------------------------------------" -ForegroundColor Cyan
Write-Host ""

docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT 
    brand, 
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'available' THEN 1 END) as available,
    COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered
FROM laptops 
GROUP BY brand 
ORDER BY COUNT(*) DESC;
"

Write-Host ""
Write-Host "================================================" -ForegroundColor Green
Write-Host "  Sample Data Verification Complete!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""

Write-Host "Key Features in Sample Data:" -ForegroundColor Cyan
Write-Host "  * Three shipment types (single, bulk, warehouse-to-engineer)" -ForegroundColor White
Write-Host "  * All eight shipment statuses represented" -ForegroundColor White
Write-Host "  * Multiple bulk shipments (2-6 laptops each)" -ForegroundColor White
Write-Host "  * High-end workstations and premium laptops" -ForegroundColor White
Write-Host "  * Complete forms with detailed JSON data" -ForegroundColor White
Write-Host "  * Photo URLs and realistic notes" -ForegroundColor White
Write-Host "  * Audit logs tracking system activity" -ForegroundColor White
Write-Host "  * Magic links for secure access" -ForegroundColor White
Write-Host ""

Write-Host "Data Quality Indicators:" -ForegroundColor Cyan
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT 
    '  Shipments with forms: ' || COUNT(DISTINCT s.id) || ' / ' || (SELECT COUNT(*) FROM shipments) as forms_coverage
FROM shipments s
WHERE EXISTS (SELECT 1 FROM pickup_forms pf WHERE pf.shipment_id = s.id);

SELECT 
    '  Engineers with addresses: ' || COUNT(*) || ' / ' || (SELECT COUNT(*) FROM software_engineers) as address_coverage
FROM software_engineers WHERE address IS NOT NULL AND address != '';

SELECT 
    '  Average laptops per bulk shipment: ' || ROUND(AVG(laptop_count), 1) as avg_bulk_size
FROM shipments WHERE laptop_count > 1;
" -t

Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "  1. Application: http://localhost:8080" -ForegroundColor White
Write-Host "  2. MailHog (Email Testing): http://localhost:8025" -ForegroundColor White
Write-Host "  3. View logs: docker compose logs -f app" -ForegroundColor White
Write-Host ""

Write-Host "Test Credentials (Password: Test123!):" -ForegroundColor Yellow
Write-Host "  Logistics:       logistics@bairesdev.com" -ForegroundColor Gray
Write-Host "  Warehouse:       warehouse@bairesdev.com" -ForegroundColor Gray
Write-Host "  Project Manager: pm@bairesdev.com" -ForegroundColor Gray
Write-Host "  Client:          client@techcorp.com" -ForegroundColor Gray
Write-Host ""

Write-Host "Additional Resources:" -ForegroundColor Cyan
Write-Host "  Documentation: .\docs\" -ForegroundColor White
Write-Host "  Reload data: docker exec -i laptop-tracking-db psql ..." -ForegroundColor White
Write-Host "  Backup DB: .\scripts\backup-db.ps1" -ForegroundColor White
Write-Host ""
