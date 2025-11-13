# =============================================
# Reload Comprehensive Sample Data
# =============================================
# Clears existing shipments and related data, then loads fresh comprehensive data

Write-Host ""
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "  Reloading Comprehensive Sample Data" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host ""

Write-Host "Step 1: Clearing existing shipments and forms..." -ForegroundColor Cyan
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
DELETE FROM audit_logs WHERE entity_type = 'shipment';
DELETE FROM delivery_forms;
DELETE FROM reception_reports;
DELETE FROM pickup_forms;
DELETE FROM shipment_laptops;
DELETE FROM shipments;
ALTER SEQUENCE shipments_id_seq RESTART WITH 1;
ALTER SEQUENCE pickup_forms_id_seq RESTART WITH 1;
ALTER SEQUENCE reception_reports_id_seq RESTART WITH 1;
ALTER SEQUENCE delivery_forms_id_seq RESTART WITH 1;
" | Out-Null

Write-Host "[OK] Cleared existing shipments data" -ForegroundColor Green
Write-Host ""

Write-Host "Step 2: Loading comprehensive shipments..." -ForegroundColor Cyan
Get-Content scripts/enhanced-shipments-comprehensive.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev | Out-Null

Write-Host "[OK] Comprehensive data loaded" -ForegroundColor Green
Write-Host ""

Write-Host "Step 3: Verifying data..." -ForegroundColor Cyan
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "
SELECT '========================================' AS separator
UNION ALL SELECT 'COMPREHENSIVE DATA SUMMARY'
UNION ALL SELECT '========================================'
UNION ALL SELECT ''
UNION ALL SELECT 'Core Entities:'
UNION ALL SELECT '  Users: ' || COUNT(*) FROM users
UNION ALL SELECT '  Companies: ' || COUNT(*) FROM client_companies
UNION ALL SELECT '  Engineers: ' || COUNT(*) FROM software_engineers
UNION ALL SELECT '  Laptops: ' || COUNT(*) FROM laptops
UNION ALL SELECT ''
UNION ALL SELECT 'Shipment Data:'
UNION ALL SELECT '  Total Shipments: ' || COUNT(*) FROM shipments
UNION ALL SELECT '  Pickup Forms: ' || COUNT(*) FROM pickup_forms
UNION ALL SELECT '  Reception Reports: ' || COUNT(*) FROM reception_reports
UNION ALL SELECT '  Delivery Forms: ' || COUNT(*) FROM delivery_forms
UNION ALL SELECT '  Shipment-Laptop Links: ' || COUNT(*) FROM shipment_laptops
UNION ALL SELECT '  Audit Logs: ' || COUNT(*) FROM audit_logs WHERE entity_type = 'shipment'
UNION ALL SELECT ''
UNION ALL SELECT 'Shipments by Status:'
;
SELECT '  ' || status || ': ' || COUNT(*) FROM shipments GROUP BY status ORDER BY 
    CASE status
        WHEN 'delivered' THEN 1
        WHEN 'in_transit_to_engineer' THEN 2
        WHEN 'released_from_warehouse' THEN 3
        WHEN 'at_warehouse' THEN 4
        WHEN 'in_transit_to_warehouse' THEN 5
        WHEN 'picked_up_from_client' THEN 6
        WHEN 'pickup_from_client_scheduled' THEN 7
        WHEN 'pending_pickup_from_client' THEN 8
    END;
SELECT '' AS blank
UNION ALL SELECT 'Shipments by Type:'
;
SELECT '  ' || shipment_type || ': ' || COUNT(*) FROM shipments GROUP BY shipment_type;
SELECT '' AS blank
UNION ALL SELECT '========================================';
" -t

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  Reload Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Application: http://localhost:8080" -ForegroundColor White
Write-Host "MailHog: http://localhost:8025" -ForegroundColor White
Write-Host ""

