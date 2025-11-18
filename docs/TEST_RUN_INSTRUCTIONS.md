# Test Suite Execution Instructions

This document provides comprehensive instructions for running the full test suite correctly, including all prerequisites, setup steps, and troubleshooting information.

## Prerequisites

1. **Docker and Docker Compose** - The application and database must be running in Docker containers
2. **Go 1.24+** - Required for running Go tests
3. **PostgreSQL Test Database** - Must be accessible at `localhost:5432`

## Database Setup

### 1. Start Docker Containers

Ensure Docker containers are running:

```bash
docker-compose up -d
```

Verify containers are running:

```bash
docker-compose ps
```

### 2. Test Database Configuration

The test suite uses a separate test database with the following configuration:

- **Database Name**: `laptop_tracking_test`
- **Host**: `localhost`
- **Port**: `5432`
- **Username**: `postgres`
- **Password**: `password` (must match Docker PostgreSQL password)
- **SSL Mode**: `disable`

### 3. Environment Variable

Set the test database URL environment variable:

**Windows PowerShell:**
```powershell
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
```

**Linux/Mac (Bash):**
```bash
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
```

**Note**: The password must match the PostgreSQL password configured in your Docker Compose file. The default password used in the test helpers is `password`.

## Running Tests

### Critical: Sequential Execution Required

**IMPORTANT**: Tests MUST be run sequentially (`-p=1`) to avoid database conflicts. Running tests in parallel will cause race conditions and test failures.

### Full Test Suite

Run all tests sequentially:

**Windows PowerShell:**
```powershell
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
go test ./... -p=1 -v -race
```

**Linux/Mac (Bash):**
```bash
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
go test ./... -p=1 -v -race
```

### Specific Package Tests

Run tests for a specific package:

```bash
# Handlers package
go test ./internal/handlers -p=1 -v -race

# Models package
go test ./internal/models -p=1 -v -race

# Validator package
go test ./internal/validator -p=1 -v -race
```

### Specific Test Function

Run a specific test function:

```bash
go test ./internal/handlers -p=1 -v -race -run TestFunctionName
```

### With Coverage

Generate coverage report:

```bash
go test ./... -p=1 -v -race -cover
```

## Test Database Helper Configuration

The test database helper (`internal/database/testhelpers.go`) uses the following default configuration:

```go
dbURL := "postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
```

**Important**: If your Docker PostgreSQL password differs from `password`, you must either:
1. Set the `TEST_DATABASE_URL` environment variable (recommended)
2. Update the default `dbURL` in `internal/database/testhelpers.go`

## Common Issues and Solutions

### Issue: "password authentication failed"

**Cause**: Password mismatch between Docker PostgreSQL and test helper.

**Solution**: 
1. Check your Docker Compose PostgreSQL password
2. Set `TEST_DATABASE_URL` with the correct password
3. Or update the default password in `internal/database/testhelpers.go`

### Issue: "database does not exist"

**Cause**: Test database `laptop_tracking_test` hasn't been created.

**Solution**: The test database should be created automatically by migrations. Ensure:
1. Docker containers are running
2. Migrations have been applied
3. Database connection is correct

### Issue: Test failures due to race conditions

**Cause**: Tests running in parallel.

**Solution**: Always use `-p=1` flag to run tests sequentially:
```bash
go test ./... -p=1 -v -race
```

### Issue: "null value in column violates not-null constraint"

**Cause**: Test data doesn't match current database schema requirements.

**Solution**: Ensure test data includes all required fields:
- Reception reports require: `laptop_id`, `warehouse_user_id`, `photo_serial_number`, `photo_external_condition`, `photo_working_condition`, `status`
- Laptops require: `client_company_id` (for most operations)
- Forms require: `laptop_brand` (now required field)

### Issue: Date validation failures

**Cause**: Hardcoded dates in tests that are now in the past.

**Solution**: Use dynamic future dates:
```go
pickupDate := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
```

### Issue: Template not found errors

**Cause**: Tests using minimal templates instead of loaded templates.

**Solution**: Use `loadTestTemplates(t)` helper function:
```go
templates := loadTestTemplates(t)
handler := NewHandler(db, templates, nil)
```

## Test Data Requirements

### Reception Reports

When creating reception reports in tests, use the new schema:

```go
_, err = db.ExecContext(ctx,
    `INSERT INTO reception_reports (
        laptop_id, warehouse_user_id, 
        photo_serial_number, photo_external_condition, photo_working_condition, 
        status, shipment_id, client_company_id, tracking_number, 
        notes, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, NULL, NULL, NULL, $7, NOW(), NOW())`,
    laptopID, warehouseUserID, 
    "http://example.com/serial.jpg", "http://example.com/ext.jpg", "http://example.com/work.jpg",
    "approved", "Test notes",
)
```

**Key Points**:
- `laptop_id` is required (not `shipment_id`)
- Photo fields are required
- `status` should be `"approved"` for laptops that need to be set to `available`
- `shipment_id`, `client_company_id`, and `tracking_number` can be NULL

### Laptops

When creating laptops in tests:

```go
laptop := &models.Laptop{
    SerialNumber:    "SN-TEST-001",
    Brand:           "Dell",           // Required
    Model:           "XPS 15",          // Required
    CPU:             "Intel Core i7",  // Required
    RAMGB:           "16GB",            // Required
    SSDGB:           "512GB",           // Required
    Status:          models.LaptopStatusDelivered,
    ClientCompanyID: &companyID,       // Required for most operations
}
```

### Form Submissions

When testing form submissions, include all required fields:

```go
formData := url.Values{}
formData.Set("laptop_serial_number", "SN123456789")
formData.Set("laptop_brand", "Dell")        // Required
formData.Set("laptop_model", "XPS 15")     // Required
formData.Set("laptop_cpu", "Intel Core i7") // Required
formData.Set("laptop_ram_gb", "16GB")       // Required
formData.Set("laptop_ssd_gb", "512GB")      // Required
formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
// ... other fields
```

## Schema Changes Impact

### Reception Reports Migration

The reception reports table was migrated from shipment-based to laptop-based:

**Old Schema** (deprecated):
- Primary key: `shipment_id`
- Checked via: `JOIN shipment_laptops`

**New Schema** (current):
- Primary key: `laptop_id`
- Direct relationship: `reception_reports.laptop_id = laptops.id`

**Impact on Tests**:
- All reception report queries must use `laptop_id`
- Handler functions updated to check reception reports directly by `laptop_id`
- Test data creation must use new schema

## Best Practices

1. **Always run tests sequentially**: Use `-p=1` flag
2. **Set environment variable**: Always set `TEST_DATABASE_URL` before running tests
3. **Use dynamic dates**: Avoid hardcoded dates that become invalid over time
4. **Include all required fields**: Check current schema requirements before creating test data
5. **Use approved reception reports**: For laptops that need `available` status, create approved reception reports
6. **Clean up test data**: Tests should clean up after themselves (handled by `SetupTestDB`)

## Verification

After running tests, verify all packages pass:

```bash
go test ./... -p=1 -v -race 2>&1 | grep -E "(^ok|^FAIL)"
```

Expected output:
```
ok  	github.com/yourusername/laptop-tracking-system/internal/auth
ok  	github.com/yourusername/laptop-tracking-system/internal/config
ok  	github.com/yourusername/laptop-tracking-system/internal/database
ok  	github.com/yourusername/laptop-tracking-system/internal/email
ok  	github.com/yourusername/laptop-tracking-system/internal/handlers
ok  	github.com/yourusername/laptop-tracking-system/internal/jira
ok  	github.com/yourusername/laptop-tracking-system/internal/models
ok  	github.com/yourusername/laptop-tracking-system/internal/validator
ok  	github.com/yourusername/laptop-tracking-system/internal/views
ok  	github.com/yourusername/laptop-tracking-system/tests/integration
ok  	github.com/yourusername/laptop-tracking-system/tests/unit
```

## Quick Reference

### Essential Commands

```bash
# Set environment variable (PowerShell)
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"

# Set environment variable (Bash)
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"

# Run all tests sequentially
go test ./... -p=1 -v -race

# Run specific package
go test ./internal/handlers -p=1 -v -race

# Run specific test
go test ./internal/handlers -p=1 -v -race -run TestFunctionName

# Run with coverage
go test ./... -p=1 -v -race -cover
```

## Additional Resources

- `docs/TEST_DATABASE_SETUP.md` - Detailed database setup instructions
- `docs/TESTING_BEST_PRACTICES.md` - General testing guidelines
- `internal/database/testhelpers.go` - Test database helper functions

