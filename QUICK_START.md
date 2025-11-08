# Quick Start Guide

## 🚀 One-Command Startup

```powershell
.\scripts\start-with-data.ps1
```

This will:
- ✅ Start all Docker containers
- ✅ Automatically load sample data if database is empty
- ✅ Show you the login credentials

## 📝 Common Commands

### Start Application (Data Persists)
```powershell
# Normal start
docker compose up -d

# With rebuild
docker compose up -d --build

# Or use the convenience script
.\scripts\start-with-data.ps1
.\scripts\start-with-data.ps1 -Build  # with rebuild
```

### Stop Application (Data Persists)
```powershell
# ✅ Safe - keeps all data
docker compose down

# ❌ DANGER - deletes all data
docker compose down -v
```

### Fresh Start (Clean Slate)
```powershell
# Interactive prompt for safety
.\scripts\start-with-data.ps1 -Fresh
```

### View Logs
```powershell
# All services
docker compose logs -f

# Specific service
docker compose logs -f app
docker compose logs -f postgres
```

## 🔐 Test User Credentials

All test users have the password: **`Test123!`**

| Email | Role | Access Level |
|-------|------|--------------|
| `logistics@bairesdev.com` | Logistics | Full Access |
| `warehouse@bairesdev.com` | Warehouse | Medium Access |
| `client@bairesdev.com` | Client | Limited Access |
| `pm@bairesdev.com` | Project Manager | Read-Only |

## 💾 Database Management

### Create Backup
```powershell
.\scripts\backup-db.ps1
```

### Restore from Backup
```powershell
# Interactive selection
.\scripts\restore-db.ps1

# Specific file
.\scripts\restore-db.ps1 -BackupFile "db-backups/laptop_tracking_backup_20251108-153045.sql"
```

### Load Sample Data Manually
```powershell
# Test users only
Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev

# All test data (companies, laptops, shipments, etc.)
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

## 🌐 Service URLs

- **Application**: http://localhost:8080
- **MailHog (Email Testing)**: http://localhost:8025
- **Database**: localhost:5432

## 🔧 Troubleshooting

### Data Disappeared After Restart?

**Problem**: Used `docker compose down -v` which deletes volumes.

**Solution**:
```powershell
# Always use without -v flag
docker compose down

# If data is lost, restore from backup
.\scripts\restore-db.ps1
```

### Database is Empty After Restart?

**Solution**:
```powershell
# Run the initialization script
.\scripts\init-db-if-empty.ps1

# Or restart with auto-init
.\scripts\start-with-data.ps1
```

### Port Already in Use?

```powershell
# Stop all containers
docker compose down

# Check what's using the port (8080, 5432, 1025, 8025)
netstat -ano | findstr :8080

# Kill the process (replace <PID> with actual process ID)
taskkill /PID <PID> /F
```

### Can't Connect to Database?

```powershell
# Check if container is running
docker ps

# Check database logs
docker compose logs postgres

# Verify database is ready
docker exec laptop-tracking-db pg_isready -U postgres
```

### Migrations Not Running?

```powershell
# Restart the app container
docker compose restart app

# Check app logs
docker compose logs app
```

## 📊 Database Queries

### Check User Count
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"
```

### List All Users
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT id, email, role FROM users ORDER BY role;"
```

### Check Shipments
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT id, status, courier_name FROM shipments LIMIT 10;"
```

## 🛡️ Data Persistence - What's Safe?

### ✅ Safe Commands (Data Preserved)
- `docker compose up -d`
- `docker compose down` (without -v)
- `docker compose restart`
- `docker compose stop`
- `docker compose start`

### ❌ Dangerous Commands (Data Lost)
- `docker compose down -v` or `docker compose down --volumes`
- `docker volume rm bdh_postgres_data`
- `docker volume prune` (if you confirm)

## 📚 Documentation

For more details, see:
- **[Data Persistence Guide](docs/DOCKER_DATA_PERSISTENCE.md)** - Complete guide on Docker volumes and data management
- **[Database Setup](docs/DATABASE_SETUP.md)** - Database configuration and migrations
- **[Test Data README](scripts/TEST_DATA_README.md)** - Sample data documentation

## 💡 Pro Tips

1. **Always create backups before major changes:**
   ```powershell
   .\scripts\backup-db.ps1
   ```

2. **Use the convenience script for daily development:**
   ```powershell
   .\scripts\start-with-data.ps1
   ```

3. **Check logs when something isn't working:**
   ```powershell
   docker compose logs -f app
   ```

4. **Never use `-v` flag with `docker compose down` unless you want to delete everything:**
   ```powershell
   # ✅ Good
   docker compose down
   
   # ❌ Bad (deletes all data)
   docker compose down -v
   ```

5. **Keep your backups organized:**
   - Backups are stored in `db-backups/` directory
   - They're timestamped automatically
   - Restore script shows you all available backups

## 🎯 Development Workflow

### Daily Development
```powershell
# Start (data persists from yesterday)
docker compose up -d

# ... do your work ...

# Stop (data persists for tomorrow)
docker compose down
```

### After Code Changes
```powershell
# Rebuild and restart (data persists)
docker compose down
docker compose up -d --build
```

### Before Major Database Changes
```powershell
# Backup first!
.\scripts\backup-db.ps1

# Make your changes...

# If something goes wrong, restore
.\scripts\restore-db.ps1
```

### Weekly Cleanup
```powershell
# Create a backup
.\scripts\backup-db.ps1

# Optional: Remove old containers
docker system prune

# Never remove volumes unless intentional!
```

## ❓ Need Help?

1. Check the logs: `docker compose logs -f`
2. Verify containers are running: `docker ps`
3. Check database connectivity: `docker exec laptop-tracking-db pg_isready -U postgres`
4. See full documentation in `docs/` directory

