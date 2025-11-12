# Three Shipment Types Implementation - Continuation Prompt (Phase 4)

## Project Context

I'm implementing three distinct shipment types for a laptop tracking system using strict TDD methodology. The project is in progress and I need to continue from where we left off.

## What Has Been Completed

### ✅ Phase 1: Database Schema Changes (COMPLETE)
All 5 sub-phases completed with migrations applied:

**Migration 000016** - Added `shipment_type` enum column:
- Three types: `single_full_journey`, `bulk_to_warehouse`, `warehouse_to_engineer`
- Migrated existing shipments to `single_full_journey`
- Indexed for filtering by type

**Migration 000017** - Added `laptop_count` column:
- Single shipments: exactly 1 laptop
- Bulk shipments: 2+ laptops
- Check constraint ensures positive count

**Migration 000018** - Added serial number correction tracking to reception_reports:
- `expected_serial_number`, `actual_serial_number`, `serial_number_corrected`
- `correction_note`, `correction_approved_by`
- Partial index for corrections

**Code Completed:**
- `ShipmentType` enum with `IsValidShipmentType()` validation
- `ValidateEngineerAssignment()` - bulk cannot have engineer, warehouse-to-engineer must have engineer
- `GetValidStatusesForType()` - type-specific status flows
- `ValidateLaptopCount()` - enforces count rules per type
- Updated `Validate()` method with all type-specific validations
- **35+ tests, all passing**

### ✅ Phase 2: Model Layer Updates (COMPLETE)
All 3 sub-phases completed:

**Phase 2.1** - Serial Number Correction Tracking:
- Added fields to `ReceptionReport` model
- `HasSerialNumberCorrection()` and `SerialNumberCorrectionNote()` methods
- **2 tests passing**

**Phase 2.2** - Inventory Availability Queries:
- `IsAvailableForWarehouseShipment()` method on Laptop model
- Checks: status (available/at_warehouse), has reception report, not in active shipment
- **6 tests passing**

**Phase 2.3** - Laptop Status Synchronization:
- `ShouldSyncLaptopStatus()` - only single shipments sync
- `GetLaptopStatusForShipmentStatus()` - maps shipment status to laptop status
- **8 tests passing**

### ✅ Phase 3: Validator Updates (COMPLETE)
All 3 sub-phases completed with strict TDD:

**Phase 3.1** - Single Full Journey Form Validator:
- Created `SingleFullJourneyFormInput` struct
- Implemented `ValidateSingleFullJourneyForm()` function
- Required: serial number, client info, pickup details, JIRA ticket
- Optional: laptop specs, engineer name, accessories
- **8 tests passing**
- **Commit:** `feat: add single full journey form validator`

**Phase 3.2** - Bulk to Warehouse Form Validator:
- Created `BulkToWarehouseFormInput` struct
- Implemented `ValidateBulkToWarehouseForm()` function
- Required: laptop count (≥2), bulk dimensions (all positive), client info
- Validates all dimension fields > 0
- **10 tests passing**
- **Commit:** `feat: add bulk to warehouse form validator`

**Phase 3.3** - Warehouse to Engineer Form Validator:
- Created `WarehouseToEngineerFormInput` struct
- Implemented `ValidateWarehouseToEngineerForm()` function
- Required: laptop selection, engineer assignment, delivery address, JIRA ticket
- Optional: courier info (required before shipping, but not on form submission)
- **12 tests passing**
- **Commit:** `feat: add warehouse to engineer form validator`

**Reusable Helper Functions Created:**
- `validateContactInfo()` - Contact name, email, phone validation
- `validateAddress()` - Full address validation with US state/ZIP
- `validatePickupDateTime()` - Date and time slot validation
- `validateJiraTicket()` - JIRA ticket format validation

**Total Tests Added:** 80+ tests across all phases, all passing
**Database State:** 3 migrations applied (000016, 000017, 000018)
**Total Commits:** 6 commits (3 for Phase 1-2, 3 for Phase 3)

## Implementation Decisions (User-Approved)

1. **Status Flows:**
   - `single_full_journey`: 8 statuses (full flow from pending_pickup → delivered)
   - `bulk_to_warehouse`: 5 statuses (ends at `at_warehouse`)
   - `warehouse_to_engineer`: 3 statuses (starts from `released_from_warehouse`)

2. **Laptop Assignment for single_full_journey:**
   - Serial number: text input (REQUIRED)
   - Engineer name: text field (OPTIONAL - assignable anytime before `released_from_warehouse`)
   - Specifications: textarea (OPTIONAL)
   - Laptop record auto-created on form submission

3. **Serial Number Corrections:**
   - Verify matches pickup form
   - Allow corrections with note/flag
   - Only Logistics users can approve corrections

4. **Bulk Shipments:**
   - Track count only during pickup
   - Bulk dimensions MANDATORY (length, width, height, weight all > 0)
   - Create laptop records during warehouse reception (when serial numbers known)

5. **Warehouse-to-Engineer Shipments:**
   - Must select from available inventory
   - Laptop must have completed reception report
   - Engineer assignment REQUIRED
   - Starts at `released_from_warehouse` status

6. **UI Navigation:**
   - Three separate "Create Shipment" buttons for better UX

## What Needs to Be Done Next

### 🔄 Phase 4: Handler Layer Updates (Days 9-12) - NEXT UP

Need to update/create handlers to use new validators and support three shipment types with TDD:

#### Phase 4.1: Update Pickup Form Handler for Single Full Journey
**File:** `internal/handlers/pickup_form.go`

Key changes needed:
- Update `PickupFormSubmit()` to accept `shipment_type` parameter
- Branch logic based on shipment type
- For `single_full_journey`:
  - Use `validator.ValidateSingleFullJourneyForm()`
  - Parse laptop details from form (serial number, specs, engineer name)
  - Auto-create laptop record with status `in_transit_to_warehouse`
  - Link laptop to shipment via `shipment_laptops` table
  - Set `laptop_count = 1`
  - Create shipment with type `single_full_journey`
  - Send notification email to client

Tests needed:
- Valid form submission creates shipment with correct type
- Laptop record is auto-created with serial number
- Engineer can be null (assigned later)
- Validation errors are handled properly
- Email notification is sent

#### Phase 4.2: Create Bulk to Warehouse Form Handler
**File:** `internal/handlers/bulk_shipment_form.go` (NEW) or extend `pickup_form.go`

Key changes needed:
- Use `validator.ValidateBulkToWarehouseForm()`
- Parse laptop count and bulk dimensions
- Set `laptop_count` to form value (≥2)
- Do NOT create laptop records (created during warehouse reception)
- Create shipment with type `bulk_to_warehouse`
- Store bulk dimensions in shipment record
- Send notification email

Tests needed:
- Valid bulk form creates shipment with correct count
- No laptop records created initially
- Bulk dimensions are stored correctly
- Cannot progress past `at_warehouse` status

#### Phase 4.3: Create Warehouse to Engineer Form Handler
**File:** `internal/handlers/warehouse_to_engineer_form.go` (NEW)

Key changes needed:
- Endpoint to GET available laptops for dropdown
- Use `validator.ValidateWarehouseToEngineerForm()`
- Verify laptop has reception report
- Verify laptop is actually available
- Create shipment with type `warehouse_to_engineer`
- Set initial status to `released_from_warehouse`
- Set `laptop_count = 1`
- Link laptop to shipment
- Update laptop status to `in_transit_to_engineer`
- REQUIRE engineer assignment
- Send notification email

Tests needed:
- Only available laptops with reception reports are selectable
- Laptop must have reception report
- Shipment starts at `released_from_warehouse`
- Engineer assignment is required
- Laptop status updates correctly

#### Phase 4.4: Update Shipments List Handler
**File:** `internal/handlers/shipments.go`

Key changes needed:
- Accept `type` query parameter for filtering
- Update SQL query to filter by `shipment_type` if provided
- Return shipment type in JSON/template data
- Display type badge/indicator in list view

Tests needed:
- Filter by `single_full_journey` returns only those shipments
- Filter by `bulk_to_warehouse` returns only bulk shipments
- Filter by `warehouse_to_engineer` returns only warehouse shipments
- No filter returns all shipments

#### Phase 4.5: Update Shipment Detail Handler
**File:** `internal/handlers/shipments.go`

Key changes needed:
- Include `shipment_type` in shipment detail response
- Show type-specific allowed status transitions
- Display laptop details for single shipments
- Display laptop count for bulk shipments
- Show available inventory for warehouse-to-engineer

Tests needed:
- Single shipment shows laptop details
- Bulk shipment shows laptop count
- Status transitions respect type constraints

### Phase 5: Templates & UI (Days 13-15)
- Create/update three form templates
- Add three "Create Shipment" buttons to Dashboard and Shipments page
- Update shipment list with type indicators

### Phase 6: Integration & Testing (Days 16-18)
- End-to-end tests for each type
- Status transition tests
- Laptop sync tests
- Serial number correction workflow tests

### Phase 7: Documentation & Cleanup (Day 18)
- Update README with three shipment types
- Mark implementation complete in plan.md
- Create user guide

## Project Structure

```
internal/
├── models/
│   ├── shipment.go          ✅ Updated with type-specific logic
│   ├── shipment_test.go     ✅ 35+ tests added
│   ├── laptop.go            ✅ Updated with availability logic
│   ├── laptop_test.go       ✅ 6 tests added
│   ├── forms.go             ✅ Updated with serial tracking
│   └── forms_test.go        ✅ 2 tests added
├── validator/
│   ├── single_shipment_form.go       ✅ NEW - Single journey validator
│   ├── single_shipment_form_test.go  ✅ NEW - 8 tests
│   ├── bulk_shipment_form.go         ✅ NEW - Bulk validator
│   ├── bulk_shipment_form_test.go    ✅ NEW - 10 tests
│   ├── warehouse_to_engineer_form.go      ✅ NEW - Warehouse validator
│   ├── warehouse_to_engineer_form_test.go ✅ NEW - 12 tests
│   └── pickup_form.go                ✅ Existing (legacy)
└── handlers/
    ├── pickup_form.go       🔄 NEXT - Update for shipment types
    ├── shipments.go         🔄 NEXT - Add type filtering
    └── ...

migrations/
├── 000016_add_shipment_type.*                     ✅ Applied
├── 000017_add_laptop_count_to_shipments.*         ✅ Applied
└── 000018_add_serial_number_tracking_to_*.*       ✅ Applied
```

## Reference Documents

- **Full TDD Plan:** `docs/THREE_SHIPMENT_TYPES_TDD_PLAN.md` (2100 lines)
- **Original Continuation:** `docs/CONTINUATION_PROMPT.md`
- **Current Models:** All in `internal/models/`
- **Current Validators:** All in `internal/validator/`
- **Current Handlers:** All in `internal/handlers/`

## Development Environment

- **Language:** Go 1.22+
- **Database:** PostgreSQL in Docker container `laptop-tracking-db`
- **Database Name:** `laptop_tracking_dev`
- **Test Command:** `go test -v ./internal/handlers -run TestName`
- **Migration Command:** `$env:DB_URL="postgresql://postgres:password@localhost:5432/laptop_tracking_dev?sslmode=disable" ; make migrate-up`
- **Run App:** Docker containers running
- **Shell:** PowerShell

## TDD Workflow to Follow

```
🟥 RED:   Write failing test first
🟩 GREEN: Implement minimal code to pass
✅ Commit only after tests pass
```

**IMPORTANT:** Each feature must follow strict TDD. No implementation without a failing test first.

## Git Commits So Far

1. `feat: add shipment_type enum and column to shipments table`
2. `feat: add engineer assignment validation based on shipment type`
3. `feat: implement type-specific status flow validation`
4. `feat: add laptop_count field with type-specific validation`
5. `feat: update shipment validation to include type-specific rules`
6. `feat: add serial number correction tracking to reception reports`
7. `feat: add inventory availability queries for warehouse-to-engineer shipments`
8. `feat: add laptop status synchronization logic for shipment types`
9. `feat: add single full journey form validator`
10. `feat: add bulk to warehouse form validator`
11. `feat: add warehouse to engineer form validator`

## Request

Please continue the implementation from **Phase 4: Handler Layer Updates**. Follow the TDD methodology strictly:

1. Start with Phase 4.1: Update Pickup Form Handler for Single Full Journey
2. Create tests first (RED)
3. Implement to pass (GREEN)
4. Commit with descriptive message
5. Move to next sub-phase

All previous phases are complete and working. The codebase is ready for Phase 4.

**Key Points for Phase 4:**
- Update existing pickup form handler to support `single_full_journey`
- Auto-create laptop records for single shipments
- Create new handlers for bulk and warehouse-to-engineer types
- Add shipment type filtering to list/detail views
- Follow strict TDD: test first, then implement
- Commit after each sub-phase passes

The full detailed plan is in `docs/THREE_SHIPMENT_TYPES_TDD_PLAN.md` starting at line 1572 (Phase 4).

