# Laptop Availability Investigation

## Issue
The warehouse-to-engineer form shows "0 laptop(s) available for shipment" even when laptops exist in the database.

## Root Cause Analysis

### The Query Logic
The `WarehouseToEngineerFormPage` handler uses this query to find available laptops:

```sql
SELECT DISTINCT l.id, l.serial_number, l.sku, l.brand, l.model, l.specs,
       l.status, l.client_company_id, l.software_engineer_id,
       l.created_at, l.updated_at,
       cc.name as client_company_name
FROM laptops l
LEFT JOIN client_companies cc ON cc.id = l.client_company_id
WHERE l.status IN ('available', 'at_warehouse')
  -- Must have a reception report
  AND EXISTS (
      SELECT 1 FROM reception_reports rr
      JOIN shipments s ON s.id = rr.shipment_id
      JOIN shipment_laptops sl ON sl.shipment_id = s.id
      WHERE sl.laptop_id = l.id
  )
  -- Must not be in any active shipment (except bulk shipments at warehouse)
  AND NOT EXISTS (
      SELECT 1 FROM shipment_laptops sl
      JOIN shipments s ON s.id = sl.shipment_id
      WHERE sl.laptop_id = l.id
        AND s.status NOT IN ('delivered', 'at_warehouse')
  )
ORDER BY l.created_at DESC
```

### Three Conditions Required
For a laptop to be available, ALL three must be true:

1. ✅ **Status Check**: `l.status IN ('available', 'at_warehouse')`
2. ✅ **Reception Report**: Must have a reception report
3. ✅ **Not in Active Shipment**: Must NOT be in shipments with status other than 'delivered' or 'at_warehouse'

## Potential Causes

### 1. Laptop Status Not Set Correctly
**Problem**: When laptops arrive via bulk shipment and reception report is created, the individual laptop status might not be updated to 'at_warehouse' or 'available'.

**Solution**: Ensure reception report handler updates laptop statuses:

```go
// In reception report handler
_, err = tx.ExecContext(ctx,
    `UPDATE laptops 
     SET status = $1, updated_at = $2
     WHERE id = ANY($3)`,
    models.LaptopStatusAtWarehouse,
    time.Now(),
    pq.Array(laptopIDs),
)
```

### 2. Missing Reception Reports
**Problem**: Laptops were added to inventory but reception reports were never created for the shipments they arrived in.

**Solution**: 
- For existing data, manually create reception reports for historical shipments
- Ensure bulk shipments require reception reports before laptops can be used
- Add validation to prevent laptop status updates without reception reports

### 3. Laptops Still Linked to Active Shipments
**Problem**: Laptops are still associated with the original bulk shipment that hasn't been marked as 'at_warehouse' or 'delivered'.

**Solution**: Ensure bulk shipment status is updated when reception report is created:

```go
// Update shipment status to 'at_warehouse' when reception report is created
_, err = tx.ExecContext(ctx,
    `UPDATE shipments 
     SET status = $1, updated_at = $2
     WHERE id = $3`,
    models.ShipmentStatusAtWarehouse,
    time.Now(),
    shipmentID,
)
```

## Diagnostic Steps

### Run the Diagnostic SQL
Execute `scripts/diagnose-laptop-availability.sql` against your database to identify which condition is failing:

```powershell
# From Docker container
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/diagnose-laptop-availability.sql
```

### Check Each Condition

**Query 1**: Shows all laptops and status
**Query 2**: Shows which laptops have reception reports  
**Query 3**: Shows which laptops are in active shipments
**Query 4**: Complete availability check with individual condition results
**Query 5**: Final list of available laptops (should match the form)

## Recommended Solutions

### Short-term Fix (Data Correction)
Run this SQL to update existing laptops that should be available:

```sql
-- Update laptop statuses for laptops with reception reports
UPDATE laptops l
SET status = 'at_warehouse', updated_at = NOW()
WHERE l.id IN (
    SELECT DISTINCT sl.laptop_id
    FROM shipment_laptops sl
    JOIN reception_reports rr ON rr.shipment_id = sl.shipment_id
    JOIN shipments s ON s.id = sl.shipment_id
    WHERE s.shipment_type = 'bulk_to_warehouse'
      AND s.status = 'at_warehouse'
      AND l.status NOT IN ('at_warehouse', 'available', 'delivered')
);

-- Ensure bulk shipments with reception reports are marked 'at_warehouse'
UPDATE shipments s
SET status = 'at_warehouse', updated_at = NOW()
WHERE s.shipment_type = 'bulk_to_warehouse'
  AND EXISTS (
      SELECT 1 FROM reception_reports rr
      WHERE rr.shipment_id = s.id
  )
  AND s.status NOT IN ('at_warehouse', 'delivered');
```

### Long-term Fix (Code Changes)

#### 1. Update Reception Report Handler
Ensure it updates laptop statuses when creating reception reports:

```go
// File: internal/handlers/reception_report.go

// After creating reception report, update all laptop statuses
for _, detail := range reportDetails {
    _, err = tx.ExecContext(r.Context(),
        `UPDATE laptops 
         SET status = $1, updated_at = $2
         WHERE serial_number = $3`,
        models.LaptopStatusAtWarehouse,
        time.Now(),
        detail.SerialNumber,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to update laptop status: %w", err)
    }
}
```

#### 2. Add Validation
Add a validation function to check laptop availability before displaying the form:

```go
func (h *PickupFormHandler) getAvailableLaptopsCount(ctx context.Context) (int, error) {
    var count int
    err := h.DB.QueryRowContext(ctx, `
        SELECT COUNT(DISTINCT l.id)
        FROM laptops l
        WHERE l.status IN ('available', 'at_warehouse')
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
          )
    `).Scan(&count)
    return count, err
}
```

## Testing

After applying fixes, verify:

1. ✅ Laptops show in the warehouse-to-engineer form dropdown
2. ✅ Count displays correctly: "X laptop(s) available for shipment"
3. ✅ Only truly available laptops appear (not in active shipments)
4. ✅ New bulk shipments → reception reports → automatically make laptops available

## Next Steps

1. Run diagnostic SQL to identify the specific issue
2. Apply appropriate short-term fix (SQL update)
3. Implement long-term code fixes if needed
4. Add test coverage for laptop availability logic
5. Document the workflow in user guides

