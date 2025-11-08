-- Sample Data for Laptop Tracking System Development
-- This script populates the database with realistic test data
-- Password for all users: "password123"

-- Clear existing data (in reverse order of dependencies)
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

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE client_companies_id_seq RESTART WITH 1;
ALTER SEQUENCE software_engineers_id_seq RESTART WITH 1;
ALTER SEQUENCE laptops_id_seq RESTART WITH 1;
ALTER SEQUENCE shipments_id_seq RESTART WITH 1;
ALTER SEQUENCE pickup_forms_id_seq RESTART WITH 1;
ALTER SEQUENCE reception_reports_id_seq RESTART WITH 1;
ALTER SEQUENCE delivery_forms_id_seq RESTART WITH 1;

-- ============================================
-- USERS (Password: "password123")
-- ============================================
-- Bcrypt hash for "password123": $2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h
INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES
-- Logistics users
('logistics@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'logistics', NOW(), NOW()),
('logistics2@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'logistics', NOW(), NOW()),

-- Client users
('client1@techcorp.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'client', NOW(), NOW()),
('client2@innovate.io', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'client', NOW(), NOW()),
('client3@globaltech.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'client', NOW(), NOW()),

-- Warehouse users
('warehouse@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'warehouse', NOW(), NOW()),
('warehouse2@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'warehouse', NOW(), NOW()),

-- Project Manager users
('pm@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'project_manager', NOW(), NOW()),
('pm2@bairesdev.com', '$2a$10$rKJ0VqZ0yX4YN3xhVZGUXO7yqK5zB8qGwJZ5FqB0X0XqZ5yX4YN3h', 'project_manager', NOW(), NOW());

-- ============================================
-- CLIENT COMPANIES
-- ============================================
INSERT INTO client_companies (name, contact_info, created_at) VALUES
('TechCorp International', '{"email":"contact@techcorp.com","phone":"+1-555-0100","address":"100 Tech Plaza, San Francisco, CA 94105"}', NOW()),
('Innovate Solutions', '{"email":"info@innovate.io","phone":"+1-555-0200","address":"200 Innovation Way, Austin, TX 78701"}', NOW()),
('Global Tech Services', '{"email":"support@globaltech.com","phone":"+1-555-0300","address":"300 Global Drive, Seattle, WA 98101"}', NOW()),
('Digital Dynamics', '{"email":"hello@digitaldynamics.com","phone":"+1-555-0400","address":"400 Digital Blvd, Boston, MA 02101"}', NOW()),
('Cloud Computing Corp', '{"email":"contact@cloudcomputing.com","phone":"+1-555-0500","address":"500 Cloud Street, Denver, CO 80202"}', NOW());

-- ============================================
-- SOFTWARE ENGINEERS
-- ============================================
INSERT INTO software_engineers (name, email, phone, address, created_at) VALUES
-- Engineers for TechCorp
('Alice Johnson', 'alice.johnson@techcorp.com', '+1-555-1001', '101 Main St, Apt 5B, New York, NY 10001', NOW()),
('Bob Smith', 'bob.smith@techcorp.com', '+1-555-1002', '202 Oak Ave, Unit 12, Los Angeles, CA 90001', NOW()),

-- Engineers for Innovate Solutions
('Carol Martinez', 'carol.martinez@innovate.io', '+1-555-2001', '303 Pine Rd, Suite 7, Chicago, IL 60601', NOW()),
('David Chen', 'david.chen@innovate.io', '+1-555-2002', '404 Elm St, Apt 3A, Houston, TX 77001', NOW()),

-- Engineers for Global Tech
('Emma Wilson', 'emma.wilson@globaltech.com', '+1-555-3001', '505 Maple Dr, Miami, FL 33101', NOW()),
('Frank Rodriguez', 'frank.rodriguez@globaltech.com', '+1-555-3002', '606 Birch Ln, Phoenix, AZ 85001', NOW()),

-- Engineers for Digital Dynamics
('Grace Lee', 'grace.lee@digitaldynamics.com', '+1-555-4001', '707 Cedar Ave, Philadelphia, PA 19101', NOW()),
('Henry Patel', 'henry.patel@digitaldynamics.com', '+1-555-4002', '808 Walnut St, San Antonio, TX 78201', NOW()),

-- Engineers for Cloud Computing
('Isabel Garcia', 'isabel.garcia@cloudcomputing.com', '+1-555-5001', '909 Spruce Way, San Diego, CA 92101', NOW()),
('Jack Thompson', 'jack.thompson@cloudcomputing.com', '+1-555-5002', '1010 Ash Blvd, Dallas, TX 75201', NOW());

-- ============================================
-- LAPTOPS
-- ============================================
INSERT INTO laptops (serial_number, brand, model, specs, status, created_at) VALUES
-- Dell laptops
('DELL-SN-001', 'Dell', 'XPS 15 9520', '{"cpu":"Intel i7-12700H","ram":"32GB DDR5","storage":"1TB NVMe SSD","display":"15.6\" 4K OLED"}', 'available', NOW()),
('DELL-SN-002', 'Dell', 'XPS 15 9520', '{"cpu":"Intel i7-12700H","ram":"32GB DDR5","storage":"1TB NVMe SSD","display":"15.6\" 4K OLED"}', 'in_transit', NOW()),
('DELL-SN-003', 'Dell', 'XPS 13 Plus', '{"cpu":"Intel i5-1240P","ram":"16GB LPDDR5","storage":"512GB NVMe SSD","display":"13.4\" FHD+"}', 'available', NOW()),
('DELL-SN-004', 'Dell', 'Latitude 9430', '{"cpu":"Intel i7-1265U","ram":"16GB DDR5","storage":"512GB NVMe SSD","display":"14\" FHD"}', 'assigned', NOW()),

-- HP laptops
('HP-SN-001', 'HP', 'EliteBook 850 G9', '{"cpu":"Intel i7-1255U","ram":"32GB DDR5","storage":"1TB NVMe SSD","display":"15.6\" FHD"}', 'available', NOW()),
('HP-SN-002', 'HP', 'ZBook Studio G9', '{"cpu":"Intel i7-12800H","ram":"64GB DDR5","storage":"2TB NVMe SSD","display":"15.6\" 4K DreamColor","gpu":"NVIDIA RTX A2000"}', 'in_transit', NOW()),
('HP-SN-003', 'HP', 'EliteBook 840 G9', '{"cpu":"Intel i5-1235U","ram":"16GB DDR4","storage":"512GB NVMe SSD","display":"14\" FHD"}', 'assigned', NOW()),

-- Lenovo laptops
('LENOVO-SN-001', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', '{"cpu":"Intel i7-1260P","ram":"32GB LPDDR5","storage":"1TB NVMe SSD","display":"14\" WUXGA"}', 'available', NOW()),
('LENOVO-SN-002', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', '{"cpu":"Intel i7-1260P","ram":"32GB LPDDR5","storage":"1TB NVMe SSD","display":"14\" WUXGA"}', 'at_warehouse', NOW()),
('LENOVO-SN-003', 'Lenovo', 'ThinkPad P1 Gen 5', '{"cpu":"Intel i9-12900H","ram":"64GB DDR5","storage":"2TB NVMe SSD","display":"16\" 4K OLED","gpu":"NVIDIA RTX A5500"}', 'assigned', NOW()),

-- MacBook laptops
('APPLE-SN-001', 'Apple', 'MacBook Pro 16" M2 Max', '{"cpu":"Apple M2 Max","ram":"32GB Unified","storage":"1TB SSD","display":"16.2\" Liquid Retina XDR"}', 'available', NOW()),
('APPLE-SN-002', 'Apple', 'MacBook Pro 14" M2 Pro', '{"cpu":"Apple M2 Pro","ram":"16GB Unified","storage":"512GB SSD","display":"14.2\" Liquid Retina XDR"}', 'in_transit', NOW()),
('APPLE-SN-003', 'Apple', 'MacBook Air M2', '{"cpu":"Apple M2","ram":"16GB Unified","storage":"512GB SSD","display":"13.6\" Liquid Retina"}', 'assigned', NOW()),

-- Microsoft Surface laptops
('SURFACE-SN-001', 'Microsoft', 'Surface Laptop 5', '{"cpu":"Intel i7-1255U","ram":"32GB LPDDR5x","storage":"1TB SSD","display":"15\" PixelSense"}', 'available', NOW()),
('SURFACE-SN-002', 'Microsoft', 'Surface Laptop Studio', '{"cpu":"Intel i7-11370H","ram":"32GB DDR4","storage":"1TB SSD","display":"14.4\" PixelSense Flow","gpu":"NVIDIA RTX A2000"}', 'at_warehouse', NOW());

-- ============================================
-- SHIPMENTS
-- ============================================

-- Shipment 1: Pending Pickup (with pickup form)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (1, NULL, 'pending_pickup_from_client', 'SCOP-67702', '2024-12-20'::date, 'High priority shipment for new project kickoff', NOW() - INTERVAL '3 days', NOW());

-- Shipment 2: Picked up from client
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (2, NULL, 'picked_up_from_client', 'SCOP-67703', 'FedEx Express', '1234567890', '2024-12-15'::date, NOW() - INTERVAL '2 days', 'Contains specialized equipment', NOW() - INTERVAL '5 days', NOW());

-- Shipment 3: In transit to warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (3, NULL, 'in_transit_to_warehouse', 'SCOP-67704', 'UPS Next Day Air', '1Z9999999999999999', '2024-12-14'::date, NOW() - INTERVAL '3 days', 'Express delivery required', NOW() - INTERVAL '6 days', NOW());

-- Shipment 4: At warehouse (with reception report)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (1, NULL, 'at_warehouse', 'SCOP-67705', 'DHL Express', 'DHLTRACK123456', '2024-12-10'::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 'Awaiting assignment', NOW() - INTERVAL '7 days', NOW());

-- Shipment 5: Released from warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (4, 7, 'released_from_warehouse', 'SCOP-67706', 'FedEx Ground', 'FX987654321', '2024-12-08'::date, NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day', 'Assigned to Grace Lee', NOW() - INTERVAL '10 days', NOW());

-- Shipment 6: In transit to engineer
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (5, 9, 'in_transit_to_engineer', 'SCOP-67707', 'UPS Ground', '1Z8888888888888888', '2024-12-05'::date, NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '3 days', 'Assigned to Isabel Garcia', NOW() - INTERVAL '12 days', NOW());

-- Shipment 7: Delivered (complete with all forms)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (2, 3, 'delivered', 'SCOP-67708', 'FedEx Express', 'FX123456789', '2024-11-25'::date, NOW() - INTERVAL '15 days', NOW() - INTERVAL '13 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days', 'Successfully delivered and confirmed', NOW() - INTERVAL '18 days', NOW());

-- Shipment 8: Another delivered shipment
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (3, 5, 'delivered', 'SCOP-67709', 'DHL Express', 'DHL999888777', '2024-11-20'::date, NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '12 days', 'Bulk shipment completed', NOW() - INTERVAL '22 days', NOW());

-- ============================================
-- SHIPMENT LAPTOPS JUNCTION
-- ============================================

-- Shipment 1 (Pending Pickup) - 2 Dell XPS laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (1, 1), (1, 2);

-- Shipment 2 (Picked Up) - 1 HP ZBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (2, 6);

-- Shipment 3 (In Transit to Warehouse) - 2 Lenovo ThinkPads
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (3, 8), (3, 9);

-- Shipment 4 (At Warehouse) - 1 MacBook Pro
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (4, 11);

-- Shipment 5 (Released) - 1 Dell Latitude
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (5, 4);

-- Shipment 6 (In Transit to Engineer) - 1 Surface Laptop Studio
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (6, 15);

-- Shipment 7 (Delivered) - 1 HP EliteBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (7, 7);

-- Shipment 8 (Delivered) - 1 Lenovo ThinkPad P1
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (8, 10);

-- ============================================
-- PICKUP FORMS
-- ============================================

-- Pickup form for Shipment 1 (Pending Pickup)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(1, 3, NOW() - INTERVAL '3 days', 
 '{"contact_name":"Sarah Johnson","contact_email":"sarah.johnson@techcorp.com","contact_phone":"+1-555-0101","pickup_address":"100 Tech Plaza, Suite 1200","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"2024-12-20","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"2x USB-C charging cables, 2x laptop bags","special_instructions":"Please call 30 minutes before arrival. Building requires security check-in at main lobby."}'::jsonb);

-- Pickup form for Shipment 2 (Picked Up)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(2, 4, NOW() - INTERVAL '5 days',
 '{"contact_name":"Michael Brown","contact_email":"michael.brown@innovate.io","contact_phone":"+1-555-0202","pickup_address":"200 Innovation Way, Building B","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"2024-12-15","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"1x Docking station, 1x External monitor cable, 1x YubiKey","special_instructions":"Park in visitor parking lot C. Contact security desk for building access."}'::jsonb);

-- Pickup form for Shipment 3 (In Transit)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(3, 5, NOW() - INTERVAL '6 days',
 '{"contact_name":"Jennifer Davis","contact_email":"jennifer.davis@globaltech.com","contact_phone":"+1-555-0303","pickup_address":"300 Global Drive, Floor 5","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"2024-12-14","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":1,"assignment_type":"bulk","bulk_length":18.5,"bulk_width":14.0,"bulk_height":8.0,"bulk_weight":22.5,"include_accessories":false,"accessories_description":"","special_instructions":"Express delivery required. Equipment needed for project starting Monday."}'::jsonb);

-- Pickup form for Shipment 7 (Delivered - historical)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(7, 4, NOW() - INTERVAL '18 days',
 '{"contact_name":"Robert Wilson","contact_email":"robert.wilson@innovate.io","contact_phone":"+1-555-0205","pickup_address":"200 Innovation Way, Building A, Suite 300","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"2024-11-25","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"1x USB-C hub, 2x USB-C cables, 1x laptop sleeve","special_instructions":"Standard pickup. Contact Michael Brown upon arrival."}'::jsonb);

-- Pickup form for Shipment 8 (Delivered - bulk)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(8, 5, NOW() - INTERVAL '22 days',
 '{"contact_name":"Amanda Martinez","contact_email":"amanda.martinez@globaltech.com","contact_phone":"+1-555-0304","pickup_address":"300 Global Drive, Warehouse Loading Dock","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"2024-11-20","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":18.0,"bulk_height":12.0,"bulk_weight":35.0,"include_accessories":true,"accessories_description":"Multiple docking stations, cables, and accessories - full inventory list attached","special_instructions":"Large shipment - use loading dock entrance. Forklift available if needed."}'::jsonb);

-- ============================================
-- RECEPTION REPORTS
-- ============================================

-- Reception report for Shipment 4 (At Warehouse)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(4, 6, NOW() - INTERVAL '2 days', 'All items received in good condition. Packaging intact. MacBook Pro verified and logged into inventory.', ARRAY['/uploads/reception/shipment4_photo1.jpg', '/uploads/reception/shipment4_photo2.jpg']);

-- Reception report for Shipment 7 (Delivered - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(7, 6, NOW() - INTERVAL '13 days', 'HP EliteBook received. All accessories present and accounted for. No visible damage.', ARRAY['/uploads/reception/shipment7_photo1.jpg']);

-- Reception report for Shipment 8 (Delivered - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(8, 7, NOW() - INTERVAL '18 days', 'Bulk shipment received. Lenovo ThinkPad P1 in excellent condition. All accessories inventoried.', ARRAY['/uploads/reception/shipment8_photo1.jpg', '/uploads/reception/shipment8_photo2.jpg', '/uploads/reception/shipment8_photo3.jpg']);

-- ============================================
-- DELIVERY FORMS
-- ============================================

-- Delivery form for Shipment 7 (Delivered)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(7, 3, NOW() - INTERVAL '7 days', 'Laptop delivered successfully. Device powered on and verified working. All accessories received and functional. Engineer confirmed satisfaction.', ARRAY['/uploads/delivery/shipment7_photo1.jpg']);

-- Delivery form for Shipment 8 (Delivered)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(8, 5, NOW() - INTERVAL '12 days', 'ThinkPad P1 delivered and set up. Engineer tested graphics capabilities and confirmed all specs match requirements. Docking station configured successfully.', ARRAY['/uploads/delivery/shipment8_photo1.jpg', '/uploads/delivery/shipment8_photo2.jpg']);

-- ============================================
-- AUDIT LOGS
-- ============================================

INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details) VALUES
(1, 'shipment_created', 'shipment', 1, NOW() - INTERVAL '3 days', '{"action":"shipment_created","jira_ticket_number":"SCOP-67702","client_company_id":1}'),
(3, 'pickup_form_submitted', 'pickup_form', 1, NOW() - INTERVAL '3 days', '{"action":"pickup_form_submitted","shipment_id":1}'),
(1, 'status_updated', 'shipment', 2, NOW() - INTERVAL '2 days', '{"action":"status_updated","new_status":"picked_up_from_client"}'),
(6, 'reception_report_created', 'reception_report', 1, NOW() - INTERVAL '2 days', '{"action":"reception_report_created","shipment_id":4}'),
(1, 'engineer_assigned', 'shipment', 5, NOW() - INTERVAL '1 day', '{"action":"engineer_assigned","engineer_id":7}'),
(1, 'status_updated', 'shipment', 7, NOW() - INTERVAL '7 days', '{"action":"status_updated","new_status":"delivered"}');

-- ============================================
-- Summary
-- ============================================

SELECT 'Sample data loaded successfully!' AS message;
SELECT COUNT(*) AS total_users FROM users;
SELECT COUNT(*) AS total_companies FROM client_companies;
SELECT COUNT(*) AS total_engineers FROM software_engineers;
SELECT COUNT(*) AS total_laptops FROM laptops;
SELECT COUNT(*) AS total_shipments FROM shipments;
SELECT COUNT(*) AS total_pickup_forms FROM pickup_forms;
SELECT COUNT(*) AS total_reception_reports FROM reception_reports;
SELECT COUNT(*) AS total_delivery_forms FROM delivery_forms;

