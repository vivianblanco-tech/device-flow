# Enhanced Sample Data Documentation

## Overview

The enhanced sample data provides a comprehensive, realistic dataset for testing and demonstrating the Laptop Tracking System. This data includes bulk shipments, high-end workstations, premium laptops, and complete documentation across all shipment stages.

## Loading the Data

### Automatic Loading (Recommended)
```powershell
.\scripts\start-with-data.ps1
```
The script will automatically detect if the database is empty and prompt you to load the comprehensive sample data.

### Manual Loading
```powershell
# Load directly into Docker database
Get-Content scripts/enhanced-sample-data.sql | docker exec -i laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

### Verification
```powershell
.\scripts\verify-test-data.ps1
```
This will display detailed statistics about the loaded data.

## Data Contents

### Users (13 total)
All users have password: `Test123!`

**Logistics Team:**
- logistics@bairesdev.com
- sarah.logistics@bairesdev.com
- james.logistics@bairesdev.com

**Warehouse Team:**
- warehouse@bairesdev.com
- michael.warehouse@bairesdev.com
- jessica.warehouse@bairesdev.com

**Project Managers:**
- pm@bairesdev.com
- jennifer.pm@bairesdev.com
- david.pm@bairesdev.com

**Client Users:**
- client@techcorp.com (TechCorp International)
- admin@innovate.io (Innovate Solutions Ltd)
- purchasing@globaltech.com (Global Tech Services)
- it-manager@digitaldynamics.com (Digital Dynamics Corp)
- operations@cloudventures.com (Cloud Ventures Inc)

### Client Companies (8 total)
1. TechCorp International - San Francisco, CA
2. Innovate Solutions Ltd - Austin, TX
3. Global Tech Services - Seattle, WA
4. Digital Dynamics Corp - Boston, MA
5. Cloud Ventures Inc - Denver, CO
6. DataDrive Systems - Chicago, IL
7. NextGen Software - Portland, OR
8. Enterprise Solutions Group - New York, NY

### Software Engineers (22 total)
Engineers distributed across all client companies with realistic contact information and addresses.

### Laptops (35+ units)

**High-End Workstations:**
- **Dell Precision 5570** - Intel i9-12900H, 64GB RAM, RTX A2000, 2TB SSD
- **Dell Precision 7670** - Intel i9-12950HX, 128GB RAM, RTX A5500, 4TB SSD
- **HP ZBook Studio G9** - Intel i9-12900H, 64GB RAM, RTX A3000, 2TB SSD
- **HP ZBook Fury G9** - Intel i9-12950HX, 128GB RAM, RTX A5500, 4TB SSD
- **Lenovo ThinkPad P1 Gen 5** - Intel i9-12900H, 64GB RAM, RTX A5500, 2TB SSD
- **Lenovo ThinkPad P16 Gen 1** - Intel i9-12950HX, 128GB RAM, RTX A5500, 4TB SSD

**Premium Developer Laptops:**
- **Dell XPS 15 9520** - Intel i7-12700H, 32GB RAM, RTX 3050 Ti, 4K OLED
- **Dell XPS 13 Plus 9315** - Intel i7-1360P, 32GB RAM, FHD+ Touch
- **HP EliteBook 850/840 G9** - Intel i7-1265U, 16-32GB RAM, FHD/LTE
- **Lenovo ThinkPad X1 Carbon Gen 10** - Intel i7-1260P, 32GB RAM, WQUXGA, 5G

**Apple MacBooks:**
- **MacBook Pro 16" M2 Max** - 96GB RAM, 2TB SSD, Liquid Retina XDR
- **MacBook Pro 16" M2 Pro** - 32GB RAM, 1TB SSD
- **MacBook Pro 14" M2 Pro** - 16GB RAM, 512GB SSD
- **MacBook Air 13" M2** - 24GB RAM, 1TB SSD

**Other Premium Options:**
- **Microsoft Surface Laptop Studio** - i7-11370H, 32GB RAM, RTX 3050 Ti
- **Microsoft Surface Laptop 5** - i7-1255U, 32GB RAM
- **ASUS ZenBook Pro 15 OLED** - i9-12900H, 32GB RAM, RTX 3060
- **ASUS ROG Zephyrus G14** - Ryzen 9 6900HS, 32GB RAM, RX 6800S
- **Acer Swift X** - Ryzen 7 5800U, 16GB RAM, RTX 3050 Ti

All laptops include:
- Detailed specifications with CPU, RAM, Storage, Display, GPU
- SKU numbers for inventory management
- Realistic status assignments
- Complete accessory information

### Shipments (15 total)

#### Status Distribution:
- **Delivered:** 4 shipments (including 2 bulk shipments)
- **In Transit to Engineer:** 3 shipments (including 1 bulk: 5 MacBooks)
- **Released from Warehouse:** 1 shipment
- **At Warehouse:** 2 shipments (including 1 bulk: 4 ThinkPads)
- **In Transit to Warehouse:** 1 shipment
- **Picked Up from Client:** 1 shipment
- **Pickup Scheduled:** 2 shipments (including 1 bulk: 6 EliteBooks)
- **Pending Pickup:** 2 shipments (including 1 urgent bulk: 3 MacBooks)

#### Bulk Shipments:
1. **SCOP-80002** - 3 HP ZBook workstations (Delivered)
2. **SCOP-80004** - 5 MacBook Pro M2 Max (In Transit, high-value)
3. **SCOP-80007** - 4 ThinkPad X1 Carbon (At Warehouse)
4. **SCOP-80011** - 6 HP EliteBook (Pickup Scheduled)
5. **SCOP-80013** - 3 MacBook Pro (Pending Pickup, urgent)
6. **SCOP-80015** - 2 Acer Swift X (Delivered)

### Pickup Forms (13 forms)
Each pickup form includes:
- Contact person details (name, email, phone)
- Complete pickup address with city, state, zip
- Scheduled pickup date and time slot
- Number of laptops and boxes
- Assignment type (single or bulk)
- For bulk shipments: dimensions and weight
- Detailed accessories description
- Special instructions and handling notes

**Example Bulk Form Data:**
```json
{
  "number_of_laptops": 5,
  "number_of_boxes": 3,
  "assignment_type": "bulk",
  "bulk_length": 26.0,
  "bulk_width": 22.0,
  "bulk_height": 16.0,
  "bulk_weight": 62.0,
  "accessories_description": "5x Apple USB-C cables, 5x Magic Mouse, 5x Magic Keyboard, adapters...",
  "special_instructions": "HIGH-VALUE SHIPMENT: Total value >$20,000. Signature required..."
}
```

### Reception Reports (7 reports)
Warehouse reports include:
- Timestamp of reception
- Condition assessment
- Serial number verification
- Functionality testing results
- Accessory inventory
- Photo documentation URLs
- Storage location
- Special handling notes

**Example Reception Report:**
> "BULK RECEPTION: 3x HP ZBook Studio G9 workstations received. Heavy shipment - two-person lift required. All boxes in excellent condition, no shipping damage detected. Serial numbers verified: HP-ZBOOK-G9-001, HP-ZBOOK-G9-002, HP-ZBOOK-G9-003. Individual inspection performed on each unit... Display quality excellent - 4K DreamColor panels tested, no dead pixels on any unit. RTX A3000 GPUs tested with benchmark - all performing within spec..."

### Delivery Forms (4 forms)
Engineer delivery documentation includes:
- Delivery timestamp and location
- Engineer contact confirmation
- Package condition verification
- On-site setup and testing
- Performance verification
- Engineer satisfaction feedback
- Multiple photos of delivery and setup
- Follow-up scheduling

**Example Delivery Report:**
> "BULK DELIVERY: 3x HP ZBook Studio G9 workstations delivered to Emily Rodriguez... Sequential setup of all 3 units... DaVinci Resolve tested on all units - renders fast, no issues. Adobe Creative Cloud installed and verified. Network rendering setup tested between units. Engineer feedback: 'These are exactly what our video editing team needs. The color accuracy on the DreamColor displays is perfect...'"

### Audit Logs
Complete audit trail including:
- Shipment creation events
- Status transitions
- Form submissions
- Engineer assignments
- Tracking updates

## Test Scenarios

### 1. Bulk Shipment Workflow
**Shipment SCOP-80002 (Completed):**
- 3 HP ZBook workstations
- Complete journey from pickup to delivery
- Video editing team deployment
- High-value equipment handling
- Multi-unit testing and setup

### 2. High-Value Equipment
**Shipment SCOP-80004 (In Transit):**
- 5 MacBook Pro M2 Max
- Total value >$20,000
- White-glove service
- iOS development team
- Extra insurance and security

### 3. Urgent Bulk Shipment
**Shipment SCOP-80013 (Pending Pickup):**
- 3 MacBook Pro M2
- Emergency team scaling
- Expedited processing required
- Priority handling

### 4. New Hire Onboarding
**Shipment SCOP-80007 (At Warehouse):**
- 4 ThinkPad X1 Carbon
- Pending engineer assignments
- Scheduled for Monday onboarding
- Complete accessory packages

### 5. Standard Single Delivery
**Shipment SCOP-80001 (Delivered):**
- Dell XPS 15
- Complete standard workflow
- Developer laptop
- Successful delivery and setup

## Realistic Features

### Accessories Included
- Power adapters (original manufacturer)
- Laptop bags/backpacks/cases
- Wireless keyboards and mice
- USB-C docking stations
- Cable adapters and organizers
- External monitors (for workstations)
- Graphics tablets (for creative work)
- Specialized peripherals

### Detailed Notes
- Building access instructions
- Security protocols
- Contact information
- Time windows
- Special handling requirements
- Equipment specifications
- Testing procedures
- Engineer feedback

### Geographic Distribution
Locations across major US cities:
- San Francisco, CA
- Austin, TX
- Seattle, WA
- Boston, MA
- Denver, CO
- Chicago, IL
- Portland, OR
- New York, NY
- Los Angeles, CA
- Miami, FL
- Phoenix, AZ
- Houston, TX

## Usage Tips

### Testing Different Roles

**Logistics User:**
- Create new shipments
- Update shipment status
- Manage tracking numbers
- View all shipments

**Warehouse User:**
- Create reception reports
- Process arrivals
- Release to engineers
- Inventory management

**Project Manager:**
- Assign engineers to shipments
- View all statuses
- Monitor delays
- Resource planning

**Client User:**
- Submit pickup forms
- View company shipments
- Track deliveries
- Request new equipment

### Quick Access Commands

```powershell
# Start everything with data
.\scripts\start-with-data.ps1

# Verify what's loaded
.\scripts\verify-test-data.ps1

# Fresh start (WARNING: Deletes all data)
.\scripts\start-with-data.ps1 -Fresh

# Restart services
docker compose restart

# View logs
docker compose logs -f app

# Access database
docker exec -it laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

### SQL Queries for Testing

```sql
-- View all bulk shipments
SELECT s.id, s.jira_ticket_number, s.status, COUNT(sl.laptop_id) as laptops
FROM shipments s
JOIN shipment_laptops sl ON sl.shipment_id = s.id
GROUP BY s.id
HAVING COUNT(sl.laptop_id) > 1
ORDER BY COUNT(sl.laptop_id) DESC;

-- View high-end workstations
SELECT serial_number, brand, model, specs, status
FROM laptops
WHERE model LIKE '%Precision%' OR model LIKE '%ZBook%' OR model LIKE '%P1%' OR model LIKE '%P16%'
ORDER BY brand, model;

-- View shipments by client
SELECT cc.name, COUNT(s.id) as shipment_count
FROM client_companies cc
LEFT JOIN shipments s ON s.client_company_id = cc.id
GROUP BY cc.name
ORDER BY shipment_count DESC;

-- View recent activity
SELECT timestamp, action, entity_type, details
FROM audit_logs
ORDER BY timestamp DESC
LIMIT 20;
```

## Maintenance

### Backup Current Data
```powershell
.\scripts\backup-db.ps1
```

### Restore from Backup
```powershell
.\scripts\restore-db.ps1
```

### Reset to Enhanced Sample Data
```powershell
# Stop services
docker compose down -v

# Start fresh
.\scripts\start-with-data.ps1 -Fresh
```

## Support

For issues or questions:
1. Check application logs: `docker compose logs -f app`
2. Check database logs: `docker compose logs -f db`
3. Verify data: `.\scripts\verify-test-data.ps1`
4. Review documentation in `/docs` folder

## Summary

This enhanced sample data provides:
- ✅ **15 shipments** covering all workflow stages
- ✅ **35+ laptops** from 7 major brands
- ✅ **6 bulk shipments** with 2-6 laptops each
- ✅ **22 engineers** across 8 companies
- ✅ **Complete documentation** at every stage
- ✅ **Realistic accessories** and specifications
- ✅ **Detailed notes** and instructions
- ✅ **Geographic distribution** across US
- ✅ **Multiple user roles** for testing
- ✅ **Audit trail** of all activities

Perfect for development, testing, demonstrations, and training!

