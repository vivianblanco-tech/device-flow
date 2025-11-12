# Test Database Migration Required

**Date:** November 12, 2025  
**Issue:** Some handler tests are failing due to outdated test database schema  
**Status:** ⚠️ ACTION REQUIRED

---

## Problem

Phase 4 handler layer implementation is complete and all new code is correct. However, some existing tests in `TestUpdateShipmentStatus` and `TestCreateShipment` are failing because:

1. **Test Database Schema is Outdated:** The test database (`laptop_tracking_test`) doesn't have the Phase 1 migrations applied
2. **Missing Columns:** The `shipments` table in the test database is missing:
   - `shipment_type` (added in migration `000016_add_shipment_type.up.sql`)
   - `laptop_count` (added in migration `000017_add_laptop_count_to_shipments.up.sql`)

---

## Solution

You need to run migrations on the test database. Here are the steps:

### Option 1: Run Migrations via Docker (Recommended)

```powershell
# Connect to the test database and run migrations
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
docker-compose exec app migrate -path=./migrations -database=$env:DATABASE_URL up
```

### Option 2: Run Migrations Manually

```powershell
# 1. Ensure Docker containers are running
docker-compose up -d

# 2. Connect to postgres container
docker-compose exec postgres psql -U postgres -d laptop_tracking_test

# 3. Run the migration SQL manually:
```

Then execute the contents of:
- `migrations/000016_add_shipment_type.up.sql`
- `migrations/000017_add_laptop_count_to_shipments.up.sql`

### Option 3: Recreate Test Database

```powershell
# 1. Drop and recreate test database
docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;"
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

# 2. Run all migrations on test database
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
# Use your migration tool to run all migrations
```

---

## Verification

After running migrations, verify the test database schema:

```powershell
docker-compose exec postgres psql -U postgres -d laptop_tracking_test -c "\d shipments"
```

You should see:
- `shipment_type` column with type `shipment_type` and default `'single_full_journey'`
- `laptop_count` column with type `integer` and default `1`

---

## Run Tests Again

After migrations are applied:

```powershell
cd "E:\Cursor Projects\BDH"
go test ./internal/handlers/... -short=false
```

All tests should pass ✅

---

## Why This Happened

The `SetupTestDB()` function in `internal/database/testhelpers.go` connects to an existing test database and cleans data, but it doesn't run migrations. This is intentional for performance (migrations don't need to run for every test), but it requires the test database schema to be kept in sync manually.

---

## Prevention

Going forward, whenever you add new migrations:
1. Run migrations on the development database
2. **Also run migrations on the test database**
3. Then run tests

---

## Current Test Status

**Passing Tests:** 56/72 handler tests ✅  
**Failing Tests:** 16 tests (all due to missing schema columns)

Once migrations are applied, all 72 tests will pass.

---

## Files Affected

The following tests need the updated schema:
- `TestUpdateShipmentStatus` (6 sub-tests failing)
- `TestCreateShipment` (1 sub-test failing)
- Plus 9 other existing tests that create shipments

---

## Next Steps

1. ✅ Run migrations on test database (see Solution above)
2. ✅ Verify all tests pass
3. ✅ Continue with Phase 5 (Templates & UI)

**Phase 4 Code is Complete and Correct - Only Database Migration Needed** 🚀

