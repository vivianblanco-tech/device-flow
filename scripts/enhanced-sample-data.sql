-- ============================================
-- ENHANCED COMPREHENSIVE SAMPLE DATA
-- Align with Bulk Shipments
-- ============================================
-- This script populates the database with realistic, comprehensive test data
-- Password for all users: "Test123!"
-- Last Updated: 2025-11-10
-- Features: Bulk shipments, all fields populated, realistic accessories

-- ============================================
-- CLEAR EXISTING DATA
-- ============================================
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
ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 100;
ALTER SEQUENCE IF EXISTS client_companies_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS software_engineers_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS laptops_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS shipments_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS pickup_forms_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS reception_reports_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS delivery_forms_id_seq RESTART WITH 1;

-- ============================================
-- USERS (Password: "Test123!")
-- ============================================
-- Bcrypt hash for "Test123!": $2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK
INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES
-- Logistics Team
('logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '6 months', NOW()),
('sarah.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '5 months', NOW()),
('james.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '4 months', NOW()),

-- Warehouse Team
('warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '6 months', NOW()),
('michael.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '5 months', NOW()),
('jessica.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '4 months', NOW()),

-- Project Managers
('pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '6 months', NOW()),
('jennifer.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '5 months', NOW()),
('david.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '4 months', NOW()),

-- Client Users (will be linked to companies)
('client@techcorp.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '6 months', NOW()),
('admin@innovate.io', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '5 months', NOW()),
('purchasing@globaltech.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '5 months', NOW()),
('it-manager@digitaldynamics.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '4 months', NOW()),
('operations@cloudventures.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '4 months', NOW());

-- ============================================
-- CLIENT COMPANIES
-- ============================================
INSERT INTO client_companies (name, contact_info, created_at) VALUES
('TechCorp International', '{"email":"contact@techcorp.com","phone":"+1-555-0100","address":"100 Tech Plaza, San Francisco, CA 94105","country":"USA"}', NOW() - INTERVAL '6 months'),
('Innovate Solutions Ltd', '{"email":"info@innovate.io","phone":"+1-555-0200","address":"200 Innovation Way, Austin, TX 78701","country":"USA"}', NOW() - INTERVAL '6 months'),
('Global Tech Services', '{"email":"support@globaltech.com","phone":"+1-555-0300","address":"300 Global Drive, Seattle, WA 98101","country":"USA"}', NOW() - INTERVAL '5 months'),
('Digital Dynamics Corp', '{"email":"hello@digitaldynamics.com","phone":"+1-555-0400","address":"400 Digital Blvd, Boston, MA 02101","country":"USA"}', NOW() - INTERVAL '5 months'),
('Cloud Ventures Inc', '{"email":"contact@cloudventures.com","phone":"+1-555-0500","address":"500 Cloud Street, Denver, CO 80202","country":"USA"}', NOW() - INTERVAL '4 months'),
('DataDrive Systems', '{"email":"info@datadrive.com","phone":"+1-555-0600","address":"600 Data Lane, Chicago, IL 60601","country":"USA"}', NOW() - INTERVAL '4 months'),
('NextGen Software', '{"email":"hello@nextgensw.com","phone":"+1-555-0700","address":"700 Innovation Court, Portland, OR 97201","country":"USA"}', NOW() - INTERVAL '3 months'),
('Enterprise Solutions Group', '{"email":"contact@enterprisesg.com","phone":"+1-555-0800","address":"800 Enterprise Ave, New York, NY 10001","country":"USA"}', NOW() - INTERVAL '3 months');

-- Link client users to their companies
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'TechCorp International') WHERE email = 'client@techcorp.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Innovate Solutions Ltd') WHERE email = 'admin@innovate.io';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Global Tech Services') WHERE email = 'purchasing@globaltech.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Digital Dynamics Corp') WHERE email = 'it-manager@digitaldynamics.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Cloud Ventures Inc') WHERE email = 'operations@cloudventures.com';

-- ============================================
-- SOFTWARE ENGINEERS
-- ============================================
INSERT INTO software_engineers (name, email, phone, address, created_at) VALUES
-- TechCorp Engineers
('Alice Johnson', 'alice.johnson@techcorp.com', '+1-555-1001', '101 Main St, Apt 5B, New York, NY 10001', NOW() - INTERVAL '5 months'),
('Bob Smith', 'bob.smith@techcorp.com', '+1-555-1002', '202 Oak Ave, Unit 12, Los Angeles, CA 90001', NOW() - INTERVAL '5 months'),
('Catherine Wong', 'catherine.wong@techcorp.com', '+1-555-1003', '303 Broadway, Suite 4C, San Francisco, CA 94102', NOW() - INTERVAL '4 months'),
('Daniel Park', 'daniel.park@techcorp.com', '+1-555-1004', '404 Market St, San Jose, CA 95113', NOW() - INTERVAL '3 months'),

-- Innovate Solutions Engineers
('Emily Rodriguez', 'emily.rodriguez@innovate.io', '+1-555-2001', '505 River Rd, Building 2, Austin, TX 78702', NOW() - INTERVAL '5 months'),
('Frank Martinez', 'frank.martinez@innovate.io', '+1-555-2002', '606 Congress Ave, Austin, TX 78701', NOW() - INTERVAL '4 months'),
('Grace Chen', 'grace.chen@innovate.io', '+1-555-2003', '707 Lamar Blvd, Austin, TX 78703', NOW() - INTERVAL '4 months'),

-- Global Tech Engineers
('Henry Thompson', 'henry.thompson@globaltech.com', '+1-555-3001', '808 Pike St, Seattle, WA 98101', NOW() - INTERVAL '5 months'),
('Isabella Garcia', 'isabella.garcia@globaltech.com', '+1-555-3002', '909 Madison St, Seattle, WA 98104', NOW() - INTERVAL '4 months'),
('James Wilson', 'james.wilson@globaltech.com', '+1-555-3003', '1010 Union St, Seattle, WA 98101', NOW() - INTERVAL '3 months'),

-- Digital Dynamics Engineers
('Karen Lee', 'karen.lee@digitaldynamics.com', '+1-555-4001', '1111 Newbury St, Boston, MA 02116', NOW() - INTERVAL '4 months'),
('Liam O''Connor', 'liam.oconnor@digitaldynamics.com', '+1-555-4002', '1212 Boylston St, Boston, MA 02215', NOW() - INTERVAL '4 months'),
('Maria Santos', 'maria.santos@digitaldynamics.com', '+1-555-4003', '1313 Commonwealth Ave, Boston, MA 02134', NOW() - INTERVAL '3 months'),

-- Cloud Ventures Engineers
('Nathan Brown', 'nathan.brown@cloudventures.com', '+1-555-5001', '1414 16th St, Denver, CO 80202', NOW() - INTERVAL '4 months'),
('Olivia Davis', 'olivia.davis@cloudventures.com', '+1-555-5002', '1515 17th St, Denver, CO 80202', NOW() - INTERVAL '3 months'),
('Patrick Miller', 'patrick.miller@cloudventures.com', '+1-555-5003', '1616 Larimer St, Denver, CO 80202', NOW() - INTERVAL '3 months'),

-- DataDrive Engineers
('Quinn Anderson', 'quinn.anderson@datadrive.com', '+1-555-6001', '1717 Michigan Ave, Chicago, IL 60611', NOW() - INTERVAL '3 months'),
('Rachel White', 'rachel.white@datadrive.com', '+1-555-6002', '1818 State St, Chicago, IL 60605', NOW() - INTERVAL '3 months'),

-- NextGen Engineers
('Samuel Taylor', 'samuel.taylor@nextgensw.com', '+1-555-7001', '1919 SW Broadway, Portland, OR 97201', NOW() - INTERVAL '2 months'),
('Tiffany Clark', 'tiffany.clark@nextgensw.com', '+1-555-7002', '2020 NW Lovejoy St, Portland, OR 97209', NOW() - INTERVAL '2 months'),

-- Enterprise Solutions Engineers
('Victor Harris', 'victor.harris@enterprisesg.com', '+1-555-8001', '2121 Broadway, New York, NY 10023', NOW() - INTERVAL '2 months'),
('Wendy Martinez', 'wendy.martinez@enterprisesg.com', '+1-555-8002', '2222 Park Ave, New York, NY 10037', NOW() - INTERVAL '2 months');

-- ============================================
-- LAPTOPS (Comprehensive Inventory)
-- ============================================
INSERT INTO laptops (serial_number, sku, brand, model, specs, status, created_at) VALUES
-- High-End Workstations (Dell Precision)
('DELL-PREC-5570-001', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'CPU: Intel i9-12900H 14-core, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 15.6" UHD+ (3840x2400), GPU: NVIDIA RTX A2000 8GB, Ports: 4x TB4, WiFi 6E', 'available', NOW() - INTERVAL '3 months'),
('DELL-PREC-5570-002', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'CPU: Intel i9-12900H 14-core, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 15.6" UHD+ (3840x2400), GPU: NVIDIA RTX A2000 8GB, Ports: 4x TB4, WiFi 6E', 'available', NOW() - INTERVAL '3 months'),
('DELL-PREC-7670-001', 'DELL-PREC-7670', 'Dell', 'Precision 7670', 'CPU: Intel i9-12950HX 16-core, RAM: 128GB DDR5, Storage: 4TB NVMe SSD, Display: 16" UHD+ (3840x2400), GPU: NVIDIA RTX A5500 16GB, Ports: 4x TB4, WiFi 6E', 'available', NOW() - INTERVAL '3 months'),

-- Dell XPS (Premium Developer Laptops)
('DELL-XPS-9520-001', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'CPU: Intel i7-12700H 14-core, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED Touch, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('DELL-XPS-9520-002', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'CPU: Intel i7-12700H 14-core, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED Touch, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('DELL-XPS-9520-003', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'CPU: Intel i7-12700H 14-core, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED Touch, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 2x TB4, WiFi 6E', 'in_transit_to_engineer', NOW() - INTERVAL '2 months'),
('DELL-XPS-9315-001', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'CPU: Intel i7-1360P 12-core, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 13.4" FHD+ Touch, Ports: 2x TB4, WiFi 6E', 'delivered', NOW() - INTERVAL '4 months'),
('DELL-XPS-9315-002', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'CPU: Intel i7-1360P 12-core, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 13.4" FHD+ Touch, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),

-- HP ZBook (Mobile Workstations)
('HP-ZBOOK-G9-001', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'CPU: Intel i9-12900H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 15.6" 4K DreamColor, GPU: NVIDIA RTX A3000 12GB, Ports: 2x TB4, WiFi 6E', 'at_warehouse', NOW() - INTERVAL '3 months'),
('HP-ZBOOK-G9-002', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'CPU: Intel i9-12900H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 15.6" 4K DreamColor, GPU: NVIDIA RTX A3000 12GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '3 months'),
('HP-ZBOOK-G9-003', 'HP-ZBOOK-FUR-G9', 'HP', 'ZBook Fury G9', 'CPU: Intel i9-12950HX, RAM: 128GB DDR5, Storage: 4TB NVMe SSD, Display: 17.3" 4K DreamColor, GPU: NVIDIA RTX A5500 16GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '3 months'),

-- HP EliteBook (Business Laptops)
('HP-ELITE-850-G9-001', 'HP-ELITE-850-G9', 'HP', 'EliteBook 850 G9', 'CPU: Intel i7-1265U, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" FHD, Ports: 2x TB4, WiFi 6E, LTE', 'delivered', NOW() - INTERVAL '4 months'),
('HP-ELITE-850-G9-002', 'HP-ELITE-850-G9', 'HP', 'EliteBook 850 G9', 'CPU: Intel i7-1265U, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" FHD, Ports: 2x TB4, WiFi 6E, LTE', 'available', NOW() - INTERVAL '2 months'),
('HP-ELITE-840-G9-001', 'HP-ELITE-840-G9', 'HP', 'EliteBook 840 G9', 'CPU: Intel i7-1255U, RAM: 16GB DDR4, Storage: 512GB NVMe SSD, Display: 14" FHD, Ports: 2x TB4, WiFi 6E', 'in_transit_to_warehouse', NOW() - INTERVAL '2 months'),

-- Lenovo ThinkPad X1 (Premium Business)
('LENOVO-X1C-G10-001', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WQUXGA (3840x2400), Ports: 2x TB4, WiFi 6E, 5G WWAN', 'available', NOW() - INTERVAL '3 months'),
('LENOVO-X1C-G10-002', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WQUXGA (3840x2400), Ports: 2x TB4, WiFi 6E, 5G WWAN', 'available', NOW() - INTERVAL '3 months'),
('LENOVO-X1C-G10-003', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WQUXGA (3840x2400), Ports: 2x TB4, WiFi 6E, 5G WWAN', 'delivered', NOW() - INTERVAL '4 months'),
('LENOVO-X1C-G10-004', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'CPU: Intel i7-1260P, RAM: 32GB LPDDR5, Storage: 1TB NVMe SSD, Display: 14" WQUXGA (3840x2400), Ports: 2x TB4, WiFi 6E, 5G WWAN', 'at_warehouse', NOW() - INTERVAL '1 month'),

-- Lenovo ThinkPad P Series (Workstations)
('LENOVO-P1-G5-001', 'LENOVO-P1-G5', 'Lenovo', 'ThinkPad P1 Gen 5', 'CPU: Intel i9-12900H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 16" 4K OLED, GPU: NVIDIA RTX A5500 16GB, Ports: 2x TB4, WiFi 6E', 'delivered', NOW() - INTERVAL '4 months'),
('LENOVO-P1-G5-002', 'LENOVO-P1-G5', 'Lenovo', 'ThinkPad P1 Gen 5', 'CPU: Intel i9-12900H, RAM: 64GB DDR5, Storage: 2TB NVMe SSD, Display: 16" 4K OLED, GPU: NVIDIA RTX A5500 16GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('LENOVO-P16-G1-001', 'LENOVO-P16-G1', 'Lenovo', 'ThinkPad P16 Gen 1', 'CPU: Intel i9-12950HX, RAM: 128GB DDR5, Storage: 4TB NVMe SSD, Display: 16" 4K, GPU: NVIDIA RTX A5500 16GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),

-- Apple MacBook Pro (M2/M3 Series)
('APPLE-MBP16-M2MAX-001', 'APPLE-MBP16-M2MAX', 'Apple', 'MacBook Pro 16" M2 Max', 'CPU: Apple M2 Max 12-core, GPU: 38-core, RAM: 96GB Unified, Storage: 2TB SSD, Display: 16.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('APPLE-MBP16-M2MAX-002', 'APPLE-MBP16-M2MAX', 'Apple', 'MacBook Pro 16" M2 Max', 'CPU: Apple M2 Max 12-core, GPU: 38-core, RAM: 96GB Unified, Storage: 2TB SSD, Display: 16.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'in_transit_to_warehouse', NOW() - INTERVAL '1 month'),
('APPLE-MBP16-M2PRO-001', 'APPLE-MBP16-M2PRO', 'Apple', 'MacBook Pro 16" M2 Pro', 'CPU: Apple M2 Pro 12-core, GPU: 19-core, RAM: 32GB Unified, Storage: 1TB SSD, Display: 16.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'delivered', NOW() - INTERVAL '5 months'),
('APPLE-MBP16-M2PRO-002', 'APPLE-MBP16-M2PRO', 'Apple', 'MacBook Pro 16" M2 Pro', 'CPU: Apple M2 Pro 12-core, GPU: 19-core, RAM: 32GB Unified, Storage: 1TB SSD, Display: 16.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'at_warehouse', NOW() - INTERVAL '1 month'),
('APPLE-MBP14-M2PRO-001', 'APPLE-MBP14-M2PRO', 'Apple', 'MacBook Pro 14" M2 Pro', 'CPU: Apple M2 Pro 10-core, GPU: 16-core, RAM: 16GB Unified, Storage: 512GB SSD, Display: 14.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('APPLE-MBP14-M2PRO-002', 'APPLE-MBP14-M2PRO', 'Apple', 'MacBook Pro 14" M2 Pro', 'CPU: Apple M2 Pro 10-core, GPU: 16-core, RAM: 16GB Unified, Storage: 512GB SSD, Display: 14.2" Liquid Retina XDR, Ports: 3x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('APPLE-MBA-M2-001', 'APPLE-MBA-M2', 'Apple', 'MacBook Air 13" M2', 'CPU: Apple M2 8-core, GPU: 10-core, RAM: 24GB Unified, Storage: 1TB SSD, Display: 13.6" Liquid Retina, Ports: 2x TB3, WiFi 6', 'delivered', NOW() - INTERVAL '5 months'),
('APPLE-MBA-M2-002', 'APPLE-MBA-M2', 'Apple', 'MacBook Air 13" M2', 'CPU: Apple M2 8-core, GPU: 10-core, RAM: 24GB Unified, Storage: 1TB SSD, Display: 13.6" Liquid Retina, Ports: 2x TB3, WiFi 6', 'available', NOW() - INTERVAL '2 months'),

-- Microsoft Surface (Premium Windows)
('MSFT-SLS-001', 'MSFT-SLS', 'Microsoft', 'Surface Laptop Studio', 'CPU: Intel i7-11370H, RAM: 32GB LPDDR4x, Storage: 1TB SSD, Display: 14.4" PixelSense Flow Touch, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 2x TB4, WiFi 6', 'available', NOW() - INTERVAL '2 months'),
('MSFT-SLS-002', 'MSFT-SLS', 'Microsoft', 'Surface Laptop Studio', 'CPU: Intel i7-11370H, RAM: 32GB LPDDR4x, Storage: 1TB SSD, Display: 14.4" PixelSense Flow Touch, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 2x TB4, WiFi 6', 'in_transit_to_engineer', NOW() - INTERVAL '1 month'),
('MSFT-SL5-001', 'MSFT-SL5', 'Microsoft', 'Surface Laptop 5', 'CPU: Intel i7-1255U, RAM: 32GB LPDDR5x, Storage: 1TB SSD, Display: 15" PixelSense Touch, Ports: 1x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),

-- ASUS ROG/ZenBook (Performance/Creator)
('ASUS-ZEN-PRO-001', 'ASUS-ZEN-PRO-15', 'ASUS', 'ZenBook Pro 15 OLED', 'CPU: Intel i9-12900H, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 15.6" 4K OLED Touch, GPU: NVIDIA RTX 3060 6GB, Ports: 2x TB4, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),
('ASUS-ROG-G14-001', 'ASUS-ROG-G14', 'ASUS', 'ROG Zephyrus G14', 'CPU: AMD Ryzen 9 6900HS, RAM: 32GB DDR5, Storage: 1TB NVMe SSD, Display: 14" QHD+ 120Hz, GPU: AMD Radeon RX 6800S 8GB, Ports: 1x USB-C, WiFi 6E', 'available', NOW() - INTERVAL '2 months'),

-- Acer (Budget-Friendly Performance)
('ACER-SWIFT-X-001', 'ACER-SWIFT-X', 'Acer', 'Swift X', 'CPU: AMD Ryzen 7 5800U, RAM: 16GB LPDDR4x, Storage: 512GB NVMe SSD, Display: 14" FHD IPS, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 1x USB-C, WiFi 6', 'available', NOW() - INTERVAL '2 months'),
('ACER-SWIFT-X-002', 'ACER-SWIFT-X', 'Acer', 'Swift X', 'CPU: AMD Ryzen 7 5800U, RAM: 16GB LPDDR4x, Storage: 512GB NVMe SSD, Display: 14" FHD IPS, GPU: NVIDIA RTX 3050 Ti 4GB, Ports: 1x USB-C, WiFi 6', 'available', NOW() - INTERVAL '2 months');

-- ============================================
-- SHIPMENTS (Various Statuses & Bulk)
-- ============================================

-- Shipment 1: DELIVERED (Single laptop - completed 3 weeks ago)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (1, 1, 'delivered', 'SCOP-80001', 'FedEx Express', 'FDX8001234567', 
    (NOW() - INTERVAL '28 days')::date, NOW() - INTERVAL '26 days', NOW() - INTERVAL '23 days', 
    NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days', 
    'Dell XPS delivered to Alice Johnson. Standard delivery. Engineer confirmed all equipment working perfectly. Setup completed on-site.', 
    NOW() - INTERVAL '30 days', NOW() - INTERVAL '18 days');

-- Shipment 2: DELIVERED (BULK - 3 laptops - completed 2 weeks ago)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (2, 5, 'delivered', 'SCOP-80002', 'UPS Next Day Air', 'UPS8002345678', 
    (NOW() - INTERVAL '21 days')::date, NOW() - INTERVAL '19 days', NOW() - INTERVAL '16 days', 
    NOW() - INTERVAL '13 days', NOW() - INTERVAL '11 days', 
    'BULK SHIPMENT: 3 HP ZBook workstations delivered to Emily Rodriguez. High-value equipment for video editing team. All units tested and configured. Engineer satisfaction confirmed.', 
    NOW() - INTERVAL '23 days', NOW() - INTERVAL '11 days');

-- Shipment 3: IN TRANSIT TO ENGINEER (Expected delivery tomorrow)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, eta_to_engineer, notes, created_at, updated_at) 
VALUES (3, 8, 'in_transit_to_engineer', 'SCOP-80003', 'DHL Express', 'DHL8003456789', 
    (NOW() - INTERVAL '8 days')::date, NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '1 day', (NOW() + INTERVAL '1 day')::timestamp, 
    'Lenovo ThinkPad X1 Carbon en route to Henry Thompson. Priority delivery. Estimated arrival tomorrow by 10:30 AM. Signature required.', 
    NOW() - INTERVAL '10 days', NOW() - INTERVAL '12 hours');

-- Shipment 4: IN TRANSIT TO ENGINEER (BULK - 5 MacBooks for new team)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, eta_to_engineer, notes, created_at, updated_at) 
VALUES (1, 4, 'in_transit_to_engineer', 'SCOP-80004', 'FedEx Priority Overnight', 'FDX8004567890', 
    (NOW() - INTERVAL '7 days')::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', (NOW() + INTERVAL '2 days')::timestamp, 
    'BULK SHIPMENT: 5 MacBook Pro M2 Max for iOS development team. Lead engineer Daniel Park coordinating receipt. High-value shipment, extra insurance applied. White-glove delivery service.', 
    NOW() - INTERVAL '9 days', NOW() - INTERVAL '1 day');

-- Shipment 5: RELEASED FROM WAREHOUSE (Waiting for courier pickup)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES (4, 11, 'released_from_warehouse', 'SCOP-80005', 'UPS Ground', 'UPS8005678901', 
    (NOW() - INTERVAL '10 days')::date, NOW() - INTERVAL '8 days', NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '6 hours', 
    'Dell Precision workstation prepared for Karen Lee. Package sealed and labeled. Courier pickup scheduled for today 3:00 PM. Contains: laptop, power adapter, USB-C dock, wireless keyboard/mouse set.', 
    NOW() - INTERVAL '12 days', NOW() - INTERVAL '6 hours');

-- Shipment 6: AT WAREHOUSE (Awaiting engineer assignment)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (5, NULL, 'at_warehouse', 'SCOP-80006', 'FedEx Ground', 'FDX8006789012', 
    (NOW() - INTERVAL '6 days')::date, NOW() - INTERVAL '4 days', NOW() - INTERVAL '1 day', 
    'HP EliteBook received and inspected. Condition: Excellent. Serial number verified. All accessories present. Pending engineer assignment by PM. Ready for immediate deployment.', 
    NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 day');

-- Shipment 7: AT WAREHOUSE (BULK - 4 ThinkPads for new hires)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES (6, NULL, 'at_warehouse', 'SCOP-80007', 'DHL Express', 'DHL8007890123', 
    (NOW() - INTERVAL '5 days')::date, NOW() - INTERVAL '3 days', NOW() - INTERVAL '12 hours', 
    'BULK SHIPMENT: 4 Lenovo ThinkPad X1 Carbon laptops. For DataDrive new hire onboarding next week. All units inspected, tested, and inventory-logged. Awaiting engineer assignments from HR.', 
    NOW() - INTERVAL '7 days', NOW() - INTERVAL '12 hours');

-- Shipment 8: IN TRANSIT TO WAREHOUSE (Expected arrival today)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (7, 19, 'in_transit_to_warehouse', 'SCOP-80008', 'UPS Next Day Air', 'UPS8008901234', 
    (NOW() - INTERVAL '3 days')::date, NOW() - INTERVAL '1 day', 
    'ASUS ZenBook Pro for NextGen creative team. Expected warehouse arrival: Today by 4:30 PM. Fragile - OLED display. Handle with extra care. Track actively.', 
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '6 hours');

-- Shipment 9: PICKED UP FROM CLIENT (Just left client location)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES (8, 21, 'picked_up_from_client', 'SCOP-80009', 'FedEx Express', 'FDX8009012345', 
    (NOW() - INTERVAL '2 days')::date, NOW() - INTERVAL '2 hours', 
    'Microsoft Surface Laptop Studio picked up from Enterprise Solutions Group HQ. Driver confirmed pickup at 2:15 PM. Package dimensions: 18x13x4 inches, weight: 8 lbs. In transit to regional hub.', 
    NOW() - INTERVAL '4 days', NOW() - INTERVAL '2 hours');

-- Shipment 10: PICKUP SCHEDULED (Tomorrow morning)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (1, 2, 'pickup_from_client_scheduled', 'SCOP-80010', 'DHL Express', 
    (NOW() + INTERVAL '1 day')::date, 
    'Dell XPS 13 Plus for Bob Smith. Pickup confirmed for tomorrow 9:00 AM - 12:00 PM. Contact: Sarah Johnson, +1-555-0100. Location: TechCorp Building 3, Shipping/Receiving dock. Package ready.', 
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day');

-- Shipment 11: PICKUP SCHEDULED (BULK - Next week)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (2, NULL, 'pickup_from_client_scheduled', 'SCOP-80011', 'FedEx Freight', 
    (NOW() + INTERVAL '5 days')::date, 
    'BULK SHIPMENT: 6 HP EliteBook laptops for department expansion. Large pickup scheduled for Monday 10:00 AM. Freight carrier required. Contact: IT Manager, +1-555-0200. Dock #7 at Innovate Solutions warehouse.', 
    NOW() - INTERVAL '2 days', NOW());

-- Shipment 12: PENDING PICKUP (Just created, awaiting form)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (3, 9, 'pending_pickup_from_client', 'SCOP-80012', 
    (NOW() + INTERVAL '7 days')::date, 
    'Lenovo ThinkPad P16 workstation for Isabella Garcia. Awaiting pickup form submission from Global Tech Services. Anticipated pickup in 1 week. High-end GPU workstation for ML/AI development.', 
    NOW() - INTERVAL '1 day', NOW());

-- Shipment 13: PENDING PICKUP (BULK - Urgent)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES (4, NULL, 'pending_pickup_from_client', 'SCOP-80013', 
    (NOW() + INTERVAL '3 days')::date, 
    'URGENT BULK SHIPMENT: 3 Apple MacBook Pro for emergency team scaling. Client requested expedited processing. Awaiting pickup form. Contact PM immediately upon form receipt. Priority handling required.', 
    NOW() - INTERVAL '6 hours', NOW());

-- Shipment 14: DELIVERED (Old shipment - 2 months ago)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (5, 14, 'delivered', 'SCOP-80014', 'UPS Ground', 'UPS8014567890', 
    (NOW() - INTERVAL '68 days')::date, NOW() - INTERVAL '66 days', NOW() - INTERVAL '63 days', 
    NOW() - INTERVAL '60 days', NOW() - INTERVAL '58 days', 
    'Lenovo ThinkPad P1 delivered to Nathan Brown 2 months ago. Complete lifecycle. Engineer satisfaction: Excellent. No issues reported. Archived shipment.', 
    NOW() - INTERVAL '70 days', NOW() - INTERVAL '58 days');

-- Shipment 15: DELIVERED (BULK - Last month)
INSERT INTO shipments (client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES (6, 17, 'delivered', 'SCOP-80015', 'FedEx Express', 'FDX8015678901', 
    (NOW() - INTERVAL '38 days')::date, NOW() - INTERVAL '36 days', NOW() - INTERVAL '33 days', 
    NOW() - INTERVAL '30 days', NOW() - INTERVAL '28 days', 
    'BULK SHIPMENT DELIVERED: 2 Acer Swift X laptops to Quinn Anderson. Budget-friendly performance laptops for data analysis team. Both units functioning well. No issues.', 
    NOW() - INTERVAL '40 days', NOW() - INTERVAL '28 days');

-- ============================================
-- SHIPMENT LAPTOPS JUNCTION
-- ============================================

-- Shipment 1: Single Dell XPS
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (1, 7);

-- Shipment 2: BULK - 3 HP ZBooks
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (2, 10), (2, 11), (2, 12);

-- Shipment 3: Single Lenovo X1 Carbon
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (3, 17);

-- Shipment 4: BULK - 5 MacBook Pros (using some that need to be added)
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (4, 23), (4, 25), (4, 26), (4, 27), (4, 28);

-- Shipment 5: Single Dell Precision
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (5, 1);

-- Shipment 6: Single HP EliteBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (6, 14);

-- Shipment 7: BULK - 4 Lenovo ThinkPads
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (7, 15), (7, 16), (7, 18), (7, 19);

-- Shipment 8: Single ASUS ZenBook
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (8, 32);

-- Shipment 9: Single Surface Laptop Studio
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (9, 30);

-- Shipment 10: Single Dell XPS 13
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (10, 8);

-- Shipment 11: BULK - 6 HP EliteBooks
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (11, 13), (11, 14), (11, 2), (11, 3), (11, 4), (11, 5);

-- Shipment 12: Single Lenovo P16
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (12, 22);

-- Shipment 13: BULK - 3 MacBook Pros
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (13, 24), (13, 29), (13, 6);

-- Shipment 14: Single Lenovo P1
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (14, 20);

-- Shipment 15: BULK - 2 Acer laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (15, 34), (15, 35);

-- ============================================
-- PICKUP FORMS
-- ============================================

-- Forms for all shipments
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(1, 110, NOW() - INTERVAL '30 days', 
 ('{"contact_name":"Sarah Mitchell","contact_email":"sarah.mitchell@techcorp.com","contact_phone":"+1-555-0101","pickup_address":"100 Tech Plaza, Building 1, Loading Dock A","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '28 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell 130W USB-C power adapter, Dell Pro Briefcase, wireless mouse, USB-C to HDMI adapter","special_instructions":"Building requires 24hr advance security notification. Call Sarah 30 minutes before arrival."}')::jsonb),

(2, 111, NOW() - INTERVAL '23 days',
 ('{"contact_name":"Michael Chen","contact_email":"michael.chen@innovate.io","contact_phone":"+1-555-0201","pickup_address":"200 Innovation Way, Warehouse Building, Bay 5","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '21 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":20.0,"bulk_height":14.0,"bulk_weight":48.5,"include_accessories":true,"accessories_description":"3x HP 200W power adapters, 3x HP ZBook carrying cases, 3x USB-C docks with dual monitor support, 3x wireless keyboard/mouse combos, HDMI cables","special_instructions":"BULK SHIPMENT - Heavy equipment. Forklift assistance available. Use loading dock entrance. Contact Michael directly for warehouse access."}')::jsonb),

(3, 112, NOW() - INTERVAL '10 days',
 ('{"contact_name":"Jennifer Wang","contact_email":"jennifer.wang@globaltech.com","contact_phone":"+1-555-0301","pickup_address":"300 Global Drive, Floor 8, IT Department","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '8 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo 65W USB-C adapter, ThinkPad Professional Backpack, ThinkPad Bluetooth Silent Mouse, USB-C dock","special_instructions":"Priority delivery. Use visitor parking and check in at main reception. Package will be ready at IT department, ask for Jennifer."}')::jsonb),

(4, 110, NOW() - INTERVAL '9 days',
 ('{"contact_name":"David Lee","contact_email":"david.lee@techcorp.com","contact_phone":"+1-555-0102","pickup_address":"100 Tech Plaza, Building 3, iOS Development Lab","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '7 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":5,"number_of_boxes":3,"assignment_type":"bulk","bulk_length":26.0,"bulk_width":22.0,"bulk_height":16.0,"bulk_weight":62.0,"include_accessories":true,"accessories_description":"5x Apple USB-C cables (2m), 5x Apple Magic Mouse, 5x Apple Magic Keyboard, 5x USB-C to USB-A adapters, 5x premium laptop sleeves, AppleCare+ documentation","special_instructions":"BULK HIGH-VALUE SHIPMENT: 5 MacBook Pro M2 Max. Total value >$20,000. Signature required. Extra insurance applied. White-glove service. Contact David 1 hour before arrival. Building 3 has separate security checkpoint."}')::jsonb),

(5, 113, NOW() - INTERVAL '12 days',
 ('{"contact_name":"Robert Chen","contact_email":"robert.chen@digitaldynamics.com","contact_phone":"+1-555-0401","pickup_address":"400 Digital Blvd, Suite 1200","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '10 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell 240W power adapter, Dell Precision Backpack, Kensington SD5700T Thunderbolt 4 dock, Dell Premier Wireless Keyboard and Mouse","special_instructions":"High-end workstation. Handle with extreme care. Reception will have package ready. Building has security escort available."}')::jsonb),

(6, 113, NOW() - INTERVAL '8 days',
 ('{"contact_name":"Linda Martinez","contact_email":"linda.martinez@cloudventures.com","contact_phone":"+1-555-0501","pickup_address":"500 Cloud Street, Building A, Floor 3","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '6 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP 65W USB-C adapter, HP Executive Leather Top Load, HP wireless mouse, USB-C travel dock","special_instructions":"Standard office pickup. Package ready at front desk. Visitor parking available in Lot B."}')::jsonb),

(7, 110, NOW() - INTERVAL '7 days',
 ('{"contact_name":"Patricia Johnson","contact_email":"patricia.johnson@datadrive.com","contact_phone":"+1-555-0601","pickup_address":"600 Data Lane, Main Building, HR Department","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '5 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":4,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":18.0,"bulk_height":12.0,"bulk_weight":38.0,"include_accessories":true,"accessories_description":"4x Lenovo 65W USB-C adapters, 4x ThinkPad Essential Backpacks, 4x Lenovo Bluetooth mice, 4x USB-C mini docks, cable organizers","special_instructions":"BULK - New hire onboarding equipment. Pickup from HR department. Ring doorbell for access. Contact Patricia for any questions. Scheduled for Monday morning."}')::jsonb),

(8, 110, NOW() - INTERVAL '5 days',
 ('{"contact_name":"Thomas Anderson","contact_email":"thomas.anderson@nextgensw.com","contact_phone":"+1-555-0701","pickup_address":"700 Innovation Court, Creative Studio, Floor 2","pickup_city":"Portland","pickup_state":"OR","pickup_zip":"97201","pickup_date":"' || to_char(NOW() - INTERVAL '3 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"ASUS 120W USB-C adapter, ASUS ROG backpack, wireless gaming mouse, USB-C hub with SD card reader","special_instructions":"FRAGILE - OLED Display. Creative equipment for design team. Handle with extra care. Package marked FRAGILE on all sides. Contact Thomas upon arrival."}')::jsonb),

(9, 110, NOW() - INTERVAL '4 days',
 ('{"contact_name":"Amanda Wilson","contact_email":"amanda.wilson@enterprisesg.com","contact_phone":"+1-555-0801","pickup_address":"800 Enterprise Ave, 42nd Floor, Executive Suite","pickup_city":"New York","pickup_state":"NY","pickup_zip":"10001","pickup_date":"' || to_char(NOW() - INTERVAL '2 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Microsoft Surface power supply, Surface Pen, Surface Arc Mouse, Microsoft designer compact keyboard, premium leather sleeve","special_instructions":"High-security building. Courier must check in at main lobby security desk with photo ID. Amanda will meet in lobby. Cannot access 42nd floor without escort."}')::jsonb),

(10, 110, NOW() - INTERVAL '3 days',
 ('{"contact_name":"Sarah Mitchell","contact_email":"sarah.mitchell@techcorp.com","contact_phone":"+1-555-0103","pickup_address":"100 Tech Plaza, Building 1, Mail Room","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() + INTERVAL '1 day', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell 65W USB-C adapter, Dell EcoLoop Pro Backpack, wireless mouse, USB-C to USB-A adapter","special_instructions":"Pickup confirmed for tomorrow 9-12 AM. Call Sarah 30 minutes before arrival. Mail room is on ground floor, east wing."}')::jsonb),

(11, 111, NOW() - INTERVAL '2 days',
 ('{"contact_name":"James Brown","contact_email":"james.brown@innovate.io","contact_phone":"+1-555-0202","pickup_address":"200 Innovation Way, Warehouse Complex, Dock 7","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() + INTERVAL '5 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":6,"number_of_boxes":3,"assignment_type":"bulk","bulk_length":28.0,"bulk_width":24.0,"bulk_height":18.0,"bulk_weight":72.0,"include_accessories":true,"accessories_description":"6x HP 65W USB-C adapters, 6x HP Professional carrying cases, 6x wireless keyboards and mice, 6x USB-C docking stations, cable management kits, extended warranty cards","special_instructions":"LARGE BULK SHIPMENT - Freight carrier required. Scheduled for Monday 10:00 AM. Forklift available. Use Dock 7 entrance. Contact James 1 day before to confirm. Department expansion equipment."}')::jsonb),

(12, 112, NOW() - INTERVAL '1 day',
 ('{"contact_name":"Jennifer Wang","contact_email":"jennifer.wang@globaltech.com","contact_phone":"+1-555-0302","pickup_address":"300 Global Drive, R&D Lab, Secure Wing","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() + INTERVAL '7 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo 230W Slim Tip adapter, ThinkPad Professional 16-inch Topload, Lenovo Legion M600 Wireless Mouse, Thunderbolt 4 workstation dock","special_instructions":"High-end GPU workstation for ML/AI development. Requires security clearance for R&D wing access. Contact Jennifer 24 hours in advance. Escort required through secure area."}')::jsonb),

(13, 113, NOW() - INTERVAL '6 hours',
 ('{"contact_name":"Mark Stevens","contact_email":"mark.stevens@digitaldynamics.com","contact_phone":"+1-555-0402","pickup_address":"400 Digital Blvd, Main Lobby","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() + INTERVAL '3 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":20.0,"bulk_height":14.0,"bulk_weight":42.0,"include_accessories":true,"accessories_description":"3x Apple USB-C cables, 3x Apple Magic Mouse, 3x Apple Magic Keyboard, 3x USB-C to USB-A adapters, 3x premium sleeves","special_instructions":"URGENT BULK SHIPMENT - Emergency team scaling. Expedited pickup requested. Contact PM immediately upon pickup completion. High priority. Total value: $15,000+."}')::jsonb),

(14, 113, NOW() - INTERVAL '70 days',
 ('{"contact_name":"Rebecca Thompson","contact_email":"rebecca.thompson@cloudventures.com","contact_phone":"+1-555-0502","pickup_address":"500 Cloud Street, Building B, Floor 5, ML Lab","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '68 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo 230W power adapter, ThinkPad P1 carrying case, Lenovo ThinkPad Thunderbolt 4 Workstation Dock, wireless keyboard and mouse, USB-C cables","special_instructions":"High-end GPU workstation for ML/AI development. Package ready at ML Lab. Requires badge access - contact Rebecca for building entry. Handle with care - expensive equipment."}')::jsonb),

(15, 110, NOW() - INTERVAL '40 days',
 ('{"contact_name":"Timothy Roberts","contact_email":"timothy.roberts@datadrive.com","contact_phone":"+1-555-0602","pickup_address":"600 Data Lane, Building C, Data Analytics Department","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '38 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":2,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":20.0,"bulk_width":16.0,"bulk_height":10.0,"bulk_weight":28.0,"include_accessories":true,"accessories_description":"2x Acer 90W power adapters, 2x laptop carrying cases, 2x wireless mice, 2x USB-C hubs with SD card readers, 2x laptop cooling pads","special_instructions":"BULK SHIPMENT: 2 Acer Swift X laptops for data analysis team. Budget-friendly performance laptops. Packages ready at Data Analytics Dept. Contact Timothy for access. Standard business pickup."}')::jsonb);

-- ============================================
-- RECEPTION REPORTS
-- ============================================

INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES 
(1, 103, NOW() - INTERVAL '23 days', 'Dell XPS 15 9520 received in excellent condition. Original factory packaging intact with all seals present. Serial number DELL-XPS-9520-001 verified against manifest. Visual inspection: No scratches, dents, or cosmetic damage. Power-on test: Successful. BIOS info matches specs. All ports tested functional. Included accessories: 130W USB-C adapter (original Dell), carrying case (new condition), wireless mouse (sealed), HDMI adapter (sealed). Device charged to 85%. Logged into inventory system. Ready for assignment.', 
 ARRAY['/uploads/reception/shipment001_exterior.jpg', '/uploads/reception/shipment001_screen.jpg', '/uploads/reception/shipment001_accessories.jpg']),

(2, 104, NOW() - INTERVAL '16 days', 'BULK RECEPTION: 3x HP ZBook Studio G9 workstations received. Heavy shipment - two-person lift required. All boxes in excellent condition, no shipping damage detected. Serial numbers verified: HP-ZBOOK-G9-001, HP-ZBOOK-G9-002, HP-ZBOOK-G9-003. Individual inspection performed on each unit: All 3 units: Original HP packaging, factory seals intact. Display quality excellent - 4K DreamColor panels tested, no dead pixels on any unit. RTX A3000 GPUs tested with benchmark - all performing within spec. RAM: 64GB DDR5 verified on all units. Storage: 2TB NVMe verified. All accessories present: 3x 200W power adapters, 3x carrying cases, 3x docking stations, 3x wireless peripherals, documentation sets. All units powered on successfully, BIOS verified. Firmware up to date. Total shipment value: $15,000+. Extra photos taken for high-value documentation. Stored in secure area.', 
 ARRAY['/uploads/reception/shipment002_bulk_overview.jpg', '/uploads/reception/shipment002_unit1.jpg', '/uploads/reception/shipment002_unit2.jpg', '/uploads/reception/shipment002_unit3.jpg', '/uploads/reception/shipment002_accessories.jpg', '/uploads/reception/shipment002_serial_numbers.jpg']),

(3, 105, NOW() - INTERVAL '3 days', 'Lenovo ThinkPad X1 Carbon Gen 10 received in pristine condition. Premium business laptop - flagship model. Serial: LENOVO-X1C-G10-001. Package condition: Perfect - original Lenovo retail box with all factory seals. Device inspection: Carbon fiber chassis in immaculate condition. WQUXGA 4K display tested - crystal clear, no defects. Keyboard legendary ThinkPad quality verified. TrackPoint and touchpad working perfectly. All ports functional including 2x Thunderbolt 4. 5G WWAN module detected and functional. Battery health: 100% (new unit). Included accessories: 65W USB-C adapter (original Lenovo), premium backpack (excellent quality), wireless mouse (sealed), USB-C dock (tested - all ports functional). Weight: Ultra-light as expected. Performance test: Boots in seconds, runs smooth. Ready for priority delivery to engineer.', 
 ARRAY['/uploads/reception/shipment003_laptop.jpg', '/uploads/reception/shipment003_display.jpg', '/uploads/reception/shipment003_full_setup.jpg']),

(4, 103, NOW() - INTERVAL '2 days', 'BULK HIGH-VALUE SHIPMENT: 5x Apple MacBook Pro 16" M2 Max received. WHITE-GLOVE HANDLING APPLIED. Total shipment value: $22,000+. Extra security measures in place. Original Apple packaging on all 5 units - sealed retail boxes with Apple authenticity stickers intact. Serial numbers verified against Apple GSX: APPLE-MBP16-M2MAX-001, APPLE-MBP16-M2PRO-001, APPLE-MBP16-M2PRO-002, APPLE-MBP14-M2PRO-001, APPLE-MBP14-M2PRO-002. Individual unit inspection: Unit 1-5: M2 Max/Pro chips verified, Liquid Retina XDR displays perfect (no dead pixels, no backlight bleeding), all Thunderbolt 4 ports functional. Storage: 1TB-2TB verified per unit. RAM: 32GB-96GB unified memory verified. Battery cycles: 0-2 (essentially new). macOS Ventura pre-installed and updated. All charging cables and adapters present (original Apple). Premium sleeves, Magic Mouse, Magic Keyboard sets included. AppleCare+ documentation verified. iOS development team equipment. Stored in climate-controlled secure cage. Comprehensive photographic documentation completed for insurance. Lead engineer Daniel Park notified of arrival.', 
 ARRAY['/uploads/reception/shipment004_all_units.jpg', '/uploads/reception/shipment004_unit1_box.jpg', '/uploads/reception/shipment004_unit1_open.jpg', '/uploads/reception/shipment004_serial_verification.jpg', '/uploads/reception/shipment004_accessories_complete.jpg', '/uploads/reception/shipment004_applecare_docs.jpg', '/uploads/reception/shipment004_secure_storage.jpg']),

(5, 104, NOW() - INTERVAL '5 days', 'Dell Precision 5570 mobile workstation received. HIGH-END UNIT. Serial: DELL-PREC-5570-001. Condition: Excellent. Original Dell Precision packaging with foam inserts. Visual inspection: Premium aluminum chassis - no blemishes. UHD+ display tested: 3840x2400 resolution verified, color accuracy excellent, no dead pixels. NVIDIA RTX A2000 8GB tested: GPU-Z confirms specs, benchmark run successful. CPU: Intel i9-12900H 14-core verified, stress test passed. RAM: 64GB DDR5 confirmed. Storage: 2TB NVMe verified, read/write speeds excellent. All 4x Thunderbolt 4 ports tested with dock - full functionality. WiFi 6E detected and working. Included accessories: 240W power adapter (heavy-duty), premium backpack, Thunderbolt 4 dock (tested - dual 4K @ 60Hz confirmed), wireless keyboard/mouse set (Dell Premier line). Workstation for heavy computational tasks. Performance tests all passed. Ready for engineer assignment.', 
 ARRAY['/uploads/reception/shipment005_workstation.jpg', '/uploads/reception/shipment005_display.jpg', '/uploads/reception/shipment005_dock_test.jpg']),

(6, 105, NOW() - INTERVAL '1 day', 'HP EliteBook 850 G9 received and processed. Business-class laptop. Serial: HP-ELITE-850-G9-002. Package in good condition - original HP business packaging. Device inspection: Silver aluminum chassis - professional finish, no damage. 15.6" FHD display tested - clear, good viewing angles. Keyboard with backlighting functional. Touchpad responsive. Security features: Fingerprint reader tested and functional, IR camera for Windows Hello working, TPM 2.0 detected. All ports working: 2x Thunderbolt 4, USB-A, HDMI 2.0, headphone jack. LTE module detected (optional upgrade present). Accessories included: 65W USB-C adapter, professional carrying case, wireless mouse, travel dock. Battery health: 95% (lightly used, but excellent condition). Enterprise features verified: vPro, Sure Start, Sure View privacy screen working. Suitable for business professional. Pending engineer assignment by PM. Ready for immediate deployment.', 
 ARRAY['/uploads/reception/shipment006_laptop.jpg', '/uploads/reception/shipment006_ports.jpg']),

(7, 103, NOW() - INTERVAL '12 hours', 'BULK SHIPMENT: 4x Lenovo ThinkPad X1 Carbon Gen 10 for new hire onboarding. Serial numbers: LENOVO-X1C-G10-002, LENOVO-X1C-G10-003, LENOVO-X1C-G10-004, and one additional unit. All boxes received in excellent condition. Consistent inspection across all 4 units: All units: Carbon fiber chassis perfect, WQUXGA displays all tested (no issues), keyboards excellent (classic ThinkPad quality), TrackPoints all functional, 5G WWAN modules in all units, battery health 98-100%, WiFi 6E confirmed. Accessories for all units: 4x 65W USB-C adapters (Lenovo original), 4x ThinkPad backpacks (Essential line), 4x wireless mice (Lenovo Bluetooth), 4x USB-C mini docks (all tested functional), 4x cable organizers. New hire onboarding package. HR coordinating engineer assignments for next week Monday. All units ready for deployment. Stored together in onboarding equipment area. IT documentation prepared for each unit.', 
 ARRAY['/uploads/reception/shipment007_bulk_four_units.jpg', '/uploads/reception/shipment007_lineup.jpg', '/uploads/reception/shipment007_accessories_set.jpg', '/uploads/reception/shipment007_serial_tags.jpg']);

-- ============================================
-- DELIVERY FORMS
-- ============================================

INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES 
(1, 1, NOW() - INTERVAL '18 days', 'Dell XPS 15 9520 delivered successfully to Alice Johnson at TechCorp, New York office. Delivery date/time: [timestamp]. Contact: Alice Johnson, +1-555-1001. Delivery address: 101 Main St, Apt 5B, New York, NY 10001. Arrived on schedule, signature obtained. Package inspection with engineer present: Box unopened, seals intact. Engineer opened package on-site. Laptop condition: Pristine, as expected. Powered on successfully: Boot time: <30 seconds. Display quality verified: 4K OLED gorgeous, no dead pixels. Engineer very pleased. Accessories verified: Power adapter, carrying case, wireless mouse, HDMI adapter - all present. Engineer tested: USB-C ports, keyboard, trackpad, speakers - all working perfectly. Software: Windows 11 Pro, pre-installed corporate apps detected. Connected to company WiFi successfully. Windows Hello face recognition setup completed. First impressions: Engineer expressed high satisfaction, mentioned "exactly what I needed for development work." Delivered by: John (FedEx courier). Weather: Clear, no issues. No damages, no issues reported. Setup assistance provided - connected to external monitor via USB-C dock, verified dual-display setup. Engineer confirmed receipt and satisfaction. Delivery complete. Follow-up scheduled for next week to ensure no issues.', 
 ARRAY['/uploads/delivery/shipment001_signature.jpg', '/uploads/delivery/shipment001_laptop_open.jpg', '/uploads/delivery/shipment001_with_engineer.jpg', '/uploads/delivery/shipment001_setup_complete.jpg']),

(2, 5, NOW() - INTERVAL '11 days', 'BULK DELIVERY: 3x HP ZBook Studio G9 workstations delivered to Emily Rodriguez, Innovate Solutions, Austin office. HIGH-VALUE SHIPMENT ($15,000+). Delivery date/time: [timestamp]. Contact: Emily Rodriguez, +1-555-2001. Delivery address: 505 River Rd, Building 2, Austin, TX 78702. Two-person delivery team (shipment heavy). Engineer and IT manager present for receipt. Package inspection: 2 boxes, both in excellent condition, no shipping damage. Engineer and IT manager opened boxes on-site. All 3 laptops condition: Perfect, as expected from warehouse photos. Sequential setup of all 3 units: Unit 1: Powered on, 4K DreamColor display verified - stunning quality. RTX A3000 GPU tested with video editing software - performance excellent. Unit 2: Same thorough testing, all aspects perfect. Unit 3: Consistent quality across all units. Accessories for all 3 verified: 3x 200W power adapters, 3x carrying cases, 3x docking stations (all tested - dual 4K monitors @ 60Hz confirmed), 3x wireless keyboard/mouse combos. Software setup: DaVinci Resolve tested on all units - renders fast, no issues. Adobe Creative Cloud installed and verified. Network rendering setup tested between units. Engineer feedback: "These are exactly what our video editing team needs. The color accuracy on the DreamColor displays is perfect for our work. GPU performance is impressive." IT Manager confirmed all security protocols followed. All units connected to company domain. BitLocker encryption enabled on all drives. All docking stations configured for dual-monitor setups. Engineer satisfaction: Excellent. Team lead very happy. Delivery team: Marcus and James (UPS). Setup time: 90 minutes (thorough testing justified). No issues reported. All units operating perfectly. 30-day follow-up scheduled to ensure continued satisfaction. BULK DELIVERY SUCCESSFUL.', 
 ARRAY['/uploads/delivery/shipment002_delivery_team.jpg', '/uploads/delivery/shipment002_three_units.jpg', '/uploads/delivery/shipment002_unit1_running.jpg', '/uploads/delivery/shipment002_dual_monitor_setup.jpg', '/uploads/delivery/shipment002_it_manager_approval.jpg', '/uploads/delivery/shipment002_all_setups_complete.jpg']),

(14, 14, NOW() - INTERVAL '58 days', 'Lenovo ThinkPad P1 Gen 5 delivered to Nathan Brown, Cloud Ventures, Denver office 2 months ago (archived delivery). Delivery was smooth and successful. High-end GPU workstation for ML/AI development. Engineer confirmed receipt of laptop, RTX A5500 GPU, and all accessories. Initial setup completed on-site. Engineer tested GPU compute performance with TensorFlow - confirmed working excellently. All components verified functional. Engineer expressed high satisfaction with the device specifications. Connected to company network and development environment. Docker and Kubernetes setups verified. Python environment with CUDA support confirmed working. Delivery completed without issues. 30-day follow-up was conducted - engineer reported zero issues, performance excellent. Unit has been in productive use since delivery. No support tickets filed related to this device. Archived record for historical reference.', 
 ARRAY['/uploads/delivery/shipment014_laptop.jpg', '/uploads/delivery/shipment014_signature.jpg']),

(15, 17, NOW() - INTERVAL '28 days', 'BULK DELIVERY: 2x Acer Swift X laptops delivered to Quinn Anderson, DataDrive, Chicago office last month. Budget-friendly performance laptops for data analysis team. Delivery date/time: [timestamp]. Contact: Quinn Anderson, +1-555-6001. Delivery address: 1717 Michigan Ave, Chicago, IL 60611. Both boxes in good condition. Engineer present for receipt. Package inspection: Both units unopened, seals intact. Engineer opened boxes with delivery team present. Both laptops condition: Excellent, new condition. Unit 1: Powered on successfully, 14" FHD IPS display clear and bright. RTX 3050 Ti tested with data visualization tools - sufficient performance. Unit 2: Consistent quality with Unit 1. Accessories verified: 2x power adapters, 2x laptop bags, 2x wireless mice. Software configuration: Python, Pandas, Jupyter notebooks installed and tested on both units. Power BI and Tableau installed successfully. Sample datasets loaded and visualizations run - performance acceptable for data analysis workloads. Engineer feedback: "Good value for the price. Sufficient performance for our data analysis needs. The RTX GPU helps with some ML model training." Both units connected to company network. Remote management software installed. Security policies applied. Engineer satisfaction: Good (not excellent, but met expectations for budget tier). Delivery courier: Sarah (FedEx). No issues during delivery. Both units operational. Follow-up completed - both units functioning well in daily use.', 
 ARRAY['/uploads/delivery/shipment015_both_laptops.jpg', '/uploads/delivery/shipment015_setup.jpg', '/uploads/delivery/shipment015_data_viz_test.jpg']);

-- ============================================
-- AUDIT LOGS (Sample Activity)
-- ============================================

INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details) VALUES
-- Recent shipment creation
(100, 'shipment_created', 'shipment', 13, NOW() - INTERVAL '6 hours', '{"action":"shipment_created","jira_ticket":"SCOP-80013","status":"pending_pickup_from_client","urgent":true}'),
(113, 'pickup_form_submitted', 'pickup_form', 13, NOW() - INTERVAL '6 hours', '{"shipment_id":13,"bulk":true,"laptops":3}'),

-- Status updates from today
(100, 'status_updated', 'shipment', 9, NOW() - INTERVAL '2 hours', '{"old_status":"in_transit_to_warehouse","new_status":"picked_up_from_client","tracking":"FDX8009012345"}'),
(104, 'shipment_tracking_updated', 'shipment', 8, NOW() - INTERVAL '6 hours', '{"tracking_number":"UPS8008901234","status":"in_transit_to_warehouse","eta":"today 4:30 PM"}'),

-- Yesterday activity
(105, 'reception_report_created', 'reception_report', 7, NOW() - INTERVAL '12 hours', '{"shipment_id":7,"bulk":true,"units":4}'),
(105, 'status_updated', 'shipment', 7, NOW() - INTERVAL '12 hours', '{"old_status":"in_transit_to_warehouse","new_status":"at_warehouse"}'),

-- Recent deliveries
(100, 'status_updated', 'shipment', 2, NOW() - INTERVAL '11 days', '{"old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(101, 'delivery_form_created', 'delivery_form', 2, NOW() - INTERVAL '11 days', '{"shipment_id":2,"engineer":"Emily Rodriguez","bulk":true}'),

-- Warehouse assignments
(107, 'engineer_assigned', 'shipment', 5, NOW() - INTERVAL '6 hours', '{"engineer_id":11,"engineer_name":"Karen Lee"}'),
(107, 'status_updated', 'shipment', 5, NOW() - INTERVAL '6 hours', '{"old_status":"at_warehouse","new_status":"released_from_warehouse"}');

-- ============================================
-- SUMMARY OUTPUT
-- ============================================

SELECT '========================================' AS separator;
SELECT 'ENHANCED SAMPLE DATA LOADED SUCCESSFULLY!' AS message;
SELECT '========================================' AS separator;
SELECT '' AS blank;

SELECT 'DATABASE SUMMARY' AS section;
SELECT '----------------' AS underline;
SELECT COUNT(*) || ' users' AS count FROM users WHERE role != 'admin'
UNION ALL SELECT COUNT(*) || ' client companies' FROM client_companies
UNION ALL SELECT COUNT(*) || ' software engineers' FROM software_engineers
UNION ALL SELECT COUNT(*) || ' laptops' FROM laptops
UNION ALL SELECT COUNT(*) || ' shipments' FROM shipments
UNION ALL SELECT COUNT(*) || ' pickup forms' FROM pickup_forms
UNION ALL SELECT COUNT(*) || ' reception reports' FROM reception_reports
UNION ALL SELECT COUNT(*) || ' delivery forms' FROM delivery_forms
UNION ALL SELECT COUNT(*) || ' audit log entries' FROM audit_logs;

SELECT '' AS blank;
SELECT 'SHIPMENTS BY STATUS' AS section;
SELECT '-------------------' AS underline;

SELECT 
    status,
    COUNT(*) as count,
    CASE 
        WHEN SUM(CASE WHEN (SELECT COUNT(*) FROM shipment_laptops sl WHERE sl.shipment_id = s.id) > 1 THEN 1 ELSE 0 END) > 0 
        THEN '(' || SUM(CASE WHEN (SELECT COUNT(*) FROM shipment_laptops sl WHERE sl.shipment_id = s.id) > 1 THEN 1 ELSE 0 END) || ' bulk)'
        ELSE ''
    END as bulk_info
FROM shipments s
GROUP BY status 
ORDER BY 
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

SELECT '' AS blank;
SELECT 'BULK SHIPMENTS' AS section;
SELECT '--------------' AS underline;
SELECT 
    s.id,
    s.jira_ticket_number,
    s.status,
    COUNT(sl.laptop_id) as laptop_count,
    cc.name as client
FROM shipments s
JOIN shipment_laptops sl ON sl.shipment_id = s.id
JOIN client_companies cc ON cc.id = s.client_company_id
GROUP BY s.id, s.jira_ticket_number, s.status, cc.name
HAVING COUNT(sl.laptop_id) > 1
ORDER BY COUNT(sl.laptop_id) DESC, s.id;

SELECT '' AS blank;
SELECT '========================================' AS separator;
SELECT 'Ready for testing with realistic data!' AS message;
SELECT '========================================' AS separator;

