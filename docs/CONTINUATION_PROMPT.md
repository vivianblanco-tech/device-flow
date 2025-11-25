# Three Shipment Types Implementation - Continuation Prompt

## Project Context

I'm implementing three distinct shipment types for Align using strict TDD methodology. The project is in progress and I need to continue from where we left off.

## What Has Been Completed

### ✅ Phase 1: Database Schema Changes (COMPLETE)
All 5 sub-phases completed with migrations applied:

**Migration 000016** - Added `shipment_type` enum column:
- Three types: `single_full_journey`, `bulk_to_warehouse`, `warehouse_to_engineer`
- Migrated existing shipments to `single_full_journey`

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
- 35+ tests, all passing

### ✅ Phase 2: Model Layer Updates (COMPLETE)
All 3 sub-phases completed:

**Phase 2.1** - Serial Number Correction Tracking:
- Added fields to `ReceptionReport` model
- `HasSerialNumberCorrection()` and `SerialNumberCorrectionNote()` methods
- 2 tests passing

**Phase 2.2** - Inventory Availability Queries:
- `IsAvailableForWarehouseShipment()` method on Laptop model
- Checks: status (available/at_warehouse), has reception report, not in active shipment
- 6 tests passing

**Phase 2.3** - Laptop Status Synchronization:
- `ShouldSyncLaptopStatus()` - only single shipments sync
- `GetLaptopStatusForShipmentStatus()` - maps shipment status to laptop status
- 8 tests passing

**Total Tests Added:** 50+ tests, all passing
**Database State:** 3 migrations applied (000016, 000017, 000018)

## Implementation Decisions (User-Approved)

1. **Status Flows:**
   - `single_full_journey`: 8 statuses (full flow)
   - `bulk_to_warehouse`: 5 statuses (ends at `at_warehouse`)
   - `warehouse_to_engineer`: 3 statuses (starts from `released_from_warehouse`)

2. **Laptop Assignment:**
   - Serial number: text input
   - Engineer name: text field (assignable anytime before `released_from_warehouse`)
   - Specifications: textarea
   - Laptop record auto-created on form submission for single_full_journey

3. **Serial Number Corrections:**
   - Verify matches pickup form
   - Allow corrections with note/flag
   - Only Logistics users can approve corrections

4. **Bulk Shipments:**
   - Track count only during pickup
   - Create laptop records during warehouse reception (when serial numbers known)

5. **UI Navigation:**
   - Three separate "Create Shipment" buttons for better UX

## What Needs to Be Done Next

### 🔄 Phase 3: Validator Updates (Days 7-8) - NEXT UP

Need to create three separate form validators with TDD:

#### Phase 3.1: Single Full Journey Form Validator
**File:** `internal/validator/single_shipment_form.go`

Required fields:
- Client company, contact info, pickup address/date/time, JIRA ticket
- **Laptop serial number** (required)
- Laptop specs (optional)
- Engineer name (optional)
- Accessories (optional with description if checked)

Validation rules:
- All base pickup form validations
- Serial number cannot be empty
- If accessories included, description required

#### Phase 3.2: Bulk to Warehouse Form Validator
**File:** `internal/validator/bulk_shipment_form.go`

Required fields:
- Client company, contact info, pickup address/date/time, JIRA ticket
- Number of laptops (must be >= 2)
- **Bulk dimensions** (length, width, height, weight - all required and positive)
- Accessories (optional)

Validation rules:
- All base pickup form validations
- Laptop count >= 2
- All bulk dimensions > 0

#### Phase 3.3: Warehouse to Engineer Form Validator
**File:** `internal/validator/warehouse_to_engineer_form.go`

Required fields:
- **Laptop selection** (from available inventory)
- **Software engineer** (required)
- Engineer address (full address required)
- Courier info (optional initially)
- JIRA ticket

Validation rules:
- Laptop ID must be provided
- Engineer must be selected or named
- Full delivery address required

### Phase 4: Handler Layer (Days 9-12)
- Update existing pickup form handler to support single_full_journey
- Create bulk shipment form handler
- Create warehouse-to-engineer form handler
- Update shipment list/detail handlers with type filtering

### Phase 5: Templates & UI (Days 13-15)
- Create/update three form templates
- Add three "Create Shipment" buttons to Dashboard and Shipments page
- Update shipment list with type indicators

### Phase 6: Integration & Testing (Days 16-18)
- End-to-end tests for each type
- Status transition tests
- Laptop sync tests
- Serial number correction workflow tests

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
│   └── pickup_form.go       🔄 Needs splitting into 3 validators
└── handlers/
    ├── pickup_form.go       🔄 Needs update for shipment types
    └── shipments.go         🔄 Needs type filtering

migrations/
├── 000016_add_shipment_type.*                     ✅ Applied
├── 000017_add_laptop_count_to_shipments.*         ✅ Applied
└── 000018_add_serial_number_tracking_to_*.*       ✅ Applied
```

## Reference Documents

- **Full TDD Plan:** `docs/THREE_SHIPMENT_TYPES_TDD_PLAN.md` (2100 lines)
- **Current Models:** All in `internal/models/`
- **Existing Validators:** `internal/validator/`

## Development Environment

- **Language:** Go 1.22+
- **Database:** PostgreSQL in Docker container `laptop-tracking-db`
- **Database Name:** `laptop_tracking_dev`
- **Test Command:** `go test -v ./internal/models -run TestName`
- **Migration Command:** `$env:DB_URL="postgresql://postgres:password@localhost:5432/laptop_tracking_dev?sslmode=disable" ; make migrate-up`

## TDD Workflow to Follow

```
🟥 RED:   Write failing test first
🟩 GREEN: Implement minimal code to pass
✅ Commit only after tests pass
```

**IMPORTANT:** Each feature must follow strict TDD. No implementation without a failing test first.

## Request

Please continue the implementation from **Phase 3: Validator Updates**. Follow the TDD methodology strictly:

1. Start with Phase 3.1: Single Full Journey Form Validator
2. Create tests first (RED)
3. Implement to pass (GREEN)
4. Commit with descriptive message
5. Move to next phase

All previous phases are complete and working. The codebase is ready for Phase 3.

