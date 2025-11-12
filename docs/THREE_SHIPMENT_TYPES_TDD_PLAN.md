# Three Shipment Types - TDD Implementation Plan

**Created:** November 12, 2025  
**Status:** Awaiting Approval  
**Estimated Duration:** 12-18 days

---

## Executive Summary

This plan implements three distinct shipment types to replace the current single shipment flow:

1. **`single_full_journey`** (Client → Warehouse → Engineer) - Priority: HIGHEST
   - Exactly 1 laptop per shipment
   - Full 8-status flow
   - Laptop details captured upfront, auto-created
   - Engineer can be assigned anytime before warehouse release

2. **`bulk_to_warehouse`** (Client → Warehouse only)
   - 2+ laptops per shipment
   - Status flow ends at `at_warehouse`
   - Bulk dimensions mandatory
   - Laptops created during warehouse reception

3. **`warehouse_to_engineer`** (Warehouse → Engineer only)
   - Exactly 1 laptop per shipment
   - Starts from `released_from_warehouse` (3-status flow)
   - Requires available laptop with completed reception report
   - Logistics selects from inventory

---

## Design Decisions (User-Approved)

1. **Status Flows:**
   - `warehouse_to_engineer`: `released_from_warehouse` → `in_transit_to_engineer` → `delivered`
   
2. **Laptop Assignment (`single_full_journey`):**
   - Serial number: text input
   - Software Engineer name: text field (assignable anytime before `released_from_warehouse`)
   - Specifications: textarea
   - Laptop record auto-created on form submission

3. **Reception Report Verification:**
   - Verify serial number matches pickup form
   - Allow corrections with note/flag
   - Only Logistics users can update/correct serial numbers

4. **Bulk Shipment Laptops:**
   - Track count only during pickup form
   - Create laptop records during warehouse reception (when serial numbers known)

5. **UI Navigation:**
   - Three separate "Create Shipment" buttons on Dashboard and Shipments page
   - Clear labels: "+ Single Shipment", "+ Bulk to Warehouse", "+ Warehouse to Engineer"

6. **Backward Compatibility:**
   - Migrate existing shipments to `single_full_journey` type
   - Add `shipment_type` column with migration

7. **Current Pickup Form:**
   - Update existing form to handle `single_full_journey` type
   - Remove bulk toggle, add laptop details section

---

## Phase 1: Database Schema Changes (Days 1-3)

### 1.1 Create Shipment Type Enum & Migration

#### 🟥 RED: Test shipment type enum validation
**File:** `internal/models/shipment_test.go`
```go
func TestShipmentType_Validation(t *testing.T) {
    validTypes := []ShipmentType{
        ShipmentTypeSingleFullJourney,
        ShipmentTypeBulkToWarehouse,
        ShipmentTypeWarehouseToEngineer,
    }
    
    for _, shipmentType := range validTypes {
        if !IsValidShipmentType(shipmentType) {
            t.Errorf("Expected %s to be valid", shipmentType)
        }
    }
    
    // Test invalid type
    if IsValidShipmentType("invalid_type") {
        t.Error("Expected invalid_type to be invalid")
    }
}
```

#### 🟩 GREEN: Implement shipment type enum
**File:** `internal/models/shipment.go`
```go
// ShipmentType represents the type of shipment
type ShipmentType string

const (
    ShipmentTypeSingleFullJourney    ShipmentType = "single_full_journey"
    ShipmentTypeBulkToWarehouse      ShipmentType = "bulk_to_warehouse"
    ShipmentTypeWarehouseToEngineer  ShipmentType = "warehouse_to_engineer"
)

// IsValidShipmentType checks if a given shipment type is valid
func IsValidShipmentType(shipmentType ShipmentType) bool {
    switch shipmentType {
    case ShipmentTypeSingleFullJourney,
         ShipmentTypeBulkToWarehouse,
         ShipmentTypeWarehouseToEngineer:
        return true
    }
    return false
}
```

#### Update Shipment struct
Add `ShipmentType` field to `Shipment` struct in `internal/models/shipment.go`:
```go
type Shipment struct {
    ID                  int64           `json:"id" db:"id"`
    ShipmentType        ShipmentType    `json:"shipment_type" db:"shipment_type"`  // NEW
    ClientCompanyID     int64           `json:"client_company_id" db:"client_company_id"`
    // ... rest of fields
}
```

#### Create migration
**File:** `migrations/000016_add_shipment_type.up.sql`
```sql
-- Create shipment_type enum
CREATE TYPE shipment_type AS ENUM (
    'single_full_journey',
    'bulk_to_warehouse',
    'warehouse_to_engineer'
);

-- Add shipment_type column to shipments table
ALTER TABLE shipments ADD COLUMN shipment_type shipment_type;

-- Set default for existing shipments
UPDATE shipments SET shipment_type = 'single_full_journey' WHERE shipment_type IS NULL;

-- Make column NOT NULL after setting defaults
ALTER TABLE shipments ALTER COLUMN shipment_type SET NOT NULL;

-- Set default for new rows
ALTER TABLE shipments ALTER COLUMN shipment_type SET DEFAULT 'single_full_journey';

-- Add index for filtering by type
CREATE INDEX idx_shipments_type ON shipments(shipment_type);

-- Add comment
COMMENT ON COLUMN shipments.shipment_type IS 'Type of shipment: single_full_journey, bulk_to_warehouse, or warehouse_to_engineer';
```

**File:** `migrations/000016_add_shipment_type.down.sql`
```sql
-- Remove index
DROP INDEX IF EXISTS idx_shipments_type;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS shipment_type;

-- Drop enum type
DROP TYPE IF EXISTS shipment_type;
```

**Commit:** `feat: add shipment_type enum and column to shipments table`

---

### 1.2 Add Laptop Assignment Flexibility

#### 🟥 RED: Test engineer assignment validation
**File:** `internal/models/shipment_test.go`
```go
func TestShipment_EngineerAssignmentRules(t *testing.T) {
    tests := []struct {
        name          string
        shipmentType  ShipmentType
        status        ShipmentStatus
        engineerID    *int64
        shouldBeValid bool
    }{
        {
            name:          "single_full_journey can have engineer assigned before release",
            shipmentType:  ShipmentTypeSingleFullJourney,
            status:        ShipmentStatusAtWarehouse,
            engineerID:    nil,
            shouldBeValid: true,
        },
        {
            name:          "warehouse_to_engineer must have engineer assigned",
            shipmentType:  ShipmentTypeWarehouseToEngineer,
            status:        ShipmentStatusReleasedFromWarehouse,
            engineerID:    nil,
            shouldBeValid: false,
        },
        {
            name:          "bulk_to_warehouse cannot have engineer assigned",
            shipmentType:  ShipmentTypeBulkToWarehouse,
            status:        ShipmentStatusAtWarehouse,
            engineerID:    intPtr(1),
            shouldBeValid: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &Shipment{
                ShipmentType:       tt.shipmentType,
                Status:             tt.status,
                SoftwareEngineerID: tt.engineerID,
            }
            
            err := s.ValidateEngineerAssignment()
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid && err == nil {
                t.Error("Expected error, got nil")
            }
        })
    }
}
```

#### 🟩 GREEN: Implement engineer assignment validation
**File:** `internal/models/shipment.go`
```go
// ValidateEngineerAssignment validates engineer assignment based on shipment type
func (s *Shipment) ValidateEngineerAssignment() error {
    switch s.ShipmentType {
    case ShipmentTypeBulkToWarehouse:
        // Bulk shipments cannot have engineer assigned
        if s.SoftwareEngineerID != nil {
            return errors.New("bulk_to_warehouse shipments cannot have software engineer assigned")
        }
    case ShipmentTypeWarehouseToEngineer:
        // Warehouse-to-engineer shipments must have engineer assigned
        if s.SoftwareEngineerID == nil {
            return errors.New("warehouse_to_engineer shipments must have software engineer assigned")
        }
    case ShipmentTypeSingleFullJourney:
        // Single full journey can be assigned anytime (optional validation here)
        // No error - engineer can be nil or assigned
    }
    return nil
}
```

**Commit:** `feat: add engineer assignment validation based on shipment type`

---

### 1.3 Update Status Flow Validation

#### 🟥 RED: Test type-specific status flows
**File:** `internal/models/shipment_test.go`
```go
func TestShipment_TypeSpecificStatusFlows(t *testing.T) {
    tests := []struct {
        name         string
        shipmentType ShipmentType
        currentStatus ShipmentStatus
        nextStatus   ShipmentStatus
        shouldBeValid bool
    }{
        {
            name:          "single_full_journey allows full flow to delivered",
            shipmentType:  ShipmentTypeSingleFullJourney,
            currentStatus: ShipmentStatusInTransitToEngineer,
            nextStatus:    ShipmentStatusDelivered,
            shouldBeValid: true,
        },
        {
            name:          "bulk_to_warehouse cannot go past at_warehouse",
            shipmentType:  ShipmentTypeBulkToWarehouse,
            currentStatus: ShipmentStatusAtWarehouse,
            nextStatus:    ShipmentStatusReleasedFromWarehouse,
            shouldBeValid: false,
        },
        {
            name:          "warehouse_to_engineer starts from released_from_warehouse",
            shipmentType:  ShipmentTypeWarehouseToEngineer,
            currentStatus: ShipmentStatusReleasedFromWarehouse,
            nextStatus:    ShipmentStatusInTransitToEngineer,
            shouldBeValid: true,
        },
        {
            name:          "warehouse_to_engineer cannot have earlier statuses",
            shipmentType:  ShipmentTypeWarehouseToEngineer,
            currentStatus: ShipmentStatusPendingPickup,
            nextStatus:    ShipmentStatusPickupScheduled,
            shouldBeValid: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &Shipment{
                ShipmentType: tt.shipmentType,
                Status:       tt.currentStatus,
            }
            
            isValid := s.IsValidStatusTransition(tt.nextStatus)
            if tt.shouldBeValid && !isValid {
                t.Error("Expected transition to be valid")
            }
            if !tt.shouldBeValid && isValid {
                t.Error("Expected transition to be invalid")
            }
        })
    }
}
```

#### 🟩 GREEN: Implement type-specific status flow validation
**File:** `internal/models/shipment.go`
```go
// GetValidStatusesForType returns valid statuses for a shipment type
func GetValidStatusesForType(shipmentType ShipmentType) []ShipmentStatus {
    switch shipmentType {
    case ShipmentTypeSingleFullJourney:
        return []ShipmentStatus{
            ShipmentStatusPendingPickup,
            ShipmentStatusPickupScheduled,
            ShipmentStatusPickedUpFromClient,
            ShipmentStatusInTransitToWarehouse,
            ShipmentStatusAtWarehouse,
            ShipmentStatusReleasedFromWarehouse,
            ShipmentStatusInTransitToEngineer,
            ShipmentStatusDelivered,
        }
    case ShipmentTypeBulkToWarehouse:
        return []ShipmentStatus{
            ShipmentStatusPendingPickup,
            ShipmentStatusPickupScheduled,
            ShipmentStatusPickedUpFromClient,
            ShipmentStatusInTransitToWarehouse,
            ShipmentStatusAtWarehouse,
        }
    case ShipmentTypeWarehouseToEngineer:
        return []ShipmentStatus{
            ShipmentStatusReleasedFromWarehouse,
            ShipmentStatusInTransitToEngineer,
            ShipmentStatusDelivered,
        }
    default:
        return []ShipmentStatus{}
    }
}

// IsValidStatusForType checks if a status is valid for the shipment type
func (s *Shipment) IsValidStatusForType(status ShipmentStatus) bool {
    validStatuses := GetValidStatusesForType(s.ShipmentType)
    for _, validStatus := range validStatuses {
        if status == validStatus {
            return true
        }
    }
    return false
}

// Update IsValidStatusTransition to consider shipment type
func (s *Shipment) IsValidStatusTransition(newStatus ShipmentStatus) bool {
    // First check if the new status is valid for this shipment type
    if !s.IsValidStatusForType(newStatus) {
        return false
    }
    
    // Then check if it's the immediate next status
    allowedStatuses := s.GetNextAllowedStatuses()
    for _, allowed := range allowedStatuses {
        if allowed == newStatus {
            return true
        }
    }
    
    return false
}

// Update GetNextAllowedStatuses to consider shipment type
func (s *Shipment) GetNextAllowedStatuses() []ShipmentStatus {
    // Get valid statuses for this type
    validStatuses := GetValidStatusesForType(s.ShipmentType)
    
    // Find current status index
    currentIndex := -1
    for i, status := range validStatuses {
        if status == s.Status {
            currentIndex = i
            break
        }
    }
    
    // Return next status if available
    if currentIndex >= 0 && currentIndex < len(validStatuses)-1 {
        return []ShipmentStatus{validStatuses[currentIndex+1]}
    }
    
    return []ShipmentStatus{} // No next status available
}
```

**Commit:** `feat: implement type-specific status flow validation`

---

### 1.4 Add Laptop Count Tracking

#### 🟥 RED: Test laptop count field
**File:** `internal/models/shipment_test.go`
```go
func TestShipment_LaptopCountTracking(t *testing.T) {
    tests := []struct {
        name          string
        shipmentType  ShipmentType
        laptopCount   int
        shouldBeValid bool
    }{
        {
            name:          "single_full_journey must have count = 1",
            shipmentType:  ShipmentTypeSingleFullJourney,
            laptopCount:   1,
            shouldBeValid: true,
        },
        {
            name:          "single_full_journey cannot have count > 1",
            shipmentType:  ShipmentTypeSingleFullJourney,
            laptopCount:   2,
            shouldBeValid: false,
        },
        {
            name:          "bulk_to_warehouse must have count >= 2",
            shipmentType:  ShipmentTypeBulkToWarehouse,
            laptopCount:   2,
            shouldBeValid: true,
        },
        {
            name:          "bulk_to_warehouse cannot have count = 1",
            shipmentType:  ShipmentTypeBulkToWarehouse,
            laptopCount:   1,
            shouldBeValid: false,
        },
        {
            name:          "warehouse_to_engineer must have count = 1",
            shipmentType:  ShipmentTypeWarehouseToEngineer,
            laptopCount:   1,
            shouldBeValid: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &Shipment{
                ShipmentType: tt.shipmentType,
                LaptopCount:  tt.laptopCount,
            }
            
            err := s.ValidateLaptopCount()
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid && err == nil {
                t.Error("Expected error, got nil")
            }
        })
    }
}
```

#### 🟩 GREEN: Add laptop count field and validation
**File:** `internal/models/shipment.go`
```go
type Shipment struct {
    // ... existing fields ...
    LaptopCount  int  `json:"laptop_count" db:"laptop_count"`  // NEW
    // ... rest of fields ...
}

// ValidateLaptopCount validates laptop count based on shipment type
func (s *Shipment) ValidateLaptopCount() error {
    switch s.ShipmentType {
    case ShipmentTypeSingleFullJourney, ShipmentTypeWarehouseToEngineer:
        if s.LaptopCount != 1 {
            return errors.New("single shipments must have exactly 1 laptop")
        }
    case ShipmentTypeBulkToWarehouse:
        if s.LaptopCount < 2 {
            return errors.New("bulk shipments must have at least 2 laptops")
        }
    }
    return nil
}
```

#### Create migration
**File:** `migrations/000017_add_laptop_count_to_shipments.up.sql`
```sql
-- Add laptop_count column
ALTER TABLE shipments ADD COLUMN laptop_count INTEGER;

-- Set default count based on shipment type
UPDATE shipments SET laptop_count = 1 WHERE shipment_type = 'single_full_journey';
UPDATE shipments SET laptop_count = 1 WHERE shipment_type = 'warehouse_to_engineer';

-- For bulk shipments without laptops, set to 0 temporarily
UPDATE shipments SET laptop_count = 0 WHERE shipment_type = 'bulk_to_warehouse' AND laptop_count IS NULL;

-- Count existing laptops for bulk shipments that have them
UPDATE shipments s
SET laptop_count = (
    SELECT COUNT(*) FROM shipment_laptops sl WHERE sl.shipment_id = s.id
)
WHERE shipment_type = 'bulk_to_warehouse' AND laptop_count = 0;

-- Make column NOT NULL
ALTER TABLE shipments ALTER COLUMN laptop_count SET NOT NULL;

-- Add check constraint
ALTER TABLE shipments ADD CONSTRAINT chk_laptop_count_positive CHECK (laptop_count > 0);

-- Add comment
COMMENT ON COLUMN shipments.laptop_count IS 'Number of laptops in this shipment (must be 1 for single, 2+ for bulk)';
```

**File:** `migrations/000017_add_laptop_count_to_shipments.down.sql`
```sql
-- Remove constraint
ALTER TABLE shipments DROP CONSTRAINT IF EXISTS chk_laptop_count_positive;

-- Remove column
ALTER TABLE shipments DROP COLUMN IF EXISTS laptop_count;
```

**Commit:** `feat: add laptop_count field with type-specific validation`

---

### 1.5 Update Shipment Validation

#### 🟥 RED: Test complete shipment validation
**File:** `internal/models/shipment_test.go`
```go
func TestShipment_ValidateWithType(t *testing.T) {
    tests := []struct {
        name          string
        shipment      Shipment
        shouldBeValid bool
        errorContains string
    }{
        {
            name: "valid single_full_journey",
            shipment: Shipment{
                ShipmentType:     ShipmentTypeSingleFullJourney,
                ClientCompanyID:  1,
                LaptopCount:      1,
                Status:           ShipmentStatusPendingPickup,
                JiraTicketNumber: "SCOP-12345",
            },
            shouldBeValid: true,
        },
        {
            name: "invalid single_full_journey - wrong laptop count",
            shipment: Shipment{
                ShipmentType:     ShipmentTypeSingleFullJourney,
                ClientCompanyID:  1,
                LaptopCount:      2,
                Status:           ShipmentStatusPendingPickup,
                JiraTicketNumber: "SCOP-12345",
            },
            shouldBeValid: false,
            errorContains: "exactly 1 laptop",
        },
        {
            name: "valid bulk_to_warehouse",
            shipment: Shipment{
                ShipmentType:     ShipmentTypeBulkToWarehouse,
                ClientCompanyID:  1,
                LaptopCount:      5,
                Status:           ShipmentStatusAtWarehouse,
                JiraTicketNumber: "SCOP-12345",
            },
            shouldBeValid: true,
        },
        {
            name: "invalid warehouse_to_engineer - missing engineer",
            shipment: Shipment{
                ShipmentType:       ShipmentTypeWarehouseToEngineer,
                ClientCompanyID:    1,
                LaptopCount:        1,
                Status:             ShipmentStatusReleasedFromWarehouse,
                JiraTicketNumber:   "SCOP-12345",
                SoftwareEngineerID: nil,
            },
            shouldBeValid: false,
            errorContains: "must have software engineer assigned",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.shipment.Validate()
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid {
                if err == nil {
                    t.Error("Expected error, got nil")
                } else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
                    t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
                }
            }
        })
    }
}
```

#### 🟩 GREEN: Update Validate method
**File:** `internal/models/shipment.go`
```go
// Validate validates the Shipment model
func (s *Shipment) Validate() error {
    // Shipment type validation
    if s.ShipmentType == "" {
        return errors.New("shipment type is required")
    }
    if !IsValidShipmentType(s.ShipmentType) {
        return errors.New("invalid shipment type")
    }

    // Client company validation
    if s.ClientCompanyID == 0 {
        return errors.New("client company ID is required")
    }

    // Status validation
    if s.Status == "" {
        return errors.New("status is required")
    }
    if !IsValidShipmentStatus(s.Status) {
        return errors.New("invalid status")
    }
    
    // Validate status is valid for this shipment type
    if !s.IsValidStatusForType(s.Status) {
        return fmt.Errorf("status %s is not valid for shipment type %s", s.Status, s.ShipmentType)
    }

    // JIRA ticket validation
    if s.JiraTicketNumber == "" {
        return errors.New("JIRA ticket number is required")
    }
    if !IsValidJiraTicketFormat(s.JiraTicketNumber) {
        return errors.New("JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)")
    }

    // Laptop count validation
    if err := s.ValidateLaptopCount(); err != nil {
        return err
    }

    // Engineer assignment validation
    if err := s.ValidateEngineerAssignment(); err != nil {
        return err
    }

    return nil
}
```

**Commit:** `feat: update shipment validation to include type-specific rules`

---

## Phase 2: Model Layer Updates (Days 4-6)

### 2.1 Add Laptop Serial Number Correction Tracking

#### 🟥 RED: Test serial number correction tracking
**File:** `internal/models/reception_report_test.go`
```go
func TestReceptionReport_SerialNumberCorrection(t *testing.T) {
    report := &ReceptionReport{
        ShipmentID:              1,
        WarehouseUserID:         1,
        ExpectedSerialNumber:    "ABC123",
        ActualSerialNumber:      "ABC456",
        SerialNumberCorrected:   true,
        CorrectionNote:          "Serial number mismatch - updated to match physical device",
        CorrectionApprovedBy:    intPtr(2),
    }
    
    if !report.HasSerialNumberCorrection() {
        t.Error("Expected serial number correction to be detected")
    }
    
    if report.SerialNumberCorrectionNote() == "" {
        t.Error("Expected correction note to be present")
    }
}

func TestReceptionReport_NoCorrection(t *testing.T) {
    report := &ReceptionReport{
        ShipmentID:           1,
        WarehouseUserID:      1,
        ExpectedSerialNumber: "ABC123",
        ActualSerialNumber:   "ABC123",
        SerialNumberCorrected: false,
    }
    
    if report.HasSerialNumberCorrection() {
        t.Error("Expected no serial number correction")
    }
}
```

#### 🟩 GREEN: Add correction tracking to ReceptionReport model
**File:** `internal/models/reception_report.go`
```go
type ReceptionReport struct {
    // ... existing fields ...
    ExpectedSerialNumber  string     `json:"expected_serial_number,omitempty" db:"expected_serial_number"`
    ActualSerialNumber    string     `json:"actual_serial_number" db:"actual_serial_number"`
    SerialNumberCorrected bool       `json:"serial_number_corrected" db:"serial_number_corrected"`
    CorrectionNote        string     `json:"correction_note,omitempty" db:"correction_note"`
    CorrectionApprovedBy  *int64     `json:"correction_approved_by,omitempty" db:"correction_approved_by"`
    // ... rest of fields ...
}

// HasSerialNumberCorrection returns true if serial number was corrected
func (r *ReceptionReport) HasSerialNumberCorrection() bool {
    return r.SerialNumberCorrected
}

// SerialNumberCorrectionNote returns the correction note
func (r *ReceptionReport) SerialNumberCorrectionNote() string {
    return r.CorrectionNote
}
```

#### Create migration
**File:** `migrations/000018_add_serial_number_tracking_to_reception_reports.up.sql`
```sql
-- Add serial number tracking columns to reception_reports
ALTER TABLE reception_reports ADD COLUMN expected_serial_number VARCHAR(255);
ALTER TABLE reception_reports ADD COLUMN actual_serial_number VARCHAR(255);
ALTER TABLE reception_reports ADD COLUMN serial_number_corrected BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE reception_reports ADD COLUMN correction_note TEXT;
ALTER TABLE reception_reports ADD COLUMN correction_approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL;

-- Add index for finding corrected serial numbers
CREATE INDEX idx_reception_reports_serial_corrected ON reception_reports(serial_number_corrected) WHERE serial_number_corrected = TRUE;

-- Add comments
COMMENT ON COLUMN reception_reports.expected_serial_number IS 'Serial number from pickup form';
COMMENT ON COLUMN reception_reports.actual_serial_number IS 'Actual serial number received at warehouse';
COMMENT ON COLUMN reception_reports.serial_number_corrected IS 'Whether serial number was corrected from expected';
COMMENT ON COLUMN reception_reports.correction_note IS 'Note explaining why serial number was corrected';
COMMENT ON COLUMN reception_reports.correction_approved_by IS 'User ID (Logistics) who approved the correction';
```

**File:** `migrations/000018_add_serial_number_tracking_to_reception_reports.down.sql`
```sql
-- Remove index
DROP INDEX IF EXISTS idx_reception_reports_serial_corrected;

-- Remove columns
ALTER TABLE reception_reports DROP COLUMN IF EXISTS correction_approved_by;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS correction_note;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS serial_number_corrected;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS actual_serial_number;
ALTER TABLE reception_reports DROP COLUMN IF EXISTS expected_serial_number;
```

**Commit:** `feat: add serial number correction tracking to reception reports`

---

### 2.2 Add Inventory Availability Queries

#### 🟥 RED: Test available laptop queries for warehouse-to-engineer
**File:** `internal/models/laptop_test.go`
```go
func TestLaptop_GetAvailableLaptopsForWarehouseShipment(t *testing.T) {
    // This test would need database setup
    // Testing the query logic
    
    // Available laptop must have:
    // 1. status = 'available' OR 'at_warehouse'
    // 2. Has a completed reception report
    // 3. Not currently in any active shipment
    // 4. ClientCompanyID can be null or assigned
}

func TestLaptop_IsAvailableForWarehouseShipment(t *testing.T) {
    tests := []struct {
        name            string
        laptop          Laptop
        hasReception    bool
        inActiveShipment bool
        shouldBeAvailable bool
    }{
        {
            name: "available laptop with reception",
            laptop: Laptop{
                Status: LaptopStatusAvailable,
            },
            hasReception:      true,
            inActiveShipment:  false,
            shouldBeAvailable: true,
        },
        {
            name: "available laptop without reception",
            laptop: Laptop{
                Status: LaptopStatusAvailable,
            },
            hasReception:      false,
            inActiveShipment:  false,
            shouldBeAvailable: false,
        },
        {
            name: "laptop in active shipment",
            laptop: Laptop{
                Status: LaptopStatusAvailable,
            },
            hasReception:      true,
            inActiveShipment:  true,
            shouldBeAvailable: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            available := tt.laptop.IsAvailableForWarehouseShipment(tt.hasReception, tt.inActiveShipment)
            if available != tt.shouldBeAvailable {
                t.Errorf("Expected available=%v, got %v", tt.shouldBeAvailable, available)
            }
        })
    }
}
```

#### 🟩 GREEN: Add availability helper methods
**File:** `internal/models/laptop.go`
```go
// IsAvailableForWarehouseShipment checks if laptop can be used for warehouse-to-engineer shipment
func (l *Laptop) IsAvailableForWarehouseShipment(hasReceptionReport bool, inActiveShipment bool) bool {
    // Must be available or at warehouse
    if l.Status != LaptopStatusAvailable && l.Status != LaptopStatusAtWarehouse {
        return false
    }
    
    // Must have completed reception report
    if !hasReceptionReport {
        return false
    }
    
    // Must not be in an active shipment
    if inActiveShipment {
        return false
    }
    
    return true
}
```

**File:** `internal/models/laptop_queries.go` (NEW FILE)
```go
package models

import (
    "context"
    "database/sql"
)

// GetAvailableLaptopsForWarehouseShipment returns laptops available for warehouse-to-engineer shipments
func GetAvailableLaptopsForWarehouseShipment(ctx context.Context, db *sql.DB) ([]Laptop, error) {
    query := `
        SELECT DISTINCT l.id, l.serial_number, l.sku, l.brand, l.model, l.specs, 
               l.status, l.client_company_id, l.software_engineer_id, 
               l.created_at, l.updated_at,
               cc.name as client_company_name,
               se.name as software_engineer_name,
               se.employee_id
        FROM laptops l
        LEFT JOIN client_companies cc ON cc.id = l.client_company_id
        LEFT JOIN software_engineers se ON se.id = l.software_engineer_id
        WHERE l.status IN ('available', 'at_warehouse')
          -- Must have a reception report
          AND EXISTS (
              SELECT 1 FROM reception_reports rr
              JOIN shipments s ON s.id = rr.shipment_id
              JOIN shipment_laptops sl ON sl.shipment_id = s.id
              WHERE sl.laptop_id = l.id
          )
          -- Must not be in any active shipment (excluding delivered)
          AND NOT EXISTS (
              SELECT 1 FROM shipment_laptops sl
              JOIN shipments s ON s.id = sl.shipment_id
              WHERE sl.laptop_id = l.id
                AND s.status != 'delivered'
          )
        ORDER BY l.created_at DESC
    `
    
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var laptops []Laptop
    for rows.Next() {
        var l Laptop
        var clientCompanyName sql.NullString
        var softwareEngineerName sql.NullString
        var employeeID sql.NullString
        
        err := rows.Scan(
            &l.ID, &l.SerialNumber, &l.SKU, &l.Brand, &l.Model, &l.Specs,
            &l.Status, &l.ClientCompanyID, &l.SoftwareEngineerID,
            &l.CreatedAt, &l.UpdatedAt,
            &clientCompanyName, &softwareEngineerName, &employeeID,
        )
        if err != nil {
            continue
        }
        
        l.ClientCompanyName = clientCompanyName.String
        l.SoftwareEngineerName = softwareEngineerName.String
        l.EmployeeID = employeeID.String
        
        laptops = append(laptops, l)
    }
    
    return laptops, nil
}
```

**Commit:** `feat: add inventory availability queries for warehouse-to-engineer shipments`

---

### 2.3 Add Laptop Status Synchronization

#### 🟥 RED: Test laptop status sync with shipment type
**File:** `internal/models/shipment_test.go`
```go
func TestShipment_SyncLaptopStatusOnUpdate(t *testing.T) {
    tests := []struct {
        name             string
        shipmentType     ShipmentType
        shipmentStatus   ShipmentStatus
        expectedLaptopStatus LaptopStatus
        shouldSync       bool
    }{
        {
            name:                 "single_full_journey syncs laptop status",
            shipmentType:         ShipmentTypeSingleFullJourney,
            shipmentStatus:       ShipmentStatusInTransitToWarehouse,
            expectedLaptopStatus: LaptopStatusInTransitToWarehouse,
            shouldSync:           true,
        },
        {
            name:           "bulk_to_warehouse does not sync laptop status",
            shipmentType:   ShipmentTypeBulkToWarehouse,
            shipmentStatus: ShipmentStatusInTransitToWarehouse,
            shouldSync:     false,
        },
        {
            name:                 "warehouse_to_engineer syncs laptop status",
            shipmentType:         ShipmentTypeWarehouseToEngineer,
            shipmentStatus:       ShipmentStatusInTransitToEngineer,
            expectedLaptopStatus: LaptopStatusInTransitToEngineer,
            shouldSync:           true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &Shipment{
                ShipmentType: tt.shipmentType,
                Status:       tt.shipmentStatus,
            }
            
            shouldSync := s.ShouldSyncLaptopStatus()
            if shouldSync != tt.shouldSync {
                t.Errorf("Expected shouldSync=%v, got %v", tt.shouldSync, shouldSync)
            }
            
            if tt.shouldSync {
                laptopStatus := s.GetLaptopStatusForShipmentStatus()
                if laptopStatus != tt.expectedLaptopStatus {
                    t.Errorf("Expected laptop status %s, got %s", tt.expectedLaptopStatus, laptopStatus)
                }
            }
        })
    }
}
```

#### 🟩 GREEN: Implement laptop status synchronization
**File:** `internal/models/shipment.go`
```go
// ShouldSyncLaptopStatus returns true if laptop status should sync with shipment status
func (s *Shipment) ShouldSyncLaptopStatus() bool {
    // Only single shipments sync laptop status
    return s.ShipmentType == ShipmentTypeSingleFullJourney ||
           s.ShipmentType == ShipmentTypeWarehouseToEngineer
}

// GetLaptopStatusForShipmentStatus returns the corresponding laptop status
func (s *Shipment) GetLaptopStatusForShipmentStatus() LaptopStatus {
    if !s.ShouldSyncLaptopStatus() {
        return "" // No sync for bulk shipments
    }
    
    switch s.Status {
    case ShipmentStatusPendingPickup, ShipmentStatusPickupScheduled,
         ShipmentStatusPickedUpFromClient, ShipmentStatusInTransitToWarehouse:
        return LaptopStatusInTransitToWarehouse
    case ShipmentStatusAtWarehouse:
        return LaptopStatusAtWarehouse
    case ShipmentStatusReleasedFromWarehouse, ShipmentStatusInTransitToEngineer:
        return LaptopStatusInTransitToEngineer
    case ShipmentStatusDelivered:
        return LaptopStatusDelivered
    default:
        return LaptopStatusAvailable
    }
}
```

**Commit:** `feat: add laptop status synchronization logic for shipment types`

---

## Phase 3: Validator Updates (Days 7-8)

### 3.1 Create Single Full Journey Form Validator

#### 🟥 RED: Test single full journey form validation
**File:** `internal/validator/single_shipment_form_test.go` (NEW FILE)
```go
package validator

import (
    "testing"
)

func TestValidateSingleFullJourneyForm(t *testing.T) {
    tests := []struct {
        name          string
        input         SingleFullJourneyFormInput
        shouldBeValid bool
        errorContains string
    }{
        {
            name: "valid single full journey form",
            input: SingleFullJourneyFormInput{
                ClientCompanyID:     1,
                ContactName:         "John Doe",
                ContactEmail:        "john@company.com",
                ContactPhone:        "+1-555-0123",
                PickupAddress:       "123 Main St",
                PickupCity:          "New York",
                PickupState:         "NY",
                PickupZip:           "10001",
                PickupDate:          "2025-11-15",
                PickupTimeSlot:      "morning",
                JiraTicketNumber:    "SCOP-12345",
                LaptopSerialNumber:  "ABC123456",
                LaptopSpecs:         "Dell XPS 15, 16GB RAM, 512GB SSD",
                EngineerName:        "Jane Smith",
            },
            shouldBeValid: true,
        },
        {
            name: "missing laptop serial number",
            input: SingleFullJourneyFormInput{
                ClientCompanyID:  1,
                ContactName:      "John Doe",
                ContactEmail:     "john@company.com",
                ContactPhone:     "+1-555-0123",
                PickupAddress:    "123 Main St",
                PickupCity:       "New York",
                PickupState:      "NY",
                PickupZip:        "10001",
                PickupDate:       "2025-11-15",
                PickupTimeSlot:   "morning",
                JiraTicketNumber: "SCOP-12345",
                LaptopSpecs:      "Dell XPS 15",
                EngineerName:     "Jane Smith",
            },
            shouldBeValid: false,
            errorContains: "serial number is required",
        },
        {
            name: "engineer name optional",
            input: SingleFullJourneyFormInput{
                ClientCompanyID:    1,
                ContactName:        "John Doe",
                ContactEmail:       "john@company.com",
                ContactPhone:       "+1-555-0123",
                PickupAddress:      "123 Main St",
                PickupCity:         "New York",
                PickupState:        "NY",
                PickupZip:          "10001",
                PickupDate:         "2025-11-15",
                PickupTimeSlot:     "morning",
                JiraTicketNumber:   "SCOP-12345",
                LaptopSerialNumber: "ABC123456",
                LaptopSpecs:        "Dell XPS 15",
                EngineerName:       "", // Optional
            },
            shouldBeValid: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateSingleFullJourneyForm(tt.input)
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid {
                if err == nil {
                    t.Error("Expected error, got nil")
                } else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
                    t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
                }
            }
        })
    }
}
```

#### 🟩 GREEN: Implement single full journey form validator
**File:** `internal/validator/single_shipment_form.go` (NEW FILE)
```go
package validator

import (
    "errors"
    "regexp"
    "strings"
)

// SingleFullJourneyFormInput represents the pickup form data for single full journey shipments
type SingleFullJourneyFormInput struct {
    ClientCompanyID     int64
    ContactName         string
    ContactEmail        string
    ContactPhone        string
    PickupAddress       string
    PickupCity          string
    PickupState         string
    PickupZip           string
    PickupDate          string
    PickupTimeSlot      string
    JiraTicketNumber    string
    SpecialInstructions string
    
    // Laptop information (required)
    LaptopSerialNumber  string
    LaptopSpecs         string
    
    // Engineer assignment (optional - can be assigned later)
    EngineerName        string
    
    // Accessories (optional)
    IncludeAccessories     bool
    AccessoriesDescription string
}

// ValidateSingleFullJourneyForm validates the single full journey pickup form
func ValidateSingleFullJourneyForm(input SingleFullJourneyFormInput) error {
    // Client company validation
    if input.ClientCompanyID == 0 {
        return errors.New("client company is required")
    }
    
    // Contact information validation
    if err := validateContactInfo(input.ContactName, input.ContactEmail, input.ContactPhone); err != nil {
        return err
    }
    
    // Pickup address validation
    if err := validateAddress(input.PickupAddress, input.PickupCity, input.PickupState, input.PickupZip); err != nil {
        return err
    }
    
    // Pickup date and time validation
    if err := validatePickupDateTime(input.PickupDate, input.PickupTimeSlot); err != nil {
        return err
    }
    
    // JIRA ticket validation
    if err := validateJiraTicket(input.JiraTicketNumber); err != nil {
        return err
    }
    
    // Laptop serial number validation (REQUIRED)
    if strings.TrimSpace(input.LaptopSerialNumber) == "" {
        return errors.New("laptop serial number is required")
    }
    
    // Laptop specifications (optional but recommended)
    // No validation - can be empty
    
    // Engineer name (optional - can be assigned later)
    // No validation - can be empty
    
    // Accessories validation
    if input.IncludeAccessories && strings.TrimSpace(input.AccessoriesDescription) == "" {
        return errors.New("accessories description is required when including accessories")
    }
    
    return nil
}
```

**Commit:** `feat: add single full journey form validator`

---

### 3.2 Create Bulk to Warehouse Form Validator

#### 🟥 RED: Test bulk to warehouse form validation
**File:** `internal/validator/bulk_shipment_form_test.go` (NEW FILE)
```go
package validator

import (
    "testing"
)

func TestValidateBulkToWarehouseForm(t *testing.T) {
    tests := []struct {
        name          string
        input         BulkToWarehouseFormInput
        shouldBeValid bool
        errorContains string
    }{
        {
            name: "valid bulk form",
            input: BulkToWarehouseFormInput{
                ClientCompanyID:  1,
                ContactName:      "John Doe",
                ContactEmail:     "john@company.com",
                ContactPhone:     "+1-555-0123",
                PickupAddress:    "123 Main St",
                PickupCity:       "New York",
                PickupState:      "NY",
                PickupZip:        "10001",
                PickupDate:       "2025-11-15",
                PickupTimeSlot:   "morning",
                JiraTicketNumber: "SCOP-12345",
                NumberOfLaptops:  5,
                BulkLength:       30.0,
                BulkWidth:        20.0,
                BulkHeight:       15.0,
                BulkWeight:       50.0,
            },
            shouldBeValid: true,
        },
        {
            name: "invalid - laptop count too low",
            input: BulkToWarehouseFormInput{
                ClientCompanyID:  1,
                ContactName:      "John Doe",
                ContactEmail:     "john@company.com",
                ContactPhone:     "+1-555-0123",
                PickupAddress:    "123 Main St",
                PickupCity:       "New York",
                PickupState:      "NY",
                PickupZip:        "10001",
                PickupDate:       "2025-11-15",
                PickupTimeSlot:   "morning",
                JiraTicketNumber: "SCOP-12345",
                NumberOfLaptops:  1, // Too low for bulk
                BulkLength:       30.0,
                BulkWidth:        20.0,
                BulkHeight:       15.0,
                BulkWeight:       50.0,
            },
            shouldBeValid: false,
            errorContains: "at least 2 laptops",
        },
        {
            name: "invalid - missing bulk dimensions",
            input: BulkToWarehouseFormInput{
                ClientCompanyID:  1,
                ContactName:      "John Doe",
                ContactEmail:     "john@company.com",
                ContactPhone:     "+1-555-0123",
                PickupAddress:    "123 Main St",
                PickupCity:       "New York",
                PickupState:      "NY",
                PickupZip:        "10001",
                PickupDate:       "2025-11-15",
                PickupTimeSlot:   "morning",
                JiraTicketNumber: "SCOP-12345",
                NumberOfLaptops:  5,
                BulkLength:       0, // Missing
                BulkWidth:        20.0,
                BulkHeight:       15.0,
                BulkWeight:       50.0,
            },
            shouldBeValid: false,
            errorContains: "bulk dimensions are required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateBulkToWarehouseForm(tt.input)
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid {
                if err == nil {
                    t.Error("Expected error, got nil")
                } else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
                    t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
                }
            }
        })
    }
}
```

#### 🟩 GREEN: Implement bulk to warehouse form validator
**File:** `internal/validator/bulk_shipment_form.go` (NEW FILE)
```go
package validator

import (
    "errors"
)

// BulkToWarehouseFormInput represents the pickup form data for bulk shipments
type BulkToWarehouseFormInput struct {
    ClientCompanyID     int64
    ContactName         string
    ContactEmail        string
    ContactPhone        string
    PickupAddress       string
    PickupCity          string
    PickupState         string
    PickupZip           string
    PickupDate          string
    PickupTimeSlot      string
    JiraTicketNumber    string
    SpecialInstructions string
    
    // Laptop count (must be >= 2)
    NumberOfLaptops  int
    
    // Bulk dimensions (REQUIRED)
    BulkLength  float64
    BulkWidth   float64
    BulkHeight  float64
    BulkWeight  float64
    
    // Accessories (optional)
    IncludeAccessories     bool
    AccessoriesDescription string
}

// ValidateBulkToWarehouseForm validates the bulk to warehouse pickup form
func ValidateBulkToWarehouseForm(input BulkToWarehouseFormInput) error {
    // Client company validation
    if input.ClientCompanyID == 0 {
        return errors.New("client company is required")
    }
    
    // Contact information validation
    if err := validateContactInfo(input.ContactName, input.ContactEmail, input.ContactPhone); err != nil {
        return err
    }
    
    // Pickup address validation
    if err := validateAddress(input.PickupAddress, input.PickupCity, input.PickupState, input.PickupZip); err != nil {
        return err
    }
    
    // Pickup date and time validation
    if err := validatePickupDateTime(input.PickupDate, input.PickupTimeSlot); err != nil {
        return err
    }
    
    // JIRA ticket validation
    if err := validateJiraTicket(input.JiraTicketNumber); err != nil {
        return err
    }
    
    // Laptop count validation (must be >= 2 for bulk)
    if input.NumberOfLaptops < 2 {
        return errors.New("bulk shipments must have at least 2 laptops")
    }
    
    // Bulk dimensions validation (REQUIRED for bulk shipments)
    if input.BulkLength <= 0 || input.BulkWidth <= 0 || input.BulkHeight <= 0 || input.BulkWeight <= 0 {
        return errors.New("bulk dimensions (length, width, height, weight) are required and must be positive")
    }
    
    // Accessories validation
    if input.IncludeAccessories && strings.TrimSpace(input.AccessoriesDescription) == "" {
        return errors.New("accessories description is required when including accessories")
    }
    
    return nil
}
```

**Commit:** `feat: add bulk to warehouse form validator`

---

### 3.3 Create Warehouse to Engineer Form Validator

#### 🟥 RED: Test warehouse to engineer form validation
**File:** `internal/validator/warehouse_to_engineer_form_test.go` (NEW FILE)
```go
package validator

import (
    "testing"
)

func TestValidateWarehouseToEngineerForm(t *testing.T) {
    tests := []struct {
        name          string
        input         WarehouseToEngineerFormInput
        shouldBeValid bool
        errorContains string
    }{
        {
            name: "valid warehouse to engineer form",
            input: WarehouseToEngineerFormInput{
                LaptopID:            1,
                SoftwareEngineerID:  5,
                EngineerName:        "Jane Smith",
                EngineerEmail:       "jane@bairesdev.com",
                EngineerAddress:     "456 Tech Ave",
                EngineerCity:        "San Francisco",
                EngineerState:       "CA",
                EngineerZip:         "94102",
                CourierName:         "FedEx",
                TrackingNumber:      "123456789",
                JiraTicketNumber:    "SCOP-12345",
            },
            shouldBeValid: true,
        },
        {
            name: "invalid - missing laptop ID",
            input: WarehouseToEngineerFormInput{
                LaptopID:            0, // Missing
                SoftwareEngineerID:  5,
                EngineerName:        "Jane Smith",
                EngineerEmail:       "jane@bairesdev.com",
                EngineerAddress:     "456 Tech Ave",
                EngineerCity:        "San Francisco",
                EngineerState:       "CA",
                EngineerZip:         "94102",
                CourierName:         "FedEx",
                TrackingNumber:      "123456789",
                JiraTicketNumber:    "SCOP-12345",
            },
            shouldBeValid: false,
            errorContains: "laptop selection is required",
        },
        {
            name: "invalid - missing engineer",
            input: WarehouseToEngineerFormInput{
                LaptopID:         1,
                EngineerAddress:  "456 Tech Ave",
                EngineerCity:     "San Francisco",
                EngineerState:    "CA",
                EngineerZip:      "94102",
                CourierName:      "FedEx",
                TrackingNumber:   "123456789",
                JiraTicketNumber: "SCOP-12345",
            },
            shouldBeValid: false,
            errorContains: "software engineer is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateWarehouseToEngineerForm(tt.input)
            if tt.shouldBeValid && err != nil {
                t.Errorf("Expected valid, got error: %v", err)
            }
            if !tt.shouldBeValid {
                if err == nil {
                    t.Error("Expected error, got nil")
                } else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
                    t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
                }
            }
        })
    }
}
```

#### 🟩 GREEN: Implement warehouse to engineer form validator
**File:** `internal/validator/warehouse_to_engineer_form.go` (NEW FILE)
```go
package validator

import (
    "errors"
    "strings"
)

// WarehouseToEngineerFormInput represents the form data for warehouse-to-engineer shipments
type WarehouseToEngineerFormInput struct {
    LaptopID            int64
    SoftwareEngineerID  int64
    EngineerName        string
    EngineerEmail       string
    EngineerAddress     string
    EngineerCity        string
    EngineerState       string
    EngineerZip         string
    CourierName         string
    TrackingNumber      string
    JiraTicketNumber    string
    SpecialInstructions string
}

// ValidateWarehouseToEngineerForm validates the warehouse-to-engineer shipment form
func ValidateWarehouseToEngineerForm(input WarehouseToEngineerFormInput) error {
    // Laptop selection validation (REQUIRED)
    if input.LaptopID == 0 {
        return errors.New("laptop selection is required")
    }
    
    // Software engineer validation (REQUIRED)
    if input.SoftwareEngineerID == 0 && strings.TrimSpace(input.EngineerName) == "" {
        return errors.New("software engineer is required")
    }
    
    // Engineer address validation
    if err := validateAddress(input.EngineerAddress, input.EngineerCity, input.EngineerState, input.EngineerZip); err != nil {
        return err
    }
    
    // Courier information validation (optional initially, required before shipping)
    // For form submission, we'll make it optional
    
    // JIRA ticket validation
    if err := validateJiraTicket(input.JiraTicketNumber); err != nil {
        return err
    }
    
    return nil
}
```

**Commit:** `feat: add warehouse to engineer form validator`

---

## Phase 4: Handler Layer (Days 9-12)

### 4.1 Update Pickup Form Handler for Single Full Journey

#### 🟥 RED: Test single full journey form submission
**File:** `internal/handlers/pickup_form_test.go`
```go
func TestPickupFormHandler_SubmitSingleFullJourney(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Create test data
    companyID := createTestClientCompany(t, db, "Test Company")
    userID := createTestUser(t, db, "client@test.com", models.RoleClient)
    
    // Create handler
    handler := NewPickupFormHandler(db, nil, nil)
    
    // Create form data
    formData := url.Values{
        "shipment_type":       {string(models.ShipmentTypeSingleFullJourney)},
        "client_company_id":   {strconv.FormatInt(companyID, 10)},
        "contact_name":        {"John Doe"},
        "contact_email":       {"john@test.com"},
        "contact_phone":       {"+1-555-0123"},
        "pickup_address":      {"123 Main St"},
        "pickup_city":         {"New York"},
        "pickup_state":        {"NY"},
        "pickup_zip":          {"10001"},
        "pickup_date":         {"2025-11-15"},
        "pickup_time_slot":    {"morning"},
        "jira_ticket_number":  {"SCOP-12345"},
        "laptop_serial_number": {"ABC123456"},
        "laptop_specs":        {"Dell XPS 15, 16GB RAM"},
        "engineer_name":       {"Jane Smith"},
    }
    
    // Create request
    req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Role: models.RoleClient}))
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call handler
    handler.PickupFormSubmit(w, req)
    
    // Check response
    if w.Code != http.StatusSeeOther {
        t.Errorf("Expected status 303, got %d", w.Code)
    }
    
    // Verify shipment was created
    var shipmentID int64
    var shipmentType models.ShipmentType
    var laptopCount int
    err := db.QueryRow(`
        SELECT id, shipment_type, laptop_count 
        FROM shipments 
        WHERE client_company_id = $1 AND jira_ticket_number = $2
    `, companyID, "SCOP-12345").Scan(&shipmentID, &shipmentType, &laptopCount)
    
    if err != nil {
        t.Fatalf("Shipment not created: %v", err)
    }
    
    if shipmentType != models.ShipmentTypeSingleFullJourney {
        t.Errorf("Expected shipment type %s, got %s", models.ShipmentTypeSingleFullJourney, shipmentType)
    }
    
    if laptopCount != 1 {
        t.Errorf("Expected laptop count 1, got %d", laptopCount)
    }
    
    // Verify laptop was created
    var laptopSerialNumber string
    err = db.QueryRow(`
        SELECT l.serial_number 
        FROM laptops l
        JOIN shipment_laptops sl ON sl.laptop_id = l.id
        WHERE sl.shipment_id = $1
    `, shipmentID).Scan(&laptopSerialNumber)
    
    if err != nil {
        t.Fatalf("Laptop not created: %v", err)
    }
    
    if laptopSerialNumber != "ABC123456" {
        t.Errorf("Expected serial number ABC123456, got %s", laptopSerialNumber)
    }
}
```

#### 🟩 GREEN: Update pickup form handler
**File:** `internal/handlers/pickup_form.go`

Update the `PickupFormSubmit` method to handle shipment types. This will be a significant update to the existing handler.

Key changes:
1. Add `shipment_type` form field handling
2. Branch logic based on shipment type
3. For `single_full_journey`: create laptop record automatically
4. Update validation to use new type-specific validators

**Commit:** `feat: update pickup form handler to support single full journey shipments`

---

### 4.2 Create Bulk to Warehouse Form Handler

#### 🟥 RED: Test bulk to warehouse form submission
**File:** `internal/handlers/bulk_shipment_form_test.go` (NEW FILE)
```go
func TestBulkShipmentFormHandler_Submit(t *testing.T) {
    // Similar structure to single shipment test
    // Key differences:
    // - NumberOfLaptops >= 2
    // - Bulk dimensions required
    // - No laptop records created initially
    // - Verify laptop_count is set correctly
}
```

#### 🟩 GREEN: Create bulk shipment form handler
**File:** `internal/handlers/bulk_shipment_form.go` (NEW FILE OR extend pickup_form.go)

**Commit:** `feat: add bulk to warehouse shipment form handler`

---

### 4.3 Create Warehouse to Engineer Form Handler

#### 🟥 RED: Test warehouse to engineer form submission
**File:** `internal/handlers/warehouse_to_engineer_form_test.go` (NEW FILE)
```go
func TestWarehouseToEngineerFormHandler_Submit(t *testing.T) {
    // Setup: Create laptop with status 'available' and reception report
    // Verify only available laptops are shown in dropdown
    // Verify shipment is created with correct type
    // Verify laptop status is updated
    // Verify shipment starts at 'released_from_warehouse' status
}
```

#### 🟩 GREEN: Create warehouse to engineer form handler
**File:** `internal/handlers/warehouse_to_engineer_form.go` (NEW FILE)

**Commit:** `feat: add warehouse to engineer shipment form handler`

---

### 4.4 Update Shipments List Handler

#### 🟥 RED: Test shipment list filtering by type
**File:** `internal/handlers/shipments_test.go`
```go
func TestShipmentsHandler_ListWithTypeFilter(t *testing.T) {
    // Create shipments of different types
    // Test filtering by shipment_type
    // Verify correct shipments are returned
}
```

#### 🟩 GREEN: Add type filter to shipments list
**File:** `internal/handlers/shipments.go`

Update `ShipmentsList` method to:
1. Accept `type` query parameter
2. Filter shipments by type
3. Display shipment type in list view

**Commit:** `feat: add shipment type filtering to shipments list`

---

### 4.5 Update Shipment Detail Handler

#### 🟥 RED: Test shipment detail displays type information
**File:** `internal/handlers/shipments_test.go`
```go
func TestShipmentsHandler_DetailShowsType(t *testing.T) {
    // Create shipment with specific type
    // Verify type is displayed
    // Verify type-specific information is shown
}
```

#### 🟩 GREEN: Update shipment detail view
**File:** `internal/handlers/shipments.go`

Update `ShipmentDetail` method to include shipment type in response data.

**Commit:** `feat: add shipment type display to shipment detail view`

---

## Phase 5: Templates & UI (Days 13-15)

### 5.1 Update/Create Single Full Journey Form Template

#### 🟥 RED: Test single form template renders correctly
Create basic test to verify template renders without errors.

#### 🟩 GREEN: Create/update template
**File:** `templates/pages/single-shipment-form.html` (NEW) or update existing

Based on current `shipment-pickup-form.html`, but:
- Remove bulk toggle
- Remove "number of boxes" field
- Always set laptop count to 1
- Add laptop details section:
  - Serial number (text input, required)
  - Specifications (textarea, optional)
  - Engineer name (text input, optional)
- Keep accessories section

**Commit:** `feat: create single full journey shipment form template`

---

### 5.2 Create Bulk to Warehouse Form Template

#### 🟥 RED: Test bulk form template renders correctly

#### 🟩 GREEN: Create template
**File:** `templates/pages/bulk-shipment-form.html` (NEW)

Similar to current form, but:
- Bulk dimensions MANDATORY (not toggled)
- Laptop count >= 2 (required)
- No engineer assignment section
- Clear indication this is bulk-only

**Commit:** `feat: create bulk to warehouse shipment form template`

---

### 5.3 Create Warehouse to Engineer Form Template

#### 🟥 RED: Test warehouse-to-engineer form template renders correctly

#### 🟩 GREEN: Create template
**File:** `templates/pages/warehouse-to-engineer-form.html` (NEW)

Key elements:
- Laptop selection dropdown (populated from available inventory)
- Display laptop details (read-only: serial number, specs, client company)
- Engineer selection/creation section
- Delivery address
- Courier information
- Tracking number

**Commit:** `feat: create warehouse to engineer shipment form template`

---

### 5.4 Update Dashboard with Three Create Buttons

#### 🟥 RED: Test dashboard displays three create buttons

#### 🟩 GREEN: Update dashboard template
**File:** `templates/pages/dashboard.html`

Add three prominent buttons in the header or actions section:
```html
<div class="flex space-x-4">
    <a href="/shipments/create/single" class="btn btn-primary">
        + Single Shipment
    </a>
    <a href="/shipments/create/bulk" class="btn btn-secondary">
        + Bulk to Warehouse
    </a>
    <a href="/shipments/create/warehouse-to-engineer" class="btn btn-secondary">
        + Warehouse to Engineer
    </a>
</div>
```

**Commit:** `feat: add three shipment type creation buttons to dashboard`

---

### 5.5 Update Shipments List Page

#### 🟥 RED: Test shipments list shows type indicators

#### 🟩 GREEN: Update shipments list template
**File:** `templates/pages/shipments.html`

Updates:
1. Add three create buttons (same as dashboard)
2. Add shipment type column/badge
3. Add type filter dropdown
4. Display type-specific information (e.g., laptop count for bulk)

**Commit:** `feat: update shipments list with type indicators and filters`

---

### 5.6 Update Shipment Detail Page

#### 🟥 RED: Test shipment detail shows type-specific information

#### 🟩 GREEN: Update shipment detail template
**File:** `templates/pages/shipment-detail.html`

Updates:
1. Display shipment type prominently
2. Show type-specific status flow (only relevant statuses)
3. Show laptop details for single shipments
4. Show laptop count for bulk shipments
5. Show available actions based on type

**Commit:** `feat: update shipment detail page with type-specific information`

---

## Phase 6: Integration & Testing (Days 16-18)

### 6.1 Integration Tests for Each Shipment Type

#### 🟥 RED: Write end-to-end test for single full journey
**File:** `tests/integration/single_shipment_flow_test.go` (NEW)
```go
func TestSingleFullJourneyFlow(t *testing.T) {
    // 1. Create shipment via form submission
    // 2. Verify laptop is auto-created
    // 3. Progress through each status
    // 4. Verify laptop status syncs
    // 5. Create reception report (verify serial number)
    // 6. Assign engineer (if not already assigned)
    // 7. Release from warehouse
    // 8. Mark in transit to engineer
    // 9. Create delivery form
    // 10. Verify final status is delivered
    // 11. Verify laptop status is delivered
}
```

#### 🟩 GREEN: Implement integration tests for all three types

**Commit:** `test: add integration tests for all three shipment types`

---

### 6.2 Test Status Transition Restrictions

#### 🟥 RED: Test type-specific status restrictions

**File:** `tests/integration/status_transition_test.go` (NEW)
```go
func TestBulkShipmentCannotProgressPastWarehouse(t *testing.T) {
    // Create bulk shipment
    // Progress to at_warehouse
    // Attempt to update to released_from_warehouse
    // Verify error/rejection
}

func TestWarehouseToEngineerStartsAtReleased(t *testing.T) {
    // Create warehouse-to-engineer shipment
    // Verify initial status is released_from_warehouse
    // Verify cannot set to earlier statuses
}
```

#### 🟩 GREEN: Ensure status restrictions work correctly

**Commit:** `test: verify type-specific status transition restrictions`

---

### 6.3 Test Laptop Status Synchronization

#### 🟥 RED: Test laptop sync for single shipments only

**File:** `tests/integration/laptop_sync_test.go` (NEW)
```go
func TestLaptopStatusSyncsForSingleShipments(t *testing.T) {
    // Create single_full_journey shipment
    // Update shipment status
    // Verify laptop status updates automatically
}

func TestLaptopStatusDoesNotSyncForBulkShipments(t *testing.T) {
    // Create bulk shipment (after reception with laptops)
    // Update shipment status
    // Verify laptop statuses remain unchanged
}
```

#### 🟩 GREEN: Implement laptop sync tests

**Commit:** `test: verify laptop status synchronization logic`

---

### 6.4 Test Inventory Availability for Warehouse-to-Engineer

#### 🟥 RED: Test only available laptops with reception reports are selectable

**File:** `tests/integration/inventory_availability_test.go` (NEW)
```go
func TestWarehouseToEngineerOnlyShowsAvailableLaptops(t *testing.T) {
    // Create laptop without reception report
    // Create laptop with reception report but in active shipment
    // Create laptop with reception report and available
    // Query available laptops
    // Verify only the third laptop is returned
}
```

#### 🟩 GREEN: Implement inventory availability tests

**Commit:** `test: verify inventory availability queries for warehouse-to-engineer`

---

### 6.5 Test Serial Number Correction Workflow

#### 🟥 RED: Test serial number correction for single shipments

**File:** `tests/integration/serial_correction_test.go` (NEW)
```go
func TestSerialNumberCorrectionWorkflow(t *testing.T) {
    // Create single shipment with serial "ABC123"
    // Warehouse receives with serial "XYZ789"
    // Verify correction is flagged
    // Logistics user approves correction
    // Verify laptop serial is updated
    // Verify correction note is saved
}

func TestOnlyLogisticsCanApproveCorrections(t *testing.T) {
    // Create correction scenario
    // Attempt approval as warehouse user (should fail)
    // Attempt approval as logistics user (should succeed)
}
```

#### 🟩 GREEN: Implement serial correction tests

**Commit:** `test: verify serial number correction workflow with role restrictions`

---

## Phase 7: Documentation & Cleanup (Day 18)

### 7.1 Update README

**File:** `readme.md`

Updates:
1. Document three shipment types
2. Update process flow section
3. Add examples for each type
4. Update features list

**Commit:** `docs: update README with three shipment types documentation`

---

### 7.2 Update Plan.md

**File:** `docs/plan.md`

Mark Phase 7 (or new phase) as complete with shipment types implementation.

**Commit:** `docs: mark shipment types implementation as complete`

---

### 7.3 Create User Guide

**File:** `docs/SHIPMENT_TYPES_GUIDE.md` (NEW)

Comprehensive guide explaining:
1. When to use each shipment type
2. How to create each type
3. Status flows for each type
4. Special considerations (engineer assignment, laptop sync, etc.)

**Commit:** `docs: add comprehensive shipment types user guide`

---

## Summary

**Total Phases:** 7  
**Total Days:** 16-18  
**Total Files:** ~30 new/modified  
**Total Migrations:** 3 new  
**Total Tests:** 40+ new test cases

### Key Deliverables

1. ✅ Three distinct shipment types with full validation
2. ✅ Type-specific status flows
3. ✅ Automated laptop record creation for single shipments
4. ✅ Serial number correction tracking
5. ✅ Inventory availability for warehouse-to-engineer
6. ✅ Three separate form templates
7. ✅ Updated dashboard and list views with three create buttons
8. ✅ Comprehensive test coverage
9. ✅ Full documentation

### Risk Mitigation

- **Database migrations:** Test thoroughly in development before production
- **Backward compatibility:** Migration sets existing shipments to `single_full_journey`
- **User experience:** Clear labeling and separate buttons prevent confusion
- **Data integrity:** Type-specific validation prevents invalid state transitions

---

## Approval Checklist

Before proceeding, please confirm:

- [ ] Scope and phases are acceptable
- [ ] Timeline (16-18 days) is acceptable
- [ ] TDD approach with RED/GREEN/REFACTOR cycles is understood
- [ ] All design decisions are approved
- [ ] Ready to begin Phase 1: Database Schema Changes

**Awaiting your approval to proceed with implementation.**

