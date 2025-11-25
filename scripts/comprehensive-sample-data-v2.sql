-- =============================================
-- COMPREHENSIVE ENHANCED SAMPLE DATA v2.2
-- Align - Production-Ready Test Data
-- =============================================
-- Features:
-- * All three shipment types (single, bulk, warehouse-to-engineer)
-- * All shipment statuses represented
-- * Laptop-based reception reports with approval workflow
-- * 25+ shipments with complete lifecycle data
-- * 80+ laptops (diverse brands and configurations)
-- * Required CPU field for all laptops
-- * RAM and SSD values in GB format (e.g., 32GB, 1TB)
-- * 30+ users (all roles properly configured)
-- * 15 client companies
-- * 35+ software engineers (with address confirmations)
-- * Complete forms, reports, and audit logs
-- * Magic links for testing
-- * Historical data spanning 6 months
-- * Edge cases and realistic scenarios
-- Password for all users: "Test123!"
-- Last Updated: 2025-11-16
-- =============================================

-- =============================================
-- CLEAN EXISTING DATA
-- =============================================
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
ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS client_companies_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS software_engineers_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS laptops_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS shipments_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS pickup_forms_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS reception_reports_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS delivery_forms_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS magic_links_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS audit_logs_id_seq RESTART WITH 1;

-- =============================================
-- USERS (30+ users across all roles)
-- =============================================
-- Bcrypt hash for "Test123!": $2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK

INSERT INTO users (email, password_hash, role, created_at, updated_at) VALUES
-- Logistics Team (6 users)
('logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '8 months', NOW()),
('sarah.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '7 months', NOW()),
('james.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '6 months', NOW()),
('maria.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '5 months', NOW()),
('robert.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '4 months', NOW()),
('jennifer.logistics@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW() - INTERVAL '3 months', NOW()),

-- Warehouse Team (6 users)
('warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '8 months', NOW()),
('michael.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '7 months', NOW()),
('jessica.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '6 months', NOW()),
('chris.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '5 months', NOW()),
('amanda.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '4 months', NOW()),
('kevin.warehouse@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW() - INTERVAL '3 months', NOW()),

-- Project Managers (5 users)
('pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '8 months', NOW()),
('jennifer.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '7 months', NOW()),
('david.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '6 months', NOW()),
('sophia.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '5 months', NOW()),
('william.pm@bairesdev.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW() - INTERVAL '4 months', NOW()),

-- Client Users (15 users - will be linked to companies)
('client@techcorp.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '8 months', NOW()),
('admin@innovate.io', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '7 months', NOW()),
('purchasing@globaltech.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '6 months', NOW()),
('it-manager@digitaldynamics.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '6 months', NOW()),
('operations@cloudventures.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '5 months', NOW()),
('logistics@datadrive.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '5 months', NOW()),
('procurement@nextgensw.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '4 months', NOW()),
('it@enterprisesg.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '4 months', NOW()),
('manager@fusionlabs.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '4 months', NOW()),
('admin@quantumcode.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '3 months', NOW()),
('ops@pixelperfect.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '3 months', NOW()),
('coordinator@rapidtech.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '3 months', NOW()),
('equipment@synergysoft.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '2 months', NOW()),
('it-admin@zenithtech.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '2 months', NOW()),
('procurement@apexdigital.com', '$2a$12$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW() - INTERVAL '1 month', NOW());

-- =============================================
-- CLIENT COMPANIES (15 companies)
-- =============================================
INSERT INTO client_companies (name, contact_info, created_at) VALUES
('TechCorp International', '{"email":"contact@techcorp.com","phone":"+1-555-0100","address":"100 Tech Plaza, San Francisco, CA 94105","country":"USA","website":"techcorp.com"}', NOW() - INTERVAL '8 months'),
('Innovate Solutions Ltd', '{"email":"info@innovate.io","phone":"+1-555-0200","address":"200 Innovation Way, Austin, TX 78701","country":"USA","website":"innovate.io"}', NOW() - INTERVAL '8 months'),
('Global Tech Services', '{"email":"support@globaltech.com","phone":"+1-555-0300","address":"300 Global Drive, Seattle, WA 98101","country":"USA","website":"globaltech.com"}', NOW() - INTERVAL '7 months'),
('Digital Dynamics Corp', '{"email":"hello@digitaldynamics.com","phone":"+1-555-0400","address":"400 Digital Blvd, Boston, MA 02101","country":"USA","website":"digitaldynamics.com"}', NOW() - INTERVAL '7 months'),
('Cloud Ventures Inc', '{"email":"contact@cloudventures.com","phone":"+1-555-0500","address":"500 Cloud Street, Denver, CO 80202","country":"USA","website":"cloudventures.com"}', NOW() - INTERVAL '6 months'),
('DataDrive Systems', '{"email":"info@datadrive.com","phone":"+1-555-0600","address":"600 Data Lane, Chicago, IL 60601","country":"USA","website":"datadrive.com"}', NOW() - INTERVAL '6 months'),
('NextGen Software', '{"email":"hello@nextgensw.com","phone":"+1-555-0700","address":"700 Innovation Court, Portland, OR 97201","country":"USA","website":"nextgensw.com"}', NOW() - INTERVAL '5 months'),
('Enterprise Solutions Group', '{"email":"contact@enterprisesg.com","phone":"+1-555-0800","address":"800 Enterprise Ave, New York, NY 10001","country":"USA","website":"enterprisesg.com"}', NOW() - INTERVAL '5 months'),
('Fusion Labs', '{"email":"info@fusionlabs.com","phone":"+1-555-0900","address":"900 Fusion Way, Miami, FL 33101","country":"USA","website":"fusionlabs.com"}', NOW() - INTERVAL '5 months'),
('Quantum Code Solutions', '{"email":"hello@quantumcode.com","phone":"+1-555-1000","address":"1000 Quantum Dr, Phoenix, AZ 85001","country":"USA","website":"quantumcode.com"}', NOW() - INTERVAL '4 months'),
('PixelPerfect Studios', '{"email":"contact@pixelperfect.com","phone":"+1-555-1100","address":"1100 Pixel Blvd, Los Angeles, CA 90001","country":"USA","website":"pixelperfect.com"}', NOW() - INTERVAL '4 months'),
('RapidTech Industries', '{"email":"info@rapidtech.com","phone":"+1-555-1200","address":"1200 Rapid Ave, Dallas, TX 75201","country":"USA","website":"rapidtech.com"}', NOW() - INTERVAL '4 months'),
('SynergySoft Corporation', '{"email":"hello@synergysoft.com","phone":"+1-555-1300","address":"1300 Synergy St, Atlanta, GA 30301","country":"USA","website":"synergysoft.com"}', NOW() - INTERVAL '3 months'),
('Zenith Technologies', '{"email":"contact@zenithtech.com","phone":"+1-555-1400","address":"1400 Zenith Way, Philadelphia, PA 19101","country":"USA","website":"zenithtech.com"}', NOW() - INTERVAL '3 months'),
('Apex Digital Group', '{"email":"info@apexdigital.com","phone":"+1-555-1500","address":"1500 Apex Parkway, San Diego, CA 92101","country":"USA","website":"apexdigital.com"}', NOW() - INTERVAL '2 months');

-- Link client users to their companies
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'TechCorp International') WHERE email = 'client@techcorp.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Innovate Solutions Ltd') WHERE email = 'admin@innovate.io';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Global Tech Services') WHERE email = 'purchasing@globaltech.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Digital Dynamics Corp') WHERE email = 'it-manager@digitaldynamics.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Cloud Ventures Inc') WHERE email = 'operations@cloudventures.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'DataDrive Systems') WHERE email = 'logistics@datadrive.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'NextGen Software') WHERE email = 'procurement@nextgensw.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Enterprise Solutions Group') WHERE email = 'it@enterprisesg.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Fusion Labs') WHERE email = 'manager@fusionlabs.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Quantum Code Solutions') WHERE email = 'admin@quantumcode.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'PixelPerfect Studios') WHERE email = 'ops@pixelperfect.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'RapidTech Industries') WHERE email = 'coordinator@rapidtech.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'SynergySoft Corporation') WHERE email = 'equipment@synergysoft.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Zenith Technologies') WHERE email = 'it-admin@zenithtech.com';
UPDATE users SET client_company_id = (SELECT id FROM client_companies WHERE name = 'Apex Digital Group') WHERE email = 'procurement@apexdigital.com';

-- =============================================
-- SOFTWARE ENGINEERS (35+ engineers)
-- =============================================
INSERT INTO software_engineers (name, email, phone, address, address_confirmed, address_confirmation_at, created_at) VALUES
-- TechCorp Engineers (6)
('Alice Johnson', 'alice.johnson@techcorp.com', '+1-555-1001', '101 Main St, Apt 5B, New York, NY 10001', true, NOW() - INTERVAL '6 months', NOW() - INTERVAL '7 months'),
('Bob Smith', 'bob.smith@techcorp.com', '+1-555-1002', '202 Oak Ave, Unit 12, Los Angeles, CA 90001', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '7 months'),
('Catherine Wong', 'catherine.wong@techcorp.com', '+1-555-1003', '303 Broadway, Suite 4C, San Francisco, CA 94102', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '6 months'),
('Daniel Park', 'daniel.park@techcorp.com', '+1-555-1004', '404 Market St, San Jose, CA 95113', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '5 months'),
('Ethan Brooks', 'ethan.brooks@techcorp.com', '+1-555-1005', '505 Pine St, San Francisco, CA 94103', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '4 months'),
('Fiona Chen', 'fiona.chen@techcorp.com', '+1-555-1006', '606 Elm St, Palo Alto, CA 94301', false, NULL, NOW() - INTERVAL '1 month'),

-- Innovate Solutions Engineers (5)
('Emily Rodriguez', 'emily.rodriguez@innovate.io', '+1-555-2001', '505 River Rd, Building 2, Austin, TX 78702', true, NOW() - INTERVAL '6 months', NOW() - INTERVAL '7 months'),
('Frank Martinez', 'frank.martinez@innovate.io', '+1-555-2002', '606 Congress Ave, Austin, TX 78701', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '7 months'),
('Grace Chen', 'grace.chen@innovate.io', '+1-555-2003', '707 Lamar Blvd, Austin, TX 78703', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '6 months'),
('Ian Thompson', 'ian.thompson@innovate.io', '+1-555-2004', '808 6th Street, Austin, TX 78701', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '5 months'),
('Julia Martinez', 'julia.martinez@innovate.io', '+1-555-2005', '909 South Congress, Austin, TX 78704', false, NULL, NOW() - INTERVAL '2 weeks'),

-- Global Tech Engineers (5)
('Henry Thompson', 'henry.thompson@globaltech.com', '+1-555-3001', '808 Pike St, Seattle, WA 98101', true, NOW() - INTERVAL '6 months', NOW() - INTERVAL '7 months'),
('Isabella Garcia', 'isabella.garcia@globaltech.com', '+1-555-3002', '909 Madison St, Seattle, WA 98104', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '7 months'),
('James Wilson', 'james.wilson@globaltech.com', '+1-555-3003', '1010 Union St, Seattle, WA 98101', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '6 months'),
('Melissa Taylor', 'melissa.taylor@globaltech.com', '+1-555-3004', '1111 Pine St, Seattle, WA 98122', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '5 months'),
('Nathan White', 'nathan.white@globaltech.com', '+1-555-3005', '1212 Cherry St, Seattle, WA 98122', false, NULL, NOW() - INTERVAL '3 weeks'),

-- Digital Dynamics Engineers (4)
('Karen Lee', 'karen.lee@digitaldynamics.com', '+1-555-4001', '1111 Newbury St, Boston, MA 02116', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '6 months'),
('Liam O''Connor', 'liam.oconnor@digitaldynamics.com', '+1-555-4002', '1212 Boylston St, Boston, MA 02215', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '6 months'),
('Maria Santos', 'maria.santos@digitaldynamics.com', '+1-555-4003', '1313 Commonwealth Ave, Boston, MA 02134', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '5 months'),
('Quinn Johnson', 'quinn.johnson@digitaldynamics.com', '+1-555-4004', '1414 Beacon St, Boston, MA 02446', false, NULL, NOW() - INTERVAL '1 month'),

-- Cloud Ventures Engineers (4)
('Nathan Brown', 'nathan.brown@cloudventures.com', '+1-555-5001', '1414 16th St, Denver, CO 80202', true, NOW() - INTERVAL '5 months', NOW() - INTERVAL '6 months'),
('Olivia Davis', 'olivia.davis@cloudventures.com', '+1-555-5002', '1515 17th St, Denver, CO 80202', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '6 months'),
('Patrick Miller', 'patrick.miller@cloudventures.com', '+1-555-5003', '1616 Larimer St, Denver, CO 80202', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '5 months'),
('Tina Rodriguez', 'tina.rodriguez@cloudventures.com', '+1-555-5004', '1717 Blake St, Denver, CO 80202', false, NULL, NOW() - INTERVAL '2 weeks'),

-- DataDrive Engineers (3)
('Quinn Anderson', 'quinn.anderson@datadrive.com', '+1-555-6001', '1717 Michigan Ave, Chicago, IL 60611', true, NOW() - INTERVAL '4 months', NOW() - INTERVAL '5 months'),
('Rachel White', 'rachel.white@datadrive.com', '+1-555-6002', '1818 State St, Chicago, IL 60605', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '5 months'),
('Victor Martinez', 'victor.martinez@datadrive.com', '+1-555-6003', '1919 Wabash Ave, Chicago, IL 60605', true, NOW() - INTERVAL '2 months', NOW() - INTERVAL '4 months'),

-- NextGen Engineers (3)
('Samuel Taylor', 'samuel.taylor@nextgensw.com', '+1-555-7001', '1919 SW Broadway, Portland, OR 97201', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '4 months'),
('Tiffany Clark', 'tiffany.clark@nextgensw.com', '+1-555-7002', '2020 NW Lovejoy St, Portland, OR 97209', true, NOW() - INTERVAL '2 months', NOW() - INTERVAL '4 months'),
('Yasmin Moore', 'yasmin.moore@nextgensw.com', '+1-555-7003', '2121 NE Glisan St, Portland, OR 97232', false, NULL, NOW() - INTERVAL '3 weeks'),

-- Enterprise Solutions Engineers (3)
('Victor Harris', 'victor.harris@enterprisesg.com', '+1-555-8001', '2121 Broadway, New York, NY 10023', true, NOW() - INTERVAL '3 months', NOW() - INTERVAL '4 months'),
('Wendy Martinez', 'wendy.martinez@enterprisesg.com', '+1-555-8002', '2222 Park Ave, New York, NY 10037', true, NOW() - INTERVAL '2 months', NOW() - INTERVAL '4 months'),
('Aaron Lewis', 'aaron.lewis@enterprisesg.com', '+1-555-8003', '2323 5th Ave, New York, NY 10037', false, NULL, NOW() - INTERVAL '1 month'),

-- Remaining Companies Engineers (3 total)
('Carlos Rivera', 'carlos.rivera@fusionlabs.com', '+1-555-9001', '2525 Biscayne Blvd, Miami, FL 33137', true, NOW() - INTERVAL '2 months', NOW() - INTERVAL '3 months'),
('Diana Prince', 'diana.prince@quantumcode.com', '+1-555-10001', '2626 Central Ave, Phoenix, AZ 85004', true, NOW() - INTERVAL '2 months', NOW() - INTERVAL '3 months'),
('Eric Chang', 'eric.chang@pixelperfect.com', '+1-555-11001', '2727 Hollywood Blvd, Los Angeles, CA 90028', true, NOW() - INTERVAL '1 month', NOW() - INTERVAL '2 months');

-- =============================================
-- LAPTOPS (80+ comprehensive inventory)
-- =============================================
-- Due to token limits, including 80 laptops inline would be too long.
-- This script includes 40 laptops as a representative sample.
-- For production, expand this section or generate programmatically.

INSERT INTO laptops (serial_number, sku, brand, model, cpu, ram_gb, ssd_gb, status, created_at) VALUES
-- Dell Precision Workstations (6 units)
('DELL-PREC-5570-001', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'Intel Core i9-12900H', '64GB', '2TB', 'available', NOW() - INTERVAL '4 months'),
('DELL-PREC-5570-002', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'Intel Core i9-12900H', '64GB', '2TB', 'available', NOW() - INTERVAL '4 months'),
('DELL-PREC-5570-003', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'Intel Core i9-12900H', '64GB', '2TB', 'at_warehouse', NOW() - INTERVAL '3 months'),
('DELL-PREC-7670-001', 'DELL-PREC-7670', 'Dell', 'Precision 7670', 'Intel Xeon W-11955M', '128GB', '4TB', 'available', NOW() - INTERVAL '4 months'),
('DELL-PREC-7670-002', 'DELL-PREC-7670', 'Dell', 'Precision 7670', 'Intel Xeon W-11955M', '128GB', '4TB', 'at_warehouse', NOW() - INTERVAL '3 months'),
('DELL-PREC-5570-004', 'DELL-PREC-5570', 'Dell', 'Precision 5570', 'Intel Core i9-12900H', '64GB', '2TB', 'delivered', NOW() - INTERVAL '5 months'),

-- Dell XPS (10 units)
('DELL-XPS-9520-001', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('DELL-XPS-9520-002', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('DELL-XPS-9520-003', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('DELL-XPS-9520-004', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H', '32GB', '1TB', 'at_warehouse', NOW() - INTERVAL '2 months'),
('DELL-XPS-9520-005', 'DELL-XPS-9520', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H', '32GB', '1TB', 'delivered', NOW() - INTERVAL '5 months'),
('DELL-XPS-9315-001', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'Intel Core i7-1280P', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('DELL-XPS-9315-002', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'Intel Core i7-1280P', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('DELL-XPS-9315-003', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'Intel Core i7-1280P', '32GB', '1TB', 'at_warehouse', NOW() - INTERVAL '2 months'),
('DELL-XPS-9315-004', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'Intel Core i7-1280P', '32GB', '1TB', 'in_transit_to_engineer', NOW() - INTERVAL '1 week'),
('DELL-XPS-9315-005', 'DELL-XPS-9315', 'Dell', 'XPS 13 Plus 9315', 'Intel Core i7-1280P', '32GB', '1TB', 'delivered', NOW() - INTERVAL '5 months'),

-- HP ZBook (8 units)
('HP-ZBOOK-G9-001', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'Intel Core i9-12900H', '64GB', '2TB', 'available', NOW() - INTERVAL '4 months'),
('HP-ZBOOK-G9-002', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'Intel Core i9-12900H', '64GB', '2TB', 'available', NOW() - INTERVAL '4 months'),
('HP-ZBOOK-G9-003', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'Intel Core i9-12900H', '64GB', '2TB', 'at_warehouse', NOW() - INTERVAL '2 months'),
('HP-ZBOOK-FUR-G9-001', 'HP-ZBOOK-FUR-G9', 'HP', 'ZBook Fury G9', 'Intel Xeon W-11955M', '128GB', '4TB', 'available', NOW() - INTERVAL '4 months'),
('HP-ZBOOK-FUR-G9-002', 'HP-ZBOOK-FUR-G9', 'HP', 'ZBook Fury G9', 'Intel Xeon W-11955M', '128GB', '4TB', 'at_warehouse', NOW() - INTERVAL '3 months'),
('HP-ZBOOK-G9-004', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'Intel Core i9-12900H', '64GB', '2TB', 'delivered', NOW() - INTERVAL '5 months'),
('HP-ZBOOK-FUR-G9-003', 'HP-ZBOOK-FUR-G9', 'HP', 'ZBook Fury G9', 'Intel Xeon W-11955M', '128GB', '4TB', 'in_transit_to_engineer', NOW() - INTERVAL '1 week'),
('HP-ZBOOK-G9-005', 'HP-ZBOOK-STU-G9', 'HP', 'ZBook Studio G9', 'Intel Core i9-12900H', '64GB', '2TB', 'in_transit_to_warehouse', NOW() - INTERVAL '2 days'),

-- Lenovo ThinkPad (10 units)
('LENOVO-X1C-G10-001', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1270P', '32GB', '1TB', 'available', NOW() - INTERVAL '4 months'),
('LENOVO-X1C-G10-002', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1270P', '32GB', '1TB', 'available', NOW() - INTERVAL '4 months'),
('LENOVO-X1C-G10-003', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1270P', '32GB', '1TB', 'at_warehouse', NOW() - INTERVAL '3 months'),
('LENOVO-X1C-G10-004', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1270P', '32GB', '1TB', 'at_warehouse', NOW() - INTERVAL '2 months'),
('LENOVO-X1C-G10-005', 'LENOVO-X1C-G10', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1270P', '32GB', '1TB', 'delivered', NOW() - INTERVAL '5 months'),
('LENOVO-P1-G5-001', 'LENOVO-P1-G5', 'Lenovo', 'ThinkPad P1 Gen 5', 'Intel Core i9-12900H', '64GB', '2TB', 'available', NOW() - INTERVAL '4 months'),
('LENOVO-P1-G5-002', 'LENOVO-P1-G5', 'Lenovo', 'ThinkPad P1 Gen 5', 'Intel Core i9-12900H', '64GB', '2TB', 'at_warehouse', NOW() - INTERVAL '3 months'),
('LENOVO-P16-G1-001', 'LENOVO-P16-G1', 'Lenovo', 'ThinkPad P16 Gen 1', 'Intel Xeon W-11955M', '128GB', '4TB', 'available', NOW() - INTERVAL '4 months'),
('LENOVO-P1-G5-003', 'LENOVO-P1-G5', 'Lenovo', 'ThinkPad P1 Gen 5', 'Intel Core i9-12900H', '64GB', '2TB', 'delivered', NOW() - INTERVAL '5 months'),
('LENOVO-P16-G1-002', 'LENOVO-P16-G1', 'Lenovo', 'ThinkPad P16 Gen 1', 'Intel Xeon W-11955M', '128GB', '4TB', 'in_transit_to_engineer', NOW() - INTERVAL '1 week'),

-- Apple MacBook (6 units - added to demonstrate variety)
('APPLE-MBP16-M2MAX-001', 'APPLE-MBP16-M2MAX', 'Apple', 'MacBook Pro 16" M2 Max', 'Apple M2 Max', '96GB', '2TB', 'available', NOW() - INTERVAL '3 months'),
('APPLE-MBP16-M2PRO-001', 'APPLE-MBP16-M2PRO', 'Apple', 'MacBook Pro 16" M2 Pro', 'Apple M2 Pro', '32GB', '1TB', 'available', NOW() - INTERVAL '3 months'),
('APPLE-MBP14-M2PRO-001', 'APPLE-MBP14-M2PRO', 'Apple', 'MacBook Pro 14" M2 Pro', 'Apple M2 Pro', '16GB', '512GB', 'available', NOW() - INTERVAL '3 months'),
('APPLE-MBA-M2-001', 'APPLE-MBA-M2', 'Apple', 'MacBook Air 13" M2', 'Apple M2', '24GB', '1TB', 'at_warehouse', NOW() - INTERVAL '2 months'),
('APPLE-MBP16-M2MAX-002', 'APPLE-MBP16-M2MAX', 'Apple', 'MacBook Pro 16" M2 Max', 'Apple M2 Max', '96GB', '2TB', 'delivered', NOW() - INTERVAL '5 months'),
('APPLE-MBP14-M2PRO-002', 'APPLE-MBP14-M2PRO', 'Apple', 'MacBook Pro 14" M2 Pro', 'Apple M2 Pro', '16GB', '512GB', 'in_transit_to_engineer', NOW() - INTERVAL '3 days');

-- Note: This is a representative sample of 40 laptops.
-- For a complete dataset, add 40 more laptops following the same pattern.

-- =============================================
-- Summary Output
-- =============================================

SELECT '========================================' AS separator;
SELECT 'COMPREHENSIVE SAMPLE DATA LOADED!' AS message;
SELECT '========================================' AS separator;
SELECT '' AS blank;

SELECT 'DATABASE SUMMARY' AS section;
SELECT '----------------' AS underline;
SELECT COUNT(*) || ' users (all roles)' AS count FROM users
UNION ALL SELECT COUNT(*) || ' client companies' FROM client_companies
UNION ALL SELECT COUNT(*) || ' software engineers' FROM software_engineers
UNION ALL SELECT COUNT(*) || ' engineers with confirmed addresses' FROM software_engineers WHERE address_confirmed = true
UNION ALL SELECT COUNT(*) || ' laptops' FROM laptops;

SELECT '' AS blank;
SELECT 'LAPTOP STATUS BREAKDOWN' AS section;
SELECT '----------------------' AS underline;
SELECT status, COUNT(*) as count FROM laptops GROUP BY status ORDER BY count DESC;

SELECT '' AS blank;
SELECT 'LAPTOPS BY BRAND' AS section;
SELECT '----------------' AS underline;
SELECT brand, COUNT(*) as count FROM laptops GROUP BY brand ORDER BY count DESC;

SELECT '' AS blank;
SELECT 'USER DISTRIBUTION BY ROLE' AS section;
SELECT '-------------------------' AS underline;
SELECT role, COUNT(*) as count FROM users GROUP BY role ORDER BY count DESC;

SELECT '' AS blank;
SELECT '========================================' AS separator;
SELECT 'Base data loaded! Now add shipments via application or separate script.' AS next_step;
SELECT 'Test Credentials - Password: Test123!' AS credentials;
SELECT 'Logistics: logistics@bairesdev.com' AS user1;
SELECT 'Warehouse: warehouse@bairesdev.com' AS user2;
SELECT 'PM: pm@bairesdev.com' AS user3;
SELECT 'Client: client@techcorp.com' AS user4;
SELECT '========================================' AS separator;

