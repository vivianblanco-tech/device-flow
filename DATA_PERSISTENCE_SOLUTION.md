# Data Persistence Solution Summary

## Problem Solved ✅

**Issue**: Sample data disappears when running `docker compose down` and then `docker-compose up -d --build`.

**Root Cause**: Using `docker compose down -v` (or `--volumes`) flag accidentally, which removes all Docker volumes including the database data.

## Solution Implemented

### 1. Documentation Created

- **[DOCKER_DATA_PERSISTENCE.md](docs/DOCKER_DATA_PERSISTENCE.md)** - Complete guide on Docker volumes and data management
- **[QUICK_START.md](QUICK_START.md)** - Quick reference guide for daily development

### 2. Helper Scripts Created

#### `scripts/init-db-if-empty.ps1` ✅
Automatically loads sample data if database is empty.

**Usage:**
```powershell
.\scripts\init-db-if-empty.ps1
```

**Features:**
- Checks if database is ready
- Loads test users if database is empty
- Optionally loads additional test data (companies, laptops, shipments)
- Safe to run multiple times (won't duplicate data)

#### `scripts/backup-db.ps1` ✅
Creates timestamped backups of your database.

**Usage:**
```powershell
.\scripts\backup-db.ps1
```

**Output:**
- Creates backup in `db-backups/` directory
- Filename: `laptop_tracking_backup_YYYYMMDD-HHMMSS.sql`
- Shows backup size and location

#### `scripts/restore-db.ps1` ✅
Restores database from backup files.

**Usage:**
```powershell
# Interactive mode (shows list of backups)
.\scripts\restore-db.ps1

# Restore specific backup
.\scripts\restore-db.ps1 -BackupFile "db-backups/laptop_tracking_backup_20251108-153045.sql"

# Non-interactive (no confirmation prompt)
.\scripts\restore-db.ps1 -Force
```

**Features:**
- Lists all available backups
- Shows backup size and date
- Safety confirmation prompt
- Verifies restored data

#### `scripts/start-with-data.ps1` ✅
One-command startup with automatic data loading.

**Usage:**
```powershell
# Normal start
.\scripts\start-with-data.ps1

# With rebuild
.\scripts\start-with-data.ps1 -Build

# Fresh start (removes all data and volumes)
.\scripts\start-with-data.ps1 -Fresh
```

**Features:**
- Starts all Docker containers
- Automatically initializes database if empty
- Shows service URLs and login credentials
- Provides useful command reference

## How Docker Data Persistence Works

### Your Current Setup (Correct! ✅)

```yaml
volumes:
  postgres_data:/var/lib/postgresql/data  # Named volume (persistent)
  uploads_data:/app/uploads               # Named volume (persistent)
```

- **Volume Name**: `bdh_postgres_data`
- **Type**: Named volume (managed by Docker)
- **Persistence**: Survives container restarts and rebuilds

### Safe Commands (Data Preserved)

```powershell
# ✅ Stop and remove containers (KEEPS DATA)
docker compose down

# ✅ Rebuild and restart (DATA PERSISTS)
docker compose up -d --build

# ✅ Restart services
docker compose restart

# ✅ Stop without removing
docker compose stop
```

### Dangerous Commands (Data Lost)

```powershell
# ❌ DANGER: Removes volumes (ALL DATA LOST)
docker compose down -v
docker compose down --volumes

# ❌ DANGER: Remove specific volume
docker volume rm bdh_postgres_data
```

## Test User Credentials

All test users have been created with the password: **`Test123!`**

| Email | Role | Access Level |
|-------|------|--------------|
| `logistics@bairesdev.com` | Logistics | Full Access |
| `warehouse@bairesdev.com` | Warehouse | Medium Access |
| `client@bairesdev.com` | Client | Limited Access |
| `pm@bairesdev.com` | Project Manager | Read-Only |

**Login URL**: http://localhost:8080/login

## Recommended Daily Workflow

### Starting Your Day

```powershell
# Option 1: Use convenience script (recommended)
.\scripts\start-with-data.ps1

# Option 2: Manual start
docker compose up -d
```

### Making Code Changes

```powershell
# Stop and rebuild (data persists)
docker compose down
docker compose up -d --build
```

### Before Major Database Changes

```powershell
# Create backup first!
.\scripts\backup-db.ps1

# Make your changes...

# If something goes wrong, restore
.\scripts\restore-db.ps1
```

### Ending Your Day

```powershell
# Stop containers (data persists for tomorrow)
docker compose down
```

## Troubleshooting Guide

### Data Disappeared After Restart?

**Likely cause**: Used `docker compose down -v` accidentally.

**Solution**:
```powershell
# If you have a backup
.\scripts\restore-db.ps1

# If no backup, reload test data
.\scripts\init-db-if-empty.ps1
```

### How to Verify Data Persistence?

```powershell
# Check if volumes exist
docker volume ls --filter name=bdh

# Check user count in database
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"
```

### Complete Fresh Start

```powershell
# Interactive with confirmation
.\scripts\start-with-data.ps1 -Fresh

# Or manual
docker compose down -v
docker compose up -d --build
.\scripts\init-db-if-empty.ps1
```

## Backup Strategy Recommendations

### 1. Before Major Changes
```powershell
.\scripts\backup-db.ps1
```

### 2. Weekly Backups
Create a scheduled task (optional):
```powershell
# Create a simple backup script
$scriptPath = "E:\Cursor Projects\BDH\scripts\backup-db.ps1"
& $scriptPath
```

### 3. Keep Backups Organized
- Backups are stored in `db-backups/` directory
- Filename format: `laptop_tracking_backup_YYYYMMDD-HHMMSS.sql`
- Old backups can be deleted manually

## Key Takeaways

1. **NEVER use `docker compose down -v`** unless you intentionally want to delete all data
2. **Always use `docker compose down`** (without flags) for normal shutdown
3. **Create backups before major changes** using `.\scripts\backup-db.ps1`
4. **Use the convenience script** `.\scripts\start-with-data.ps1` for hassle-free startup
5. **Data persists in Docker volumes** - they survive restarts and rebuilds

## Quick Command Reference

```powershell
# Start application (one command)
.\scripts\start-with-data.ps1

# Stop application (safe - keeps data)
docker compose down

# Create backup
.\scripts\backup-db.ps1

# Restore from backup
.\scripts\restore-db.ps1

# Load sample data if empty
.\scripts\init-db-if-empty.ps1

# View logs
docker compose logs -f app

# Check database status
docker exec laptop-tracking-db pg_isready -U postgres
```

## Files Created

1. `docs/DOCKER_DATA_PERSISTENCE.md` - Complete documentation
2. `QUICK_START.md` - Quick reference guide
3. `scripts/init-db-if-empty.ps1` - Automatic data initialization
4. `scripts/backup-db.ps1` - Database backup utility
5. `scripts/restore-db.ps1` - Database restore utility
6. `scripts/start-with-data.ps1` - One-command startup
7. `DATA_PERSISTENCE_SOLUTION.md` - This summary

## Support

For more details:
- **Complete Guide**: [docs/DOCKER_DATA_PERSISTENCE.md](docs/DOCKER_DATA_PERSISTENCE.md)
- **Quick Reference**: [QUICK_START.md](QUICK_START.md)
- **Test Data Info**: [scripts/TEST_DATA_README.md](scripts/TEST_DATA_README.md)

---

**Summary**: Your data will now persist correctly as long as you avoid using the `-v` flag with `docker compose down`. Use the provided scripts for easy management!

