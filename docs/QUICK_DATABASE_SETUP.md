# Quick Database Setup - Choose Your Method

## 🚀 Method 1: Docker (Easiest - Recommended)

### Step 1: Start Docker Desktop
1. Open Docker Desktop application
2. Wait for it to fully start (icon in system tray should be green)

### Step 2: Start PostgreSQL Container
```powershell
# Navigate to project directory
cd "e:\Cursor Projects\BDH"

# Start PostgreSQL
docker-compose up -d postgres

# Wait ~10 seconds for PostgreSQL to initialize
Start-Sleep -Seconds 10

# Verify it's running
docker ps
```

### Step 3: Create Test Database
```powershell
# Create test database
docker exec -it laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
```

### Step 4: Run Migrations
```powershell
# Development database
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/laptop_tracking_dev?sslmode=disable" up

# Test database
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable" up
```

### Step 5: Run Tests
```powershell
# Run all tests including integration tests
go test ./... -v
```

---

## 💻 Method 2: Local PostgreSQL (If Already Installed)

### Step 1: Check PostgreSQL is Running
```powershell
# Check service status
Get-Service -Name postgresql*

# If not running, start it
Start-Service postgresql*

# Test connection
psql -U postgres -c "SELECT version();"
```

### Step 2: Run Automated Setup Script
```powershell
# Run the setup script
.\scripts\setup-database.ps1
```

This will:
- Create development database (`laptop_tracking_dev`)
- Create test database (`laptop_tracking_test`)
- Run all migrations
- Set up .env file

### Step 3: Run Tests
```powershell
go test ./... -v
```

---

## 🔧 Quick Commands Reference

### Docker Commands
```powershell
# Start database
docker-compose up -d postgres

# Stop database
docker-compose stop postgres

# View logs
docker-compose logs postgres

# Access psql
docker exec -it laptop-tracking-db psql -U postgres -d laptop_tracking_dev

# Remove everything (WARNING: deletes data)
docker-compose down -v
```

### Test Commands
```powershell
# All tests (needs database)
go test ./... -v

# Unit tests only (no database needed)
go test ./... -v -short

# Specific package
go test ./internal/models/... -v
go test ./internal/auth/... -v
go test ./internal/validator/... -v
```

### Database Info
- **Host**: localhost
- **Port**: 5432
- **User**: postgres
- **Password**: postgres
- **Dev Database**: laptop_tracking_dev
- **Test Database**: laptop_tracking_test

---

## ✅ Verification

After setup, verify everything works:

```powershell
# 1. Check databases exist
docker exec -it laptop-tracking-db psql -U postgres -c "\l"

# 2. Check tables in dev database
docker exec -it laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "\dt"

# 3. Run tests
go test ./... -v

# 4. Check test summary
Get-Content docs\TEST_STATUS_SUMMARY.md
```

---

## 🆘 Troubleshooting

### Docker Desktop Not Running
- **Solution**: Open Docker Desktop and wait for it to start

### Port 5432 Already in Use
```powershell
# Stop local PostgreSQL
Stop-Service postgresql*

# Or change Docker port in docker-compose.yml to 5433
```

### Migration Command Not Found
```powershell
# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Add to PATH: %USERPROFILE%\go\bin
```

### Password Issues
- Docker uses: `postgres/postgres`
- Check your `.env` file matches

---

## 📊 Expected Test Results

After setup, you should see:
- ✅ Phase 1: 133/133 tests passing
- ✅ Phase 2: 22/22 tests passing (5 unit + 17 integration)
- ✅ Phase 3: 50/50 tests passing (48 validator + 2 handler)
- ✅ **Total: 205+ tests passing**

Integration tests will no longer skip once database is connected!














