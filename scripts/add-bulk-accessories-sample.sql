-- Add sample pickup forms with bulk dimensions and accessories
-- This file includes ALL fields that the pickup form template expects

-- Shipment 42: Bulk assignment with accessories (complete data)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(42, 1, NOW(), '{
  "company_name": "Tech Innovations Inc.",
  "contact_person": "Sarah Johnson",
  "contact_email": "sarah.johnson@techinnovations.com",
  "contact_phone": "+1-555-1234",
  "pickup_address": "1500 Innovation Drive, Building A",
  "pickup_city": "San Francisco",
  "pickup_state": "CA",
  "pickup_zip": "94103",
  "preferred_date": "2024-12-15",
  "pickup_time_slot": "morning",
  "num_laptops": 10,
  "number_of_boxes": 3,
  "assignment_type": "bulk",
  "bulk_length": 24.5,
  "bulk_width": 18.0,
  "bulk_height": 12.5,
  "bulk_weight": 45.0,
  "include_accessories": true,
  "accessories_description": "10x USB-C chargers, 5x laptop bags, 10x wireless mice, 3x docking stations",
  "special_instructions": "Fragile - Handle with care. Use freight elevator on west side."
}'::jsonb)
ON CONFLICT (shipment_id) DO UPDATE SET
  form_data = EXCLUDED.form_data,
  submitted_at = EXCLUDED.submitted_at;

-- Shipment 43: Another bulk assignment with different dimensions (complete data)
INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data) VALUES 
(43, 1, NOW(), '{
  "company_name": "Global Solutions Corp",
  "contact_person": "Michael Chen",
  "contact_email": "mchen@globalsolutions.com",
  "contact_phone": "+1-555-9999",
  "pickup_address": "2000 Enterprise Boulevard, Suite 500",
  "pickup_city": "Seattle",
  "pickup_state": "WA",
  "pickup_zip": "98101",
  "preferred_date": "2024-12-20",
  "pickup_time_slot": "afternoon",
  "num_laptops": 15,
  "number_of_boxes": 5,
  "assignment_type": "bulk",
  "bulk_length": 30.0,
  "bulk_width": 20.0,
  "bulk_height": 15.0,
  "bulk_weight": 62.5,
  "include_accessories": true,
  "accessories_description": "15x Power adapters, 15x Laptop sleeves, 10x External keyboards, 5x USB hubs",
  "special_instructions": "Schedule pickup after 2 PM. Contact security desk upon arrival."
}'::jsonb)
ON CONFLICT (shipment_id) DO UPDATE SET
  form_data = EXCLUDED.form_data,
  submitted_at = EXCLUDED.submitted_at;

