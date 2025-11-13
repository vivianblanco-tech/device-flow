-- Diagnostic query to identify why laptops aren't showing as available
-- for warehouse-to-engineer shipments

-- 1. Show all laptops and their current status
SELECT 
    l.id,
    l.serial_number,
    l.status,
    l.client_company_id,
    cc.name as company_name,
    CASE 
        WHEN l.status IN ('available', 'at_warehouse') THEN '✓'
        ELSE '✗'
    END as status_ok
FROM laptops l
LEFT JOIN client_companies cc ON cc.id = l.client_company_id
ORDER BY l.id;

-- 2. Check which laptops have reception reports
SELECT 
    l.id,
    l.serial_number,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM reception_reports rr
            JOIN shipment_laptops sl ON sl.shipment_id = rr.shipment_id
            WHERE sl.laptop_id = l.id
        ) THEN '✓ Has reception report'
        ELSE '✗ No reception report'
    END as reception_status
FROM laptops l
ORDER BY l.id;

-- 3. Check which laptops are in active shipments
SELECT 
    l.id,
    l.serial_number,
    s.id as shipment_id,
    s.status as shipment_status,
    s.shipment_type,
    CASE 
        WHEN s.status NOT IN ('delivered', 'at_warehouse') THEN '✗ In active shipment'
        ELSE '✓ Shipment inactive/at warehouse'
    END as active_shipment_check
FROM laptops l
JOIN shipment_laptops sl ON sl.laptop_id = l.id
JOIN shipments s ON s.id = sl.shipment_id
ORDER BY l.id, s.id;

-- 4. Complete availability check - this replicates the actual query logic
SELECT 
    l.id,
    l.serial_number,
    l.brand,
    l.model,
    l.status,
    -- Check 1: Status
    CASE WHEN l.status IN ('available', 'at_warehouse') THEN '✓' ELSE '✗' END as status_check,
    -- Check 2: Has reception report
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM reception_reports rr
            JOIN shipment_laptops sl ON sl.shipment_id = rr.shipment_id
            WHERE sl.laptop_id = l.id
        ) THEN '✓' 
        ELSE '✗' 
    END as reception_check,
    -- Check 3: Not in active shipment
    CASE 
        WHEN NOT EXISTS (
            SELECT 1 FROM shipment_laptops sl
            JOIN shipments s ON s.id = sl.shipment_id
            WHERE sl.laptop_id = l.id
              AND s.status NOT IN ('delivered', 'at_warehouse')
        ) THEN '✓'
        ELSE '✗'
    END as not_active_check,
    -- Overall availability
    CASE 
        WHEN l.status IN ('available', 'at_warehouse')
         AND EXISTS (
            SELECT 1 FROM reception_reports rr
            JOIN shipment_laptops sl ON sl.shipment_id = rr.shipment_id
            WHERE sl.laptop_id = l.id
         )
         AND NOT EXISTS (
            SELECT 1 FROM shipment_laptops sl
            JOIN shipments s ON s.id = sl.shipment_id
            WHERE sl.laptop_id = l.id
              AND s.status NOT IN ('delivered', 'at_warehouse')
         ) THEN '✓ AVAILABLE'
        ELSE '✗ NOT AVAILABLE'
    END as final_availability
FROM laptops l
ORDER BY l.id;

-- 5. Show laptops that SHOULD be available (for debugging)
SELECT 
    l.id,
    l.serial_number,
    l.brand,
    l.model,
    l.status,
    cc.name as client_company_name
FROM laptops l
LEFT JOIN client_companies cc ON cc.id = l.client_company_id
WHERE l.status IN ('available', 'at_warehouse')
  AND EXISTS (
      SELECT 1 FROM reception_reports rr
      JOIN shipments s ON s.id = rr.shipment_id
      JOIN shipment_laptops sl ON sl.shipment_id = s.id
      WHERE sl.laptop_id = l.id
  )
  AND NOT EXISTS (
      SELECT 1 FROM shipment_laptops sl
      JOIN shipments s ON s.id = sl.shipment_id
      WHERE sl.laptop_id = l.id
        AND s.status NOT IN ('delivered', 'at_warehouse')
  )
ORDER BY l.created_at DESC;

