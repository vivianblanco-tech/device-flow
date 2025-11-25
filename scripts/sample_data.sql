-- ============================================
-- COMPREHENSIVE SAMPLE DATA FOR ALIGN
-- ============================================
-- This script populates the database with realistic test data
-- Password for all users: "password123"
-- Last Updated: 2025-11-09

-- ============================================
-- CLEAR EXISTING DATA
-- ============================================
-- Clear in reverse order of dependencies
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
-- Bcrypt hash for "password123": $2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6
INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES
-- Logistics users
('logistics@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'logistics', NOW(), NOW()),
('logistics2@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'logistics', NOW(), NOW()),
('sarah.logistics@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'logistics', NOW(), NOW()),

-- Client users
('client1@techcorp.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),
('client2@innovate.io', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),
('client3@globaltech.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),
('admin@digitaldynamics.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),
('operations@cloudcomputing.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),
('purchasing@retailsolutions.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'client', NOW(), NOW()),

-- Warehouse users
('warehouse@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'warehouse', NOW(), NOW()),
('warehouse2@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'warehouse', NOW(), NOW()),
('michael.warehouse@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'warehouse', NOW(), NOW()),

-- Project Manager users
('pm@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'project_manager', NOW(), NOW()),
('pm2@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'project_manager', NOW(), NOW()),
('jennifer.pm@bairesdev.com', '$2a$10$lmjhwGnRTf7LvgahdoEmteieCMFpTtXzJXmgWrM.hWvhVmm1J9hE6', 'project_manager', NOW(), NOW());

-- ============================================
-- CLIENT COMPANIES
-- ============================================
INSERT INTO client_companies (name, contact_info, created_at) VALUES
('TechCorp International', '{"email":"contact@techcorp.com","phone":"+1-555-0100","address":"100 Tech Plaza, San Francisco, CA 94105"}', NOW()),
('Innovate Solutions', '{"email":"info@innovate.io","phone":"+1-555-0200","address":"200 Innovation Way, Austin, TX 78701"}', NOW()),
('Global Tech Services', '{"email":"support@globaltech.com","phone":"+1-555-0300","address":"300 Global Drive, Seattle, WA 98101"}', NOW()),
('Digital Dynamics', '{"email":"hello@digitaldynamics.com","phone":"+1-555-0400","address":"400 Digital Blvd, Boston, MA 02101"}', NOW()),
('Cloud Computing Corp', '{"email":"contact@cloudcomputing.com","phone":"+1-555-0500","address":"500 Cloud Street, Denver, CO 80202"}', NOW()),
('Retail Solutions Inc', '{"email":"info@retailsolutions.com","phone":"+1-555-0600","address":"600 Retail Row, Chicago, IL 60601"}', NOW()),
('FinTech Partners', '{"email":"hello@fintechpartners.com","phone":"+1-555-0700","address":"700 Finance Ave, New York, NY 10001"}', NOW()),
('Healthcare Systems', '{"email":"support@healthcaresys.com","phone":"+1-555-0800","address":"800 Medical Plaza, Houston, TX 77001"}', NOW());

-- ============================================
-- SOFTWARE ENGINEERS
-- ============================================
INSERT INTO software_engineers (name, email, phone, address, created_at) VALUES
-- Engineers for TechCorp
('Alice Johnson', 'alice.johnson@techcorp.com', '+1-555-1001', '101 Main St, Apt 5B, New York, NY 10001', NOW()),
('Bob Smith', 'bob.smith@techcorp.com', '+1-555-1002', '202 Oak Ave, Unit 12, Los Angeles, CA 90001', NOW()),
('Catherine Wong', 'catherine.wong@techcorp.com', '+1-555-1003', '303 Broadway, Suite 4C, San Francisco, CA 94102', NOW()),

-- Engineers for Innovate Solutions
('Carol Martinez', 'carol.martinez@innovate.io', '+1-555-2001', '303 Pine Rd, Suite 7, Chicago, IL 60601', NOW()),
('David Chen', 'david.chen@innovate.io', '+1-555-2002', '404 Elm St, Apt 3A, Houston, TX 77001', NOW()),
('Emily Rodriguez', 'emily.rodriguez@innovate.io', '+1-555-2003', '505 River Rd, Building 2, Austin, TX 78702', NOW()),

-- Engineers for Global Tech
('Emma Wilson', 'emma.wilson@globaltech.com', '+1-555-3001', '505 Maple Dr, Miami, FL 33101', NOW()),
('Frank Rodriguez', 'frank.rodriguez@globaltech.com', '+1-555-3002', '606 Birch Ln, Phoenix, AZ 85001', NOW()),
('George Kim', 'george.kim@globaltech.com', '+1-555-3003', '707 Pacific Ave, Seattle, WA 98101', NOW()),

-- Engineers for Digital Dynamics
('Grace Lee', 'grace.lee@digitaldynamics.com', '+1-555-4001', '707 Cedar Ave, Philadelphia, PA 19101', NOW()),
('Henry Patel', 'henry.patel@digitaldynamics.com', '+1-555-4002', '808 Walnut St, San Antonio, TX 78201', NOW()),
('Isabella Santos', 'isabella.santos@digitaldynamics.com', '+1-555-4003', '909 Commonwealth Ave, Boston, MA 02101', NOW()),

-- Engineers for Cloud Computing
('Isabel Garcia', 'isabel.garcia@cloudcomputing.com', '+1-555-5001', '909 Spruce Way, San Diego, CA 92101', NOW()),
('Jack Thompson', 'jack.thompson@cloudcomputing.com', '+1-555-5002', '1010 Ash Blvd, Dallas, TX 75201', NOW()),
('Karen Nguyen', 'karen.nguyen@cloudcomputing.com', '+1-555-5003', '1111 Mountain View, Denver, CO 80202', NOW()),

-- Engineers for Retail Solutions
('Liam O''Connor', 'liam.oconnor@retailsolutions.com', '+1-555-6001', '1212 State St, Milwaukee, WI 53201', NOW()),
('Mia Anderson', 'mia.anderson@retailsolutions.com', '+1-555-6002', '1313 Lake Shore Dr, Chicago, IL 60601', NOW()),

-- Engineers for FinTech Partners
('Nathan Brown', 'nathan.brown@fintechpartners.com', '+1-555-7001', '1414 Wall St, New York, NY 10005', NOW()),
('Olivia Martinez', 'olivia.martinez@fintechpartners.com', '+1-555-7002', '1515 Broadway, New York, NY 10036', NOW()),

-- Engineers for Healthcare Systems
('Peter Johnson', 'peter.johnson@healthcaresys.com', '+1-555-8001', '1616 Medical Center Dr, Houston, TX 77030', NOW()),
('Quinn Davis', 'quinn.davis@healthcaresys.com', '+1-555-8002', '1717 Hospital Blvd, Houston, TX 77054', NOW());

-- ============================================
-- LAPTOPS
-- ============================================
INSERT INTO laptops (serial_number, brand, model, specs, status, created_at) VALUES
-- Dell laptops
('DELL-XPS-001', 'Dell', 'XPS 15 9520', 'CPU: Intel i7-12700H, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED, GPU: NVIDIA RTX 3050 Ti', 'available', NOW()),
('DELL-XPS-002', 'Dell', 'XPS 15 9520', 'CPU: Intel i7-12700H, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED, GPU: NVIDIA RTX 3050 Ti', 'in_transit_to_engineer', NOW()),
('DELL-XPS-003', 'Dell', 'XPS 13 Plus', 'CPU: Intel i5-1240P, RAM: 16GB LPDDR5, Storage: 512GB NVMe SSD, Display: 13.4" FHD+', 'available', NOW()),
('DELL-LAT-004', 'Dell', 'Latitude 9430', 'CPU: Intel i7-1265U, RAM: 16GB DDR5, Storage: 512GB NVMe SSD, Display: 14" FHD', 'delivered', NOW()),
('DELL-LAT-005', 'Dell', 'Latitude 7430', 'CPU: Intel i5-1245U, RAM: 16GB DDR4, Storage: 256GB NVMe SSD, Display: 14" FHD', 'in_transit_to_warehouse', NOW()),
('DELL-PRE-006', 'Dell', 'Precision 5570', 'CPU: Intel i7-12800H, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" FHD+, GPU: NVIDIA RTX A2000', 'available', NOW()),

-- HP laptops
('HP-ELITE-001', 'HP', 'EliteBook 850 G9', 'CPU: Intel i7-1255U, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" FHD', 'available', NOW()),
('HP-ZBOOK-002', 'HP', 'ZBook Studio G9', 'CPU: Intel i7-12800H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 15.6" 4K DreamColor, GPU: NVIDIA RTX A2000', 'in_transit_to_warehouse', NOW()),
('HP-ELITE-003', 'HP', 'EliteBook 840 G9', 'CPU: Intel i5-1235U, RAM: 16GB DDR4, Storage: 512GB NVMe SSD, Display: 14" FHD', 'delivered', NOW()),
('HP-ELITE-004', 'HP', 'EliteBook 860 G9', 'CPU: Intel i7-1265U, RAM: 32GB DDR5, Storage: 512GB NVMe SSD, Display: 16" FHD+', 'at_warehouse', NOW()),
('HP-ZBOOK-005', 'HP', 'ZBook Firefly G9', 'CPU: Intel i7-1255U, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 14" FHD, GPU: NVIDIA T550', 'available', NOW()),

-- Lenovo laptops
('LENOVO-X1-001', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WUXGA', 'available', NOW()),
('LENOVO-X1-002', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WUXGA', 'at_warehouse', NOW()),
('LENOVO-P1-003', 'Lenovo', 'ThinkPad P1 Gen 5', 'CPU: Intel i9-12900H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 16" 4K OLED, GPU: NVIDIA RTX A5500', 'delivered', NOW()),
('LENOVO-X1-004', 'Lenovo', 'ThinkPad X1 Yoga Gen 7', 'CPU: Intel i7-1260P, RAM: 16GB LPDDR5, Storage: 512GB NVMe SSD, Display: 14" FHD+ Touch', 'available', NOW()),
('LENOVO-T14-005', 'Lenovo', 'ThinkPad T14 Gen 3', 'CPU: Intel i5-1240P, RAM: 16GB DDR4, Storage: 512GB NVMe SSD, Display: 14" FHD', 'in_transit_to_engineer', NOW()),

-- Apple MacBook laptops
('APPLE-MBP-001', 'Apple', 'MacBook Pro 16" M2 Max', 'CPU: Apple M2 Max, RAM: 32GB Unified, Storage: 1TB SSD, Display: 16.2" Liquid Retina XDR', 'available', NOW()),
('APPLE-MBP-002', 'Apple', 'MacBook Pro 14" M2 Pro', 'CPU: Apple M2 Pro, RAM: 16GB Unified, Storage: 512GB SSD, Display: 14.2" Liquid Retina XDR', 'in_transit_to_warehouse', NOW()),
('APPLE-MBA-003', 'Apple', 'MacBook Air M2', 'CPU: Apple M2, RAM: 16GB Unified, Storage: 512GB SSD, Display: 13.6" Liquid Retina', 'delivered', NOW()),
('APPLE-MBP-004', 'Apple', 'MacBook Pro 16" M2 Pro', 'CPU: Apple M2 Pro, RAM: 32GB Unified, Storage: 1TB SSD, Display: 16.2" Liquid Retina XDR', 'at_warehouse', NOW()),
('APPLE-MBA-005', 'Apple', 'MacBook Air M2', 'CPU: Apple M2, RAM: 24GB Unified, Storage: 1TB SSD, Display: 13.6" Liquid Retina', 'available', NOW()),

-- Microsoft Surface laptops
('SURFACE-LAP-001', 'Microsoft', 'Surface Laptop 5', 'CPU: Intel i7-1255U, RAM: 32GB LPDDR5x, Storage: 1TB SSD, Display: 15" PixelSense', 'available', NOW()),
('SURFACE-STU-002', 'Microsoft', 'Surface Laptop Studio', 'CPU: Intel i7-11370H, RAM: 32GB DDR4, Storage: 1TB SSD, Display: 14.4" PixelSense Flow, GPU: NVIDIA RTX A2000', 'at_warehouse', NOW()),
('SURFACE-LAP-003', 'Microsoft', 'Surface Laptop 5', 'CPU: Intel i5-1235U, RAM: 16GB LPDDR5x, Storage: 512GB SSD, Display: 13.5" PixelSense', 'in_transit_to_engineer', NOW()),

-- ASUS laptops
('ASUS-ZEN-001', 'ASUS', 'ZenBook Pro 15 OLED', 'CPU: Intel i9-12900H, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED, GPU: NVIDIA RTX 3060', 'available', NOW()),
('ASUS-ROG-002', 'ASUS', 'ROG Zephyrus G14', 'CPU: AMD Ryzen 9 6900HS, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 14" QHD+, GPU: AMD Radeon RX 6800S', 'available', NOW()),

-- Acer laptops
('ACER-PRE-001', 'Acer', 'Predator Helios 300', 'CPU: Intel i7-12700H, RAM: 16GB DDR5, Storage: 512GB NVMe SSD, Display: 15.6" FHD 144Hz, GPU: NVIDIA RTX 3060', 'available', NOW()),
('ACER-SWIFT-002', 'Acer', 'Swift X', 'CPU: AMD Ryzen 7 5800U, RAM: 16GB LPDDR4x, Storage: 512GB NVMe SSD, Display: 14" FHD, GPU: NVIDIA RTX 3050 Ti', 'available', NOW());

-- ============================================
-- SHIPMENTS
-- ============================================

-- Shipment 1: Pending Pickup (recently created, pickup scheduled in future)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (1, NULL, 'pending_pickup_from_client', 'SCOP-70001', (NOW() + INTERVAL '5 days')::date, 'New project kickoff - Dell XPS laptops requested', NOW() - INTERVAL '1 day', NOW());

-- Shipment 2: Pending Pickup (urgent)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (2, NULL, 'pending_pickup_from_client', 'SCOP-70002', (NOW() + INTERVAL '2 days')::date, 'Urgent: Replacement laptop needed ASAP', NOW() - INTERVAL '6 hours', NOW());

-- Shipment 3: Pickup Scheduled (client submitted form, pickup date confirmed)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (3, NULL, 'pickup_from_client_scheduled', 'SCOP-70003', (NOW() + INTERVAL '3 days')::date, 'Bulk shipment - multiple laptops for new hires', NOW() - INTERVAL '2 days', NOW());

-- Shipment 4: Pickup Scheduled (tomorrow)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (4, NULL, 'pickup_from_client_scheduled', 'SCOP-70004', (NOW() + INTERVAL '1 day')::date, 'High-end workstation for data science team', NOW() - INTERVAL '3 days', NOW());

-- Shipment 5: Picked up from client (on its way to warehouse)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (5, NULL, 'picked_up_from_client', 'SCOP-70005', 'FedEx Express', 'FX1234567890', (NOW() - INTERVAL '1 day')::date, NOW() - INTERVAL '8 hours', 'Standard pickup completed', NOW() - INTERVAL '4 days', NOW());

-- Shipment 6: Picked up from client
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (6, NULL, 'picked_up_from_client', 'SCOP-70006', 'UPS Next Day Air', 'UPS2345678901', (NOW() - INTERVAL '2 days')::date, NOW() - INTERVAL '1 day', 'Priority shipment', NOW() - INTERVAL '5 days', NOW());

-- Shipment 7: In transit to warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (7, NULL, 'in_transit_to_warehouse', 'SCOP-70007', 'DHL Express', 'DHL3456789012', (NOW() - INTERVAL '3 days')::date, NOW() - INTERVAL '2 days', 'International shipment from overseas office', NOW() - INTERVAL '6 days', NOW());

-- Shipment 8: In transit to warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (1, NULL, 'in_transit_to_warehouse', 'SCOP-70008', 'FedEx Ground', 'FX4567890123', (NOW() - INTERVAL '4 days')::date, NOW() - INTERVAL '3 days', 'Multiple MacBooks for iOS development team', NOW() - INTERVAL '7 days', NOW());

-- Shipment 9: At warehouse (ready for assignment)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (2, NULL, 'at_warehouse', 'SCOP-70009', 'UPS Ground', 'UPS5678901234', (NOW() - INTERVAL '6 days')::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 'Awaiting engineer assignment', NOW() - INTERVAL '8 days', NOW());

-- Shipment 10: At warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (3, NULL, 'at_warehouse', 'SCOP-70010', 'FedEx Express', 'FX6789012345', (NOW() - INTERVAL '7 days')::date, NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days', 'HP ZBook for video editing', NOW() - INTERVAL '9 days', NOW());

-- Shipment 11: At warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (4, NULL, 'at_warehouse', 'SCOP-70011', 'DHL Express', 'DHL7890123456', (NOW() - INTERVAL '8 days')::date, NOW() - INTERVAL '7 days', NOW() - INTERVAL '4 days', 'ThinkPad for backend developer', NOW() - INTERVAL '10 days', NOW());

-- Shipment 12: Released from warehouse (assigned to engineer, ready to ship)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (5, 13, 'released_from_warehouse', 'SCOP-70012', 'FedEx Ground', 'FX7890123456', (NOW() - INTERVAL '10 days')::date, NOW() - INTERVAL '9 days', NOW() - INTERVAL '6 days', NOW() - INTERVAL '1 day', 'Assigned to Isabel Garcia', NOW() - INTERVAL '12 days', NOW());

-- Shipment 13: Released from warehouse
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (6, 16, 'released_from_warehouse', 'SCOP-70013', 'UPS Ground', 'UPS8901234567', (NOW() - INTERVAL '11 days')::date, NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '2 days', 'Assigned to Liam O''Connor', NOW() - INTERVAL '13 days', NOW());

-- Shipment 14: In transit to engineer
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (7, 17, 'in_transit_to_engineer', 'SCOP-70014', 'FedEx Express', 'FX8901234567', (NOW() - INTERVAL '12 days')::date, NOW() - INTERVAL '11 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '3 days', 'En route to Mia Anderson', NOW() - INTERVAL '14 days', NOW());

-- Shipment 15: In transit to engineer
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (8, 18, 'in_transit_to_engineer', 'SCOP-70015', 'DHL Express', 'DHL9012345678', (NOW() - INTERVAL '13 days')::date, NOW() - INTERVAL '12 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '4 days', 'Delivering to Nathan Brown', NOW() - INTERVAL '15 days', NOW());

-- Shipment 16: In transit to engineer
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (1, 19, 'in_transit_to_engineer', 'SCOP-70016', 'UPS Next Day Air', 'UPS9012345678', (NOW() - INTERVAL '14 days')::date, NOW() - INTERVAL '13 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '5 days', 'Priority delivery to Olivia Martinez', NOW() - INTERVAL '16 days', NOW());

-- Shipment 17: Delivered
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (2, 4, 'delivered', 'SCOP-70017', 'FedEx Express', 'FX9012345678', (NOW() - INTERVAL '20 days')::date, NOW() - INTERVAL '19 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '13 days', NOW() - INTERVAL '10 days', 'Successfully delivered to Carol Martinez', NOW() - INTERVAL '22 days', NOW());

-- Shipment 18: Delivered
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (3, 7, 'delivered', 'SCOP-70018', 'UPS Ground', 'UPS0123456789', (NOW() - INTERVAL '25 days')::date, NOW() - INTERVAL '24 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '15 days', 'Delivered to Emma Wilson', NOW() - INTERVAL '27 days', NOW());

-- Shipment 19: Delivered
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (4, 10, 'delivered', 'SCOP-70019', 'DHL Express', 'DHL1234567890', (NOW() - INTERVAL '30 days')::date, NOW() - INTERVAL '29 days', NOW() - INTERVAL '26 days', NOW() - INTERVAL '23 days', NOW() - INTERVAL '20 days', 'Complete - Grace Lee confirmed receipt', NOW() - INTERVAL '32 days', NOW());

-- Shipment 20: Delivered (bulk shipment)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (5, 14, 'delivered', 'SCOP-70020', 'FedEx Ground', 'FX0123456789', (NOW() - INTERVAL '35 days')::date, NOW() - INTERVAL '34 days', NOW() - INTERVAL '31 days', NOW() - INTERVAL '28 days', NOW() - INTERVAL '25 days', 'Bulk delivery completed - Jack Thompson', NOW() - INTERVAL '37 days', NOW());

-- ============================================
-- SHIPMENT LAPTOPS JUNCTION
-- ============================================

-- Shipment 1: 2 Dell XPS
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (1, 1), (1, 2);

-- Shipment 2: 1 HP Elite
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (2, 7);

-- Shipment 3: 3 Lenovo ThinkPads (bulk)
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (3, 12), (3, 14), (3, 15);

-- Shipment 4: 1 HP ZBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (4, 8);

-- Shipment 5: 1 MacBook Pro
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (5, 17);

-- Shipment 6: 1 Dell Precision
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (6, 6);

-- Shipment 7: 2 Apple MacBooks
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (7, 20), (7, 21);

-- Shipment 8: 2 Apple MacBooks
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (8, 16), (8, 18);

-- Shipment 9: 1 Surface Laptop Studio
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (9, 22);

-- Shipment 10: 1 HP ZBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (10, 8);

-- Shipment 11: 1 Lenovo X1 Carbon
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (11, 13);

-- Shipment 12: 1 HP Elite
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (12, 10);

-- Shipment 13: 1 ASUS ZenBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (13, 24);

-- Shipment 14: 1 Dell XPS
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (14, 3);

-- Shipment 15: 1 Surface Laptop
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (15, 23);

-- Shipment 16: 1 Lenovo ThinkPad
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (16, 14);

-- Shipment 17: 1 Dell Latitude
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (17, 4);

-- Shipment 18: 1 HP EliteBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (18, 9);

-- Shipment 19: 1 Lenovo P1
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (19, 14);

-- Shipment 20: 2 Acer laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (20, 26), (20, 27);

-- ============================================
-- PICKUP FORMS
-- ============================================
-- Every shipment must have a pickup form

-- Pickup form for Shipment 1 (Pending Pickup)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(1, 4, NOW() - INTERVAL '1 day', 
 '{"contact_name":"Sarah Johnson","contact_email":"sarah.johnson@techcorp.com","contact_phone":"+1-555-0101","pickup_address":"100 Tech Plaza, Suite 1200","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"2025-12-15","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"2x USB-C charging cables, 2x laptop bags, 2x wireless mice","special_instructions":"Please call 30 minutes before arrival. Building requires security check-in at main lobby. Ask for Sarah Johnson at reception."}'::jsonb);

-- Pickup form for Shipment 2 (Pending Pickup - Urgent)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(2, 5, NOW() - INTERVAL '6 hours',
 '{"contact_name":"Michael Brown","contact_email":"michael.brown@innovate.io","contact_phone":"+1-555-0202","pickup_address":"200 Innovation Way, Building B, Floor 3","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"2025-12-12","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"USB-C dock, external keyboard, HDMI cable","special_instructions":"URGENT: Engineer laptop failed. Replacement needed ASAP. Contact Michael directly upon arrival."}'::jsonb);

-- Pickup form for Shipment 3 (Pickup Scheduled - Bulk)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(3, 6, NOW() - INTERVAL '2 days',
 '{"contact_name":"Jennifer Davis","contact_email":"jennifer.davis@globaltech.com","contact_phone":"+1-555-0303","pickup_address":"300 Global Drive, Floor 5, IT Department","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"2025-12-13","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":18.0,"bulk_height":10.0,"bulk_weight":35.5,"include_accessories":true,"accessories_description":"Multiple charging cables, 3x laptop bags, 3x wireless keyboards and mice sets","special_instructions":"Bulk shipment for new hires. Use loading dock entrance on North side. Forklift assistance available if needed."}'::jsonb);

-- Pickup form for Shipment 4 (Pickup Scheduled)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(4, 7, NOW() - INTERVAL '3 days',
 '{"contact_name":"Robert Chen","contact_email":"robert.chen@digitaldynamics.com","contact_phone":"+1-555-0404","pickup_address":"400 Digital Blvd, Suite 800","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"2025-12-11","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"High-end docking station, 4K monitor cable, ergonomic keyboard and mouse","special_instructions":"High-value equipment. Please handle with extra care. Security escort available in building."}'::jsonb);

-- Pickup form for Shipment 5 (Picked up from client)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(5, 8, NOW() - INTERVAL '4 days',
 '{"contact_name":"Linda Martinez","contact_email":"linda.martinez@cloudcomputing.com","contact_phone":"+1-555-0505","pickup_address":"500 Cloud Street, Building A","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"2025-11-09","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"USB-C charger, laptop sleeve, wireless mouse","special_instructions":"Standard pickup. Receptionist will have package ready at front desk."}'::jsonb);

-- Pickup form for Shipment 6 (Picked up from client)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(6, 9, NOW() - INTERVAL '5 days',
 '{"contact_name":"Thomas Anderson","contact_email":"thomas.anderson@retailsolutions.com","contact_phone":"+1-555-0606","pickup_address":"600 Retail Row, Warehouse 3","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"2025-11-08","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":false,"accessories_description":"","special_instructions":"Warehouse pickup. Enter through loading dock 3. Ring bell for assistance."}'::jsonb);

-- Pickup form for Shipment 7 (In transit to warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(7, 4, NOW() - INTERVAL '6 days',
 '{"contact_name":"Amanda Wilson","contact_email":"amanda.wilson@fintechpartners.com","contact_phone":"+1-555-0707","pickup_address":"700 Finance Ave, 45th Floor","pickup_city":"New York","pickup_state":"NY","pickup_zip":"10001","pickup_date":"2025-11-07","pickup_time_slot":"afternoon","number_of_laptops":2,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"2x USB-C hubs, 2x laptop sleeves, security cables","special_instructions":"High security building. Courier must present ID at security desk. Amanda Wilson will meet in lobby."}'::jsonb);

-- Pickup form for Shipment 8 (In transit to warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(8, 4, NOW() - INTERVAL '7 days',
 '{"contact_name":"David Park","contact_email":"david.park@techcorp.com","contact_phone":"+1-555-0108","pickup_address":"100 Tech Plaza, Suite 2400","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"2025-11-06","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":2,"assignment_type":"single","include_accessories":true,"accessories_description":"2x Apple USB-C cables, 2x Apple Magic Mouse, 2x laptop cases","special_instructions":"MacBooks for iOS development. Fragile - handle with care. Contact David for building access."}'::jsonb);

-- Pickup form for Shipment 9 (At warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(9, 5, NOW() - INTERVAL '8 days',
 '{"contact_name":"Patricia Lopez","contact_email":"patricia.lopez@innovate.io","contact_phone":"+1-555-0209","pickup_address":"200 Innovation Way, Building C","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"2025-11-04","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Docking station, USB-C hub, webcam","special_instructions":"Standard office pickup. Package will be ready at IT department, Building C, 2nd floor."}'::jsonb);

-- Pickup form for Shipment 10 (At warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(10, 6, NOW() - INTERVAL '9 days',
 '{"contact_name":"Richard Kim","contact_email":"richard.kim@globaltech.com","contact_phone":"+1-555-0310","pickup_address":"300 Global Drive, Media Lab","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"2025-11-03","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"High-capacity power adapter, external GPU dock, video editing peripherals","special_instructions":"High-end workstation for video editing. Very heavy equipment. Requires two-person lift."}'::jsonb);

-- Pickup form for Shipment 11 (At warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(11, 7, NOW() - INTERVAL '10 days',
 '{"contact_name":"Susan Taylor","contact_email":"susan.taylor@digitaldynamics.com","contact_phone":"+1-555-0411","pickup_address":"400 Digital Blvd, R&D Wing","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"2025-11-02","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"USB-C dock, wireless keyboard and mouse set, laptop bag","special_instructions":"R&D department pickup. Requires escort through secure area. Contact Susan 15 minutes before arrival."}'::jsonb);

-- Pickup form for Shipment 12 (Released from warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(12, 8, NOW() - INTERVAL '12 days',
 '{"contact_name":"Mark Johnson","contact_email":"mark.johnson@cloudcomputing.com","contact_phone":"+1-555-0512","pickup_address":"500 Cloud Street, Building B, Floor 7","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"2025-10-31","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Laptop charger, mouse, USB hub","special_instructions":"Pickup from IT storage room. Building security will provide access."}'::jsonb);

-- Pickup form for Shipment 13 (Released from warehouse)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(13, 9, NOW() - INTERVAL '13 days',
 '{"contact_name":"Jessica White","contact_email":"jessica.white@retailsolutions.com","contact_phone":"+1-555-0613","pickup_address":"600 Retail Row, Office Tower","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"2025-10-30","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Laptop bag, charger, wireless accessories","special_instructions":"Office tower, 15th floor. Contact Jessica at reception."}'::jsonb);

-- Pickup form for Shipment 14 (In transit to engineer)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(14, 4, NOW() - INTERVAL '14 days',
 '{"contact_name":"Christopher Lee","contact_email":"christopher.lee@fintechpartners.com","contact_phone":"+1-555-0714","pickup_address":"700 Finance Ave, 30th Floor","pickup_city":"New York","pickup_state":"NY","pickup_zip":"10001","pickup_date":"2025-10-29","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":false,"accessories_description":"","special_instructions":"Express pickup. Time-sensitive. Security clearance required - call ahead."}'::jsonb);

-- Pickup form for Shipment 15 (In transit to engineer)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(15, 6, NOW() - INTERVAL '15 days',
 '{"contact_name":"Nancy Rodriguez","contact_email":"nancy.rodriguez@healthcaresys.com","contact_phone":"+1-555-0815","pickup_address":"800 Medical Plaza, Research Wing","pickup_city":"Houston","pickup_state":"TX","pickup_zip":"77001","pickup_date":"2025-10-28","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Medical-grade peripherals, sanitized accessories","special_instructions":"Hospital facility. Follow health protocols. Package will be in sterile packaging."}'::jsonb);

-- Pickup form for Shipment 16 (In transit to engineer)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(16, 4, NOW() - INTERVAL '16 days',
 '{"contact_name":"Daniel Garcia","contact_email":"daniel.garcia@techcorp.com","contact_phone":"+1-555-0116","pickup_address":"100 Tech Plaza, Innovation Lab","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"2025-10-27","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Development kit, specialized cables, laptop accessories","special_instructions":"Priority shipment for VIP engineer. Handle with care."}'::jsonb);

-- Pickup form for Shipment 17 (Delivered)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(17, 5, NOW() - INTERVAL '22 days',
 '{"contact_name":"Elizabeth Brown","contact_email":"elizabeth.brown@innovate.io","contact_phone":"+1-555-0217","pickup_address":"200 Innovation Way, Building A, Suite 500","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"2025-10-21","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Docking station, external monitor cable, wireless peripherals","special_instructions":"Standard office pickup. Package ready at reception desk."}'::jsonb);

-- Pickup form for Shipment 18 (Delivered)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(18, 6, NOW() - INTERVAL '27 days',
 '{"contact_name":"Charles Wilson","contact_email":"charles.wilson@globaltech.com","contact_phone":"+1-555-0318","pickup_address":"300 Global Drive, Engineering Department","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"2025-10-16","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Laptop charger, bag, USB accessories","special_instructions":"Engineering department, 4th floor. Ring doorbell for access."}'::jsonb);

-- Pickup form for Shipment 19 (Delivered)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(19, 7, NOW() - INTERVAL '32 days',
 '{"contact_name":"Barbara Martinez","contact_email":"barbara.martinez@digitaldynamics.com","contact_phone":"+1-555-0419","pickup_address":"400 Digital Blvd, Creative Studio","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"2025-10-11","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Graphics tablet, stylus, laptop accessories","special_instructions":"Creative department. Artistic equipment included. Handle with care."}'::jsonb);

-- Pickup form for Shipment 20 (Delivered - Bulk)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(20, 8, NOW() - INTERVAL '37 days',
 '{"contact_name":"William Thompson","contact_email":"william.thompson@cloudcomputing.com","contact_phone":"+1-555-0520","pickup_address":"500 Cloud Street, Warehouse Complex","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"2025-10-06","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":20.0,"bulk_height":12.0,"bulk_weight":42.0,"include_accessories":true,"accessories_description":"Multiple chargers, bags, wireless peripherals, docking stations","special_instructions":"Bulk shipment. Use loading dock. Forklift required. Heavy equipment."}'::jsonb);

-- ============================================
-- RECEPTION REPORTS
-- ============================================
-- For shipments that reached the warehouse

-- Reception report for Shipment 9 (At warehouse)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(9, 10, NOW() - INTERVAL '2 days', 'Surface Laptop Studio received in excellent condition. All packaging intact. Serial number verified. Device powered on and functioning properly. Ready for assignment.', ARRAY['/uploads/reception/shipment9_photo1.jpg', '/uploads/reception/shipment9_photo2.jpg']);

-- Reception report for Shipment 10 (At warehouse)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(10, 11, NOW() - INTERVAL '3 days', 'HP ZBook Studio G9 arrived safely. High-end workstation verified. All specs match requirements. GPU tested and operational. Packaging secure. No visible damage.', ARRAY['/uploads/reception/shipment10_photo1.jpg', '/uploads/reception/shipment10_photo2.jpg', '/uploads/reception/shipment10_photo3.jpg']);

-- Reception report for Shipment 11 (At warehouse)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(11, 12, NOW() - INTERVAL '4 days', 'Lenovo ThinkPad X1 Carbon received. Condition: Excellent. All accessories present. Battery at 85% health. BIOS updated to latest version. Logged into inventory system.', ARRAY['/uploads/reception/shipment11_photo1.jpg']);

-- Reception report for Shipment 12 (Released from warehouse - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(12, 10, NOW() - INTERVAL '6 days', 'HP EliteBook received and processed. Standard inspection completed. All tests passed. Accessories inventory checked and verified. Released for assignment to engineer.', ARRAY['/uploads/reception/shipment12_photo1.jpg', '/uploads/reception/shipment12_photo2.jpg']);

-- Reception report for Shipment 13 (Released from warehouse - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(13, 11, NOW() - INTERVAL '7 days', 'ASUS ZenBook Pro received. High-end configuration verified. OLED display tested - no dead pixels. Performance benchmarks completed successfully. Accessories complete.', ARRAY['/uploads/reception/shipment13_photo1.jpg']);

-- Reception report for Shipment 14 (In transit to engineer - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(14, 12, NOW() - INTERVAL '8 days', 'Dell XPS 13 Plus received in perfect condition. Ultrabook verified. All ports tested. Display quality excellent. Ready for immediate deployment.', ARRAY['/uploads/reception/shipment14_photo1.jpg']);

-- Reception report for Shipment 15 (In transit to engineer - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(15, 10, NOW() - INTERVAL '9 days', 'Microsoft Surface Laptop 5 received. Premium finish intact. PixelSense display tested. All accessories present including charger and documentation.', ARRAY['/uploads/reception/shipment15_photo1.jpg', '/uploads/reception/shipment15_photo2.jpg']);

-- Reception report for Shipment 16 (In transit to engineer - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(16, 11, NOW() - INTERVAL '10 days', 'Lenovo ThinkPad received with development kit. Special handling completed. All specialized equipment verified. Priority item processed immediately.', ARRAY['/uploads/reception/shipment16_photo1.jpg']);

-- Reception report for Shipment 17 (Delivered - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(17, 12, NOW() - INTERVAL '16 days', 'Dell Latitude 9430 received and inspected. Enterprise-grade device. Security features verified. Docking station tested. All accessories accounted for.', ARRAY['/uploads/reception/shipment17_photo1.jpg', '/uploads/reception/shipment17_photo2.jpg']);

-- Reception report for Shipment 18 (Delivered - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(18, 10, NOW() - INTERVAL '21 days', 'HP EliteBook 840 G9 received in excellent condition. Business-class laptop verified. All components tested. Battery health: 92%. Ready for deployment.', ARRAY['/uploads/reception/shipment18_photo1.jpg']);

-- Reception report for Shipment 19 (Delivered - historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(19, 11, NOW() - INTERVAL '26 days', 'Lenovo ThinkPad P1 Gen 5 received with graphics equipment. Workstation-grade specs verified. NVIDIA RTX A5500 GPU tested. Performance excellent. All peripherals present.', ARRAY['/uploads/reception/shipment19_photo1.jpg', '/uploads/reception/shipment19_photo2.jpg', '/uploads/reception/shipment19_photo3.jpg']);

-- Reception report for Shipment 20 (Delivered - historical bulk)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(20, 12, NOW() - INTERVAL '31 days', 'Bulk shipment received: 2x Acer laptops. Both units inspected and tested. Gaming-grade specs verified on both. All accessories inventoried. Heavy shipment processed with care.', ARRAY['/uploads/reception/shipment20_photo1.jpg', '/uploads/reception/shipment20_photo2.jpg', '/uploads/reception/shipment20_photo3.jpg', '/uploads/reception/shipment20_photo4.jpg']);

-- ============================================
-- DELIVERY FORMS
-- ============================================
-- For shipments that were delivered to engineers

-- Delivery form for Shipment 17 (Delivered)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(17, 4, NOW() - INTERVAL '10 days', 'Dell Latitude 9430 delivered successfully to Carol Martinez. Device unboxed and powered on. Engineer verified all accessories present. Laptop configured and ready for use. Engineer expressed satisfaction with the device specifications.', ARRAY['/uploads/delivery/shipment17_photo1.jpg', '/uploads/delivery/shipment17_photo2.jpg']);

-- Delivery form for Shipment 18 (Delivered)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(18, 7, NOW() - INTERVAL '15 days', 'HP EliteBook 840 G9 delivered to Emma Wilson. Complete setup performed on-site. Engineer tested keyboard, trackpad, and display. All peripherals connected and verified. Docking station configured successfully. Engineer confirmed satisfaction.', ARRAY['/uploads/delivery/shipment18_photo1.jpg', '/uploads/delivery/shipment18_photo2.jpg']);

-- Delivery form for Shipment 19 (Delivered)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(19, 10, NOW() - INTERVAL '20 days', 'Lenovo ThinkPad P1 Gen 5 delivered to Grace Lee. High-end workstation setup completed. Graphics capabilities tested with sample rendering. Engineer confirmed RTX A5500 performance meets requirements. Creative software installed and verified. Tablet and stylus paired successfully.', ARRAY['/uploads/delivery/shipment19_photo1.jpg', '/uploads/delivery/shipment19_photo2.jpg', '/uploads/delivery/shipment19_photo3.jpg']);

-- Delivery form for Shipment 20 (Delivered - Bulk)
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(20, 14, NOW() - INTERVAL '25 days', 'Bulk delivery: 2x Acer laptops delivered to Jack Thompson. Both units unboxed and tested. Gaming-grade performance verified on both systems. All docking stations and accessories connected. Engineer conducted performance tests and confirmed specifications. Multi-monitor setup configured successfully. Engineer approved both devices.', ARRAY['/uploads/delivery/shipment20_photo1.jpg', '/uploads/delivery/shipment20_photo2.jpg', '/uploads/delivery/shipment20_photo3.jpg']);

-- ============================================
-- AUDIT LOGS
-- ============================================

INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details) VALUES
-- Recent activity
(1, 'shipment_created', 'shipment', 1, NOW() - INTERVAL '1 day', '{"action":"shipment_created","jira_ticket_number":"SCOP-70001","client_company_id":1}'),
(4, 'pickup_form_submitted', 'pickup_form', 1, NOW() - INTERVAL '1 day', '{"action":"pickup_form_submitted","shipment_id":1}'),
(1, 'shipment_created', 'shipment', 2, NOW() - INTERVAL '6 hours', '{"action":"shipment_created","jira_ticket_number":"SCOP-70002","client_company_id":2,"priority":"urgent"}'),
(5, 'pickup_form_submitted', 'pickup_form', 2, NOW() - INTERVAL '6 hours', '{"action":"pickup_form_submitted","shipment_id":2}'),

-- Shipment status updates
(1, 'status_updated', 'shipment', 5, NOW() - INTERVAL '8 hours', '{"action":"status_updated","old_status":"pickup_from_client_scheduled","new_status":"picked_up_from_client"}'),
(1, 'status_updated', 'shipment', 6, NOW() - INTERVAL '1 day', '{"action":"status_updated","old_status":"pickup_from_client_scheduled","new_status":"picked_up_from_client"}'),
(1, 'status_updated', 'shipment', 7, NOW() - INTERVAL '2 days', '{"action":"status_updated","old_status":"picked_up_from_client","new_status":"in_transit_to_warehouse"}'),

-- Warehouse activity
(10, 'reception_report_created', 'reception_report', 1, NOW() - INTERVAL '2 days', '{"action":"reception_report_created","shipment_id":9}'),
(11, 'reception_report_created', 'reception_report', 2, NOW() - INTERVAL '3 days', '{"action":"reception_report_created","shipment_id":10}'),
(12, 'reception_report_created', 'reception_report', 3, NOW() - INTERVAL '4 days', '{"action":"reception_report_created","shipment_id":11}'),

-- Engineer assignments
(13, 'engineer_assigned', 'shipment', 12, NOW() - INTERVAL '1 day', '{"action":"engineer_assigned","engineer_id":13,"engineer_name":"Isabel Garcia"}'),
(13, 'engineer_assigned', 'shipment', 13, NOW() - INTERVAL '2 days', '{"action":"engineer_assigned","engineer_id":16,"engineer_name":"Liam O''Connor"}'),

-- Delivery completions
(1, 'status_updated', 'shipment', 17, NOW() - INTERVAL '10 days', '{"action":"status_updated","old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(1, 'status_updated', 'shipment', 18, NOW() - INTERVAL '15 days', '{"action":"status_updated","old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(1, 'status_updated', 'shipment', 19, NOW() - INTERVAL '20 days', '{"action":"status_updated","old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(1, 'status_updated', 'shipment', 20, NOW() - INTERVAL '25 days', '{"action":"status_updated","old_status":"in_transit_to_engineer","new_status":"delivered"}');

-- ============================================
-- SUMMARY
-- ============================================

SELECT 'Comprehensive sample data loaded successfully!' AS message;
SELECT '========================================' AS separator;
SELECT 'DATABASE SUMMARY' AS title;
SELECT '========================================' AS separator;
SELECT COUNT(*) AS total_users, 'Users' AS entity FROM users
UNION ALL
SELECT COUNT(*), 'Client Companies' FROM client_companies
UNION ALL
SELECT COUNT(*), 'Software Engineers' FROM software_engineers
UNION ALL
SELECT COUNT(*), 'Laptops' FROM laptops
UNION ALL
SELECT COUNT(*), 'Shipments' FROM shipments
UNION ALL
SELECT COUNT(*), 'Pickup Forms' FROM pickup_forms
UNION ALL
SELECT COUNT(*), 'Reception Reports' FROM reception_reports
UNION ALL
SELECT COUNT(*), 'Delivery Forms' FROM delivery_forms
UNION ALL
SELECT COUNT(*), 'Audit Logs' FROM audit_logs;

SELECT '========================================' AS separator;
SELECT 'SHIPMENTS BY STATUS' AS title;
SELECT '========================================' AS separator;
SELECT status, COUNT(*) as count FROM shipments GROUP BY status ORDER BY status;

SELECT '========================================' AS separator;
SELECT 'All data loaded and ready for testing!' AS message;
