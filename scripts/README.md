# Scripts Directory

This directory contains utility scripts for managing the Laptop Tracking System database and deployment.

## 📁 Contents Overview

### Sample Data Scripts

#### `enhanced-sample-data.sql` ⭐ **RECOMMENDED**
**Primary sample data file for development and testing.**

**Contents:**
- 14 users across all roles (Logistics, Warehouse, PM, Client)
- 8 client companies with complete contact information
- 22 software engineers with addresses
- 35+ laptops (Dell, HP, Lenovo, Apple, Microsoft, ASUS, Acer)
- 15 comprehensive shipments covering:
  - All 3 shipment types (single_full_journey, bulk_to_warehouse, warehouse_to_engineer)
  - All 8 shipment statuses
  - Historical data spanning 6 months
- Complete pickup forms, reception reports, and delivery forms
- Audit logs showing realistic system activity
- Magic links for secure access testing

**When to use:**
- ✅ Standard development and testing
- ✅ Demo environments
- ✅ Initial production testing
- ✅ Feature development
- ✅ Integration testing

**Load command:**
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data.sql
```

#### `enhanced-sample-data-comprehensive.sql`
**Extended base data for high-volume testing.**

**Contents:**
- 30 users across all roles
- 15 client companies
- 50+ software engineers
- 110+ laptops with diverse brands and models
- Proper SKUs and detailed specifications
- Realistic status distributions

**When to use:**
- ✅ Stress testing
- ✅ Performance testing with large datasets
- ✅ Multi-company scenarios
- ✅ User/role permission testing
- ✅ Inventory management testing

**Note:** This file provides base entities (users, companies, engineers, laptops) without shipments. Create shipments through the application or combine with enhanced-sample-data.sql.

**Load command:**
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data-comprehensive.sql
```

### PowerShell Scripts

#### `start-with-data.ps1` 🚀
**Main startup script - automatically initializes database with sample data.**

**Features:**
- Starts Docker containers
- Checks if database is empty
- Automatically loads sample data if needed
- Shows comprehensive status information

**Usage:**
```powershell
# Standard start
.\scripts\start-with-data.ps1

# Fresh start (removes all data and volumes)
.\scripts\start-with-data.ps1 -Fresh

# Rebuild containers
.\scripts\start-with-data.ps1 -Build
```

#### `verify-test-data.ps1` ✅
**Verifies and displays sample data statistics.**

**Shows:**
- Entity counts (users, companies, engineers, laptops, shipments)
- Shipment status breakdown
- Laptop status distribution
- Laptop brands distribution
- Bulk shipments details
- Recent shipments with details
- Data quality indicators

**Usage:**
```powershell
.\scripts\verify-test-data.ps1
```

#### `init-db-if-empty.ps1`
**Automated database initialization (called by start-with-data.ps1).**

Checks if database has data and loads enhanced-sample-data.sql if empty.

### Database Management Scripts

#### `backup-db.ps1`
**Creates a backup of the PostgreSQL database.**

**Usage:**
```powershell
.\scripts\backup-db.ps1
```

#### `restore-db.ps1`
**Restores database from a backup file.**

**Usage:**
```powershell
.\scripts\restore-db.ps1 -BackupFile "path/to/backup.sql"
```

## 🎯 Quick Start Workflows

### First Time Setup
```powershell
# 1. Start application with automatic data loading
.\scripts\start-with-data.ps1

# 2. Verify data loaded correctly
.\scripts\verify-test-data.ps1

# 3. Access application
# Open browser: http://localhost:8080
# Login: logistics@bairesdev.com / Test123!
```

### Reset Database with Fresh Data
```powershell
# Complete reset
.\scripts\start-with-data.ps1 -Fresh

# Or manual reload
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data.sql
```

### Load High-Volume Data
```powershell
# 1. Load comprehensive base data
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data-comprehensive.sql

# 2. Create shipments via the application
# Or load additional sample shipments
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data.sql
```

## 📊 Sample Data Metrics

| Metric | Enhanced (Standard) | Comprehensive |
|--------|---------------------|---------------|
| **Users** | 14 | 30 |
| **Client Companies** | 8 | 15 |
| **Software Engineers** | 22 | 50+ |
| **Laptops** | 35 | 110+ |
| **Shipments** | 15 | 0* |
| **Pickup Forms** | 15 | 0* |
| **Reception Reports** | 7 | 0* |
| **Delivery Forms** | 4 | 0* |
| **Audit Logs** | 10+ | 0* |

*Comprehensive focuses on base entities; add transactional data separately

## 🔐 Test Credentials

All users have password: **Test123!**

**Role-Based Accounts:**
- **Logistics**: logistics@bairesdev.com
- **Warehouse**: warehouse@bairesdev.com
- **Project Manager**: pm@bairesdev.com
- **Client**: client@techcorp.com, admin@innovate.io

## 📝 Sample Data Features

### Shipment Types Coverage
- ✅ **Single Full Journey**: 1 laptop, complete lifecycle (client → warehouse → engineer)
- ✅ **Bulk to Warehouse**: 2-6 laptops, stops at warehouse
- ✅ **Warehouse to Engineer**: 1 laptop, warehouse → engineer only

### Status Coverage
All 8 statuses represented:
1. pending_pickup_from_client
2. pickup_from_client_scheduled
3. picked_up_from_client
4. in_transit_to_warehouse
5. at_warehouse
6. released_from_warehouse
7. in_transit_to_engineer
8. delivered

### Realistic Data Elements
- ✅ Historical timestamps (past 6 months)
- ✅ Detailed equipment specifications
- ✅ Realistic company and engineer information
- ✅ Complete address information
- ✅ Tracking numbers and courier details
- ✅ Detailed notes and observations
- ✅ Photo URL references
- ✅ Accessories descriptions
- ✅ JSON form data
- ✅ Audit trail entries

## 🛠️ Customization

### Creating Additional Sample Data

1. **Through the Application** (Recommended):
   ```
   - Login as logistics user
   - Use web interface to create shipments
   - Ensures all business logic is followed
   - Automatically generates proper relationships
   ```

2. **SQL Script**:
   ```sql
   -- Template for adding a new shipment
   INSERT INTO shipments (
       client_company_id, 
       software_engineer_id, 
       status, 
       shipment_type,
       laptop_count,
       jira_ticket_number, 
       notes, 
       created_at, 
       updated_at
   ) VALUES (
       1,                              -- company ID
       1,                              -- engineer ID
       'pending_pickup_from_client',   -- status
       'single_full_journey',          -- type
       1,                              -- laptop count
       'SCOP-90001',                   -- JIRA ticket
       'Description here',             -- notes
       NOW(),                          -- created
       NOW()                           -- updated
   );
   ```

### Modifying Existing Data

Edit the SQL files directly, then reload:
```powershell
docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev < scripts/enhanced-sample-data.sql
```

## 🐛 Troubleshooting

### Database Connection Issues
```powershell
# Check if container is running
docker ps | findstr laptop-tracking-db

# Restart database container
docker compose restart postgres

# View database logs
docker compose logs postgres
```

### Data Loading Errors
```powershell
# Check for foreign key violations
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "\d shipments"

# Verify table structure
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "\dt"

# Clear and reload
.\scripts\start-with-data.ps1 -Fresh
```

### Script Execution Errors
```powershell
# Enable script execution (if needed)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Run with verbose output
.\scripts\start-with-data.ps1 -Verbose
```

## 📚 Additional Documentation

- **Database Setup**: See `/docs/DATABASE_SETUP.md`
- **Docker Guide**: See `/DOCKER_CHEAT_SHEET.md`
- **Testing Guide**: See `/docs/TESTING_BEST_PRACTICES.md`
- **Sample Data Details**: See `/docs/SAMPLE_DATA_ENHANCEMENT_SUMMARY.md`

## 🔄 Maintenance

### Regular Tasks
- Backup database before major changes: `.\scripts\backup-db.ps1`
- Verify data integrity: `.\scripts\verify-test-data.ps1`
- Update sample data as application evolves
- Test data loading in CI/CD pipelines

### Best Practices
- ✅ Always backup before modifications
- ✅ Test data changes in development first
- ✅ Keep sample data synchronized with schema migrations
- ✅ Document any custom data scenarios
- ✅ Use realistic but anonymized data

---

**Last Updated**: November 13, 2025  
**Maintained By**: Development Team  
**Questions?**: Check `/docs/` or create an issue

