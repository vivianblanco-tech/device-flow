# Test Data for Laptop Tracking System

This document describes the test data that has been added to the database and provides helpful queries for testing and development.

## Overview

The `create-test-data.sql` script populates the following tables with realistic test data:

1. **Client Companies** - 5 companies
2. **Software Engineers** - 10 engineers  
3. **Laptops** - 15 laptops from various brands
4. **Shipments** - 13 shipments in various statuses
5. **Shipment-Laptop Links** - Associations between shipments and laptops

## Running the Script

To populate the database with test data:

```powershell
# From the project root directory
cd "E:\Cursor Projects\BDH"
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

## Test Data Summary

### Client Companies (5 total)

| Name | Contact Email |
|------|---------------|
| TechCorp Solutions | contact@techcorp-solutions.com |
| Global Innovations Inc | info@globalinnovations.com |
| Digital Dynamics LLC | support@digitaldynamics.com |
| CloudFirst Technologies | hello@cloudfirst.tech |
| NextGen Software Group | contact@nextgensoftware.com |

### Software Engineers (10 total)

**Engineers with Confirmed Addresses (7):**
- Alex Thompson - alex.thompson@email.com
- Maria Garcia - maria.garcia@email.com
- James Wilson - james.wilson@email.com
- Emily Chen - emily.chen@email.com
- Robert Martinez - robert.martinez@email.com
- Sarah Anderson - sarah.anderson@email.com
- David Kim - david.kim@email.com

**Engineers with Unconfirmed Addresses (3):**
- Jessica Taylor - jessica.taylor@email.com
- Michael Brown - michael.brown@email.com
- Lisa Johnson - lisa.johnson@email.com

### Laptops (15 total)

**Brands:**
- Dell: 5 laptops
- Lenovo: 3 laptops
- HP: 2 laptops
- Apple: 2 laptops
- ASUS: 2 laptops
- Microsoft: 1 laptop

**Status Distribution:**
- Delivered: 5
- At Warehouse: 3
- Available: 3
- In Transit to Warehouse: 2
- In Transit to Engineer: 2

### Shipments (13 total)

**Status Distribution:**
- Delivered: 5
- At Warehouse: 3
- In Transit to Engineer: 2
- In Transit to Warehouse: 1
- Pending Pickup from Client: 1
- Picked Up from Client: 1

## Useful Test Queries

### View All Client Companies
```sql
SELECT id, name, 
       SUBSTRING(contact_info FROM 'Email: ([^\n]+)') as email
FROM client_companies 
ORDER BY name;
```

### View All Software Engineers
```sql
SELECT id, name, email, address_confirmed, created_at
FROM software_engineers
ORDER BY name;
```

### View All Laptops with Details
```sql
SELECT id, serial_number, brand, model, status, created_at
FROM laptops
ORDER BY brand, model;
```

### View All Shipments with Client and Engineer Info
```sql
SELECT 
    s.id,
    cc.name as client_company,
    se.name as software_engineer,
    s.status,
    s.courier_name,
    s.tracking_number,
    s.created_at
FROM shipments s
JOIN client_companies cc ON s.client_company_id = cc.id
LEFT JOIN software_engineers se ON s.software_engineer_id = se.id
ORDER BY s.created_at DESC;
```

### View Shipments with Laptop Details
```sql
SELECT 
    s.id as shipment_id,
    s.status as shipment_status,
    l.serial_number,
    l.brand,
    l.model,
    l.status as laptop_status,
    cc.name as client_company,
    se.name as engineer
FROM shipments s
JOIN shipment_laptops sl ON s.id = sl.shipment_id
JOIN laptops l ON sl.laptop_id = l.id
JOIN client_companies cc ON s.client_company_id = cc.id
LEFT JOIN software_engineers se ON s.software_engineer_id = se.id
ORDER BY s.id;
```

### Find Laptops Currently at Warehouse
```sql
SELECT 
    l.id,
    l.serial_number,
    l.brand,
    l.model,
    l.specs,
    s.id as shipment_id,
    cc.name as client_company
FROM laptops l
JOIN shipment_laptops sl ON l.id = sl.laptop_id
JOIN shipments s ON sl.shipment_id = s.id
JOIN client_companies cc ON s.client_company_id = cc.id
WHERE l.status = 'at_warehouse'
ORDER BY l.brand, l.model;
```

### Find Available Laptops (Not Yet in Shipment)
```sql
SELECT 
    l.id,
    l.serial_number,
    l.brand,
    l.model,
    l.specs,
    l.status
FROM laptops l
WHERE l.status = 'available'
ORDER BY l.brand, l.model;
```

### Find Shipments Awaiting Engineer Assignment
```sql
SELECT 
    s.id,
    s.status,
    cc.name as client_company,
    s.courier_name,
    s.tracking_number,
    COUNT(l.id) as laptop_count
FROM shipments s
JOIN client_companies cc ON s.client_company_id = cc.id
LEFT JOIN shipment_laptops sl ON s.id = sl.shipment_id
LEFT JOIN laptops l ON sl.laptop_id = l.id
WHERE s.software_engineer_id IS NULL
GROUP BY s.id, s.status, cc.name, s.courier_name, s.tracking_number
ORDER BY s.created_at;
```

### Find Engineers with Delivered Laptops
```sql
SELECT 
    se.id,
    se.name,
    se.email,
    COUNT(DISTINCT s.id) as shipment_count,
    COUNT(l.id) as laptop_count
FROM software_engineers se
JOIN shipments s ON se.id = s.software_engineer_id
LEFT JOIN shipment_laptops sl ON s.id = sl.shipment_id
LEFT JOIN laptops l ON sl.laptop_id = l.id
WHERE s.status = 'delivered'
GROUP BY se.id, se.name, se.email
ORDER BY laptop_count DESC;
```

### Shipment Status History (Simulated)
```sql
SELECT 
    s.id as shipment_id,
    cc.name as client,
    se.name as engineer,
    s.status,
    s.pickup_scheduled_date,
    s.picked_up_at,
    s.arrived_warehouse_at,
    s.released_warehouse_at,
    s.delivered_at
FROM shipments s
JOIN client_companies cc ON s.client_company_id = cc.id
LEFT JOIN software_engineers se ON s.software_engineer_id = se.id
ORDER BY s.created_at DESC;
```

## Test Scenarios

The test data includes various scenarios for comprehensive testing:

### Scenario 1: Complete Delivery Flow
- **Shipment ID: 16** - Dell XPS 13 delivered to Alex Thompson
- Status progression: pending → picked_up → at_warehouse → released → delivered

### Scenario 2: In Transit to Engineer
- **Shipment ID: 18** - Lenovo X1 Carbon to James Wilson
- Currently in transit, expected delivery soon

### Scenario 3: Awaiting Assignment
- **Shipment ID: 22** - Dell XPS 15 at warehouse
- No engineer assigned yet, ready for assignment

### Scenario 4: Pending Pickup from Client
- **Shipment ID: 27** - Scheduled for future pickup
- Tests the early stages of the workflow

### Scenario 5: Multiple Deliveries to Same Engineer
- Engineers can receive multiple laptops over time
- Sarah Anderson has received a MacBook Pro 14"

## Clearing Test Data

To remove all test data (WARNING: This will delete all data):

```sql
DELETE FROM shipment_laptops;
DELETE FROM shipments;
DELETE FROM laptops;
DELETE FROM software_engineers;
DELETE FROM client_companies;
```

Or to clear and re-populate:

```powershell
cd "E:\Cursor Projects\BDH"
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "DELETE FROM shipment_laptops; DELETE FROM shipments; DELETE FROM laptops; DELETE FROM software_engineers; DELETE FROM client_companies;"
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

## Notes

- The script uses `ON CONFLICT` clauses to prevent duplicate entries when run multiple times
- Timestamps are set relative to NOW() to simulate realistic timelines
- Serial numbers follow a consistent naming pattern: `{BRAND}-{MODEL}-SN{NUMBER}`
- All data is entirely fictional and generated for testing purposes

## Integration with Application

This test data is designed to work with the laptop tracking application's features:

1. **Client Portal** - Clients can log in and create shipment requests
2. **Warehouse Management** - View laptops at warehouse, assign to engineers
3. **Engineer Portal** - Engineers can confirm addresses and track deliveries
4. **Admin Dashboard** - View all shipments, manage users and companies
5. **Status Tracking** - Track shipments through the entire delivery pipeline

## Maintenance

To update or modify test data:

1. Edit the `scripts/create-test-data.sql` file
2. Clear existing test data
3. Re-run the script

For incremental updates, use standard SQL UPDATE/INSERT statements or modify the script and re-run (ON CONFLICT clauses will handle updates).

