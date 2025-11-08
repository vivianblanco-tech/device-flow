# Docker Data Persistence Guide

## Problem: Data Disappears After `docker compose down`

This guide explains why data might disappear and how to prevent it.

## Understanding Docker Compose Commands

### Safe Commands (Data Persists):
```powershell
# Stop containers, remove containers and networks, KEEP volumes
docker compose down

# Restart services (data persists)
docker compose restart

# Stop without removing
docker compose stop
```

### Destructive Commands (Data Lost):
```powershell
# ❌ DANGER: Removes volumes (ALL DATA LOST)
docker compose down -v
docker compose down --volumes

# ❌ DANGER: Remove specific volume manually
docker volume rm bdh_postgres_data
```

## Current Setup (docker-compose.yml)

Your project uses **named volumes** which provide persistent storage:

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data  # Named volume (persistent)
  - uploads_data:/app/uploads               # Named volume (persistent)
```

### Volume Locations:
- **Volume Name**: `bdh_postgres_data` (for PostgreSQL data)
- **Volume Name**: `bdh_uploads_data` (for uploaded files)
- **Inspect**: `docker volume inspect bdh_postgres_data`

## Solutions

### Solution 1: Use Correct Docker Commands (Recommended)

**ALWAYS use `docker compose down` without the `-v` flag:**

```powershell
# Stop and remove containers (KEEPS DATA)
docker compose down

# Rebuild and restart (DATA PERSISTS)
docker compose up -d --build
```

### Solution 2: Create Initialization Script

Create a script that automatically loads sample data if the database is empty:

**File: `scripts/init-db-if-empty.ps1`**

```powershell
# Check if users table has any data
$userCount = docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -t -c "SELECT COUNT(*) FROM users;" 2>$null

if ($LASTEXITCODE -eq 0 -and $userCount -match '^\s*0\s*$') {
    Write-Host "Database is empty. Loading sample data..." -ForegroundColor Yellow
    
    # Load test users
    Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
    
    # Load test data (optional)
    # Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
    
    Write-Host "Sample data loaded successfully!" -ForegroundColor Green
} else {
    Write-Host "Database already contains data. Skipping initialization." -ForegroundColor Cyan
}
```

**Usage:**
```powershell
docker compose up -d --build
./scripts/init-db-if-empty.ps1
```

### Solution 3: Use Local Directory Volume (Alternative)

Change to a local directory volume for easier access and guaranteed persistence:

**Modify `docker-compose.yml`:**

```yaml
services:
  postgres:
    # ... existing config ...
    volumes:
      - ./db-data/postgres:/var/lib/postgresql/data  # Local directory
```

**Pros:** 
- Data stored in your project folder
- Easy to backup
- Survives `docker compose down -v`

**Cons:**
- Slower on Windows/Mac
- Permission issues on some systems
- Already set up (you have `./db-data/postgres/`)

### Solution 4: Backup and Restore Scripts

Create scripts to backup and restore your database:

**File: `scripts/backup-db.ps1`**

```powershell
# Create backup directory
$backupDir = "db-backups"
if (!(Test-Path $backupDir)) {
    New-Item -ItemType Directory -Path $backupDir
}

# Generate backup filename with timestamp
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$backupFile = "$backupDir/laptop_tracking_backup_$timestamp.sql"

Write-Host "Creating database backup..." -ForegroundColor Cyan
docker exec laptop-tracking-db pg_dump -U postgres laptop_tracking_dev > $backupFile

if ($LASTEXITCODE -eq 0) {
    Write-Host "Backup created: $backupFile" -ForegroundColor Green
    
    # Show file size
    $fileSize = (Get-Item $backupFile).Length / 1KB
    Write-Host "Size: $([math]::Round($fileSize, 2)) KB" -ForegroundColor Cyan
} else {
    Write-Host "Backup failed!" -ForegroundColor Red
}
```

**File: `scripts/restore-db.ps1`**

```powershell
param(
    [Parameter(Mandatory=$false)]
    [string]$BackupFile
)

# If no backup file specified, use the latest
if ([string]::IsNullOrEmpty($BackupFile)) {
    $BackupFile = Get-ChildItem -Path "db-backups" -Filter "*.sql" | 
                  Sort-Object LastWriteTime -Descending | 
                  Select-Object -First 1 -ExpandProperty FullName
}

if (!(Test-Path $BackupFile)) {
    Write-Host "Backup file not found: $BackupFile" -ForegroundColor Red
    exit 1
}

Write-Host "Restoring database from: $BackupFile" -ForegroundColor Cyan

# Drop and recreate database
docker exec laptop-tracking-db psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_dev;"
docker exec laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;"

# Restore from backup
Get-Content $BackupFile | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database restored successfully!" -ForegroundColor Green
} else {
    Write-Host "Restore failed!" -ForegroundColor Red
}
```

**Usage:**
```powershell
# Create backup before making changes
./scripts/backup-db.ps1

# Restore latest backup
./scripts/restore-db.ps1

# Restore specific backup
./scripts/restore-db.ps1 -BackupFile "db-backups/laptop_tracking_backup_20251108-153045.sql"
```

### Solution 5: Docker Init Script (Automatic)

Add an initialization script that runs when the database is first created:

**File: `scripts/init-db.sh`** (for Docker)

```bash
#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create test users
    INSERT INTO users (email, password_hash, role, created_at, updated_at)
    VALUES 
        ('logistics@bairesdev.com', '\$2a\$12\$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'logistics', NOW(), NOW()),
        ('client@bairesdev.com', '\$2a\$12\$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'client', NOW(), NOW()),
        ('warehouse@bairesdev.com', '\$2a\$12\$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'warehouse', NOW(), NOW()),
        ('pm@bairesdev.com', '\$2a\$12\$5jhaEE3wXZtjKA/a07CHvunJymFovVivi8e1X7WX/zQCS9wmLU2yK', 'project_manager', NOW(), NOW())
    ON CONFLICT (email) DO NOTHING;
EOSQL
```

**Modify `docker-compose.yml`:**

```yaml
services:
  postgres:
    # ... existing config ...
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sh:/docker-entrypoint-initdb.d/init-db.sh  # Add this
```

**Note:** This only runs when the database is **first created** (empty volume).

## Recommended Workflow

### Daily Development:
```powershell
# Start services (data persists)
docker compose up -d

# Stop services (data persists)
docker compose down
```

### When Making Code Changes:
```powershell
# Rebuild application (data persists)
docker compose down
docker compose up -d --build
```

### When Database is Empty:
```powershell
# Start services
docker compose up -d

# Load sample data
Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
Get-Content scripts/create-test-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

### Complete Reset (Fresh Start):
```powershell
# ⚠️ WARNING: This deletes ALL data
docker compose down -v
docker volume rm bdh_postgres_data bdh_uploads_data
docker compose up -d --build

# Wait for database to be ready
Start-Sleep -Seconds 5

# Load sample data
Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

## Verify Data Persistence

Check if volumes exist:
```powershell
docker volume ls --filter name=bdh
```

Check volume contents:
```powershell
docker volume inspect bdh_postgres_data
```

Verify data in database:
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"
```

## Troubleshooting

### Data Still Disappears?

1. **Check if migrations are dropping data:**
   - Review migration files in `migrations/`
   - Ensure migrations don't have `DROP TABLE` or `TRUNCATE` commands

2. **Check if volume is actually being used:**
   ```powershell
   docker inspect laptop-tracking-db | Select-String -Pattern "Mounts" -Context 0,20
   ```

3. **Check volume driver:**
   ```powershell
   docker volume inspect bdh_postgres_data
   ```

4. **Permissions issues:**
   - On Windows, Docker Desktop must have access to the drive
   - Settings → Resources → File Sharing

### Manual Volume Management

List all volumes:
```powershell
docker volume ls
```

Inspect specific volume:
```powershell
docker volume inspect bdh_postgres_data
```

Remove unused volumes (careful!):
```powershell
docker volume prune
```

## Best Practices

✅ **DO:**
- Use `docker compose down` (no flags) for normal shutdown
- Use named volumes for persistent data (already configured)
- Create regular backups before major changes
- Use init scripts for automatic data loading
- Version control your migration files

❌ **DON'T:**
- Use `docker compose down -v` unless you want to delete data
- Modify files in the volume directory directly
- Delete volumes manually unless intentional
- Mix named volumes with bind mounts for the same data

## Quick Reference

| Command | Effect on Data | Use Case |
|---------|---------------|----------|
| `docker compose up -d` | ✅ Preserves | Start services |
| `docker compose down` | ✅ Preserves | Stop and remove containers |
| `docker compose restart` | ✅ Preserves | Quick restart |
| `docker compose stop` | ✅ Preserves | Pause containers |
| `docker compose down -v` | ❌ DELETES | Complete cleanup |
| `docker volume rm <name>` | ❌ DELETES | Manual volume deletion |

## Summary

Your current setup is correct with named volumes. The most likely cause of data loss is:

1. **Accidentally using `-v` flag** with `docker compose down`
2. **Migrations resetting tables** during startup
3. **Manual volume deletion**

**Fix:** Always use `docker compose down` without flags, and create backup/restore scripts for safety.

