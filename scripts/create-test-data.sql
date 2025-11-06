-- ============================================
-- Test Data for Laptop Tracking System
-- ============================================
-- This script creates comprehensive test data for development and testing
-- Run this after migrations and user creation

-- Start transaction
BEGIN;

-- ============================================
-- 1. CLIENT COMPANIES
-- ============================================
\echo ''
\echo '>>> Creating Client Companies...'

INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES 
    ('Google LLC', 'contact@google.com | +1-650-253-0000 | Mountain View, CA', NOW() - INTERVAL '90 days', NOW() - INTERVAL '30 days'),
    ('Microsoft Corporation', 'support@microsoft.com | +1-425-882-8080 | Redmond, WA', NOW() - INTERVAL '85 days', NOW() - INTERVAL '25 days'),
    ('Amazon Web Services', 'aws-support@amazon.com | +1-206-266-1000 | Seattle, WA', NOW() - INTERVAL '80 days', NOW() - INTERVAL '20 days'),
    ('Meta Platforms Inc', 'tech@meta.com | +1-650-543-4800 | Menlo Park, CA', NOW() - INTERVAL '75 days', NOW() - INTERVAL '15 days'),
    ('Apple Inc', 'enterprise@apple.com | +1-408-996-1010 | Cupertino, CA', NOW() - INTERVAL '70 days', NOW() - INTERVAL '10 days')
ON CONFLICT (LOWER(name)) DO NOTHING;

\echo '>>> Client companies created.'

-- ============================================
-- 2. SOFTWARE ENGINEERS
-- ============================================
\echo ''
\echo '>>> Creating Software Engineers...'

INSERT INTO software_engineers (name, email, address, phone, address_confirmed, address_confirmation_at, created_at, updated_at)
VALUES 
    -- Engineers with confirmed addresses
    ('John Smith', 'john.smith@example.com', '123 Main St, Apt 4B, New York, NY 10001', '+1-212-555-0101', TRUE, NOW() - INTERVAL '5 days', NOW() - INTERVAL '60 days', NOW() - INTERVAL '5 days'),
    ('Maria Garcia', 'maria.garcia@example.com', '456 Oak Avenue, San Francisco, CA 94102', '+1-415-555-0102', TRUE, NOW() - INTERVAL '4 days', NOW() - INTERVAL '58 days', NOW() - INTERVAL '4 days'),
    ('David Chen', 'david.chen@example.com', '789 Pine Street, Suite 301, Seattle, WA 98101', '+1-206-555-0103', TRUE, NOW() - INTERVAL '6 days', NOW() - INTERVAL '55 days', NOW() - INTERVAL '6 days'),
    ('Sarah Johnson', 'sarah.johnson@example.com', '321 Elm Drive, Austin, TX 78701', '+1-512-555-0104', TRUE, NOW() - INTERVAL '3 days', NOW() - INTERVAL '50 days', NOW() - INTERVAL '3 days'),
    ('Michael Brown', 'michael.brown@example.com', '654 Maple Court, Boston, MA 02101', '+1-617-555-0105', TRUE, NOW() - INTERVAL '7 days', NOW() - INTERVAL '48 days', NOW() - INTERVAL '7 days'),
    
    -- Engineers pending confirmation
    ('Emily Wilson', 'emily.wilson@example.com', '987 Cedar Lane, Denver, CO 80201', '+1-303-555-0106', FALSE, NULL, NOW() - INTERVAL '45 days', NOW() - INTERVAL '2 days'),
    ('James Martinez', 'james.martinez@example.com', '147 Birch Road, Miami, FL 33101', '+1-305-555-0107', FALSE, NULL, NOW() - INTERVAL '40 days', NOW() - INTERVAL '1 day'),
    ('Linda Anderson', 'linda.anderson@example.com', '258 Spruce Street, Chicago, IL 60601', '+1-312-555-0108', TRUE, NOW() - INTERVAL '8 days', NOW() - INTERVAL '35 days', NOW() - INTERVAL '8 days'),
    ('Robert Taylor', 'robert.taylor@example.com', '369 Willow Way, Portland, OR 97201', '+1-503-555-0109', TRUE, NOW() - INTERVAL '9 days', NOW() - INTERVAL '30 days', NOW() - INTERVAL '9 days'),
    ('Jennifer Lee', 'jennifer.lee@example.com', '741 Ash Boulevard, Atlanta, GA 30301', '+1-404-555-0110', FALSE, NULL, NOW() - INTERVAL '25 days', NOW() - INTERVAL '1 day'),
    
    -- More engineers
    ('William Thomas', 'william.thomas@example.com', '852 Redwood Drive, San Diego, CA 92101', '+1-619-555-0111', TRUE, NOW() - INTERVAL '10 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days'),
    ('Jessica White', 'jessica.white@example.com', '963 Cypress Street, Phoenix, AZ 85001', '+1-602-555-0112', FALSE, NULL, NOW() - INTERVAL '15 days', NOW()),
    ('Christopher Harris', 'chris.harris@example.com', '159 Palm Avenue, Las Vegas, NV 89101', '+1-702-555-0113', TRUE, NOW() - INTERVAL '12 days', NOW() - INTERVAL '18 days', NOW() - INTERVAL '12 days'),
    ('Amanda Clark', 'amanda.clark@example.com', '357 Magnolia Court, Nashville, TN 37201', '+1-615-555-0114', TRUE, NOW() - INTERVAL '11 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '11 days'),
    ('Daniel Robinson', 'daniel.robinson@example.com', '753 Dogwood Lane, Philadelphia, PA 19101', '+1-215-555-0115', FALSE, NULL, NOW() - INTERVAL '10 days', NOW())
ON CONFLICT (LOWER(email)) DO NOTHING;

\echo '>>> Software engineers created.'

-- ============================================
-- 3. LAPTOPS
-- ============================================
\echo ''
\echo '>>> Creating Laptops...'

INSERT INTO laptops (serial_number, brand, model, specs, status, created_at, updated_at)
VALUES 
    -- Dell Laptops
    ('DELL-SN-001', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H, 16GB RAM, 512GB SSD, NVIDIA RTX 3050 Ti', 'delivered', NOW() - INTERVAL '50 days', NOW() - INTERVAL '10 days'),
    ('DELL-SN-002', 'Dell', 'XPS 15 9520', 'Intel Core i7-12700H, 16GB RAM, 512GB SSD, NVIDIA RTX 3050 Ti', 'delivered', NOW() - INTERVAL '48 days', NOW() - INTERVAL '8 days'),
    ('DELL-SN-003', 'Dell', 'Latitude 7430', 'Intel Core i5-1245U, 16GB RAM, 256GB SSD', 'at_warehouse', NOW() - INTERVAL '45 days', NOW() - INTERVAL '5 days'),
    ('DELL-SN-004', 'Dell', 'Latitude 7430', 'Intel Core i5-1245U, 16GB RAM, 256GB SSD', 'in_transit_to_engineer', NOW() - INTERVAL '40 days', NOW() - INTERVAL '2 days'),
    ('DELL-SN-005', 'Dell', 'Precision 5570', 'Intel Core i9-12900H, 32GB RAM, 1TB SSD, NVIDIA RTX A2000', 'delivered', NOW() - INTERVAL '35 days', NOW() - INTERVAL '5 days'),
    
    -- HP Laptops
    ('HP-SN-001', 'HP', 'EliteBook 840 G9', 'Intel Core i7-1265U, 16GB RAM, 512GB SSD', 'delivered', NOW() - INTERVAL '55 days', NOW() - INTERVAL '15 days'),
    ('HP-SN-002', 'HP', 'EliteBook 840 G9', 'Intel Core i7-1265U, 16GB RAM, 512GB SSD', 'at_warehouse', NOW() - INTERVAL '30 days', NOW() - INTERVAL '3 days'),
    ('HP-SN-003', 'HP', 'ZBook Studio G9', 'Intel Core i7-12800H, 32GB RAM, 1TB SSD, NVIDIA RTX A2000', 'in_transit_to_warehouse', NOW() - INTERVAL '25 days', NOW() - INTERVAL '1 day'),
    ('HP-SN-004', 'HP', 'ProBook 450 G9', 'Intel Core i5-1235U, 8GB RAM, 256GB SSD', 'in_transit_to_engineer', NOW() - INTERVAL '20 days', NOW() - INTERVAL '1 day'),
    ('HP-SN-005', 'HP', 'EliteBook 1040 G9', 'Intel Core i7-1265U, 16GB RAM, 512GB SSD', 'available', NOW() - INTERVAL '15 days', NOW()),
    
    -- Lenovo Laptops
    ('LENOVO-SN-001', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1260P, 16GB RAM, 512GB SSD', 'delivered', NOW() - INTERVAL '60 days', NOW() - INTERVAL '20 days'),
    ('LENOVO-SN-002', 'Lenovo', 'ThinkPad X1 Carbon Gen 10', 'Intel Core i7-1260P, 16GB RAM, 512GB SSD', 'at_warehouse', NOW() - INTERVAL '28 days', NOW() - INTERVAL '3 days'),
    ('LENOVO-SN-003', 'Lenovo', 'ThinkPad P1 Gen 5', 'Intel Core i9-12900H, 32GB RAM, 1TB SSD, NVIDIA RTX A3000', 'in_transit_to_warehouse', NOW() - INTERVAL '22 days', NOW() - INTERVAL '1 day'),
    ('LENOVO-SN-004', 'Lenovo', 'ThinkPad T14s Gen 3', 'AMD Ryzen 7 PRO 6850U, 16GB RAM, 512GB SSD', 'delivered', NOW() - INTERVAL '18 days', NOW() - INTERVAL '3 days'),
    ('LENOVO-SN-005', 'Lenovo', 'ThinkPad X13 Gen 3', 'Intel Core i5-1235U, 16GB RAM, 256GB SSD', 'available', NOW() - INTERVAL '12 days', NOW()),
    
    -- Apple MacBooks
    ('APPLE-SN-001', 'Apple', 'MacBook Pro 16" M2 Max', 'Apple M2 Max, 32GB RAM, 1TB SSD', 'delivered', NOW() - INTERVAL '42 days', NOW() - INTERVAL '7 days'),
    ('APPLE-SN-002', 'Apple', 'MacBook Pro 14" M2 Pro', 'Apple M2 Pro, 16GB RAM, 512GB SSD', 'at_warehouse', NOW() - INTERVAL '26 days', NOW() - INTERVAL '2 days'),
    ('APPLE-SN-003', 'Apple', 'MacBook Air 13" M2', 'Apple M2, 16GB RAM, 512GB SSD', 'in_transit_to_warehouse', NOW() - INTERVAL '19 days', NOW() - INTERVAL '1 day'),
    ('APPLE-SN-004', 'Apple', 'MacBook Pro 16" M2 Max', 'Apple M2 Max, 64GB RAM, 2TB SSD', 'available', NOW() - INTERVAL '14 days', NOW()),
    ('APPLE-SN-005', 'Apple', 'MacBook Pro 14" M2 Pro', 'Apple M2 Pro, 32GB RAM, 1TB SSD', 'available', NOW() - INTERVAL '10 days', NOW()),
    
    -- Asus Laptops
    ('ASUS-SN-001', 'Asus', 'ROG Zephyrus G15', 'AMD Ryzen 9 6900HS, 32GB RAM, 1TB SSD, NVIDIA RTX 3080 Ti', 'delivered', NOW() - INTERVAL '38 days', NOW() - INTERVAL '8 days'),
    ('ASUS-SN-002', 'Asus', 'ZenBook Pro 15', 'Intel Core i7-12700H, 16GB RAM, 1TB SSD, NVIDIA RTX 3050', 'at_warehouse', NOW() - INTERVAL '24 days', NOW() - INTERVAL '2 days'),
    ('ASUS-SN-003', 'Asus', 'VivoBook S15', 'Intel Core i5-1235U, 16GB RAM, 512GB SSD', 'available', NOW() - INTERVAL '16 days', NOW()),
    
    -- Microsoft Surface
    ('MSFT-SN-001', 'Microsoft', 'Surface Laptop Studio', 'Intel Core i7-11370H, 32GB RAM, 1TB SSD, NVIDIA RTX 3050 Ti', 'in_transit_to_engineer', NOW() - INTERVAL '21 days', NOW() - INTERVAL '1 day'),
    ('MSFT-SN-002', 'Microsoft', 'Surface Laptop 5', 'Intel Core i7-1255U, 16GB RAM, 512GB SSD', 'available', NOW() - INTERVAL '11 days', NOW())
ON CONFLICT (LOWER(serial_number)) DO NOTHING;

\echo '>>> Laptops created.'

-- ============================================
-- 4. LINK CLIENT USER TO COMPANY
-- ============================================
\echo ''
\echo '>>> Linking client user to Google...'

-- Update client user to link with Google (using subquery to avoid variable conflicts)
UPDATE users 
SET client_company_id = (SELECT id FROM client_companies WHERE LOWER(name) = 'google llc')
WHERE email = 'client@bairesdev.com';

\echo '>>> Client user linked to company.'

-- ============================================
-- 4. SHIPMENTS
-- ============================================
\echo ''
\echo '>>> Creating Shipments...'

-- Shipment 1: Delivered (Complete workflow)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'google llc'),
    (SELECT id FROM software_engineers WHERE email = 'john.smith@example.com'),
    'delivered',
    'FedEx Express',
    'FDX-123456789',
    NOW() - INTERVAL '50 days',
    NOW() - INTERVAL '48 days',
    NOW() - INTERVAL '45 days',
    NOW() - INTERVAL '40 days',
    NOW() - INTERVAL '38 days',
    'Delivered successfully. Engineer confirmed receipt.',
    'LTS-101',
    NOW() - INTERVAL '52 days',
    NOW() - INTERVAL '38 days'
);

-- Shipment 2: Delivered (Multiple laptops)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'microsoft corporation'),
    (SELECT id FROM software_engineers WHERE email = 'maria.garcia@example.com'),
    'delivered',
    'UPS Next Day Air',
    'UPS-987654321',
    NOW() - INTERVAL '45 days',
    NOW() - INTERVAL '43 days',
    NOW() - INTERVAL '40 days',
    NOW() - INTERVAL '35 days',
    NOW() - INTERVAL '33 days',
    'High priority delivery completed.',
    'LTS-102',
    NOW() - INTERVAL '47 days',
    NOW() - INTERVAL '33 days'
);

-- Shipment 3: In Transit to Engineer
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'amazon web services'),
    (SELECT id FROM software_engineers WHERE email = 'david.chen@example.com'),
    'in_transit_to_engineer',
    'DHL Express',
    'DHL-456789123',
    NOW() - INTERVAL '10 days',
    NOW() - INTERVAL '8 days',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '2 days',
    'Out for delivery. Expected delivery today.',
    'LTS-103',
    NOW() - INTERVAL '12 days',
    NOW() - INTERVAL '1 day'
);

-- Shipment 4: At Warehouse (Pending assignment)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'meta platforms inc'),
    NULL,  -- No engineer assigned yet
    'at_warehouse',
    'FedEx Ground',
    'FDX-789456123',
    NOW() - INTERVAL '8 days',
    NOW() - INTERVAL '6 days',
    NOW() - INTERVAL '3 days',
    'Awaiting engineer assignment and address confirmation.',
    'LTS-104',
    NOW() - INTERVAL '10 days',
    NOW() - INTERVAL '3 days'
);

-- Shipment 5: In Transit to Warehouse
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'apple inc'),
    (SELECT id FROM software_engineers WHERE email = 'sarah.johnson@example.com'),
    'in_transit_to_warehouse',
    'UPS Ground',
    'UPS-321654987',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '3 days',
    'Expected arrival at warehouse tomorrow.',
    'LTS-105',
    NOW() - INTERVAL '7 days',
    NOW() - INTERVAL '1 day'
);

-- Shipment 6: Picked up from Client
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'google llc'),
    (SELECT id FROM software_engineers WHERE email = 'michael.brown@example.com'),
    'picked_up_from_client',
    'FedEx Priority',
    'FDX-654321789',
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '1 day',
    'Just picked up. Processing through courier network.',
    'LTS-106',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '1 day'
);

-- Shipment 7: Pending Pickup (Recently created)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name,
    pickup_scheduled_date,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'microsoft corporation'),
    (SELECT id FROM software_engineers WHERE email = 'linda.anderson@example.com'),
    'pending_pickup',
    'DHL Standard',
    NOW() + INTERVAL '2 days',
    'Pickup scheduled for next week.',
    'LTS-107',
    NOW() - INTERVAL '2 days',
    NOW()
);

-- Shipment 8: At Warehouse (Ready for release)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'amazon web services'),
    (SELECT id FROM software_engineers WHERE email = 'robert.taylor@example.com'),
    'at_warehouse',
    'FedEx Express',
    'FDX-111222333',
    NOW() - INTERVAL '15 days',
    NOW() - INTERVAL '13 days',
    NOW() - INTERVAL '10 days',
    'Waiting for engineer address confirmation.',
    'LTS-108',
    NOW() - INTERVAL '17 days',
    NOW() - INTERVAL '2 days'
);

-- Shipment 9: Delivered (Old shipment)
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'meta platforms inc'),
    (SELECT id FROM software_engineers WHERE email = 'william.thomas@example.com'),
    'delivered',
    'UPS Next Day',
    'UPS-444555666',
    NOW() - INTERVAL '60 days',
    NOW() - INTERVAL '58 days',
    NOW() - INTERVAL '55 days',
    NOW() - INTERVAL '52 days',
    NOW() - INTERVAL '50 days',
    'Completed delivery for new hire onboarding.',
    'LTS-109',
    NOW() - INTERVAL '62 days',
    NOW() - INTERVAL '50 days'
);

-- Shipment 10: Delivered
INSERT INTO shipments (
    client_company_id, software_engineer_id, status, courier_name, tracking_number,
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
    notes, jira_ticket_number, created_at, updated_at
)
VALUES (
    (SELECT id FROM client_companies WHERE LOWER(name) = 'apple inc'),
    (SELECT id FROM software_engineers WHERE email = 'chris.harris@example.com'),
    'delivered',
    'DHL Express',
    'DHL-777888999',
    NOW() - INTERVAL '35 days',
    NOW() - INTERVAL '33 days',
    NOW() - INTERVAL '30 days',
    NOW() - INTERVAL '27 days',
    NOW() - INTERVAL '25 days',
    'MacBook delivery completed.',
    'LTS-110',
    NOW() - INTERVAL '37 days',
    NOW() - INTERVAL '25 days'
);

\echo '>>> Shipments created.'

-- ============================================
-- 5. SHIPMENT-LAPTOP ASSOCIATIONS
-- ============================================
\echo ''
\echo '>>> Linking Laptops to Shipments...'

-- Link laptops to shipments using JIRA ticket numbers to find shipments
INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at)
VALUES 
    -- Shipment LTS-101 (Delivered): Dell XPS and Lenovo
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-101'), (SELECT id FROM laptops WHERE serial_number = 'DELL-SN-001'), NOW() - INTERVAL '52 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-101'), (SELECT id FROM laptops WHERE serial_number = 'LENOVO-SN-004'), NOW() - INTERVAL '52 days'),
    
    -- Shipment LTS-102 (Delivered): HP EliteBook and Apple MacBook
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-102'), (SELECT id FROM laptops WHERE serial_number = 'HP-SN-001'), NOW() - INTERVAL '47 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-102'), (SELECT id FROM laptops WHERE serial_number = 'APPLE-SN-001'), NOW() - INTERVAL '47 days'),
    
    -- Shipment LTS-103 (In Transit to Engineer): 2 laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-103'), (SELECT id FROM laptops WHERE serial_number = 'DELL-SN-004'), NOW() - INTERVAL '12 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-103'), (SELECT id FROM laptops WHERE serial_number = 'MSFT-SN-001'), NOW() - INTERVAL '12 days'),
    
    -- Shipment LTS-104 (At Warehouse): 3 laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-104'), (SELECT id FROM laptops WHERE serial_number = 'DELL-SN-003'), NOW() - INTERVAL '10 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-104'), (SELECT id FROM laptops WHERE serial_number = 'HP-SN-002'), NOW() - INTERVAL '10 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-104'), (SELECT id FROM laptops WHERE serial_number = 'LENOVO-SN-002'), NOW() - INTERVAL '10 days'),
    
    -- Shipment LTS-105 (In Transit to Warehouse): 2 laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-105'), (SELECT id FROM laptops WHERE serial_number = 'HP-SN-003'), NOW() - INTERVAL '7 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-105'), (SELECT id FROM laptops WHERE serial_number = 'LENOVO-SN-003'), NOW() - INTERVAL '7 days'),
    
    -- Shipment LTS-106 (Picked Up): 1 Apple MacBook
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-106'), (SELECT id FROM laptops WHERE serial_number = 'APPLE-SN-002'), NOW() - INTERVAL '5 days'),
    
    -- Shipment LTS-107 (Pending Pickup): 1 HP
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-107'), (SELECT id FROM laptops WHERE serial_number = 'HP-SN-004'), NOW() - INTERVAL '2 days'),
    
    -- Shipment LTS-108 (At Warehouse): 2 laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-108'), (SELECT id FROM laptops WHERE serial_number = 'ASUS-SN-002'), NOW() - INTERVAL '17 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-108'), (SELECT id FROM laptops WHERE serial_number = 'APPLE-SN-003'), NOW() - INTERVAL '17 days'),
    
    -- Shipment LTS-109 (Delivered): 2 laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-109'), (SELECT id FROM laptops WHERE serial_number = 'LENOVO-SN-001'), NOW() - INTERVAL '62 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-109'), (SELECT id FROM laptops WHERE serial_number = 'ASUS-SN-001'), NOW() - INTERVAL '62 days'),
    
    -- Shipment LTS-110 (Delivered): 2 Dell laptops
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-110'), (SELECT id FROM laptops WHERE serial_number = 'DELL-SN-002'), NOW() - INTERVAL '37 days'),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-110'), (SELECT id FROM laptops WHERE serial_number = 'DELL-SN-005'), NOW() - INTERVAL '37 days');

\echo '>>> Laptop-Shipment associations created.'

-- ============================================
-- 6. PICKUP FORMS
-- ============================================
\echo ''
\echo '>>> Creating Pickup Forms...'

-- Pickup forms using JIRA ticket numbers to reference shipments
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
VALUES (
    (SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-101'),
    (SELECT id FROM users WHERE email = 'client@bairesdev.com'),
    NOW() - INTERVAL '52 days',
    '{
        "company_name": "Google LLC",
        "contact_person": "Jane Doe",
        "contact_email": "jane.doe@google.com",
        "contact_phone": "+1-650-555-0001",
        "pickup_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043",
        "preferred_date": "2024-09-15",
        "num_laptops": 2,
        "special_instructions": "Please call 30 minutes before arrival"
    }'::jsonb
);

INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
VALUES (
    (SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-102'),
    (SELECT id FROM users WHERE email = 'logistics@bairesdev.com'),
    NOW() - INTERVAL '47 days',
    '{
        "company_name": "Microsoft Corporation",
        "contact_person": "Bob Smith",
        "contact_email": "bob.smith@microsoft.com",
        "contact_phone": "+1-425-555-0002",
        "pickup_address": "One Microsoft Way, Redmond, WA 98052",
        "preferred_date": "2024-09-20",
        "num_laptops": 2,
        "special_instructions": "Building 92, reception desk"
    }'::jsonb
);

-- More pickup forms for other shipments
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
VALUES 
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-103'), (SELECT id FROM users WHERE email = 'client@bairesdev.com'), NOW() - INTERVAL '12 days', 
     '{"company_name": "Amazon Web Services", "contact_person": "Alice Johnson", "contact_email": "alice@aws.com", "pickup_address": "410 Terry Ave N, Seattle, WA 98109", "num_laptops": 2}'::jsonb),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-104'), (SELECT id FROM users WHERE email = 'logistics@bairesdev.com'), NOW() - INTERVAL '10 days',
     '{"company_name": "Meta Platforms Inc", "contact_person": "Mark Davis", "contact_email": "mark@meta.com", "pickup_address": "1 Hacker Way, Menlo Park, CA 94025", "num_laptops": 3}'::jsonb),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-105'), (SELECT id FROM users WHERE email = 'client@bairesdev.com'), NOW() - INTERVAL '7 days',
     '{"company_name": "Apple Inc", "contact_person": "Steve Wilson", "contact_email": "steve@apple.com", "pickup_address": "One Apple Park Way, Cupertino, CA 95014", "num_laptops": 2}'::jsonb),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-106'), (SELECT id FROM users WHERE email = 'logistics@bairesdev.com'), NOW() - INTERVAL '5 days',
     '{"company_name": "Google LLC", "contact_person": "Larry Page", "contact_email": "larry@google.com", "pickup_address": "1600 Amphitheatre Parkway, Mountain View, CA", "num_laptops": 1}'::jsonb),
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-107'), (SELECT id FROM users WHERE email = 'client@bairesdev.com'), NOW() - INTERVAL '2 days',
     '{"company_name": "Microsoft Corporation", "contact_person": "Satya N", "contact_email": "satya@microsoft.com", "pickup_address": "One Microsoft Way, Redmond, WA", "num_laptops": 1}'::jsonb);

\echo '>>> Pickup forms created.'

-- ============================================
-- 7. RECEPTION REPORTS
-- ============================================
\echo ''
\echo '>>> Creating Reception Reports...'

-- Reception reports using JIRA ticket numbers
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls)
VALUES 
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-101'), (SELECT id FROM users WHERE email = 'warehouse@bairesdev.com'), NOW() - INTERVAL '45 days',
     'All items received in good condition. Dell XPS 15 and Lenovo verified. No visible damage.', 
     ARRAY['https://example.com/photos/shipment1_photo1.jpg', 'https://example.com/photos/shipment1_photo2.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-102'), (SELECT id FROM users WHERE email = 'warehouse@bairesdev.com'), NOW() - INTERVAL '40 days',
     'HP EliteBook and MacBook Pro received. Original packaging intact. Serial numbers verified.',
     ARRAY['https://example.com/photos/shipment2_photo1.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-103'), (SELECT id FROM users WHERE email = 'warehouse@bairesdev.com'), NOW() - INTERVAL '5 days',
     'Two laptops received and inspected. Ready for delivery assignment.', 
     ARRAY['https://example.com/photos/shipment3.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-104'), (SELECT id FROM users WHERE email = 'warehouse@bairesdev.com'), NOW() - INTERVAL '3 days',
     'Three laptops in excellent condition. Awaiting engineer assignment.', 
     ARRAY['https://example.com/photos/shipment4_1.jpg', 'https://example.com/photos/shipment4_2.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-108'), (SELECT id FROM users WHERE email = 'warehouse@bairesdev.com'), NOW() - INTERVAL '10 days',
     'Asus laptops received. Waiting for address confirmation from engineer.', 
     ARRAY[]::TEXT[]);

\echo '>>> Reception reports created.'

-- ============================================
-- 8. DELIVERY FORMS
-- ============================================
\echo ''
\echo '>>> Creating Delivery Forms...'

-- Delivery forms using JIRA ticket numbers
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls)
VALUES 
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-101'), (SELECT id FROM software_engineers WHERE email = 'john.smith@example.com'), NOW() - INTERVAL '38 days',
     'Laptops delivered and verified by engineer. Power adapters and all accessories included.',
     ARRAY['https://example.com/photos/delivery1_signature.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-102'), (SELECT id FROM software_engineers WHERE email = 'maria.garcia@example.com'), NOW() - INTERVAL '33 days',
     'Successful delivery. Engineer confirmed laptops are working properly.',
     ARRAY['https://example.com/photos/delivery2_signature.jpg', 'https://example.com/photos/delivery2_laptop.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-109'), (SELECT id FROM software_engineers WHERE email = 'william.thomas@example.com'), NOW() - INTERVAL '50 days',
     'Lenovo ThinkPad and Asus laptop delivered. Engineer satisfied with condition.', 
     ARRAY['https://example.com/photos/delivery9.jpg']),
    
    ((SELECT id FROM shipments WHERE jira_ticket_number = 'LTS-110'), (SELECT id FROM software_engineers WHERE email = 'chris.harris@example.com'), NOW() - INTERVAL '25 days',
     'Two Dell laptops delivered successfully. Setup completed on-site.', 
     ARRAY['https://example.com/photos/delivery10.jpg']);

\echo '>>> Delivery forms created.'

-- ============================================
-- COMMIT TRANSACTION
-- ============================================

COMMIT;

\echo ''
\echo '======================================='
\echo 'Test Data Created Successfully! ✓'
\echo '======================================='
\echo ''
\echo 'Summary:'
\echo '  - 5 Client Companies'
\echo '  - 15 Software Engineers'
\echo '  - 25 Laptops'
\echo '  - 10 Shipments (various statuses)'
\echo '  - 7 Pickup Forms'
\echo '  - 5 Reception Reports'
\echo '  - 4 Delivery Forms'
\echo ''
\echo 'Shipment Status Distribution:'
\echo '  - Delivered: 4 shipments'
\echo '  - At Warehouse: 2 shipments'
\echo '  - In Transit to Engineer: 1 shipment'
\echo '  - In Transit to Warehouse: 1 shipment'
\echo '  - Picked Up: 1 shipment'
\echo '  - Pending Pickup: 1 shipment'
\echo ''
\echo 'You can now test the application with realistic data!'
\echo ''

-- Display some sample data
\echo 'Recent Shipments:'
SELECT 
    s.id,
    cc.name as company,
    s.status,
    s.tracking_number,
    COUNT(sl.laptop_id) as laptop_count
FROM shipments s
JOIN client_companies cc ON cc.id = s.client_company_id
LEFT JOIN shipment_laptops sl ON sl.shipment_id = s.id
GROUP BY s.id, cc.name, s.status, s.tracking_number
ORDER BY s.created_at DESC
LIMIT 5;

\echo ''
