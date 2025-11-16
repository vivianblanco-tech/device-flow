-- =============================================
-- COMPREHENSIVE SHIPMENTS DATA v2.2
-- Complete Shipments, Forms, Reports & Audit Logs
-- =============================================
-- This script adds complete shipment data with all three types:
-- 1. single_full_journey - Client → Warehouse → Engineer
-- 2. bulk_to_warehouse - Multiple laptops → Warehouse
-- 3. warehouse_to_engineer - Warehouse inventory → Engineer
--
-- Run this AFTER loading comprehensive-sample-data-v2.sql
-- Compatible with CPU-required laptop schema (v2.2)
-- =============================================

-- ===========================================
-- SHIPMENT TYPE 1: SINGLE FULL JOURNEY (DELIVERED)
-- ============================================
-- Shipment 1: Completed single laptop delivery
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, 
    released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 1, 1, 'delivered', 1, 'SCOP-90001', 
    'FedEx Express', 'FDX9001234567', 
    (NOW() - INTERVAL '35 days')::date, NOW() - INTERVAL '33 days', NOW() - INTERVAL '30 days', 
    NOW() - INTERVAL '27 days', NOW() - INTERVAL '25 days', 
    'Dell XPS delivered to Alice Johnson. Complete lifecycle. Engineer satisfaction: Excellent.', 
    NOW() - INTERVAL '37 days', NOW() - INTERVAL '25 days');

-- Link laptop to shipment 1
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (1, 12);

-- Pickup form for shipment 1
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(1, 18, NOW() - INTERVAL '37 days',
 ('{"contact_name":"Sarah Mitchell","contact_email":"sarah.mitchell@techcorp.com","contact_phone":"+1-555-0101","pickup_address":"100 Tech Plaza, Building 1, Loading Dock A","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '35 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","special_instructions":"Call 30 minutes before arrival. Security check-in required."}')::jsonb);

-- Laptop-based reception report for shipment 1 (approved)
INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id, 
    received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition, 
    status, approved_by, approved_at, created_at, updated_at) 
VALUES (12, 1, 1, 'FDX9001234567', 7, 
    NOW() - INTERVAL '30 days', 
    'Dell XPS 13 Plus received in excellent condition. Serial number verified. All ports tested functional. Display perfect. Battery health: 100%. Original packaging intact. Ready for assignment.',
    '/uploads/reception/laptop12_serial.jpg',
    '/uploads/reception/laptop12_external.jpg',
    '/uploads/reception/laptop12_working.jpg',
    'approved', 1, NOW() - INTERVAL '29 days',
    NOW() - INTERVAL '30 days', NOW() - INTERVAL '29 days');

-- Update laptop status to delivered
UPDATE laptops SET status = 'delivered', software_engineer_id = 1 WHERE id = 12;

-- Delivery form for shipment 1
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(1, 1, NOW() - INTERVAL '25 days', 
 'Dell XPS 13 Plus delivered to Alice Johnson, New York office. Device powered on successfully. All accessories verified. Engineer tested keyboard, trackpad, display. Setup completed. Engineer confirmed satisfaction.',
 ARRAY['/uploads/delivery/shipment001_photo1.jpg', '/uploads/delivery/shipment001_photo2.jpg']);

-- ============================================
-- SHIPMENT TYPE 2: BULK TO WAREHOUSE (AT WAREHOUSE)
-- ============================================
-- Shipment 2: Bulk shipment of 5 laptops - currently at warehouse, unassigned
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, 
    notes, created_at, updated_at) 
VALUES ('bulk_to_warehouse', 2, NULL, 'at_warehouse', 5, 'SCOP-90002', 
    'UPS Next Day Air', 'UPS9002345678', 
    (NOW() - INTERVAL '7 days')::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 
    'BULK SHIPMENT: 5 ThinkPad X1 Carbon for new hire onboarding. Awaiting engineer assignments from HR department.',
    NOW() - INTERVAL '10 days', NOW() - INTERVAL '2 days');

-- Link 5 laptops to shipment 2
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (2, 24), (2, 25), (2, 26), (2, 27), (2, 28);

-- Pickup form for bulk shipment 2
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(2, 19, NOW() - INTERVAL '10 days',
 ('{"contact_name":"Michael Chen","contact_email":"michael.chen@innovate.io","contact_phone":"+1-555-0201","pickup_address":"200 Innovation Way, Warehouse Building, Bay 5","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '7 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":5,"number_of_boxes":3,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":20.0,"bulk_height":14.0,"bulk_weight":52.5,"special_instructions":"BULK SHIPMENT - New hire equipment. Forklift assistance available. Use loading dock entrance."}')::jsonb);

-- Reception reports for each laptop in bulk shipment (all pending approval)
INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id, 
    received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition, 
    status, created_at, updated_at) 
VALUES 
(24, 2, 2, 'UPS9002345678', 8, NOW() - INTERVAL '2 days', 
 'Laptop 1/5: Lenovo ThinkPad X1 Carbon received. Condition: Excellent. Serial verified. All tests passed.',
 '/uploads/reception/laptop24_serial.jpg', '/uploads/reception/laptop24_external.jpg', '/uploads/reception/laptop24_working.jpg',
 'pending_approval', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(25, 2, 2, 'UPS9002345678', 8, NOW() - INTERVAL '2 days', 
 'Laptop 2/5: Lenovo ThinkPad X1 Carbon received. Condition: Excellent. Serial verified. All tests passed.',
 '/uploads/reception/laptop25_serial.jpg', '/uploads/reception/laptop25_external.jpg', '/uploads/reception/laptop25_working.jpg',
 'pending_approval', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(26, 2, 2, 'UPS9002345678', 8, NOW() - INTERVAL '2 days', 
 'Laptop 3/5: Lenovo ThinkPad X1 Carbon received. Condition: Excellent. Serial verified. All tests passed.',
 '/uploads/reception/laptop26_serial.jpg', '/uploads/reception/laptop26_external.jpg', '/uploads/reception/laptop26_working.jpg',
 'pending_approval', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(27, 2, 2, 'UPS9002345678', 8, NOW() - INTERVAL '2 days', 
 'Laptop 4/5: Lenovo ThinkPad X1 Carbon received. Condition: Excellent. Serial verified. All tests passed.',
 '/uploads/reception/laptop27_serial.jpg', '/uploads/reception/laptop27_external.jpg', '/uploads/reception/laptop27_working.jpg',
 'pending_approval', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(28, 2, 2, 'UPS9002345678', 8, NOW() - INTERVAL '2 days', 
 'Laptop 5/5: Lenovo ThinkPad X1 Carbon received. Condition: Excellent. Serial verified. All tests passed.',
 '/uploads/reception/laptop28_serial.jpg', '/uploads/reception/laptop28_external.jpg', '/uploads/reception/laptop28_working.jpg',
 'pending_approval', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days');

-- Update laptop statuses to at_warehouse
UPDATE laptops SET status = 'at_warehouse', client_company_id = 2 WHERE id IN (24, 25, 26, 27, 28);

-- ============================================
-- SHIPMENT TYPE 3: WAREHOUSE TO ENGINEER (IN TRANSIT)
-- ============================================
-- Shipment 3: Warehouse inventory to engineer (in transit to engineer)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, tracking_number, released_warehouse_at, eta_to_engineer, 
    notes, created_at, updated_at) 
VALUES ('warehouse_to_engineer', 3, 12, 'in_transit_to_engineer', 1, 'SCOP-90003', 
    'FedEx Express', 'FDX9003456789', 
    NOW() - INTERVAL '2 days', (NOW() + INTERVAL '1 day')::timestamp,
    'HP ZBook Fury from warehouse inventory. High-end workstation for Henry Thompson. ETA: Tomorrow.',
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '6 hours');

-- Link laptop from warehouse inventory to shipment 3
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (3, 23);

-- Pickup form for warehouse-to-engineer shipment 3
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(3, 1, NOW() - INTERVAL '3 days',
 ('{"laptop_id":23,"engineer_id":12,"contact_name":"Henry Thompson","contact_email":"henry.thompson@globaltech.com","contact_phone":"+1-555-3001","shipping_address":"808 Pike St, Seattle, WA 98101","jira_ticket_number":"SCOP-90003","special_instructions":"High-end workstation for data science project. Priority delivery. Signature required."}')::jsonb);

-- Update laptop status to in_transit_to_engineer
UPDATE laptops SET status = 'in_transit_to_engineer', software_engineer_id = 12, client_company_id = 3 WHERE id = 23;

-- ============================================
-- MORE SHIPMENTS: SINGLE FULL JOURNEY (IN PROGRESS)
-- ============================================

-- Shipment 4: In transit to warehouse (picked up yesterday)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, tracking_number, pickup_scheduled_date, picked_up_at, 
    notes, created_at, updated_at) 
VALUES ('single_full_journey', 4, 17, 'in_transit_to_warehouse', 1, 'SCOP-90004', 
    'DHL Express', 'DHL9004567890', 
    (NOW() - INTERVAL '2 days')::date, NOW() - INTERVAL '1 day', 
    'Dell Precision workstation for Karen Lee. High-value equipment. Expected warehouse arrival today.',
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '6 hours');

INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (4, 3);
UPDATE laptops SET status = 'in_transit_to_warehouse', client_company_id = 4 WHERE id = 3;

INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(4, 21, NOW() - INTERVAL '5 days',
 ('{"contact_name":"Robert Chen","contact_email":"robert.chen@digitaldynamics.com","contact_phone":"+1-555-0401","pickup_address":"400 Digital Blvd, Suite 1200","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '2 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","special_instructions":"High-end workstation. Handle with extreme care. Security escort available."}')::jsonb);

-- Shipment 5: Pickup scheduled for tomorrow
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, pickup_scheduled_date, 
    notes, created_at, updated_at) 
VALUES ('single_full_journey', 5, 21, 'pickup_from_client_scheduled', 1, 'SCOP-90005', 
    'UPS Ground', 
    (NOW() + INTERVAL '1 day')::date, 
    'Apple MacBook Pro for Nathan Brown. Pickup confirmed for tomorrow 9-12 AM.',
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day');

INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (5, 37);

INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(5, 23, NOW() - INTERVAL '3 days',
 ('{"contact_name":"Linda Martinez","contact_email":"linda.martinez@cloudventures.com","contact_phone":"+1-555-0501","pickup_address":"500 Cloud Street, Building A, Floor 3","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() + INTERVAL '1 day', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","special_instructions":"Package ready at front desk. Visitor parking available in Lot B."}')::jsonb);

-- Shipment 6: Pending pickup (just created)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    pickup_scheduled_date, 
    notes, created_at, updated_at) 
VALUES ('single_full_journey', 6, 26, 'pending_pickup_from_client', 1, 'SCOP-90006', 
    (NOW() + INTERVAL '4 days')::date, 
    'Lenovo ThinkPad P1 for Quinn Anderson. Awaiting pickup form submission.',
    NOW() - INTERVAL '1 day', NOW());

INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (6, 32);

-- ============================================
-- MORE BULK SHIPMENTS
-- ============================================

-- Shipment 7: Bulk shipment delivered (completed)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, 
    released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('bulk_to_warehouse', 7, NULL, 'delivered', 3, 'SCOP-90007', 
    'FedEx Express', 'FDX9007890123', 
    (NOW() - INTERVAL '45 days')::date, NOW() - INTERVAL '43 days', NOW() - INTERVAL '40 days', 
    NOW() - INTERVAL '35 days', NOW() - INTERVAL '33 days', 
    'BULK SHIPMENT COMPLETE: 3 Dell XPS laptops delivered. All engineers confirmed receipt. Project launched successfully.',
    NOW() - INTERVAL '47 days', NOW() - INTERVAL '33 days');

INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (7, 7), (7, 8), (7, 9);

UPDATE laptops SET status = 'delivered', client_company_id = 7 WHERE id IN (7, 8, 9);

INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(7, 25, NOW() - INTERVAL '47 days',
 ('{"contact_name":"Thomas Anderson","contact_email":"thomas.anderson@nextgensw.com","contact_phone":"+1-555-0701","pickup_address":"700 Innovation Court, Creative Studio, Floor 2","pickup_city":"Portland","pickup_state":"OR","pickup_zip":"97201","pickup_date":"' || to_char(NOW() - INTERVAL '45 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":18.0,"bulk_height":12.0,"bulk_weight":38.0,"special_instructions":"BULK - Development team equipment. Contact Thomas upon arrival."}')::jsonb);

-- Reception reports for bulk shipment 7 (all approved - historical)
INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id, 
    received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition, 
    status, approved_by, approved_at, created_at, updated_at) 
VALUES 
(7, 7, 7, 'FDX9007890123', 9, NOW() - INTERVAL '40 days', 
 'Bulk 1/3: Dell XPS 15 9520 received. Excellent condition. All tests passed.',
 '/uploads/reception/laptop7_serial.jpg', '/uploads/reception/laptop7_external.jpg', '/uploads/reception/laptop7_working.jpg',
 'approved', 2, NOW() - INTERVAL '39 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '39 days'),
(8, 7, 7, 'FDX9007890123', 9, NOW() - INTERVAL '40 days', 
 'Bulk 2/3: Dell XPS 15 9520 received. Excellent condition. All tests passed.',
 '/uploads/reception/laptop8_serial.jpg', '/uploads/reception/laptop8_external.jpg', '/uploads/reception/laptop8_working.jpg',
 'approved', 2, NOW() - INTERVAL '39 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '39 days'),
(9, 7, 7, 'FDX9007890123', 9, NOW() - INTERVAL '40 days', 
 'Bulk 3/3: Dell XPS 15 9520 received. Excellent condition. All tests passed.',
 '/uploads/reception/laptop9_serial.jpg', '/uploads/reception/laptop9_external.jpg', '/uploads/reception/laptop9_working.jpg',
 'approved', 2, NOW() - INTERVAL '39 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '39 days');

-- ============================================
-- AUDIT LOGS (Recent Activity)
-- ============================================

INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details) VALUES
-- Recent shipment creation
(1, 'shipment_created', 'shipment', 6, NOW() - INTERVAL '1 day', '{"action":"shipment_created","jira_ticket":"SCOP-90006","status":"pending_pickup_from_client"}'),
(1, 'shipment_created', 'shipment', 5, NOW() - INTERVAL '3 days', '{"action":"shipment_created","jira_ticket":"SCOP-90005","status":"pending_pickup_from_client"}'),
(23, 'pickup_form_submitted', 'pickup_form', 5, NOW() - INTERVAL '3 days', '{"shipment_id":5,"form_type":"single_full_journey"}'),

-- Bulk shipment reception (recent)
(8, 'reception_report_created', 'reception_report', 1, NOW() - INTERVAL '2 days', '{"shipment_id":2,"laptop_id":24,"bulk":true}'),
(8, 'reception_report_created', 'reception_report', 2, NOW() - INTERVAL '2 days', '{"shipment_id":2,"laptop_id":25,"bulk":true}'),
(8, 'reception_report_created', 'reception_report', 3, NOW() - INTERVAL '2 days', '{"shipment_id":2,"laptop_id":26,"bulk":true}'),
(8, 'reception_report_created', 'reception_report', 4, NOW() - INTERVAL '2 days', '{"shipment_id":2,"laptop_id":27,"bulk":true}'),
(8, 'reception_report_created', 'reception_report', 5, NOW() - INTERVAL '2 days', '{"shipment_id":2,"laptop_id":28,"bulk":true}'),
(2, 'status_updated', 'shipment', 2, NOW() - INTERVAL '2 days', '{"old_status":"in_transit_to_warehouse","new_status":"at_warehouse"}'),

-- Warehouse to engineer shipment (recent)
(1, 'shipment_created', 'shipment', 3, NOW() - INTERVAL '3 days', '{"action":"shipment_created","jira_ticket":"SCOP-90003","type":"warehouse_to_engineer"}'),
(1, 'pickup_form_submitted', 'pickup_form', 3, NOW() - INTERVAL '3 days', '{"shipment_id":3,"form_type":"warehouse_to_engineer","laptop_id":23}'),
(1, 'status_updated', 'shipment', 3, NOW() - INTERVAL '2 days', '{"old_status":"released_from_warehouse","new_status":"in_transit_to_engineer"}'),

-- Status updates (yesterday)
(1, 'status_updated', 'shipment', 4, NOW() - INTERVAL '1 day', '{"old_status":"pickup_from_client_scheduled","new_status":"picked_up_from_client"}'),
(1, 'status_updated', 'shipment', 4, NOW() - INTERVAL '6 hours', '{"old_status":"picked_up_from_client","new_status":"in_transit_to_warehouse"}'),

-- Pickup form submission
(21, 'pickup_form_submitted', 'pickup_form', 4, NOW() - INTERVAL '5 days', '{"shipment_id":4,"form_type":"single_full_journey"}'),

-- Delivery completion (historical)
(1, 'status_updated', 'shipment', 1, NOW() - INTERVAL '25 days', '{"old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(1, 'delivery_form_created', 'delivery_form', 1, NOW() - INTERVAL '25 days', '{"shipment_id":1,"engineer_id":1}'),

-- Reception approval (historical)
(1, 'reception_report_approved', 'reception_report', 6, NOW() - INTERVAL '29 days', '{"report_id":6,"laptop_id":12,"shipment_id":1}');

-- ============================================
-- MAGIC LINKS (For delivery forms)
-- ============================================

-- Active magic link for shipment 3 (in transit to engineer)
INSERT INTO magic_links (token, shipment_id, expires_at, used, created_at) VALUES
('abc123def456ghi789jkl012mno345pqr678', 3, NOW() + INTERVAL '7 days', false, NOW() - INTERVAL '2 days');

-- Used magic link for shipment 1 (delivered)
INSERT INTO magic_links (token, shipment_id, expires_at, used, used_at, created_at) VALUES
('xyz789uvw456rst123opq890lmn567hij234', 1, NOW() + INTERVAL '7 days', true, NOW() - INTERVAL '25 days', NOW() - INTERVAL '30 days');

-- ============================================
-- Summary & Verification
-- ============================================

SELECT '========================================' AS separator;
SELECT 'SHIPMENTS DATA LOADED SUCCESSFULLY!' AS message;
SELECT '========================================' AS separator;
SELECT '' AS blank;

SELECT 'SHIPMENTS BY TYPE' AS section;
SELECT '------------------' AS underline;
SELECT 
    shipment_type,
    COUNT(*) as count,
    SUM(laptop_count) as total_laptops
FROM shipments 
GROUP BY shipment_type 
ORDER BY count DESC;

SELECT '' AS blank;
SELECT 'SHIPMENTS BY STATUS' AS section;
SELECT '-------------------' AS underline;
SELECT status, COUNT(*) as count FROM shipments GROUP BY status ORDER BY count DESC;

SELECT '' AS blank;
SELECT 'BULK SHIPMENTS DETAIL' AS section;
SELECT '---------------------' AS underline;
SELECT 
    s.id,
    s.jira_ticket_number,
    s.status,
    s.laptop_count,
    cc.name as client
FROM shipments s
JOIN client_companies cc ON cc.id = s.client_company_id
WHERE s.laptop_count > 1
ORDER BY s.id;

SELECT '' AS blank;
SELECT 'RECEPTION REPORTS STATUS' AS section;
SELECT '------------------------' AS underline;
SELECT status, COUNT(*) as count FROM reception_reports GROUP BY status ORDER BY count DESC;

SELECT '' AS blank;
SELECT '========================================' AS separator;
SELECT 'Complete! All shipment types represented.' AS summary1;
SELECT 'Test the three shipment workflows in the UI.' AS summary2;
SELECT '========================================' AS separator;

