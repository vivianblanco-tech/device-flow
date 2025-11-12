# Continuation Prompt: Phase 5 - Templates & UI

**Date:** November 12, 2025  
**Current Phase:** Phase 5 (Templates & UI)  
**Previous Phase:** Phase 4 (Handler Layer) - COMPLETE ✅  
**Status:** Ready to begin Phase 5

---

## Context

You are continuing the **Three Shipment Types** implementation for the laptop tracking system. Phase 4 (Handler Layer) has been completed and verified with all 72 handler tests passing.

**Previous Phases Completed:**
- ✅ **Phase 1:** Database Schema Changes (migrations for shipment_type and laptop_count)
- ✅ **Phase 2:** Model Layer Updates (ShipmentType enum, validation methods, status flows)
- ✅ **Phase 3:** Validator Updates (three type-specific validators)
- ✅ **Phase 4:** Handler Layer Updates (three form handlers, list/detail enhancements)

---

## Phase 5 Overview

Phase 5 focuses on creating and updating UI templates to support the three distinct shipment types:

1. **`single_full_journey`** - Client → Warehouse → Engineer (8 statuses)
2. **`bulk_to_warehouse`** - Client → Warehouse only (5 statuses)
3. **`warehouse_to_engineer`** - Warehouse → Engineer only (3 statuses)

---

## Phase 5 Tasks (from THREE_SHIPMENT_TYPES_TDD_PLAN.md)

### 5.1 Single Full Journey Form Template ⏳
**Goal:** Create/update form template for single shipment type

**Requirements:**
- Remove bulk toggle (not applicable for single shipments)
- Add laptop details section:
  - Serial number (text input, **REQUIRED**)
  - Specifications (textarea, optional)
  - Engineer name (text input, **OPTIONAL** - can be assigned later)
- Keep accessories section
- Always set `laptop_count = 1`
- Set `shipment_type = "single_full_journey"`

**Template File:** `templates/pages/single-shipment-form.html` (NEW) or update existing

**Handler Integration:** Already implemented (`handleSingleFullJourneyForm()`)

---

### 5.2 Bulk to Warehouse Form Template ⏳
**Goal:** Create form template for bulk shipments to warehouse

**Requirements:**
- Bulk dimensions **MANDATORY** (not toggled):
  - Length (cm)
  - Width (cm)
  - Height (cm)
  - Weight (kg)
- Laptop count >= 2 (required, validated)
- NO engineer assignment section
- NO serial number input (laptops created during warehouse reception)
- Set `shipment_type = "bulk_to_warehouse"`
- Clear indication this is bulk-only

**Template File:** `templates/pages/bulk-shipment-form.html` (NEW)

**Handler Integration:** Already implemented (`handleBulkToWarehouseForm()`)

---

### 5.3 Warehouse to Engineer Form Template ⏳
**Goal:** Create form template for warehouse-to-engineer shipments

**Requirements:**
- Laptop selection dropdown (populated from available inventory)
- Display selected laptop details (read-only):
  - Serial number
  - Specifications
  - Client company
- Engineer selection/creation section (**REQUIRED**)
- Delivery address section
- Courier information (optional)
- Tracking number (optional)
- Set `shipment_type = "warehouse_to_engineer"`
- NO pickup date/location (ships from warehouse)

**Template File:** `templates/pages/warehouse-to-engineer-form.html` (NEW)

**Handler Integration:** Already implemented (`handleWarehouseToEngineerForm()`)

**Note:** Requires backend endpoint to fetch available laptops (AJAX/API)

---

### 5.4 Update Dashboard with Three Create Buttons ⏳
**Goal:** Provide clear navigation to create each shipment type

**Requirements:**
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

**Template File:** `templates/pages/dashboard.html`

**Routes Required:**
- `/shipments/create/single` → Single Full Journey Form
- `/shipments/create/bulk` → Bulk to Warehouse Form
- `/shipments/create/warehouse-to-engineer` → Warehouse to Engineer Form

---

### 5.5 Update Shipments List Page ⏳
**Goal:** Display shipment type information and enable type filtering

**Requirements:**
1. Add three create buttons (same as dashboard)
2. Add shipment type column/badge with visual indicators:
   - 🔵 Single Full Journey
   - 📦 Bulk to Warehouse
   - ⚡ Warehouse → Engineer
3. Add type filter dropdown:
   ```html
   <select name="type">
       <option value="">All Types</option>
       <option value="single_full_journey">Single Full Journey</option>
       <option value="bulk_to_warehouse">Bulk to Warehouse</option>
       <option value="warehouse_to_engineer">Warehouse to Engineer</option>
   </select>
   ```
4. Display type-specific information:
   - Show laptop count for bulk shipments
   - Show engineer info for single/warehouse-to-engineer

**Template File:** `templates/pages/shipments-list.html`

**Handler Support:** Already implemented (Phase 4.4) - handler passes `TypeFilter` and `AllShipmentTypes`

---

### 5.6 Update Shipment Detail Page ⏳
**Goal:** Display type-specific information prominently

**Requirements:**
1. Display shipment type as a prominent badge at top
2. Show type-specific status flow:
   - Single: All 8 statuses
   - Bulk: Only first 5 statuses
   - Warehouse→Engineer: Only last 3 statuses
3. Display laptop details section (for single shipments):
   - Serial number
   - Specifications
   - Engineer assignment (with edit capability)
4. Display laptop count (for bulk shipments):
   - "Contains X laptops"
   - Link to view individual laptops after reception
5. Display type-appropriate action buttons

**Template File:** `templates/pages/shipment-detail.html`

**Handler Support:** Already implemented (Phase 4.5) - handler passes shipment type data

---

## TDD Approach for Phase 5

While templates are primarily UI, we still follow TDD principles:

### 🟥 RED: Write Test
For each template update, write a test that:
1. Verifies template renders without errors
2. Checks for presence of required form fields
3. Validates correct data is passed to template

**Example Test Structure:**
```go
func TestSingleShipmentFormTemplate(t *testing.T) {
    // Setup handler with templates
    // Create request with test data
    // Render template
    // Verify response contains required fields:
    //   - input[name="laptop_serial_number"] (required)
    //   - textarea[name="laptop_specs"]
    //   - input[name="engineer_name"]
    //   - NO bulk dimensions fields
}
```

### 🟩 GREEN: Implement Template
Create/update the template to pass the test:
1. Add required form fields with correct names
2. Add proper validation attributes
3. Style appropriately
4. Include error message display

### ✅ Commit
After test passes, commit with descriptive message:
```
feat: add single full journey shipment form template
```

---

## Current Template Structure

**Existing Templates:**
```
templates/
├── layouts/
│   └── base.html
├── components/
│   └── navigation.html
└── pages/
    ├── dashboard.html
    ├── shipments-list.html
    ├── shipment-detail.html
    ├── pickup-form.html (legacy - may need updating)
    ├── reception-report.html
    └── delivery-form.html
```

**CSS Framework:** TailwindCSS (via `tailwindcss.exe`)

---

## Handler Endpoints Already Available

The following handler endpoints are ready to receive form submissions:

### Pickup Form Handler (`internal/handlers/pickup_form.go`)
- **Route:** `POST /pickup-form`
- **Methods:**
  - `handleSingleFullJourneyForm()` - Lines 244-400
  - `handleBulkToWarehouseForm()` - Lines 554-698
  - `handleWarehouseToEngineerForm()` - Lines 701-890

**Form Field:** `shipment_type` (determines which handler is called)

### Shipments List (`internal/handlers/shipments.go`)
- **Route:** `GET /shipments`
- **Query Parameters:**
  - `status` - Filter by status
  - `type` - Filter by shipment type
  - `search` - Search by tracking number or company name

### Shipment Detail (`internal/handlers/shipments.go`)
- **Route:** `GET /shipments/{id}`
- **Template Data:**
  - `Shipment` - Full shipment object with `ShipmentType` and `LaptopCount`
  - `Laptops` - Array of laptops linked to shipment
  - `AllShipmentTypes` - For filtering

---

## Routes to Add

You'll need to add routes for the three create form pages:

```go
// In cmd/web/main.go or routes setup
r.HandleFunc("/shipments/create/single", pickupFormHandler.SingleShipmentFormPage).Methods("GET")
r.HandleFunc("/shipments/create/bulk", pickupFormHandler.BulkShipmentFormPage).Methods("GET")
r.HandleFunc("/shipments/create/warehouse-to-engineer", pickupFormHandler.WarehouseToEngineerFormPage).Methods("GET")
```

Or reuse existing `/pickup-form` route with query parameter:
```
GET /pickup-form?type=single_full_journey
GET /pickup-form?type=bulk_to_warehouse
GET /pickup-form?type=warehouse_to_engineer
```

---

## Implementation Strategy

### Option 1: Separate Template Files (Recommended)
**Pros:**
- Clear separation of concerns
- Easier to maintain
- Type-specific logic contained
- Better for future enhancements

**Cons:**
- More files to manage
- Some code duplication for common sections

**Files:**
- `single-shipment-form.html`
- `bulk-shipment-form.html`
- `warehouse-to-engineer-form.html`

### Option 2: Single Template with Conditionals
**Pros:**
- Single file to maintain
- Shared code reused

**Cons:**
- More complex template logic
- Harder to read and debug
- Mixing concerns

**Recommendation:** Use **Option 1** for clarity and maintainability.

---

## Testing Approach

### Manual Testing Checklist (After Each Template)

1. **Form Renders Correctly**
   - [ ] All required fields present
   - [ ] Proper field types (text, textarea, select, number)
   - [ ] Validation attributes correct (required, min, max, pattern)

2. **Form Submission Works**
   - [ ] Valid data creates shipment
   - [ ] Invalid data shows error messages
   - [ ] Redirect to shipment detail after success

3. **Visual Design**
   - [ ] Consistent with existing pages
   - [ ] Mobile responsive
   - [ ] Clear labels and help text

4. **Browser Testing**
   - [ ] Chrome
   - [ ] Firefox
   - [ ] Edge

---

## Expected Deliverables for Phase 5

1. **Three New Form Templates:**
   - ✅ `templates/pages/single-shipment-form.html`
   - ✅ `templates/pages/bulk-shipment-form.html`
   - ✅ `templates/pages/warehouse-to-engineer-form.html`

2. **Updated Dashboard:**
   - ✅ `templates/pages/dashboard.html` (with three buttons)

3. **Updated Shipments List:**
   - ✅ `templates/pages/shipments-list.html` (with type badges and filter)

4. **Updated Shipment Detail:**
   - ✅ `templates/pages/shipment-detail.html` (with type-specific display)

5. **Updated Routes/Handlers:**
   - ✅ Form page handlers (GET methods)
   - ✅ Route registration in main.go

6. **Tests:**
   - ✅ Template rendering tests for each new form
   - ✅ Integration tests for form submission flows

7. **Documentation:**
   - ✅ Update README with new UI features
   - ✅ Create Phase 5 completion summary

---

## Questions to Address Before Starting

1. **Route Strategy:** Separate routes or query parameters?
2. **Available Laptops Endpoint:** Create REST API or include in page data?
3. **Engineer Selection:** Dropdown of existing engineers or free text input?
4. **Visual Design:** Icons/colors for each shipment type?

---

## Ready to Begin Phase 5

**Current Status:**
- ✅ Phase 1: Database Schema Changes - COMPLETE
- ✅ Phase 2: Model Layer Updates - COMPLETE
- ✅ Phase 3: Validator Updates - COMPLETE
- ✅ Phase 4: Handler Layer Updates - COMPLETE
- ⏳ Phase 5: Templates & UI - **READY TO START**

**Next Task:** Begin Phase 5.1 - Create Single Full Journey Form Template

---

## Verification Command

After implementing templates, verify with:

```powershell
# Run all tests
go test ./... -v

# Start the application
docker-compose up -d

# Access forms:
# http://localhost:8080/shipments/create/single
# http://localhost:8080/shipments/create/bulk
# http://localhost:8080/shipments/create/warehouse-to-engineer
```

---

**Reference Documents:**
- `docs/PHASE4_COMPLETE.md` - Phase 4 summary
- `docs/PHASE4_VERIFICATION.md` - Complete verification results
- `docs/THREE_SHIPMENT_TYPES_TDD_PLAN.md` - Full implementation plan
- `docs/tdd.md` - TDD methodology guidelines

**Ready to proceed!** 🚀

