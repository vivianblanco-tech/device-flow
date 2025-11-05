-- =====================================================
-- Test Data for Laptop Tracking System
-- =====================================================
-- This script populates the following tables with test data:
-- 1. client_companies
-- 2. software_engineers
-- 3. laptops
-- 4. shipments
-- =====================================================

\echo ''
\echo '======================================='
\echo 'Creating Test Data for Laptop Tracking'
\echo '======================================='
\echo ''

-- =====================================================
-- 1. CLIENT COMPANIES
-- =====================================================

\echo 'Creating test client companies...'

INSERT INTO client_companies (name, contact_info, created_at, updated_at)
VALUES 
    (
        'TechCorp Solutions',
        'Contact: John Smith
Email: contact@techcorp-solutions.com
Phone: +1 (555) 123-4567
Address: 123 Tech Street, San Francisco, CA 94102',
        NOW() - INTERVAL '60 days',
        NOW() - INTERVAL '60 days'
    ),
    (
        'Global Innovations Inc',
        'Contact: Sarah Johnson
Email: info@globalinnovations.com
Phone: +1 (555) 234-5678
Address: 456 Innovation Blvd, Austin, TX 78701',
        NOW() - INTERVAL '45 days',
        NOW() - INTERVAL '45 days'
    ),
    (
        'Digital Dynamics LLC',
        'Contact: Michael Chen
Email: support@digitaldynamics.com
Phone: +1 (555) 345-6789
Address: 789 Digital Ave, Seattle, WA 98101',
        NOW() - INTERVAL '30 days',
        NOW() - INTERVAL '30 days'
    ),
    (
        'CloudFirst Technologies',
        'Contact: Emma Williams
Email: hello@cloudfirst.tech
Phone: +1 (555) 456-7890
Address: 321 Cloud Way, Boston, MA 02101',
        NOW() - INTERVAL '20 days',
        NOW() - INTERVAL '20 days'
    ),
    (
        'NextGen Software Group',
        'Contact: David Rodriguez
Email: contact@nextgensoftware.com
Phone: +1 (555) 567-8901
Address: 555 Future Pkwy, Denver, CO 80202',
        NOW() - INTERVAL '15 days',
        NOW() - INTERVAL '15 days'
    )
ON CONFLICT ((LOWER(name))) DO UPDATE
SET 
    contact_info = EXCLUDED.contact_info,
    updated_at = EXCLUDED.updated_at;

\echo 'Created 5 client companies'

-- =====================================================
-- 2. SOFTWARE ENGINEERS
-- =====================================================

\echo 'Creating test software engineers...'

INSERT INTO software_engineers (name, email, address, phone, address_confirmed, address_confirmation_at, created_at, updated_at)
VALUES 
    (
        'Alex Thompson',
        'alex.thompson@email.com',
        '100 Main St, Apt 4B, New York, NY 10001',
        '+1 (555) 111-2222',
        TRUE,
        NOW() - INTERVAL '25 days',
        NOW() - INTERVAL '30 days',
        NOW() - INTERVAL '25 days'
    ),
    (
        'Maria Garcia',
        'maria.garcia@email.com',
        '250 Oak Avenue, Unit 12, Los Angeles, CA 90001',
        '+1 (555) 222-3333',
        TRUE,
        NOW() - INTERVAL '20 days',
        NOW() - INTERVAL '28 days',
        NOW() - INTERVAL '20 days'
    ),
    (
        'James Wilson',
        'james.wilson@email.com',
        '789 Pine Street, Chicago, IL 60601',
        '+1 (555) 333-4444',
        TRUE,
        NOW() - INTERVAL '15 days',
        NOW() - INTERVAL '25 days',
        NOW() - INTERVAL '15 days'
    ),
    (
        'Emily Chen',
        'emily.chen@email.com',
        '456 Elm Drive, Miami, FL 33101',
        '+1 (555) 444-5555',
        TRUE,
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '20 days',
        NOW() - INTERVAL '10 days'
    ),
    (
        'Robert Martinez',
        'robert.martinez@email.com',
        '321 Maple Road, Phoenix, AZ 85001',
        '+1 (555) 555-6666',
        TRUE,
        NOW() - INTERVAL '8 days',
        NOW() - INTERVAL '18 days',
        NOW() - INTERVAL '8 days'
    ),
    (
        'Sarah Anderson',
        'sarah.anderson@email.com',
        '654 Cedar Lane, Philadelphia, PA 19101',
        '+1 (555) 666-7777',
        TRUE,
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '15 days',
        NOW() - INTERVAL '5 days'
    ),
    (
        'David Kim',
        'david.kim@email.com',
        '987 Birch Court, San Diego, CA 92101',
        '+1 (555) 777-8888',
        TRUE,
        NOW() - INTERVAL '3 days',
        NOW() - INTERVAL '12 days',
        NOW() - INTERVAL '3 days'
    ),
    (
        'Jessica Taylor',
        'jessica.taylor@email.com',
        '147 Willow Way, Dallas, TX 75201',
        '+1 (555) 888-9999',
        FALSE,
        NULL,
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '10 days'
    ),
    (
        'Michael Brown',
        'michael.brown@email.com',
        '258 Spruce Street, San Jose, CA 95101',
        '+1 (555) 999-0000',
        FALSE,
        NULL,
        NOW() - INTERVAL '7 days',
        NOW() - INTERVAL '7 days'
    ),
    (
        'Lisa Johnson',
        'lisa.johnson@email.com',
        '369 Redwood Avenue, Austin, TX 78701',
        '+1 (555) 000-1111',
        FALSE,
        NULL,
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '5 days'
    )
ON CONFLICT ((LOWER(email))) DO UPDATE
SET 
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    phone = EXCLUDED.phone,
    address_confirmed = EXCLUDED.address_confirmed,
    address_confirmation_at = EXCLUDED.address_confirmation_at,
    updated_at = EXCLUDED.updated_at;

\echo 'Created 10 software engineers'

-- =====================================================
-- 3. LAPTOPS
-- =====================================================

\echo 'Creating test laptops...'

INSERT INTO laptops (serial_number, brand, model, specs, status, created_at, updated_at)
VALUES 
    -- Dell Laptops
    (
        'DELL-XPS13-SN001',
        'Dell',
        'XPS 13 9310',
        'Intel Core i7-1165G7, 16GB RAM, 512GB SSD, 13.3" FHD Display',
        'delivered',
        NOW() - INTERVAL '25 days',
        NOW() - INTERVAL '5 days'
    ),
    (
        'DELL-XPS15-SN002',
        'Dell',
        'XPS 15 9520',
        'Intel Core i7-12700H, 32GB RAM, 1TB SSD, 15.6" 4K Display, NVIDIA RTX 3050 Ti',
        'at_warehouse',
        NOW() - INTERVAL '20 days',
        NOW() - INTERVAL '8 days'
    ),
    (
        'DELL-LAT7420-SN003',
        'Dell',
        'Latitude 7420',
        'Intel Core i5-1145G7, 16GB RAM, 256GB SSD, 14" FHD Display',
        'delivered',
        NOW() - INTERVAL '22 days',
        NOW() - INTERVAL '3 days'
    ),
    -- Lenovo Laptops
    (
        'LNVO-X1C9-SN004',
        'Lenovo',
        'ThinkPad X1 Carbon Gen 9',
        'Intel Core i7-1185G7, 16GB RAM, 512GB SSD, 14" FHD Display',
        'in_transit_to_engineer',
        NOW() - INTERVAL '18 days',
        NOW() - INTERVAL '2 days'
    ),
    (
        'LNVO-P1G4-SN005',
        'Lenovo',
        'ThinkPad P1 Gen 4',
        'Intel Core i9-11950H, 32GB RAM, 1TB SSD, 16" 4K Display, NVIDIA RTX A2000',
        'delivered',
        NOW() - INTERVAL '15 days',
        NOW() - INTERVAL '1 day'
    ),
    (
        'LNVO-T14G2-SN006',
        'Lenovo',
        'ThinkPad T14 Gen 2',
        'AMD Ryzen 7 PRO 5850U, 16GB RAM, 512GB SSD, 14" FHD Display',
        'at_warehouse',
        NOW() - INTERVAL '12 days',
        NOW() - INTERVAL '6 days'
    ),
    -- HP Laptops
    (
        'HP-ELITE840-SN007',
        'HP',
        'EliteBook 840 G8',
        'Intel Core i7-1185G7, 16GB RAM, 512GB SSD, 14" FHD Display',
        'delivered',
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '2 days'
    ),
    (
        'HP-ZB15G8-SN008',
        'HP',
        'ZBook 15 G8',
        'Intel Core i9-11950H, 32GB RAM, 1TB SSD, 15.6" 4K Display, NVIDIA RTX A3000',
        'in_transit_to_warehouse',
        NOW() - INTERVAL '8 days',
        NOW() - INTERVAL '1 day'
    ),
    -- MacBook Laptops
    (
        'APPLE-MBP14-SN009',
        'Apple',
        'MacBook Pro 14" 2021',
        'Apple M1 Pro, 16GB RAM, 512GB SSD, 14.2" Liquid Retina XDR Display',
        'delivered',
        NOW() - INTERVAL '6 days',
        NOW() - INTERVAL '1 day'
    ),
    (
        'APPLE-MBP16-SN010',
        'Apple',
        'MacBook Pro 16" 2021',
        'Apple M1 Max, 32GB RAM, 1TB SSD, 16.2" Liquid Retina XDR Display',
        'at_warehouse',
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '4 days'
    ),
    -- ASUS Laptops
    (
        'ASUS-ZEN14-SN011',
        'ASUS',
        'ZenBook 14 UX425',
        'Intel Core i7-1165G7, 16GB RAM, 512GB SSD, 14" FHD Display',
        'available',
        NOW() - INTERVAL '4 days',
        NOW() - INTERVAL '4 days'
    ),
    (
        'ASUS-ROG15-SN012',
        'ASUS',
        'ROG Zephyrus G15',
        'AMD Ryzen 9 5900HS, 32GB RAM, 1TB SSD, 15.6" QHD Display, NVIDIA RTX 3070',
        'available',
        NOW() - INTERVAL '3 days',
        NOW() - INTERVAL '3 days'
    ),
    -- Microsoft Surface
    (
        'MSFT-SL4-SN013',
        'Microsoft',
        'Surface Laptop 4',
        'Intel Core i7-1185G7, 16GB RAM, 512GB SSD, 13.5" PixelSense Display',
        'in_transit_to_engineer',
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '1 day'
    ),
    -- Additional Dell laptops
    (
        'DELL-PRE5550-SN014',
        'Dell',
        'Precision 5550',
        'Intel Core i7-10875H, 32GB RAM, 1TB SSD, 15.6" 4K Display, NVIDIA Quadro T2000',
        'in_transit_to_warehouse',
        NOW() - INTERVAL '7 days',
        NOW() - INTERVAL '2 days'
    ),
    (
        'DELL-INS7590-SN015',
        'Dell',
        'Inspiron 15 7590',
        'Intel Core i7-9750H, 16GB RAM, 512GB SSD, 15.6" FHD Display, NVIDIA GTX 1650',
        'available',
        NOW() - INTERVAL '9 days',
        NOW() - INTERVAL '9 days'
    )
ON CONFLICT (serial_number) DO UPDATE
SET 
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    specs = EXCLUDED.specs,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at;

\echo 'Created 15 laptops'

-- =====================================================
-- 4. SHIPMENTS
-- =====================================================

\echo 'Creating test shipments...'

-- Get IDs for reference
DO $$
DECLARE
    techcorp_id BIGINT;
    global_id BIGINT;
    digital_id BIGINT;
    cloudfirst_id BIGINT;
    nextgen_id BIGINT;
    
    alex_id BIGINT;
    maria_id BIGINT;
    james_id BIGINT;
    emily_id BIGINT;
    robert_id BIGINT;
    sarah_id BIGINT;
    david_id BIGINT;
    
    dell_xps13_id BIGINT;
    dell_lat_id BIGINT;
    lnvo_x1c_id BIGINT;
    lnvo_p1_id BIGINT;
    hp_elite_id BIGINT;
    apple_mbp14_id BIGINT;
    dell_xps15_id BIGINT;
    lnvo_t14_id BIGINT;
    apple_mbp16_id BIGINT;
    hp_zbook_id BIGINT;
    msft_sl4_id BIGINT;
    dell_pre_id BIGINT;
BEGIN
    -- Get client company IDs
    SELECT id INTO techcorp_id FROM client_companies WHERE LOWER(name) = 'techcorp solutions';
    SELECT id INTO global_id FROM client_companies WHERE LOWER(name) = 'global innovations inc';
    SELECT id INTO digital_id FROM client_companies WHERE LOWER(name) = 'digital dynamics llc';
    SELECT id INTO cloudfirst_id FROM client_companies WHERE LOWER(name) = 'cloudfirst technologies';
    SELECT id INTO nextgen_id FROM client_companies WHERE LOWER(name) = 'nextgen software group';
    
    -- Get software engineer IDs
    SELECT id INTO alex_id FROM software_engineers WHERE LOWER(email) = 'alex.thompson@email.com';
    SELECT id INTO maria_id FROM software_engineers WHERE LOWER(email) = 'maria.garcia@email.com';
    SELECT id INTO james_id FROM software_engineers WHERE LOWER(email) = 'james.wilson@email.com';
    SELECT id INTO emily_id FROM software_engineers WHERE LOWER(email) = 'emily.chen@email.com';
    SELECT id INTO robert_id FROM software_engineers WHERE LOWER(email) = 'robert.martinez@email.com';
    SELECT id INTO sarah_id FROM software_engineers WHERE LOWER(email) = 'sarah.anderson@email.com';
    SELECT id INTO david_id FROM software_engineers WHERE LOWER(email) = 'david.kim@email.com';
    
    -- Get laptop IDs
    SELECT id INTO dell_xps13_id FROM laptops WHERE serial_number = 'DELL-XPS13-SN001';
    SELECT id INTO dell_lat_id FROM laptops WHERE serial_number = 'DELL-LAT7420-SN003';
    SELECT id INTO lnvo_x1c_id FROM laptops WHERE serial_number = 'LNVO-X1C9-SN004';
    SELECT id INTO lnvo_p1_id FROM laptops WHERE serial_number = 'LNVO-P1G4-SN005';
    SELECT id INTO hp_elite_id FROM laptops WHERE serial_number = 'HP-ELITE840-SN007';
    SELECT id INTO apple_mbp14_id FROM laptops WHERE serial_number = 'APPLE-MBP14-SN009';
    SELECT id INTO dell_xps15_id FROM laptops WHERE serial_number = 'DELL-XPS15-SN002';
    SELECT id INTO lnvo_t14_id FROM laptops WHERE serial_number = 'LNVO-T14G2-SN006';
    SELECT id INTO apple_mbp16_id FROM laptops WHERE serial_number = 'APPLE-MBP16-SN010';
    SELECT id INTO hp_zbook_id FROM laptops WHERE serial_number = 'HP-ZB15G8-SN008';
    SELECT id INTO msft_sl4_id FROM laptops WHERE serial_number = 'MSFT-SL4-SN013';
    SELECT id INTO dell_pre_id FROM laptops WHERE serial_number = 'DELL-PRE5550-SN014';
    
    -- Shipment 1: Delivered (Dell XPS 13 to Alex)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        techcorp_id, alex_id, 'delivered', 'FedEx', 'FDX1234567890',
        NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days', NOW() - INTERVAL '22 days',
        NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days',
        'Standard delivery, all documentation verified',
        NOW() - INTERVAL '26 days', NOW() - INTERVAL '18 days'
    );
    
    -- Shipment 2: Delivered (Dell Latitude to Maria)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        global_id, maria_id, 'delivered', 'UPS', 'UPS9876543210',
        NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days', NOW() - INTERVAL '19 days',
        NOW() - INTERVAL '17 days', NOW() - INTERVAL '15 days',
        'Express delivery requested by client',
        NOW() - INTERVAL '23 days', NOW() - INTERVAL '15 days'
    );
    
    -- Shipment 3: In Transit to Engineer (Lenovo X1 Carbon to James)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        digital_id, james_id, 'in_transit_to_engineer', 'DHL', 'DHL5678901234',
        NOW() - INTERVAL '18 days', NOW() - INTERVAL '17 days', NOW() - INTERVAL '15 days',
        NOW() - INTERVAL '2 days', NULL,
        'Engineer confirmed address, shipment dispatched',
        NOW() - INTERVAL '19 days', NOW() - INTERVAL '2 days'
    );
    
    -- Shipment 4: Delivered (Lenovo P1 to Emily)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        cloudfirst_id, emily_id, 'delivered', 'FedEx', 'FDX2345678901',
        NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '12 days',
        NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days',
        'High-value shipment, signature required',
        NOW() - INTERVAL '16 days', NOW() - INTERVAL '8 days'
    );
    
    -- Shipment 5: Delivered (HP EliteBook to Robert)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        nextgen_id, robert_id, 'delivered', 'UPS', 'UPS1234567890',
        NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '7 days',
        NOW() - INTERVAL '5 days', NOW() - INTERVAL '3 days',
        'Standard processing',
        NOW() - INTERVAL '11 days', NOW() - INTERVAL '3 days'
    );
    
    -- Shipment 6: Delivered (MacBook Pro 14 to Sarah)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        techcorp_id, sarah_id, 'delivered', 'FedEx', 'FDX3456789012',
        NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days',
        NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day',
        'Premium device, extra care taken during shipping',
        NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day'
    );
    
    -- Shipment 7: At Warehouse (Dell XPS 15)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        global_id, NULL, 'at_warehouse', 'DHL', 'DHL6789012345',
        NOW() - INTERVAL '8 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days',
        NULL, NULL,
        'Awaiting engineer assignment',
        NOW() - INTERVAL '9 days', NOW() - INTERVAL '5 days'
    );
    
    -- Shipment 8: At Warehouse (Lenovo T14)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        digital_id, NULL, 'at_warehouse', 'UPS', 'UPS2345678901',
        NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days',
        NULL, NULL,
        'Standard laptop, ready for assignment',
        NOW() - INTERVAL '7 days', NOW() - INTERVAL '4 days'
    );
    
    -- Shipment 9: At Warehouse (MacBook Pro 16)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        cloudfirst_id, NULL, 'at_warehouse', 'FedEx', 'FDX4567890123',
        NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days',
        NULL, NULL,
        'High-spec device, needs careful handling',
        NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days'
    );
    
    -- Shipment 10: In Transit to Warehouse (HP ZBook)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        nextgen_id, NULL, 'in_transit_to_warehouse', 'DHL', 'DHL7890123456',
        NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', NULL,
        NULL, NULL,
        'En route to warehouse, ETA 2 days',
        NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day'
    );
    
    -- Shipment 11: In Transit to Engineer (Surface Laptop to David)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        techcorp_id, david_id, 'in_transit_to_engineer', 'FedEx', 'FDX5678901234',
        NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '1 day', NULL,
        'Priority shipment, delivery expected tomorrow',
        NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day'
    );
    
    -- Shipment 12: Pending Pickup from Client
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        global_id, NULL, 'pending_pickup_from_client', 'UPS', NULL,
        NOW() + INTERVAL '2 days', NULL, NULL,
        NULL, NULL,
        'Pickup scheduled for future date',
        NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'
    );
    
    -- Shipment 13: Picked Up from Client (Dell Precision)
    INSERT INTO shipments (
        client_company_id, software_engineer_id, status, courier_name, tracking_number,
        pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
        notes, created_at, updated_at
    ) VALUES (
        digital_id, NULL, 'picked_up_from_client', 'DHL', 'DHL8901234567',
        NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', NULL,
        NULL, NULL,
        'Picked up from client, in transit',
        NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day'
    );
    
END $$;

\echo 'Created 13 shipments with various statuses'

-- =====================================================
-- SHIPMENT-LAPTOP JUNCTION TABLE
-- =====================================================

\echo 'Linking laptops to shipments...'

-- Link laptops to shipments based on serial numbers and shipment details
INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at)
SELECT s.id, l.id, s.created_at
FROM shipments s
JOIN laptops l ON (
    (s.status = 'delivered' AND l.status = 'delivered' AND l.serial_number = 'DELL-XPS13-SN001' AND s.courier_name = 'FedEx' AND s.tracking_number = 'FDX1234567890')
    OR (s.status = 'delivered' AND l.status = 'delivered' AND l.serial_number = 'DELL-LAT7420-SN003' AND s.courier_name = 'UPS' AND s.tracking_number = 'UPS9876543210')
    OR (s.status = 'in_transit_to_engineer' AND l.status = 'in_transit_to_engineer' AND l.serial_number = 'LNVO-X1C9-SN004')
    OR (s.status = 'delivered' AND l.status = 'delivered' AND l.serial_number = 'LNVO-P1G4-SN005')
    OR (s.status = 'delivered' AND l.status = 'delivered' AND l.serial_number = 'HP-ELITE840-SN007')
    OR (s.status = 'delivered' AND l.status = 'delivered' AND l.serial_number = 'APPLE-MBP14-SN009')
    OR (s.status = 'at_warehouse' AND l.status = 'at_warehouse' AND l.serial_number = 'DELL-XPS15-SN002')
    OR (s.status = 'at_warehouse' AND l.status = 'at_warehouse' AND l.serial_number = 'LNVO-T14G2-SN006')
    OR (s.status = 'at_warehouse' AND l.status = 'at_warehouse' AND l.serial_number = 'APPLE-MBP16-SN010')
    OR (s.status = 'in_transit_to_warehouse' AND l.status = 'in_transit_to_warehouse' AND l.serial_number = 'HP-ZB15G8-SN008')
    OR (s.status = 'in_transit_to_engineer' AND l.status = 'in_transit_to_engineer' AND l.serial_number = 'MSFT-SL4-SN013')
    OR (s.status = 'picked_up_from_client' AND l.status = 'in_transit_to_warehouse' AND l.serial_number = 'DELL-PRE5550-SN014')
)
ON CONFLICT (shipment_id, laptop_id) DO NOTHING;

\echo 'Linked laptops to shipments'

-- =====================================================
-- DISPLAY SUMMARY
-- =====================================================

\echo ''
\echo '======================================='
\echo 'Test Data Creation Complete!'
\echo '======================================='
\echo ''

-- Display summary statistics
SELECT 'Client Companies' as entity, COUNT(*) as count FROM client_companies
UNION ALL
SELECT 'Software Engineers', COUNT(*) FROM software_engineers
UNION ALL
SELECT 'Laptops', COUNT(*) FROM laptops
UNION ALL
SELECT 'Shipments', COUNT(*) FROM shipments
UNION ALL
SELECT 'Shipment-Laptop Links', COUNT(*) FROM shipment_laptops;

\echo ''
\echo '--- Shipment Status Breakdown ---'
SELECT status, COUNT(*) as count
FROM shipments
GROUP BY status
ORDER BY count DESC;

\echo ''
\echo '--- Laptop Status Breakdown ---'
SELECT status, COUNT(*) as count
FROM laptops
GROUP BY status
ORDER BY count DESC;

\echo ''
\echo '--- Laptops by Brand ---'
SELECT brand, COUNT(*) as count
FROM laptops
GROUP BY brand
ORDER BY count DESC;

\echo ''
\echo '======================================='
\echo 'You can now test the application!'
\echo '======================================='
\echo ''

