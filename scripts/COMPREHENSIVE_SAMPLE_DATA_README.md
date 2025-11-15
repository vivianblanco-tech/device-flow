# Comprehensive Sample Data v2.1 Documentation

## Overview

The comprehensive sample data v2.1 provides production-ready test data that covers **all current functionalities** of the Laptop Tracking System, including:

- ✅ **All 3 shipment types** (single full journey, bulk to warehouse, warehouse to engineer)
- ✅ **All 8 shipment statuses** (complete workflow representation)
- ✅ **Laptop-based reception reports** with approval workflow
- ✅ **Address confirmation** tracking for software engineers  
- ✅ **Magic links** for secure delivery confirmation
- ✅ **Complete audit trail** of all activities
- ✅ **Realistic data volume** suitable for demonstrations and testing

## Loading the Data

### Quick Start (Recommended)

```powershell
# Ensure Docker is running
docker-compose up -d

# Load comprehensive sample data
.\scripts\load-sample-data.ps1
```

The script will:
1. Load base data (users, companies, engineers, laptops)
2. Load shipments data (shipments, forms, reports, audit logs)
3. Display a summary of loaded data

### Verification

```powershell
.\scripts\verify-test-data.ps1
```

This shows detailed statistics and confirms all data loaded correctly.

### Fresh Start

```powershell
# WARNING: This deletes all data!
.\scripts\start-with-data.ps1 -Fresh
```

## Data Contents

### 1. Users (32 total)

**All passwords:** `Test123!`

#### Logistics Team (6 users)
- `logistics@bairesdev.com` (primary)
- `sarah.logistics@bairesdev.com`
- `james.logistics@bairesdev.com`
- `maria.logistics@bairesdev.com`
- `robert.logistics@bairesdev.com`
- `jennifer.logistics@bairesdev.com`

#### Warehouse Team (6 users)
- `warehouse@bairesdev.com` (primary)
- `michael.warehouse@bairesdev.com`
- `jessica.warehouse@bairesdev.com`
- `chris.warehouse@bairesdev.com`
- `amanda.warehouse@bairesdev.com`
- `kevin.warehouse@bairesdev.com`

#### Project Managers (5 users)
- `pm@bairesdev.com` (primary)
- `jennifer.pm@bairesdev.com`
- `david.pm@bairesdev.com`
- `sophia.pm@bairesdev.com`
- `william.pm@bairesdev.com`

#### Client Users (15 users)
Each client company has at least one user linked via `client_company_id`:
- `client@techcorp.com` → TechCorp International
- `admin@innovate.io` → Innovate Solutions Ltd
- `purchasing@globaltech.com` → Global Tech Services
- And 12 more across other companies...

### 2. Client Companies (15 total)

1. **TechCorp International** - San Francisco, CA
2. **Innovate Solutions Ltd** - Austin, TX
3. **Global Tech Services** - Seattle, WA
4. **Digital Dynamics Corp** - Boston, MA
5. **Cloud Ventures Inc** - Denver, CO
6. **DataDrive Systems** - Chicago, IL
7. **NextGen Software** - Portland, OR
8. **Enterprise Solutions Group** - New York, NY
9. **Fusion Labs** - Miami, FL
10. **Quantum Code Solutions** - Phoenix, AZ
11. **PixelPerfect Studios** - Los Angeles, CA
12. **RapidTech Industries** - Dallas, TX
13. **SynergySoft Corporation** - Atlanta, GA
14. **Zenith Technologies** - Philadelphia, PA
15. **Apex Digital Group** - San Diego, CA

Each company includes complete contact information (email, phone, address, website).

### 3. Software Engineers (35+ total)

Engineers are distributed across all client companies with:
- Full contact details (name, email, phone, address)
- **Address confirmation tracking** (`address_confirmed` field)
- Realistic geographic distribution across major US cities

**Address Confirmation Status:**
- ~75% have confirmed addresses (ready for delivery)
- ~25% pending confirmation (realistic scenario)

### 4. Laptops (40+ units)

#### Brand Distribution:
- **Dell:** Precision workstations, XPS series
- **HP:** ZBook workstations, EliteBook business laptops
- **Lenovo:** ThinkPad X1 Carbon, P-series workstations
- **Apple:** MacBook Pro (M2 Max/Pro), MacBook Air M2

#### Status Distribution:
- `available` - Available in inventory
- `at_warehouse` - Received at warehouse, pending assignment
- `in_transit_to_warehouse` - On the way to warehouse
- `in_transit_to_engineer` - On the way to engineer
- `delivered` - Successfully delivered to engineer

#### Technical Details:
Each laptop includes:
- **Serial number** (unique identifier)
- **SKU** (stock keeping unit for inventory)
- **Brand and model**
- **RAM** (in GB)
- **SSD** (in GB)
- **Status** (current lifecycle stage)
- **Client company assignment** (when applicable)
- **Engineer assignment** (when delivered/in-transit)

### 5. Shipments (7+ complete workflows)

#### Three Shipment Types:

##### Type 1: Single Full Journey
**Description:** One laptop travels from client → warehouse → engineer (complete lifecycle)

**Example:** Shipment SCOP-90001 (Delivered)
- Dell XPS delivered to Alice Johnson
- Complete pickup form with client details
- Laptop-based reception report (approved)
- Delivery form with engineer confirmation
- Full audit trail from creation to delivery

##### Type 2: Bulk to Warehouse
**Description:** Multiple laptops (2+) from client to warehouse only

**Example:** Shipment SCOP-90002 (At Warehouse)
- 5x Lenovo ThinkPad X1 Carbon
- Bulk pickup form with dimensions and weight
- 5 individual laptop-based reception reports
- All reports pending logistics approval
- Awaiting engineer assignments

##### Type 3: Warehouse to Engineer
**Description:** Single laptop from warehouse inventory directly to engineer

**Example:** Shipment SCOP-90003 (In Transit to Engineer)
- HP ZBook Fury from warehouse stock
- Assigned to Henry Thompson
- Currently in transit (ETA: Tomorrow)
- No client pickup phase (starts at warehouse)

#### All 8 Status Stages Represented:

1. **pending_pickup_from_client** - Awaiting pickup form
2. **pickup_from_client_scheduled** - Form submitted, pickup scheduled
3. **picked_up_from_client** - Picked up by courier (just happened)
4. **in_transit_to_warehouse** - En route to warehouse
5. **at_warehouse** - Received at warehouse
6. **released_from_warehouse** - Assigned and released
7. **in_transit_to_engineer** - En route to engineer
8. **delivered** - Successfully delivered

### 6. Pickup Forms (Complete JSON Data)

Each pickup form includes detailed information:

**Single Shipment Form:**
```json
{
  "contact_name": "Sarah Mitchell",
  "contact_email": "sarah.mitchell@techcorp.com",
  "contact_phone": "+1-555-0101",
  "pickup_address": "100 Tech Plaza, Building 1, Loading Dock A",
  "pickup_city": "San Francisco",
  "pickup_state": "CA",
  "pickup_zip": "94105",
  "pickup_date": "2025-11-10",
  "pickup_time_slot": "morning",
  "special_instructions": "Call 30 minutes before arrival..."
}
```

**Bulk Shipment Form:**
```json
{
  "number_of_laptops": 5,
  "number_of_boxes": 3,
  "assignment_type": "bulk",
  "bulk_length": 24.0,
  "bulk_width": 20.0,
  "bulk_height": 14.0,
  "bulk_weight": 52.5,
  "special_instructions": "BULK SHIPMENT - Forklift required..."
}
```

**Warehouse-to-Engineer Form:**
```json
{
  "laptop_id": 23,
  "engineer_id": 12,
  "contact_name": "Henry Thompson",
  "shipping_address": "808 Pike St, Seattle, WA 98101",
  "special_instructions": "Priority delivery. Signature required."
}
```

### 7. Laptop-Based Reception Reports (NEW!)

**Key Features:**
- ✅ One report per laptop (not per shipment)
- ✅ Three required photos per laptop:
  - Serial number photo
  - External condition photo
  - Working condition photo
- ✅ Approval workflow:
  - Warehouse creates report → `pending_approval` status
  - Logistics approves report → `approved` status
  - Only approved reports allow laptop progression

**Example Report:**
```
Laptop: Dell XPS 13 Plus (DELL-XPS-9315-002)
Status: Approved
Warehouse User: warehouse@bairesdev.com
Received: 2 days ago
Notes: "Serial number verified. All ports functional. Display perfect. Battery 100%."
Photos:
  - /uploads/reception/laptop12_serial.jpg
  - /uploads/reception/laptop12_external.jpg
  - /uploads/reception/laptop12_working.jpg
Approved By: logistics@bairesdev.com
Approved At: 1 day ago
```

**Pending Approval Example:**
- Bulk shipment SCOP-90002 has 5 reception reports
- All 5 are in `pending_approval` status
- Logistics user must review and approve each before assignment

### 8. Delivery Forms

Each delivered shipment includes:
- Engineer confirmation
- Delivery timestamp
- On-site setup notes
- Photo documentation
- Engineer satisfaction feedback

### 9. Audit Logs (Complete Activity Trail)

Tracks all system activities:
- Shipment creation
- Form submissions
- Status transitions
- Engineer assignments
- Report approvals
- Delivery confirmations

### 10. Magic Links (Secure Access)

**Active Link:**
- Token: `abc123def456ghi789jkl012mno345pqr678`
- For shipment: SCOP-90003
- Expires: 7 days from creation
- Used: No (still active)

**Used Link (Historical):**
- Token: `xyz789uvw456rst123opq890lmn567hij234`
- For shipment: SCOP-90001
- Used at: Delivery confirmation
- Status: Expired/Used

## Test Scenarios

### Scenario 1: Complete Single Journey

**Shipment:** SCOP-90001 (Delivered)

**Steps to Verify:**
1. Login as `logistics@bairesdev.com` (Password: `Test123!`)
2. Go to Shipments → View shipment SCOP-90001
3. Observe complete lifecycle:
   - ✓ Pickup form submitted
   - ✓ Picked up from client
   - ✓ In transit to warehouse
   - ✓ Received at warehouse (with reception report)
   - ✓ Reception report approved
   - ✓ Released from warehouse
   - ✓ In transit to engineer
   - ✓ Delivered (with delivery form)
4. Check laptop DELL-XPS-9315-002 → Status should be `delivered`

### Scenario 2: Bulk Shipment with Pending Approvals

**Shipment:** SCOP-90002 (At Warehouse)

**Steps to Verify:**
1. Login as `warehouse@bairesdev.com`
2. View shipment SCOP-90002
3. See 5 ThinkPad laptops at warehouse
4. Check reception reports → All 5 are `pending_approval`
5. Login as `logistics@bairesdev.com`
6. Go to Reception Reports → See 5 pending reports
7. Approve reports one by one
8. After approval, laptops can be assigned to engineers

### Scenario 3: Warehouse to Engineer Direct

**Shipment:** SCOP-90003 (In Transit to Engineer)

**Steps to Verify:**
1. Login as `logistics@bairesdev.com`
2. View shipment SCOP-90003
3. Notice:
   - Shipment type: `warehouse_to_engineer`
   - No client pickup phase
   - Laptop came from warehouse inventory
   - Directly assigned to Henry Thompson
   - Currently in transit (ETA shown)
4. Check laptop HP-ZBOOK-FUR-G9-001 → Status `in_transit_to_engineer`

### Scenario 4: In-Progress Workflows

**Multiple Active Shipments:**
- **SCOP-90004** - In transit to warehouse (picked up yesterday)
- **SCOP-90005** - Pickup scheduled for tomorrow
- **SCOP-90006** - Pending pickup (just created, no form yet)

**Test Workflow Progression:**
1. Login as `client@techcorp.com`
2. Submit pickup form for SCOP-90006
3. Status should change to `pickup_from_client_scheduled`
4. Login as `logistics@bairesdev.com`
5. Update status to simulate pickup
6. Continue workflow through all stages

### Scenario 5: Reception Report Approval Workflow

**Test the NEW laptop-based system:**
1. Login as `warehouse@bairesdev.com`
2. Navigate to a shipment at warehouse
3. Create reception report for each laptop:
   - Upload 3 required photos
   - Add inspection notes
   - Submit → Status becomes `pending_approval`
4. Login as `logistics@bairesdev.com`
5. Go to Reception Reports list
6. See reports awaiting approval
7. Review each report:
   - View photos
   - Read inspection notes
   - Approve or request changes
8. After approval → Laptop becomes available for assignment

## Role-Based Testing

### Logistics User
**Can access:**
- ✓ All shipments (all companies, all types)
- ✓ Create new shipments (all three types)
- ✓ Update shipment statuses
- ✓ Approve reception reports
- ✓ Assign engineers to shipments
- ✓ View inventory
- ✓ View dashboard with charts
- ✓ Calendar view

**Test:**
1. Login as `logistics@bairesdev.com`
2. Dashboard → See statistics for all shipments
3. Create Single Shipment → Select client, engineer, JIRA ticket
4. Create Bulk Shipment → Select multiple laptops
5. Create Warehouse-to-Engineer → Select from inventory
6. Approve pending reception reports
7. Update shipment statuses

### Warehouse User
**Can access:**
- ✓ Shipments at/near warehouse
- ✓ Create laptop-based reception reports
- ✓ View inventory at warehouse
- ✓ Cannot approve own reports (requires logistics)

**Test:**
1. Login as `warehouse@bairesdev.com`
2. View shipments → Only see in-transit, at-warehouse, released
3. Select shipment at warehouse → View laptops
4. Create reception report for each laptop:
   - Upload serial number photo
   - Upload external condition photo
   - Upload working condition photo
   - Add inspection notes
5. Submit → Report enters `pending_approval` queue
6. Cannot approve own reports

### Project Manager
**Can access:**
- ✓ View all shipments (read-only)
- ✓ View dashboard and reports
- ✓ View inventory
- ✓ Assign engineers to shipments
- ✓ Cannot create shipments
- ✓ Cannot update statuses

**Test:**
1. Login as `pm@bairesdev.com`
2. Dashboard → View analytics
3. Shipments list → View all, filter by status
4. Shipment detail → Assign engineer
5. Cannot create new shipments (button hidden)
6. Cannot update statuses (form disabled)

### Client User
**Can access:**
- ✓ Only own company's shipments
- ✓ Submit pickup forms
- ✓ Track own shipments
- ✓ Cannot see other companies' data

**Test:**
1. Login as `client@techcorp.com`
2. Shipments → Only see TechCorp shipments
3. Create pickup form for pending shipment
4. View shipment status and tracking
5. Cannot access other companies' shipments
6. Cannot create shipments (logistics only)

## Dashboard & Analytics Testing

**Sample data provides realistic metrics:**

### Shipment Statistics
- Total shipments: 7+
- By status (distributed across all 8 statuses)
- By type (all 3 types represented)
- Average delivery time (calculated from timestamps)

### Laptop Inventory
- Total laptops: 40+
- By brand (Dell, HP, Lenovo, Apple)
- By status (distributed realistically)
- Available vs. assigned

### Workflow Metrics
- Pickup forms submitted: 100% coverage
- Reception reports: Mix of approved and pending
- Delivery confirmations: For all delivered shipments
- Average warehouse processing time

### Charts (if implemented)
- Shipments by status (pie/donut chart)
- Shipments over time (line chart)
- Laptops by brand (bar chart)
- Delivery success rate

## Calendar Testing

The sample data includes realistic date distribution:
- **Historical deliveries** (1-2 months ago)
- **Recent activity** (past week)
- **Current shipments** (today/this week)
- **Scheduled pickups** (tomorrow, next few days)
- **ETAs for in-transit** (1-2 days out)

**Test Calendar View:**
1. Login as logistics/PM user
2. Go to Calendar view
3. Should see:
   - Past deliveries
   - Current shipments in progress
   - Scheduled pickups (future dates)
   - ETAs for in-transit shipments

## Database Query Examples

### View all bulk shipments
```sql
SELECT s.id, s.jira_ticket_number, s.status, s.laptop_count, cc.name
FROM shipments s
JOIN client_companies cc ON cc.id = s.client_company_id
WHERE s.laptop_count > 1
ORDER BY s.created_at DESC;
```

### View reception reports pending approval
```sql
SELECT 
    rr.id,
    l.serial_number,
    l.brand,
    l.model,
    rr.received_at,
    u.email as warehouse_user
FROM reception_reports rr
JOIN laptops l ON l.id = rr.laptop_id
JOIN users u ON u.id = rr.warehouse_user_id
WHERE rr.status = 'pending_approval'
ORDER BY rr.received_at DESC;
```

### View shipments by type
```sql
SELECT 
    shipment_type,
    COUNT(*) as count,
    SUM(laptop_count) as total_laptops,
    ROUND(AVG(laptop_count), 1) as avg_laptops_per_shipment
FROM shipments
GROUP BY shipment_type
ORDER BY count DESC;
```

### View engineers with confirmed vs. unconfirmed addresses
```sql
SELECT 
    address_confirmed,
    COUNT(*) as engineer_count,
    ROUND(100.0 * COUNT(*) / (SELECT COUNT(*) FROM software_engineers), 1) as percentage
FROM software_engineers
GROUP BY address_confirmed;
```

### View audit log activity (last 24 hours)
```sql
SELECT 
    timestamp,
    u.email as user_email,
    action,
    entity_type,
    entity_id,
    details
FROM audit_logs al
JOIN users u ON u.id = al.user_id
WHERE timestamp > NOW() - INTERVAL '24 hours'
ORDER BY timestamp DESC;
```

## Troubleshooting

### Data not loading
```powershell
# Check Docker is running
docker ps

# Check database container
docker ps | findstr laptop-tracking-db

# Start services
docker-compose up -d

# Retry loading
.\scripts\load-sample-data.ps1
```

### Verification fails
```powershell
# Check database connection
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT 1;"

# View logs
docker compose logs -f postgres

# Check for error messages
docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"
```

### Reset and reload
```powershell
# Complete fresh start
docker compose down -v
docker compose up -d

# Wait for database to initialize
Start-Sleep -Seconds 10

# Load data
.\scripts\load-sample-data.ps1
```

## File Structure

```
scripts/
├── comprehensive-sample-data-v2.sql          # Base data (users, companies, engineers, laptops)
├── comprehensive-shipments-data-v2.sql       # Shipments, forms, reports, audit logs
├── load-sample-data.ps1                      # Loading script (updated)
├── verify-test-data.ps1                      # Verification script (updated)
├── COMPREHENSIVE_SAMPLE_DATA_README.md       # This documentation
└── [other legacy scripts...]
```

## Summary

✅ **32 users** across all roles  
✅ **15 client companies** with complete details  
✅ **35+ software engineers** with address tracking  
✅ **40+ laptops** (Dell, HP, Lenovo, Apple)  
✅ **7+ shipments** covering all workflows  
✅ **All 3 shipment types** implemented  
✅ **All 8 status stages** represented  
✅ **Complete forms** with realistic JSON data  
✅ **Laptop-based reception reports** with approval workflow  
✅ **Delivery confirmations** with photos  
✅ **Audit trail** of all activities  
✅ **Magic links** for secure access  
✅ **Realistic timestamps** spanning 6 months  

**Perfect for:**
- Development and testing
- Demonstrations and presentations
- User training
- QA and validation
- Performance testing
- UI/UX evaluation

---

**Last Updated:** 2025-11-15  
**Version:** 2.1  
**Compatible with:** Latest application (laptop-based reception reports, three shipment types)

