# Docker Data Persistence - Cheat Sheet

## 🚨 THE GOLDEN RULE 🚨

```
✅ docker compose down     (Safe - keeps your data)
❌ docker compose down -v  (Dangerous - deletes everything!)
```

---

## 🚀 One-Command Solutions

| What You Want | Command |
|---------------|---------|
| Start everything | `.\scripts\start-with-data.ps1` |
| Stop safely | `docker compose down` |
| Rebuild app | `docker compose down && docker compose up -d --build` |
| Create backup | `.\scripts\backup-db.ps1` |
| Restore backup | `.\scripts\restore-db.ps1` |
| Load sample data | `.\scripts\init-db-if-empty.ps1` |

---

## ✅ Safe Commands (Data Persists)

```powershell
docker compose up -d                    # Start services
docker compose down                     # Stop services (DATA SAFE)
docker compose up -d --build           # Rebuild (DATA SAFE)
docker compose restart                  # Quick restart (DATA SAFE)
docker compose stop                     # Pause (DATA SAFE)
.\scripts\start-with-data.ps1          # Convenience script (DATA SAFE)
```

---

## ❌ Dangerous Commands (Data Lost)

```powershell
docker compose down -v                  # ⚠️ DELETES VOLUMES!
docker compose down --volumes           # ⚠️ DELETES VOLUMES!
docker volume rm bdh_postgres_data     # ⚠️ DELETES DATABASE!
docker volume prune                     # ⚠️ DELETES UNUSED VOLUMES!
.\scripts\start-with-data.ps1 -Fresh   # ⚠️ FRESH START (after confirmation)
```

---

## 🔐 Test Login Credentials

**Password for ALL users: `Test123!`**

```
Email: logistics@bairesdev.com  →  Role: Logistics (Full Access)
Email: warehouse@bairesdev.com  →  Role: Warehouse (Medium Access)
Email: client@bairesdev.com     →  Role: Client (Limited Access)
Email: pm@bairesdev.com         →  Role: Project Manager (Read-Only)
```

**Login**: http://localhost:8080/login

---

## 📊 Quick Diagnostics

```powershell
# Check if containers are running
docker ps

# Check if volumes exist
docker volume ls --filter name=bdh

# Check database has data
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"

# View application logs
docker compose logs -f app

# Check database health
docker exec laptop-tracking-db pg_isready -U postgres
```

---

## 🔄 Common Workflows

### Daily Development
```powershell
# Morning: Start
.\scripts\start-with-data.ps1

# ... do your work ...

# Evening: Stop
docker compose down
```

### After Code Changes
```powershell
docker compose down
docker compose up -d --build
```

### Before Major Changes
```powershell
# Backup first!
.\scripts\backup-db.ps1

# Make changes...

# Something went wrong? Restore!
.\scripts\restore-db.ps1
```

### Complete Fresh Start
```powershell
# Option 1: Safe with confirmation
.\scripts\start-with-data.ps1 -Fresh

# Option 2: Manual
docker compose down -v
docker compose up -d --build
.\scripts\init-db-if-empty.ps1
```

---

## 🆘 Emergency Recovery

### Data Disappeared!

**Step 1**: Check if you have backups
```powershell
Get-ChildItem db-backups
```

**Step 2**: Restore latest backup
```powershell
.\scripts\restore-db.ps1
```

**Step 3**: No backup? Reload sample data
```powershell
.\scripts\init-db-if-empty.ps1
```

### Can't Login?

**Step 1**: Check if users exist
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT email, role FROM users;"
```

**Step 2**: Recreate users
```powershell
Get-Content scripts/create-test-users-all-roles.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

### Database Connection Failed?

```powershell
# Check if PostgreSQL is running
docker compose logs postgres

# Restart database
docker compose restart postgres

# Full restart
docker compose down
docker compose up -d
```

---

## 📁 Important Directories

```
db-backups/          → Database backup files (*.sql)
db-data/             → Local database data (if using local volume)
scripts/             → Helper scripts
docs/                → Documentation
uploads/             → Uploaded files (delivery, reception)
migrations/          → Database migration files
```

---

## 🎯 Remember

1. **`docker compose down`** is SAFE ✅
2. **`docker compose down -v`** DELETES DATA ❌
3. **Backup before major changes** 💾
4. **Use the helper scripts** 🚀
5. **Data lives in Docker volumes** 📦

---

## 📚 Full Documentation

- **Complete Guide**: [docs/DOCKER_DATA_PERSISTENCE.md](docs/DOCKER_DATA_PERSISTENCE.md)
- **Quick Start**: [QUICK_START.md](QUICK_START.md)
- **Solution Summary**: [DATA_PERSISTENCE_SOLUTION.md](DATA_PERSISTENCE_SOLUTION.md)

---

**Created**: 2025-11-08  
**Purpose**: Prevent accidental data loss in Docker development environment

