-- =============================================
-- COMPREHENSIVE SHIPMENTS & FORMS DATA
-- Align - Production-Ready Shipments
-- =============================================
-- This file adds 40+ shipments with complete forms and reports
-- to complement the enhanced-sample-data-comprehensive.sql base data
-- Password for all users: "Test123!"
-- Last Updated: 2025-11-13
-- =============================================

-- =============================================
-- SHIPMENTS (40+ covering all types and statuses)
-- =============================================

-- DELIVERED SHIPMENTS (Historical - 10 shipments)
-- Shipment 1: Single laptop - Delivered 60 days ago
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 1, 1, 'delivered', 1, 'SCOP-90001', 'FedEx Express', 'FDX9001234567', 
    (NOW() - INTERVAL '65 days')::date, NOW() - INTERVAL '63 days', NOW() - INTERVAL '60 days', 
    NOW() - INTERVAL '57 days', NOW() - INTERVAL '55 days', 
    'Dell Precision 5570 delivered to Alice Johnson. Complete lifecycle. Excellent condition throughout.', 
    NOW() - INTERVAL '67 days', NOW() - INTERVAL '55 days');

-- Shipment 2: Bulk - 5 laptops delivered 50 days ago
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 2, 5, 'delivered', 1, 'SCOP-90002', 'UPS Next Day Air', 'UPS9002345678', 
    (NOW() - INTERVAL '55 days')::date, NOW() - INTERVAL '53 days', NOW() - INTERVAL '50 days', 
    NOW() - INTERVAL '47 days', NOW() - INTERVAL '45 days', 
    'HP ZBook Studio G9 delivered to Emily Rodriguez. Video editing workstation performing excellently.', 
    NOW() - INTERVAL '57 days', NOW() - INTERVAL '45 days');

-- Shipment 3: Single - 45 days ago
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 3, 8, 'delivered', 1, 'SCOP-90003', 'DHL Express', 'DHL9003456789', 
    (NOW() - INTERVAL '50 days')::date, NOW() - INTERVAL '48 days', NOW() - INTERVAL '45 days', 
    NOW() - INTERVAL '42 days', NOW() - INTERVAL '40 days', 
    'Lenovo ThinkPad X1 Carbon Gen 10 delivered to Henry Thompson. Premium ultrabook for executive use.', 
    NOW() - INTERVAL '52 days', NOW() - INTERVAL '40 days');

-- Shipment 4: Single - 40 days ago
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 4, 11, 'delivered', 1, 'SCOP-90004', 'FedEx Priority', 'FDX9004567890', 
    (NOW() - INTERVAL '45 days')::date, NOW() - INTERVAL '43 days', NOW() - INTERVAL '40 days', 
    NOW() - INTERVAL '37 days', NOW() - INTERVAL '35 days', 
    'Apple MacBook Pro 16" M2 Max delivered to Karen Lee. High-end development machine.', 
    NOW() - INTERVAL '47 days', NOW() - INTERVAL '35 days');

-- Shipment 5: Single - 35 days ago
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES ('single_full_journey', 5, 14, 'delivered', 1, 'SCOP-90005', 'UPS Ground', 'UPS9005678901', 
    (NOW() - INTERVAL '40 days')::date, NOW() - INTERVAL '38 days', NOW() - INTERVAL '35 days', 
    NOW() - INTERVAL '32 days', NOW() - INTERVAL '30 days', 
    'Lenovo ThinkPad P1 Gen 5 delivered to Nathan Brown. GPU workstation for ML/AI development.', 
    NOW() - INTERVAL '42 days', NOW() - INTERVAL '30 days');

-- Shipments 6-10: More delivered (various dates)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 6, 17, 'delivered', 1, 'SCOP-90006', 'DHL Express', 'DHL9006789012', 
    (NOW() - INTERVAL '38 days')::date, NOW() - INTERVAL '36 days', NOW() - INTERVAL '33 days', 
    NOW() - INTERVAL '30 days', NOW() - INTERVAL '28 days', 
    'Dell XPS 15 9520 delivered to Samuel Taylor. Premium developer laptop.', 
    NOW() - INTERVAL '40 days', NOW() - INTERVAL '28 days'),
('single_full_journey', 7, 19, 'delivered', 1, 'SCOP-90007', 'FedEx Express', 'FDX9007890123', 
    (NOW() - INTERVAL '35 days')::date, NOW() - INTERVAL '33 days', NOW() - INTERVAL '30 days', 
    NOW() - INTERVAL '27 days', NOW() - INTERVAL '25 days', 
    'HP EliteBook 850 G9 delivered to Victor Harris. Business laptop for enterprise use.', 
    NOW() - INTERVAL '37 days', NOW() - INTERVAL '25 days'),
('single_full_journey', 8, 22, 'delivered', 1, 'SCOP-90008', 'UPS Next Day Air', 'UPS9008901234', 
    (NOW() - INTERVAL '32 days')::date, NOW() - INTERVAL '30 days', NOW() - INTERVAL '27 days', 
    NOW() - INTERVAL '24 days', NOW() - INTERVAL '22 days', 
    'Microsoft Surface Laptop Studio delivered successfully.', 
    NOW() - INTERVAL '34 days', NOW() - INTERVAL '22 days'),
('single_full_journey', 1, 3, 'delivered', 1, 'SCOP-90009', 'DHL Express', 'DHL9009012345', 
    (NOW() - INTERVAL '30 days')::date, NOW() - INTERVAL '28 days', NOW() - INTERVAL '25 days', 
    NOW() - INTERVAL '22 days', NOW() - INTERVAL '20 days', 
    'ASUS ZenBook Pro 15 OLED delivered to Catherine Wong. Creative workstation.', 
    NOW() - INTERVAL '32 days', NOW() - INTERVAL '20 days'),
('single_full_journey', 2, 6, 'delivered', 1, 'SCOP-90010', 'FedEx Ground', 'FDX9010123456', 
    (NOW() - INTERVAL '28 days')::date, NOW() - INTERVAL '26 days', NOW() - INTERVAL '23 days', 
    NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days', 
    'Lenovo ThinkPad X1 Carbon delivered to Frank Martinez.', 
    NOW() - INTERVAL '30 days', NOW() - INTERVAL '18 days');

-- IN TRANSIT TO ENGINEER (5 shipments - arriving soon)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, eta_to_engineer, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 3, 9, 'in_transit_to_engineer', 1, 'SCOP-90011', 'UPS Next Day Air', 'UPS9011234567', 
    (NOW() - INTERVAL '10 days')::date, NOW() - INTERVAL '8 days', NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '2 days', (NOW() + INTERVAL '1 day')::timestamp, 
    'Dell Precision 7670 workstation en route to James Wilson. High-end GPU machine. ETA: Tomorrow.', 
    NOW() - INTERVAL '12 days', NOW() - INTERVAL '1 day'),
('single_full_journey', 4, 12, 'in_transit_to_engineer', 1, 'SCOP-90012', 'FedEx Priority Overnight', 'FDX9012345678', 
    (NOW() - INTERVAL '9 days')::date, NOW() - INTERVAL '7 days', NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '1 day', (NOW() + INTERVAL '12 hours')::timestamp, 
    'Apple MacBook Pro 14" M2 Pro. Expected delivery today afternoon.', 
    NOW() - INTERVAL '11 days', NOW() - INTERVAL '6 hours'),
('single_full_journey', 5, 15, 'in_transit_to_engineer', 1, 'SCOP-90013', 'DHL Express', 'DHL9013456789', 
    (NOW() - INTERVAL '8 days')::date, NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '18 hours', (NOW() + INTERVAL '2 days')::timestamp, 
    'HP ZBook Fury G9 in transit. Professional workstation.', 
    NOW() - INTERVAL '10 days', NOW() - INTERVAL '18 hours'),
('single_full_journey', 6, 18, 'in_transit_to_engineer', 1, 'SCOP-90014', 'UPS Ground', 'UPS9014567890', 
    (NOW() - INTERVAL '11 days')::date, NOW() - INTERVAL '9 days', NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '3 days', (NOW() + INTERVAL '1 day')::timestamp, 
    'Lenovo ThinkPad P16 Gen 1. Arriving tomorrow.', 
    NOW() - INTERVAL '13 days', NOW() - INTERVAL '2 days'),
('single_full_journey', 7, 20, 'in_transit_to_engineer', 1, 'SCOP-90015', 'FedEx Express', 'FDX9015678901', 
    (NOW() - INTERVAL '7 days')::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '12 hours', (NOW() + INTERVAL '3 days')::timestamp, 
    'Dell XPS 13 Plus 9315. Standard delivery.', 
    NOW() - INTERVAL '9 days', NOW() - INTERVAL '12 hours');

-- RELEASED FROM WAREHOUSE (3 shipments - ready for courier pickup)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 8, 21, 'released_from_warehouse', 1, 'SCOP-90016', 'UPS Ground', 'UPS9016789012', 
    (NOW() - INTERVAL '8 days')::date, NOW() - INTERVAL '6 days', NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '6 hours', 
    'HP EliteBook 840 G9 packaged and ready. Courier pickup scheduled for today 2:00 PM.', 
    NOW() - INTERVAL '10 days', NOW() - INTERVAL '6 hours'),
('single_full_journey', 1, 4, 'released_from_warehouse', 1, 'SCOP-90017', 'FedEx Express', 'FDX9017890123', 
    (NOW() - INTERVAL '7 days')::date, NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '4 hours', 
    'Apple MacBook Air M2 sealed and ready for pickup.', 
    NOW() - INTERVAL '9 days', NOW() - INTERVAL '4 hours'),
('single_full_journey', 2, 7, 'released_from_warehouse', 1, 'SCOP-90018', 'DHL Express', 'DHL9018901234', 
    (NOW() - INTERVAL '6 days')::date, NOW() - INTERVAL '4 days', NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '3 hours', 
    'Microsoft Surface Laptop 5 ready. Priority shipment.', 
    NOW() - INTERVAL '8 days', NOW() - INTERVAL '3 hours');

-- AT WAREHOUSE (8 shipments - awaiting assignment or inspection)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, arrived_warehouse_at, notes, created_at, updated_at) 
VALUES 
('bulk_to_warehouse', 3, NULL, 'at_warehouse', 4, 'SCOP-90019', 'UPS Freight', 'UPS9019012345', 
    (NOW() - INTERVAL '5 days')::date, NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day', 
    'BULK: 4 Lenovo ThinkPad X1 Carbon laptops received. All units inspected and logged. Awaiting engineer assignments.', 
    NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day'),
('bulk_to_warehouse', 4, NULL, 'at_warehouse', 3, 'SCOP-90020', 'FedEx Freight', 'FDX9020123456', 
    (NOW() - INTERVAL '4 days')::date, NOW() - INTERVAL '2 days', NOW() - INTERVAL '18 hours', 
    'BULK: 3 Dell XPS 15 laptops. All in excellent condition. Ready for distribution.', 
    NOW() - INTERVAL '6 days', NOW() - INTERVAL '18 hours'),
('single_full_journey', 5, NULL, 'at_warehouse', 1, 'SCOP-90021', 'DHL Express', 'DHL9021234567', 
    (NOW() - INTERVAL '6 days')::date, NOW() - INTERVAL '4 days', NOW() - INTERVAL '2 days', 
    'HP ZBook Studio G9 received. Pending engineer assignment.', 
    NOW() - INTERVAL '8 days', NOW() - INTERVAL '2 days'),
('single_full_journey', 6, NULL, 'at_warehouse', 1, 'SCOP-90022', 'UPS Ground', 'UPS9022345678', 
    (NOW() - INTERVAL '5 days')::date, NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day', 
    'Acer Swift X received in good condition.', 
    NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day'),
('bulk_to_warehouse', 7, NULL, 'at_warehouse', 5, 'SCOP-90023', 'FedEx Freight', 'FDX9023456789', 
    (NOW() - INTERVAL '3 days')::date, NOW() - INTERVAL '1 day', NOW() - INTERVAL '12 hours', 
    'BULK: 5 HP EliteBook 840 G9 laptops for department expansion. All units tested.', 
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '12 hours'),
('single_full_journey', 8, NULL, 'at_warehouse', 1, 'SCOP-90024', 'DHL Express', 'DHL9024567890', 
    (NOW() - INTERVAL '4 days')::date, NOW() - INTERVAL '2 days', NOW() - INTERVAL '18 hours', 
    'ASUS ROG Zephyrus G14 gaming laptop received.', 
    NOW() - INTERVAL '6 days', NOW() - INTERVAL '18 hours'),
('bulk_to_warehouse', 1, NULL, 'at_warehouse', 6, 'SCOP-90025', 'UPS Freight', 'UPS9025678901', 
    (NOW() - INTERVAL '2 days')::date, NOW() - INTERVAL '18 hours', NOW() - INTERVAL '6 hours', 
    'BULK: 6 Apple MacBook Pro units for iOS development team. High-value shipment.', 
    NOW() - INTERVAL '4 days', NOW() - INTERVAL '6 hours'),
('single_full_journey', 2, NULL, 'at_warehouse', 1, 'SCOP-90026', 'FedEx Ground', 'FDX9026789012', 
    (NOW() - INTERVAL '3 days')::date, NOW() - INTERVAL '1 day', NOW() - INTERVAL '8 hours', 
    'Lenovo ThinkPad P1 Gen 5 workstation received.', 
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '8 hours');

-- IN TRANSIT TO WAREHOUSE (5 shipments - on the way)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 3, 10, 'in_transit_to_warehouse', 1, 'SCOP-90027', 'UPS Next Day Air', 'UPS9027890123', 
    (NOW() - INTERVAL '3 days')::date, NOW() - INTERVAL '1 day', 
    'Dell Precision 5570 in transit to warehouse. Expected arrival tomorrow morning.', 
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day'),
('bulk_to_warehouse', 4, NULL, 'in_transit_to_warehouse', 3, 'SCOP-90028', 'FedEx Freight', 'FDX9028901234', 
    (NOW() - INTERVAL '2 days')::date, NOW() - INTERVAL '18 hours', 
    'BULK: 3 Apple MacBook Air M2 laptops. ETA warehouse: Today.', 
    NOW() - INTERVAL '4 days', NOW() - INTERVAL '18 hours'),
('single_full_journey', 5, 16, 'in_transit_to_warehouse', 1, 'SCOP-90029', 'DHL Express', 'DHL9029012345', 
    (NOW() - INTERVAL '4 days')::date, NOW() - INTERVAL '2 days', 
    'HP EliteBook 850 G9 on the way. Track actively.', 
    NOW() - INTERVAL '6 days', NOW() - INTERVAL '2 days'),
('single_full_journey', 6, NULL, 'in_transit_to_warehouse', 1, 'SCOP-90030', 'UPS Ground', 'UPS9030123456', 
    (NOW() - INTERVAL '5 days')::date, NOW() - INTERVAL '3 days', 
    'Microsoft Surface Laptop Studio in transit.', 
    NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days'),
('bulk_to_warehouse', 7, NULL, 'in_transit_to_warehouse', 4, 'SCOP-90031', 'FedEx Freight', 'FDX9031234567', 
    (NOW() - INTERVAL '1 day')::date, NOW() - INTERVAL '6 hours', 
    'BULK: 4 Dell XPS 13 Plus laptops. Arriving soon.', 
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '6 hours');

-- PICKED UP FROM CLIENT (4 shipments - just collected)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, tracking_number, 
    pickup_scheduled_date, picked_up_at, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 8, 23, 'picked_up_from_client', 1, 'SCOP-90032', 'UPS Ground', 'UPS9032345678', 
    NOW()::date, NOW() - INTERVAL '3 hours', 
    'Lenovo ThinkPad X1 Carbon picked up from Enterprise Solutions HQ at 10:30 AM.', 
    NOW() - INTERVAL '2 days', NOW() - INTERVAL '3 hours'),
('single_full_journey', 1, 2, 'picked_up_from_client', 1, 'SCOP-90033', 'FedEx Express', 'FDX9033456789', 
    NOW()::date, NOW() - INTERVAL '2 hours', 
    'HP ZBook Fury G9 picked up 2 hours ago. Priority express shipping.', 
    NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 hours'),
('bulk_to_warehouse', 2, NULL, 'picked_up_from_client', 3, 'SCOP-90034', 'DHL Freight', 'DHL9034567890', 
    (NOW() - INTERVAL '1 day')::date, NOW() - INTERVAL '4 hours', 
    'BULK: 3 Apple MacBook Pro picked up. Package confirmed in excellent condition.', 
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '4 hours'),
('single_full_journey', 3, 13, 'picked_up_from_client', 1, 'SCOP-90035', 'UPS Next Day Air', 'UPS9035678901', 
    NOW()::date, NOW() - INTERVAL '1 hour', 
    'Dell XPS 15 9520 just picked up. Express delivery.', 
    NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 hour');

-- PICKUP SCHEDULED (5 shipments - scheduled for pickup)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 4, 24, 'pickup_from_client_scheduled', 1, 'SCOP-90036', 'FedEx Priority', 
    (NOW() + INTERVAL '1 day')::date, 
    'ASUS ZenBook Pro 15 OLED. Pickup tomorrow 9:00 AM - 12:00 PM. Contact confirmed.', 
    NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'),
('bulk_to_warehouse', 5, NULL, 'pickup_from_client_scheduled', 4, 'SCOP-90037', 'UPS Freight', 
    (NOW() + INTERVAL '2 days')::date, 
    'BULK: 4 Lenovo ThinkPad P16 workstations. Large pickup scheduled for day after tomorrow.', 
    NOW() - INTERVAL '1 day', NOW()),
('single_full_journey', 6, 25, 'pickup_from_client_scheduled', 1, 'SCOP-90038', 'DHL Express', 
    (NOW() + INTERVAL '1 day')::date, 
    'Apple MacBook Pro 16" M2 Max. Pickup confirmed for tomorrow afternoon.', 
    NOW() - INTERVAL '12 hours', NOW()),
('single_full_journey', 7, 26, 'pickup_from_client_scheduled', 1, 'SCOP-90039', 'UPS Ground', 
    (NOW() + INTERVAL '3 days')::date, 
    'HP EliteBook 850 G9. Pickup scheduled for 3 days from now.', 
    NOW() - INTERVAL '6 hours', NOW()),
('bulk_to_warehouse', 8, NULL, 'pickup_from_client_scheduled', 5, 'SCOP-90040', 'FedEx Freight', 
    (NOW() + INTERVAL '4 days')::date, 
    'BULK: 5 Microsoft Surface Laptop 5 units. Large pickup next week.', 
    NOW() - INTERVAL '3 hours', NOW());

-- PENDING PICKUP (5 shipments - awaiting forms)
INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, 
    pickup_scheduled_date, notes, created_at, updated_at) 
VALUES 
('single_full_journey', 1, 27, 'pending_pickup_from_client', 1, 'SCOP-90041', 
    (NOW() + INTERVAL '5 days')::date, 
    'Dell Precision 7670 workstation. Awaiting pickup form submission from client.', 
    NOW() - INTERVAL '1 day', NOW()),
('bulk_to_warehouse', 2, NULL, 'pending_pickup_from_client', 6, 'SCOP-90042', 
    (NOW() + INTERVAL '7 days')::date, 
    'BULK: 6 Acer Swift X laptops. Form needed urgently for budget approval.', 
    NOW() - INTERVAL '2 days', NOW()),
('single_full_journey', 3, 28, 'pending_pickup_from_client', 1, 'SCOP-90043', 
    (NOW() + INTERVAL '6 days')::date, 
    'Lenovo ThinkPad X1 Carbon. Client notified to submit pickup form.', 
    NOW() - INTERVAL '12 hours', NOW()),
('bulk_to_warehouse', 4, NULL, 'pending_pickup_from_client', 3, 'SCOP-90044', 
    (NOW() + INTERVAL '8 days')::date, 
    'BULK: 3 HP ZBook Studio G9 workstations. Pending form - follow up needed.', 
    NOW() - INTERVAL '6 hours', NOW()),
('single_full_journey', 5, 29, 'pending_pickup_from_client', 1, 'SCOP-90045', 
    (NOW() + INTERVAL '10 days')::date, 
    'Apple MacBook Air M2. New hire equipment - urgent processing requested.', 
    NOW() - INTERVAL '3 hours', NOW());

-- =============================================
-- SHIPMENT LAPTOPS JUNCTION TABLE
-- =============================================
-- Link shipments to specific laptops from inventory

-- Delivered shipments (1-10) - using delivered laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(1, 8),   -- Dell Precision to Alice
(2, 39),  -- HP ZBook to Emily
(3, 48),  -- Lenovo X1 to Henry
(4, 83),  -- Apple MBP to Karen
(5, 65),  -- Lenovo P1 to Nathan
(6, 13),  -- Dell XPS to Samuel
(7, 37),  -- HP EliteBook to Victor
(8, 96),  -- Surface to engineer 22
(9, 99),  -- ASUS to Catherine
(10, 48); -- Lenovo to Frank

-- In transit to engineer (11-15) - using in_transit_to_engineer laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(11, 10),  -- Dell Precision workstation
(12, 85),  -- Apple MBP 14"
(13, 40),  -- HP ZBook Fury
(14, 66),  -- Lenovo P16
(15, 19);  -- Dell XPS 13

-- Released from warehouse (16-18) - using at_warehouse laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(16, 34),  -- HP EliteBook
(17, 82),  -- Apple MBA
(18, 95);  -- Surface Laptop

-- At warehouse shipments (19-26)
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
-- BULK: 4 Lenovo ThinkPads (shipment 19)
(19, 44), (19, 45), (19, 46), (19, 47),
-- BULK: 3 Dell XPS (shipment 20)
(20, 14), (20, 15), (20, 16),
-- Single HP ZBook (shipment 21)
(21, 31),
-- Single Acer (shipment 22)
(22, 107),
-- BULK: 5 HP EliteBooks (shipment 23)
(23, 32), (23, 33), (23, 35), (23, 36), (23, 38),
-- Single ASUS (shipment 24)
(24, 100),
-- BULK: 6 Apple MBPs (shipment 25)
(25, 71), (25, 72), (25, 73), (25, 74), (25, 75), (25, 76),
-- Single Lenovo P1 (shipment 26)
(26, 58);

-- In transit to warehouse (27-31)
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(27, 7),   -- Dell Precision
-- BULK: 3 Apple MBAs (shipment 28)
(28, 81), (28, 84), (28, 86),
(29, 41),  -- HP EliteBook
(30, 93),  -- Surface
-- BULK: 4 Dell XPS 13 (shipment 31)
(31, 17), (31, 20), (31, 21), (31, 22);

-- Picked up from client (32-35) - using in_transit_to_warehouse laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(32, 49),  -- Lenovo X1 Carbon
(33, 42),  -- HP ZBook Fury
-- BULK: 3 Apple MBPs (shipment 34)
(34, 70), (34, 77), (34, 78),
(35, 11);  -- Dell XPS

-- Pickup scheduled (36-40) - using available laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(36, 98),  -- ASUS ZenBook
-- BULK: 4 Lenovo P16 (shipment 37)
(37, 59), (37, 60), (37, 61), (37, 62),
(38, 67),  -- Apple MBP 16" M2 Max
(39, 43),  -- HP EliteBook
-- BULK: 5 Surface Laptops (shipment 40)
(40, 91), (40, 92), (40, 94), (40, 95), (40, 97);

-- Pending pickup (41-45) - using available laptops
INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES
(41, 4),   -- Dell Precision
-- BULK: 6 Acer Swift X (shipment 42)
(42, 101), (42, 102), (42, 103), (42, 104), (42, 105), (42, 108),
(43, 50),  -- Lenovo X1 Carbon
-- BULK: 3 HP ZBooks (shipment 44)
(44, 26), (44, 27), (44, 28),
(45, 79);  -- Apple MBA

-- =============================================
-- RECEPTION REPORTS
-- =============================================

-- Reception reports for delivered shipments (historical)
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES
(1, 1008, NOW() - INTERVAL '60 days', 
'Dell Precision 5570 received in excellent condition. Serial number verified. All original packaging intact with factory seals. Visual inspection shows no cosmetic damage. Power-on test successful. BIOS info matches specifications exactly. GPU (NVIDIA RTX A2000) tested with benchmark - performing within spec. All ports tested and functional. Included accessories: 240W power adapter, carrying case, wireless keyboard/mouse set. Device charged to 92%. Logged into inventory system. Ready for assignment.',
ARRAY['/uploads/reception/ship001_ext.jpg', '/uploads/reception/ship001_screen.jpg', '/uploads/reception/ship001_accessories.jpg']),

(2, 1009, NOW() - INTERVAL '50 days',
'HP ZBook Studio G9 workstation received. High-end unit for video editing. Original HP packaging, all seals intact. Display quality excellent - 4K DreamColor panel tested, no dead pixels. NVIDIA RTX A3000 GPU tested - benchmark scores excellent. RAM: 64GB DDR5 verified. Storage: 2TB NVMe verified. All accessories present: 200W power adapter, premium carrying case, USB-C dock (tested - dual 4K @ 60Hz working), wireless peripherals. Unit powered on successfully, firmware up to date.',
ARRAY['/uploads/reception/ship002_unit.jpg', '/uploads/reception/ship002_display.jpg', '/uploads/reception/ship002_accessories.jpg']),

(3, 1010, NOW() - INTERVAL '45 days',
'Lenovo ThinkPad X1 Carbon Gen 10 received in pristine condition. Premium business laptop - flagship model. Carbon fiber chassis immaculate. WQUXGA 4K display tested - crystal clear, no defects. Keyboard legendary ThinkPad quality verified. TrackPoint and touchpad working perfectly. All Thunderbolt 4 ports functional. 5G WWAN module detected and working. Battery health: 100% (new unit). Included: 65W USB-C adapter, premium backpack, wireless mouse, USB-C dock (all ports functional).',
ARRAY['/uploads/reception/ship003_laptop.jpg', '/uploads/reception/ship003_display.jpg']),

(4, 1008, NOW() - INTERVAL '40 days',
'Apple MacBook Pro 16" M2 Max received. High-value unit ($4,500+). Original Apple retail box sealed with authenticity stickers intact. Serial number verified against Apple GSX system. M2 Max chip verified - 12-core CPU, 38-core GPU. Liquid Retina XDR display perfect (no dead pixels, no backlight bleeding). All 3 Thunderbolt 4 ports functional. Storage: 2TB verified. RAM: 96GB unified memory verified. Battery cycles: 0 (brand new). macOS Ventura pre-installed and updated. Charging cable and adapter present. Stored in secure area.',
ARRAY['/uploads/reception/ship004_box.jpg', '/uploads/reception/ship004_open.jpg', '/uploads/reception/ship004_serial.jpg']),

(5, 1009, NOW() - INTERVAL '35 days',
'Lenovo ThinkPad P1 Gen 5 mobile workstation received. GPU workstation for ML/AI development. Original Lenovo packaging intact. 16" 4K OLED display tested - stunning quality, no issues. NVIDIA RTX A5500 16GB GPU tested: CUDA detected, benchmark successful. CPU: Intel i9-12900H verified. RAM: 64GB DDR5 confirmed. Storage: 2TB NVMe verified. All Thunderbolt 4 ports tested. Included: 230W power adapter, workstation dock (tested), wireless peripherals. Performance tests passed. Ready for deployment.',
ARRAY['/uploads/reception/ship005_workstation.jpg', '/uploads/reception/ship005_gpu_test.jpg']),

(6, 1010, NOW() - INTERVAL '33 days',
'Dell XPS 15 9520 received in excellent condition. Premium developer laptop. 4K OLED touch display gorgeous - tested, no dead pixels. NVIDIA RTX 3050 Ti GPU functional. All USB-C/Thunderbolt 4 ports working. Included accessories: power adapter, carrying case, wireless mouse. Setup and ready.',
ARRAY['/uploads/reception/ship006_laptop.jpg']),

(7, 1008, NOW() - INTERVAL '30 days',
'HP EliteBook 850 G9 business laptop received. Professional finish, no damage. 15.6" FHD display clear. Security features tested: fingerprint reader working, IR camera for Windows Hello functional. LTE module detected. All accessories present. Enterprise ready.',
ARRAY['/uploads/reception/ship007_laptop.jpg']),

(8, 1009, NOW() - INTERVAL '27 days',
'Microsoft Surface Laptop Studio received in good condition. 14.4" PixelSense Flow touch display working perfectly. NVIDIA RTX 3050 Ti tested - functional. Hinge mechanism smooth. All accessories present including Surface Pen. Ready for creative work.',
ARRAY['/uploads/reception/ship008_surface.jpg']),

(9, 1010, NOW() - INTERVAL '25 days',
'ASUS ZenBook Pro 15 OLED received. 15.6" 4K OLED touch display tested - excellent quality. NVIDIA RTX 3060 GPU working well. Premium build quality. All accessories included. Creative workstation ready.',
ARRAY['/uploads/reception/ship009_asus.jpg']),

(10, 1008, NOW() - INTERVAL '23 days',
'Lenovo ThinkPad X1 Carbon received in perfect condition. All tests passed. Premium ultrabook ready for deployment. No issues found.',
ARRAY['/uploads/reception/ship010_lenovo.jpg']);

-- Reception reports for recent arrivals at warehouse
INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls) VALUES
(19, 1009, NOW() - INTERVAL '1 day',
'BULK RECEPTION: 4x Lenovo ThinkPad X1 Carbon Gen 10 laptops received. All boxes in excellent condition, no shipping damage. Serial numbers verified for all units. Individual inspection performed on each: All 4 units have carbon fiber chassis in perfect condition, WQUXGA displays tested (no dead pixels on any), keyboards excellent (classic ThinkPad quality), 5G WWAN in all units, battery health 98-100% on all. Accessories for all: 4x 65W adapters, 4x backpacks, 4x wireless mice, 4x USB-C docks. All units ready for assignment.',
ARRAY['/uploads/reception/ship019_bulk.jpg', '/uploads/reception/ship019_lineup.jpg']),

(20, 1010, NOW() - INTERVAL '18 hours',
'BULK RECEPTION: 3x Dell XPS 15 9520 received. All units pristine. 4K OLED displays on all three tested - gorgeous. RTX 3050 Ti GPUs all functional. All accessories present for each unit. Ready for distribution.',
ARRAY['/uploads/reception/ship020_bulk.jpg']),

(21, 1008, NOW() - INTERVAL '2 days',
'HP ZBook Studio G9 received. Workstation-class laptop in excellent condition. 4K DreamColor display perfect. RTX A3000 GPU tested successfully. All accessories present. Ready for assignment.',
ARRAY['/uploads/reception/ship021_zbook.jpg']),

(22, 1009, NOW() - INTERVAL '1 day',
'Acer Swift X budget performance laptop received in good condition. 14" FHD IPS display clear. RTX 3050 Ti working. Good value unit ready for use.',
ARRAY['/uploads/reception/ship022_acer.jpg']),

(23, 1010, NOW() - INTERVAL '12 hours',
'BULK RECEPTION: 5x HP EliteBook 840 G9 for department expansion. All units received in excellent condition. Security features tested on all: fingerprint readers working, IR cameras functional. All units ready for deployment.',
ARRAY['/uploads/reception/ship023_bulk.jpg', '/uploads/reception/ship023_elitebooks.jpg']),

(24, 1008, NOW() - INTERVAL '18 hours',
'ASUS ROG Zephyrus G14 gaming/development laptop received. AMD Ryzen 9 6900HS tested. 14" QHD+ 120Hz display excellent. Radeon RX 6800S GPU benchmarked successfully. Ready for assignment.',
ARRAY['/uploads/reception/ship024_rog.jpg']),

(25, 1009, NOW() - INTERVAL '6 hours',
'HIGH-VALUE BULK RECEPTION: 6x Apple MacBook Pro units received. Total value: $18,000+. White-glove handling applied. Original Apple packaging on all units - sealed retail boxes with authenticity stickers intact. Serial numbers verified against Apple GSX for all 6 units. Individual inspection: All units have M2 Pro/Max chips verified, Liquid Retina XDR displays perfect (no dead pixels on any), all Thunderbolt ports functional. Storage and RAM verified per unit. Battery cycles: 0-2 (essentially new on all). All charging cables and adapters present. Stored in climate-controlled secure cage. Comprehensive photos taken.',
ARRAY['/uploads/reception/ship025_all.jpg', '/uploads/reception/ship025_boxes.jpg', '/uploads/reception/ship025_secure.jpg']),

(26, 1010, NOW() - INTERVAL '8 hours',
'Lenovo ThinkPad P1 Gen 5 workstation received. 16" 4K OLED display perfect. RTX A5500 GPU tested - excellent performance. All accessories present. High-end workstation ready.',
ARRAY['/uploads/reception/ship026_p1.jpg']);

-- =============================================
-- PICKUP FORMS
-- =============================================

-- Pickup forms for all shipments that have been picked up or scheduled
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES
-- Delivered shipments pickup forms
(1, 1023, NOW() - INTERVAL '67 days',
('{"contact_name":"Sarah Mitchell","contact_email":"sarah.mitchell@techcorp.com","contact_phone":"+1-555-0101","pickup_address":"100 Tech Plaza, Building 1, Loading Dock A","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '65 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell 240W power adapter, carrying case, wireless keyboard/mouse set","special_instructions":"Building requires 24hr advance notification. Call Sarah before arrival."}')::jsonb),

(2, 1024, NOW() - INTERVAL '57 days',
('{"contact_name":"Michael Chen","contact_email":"michael.chen@innovate.io","contact_phone":"+1-555-0201","pickup_address":"200 Innovation Way, Warehouse Building, Bay 5","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '55 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP 200W adapter, ZBook carrying case, USB-C dock, wireless peripherals","special_instructions":"Video editing workstation. Handle with care."}')::jsonb),

(3, 1025, NOW() - INTERVAL '52 days',
('{"contact_name":"Jennifer Wang","contact_email":"jennifer.wang@globaltech.com","contact_phone":"+1-555-0301","pickup_address":"300 Global Drive, Floor 8, IT Department","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '50 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo 65W adapter, ThinkPad backpack, wireless mouse, USB-C dock","special_instructions":"Executive laptop. Check in at main reception."}')::jsonb),

(4, 1026, NOW() - INTERVAL '47 days',
('{"contact_name":"Robert Chen","contact_email":"robert.chen@digitaldynamics.com","contact_phone":"+1-555-0401","pickup_address":"400 Digital Blvd, Suite 1200","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '45 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Apple USB-C cable, Magic Mouse, Magic Keyboard, USB-C adapters, premium sleeve","special_instructions":"HIGH-VALUE: Apple MacBook Pro M2 Max. Extra insurance. Signature required."}')::jsonb),

(5, 1027, NOW() - INTERVAL '42 days',
('{"contact_name":"Linda Martinez","contact_email":"linda.martinez@cloudventures.com","contact_phone":"+1-555-0501","pickup_address":"500 Cloud Street, Building A, Floor 3","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '40 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo 230W adapter, workstation dock, wireless peripherals","special_instructions":"GPU workstation for ML/AI. Fragile equipment - handle with care."}')::jsonb);

-- Add more pickup forms for recent shipments...
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES
(11, 1025, NOW() - INTERVAL '12 days',
('{"contact_name":"Global Tech Logistics","contact_email":"logistics@globaltech.com","contact_phone":"+1-555-0302","pickup_address":"300 Global Drive, Secure Wing","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '10 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell 240W adapter, workstation dock, premium peripherals","special_instructions":"High-end GPU workstation. Requires security escort."}')::jsonb),

(19, 1025, NOW() - INTERVAL '7 days',
('{"contact_name":"IT Department","contact_email":"it@globaltech.com","contact_phone":"+1-555-0303","pickup_address":"300 Global Drive, Main Building","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '5 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":4,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":18.0,"bulk_height":12.0,"bulk_weight":38.0,"include_accessories":true,"accessories_description":"4x Lenovo adapters, 4x backpacks, 4x wireless mice, 4x docks","special_instructions":"BULK: 4 ThinkPad X1 Carbon laptops. New hire onboarding equipment."}')::jsonb);

-- More pickup forms for shipments 6-10, 12-18, 20-40
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES
(6, 1028, NOW() - INTERVAL '40 days',
('{"contact_name":"TechServices Contact","contact_email":"contact@techservices.com","contact_phone":"+1-555-0601","pickup_address":"600 Tech Services Ave","pickup_city":"Portland","pickup_state":"OR","pickup_zip":"97201","pickup_date":"' || to_char(NOW() - INTERVAL '38 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell power adapter, carrying case","special_instructions":"Ground floor loading dock."}')::jsonb),

(7, 1029, NOW() - INTERVAL '37 days',
('{"contact_name":"Data Insights Manager","contact_email":"manager@datainsights.com","contact_phone":"+1-555-0701","pickup_address":"700 Data Drive","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '35 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP EliteBook adapter and accessories","special_instructions":"Security clearance required."}')::jsonb),

(8, 1030, NOW() - INTERVAL '34 days',
('{"contact_name":"Enterprise Contact","contact_email":"contact@enterprisesolutions.com","contact_phone":"+1-555-0801","pickup_address":"800 Enterprise Way","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char(NOW() - INTERVAL '32 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Surface accessories","special_instructions":"Building reception."}')::jsonb),

(9, 1031, NOW() - INTERVAL '32 days',
('{"contact_name":"Creative Studios","contact_email":"studio@creativestudios.com","contact_phone":"+1-555-0901","pickup_address":"900 Creative Blvd","pickup_city":"Los Angeles","pickup_state":"CA","pickup_zip":"90001","pickup_date":"' || to_char(NOW() - INTERVAL '30 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"ASUS accessories","special_instructions":"Studio entrance."}')::jsonb),

(10, 1023, NOW() - INTERVAL '30 days',
('{"contact_name":"TechCorp Logistics","contact_email":"logistics@techcorp.com","contact_phone":"+1-555-0102","pickup_address":"100 Tech Plaza, Building 2","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '28 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo adapter and accessories","special_instructions":"Standard pickup."}')::jsonb),

(12, 1026, NOW() - INTERVAL '11 days',
('{"contact_name":"Digital Ops","contact_email":"ops@digitaldynamics.com","contact_phone":"+1-555-0402","pickup_address":"400 Digital Blvd, Suite 1100","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '9 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Apple accessories","special_instructions":"Priority pickup."}')::jsonb),

(13, 1027, NOW() - INTERVAL '10 days',
('{"contact_name":"Cloud Ventures IT","contact_email":"it@cloudventures.com","contact_phone":"+1-555-0502","pickup_address":"500 Cloud Street, Building B","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '8 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP workstation accessories","special_instructions":"IT department pickup."}')::jsonb),

(14, 1028, NOW() - INTERVAL '13 days',
('{"contact_name":"Tech Services Warehouse","contact_email":"warehouse@techservices.com","contact_phone":"+1-555-0602","pickup_address":"600 Tech Services Ave, Dock 3","pickup_city":"Portland","pickup_state":"OR","pickup_zip":"97201","pickup_date":"' || to_char(NOW() - INTERVAL '11 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo workstation accessories","special_instructions":"Warehouse bay 3."}')::jsonb),

(15, 1029, NOW() - INTERVAL '9 days',
('{"contact_name":"Data Insights Logistics","contact_email":"logistics@datainsights.com","contact_phone":"+1-555-0702","pickup_address":"700 Data Drive, Loading Bay","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '7 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell XPS accessories","special_instructions":"Loading bay access required."}')::jsonb),

(16, 1030, NOW() - INTERVAL '10 days',
('{"contact_name":"Enterprise IT","contact_email":"it@enterprisesolutions.com","contact_phone":"+1-555-0802","pickup_address":"800 Enterprise Way, Floor 5","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char(NOW() - INTERVAL '8 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP EliteBook accessories","special_instructions":"Floor 5 IT department."}')::jsonb),

(17, 1023, NOW() - INTERVAL '9 days',
('{"contact_name":"TechCorp Main","contact_email":"main@techcorp.com","contact_phone":"+1-555-0103","pickup_address":"100 Tech Plaza, Main Entrance","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '7 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Apple MacBook Air accessories","special_instructions":"Main reception desk."}')::jsonb),

(18, 1024, NOW() - INTERVAL '8 days',
('{"contact_name":"Innovate Solutions","contact_email":"logistics@innovate.io","contact_phone":"+1-555-0202","pickup_address":"200 Innovation Way, Main Building","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '6 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Surface Laptop accessories","special_instructions":"Main building reception."}')::jsonb),

(20, 1026, NOW() - INTERVAL '6 days',
('{"contact_name":"Digital Dynamics Bulk","contact_email":"bulk@digitaldynamics.com","contact_phone":"+1-555-0403","pickup_address":"400 Digital Blvd, Warehouse","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '4 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":20.0,"bulk_width":16.0,"bulk_height":10.0,"bulk_weight":28.0,"include_accessories":true,"accessories_description":"3x Dell XPS adapters and accessories","special_instructions":"BULK: 3 Dell XPS laptops."}')::jsonb),

(21, 1027, NOW() - INTERVAL '8 days',
('{"contact_name":"Cloud Ventures Ops","contact_email":"ops@cloudventures.com","contact_phone":"+1-555-0503","pickup_address":"500 Cloud Street, Ops Floor","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '6 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP ZBook accessories","special_instructions":"Operations floor."}')::jsonb),

(22, 1028, NOW() - INTERVAL '7 days',
('{"contact_name":"Tech Services Main","contact_email":"main@techservices.com","contact_phone":"+1-555-0603","pickup_address":"600 Tech Services Ave","pickup_city":"Portland","pickup_state":"OR","pickup_zip":"97201","pickup_date":"' || to_char(NOW() - INTERVAL '5 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Acer accessories","special_instructions":"Standard pickup."}')::jsonb),

(23, 1029, NOW() - INTERVAL '5 days',
('{"contact_name":"Data Insights Bulk","contact_email":"bulk@datainsights.com","contact_phone":"+1-555-0703","pickup_address":"700 Data Drive, Warehouse","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '3 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":5,"number_of_boxes":3,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":18.0,"bulk_height":14.0,"bulk_weight":45.0,"include_accessories":true,"accessories_description":"5x HP EliteBook adapters and accessories","special_instructions":"BULK: 5 HP EliteBook laptops for expansion."}')::jsonb),

(24, 1030, NOW() - INTERVAL '6 days',
('{"contact_name":"Enterprise Gaming Div","contact_email":"gaming@enterprisesolutions.com","contact_phone":"+1-555-0803","pickup_address":"800 Enterprise Way, R&D","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char(NOW() - INTERVAL '4 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"ASUS ROG accessories","special_instructions":"R&D department."}')::jsonb),

(25, 1023, NOW() - INTERVAL '4 days',
('{"contact_name":"TechCorp iOS Team","contact_email":"ios@techcorp.com","contact_phone":"+1-555-0104","pickup_address":"100 Tech Plaza, iOS Development Floor","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '2 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":6,"number_of_boxes":4,"assignment_type":"bulk","bulk_length":26.0,"bulk_width":20.0,"bulk_height":16.0,"bulk_weight":65.0,"include_accessories":true,"accessories_description":"6x Apple MacBook Pro adapters, cables, accessories - HIGH VALUE","special_instructions":"HIGH-VALUE BULK: 6 MacBook Pro units for iOS team. Extra insurance and security."}')::jsonb),

(26, 1024, NOW() - INTERVAL '5 days',
('{"contact_name":"Innovate Solutions Dev","contact_email":"dev@innovate.io","contact_phone":"+1-555-0203","pickup_address":"200 Innovation Way, Dev Building","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '3 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo P1 workstation accessories","special_instructions":"Development building."}')::jsonb),

(27, 1025, NOW() - INTERVAL '5 days',
('{"contact_name":"Global Tech Hardware","contact_email":"hardware@globaltech.com","contact_phone":"+1-555-0304","pickup_address":"300 Global Drive, Hardware Dept","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '3 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell Precision accessories","special_instructions":"Hardware department."}')::jsonb),

(28, 1026, NOW() - INTERVAL '4 days',
('{"contact_name":"Digital Dynamics Apple","contact_email":"apple@digitaldynamics.com","contact_phone":"+1-555-0404","pickup_address":"400 Digital Blvd, Apple Dev Floor","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char(NOW() - INTERVAL '2 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":16.0,"bulk_height":12.0,"bulk_weight":32.0,"include_accessories":true,"accessories_description":"3x MacBook Air adapters and accessories","special_instructions":"BULK: 3 MacBook Air units."}')::jsonb),

(29, 1027, NOW() - INTERVAL '6 days',
('{"contact_name":"Cloud Ventures Hardware","contact_email":"hardware@cloudventures.com","contact_phone":"+1-555-0504","pickup_address":"500 Cloud Street, Hardware Bay","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char(NOW() - INTERVAL '4 days', 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP EliteBook accessories","special_instructions":"Hardware bay."}')::jsonb),

(30, 1030, NOW() - INTERVAL '7 days',
('{"contact_name":"Enterprise Microsoft Div","contact_email":"microsoft@enterprisesolutions.com","contact_phone":"+1-555-0804","pickup_address":"800 Enterprise Way, Microsoft Wing","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char(NOW() - INTERVAL '5 days', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Surface Laptop Studio accessories","special_instructions":"Microsoft development wing."}')::jsonb),

(31, 1029, NOW() - INTERVAL '3 days',
('{"contact_name":"Data Insights XPS Team","contact_email":"xps@datainsights.com","contact_phone":"+1-555-0704","pickup_address":"700 Data Drive, Team Floor","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char(NOW() - INTERVAL '1 day', 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":4,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":18.0,"bulk_height":10.0,"bulk_weight":36.0,"include_accessories":true,"accessories_description":"4x Dell XPS 13 adapters and accessories","special_instructions":"BULK: 4 Dell XPS 13 Plus laptops."}')::jsonb),

(32, 1030, NOW() - INTERVAL '2 days',
('{"contact_name":"Enterprise Lenovo Team","contact_email":"lenovo@enterprisesolutions.com","contact_phone":"+1-555-0805","pickup_address":"800 Enterprise Way, Lenovo Dept","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char(NOW() - INTERVAL '3 hours', 'YYYY-MM-DD HH24:MI') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Lenovo X1 Carbon accessories","special_instructions":"Just picked up."}')::jsonb),

(33, 1023, NOW() - INTERVAL '1 day',
('{"contact_name":"TechCorp Workstation Team","contact_email":"workstation@techcorp.com","contact_phone":"+1-555-0105","pickup_address":"100 Tech Plaza, Workstation Lab","pickup_city":"San Francisco","pickup_state":"CA","pickup_zip":"94105","pickup_date":"' || to_char(NOW() - INTERVAL '2 hours', 'YYYY-MM-DD HH24:MI') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP ZBook Fury accessories","special_instructions":"Express priority."}')::jsonb),

(34, 1024, NOW() - INTERVAL '3 days',
('{"contact_name":"Innovate Apple Dev","contact_email":"appledev@innovate.io","contact_phone":"+1-555-0204","pickup_address":"200 Innovation Way, Apple Building","pickup_city":"Austin","pickup_state":"TX","pickup_zip":"78701","pickup_date":"' || to_char(NOW() - INTERVAL '4 hours', 'YYYY-MM-DD HH24:MI') || '","pickup_time_slot":"morning","number_of_laptops":3,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":22.0,"bulk_width":17.0,"bulk_height":11.0,"bulk_weight":34.0,"include_accessories":true,"accessories_description":"3x MacBook Pro adapters - HIGH VALUE","special_instructions":"BULK: 3 MacBook Pro - just picked up."}')::jsonb),

(35, 1025, NOW() - INTERVAL '2 days',
('{"contact_name":"Global Tech Premium","contact_email":"premium@globaltech.com","contact_phone":"+1-555-0305","pickup_address":"300 Global Drive, Premium Suite","pickup_city":"Seattle","pickup_state":"WA","pickup_zip":"98101","pickup_date":"' || to_char(NOW() - INTERVAL '1 hour', 'YYYY-MM-DD HH24:MI') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"Dell XPS accessories","special_instructions":"Just collected - express delivery."}')::jsonb),

(36, 1026, NOW() - INTERVAL '2 days',
('{"contact_name":"Digital Dynamics ASUS","contact_email":"asus@digitaldynamics.com","contact_phone":"+1-555-0405","pickup_address":"400 Digital Blvd, Creative Suite","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char((NOW() + INTERVAL '1 day'), 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"ASUS ZenBook Pro accessories","special_instructions":"Pickup tomorrow 9 AM."}')::jsonb),

(37, 1027, NOW() - INTERVAL '1 day',
('{"contact_name":"Cloud Ventures Lenovo","contact_email":"lenovo@cloudventures.com","contact_phone":"+1-555-0505","pickup_address":"500 Cloud Street, Lenovo Lab","pickup_city":"Denver","pickup_state":"CO","pickup_zip":"80202","pickup_date":"' || to_char((NOW() + INTERVAL '2 days'), 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":4,"number_of_boxes":2,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":18.0,"bulk_height":12.0,"bulk_weight":42.0,"include_accessories":true,"accessories_description":"4x Lenovo P16 workstation accessories","special_instructions":"BULK: 4 Lenovo P16 workstations - pickup day after tomorrow."}')::jsonb),

(38, 1026, NOW() - INTERVAL '12 hours',
('{"contact_name":"Digital Dynamics Apple Pro","contact_email":"applepro@digitaldynamics.com","contact_phone":"+1-555-0406","pickup_address":"400 Digital Blvd, iOS Lab","pickup_city":"Boston","pickup_state":"MA","pickup_zip":"02101","pickup_date":"' || to_char((NOW() + INTERVAL '1 day'), 'YYYY-MM-DD') || '","pickup_time_slot":"afternoon","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"MacBook Pro 16 M2 Max accessories - HIGH VALUE","special_instructions":"HIGH VALUE: Tomorrow afternoon pickup confirmed."}')::jsonb),

(39, 1029, NOW() - INTERVAL '6 hours',
('{"contact_name":"Data Insights HP","contact_email":"hp@datainsights.com","contact_phone":"+1-555-0705","pickup_address":"700 Data Drive, HP Dept","pickup_city":"Chicago","pickup_state":"IL","pickup_zip":"60601","pickup_date":"' || to_char((NOW() + INTERVAL '3 days'), 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":1,"number_of_boxes":1,"assignment_type":"single","include_accessories":true,"accessories_description":"HP EliteBook accessories","special_instructions":"Pickup in 3 days."}')::jsonb),

(40, 1030, NOW() - INTERVAL '3 hours',
('{"contact_name":"Enterprise Surface Team","contact_email":"surface@enterprisesolutions.com","contact_phone":"+1-555-0806","pickup_address":"800 Enterprise Way, Surface Lab","pickup_city":"Miami","pickup_state":"FL","pickup_zip":"33101","pickup_date":"' || to_char((NOW() + INTERVAL '4 days'), 'YYYY-MM-DD') || '","pickup_time_slot":"morning","number_of_laptops":5,"number_of_boxes":3,"assignment_type":"bulk","bulk_length":24.0,"bulk_width":19.0,"bulk_height":13.0,"bulk_weight":40.0,"include_accessories":true,"accessories_description":"5x Surface Laptop 5 adapters and pens","special_instructions":"BULK: 5 Surface Laptop 5 units - pickup next week."}')::jsonb);

-- =============================================
-- DELIVERY FORMS
-- =============================================

-- Delivery forms for completed deliveries
INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls) VALUES
(1, 1, NOW() - INTERVAL '55 days',
'Dell Precision 5570 delivered successfully to Alice Johnson at TechCorp SF office. Delivery completed on schedule. Package inspection with engineer present: Box unopened, seals intact. Engineer opened package on-site. Laptop condition: Pristine. Powered on successfully - boot time under 30 seconds. Display quality verified - UHD+ gorgeous. Engineer very pleased with the workstation. All accessories verified present: power adapter, carrying case, wireless keyboard/mouse set. Tested: All USB-C ports, keyboard, trackpad, speakers - all working perfectly. Connected to company WiFi successfully. Setup assistance provided - connected to dual monitors via Thunderbolt dock. Engineer confirmed receipt and high satisfaction.',
ARRAY['/uploads/delivery/ship001_sig.jpg', '/uploads/delivery/ship001_setup.jpg']),

(2, 5, NOW() - INTERVAL '45 days',
'HP ZBook Studio G9 delivered to Emily Rodriguez at Innovate Solutions Austin office. Video editing workstation. Engineer and IT manager present. Package opened on-site - perfect condition. 4K DreamColor display verified - stunning quality. RTX A3000 GPU tested with video editing software - performance excellent. All accessories verified. Software setup: DaVinci Resolve tested - renders fast, no issues. Adobe Creative Cloud installed and verified. Engineer feedback: "Exactly what our video team needs. Color accuracy on DreamColor display is perfect." Setup complete with dual 4K monitors. Engineer satisfaction: Excellent.',
ARRAY['/uploads/delivery/ship002_sig.jpg', '/uploads/delivery/ship002_dual_monitors.jpg']),

(3, 8, NOW() - INTERVAL '40 days',
'Lenovo ThinkPad X1 Carbon Gen 10 delivered to Henry Thompson at Global Tech Seattle. Executive ultrabook. Delivery smooth. WQUXGA 4K display verified - crystal clear. 5G WWAN tested - connected successfully. All accessories present. Engineer confirmed all features working. Setup complete. Satisfaction: Excellent.',
ARRAY['/uploads/delivery/ship003_sig.jpg']),

(4, 11, NOW() - INTERVAL '35 days',
'Apple MacBook Pro 16" M2 Max delivered to Karen Lee at Digital Dynamics Boston. High-end development machine. Engineer confirmed receipt. M2 Max performance tested - excellent. Display perfect. All accessories verified. Connected to company network. Development environment setup verified. Xcode installed and tested. Engineer expressed high satisfaction. Premium device for iOS development.',
ARRAY['/uploads/delivery/ship004_sig.jpg', '/uploads/delivery/ship004_setup.jpg']);

-- =============================================
-- AUDIT LOGS
-- =============================================

-- Recent activity audit logs
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details) VALUES
-- Shipment creations
(1000, 'shipment_created', 'shipment', 45, NOW() - INTERVAL '3 hours', '{"jira_ticket":"SCOP-90045","type":"single_full_journey","status":"pending_pickup_from_client"}'),
(1000, 'shipment_created', 'shipment', 44, NOW() - INTERVAL '6 hours', '{"jira_ticket":"SCOP-90044","type":"bulk_to_warehouse","laptop_count":3}'),
(1000, 'shipment_created', 'shipment', 43, NOW() - INTERVAL '12 hours', '{"jira_ticket":"SCOP-90043","type":"single_full_journey"}'),

-- Recent pickups
(1000, 'status_updated', 'shipment', 35, NOW() - INTERVAL '1 hour', '{"old_status":"pickup_from_client_scheduled","new_status":"picked_up_from_client","tracking":"UPS9035678901"}'),
(1000, 'status_updated', 'shipment', 32, NOW() - INTERVAL '3 hours', '{"old_status":"pickup_from_client_scheduled","new_status":"picked_up_from_client","tracking":"UPS9032345678"}'),

-- Recent warehouse arrivals
(1009, 'reception_report_created', 'reception_report', 19, NOW() - INTERVAL '1 day', '{"shipment_id":19,"bulk":true,"units":4}'),
(1009, 'status_updated', 'shipment', 19, NOW() - INTERVAL '1 day', '{"old_status":"in_transit_to_warehouse","new_status":"at_warehouse"}'),
(1009, 'reception_report_created', 'reception_report', 25, NOW() - INTERVAL '6 hours', '{"shipment_id":25,"bulk":true,"units":6,"high_value":true}'),

-- Engineer assignments
(1015, 'engineer_assigned', 'shipment', 16, NOW() - INTERVAL '6 hours', '{"engineer_id":21,"engineer_name":"Wendy Martinez"}'),
(1015, 'status_updated', 'shipment', 16, NOW() - INTERVAL '6 hours', '{"old_status":"at_warehouse","new_status":"released_from_warehouse"}'),

-- Recent deliveries
(1000, 'status_updated', 'shipment', 10, NOW() - INTERVAL '18 days', '{"old_status":"in_transit_to_engineer","new_status":"delivered"}'),
(1001, 'delivery_form_created', 'delivery_form', 10, NOW() - INTERVAL '18 days', '{"shipment_id":10,"engineer_id":6}');

-- =============================================
-- SUMMARY OUTPUT
-- =============================================

SELECT '========================================' AS separator;
SELECT 'COMPREHENSIVE SHIPMENTS DATA LOADED!' AS message;
SELECT '========================================' AS separator;
SELECT '' AS blank;

SELECT 'SHIPMENTS SUMMARY' AS section;
SELECT '-----------------' AS underline;
SELECT 'Total shipments: ' || COUNT(*) as total FROM shipments WHERE id >= 1;
SELECT 'By status:' as breakdown;

SELECT 
    '  ' || status || ': ' || COUNT(*) as status_breakdown
FROM shipments 
WHERE id >= 1
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
SELECT 'By type:' as type_breakdown;
SELECT 
    '  ' || shipment_type || ': ' || COUNT(*) as type_count
FROM shipments
WHERE id >= 1
GROUP BY shipment_type;

SELECT '' AS blank;
SELECT 'FORMS & REPORTS' AS section;
SELECT '---------------' AS underline;
SELECT COUNT(*) || ' pickup forms' FROM pickup_forms WHERE shipment_id >= 1;
SELECT COUNT(*) || ' reception reports' FROM reception_reports WHERE shipment_id >= 1;
SELECT COUNT(*) || ' delivery forms' FROM delivery_forms WHERE shipment_id >= 1;
SELECT COUNT(*) || ' audit log entries' FROM audit_logs WHERE entity_id >= 1 AND entity_type = 'shipment';

SELECT '' AS blank;
SELECT '========================================' AS separator;
SELECT 'Ready for comprehensive testing!' AS message;
SELECT 'Application: http://localhost:8080' AS next_step;
SELECT '========================================' AS separator;


