# Test Schema Fix Summary

**Date**: November 16, 2025

## Issue Description

The test database schema was out of sync with the code, causing 21 tests to fail with the error:
```
pq: column "cpu" of relation "laptops" does not exist
```

## Root Cause

The test database (`laptop_tracking_test`) was at migration version 21, missing migration 22 which adds the `cpu` column to the `laptops` table.

## Solution Applied

Applied migration 22 to the test database using the following command:

```powershell
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $env:DATABASE_URL up
```

## Migration Details

**Migration**: `000022_add_cpu_to_laptops.up.sql`

**Changes**:
- Added `cpu` column to `laptops` table as `TEXT NOT NULL DEFAULT ''`
- Updated existing records to set `cpu = 'Unknown'` for empty values

## Verification

### Schema Version
- **Before**: Version 21
- **After**: Version 22

### Database Schema
Both test and dev databases now have the complete schema:

```sql
Column: cpu
Type: text
Nullable: not null
Default: ''::text
```

### Test Results
All previously failing tests now pass:

✅ **Previously Failing Tests** (now passing):
- `TestGetLaptopCountsByStatus`
- `TestGetAvailableLaptopCount`
- `TestGetDashboardStats`
- `TestGetAllLaptopsWithReceptionReportInfo`
- `TestGetAllLaptopsForLogisticsUsersIncludesReceptionReports`
- `TestGetAllLaptops`
- `TestGetAllLaptopsWithFilter`
- `TestSearchLaptops`
- `TestGetLaptopByID`
- `TestCreateLaptop`
- `TestCreateLaptopDuplicateSerial`
- `TestUpdateLaptop`
- `TestDeleteLaptop`
- `TestGetLaptopsByStatus`
- `TestGetAllLaptopsWithJoins`
- `TestGetLaptopByIDWithJoins`
- `TestGetAllLaptopsHandlesNullFields`
- `TestGetAllLaptops_WarehouseRoleFilter`
- `TestGetAllLaptops_NonWarehouseRoleFilter`
- `TestGetAllLaptops_WarehouseRoleWithStatusFilter`
- `TestGetAllLaptops_NoRoleFilter`

✅ **New Test** (from recent feature):
- `TestGetLaptopStatusesForNewLaptop` - Verifies only "Received at Warehouse" status is available for new laptops

### Test Suite Status

```bash
# Full models package test suite
go test ./internal/models/...
# Result: ok (1.883s)
```

## Database Status

### Test Database (`laptop_tracking_test`)
- ✅ Schema version: 22
- ✅ CPU column present
- ✅ All migrations applied
- ✅ All tests passing

### Development Database (`laptop_tracking_dev`)
- ✅ Schema version: 22
- ✅ CPU column present
- ✅ All migrations applied
- ✅ Application running correctly

## Commands for Future Reference

### Check Database Schema Version
```powershell
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "SELECT * FROM schema_migrations;"
```

### Check Laptop Table Schema
```powershell
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "\d laptops"
```

### Apply Pending Migrations
```powershell
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $env:DATABASE_URL up
```

### Run Tests with Test Database
```powershell
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
go test ./internal/models/...
```

## Prevention

To prevent this issue in the future:

1. **Always run migrations on both databases** when new migrations are added:
   - Development database: `laptop_tracking_dev`
   - Test database: `laptop_tracking_test`

2. **Verify schema version** before running tests:
   ```powershell
   docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "SELECT * FROM schema_migrations;"
   ```

3. **Add to CI/CD pipeline**: Ensure migrations run on test database before tests:
   ```yaml
   - name: Run test database migrations
     run: migrate -path migrations -database $TEST_DATABASE_URL up
   ```

4. **Include in setup scripts**: Update setup scripts to automatically sync both databases.

## Related Files

- Migration: `migrations/000022_add_cpu_to_laptops.up.sql`
- Test helper: `internal/database/testhelpers.go`
- Documentation: `docs/TEST_DATABASE_SETUP.md`

## Status

✅ **RESOLVED** - All tests passing, both databases in sync at migration version 22.

---

**Fixed by**: AI Assistant  
**Verified**: All model tests passing (1.883s)  
**Docker Environment**: PostgreSQL 15-alpine in Docker container

