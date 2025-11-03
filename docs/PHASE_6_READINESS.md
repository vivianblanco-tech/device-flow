# Phase 6 Readiness Checklist

## Prerequisites Check

Before starting Phase 6 (Dashboard & Visualization), verify all prerequisites are met.

### ✅ Completed Requirements

- [x] **Phase 0**: Project Setup & Infrastructure
- [x] **Phase 1**: Database Schema & Core Models (97.7% coverage)
- [x] **Phase 2**: Authentication System (RBAC, OAuth, Magic Links)
- [x] **Phase 3**: Core Forms & Workflows (Pickup, Reception, Delivery)
- [x] **Phase 4**: JIRA Integration (Create, Update, Sync)
- [x] **Phase 5**: Email Notifications (6 templates)

### ⚠️ Optional Improvements (Recommended)

- [ ] **Test Database Setup**: Configure `laptop_tracking_test` database
  - See: `docs/TEST_DATABASE_SETUP.md`
  - Impact: Enables 40 integration tests to run
  - Time: 10 minutes

- [ ] **Template Layouts**: Create shared layout templates
  - Current: Each page has embedded HTML structure
  - Recommended: Extract to `templates/layouts/base.html`
  - Impact: Reduces code duplication
  - Time: 30 minutes

- [ ] **Error Handling**: Add global error handler
  - Current: Each handler manages errors independently
  - Recommended: Create middleware for consistent error pages
  - Impact: Better user experience
  - Time: 1 hour

## Phase 6 Overview

**Goal**: Create dashboard with statistics, charts, calendar, and inventory management.

**Estimated Timeline**: 2-3 days

**Key Deliverables**:
1. Dashboard with key metrics
2. Data visualization charts
3. Calendar view for schedules
4. Inventory management UI

## Technical Readiness

### Database Queries Required

The following queries will be needed for Phase 6:

#### Dashboard Statistics

```sql
-- Total shipments by status
SELECT status, COUNT(*) FROM shipments GROUP BY status;

-- Average delivery time
SELECT AVG(delivered_at - created_at) FROM shipments WHERE status = 'delivered';

-- Shipments created in last 30 days
SELECT COUNT(*) FROM shipments WHERE created_at > NOW() - INTERVAL '30 days';

-- Inventory counts
SELECT status, COUNT(*) FROM laptops GROUP BY status;
```

#### Calendar Events

```sql
-- Upcoming pickups
SELECT id, pickup_date, client_company_id FROM shipments 
WHERE status IN ('pending_pickup', 'picked_up_from_client') 
AND pickup_date >= CURRENT_DATE;

-- Upcoming deliveries
SELECT id, delivery_date, software_engineer_id FROM shipments 
WHERE status IN ('released_from_warehouse', 'in_transit_to_engineer')
AND delivery_date >= CURRENT_DATE;
```

#### Inventory Management

```sql
-- Available laptops
SELECT * FROM laptops WHERE status = 'available' ORDER BY created_at DESC;

-- Laptops in use
SELECT l.*, s.id as shipment_id FROM laptops l
JOIN shipment_laptops sl ON l.id = sl.laptop_id
JOIN shipments s ON sl.shipment_id = s.id
WHERE s.status NOT IN ('delivered');
```

### Required Go Packages

Additional dependencies for Phase 6:

```go
// For date/time handling
import "time"

// For JSON API responses
import "encoding/json"

// For aggregation queries (already imported)
import "database/sql"

// No new external dependencies required!
```

### Frontend Libraries Needed

For charts and calendar:

**Option 1: Chart.js (Recommended)**
- Lightweight (~60KB)
- Easy to use
- Good documentation
- MIT License

```html
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
```

**Option 2: FullCalendar (Recommended for Calendar)**
- Feature-rich
- Mobile responsive
- MIT License

```html
<link href='https://cdn.jsdelivr.net/npm/fullcalendar@6.1.9/index.global.min.css' rel='stylesheet' />
<script src='https://cdn.jsdelivr.net/npm/fullcalendar@6.1.9/index.global.min.js'></script>
```

**Option 3: Pure CSS Calendar (Lightweight Alternative)**
- No external dependencies
- Custom implementation with Go templates
- Full control

## Phase 6 Task Breakdown

### 6.1 Dashboard Statistics (Day 1)

**Files to Create**:
- `internal/handlers/dashboard.go` - HTTP handler
- `internal/handlers/dashboard_test.go` - Tests
- `internal/models/statistics.go` - Statistics queries
- `internal/models/statistics_test.go` - Tests
- `templates/pages/dashboard.html` - UI

**TDD Approach**:
1. 🟥 RED: Test for fetching total shipments by status
2. 🟩 GREEN: Implement query
3. 🟥 RED: Test for average delivery time calculation
4. 🟩 GREEN: Implement calculation
5. 🟥 RED: Test for inventory counts
6. 🟩 GREEN: Implement query
7. 🟥 RED: Test for dashboard handler
8. 🟩 GREEN: Implement handler and template

**Estimated Tests**: 8-10 test cases

### 6.2 Data Visualization (Day 1-2)

**Files to Create**:
- `internal/handlers/api_stats.go` - JSON API endpoints
- `internal/handlers/api_stats_test.go` - Tests
- `static/js/charts.js` - Chart rendering logic
- Update `templates/pages/dashboard.html` - Add chart containers

**TDD Approach**:
1. 🟥 RED: Test for chart data formatting (shipments over time)
2. 🟩 GREEN: Implement data aggregation
3. 🟥 RED: Test for status breakdown data
4. 🟩 GREEN: Implement status counts
5. 🟥 RED: Test for API endpoint response format
6. 🟩 GREEN: Implement JSON API handler

**Charts to Implement**:
- Line chart: Shipments over time (last 30 days)
- Donut chart: Status breakdown
- Bar chart: Average delivery time by month

**Estimated Tests**: 6-8 test cases

### 6.3 Calendar View (Day 2)

**Files to Create**:
- `internal/handlers/calendar.go` - HTTP handler
- `internal/handlers/calendar_test.go` - Tests
- `internal/models/calendar.go` - Calendar event queries
- `internal/models/calendar_test.go` - Tests
- `templates/pages/calendar.html` - UI
- `static/js/calendar.js` - Calendar logic

**TDD Approach**:
1. 🟥 RED: Test for fetching upcoming pickups
2. 🟩 GREEN: Implement pickup events query
3. 🟥 RED: Test for fetching upcoming deliveries
4. 🟩 GREEN: Implement delivery events query
5. 🟥 RED: Test for event date filtering
6. 🟩 GREEN: Implement date range filter
7. 🟥 RED: Test for calendar handler
8. 🟩 GREEN: Implement handler and template

**Event Types**:
- 🔵 Pickup scheduled
- 🟡 In transit to warehouse
- 🟢 In transit to engineer
- 🟣 Delivery scheduled

**Estimated Tests**: 8-10 test cases

### 6.4 Inventory Management (Day 2-3)

**Files to Create**:
- `internal/handlers/inventory.go` - CRUD handlers
- `internal/handlers/inventory_test.go` - Tests
- `internal/validator/laptop.go` - Laptop validation
- `internal/validator/laptop_test.go` - Tests
- `templates/pages/inventory-list.html` - List view
- `templates/pages/inventory-form.html` - Add/Edit form

**TDD Approach**:
1. 🟥 RED: Test for inventory list query with filters
2. 🟩 GREEN: Implement filtered list query
3. 🟥 RED: Test for laptop creation validation
4. 🟩 GREEN: Implement validator
5. 🟥 RED: Test for add laptop handler
6. 🟩 GREEN: Implement create handler
7. 🟥 RED: Test for edit laptop handler
8. 🟩 GREEN: Implement update handler
9. 🟥 RED: Test for delete laptop handler (soft delete)
10. 🟩 GREEN: Implement delete handler

**Features**:
- List all laptops with status filter
- Search by serial number
- Add new laptop to inventory
- Edit laptop details
- Retire laptop (soft delete)
- View laptop history

**Estimated Tests**: 12-15 test cases

## Total Estimates for Phase 6

- **Test Cases**: 34-43 tests
- **Production Files**: 12 new files
- **Test Files**: 6 new files
- **Lines of Code**: ~1,500 production + ~1,800 test
- **Time**: 2-3 days with TDD

## Success Criteria

Phase 6 will be considered complete when:

- ✅ Dashboard displays key metrics (shipments, inventory, delivery times)
- ✅ Charts render correctly with real data
- ✅ Calendar shows upcoming events with color coding
- ✅ Inventory management allows CRUD operations
- ✅ All new tests pass (34-43 test cases)
- ✅ Test coverage > 80% on new code
- ✅ All features accessible by correct roles
- ✅ Responsive design works on mobile
- ✅ Documentation updated

## Risks & Mitigations

### Risk 1: Chart Library Integration Issues
**Impact**: Medium  
**Probability**: Low  
**Mitigation**: Use well-established libraries (Chart.js), test in multiple browsers

### Risk 2: Performance with Large Datasets
**Impact**: High  
**Probability**: Medium  
**Mitigation**: Add pagination, implement caching, create database indexes

### Risk 3: Calendar Complexity
**Impact**: Medium  
**Probability**: Medium  
**Mitigation**: Start with simple month view, add features incrementally

### Risk 4: Time Estimation
**Impact**: Low  
**Probability**: Medium  
**Mitigation**: Follow TDD strictly, break into smaller tasks, commit frequently

## Pre-Phase Checklist

Before starting Phase 6, ensure:

- [ ] All Phase 0-5 code is committed
- [ ] No pending changes in working directory
- [ ] Database is up to date with migrations
- [ ] Test database is configured (optional but recommended)
- [ ] Development environment is working
- [ ] Tailwind CSS is compiled
- [ ] Documentation is up to date

## Starting Phase 6

When ready to begin:

1. Review this checklist
2. Review `docs/plan.md` Phase 6 section
3. Start with 6.1: Dashboard Statistics
4. Follow TDD workflow (red/green/refactor)
5. Commit after each GREEN step
6. Update `docs/plan.md` as tasks complete

---

## Ready to Proceed?

**Current Project Status**: ✅ READY FOR PHASE 6

All prerequisites are met. The foundation is solid with:
- ✅ Complete data models
- ✅ Working authentication
- ✅ Functional forms
- ✅ JIRA integration
- ✅ Email notifications

You can confidently start Phase 6 with the knowledge that the core system is stable and well-tested.

---

**Last Updated**: November 3, 2025  
**Next Phase**: Phase 6 - Dashboard & Visualization  
**Estimated Start**: Ready to begin immediately

