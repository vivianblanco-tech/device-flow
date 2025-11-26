-- =============================================
-- COMPREHENSIVE SHIPMENTS DATA v4.0
-- Complete Shipments, Forms, Reports & Audit Logs
-- ~100 shipments across 3 types with proper timing
-- =============================================
-- This script adds complete shipment data with all three types:
-- 1. single_full_journey - Client → Warehouse → Engineer (~40 shipments)
-- 2. bulk_to_warehouse - Multiple laptops → Warehouse (~35 shipments)
-- 3. warehouse_to_engineer - Warehouse inventory → Engineer (~25 shipments)
--
-- Average Delivery Time Target: ~2.5-2.9 days
-- Timeline: Spread over 8 months (most recent to 8 months ago)
-- Run this AFTER loading comprehensive-sample-data-v4.sql
-- =============================================

-- Helper function to generate tracking numbers
DO $$
DECLARE
    shipment_id_counter INT := 1;
    laptop_id_counter INT := 1;
    max_laptop_id INT;
    jira_counter INT := 90001;
    courier_names TEXT[] := ARRAY['FedEx Express', 'UPS Next Day Air', 'DHL Express', 'FedEx Ground', 'UPS Ground'];
    tracking_prefixes TEXT[] := ARRAY['FDX', 'UPS', 'DHL', 'FDX', 'UPS'];
    delivery_days FLOAT;
    pickup_date TIMESTAMP;
    delivered_date TIMESTAMP;
    base_date TIMESTAMP;
    days_offset INT;
    courier_idx INT;
    tracking_num TEXT;
    engineer_id INT;
    client_id INT;
    laptop_ids INT[];
    bulk_size INT;
    status_text TEXT;
    i INT;
    j INT;
BEGIN
    -- Get max laptop ID to ensure we don't exceed available laptops
    SELECT COALESCE(MAX(id), 0) INTO max_laptop_id FROM laptops;
    
    -- ============================================
    -- SINGLE FULL JOURNEY SHIPMENTS (~40 shipments)
    -- ============================================
    -- Distribution: ~55% delivered, ~15% in_transit, ~15% at_warehouse, ~15% pending/scheduled
    
    FOR i IN 1..40 LOOP
        -- Calculate base date (spread over 8 months)
        days_offset := (i - 1) * 6; -- Spread shipments every ~6 days
        base_date := NOW() - (days_offset || ' days')::INTERVAL;
        
        -- Select random client (1-15) and engineer (1-35)
        client_id := 1 + ((i - 1) % 15);
        engineer_id := 1 + ((i - 1) % 35);
        
        -- Select courier
        courier_idx := 1 + ((i - 1) % 5);
        tracking_num := tracking_prefixes[courier_idx] || (900000 + jira_counter)::TEXT;
        
        -- Determine status and dates based on position
        IF i <= 22 THEN
            -- Delivered shipments (55% = 22 shipments)
            -- Average delivery time: ~2.15 days (under 3 days)
            delivery_days := 1.8 + (RANDOM() * 0.7); -- Random between 1.8-2.5 days
            pickup_date := base_date - (delivery_days + 1)::INT * INTERVAL '1 day';
            delivered_date := pickup_date + (delivery_days || ' days')::INTERVAL;
            
            status_text := 'delivered';
            
            -- Insert shipment
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
                released_warehouse_at, delivered_at, notes, created_at, updated_at)
            VALUES ('single_full_journey', client_id, engineer_id, status_text::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num,
                pickup_date::date, pickup_date, pickup_date + INTERVAL '2 days',
                pickup_date + INTERVAL '2 days' + INTERVAL '12 hours', delivered_date,
                'Single laptop delivery completed. Engineer confirmed receipt.',
                pickup_date - INTERVAL '2 days', delivered_date)
            RETURNING id INTO shipment_id_counter;
            
            -- Link laptop (use available laptops, cycling through)
            INSERT INTO shipment_laptops (shipment_id, laptop_id)
            VALUES (shipment_id_counter, laptop_id_counter);
            
            -- Update laptop status
            UPDATE laptops SET status = 'delivered', software_engineer_id = engineer_id
            WHERE id = laptop_id_counter;
            
            -- Pickup form
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5),
                pickup_date - INTERVAL '2 days',
                jsonb_build_object(
                    'contact_name', 'Contact ' || i,
                    'contact_email', 'contact' || i || '@client.com',
                    'contact_phone', '+1-555-' || LPAD(i::TEXT, 4, '0'),
                    'pickup_address', 'Address ' || i,
                    'pickup_city', 'City',
                    'pickup_state', 'CA',
                    'pickup_zip', '90001',
                    'pickup_date', to_char(pickup_date::date, 'YYYY-MM-DD'),
                    'pickup_time_slot', CASE WHEN i % 2 = 0 THEN 'morning' ELSE 'afternoon' END,
                    'special_instructions', 'Standard pickup procedure'
                ));
            
            -- Reception report (approved)
            INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id,
                received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition,
                status, approved_by, approved_at, created_at, updated_at)
            SELECT laptop_id_counter, shipment_id_counter, client_id, tracking_num, 7 + ((i - 1) % 6),
                pickup_date + INTERVAL '2 days',
                'Laptop received in excellent condition. All tests passed.',
                '/uploads/reception/laptop' || laptop_id_counter || '_serial.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_external.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_working.jpg',
                'approved', 1 + ((i - 1) % 6), pickup_date + INTERVAL '2 days' + INTERVAL '6 hours',
                pickup_date + INTERVAL '2 days', pickup_date + INTERVAL '2 days' + INTERVAL '6 hours'
            WHERE NOT EXISTS (SELECT 1 FROM reception_reports WHERE laptop_id = laptop_id_counter);
            
            -- Delivery form
            INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls)
            VALUES (shipment_id_counter, engineer_id, delivered_date,
                'Device delivered successfully. Engineer confirmed satisfaction.',
                ARRAY['/uploads/delivery/shipment' || LPAD(shipment_id_counter::TEXT, 3, '0') || '_photo1.jpg',
                      '/uploads/delivery/shipment' || LPAD(shipment_id_counter::TEXT, 3, '0') || '_photo2.jpg']);
        
        ELSIF i <= 28 THEN
            -- In transit to engineer (15% = 6 shipments)
            status_text := 'in_transit_to_engineer';
            pickup_date := base_date - INTERVAL '5 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
                released_warehouse_at, eta_to_engineer, notes, created_at, updated_at)
            VALUES ('single_full_journey', client_id, engineer_id, status_text::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num,
                pickup_date::date, pickup_date, pickup_date + INTERVAL '2 days',
                pickup_date + INTERVAL '2 days' + INTERVAL '12 hours', base_date + INTERVAL '1 day',
                'Shipment in transit to engineer. Expected delivery tomorrow.',
                pickup_date - INTERVAL '2 days', NOW())
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
            UPDATE laptops SET status = 'in_transit_to_engineer', software_engineer_id = engineer_id WHERE id = laptop_id_counter;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5), pickup_date - INTERVAL '2 days',
                jsonb_build_object('contact_name', 'Contact ' || i, 'pickup_date', to_char(pickup_date::date, 'YYYY-MM-DD')));
            
            INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id,
                received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition,
                status, created_at, updated_at)
            SELECT laptop_id_counter, shipment_id_counter, client_id, tracking_num, 7 + ((i - 1) % 6),
                pickup_date + INTERVAL '2 days', 'Received and approved.',
                '/uploads/reception/laptop' || laptop_id_counter || '_serial.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_external.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_working.jpg',
                'approved', pickup_date + INTERVAL '2 days', pickup_date + INTERVAL '2 days'
            WHERE NOT EXISTS (SELECT 1 FROM reception_reports WHERE laptop_id = laptop_id_counter);
        
        ELSIF i <= 34 THEN
            -- At warehouse (15% = 6 shipments)
            status_text := 'at_warehouse';
            pickup_date := base_date - INTERVAL '3 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
                notes, created_at, updated_at)
            VALUES ('single_full_journey', client_id, engineer_id, status_text::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num,
                pickup_date::date, pickup_date, base_date - INTERVAL '1 day',
                'Awaiting engineer assignment.',
                pickup_date - INTERVAL '2 days', NOW())
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
            UPDATE laptops SET status = 'at_warehouse' WHERE id = laptop_id_counter;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5), pickup_date - INTERVAL '2 days',
                jsonb_build_object('contact_name', 'Contact ' || i, 'pickup_date', to_char(pickup_date::date, 'YYYY-MM-DD')));
            
            INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id,
                received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition,
                status, created_at, updated_at)
            SELECT laptop_id_counter, shipment_id_counter, client_id, tracking_num, 7 + ((i - 1) % 6),
                base_date - INTERVAL '1 day', 'Received at warehouse.',
                '/uploads/reception/laptop' || laptop_id_counter || '_serial.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_external.jpg',
                '/uploads/reception/laptop' || laptop_id_counter || '_working.jpg',
                'pending_approval', base_date - INTERVAL '1 day', base_date - INTERVAL '1 day'
            WHERE NOT EXISTS (SELECT 1 FROM reception_reports WHERE laptop_id = laptop_id_counter);
        
        ELSE
            -- Pending/scheduled (15% = 6 shipments)
            IF i % 2 = 0 THEN
                status_text := 'pickup_from_client_scheduled';
                pickup_date := base_date + INTERVAL '1 day';
            ELSE
                status_text := 'pending_pickup_from_client';
                pickup_date := base_date + INTERVAL '3 days';
            END IF;
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, notes, created_at, updated_at)
            VALUES ('single_full_journey', client_id, engineer_id, status_text::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date::date,
                'Awaiting pickup from client.',
                base_date - INTERVAL '1 day', NOW())
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
        END IF;
        
        -- shipment_id_counter is set by RETURNING clause, no need to increment
        -- Increment and wrap around if exceeds max
        laptop_id_counter := laptop_id_counter + 1;
        IF laptop_id_counter > max_laptop_id THEN
            laptop_id_counter := 1;
        END IF;
        jira_counter := jira_counter + 1;
    END LOOP;
    
    -- ============================================
    -- BULK TO WAREHOUSE SHIPMENTS (~35 shipments)
    -- ============================================
    -- Distribution: ~40% delivered (historical), ~30% at_warehouse, ~20% in_transit, ~10% pending
    
    FOR i IN 1..35 LOOP
        days_offset := (i - 1) * 7;
        base_date := NOW() - (days_offset || ' days')::INTERVAL;
        client_id := 1 + ((i - 1) % 15);
        courier_idx := 1 + ((i - 1) % 5);
        tracking_num := tracking_prefixes[courier_idx] || (900000 + jira_counter)::TEXT;
        
        -- Bulk size: 2-6 laptops
        bulk_size := 2 + ((i - 1) % 5);
        laptop_ids := ARRAY[]::INT[];
        
        -- Collect laptop IDs (cycle through available laptops)
        FOR j IN 1..bulk_size LOOP
            laptop_ids := array_append(laptop_ids, laptop_id_counter);
            -- Increment and wrap around if exceeds max
            laptop_id_counter := laptop_id_counter + 1;
            IF laptop_id_counter > max_laptop_id THEN
                laptop_id_counter := 1;
            END IF;
        END LOOP;
        
        IF i <= 14 THEN
            -- Delivered bulk shipments (40% = 14 shipments)
            -- Average delivery time: ~2.15 days (under 3 days)
            delivery_days := 1.8 + (RANDOM() * 0.7); -- Random between 1.8-2.5 days
            pickup_date := base_date - (delivery_days + 3)::INT * INTERVAL '1 day';
            delivered_date := pickup_date + (delivery_days || ' days')::INTERVAL;
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
                released_warehouse_at, delivered_at, notes, created_at, updated_at)
            VALUES ('bulk_to_warehouse', client_id, NULL, 'delivered'::shipment_status, bulk_size, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date::date, pickup_date, pickup_date + INTERVAL '2 days',
                pickup_date + INTERVAL '2 days' + INTERVAL '12 hours', delivered_date,
                'Bulk shipment delivered. All laptops assigned to engineers.',
                pickup_date - INTERVAL '2 days', delivered_date)
            RETURNING id INTO shipment_id_counter;
            
            -- Link laptops
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, j);
                UPDATE laptops SET status = 'delivered' WHERE id = j;
            END LOOP;
            
            -- Pickup form
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5), pickup_date - INTERVAL '2 days',
                jsonb_build_object('number_of_laptops', bulk_size, 'assignment_type', 'bulk',
                    'pickup_date', to_char(pickup_date::date, 'YYYY-MM-DD')));
            
            -- Reception reports (all approved)
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id,
                    received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition,
                    status, approved_by, approved_at, created_at, updated_at)
                SELECT j, shipment_id_counter, client_id, tracking_num, 7 + ((i - 1) % 6),
                    pickup_date + INTERVAL '2 days', 'Bulk laptop received.',
                    '/uploads/reception/laptop' || j || '_serial.jpg',
                    '/uploads/reception/laptop' || j || '_external.jpg',
                    '/uploads/reception/laptop' || j || '_working.jpg',
                    'approved', 1 + ((i - 1) % 6), pickup_date + INTERVAL '2 days' + INTERVAL '6 hours',
                    pickup_date + INTERVAL '2 days', pickup_date + INTERVAL '2 days' + INTERVAL '6 hours'
                WHERE NOT EXISTS (SELECT 1 FROM reception_reports WHERE laptop_id = j);
            END LOOP;
        
        ELSIF i <= 25 THEN
            -- At warehouse (30% = 11 shipments)
            pickup_date := base_date - INTERVAL '5 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, arrived_warehouse_at,
                notes, created_at, updated_at)
            VALUES ('bulk_to_warehouse', client_id, NULL, 'at_warehouse'::shipment_status, bulk_size, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date::date, pickup_date, base_date - INTERVAL '2 days',
                'Bulk shipment at warehouse. Awaiting engineer assignments.',
                pickup_date - INTERVAL '2 days', NOW())
            RETURNING id INTO shipment_id_counter;
            
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, j);
                UPDATE laptops SET status = 'at_warehouse' WHERE id = j;
            END LOOP;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5), pickup_date - INTERVAL '2 days',
                jsonb_build_object('number_of_laptops', bulk_size, 'assignment_type', 'bulk'));
            
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO reception_reports (laptop_id, shipment_id, client_company_id, tracking_number, warehouse_user_id,
                    received_at, notes, photo_serial_number, photo_external_condition, photo_working_condition,
                    status, created_at, updated_at)
                SELECT j, shipment_id_counter, client_id, tracking_num, 7 + ((i - 1) % 6),
                    base_date - INTERVAL '2 days', 'Bulk laptop received.',
                    '/uploads/reception/laptop' || j || '_serial.jpg',
                    '/uploads/reception/laptop' || j || '_external.jpg',
                    '/uploads/reception/laptop' || j || '_working.jpg',
                    'pending_approval', base_date - INTERVAL '2 days', base_date - INTERVAL '2 days'
                WHERE NOT EXISTS (SELECT 1 FROM reception_reports WHERE laptop_id = j);
            END LOOP;
        
        ELSIF i <= 32 THEN
            -- In transit (20% = 7 shipments)
            pickup_date := base_date - INTERVAL '2 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, picked_up_at, notes, created_at, updated_at)
            VALUES ('bulk_to_warehouse', client_id, NULL, 'in_transit_to_warehouse'::shipment_status, bulk_size, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date::date, pickup_date,
                'Bulk shipment in transit to warehouse.',
                pickup_date - INTERVAL '2 days', NOW())
            RETURNING id INTO shipment_id_counter;
            
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, j);
                UPDATE laptops SET status = 'in_transit_to_warehouse' WHERE id = j;
            END LOOP;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 18 + ((i - 1) % 5), pickup_date - INTERVAL '2 days',
                jsonb_build_object('number_of_laptops', bulk_size, 'assignment_type', 'bulk'));
        
        ELSE
            -- Pending (10% = 3 shipments)
            pickup_date := base_date + INTERVAL '2 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, pickup_scheduled_date, notes, created_at, updated_at)
            VALUES ('bulk_to_warehouse', client_id, NULL, 'pending_pickup_from_client'::shipment_status, bulk_size, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date::date,
                'Bulk shipment pending pickup.',
                base_date - INTERVAL '1 day', NOW())
            RETURNING id INTO shipment_id_counter;
            
            FOREACH j IN ARRAY laptop_ids LOOP
                INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, j);
            END LOOP;
        END IF;
        
        -- shipment_id_counter is set by RETURNING clause, no need to increment
        jira_counter := jira_counter + 1;
    END LOOP;
    
    -- ============================================
    -- WAREHOUSE TO ENGINEER SHIPMENTS (~25 shipments)
    -- ============================================
    -- Distribution: ~60% delivered, ~25% in_transit, ~15% released
    
    FOR i IN 1..25 LOOP
        days_offset := (i - 1) * 8;
        base_date := NOW() - (days_offset || ' days')::INTERVAL;
        client_id := 1 + ((i - 1) % 15);
        engineer_id := 1 + ((i - 1) % 35);
        courier_idx := 1 + ((i - 1) % 5);
        tracking_num := tracking_prefixes[courier_idx] || (900000 + jira_counter)::TEXT;
        
        IF i <= 15 THEN
            -- Delivered (60% = 15 shipments)
            -- Average delivery time: 2.5-2.9 days from warehouse release
            delivery_days := 1.8 + (RANDOM() * 0.7); -- Random between 1.8-2.5 days
            pickup_date := base_date - (delivery_days + 1)::INT * INTERVAL '1 day';
            delivered_date := pickup_date + (delivery_days || ' days')::INTERVAL;
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, released_warehouse_at, delivered_at, notes, created_at, updated_at)
            VALUES ('warehouse_to_engineer', client_id, engineer_id, 'delivered'::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date, delivered_date,
                'Warehouse inventory delivered to engineer.',
                pickup_date - INTERVAL '1 day', delivered_date)
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
            UPDATE laptops SET status = 'delivered', software_engineer_id = engineer_id WHERE id = laptop_id_counter;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 1 + ((i - 1) % 6), pickup_date - INTERVAL '1 day',
                jsonb_build_object('engineer_id', engineer_id, 'laptop_id', laptop_id_counter));
            
            INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls)
            VALUES (shipment_id_counter, engineer_id, delivered_date, 'Delivered from warehouse.',
                ARRAY['/uploads/delivery/shipment' || LPAD(shipment_id_counter::TEXT, 3, '0') || '_photo1.jpg']);
        
        ELSIF i <= 21 THEN
            -- In transit (25% = 6 shipments)
            pickup_date := base_date - INTERVAL '2 days';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, released_warehouse_at, eta_to_engineer, notes, created_at, updated_at)
            VALUES ('warehouse_to_engineer', client_id, engineer_id, 'in_transit_to_engineer'::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date, base_date + INTERVAL '1 day',
                'In transit from warehouse to engineer.',
                pickup_date - INTERVAL '1 day', NOW())
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
            UPDATE laptops SET status = 'in_transit_to_engineer', software_engineer_id = engineer_id WHERE id = laptop_id_counter;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 1 + ((i - 1) % 6), pickup_date - INTERVAL '1 day',
                jsonb_build_object('engineer_id', engineer_id, 'laptop_id', laptop_id_counter));
        
        ELSE
            -- Released (15% = 4 shipments)
            pickup_date := base_date - INTERVAL '1 day';
            
            INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number,
                courier_name, tracking_number, released_warehouse_at, notes, created_at, updated_at)
            VALUES ('warehouse_to_engineer', client_id, engineer_id, 'released_from_warehouse'::shipment_status, 1, 'SCOP-' || jira_counter,
                courier_names[courier_idx], tracking_num, pickup_date,
                'Released from warehouse, awaiting courier pickup.',
                pickup_date - INTERVAL '1 day', NOW())
            RETURNING id INTO shipment_id_counter;
            
            INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES (shipment_id_counter, laptop_id_counter);
            UPDATE laptops SET status = 'in_transit_to_engineer', software_engineer_id = engineer_id WHERE id = laptop_id_counter;
            
            INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
            VALUES (shipment_id_counter, 1 + ((i - 1) % 6), pickup_date - INTERVAL '1 day',
                jsonb_build_object('engineer_id', engineer_id, 'laptop_id', laptop_id_counter));
        END IF;
        
        -- shipment_id_counter is set by RETURNING clause, no need to increment
        -- Increment and wrap around if exceeds max
        laptop_id_counter := laptop_id_counter + 1;
        IF laptop_id_counter > max_laptop_id THEN
            laptop_id_counter := 1;
        END IF;
        jira_counter := jira_counter + 1;
    END LOOP;
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'Error in shipments generation: %', SQLERRM;
        RAISE;
END $$;

-- ============================================
-- AUDIT LOGS (Sample recent activity)
-- ============================================
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
SELECT 
    1 + ((ROW_NUMBER() OVER ()) % 6),
    CASE (ROW_NUMBER() OVER ()) % 5
        WHEN 0 THEN 'shipment_created'
        WHEN 1 THEN 'status_updated'
        WHEN 2 THEN 'reception_report_created'
        WHEN 3 THEN 'pickup_form_submitted'
        ELSE 'delivery_form_created'
    END,
    'shipment',
    id,
    created_at + INTERVAL '1 hour',
    jsonb_build_object('action', 'shipment_created', 'jira_ticket', jira_ticket_number)
FROM shipments
ORDER BY created_at DESC
LIMIT 50;

-- ============================================
-- MAGIC LINKS (For active shipments)
-- ============================================
INSERT INTO magic_links (token, shipment_id, expires_at, used_at, user_id, created_at)
SELECT 
    'magic' || LPAD(id::TEXT, 6, '0') || MD5(id::TEXT || jira_ticket_number),
    id,
    NOW() + INTERVAL '7 days',
    CASE WHEN status = 'delivered' THEN created_at + INTERVAL '1 day' ELSE NULL END,
    1, -- Use first logistics user
    created_at
FROM shipments
WHERE status IN ('in_transit_to_engineer', 'delivered')
LIMIT 20;

-- ============================================
-- Summary & Verification
-- ============================================

SELECT '========================================' AS separator;
SELECT 'SHIPMENTS DATA v4.0 LOADED SUCCESSFULLY!' AS message;
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
SELECT 'AVERAGE DELIVERY TIME' AS section;
SELECT '---------------------' AS underline;
SELECT 
    ROUND(AVG(EXTRACT(EPOCH FROM (delivered_at - picked_up_at)) / 86400)::numeric, 2) as avg_delivery_days,
    COUNT(*) as delivered_count
FROM shipments
WHERE status = 'delivered' 
  AND picked_up_at IS NOT NULL 
  AND delivered_at IS NOT NULL;

SELECT '' AS blank;
SELECT 'SHIPMENTS WITH COMPLETE DATA' AS section;
SELECT '---------------------------' AS underline;
SELECT 
    COUNT(*) FILTER (WHERE courier_name IS NOT NULL AND tracking_number IS NOT NULL) as with_courier_tracking,
    COUNT(*) FILTER (WHERE courier_name IS NULL OR tracking_number IS NULL) as missing_courier_tracking,
    COUNT(*) as total
FROM shipments;

SELECT '' AS blank;
SELECT 'RECEPTION REPORTS STATUS' AS section;
SELECT '------------------------' AS underline;
SELECT status, COUNT(*) as count FROM reception_reports GROUP BY status ORDER BY count DESC;

SELECT '' AS blank;
SELECT '========================================' AS separator;
SELECT 'Complete! ~100 shipments loaded across all types.' AS summary1;
SELECT 'Average delivery time should be ~2.5-2.9 days.' AS summary2;
SELECT 'All shipments have complete field data.' AS summary3;
SELECT '========================================' AS separator;

