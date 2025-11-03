# Quick Start - Test Data

## TL;DR

```powershell
# Navigate to project directory
cd "E:\Cursor Projects\BDH"

# Populate database with test data
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev

# Verify test data
.\scripts\verify-test-data.ps1
```

## What Was Created

✅ **5 Client Companies**  
✅ **10 Software Engineers** (7 with confirmed addresses, 3 pending)  
✅ **15 Laptops** (Dell, Lenovo, HP, Apple, ASUS, Microsoft)  
✅ **13 Shipments** (various statuses throughout the pipeline)  
✅ **32 Shipment-Laptop Links**  

## Test Data Highlights

### Realistic Scenarios

- **Complete deliveries**: 5 shipments fully delivered to engineers
- **In progress**: 2 shipments in transit to engineers
- **At warehouse**: 3 shipments waiting for engineer assignment
- **In transit**: 1 shipment on the way to warehouse
- **Pending**: 1 shipment scheduled for future pickup
- **Available laptops**: 3 laptops not yet in any shipment

### Sample Test Cases

1. **Alex Thompson** received a Dell XPS 13 ✅
2. **Maria Garcia** received a Dell Latitude 7420 ✅
3. **James Wilson** has a Lenovo X1 Carbon in transit 🚚
4. **3 shipments** at warehouse awaiting engineer assignment 📦
5. **3 available laptops** ready to be added to new shipments 💻

## Quick Commands

### View All Shipments
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT id, status FROM shipments ORDER BY id;"
```

### View All Laptops
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT serial_number, brand, model, status FROM laptops ORDER BY brand;"
```

### View Engineers with Laptops
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT se.name, COUNT(*) as laptop_count FROM software_engineers se JOIN shipments s ON se.id = s.software_engineer_id WHERE s.status = 'delivered' GROUP BY se.name;"
```

### Clear All Test Data
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "DELETE FROM shipment_laptops; DELETE FROM shipments; DELETE FROM laptops; DELETE FROM software_engineers; DELETE FROM client_companies;"
```

### Reload Fresh Test Data
```powershell
# Clear existing
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "DELETE FROM shipment_laptops; DELETE FROM shipments; DELETE FROM laptops; DELETE FROM software_engineers; DELETE FROM client_companies;"

# Reload
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

## Files Created

- `scripts/create-test-data.sql` - Main SQL script with all test data
- `scripts/verify-test-data.ps1` - PowerShell script to verify data
- `scripts/TEST_DATA_README.md` - Comprehensive documentation
- `scripts/QUICK_START_TEST_DATA.md` - This quick reference guide

## Test Data Design

The test data is designed to cover:

✅ Multiple client companies  
✅ Engineers at various stages (confirmed/unconfirmed addresses)  
✅ Laptops from different brands with various specs  
✅ Shipments at all stages of the delivery pipeline  
✅ Realistic timestamps showing progression over time  
✅ Proper foreign key relationships between all tables  

## Integration Testing

Use this data to test:

- Client portal functionality
- Warehouse management features
- Engineer address confirmation workflow
- Shipment status tracking
- Admin dashboard views
- Email notifications
- Reporting and analytics

## Support

For detailed documentation and SQL queries, see:
- `scripts/TEST_DATA_README.md` - Full documentation with sample queries

## Status Summary

| Entity | Count | Status |
|--------|-------|--------|
| Client Companies | 5 | ✅ Created |
| Software Engineers | 10 | ✅ Created |
| Laptops | 15 | ✅ Created |
| Shipments | 13 | ✅ Created |
| Shipment-Laptop Links | 32 | ✅ Created |

**Last Updated**: 2025-11-03

