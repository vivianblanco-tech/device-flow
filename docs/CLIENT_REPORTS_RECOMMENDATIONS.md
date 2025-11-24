# Client Reports Recommendations

## Overview

This document outlines recommended reports that can be created for Client users to provide insights into their inventory and shipments. These reports leverage the existing data models and respect role-based access control (Clients can only see their own company's data).

## Current Client Access

Based on the codebase analysis, Client users currently have access to:
- **Shipments**: View their company's shipments with filtering by status, type, and search
- **Inventory**: View their company's laptops with filtering by status and search
- **Pickup Forms**: Submit and view pickup forms for their company
- **Calendar**: View shipment schedules

## Recommended Reports

### 1. Inventory Reports

#### 1.1 Inventory Summary Report
**Purpose**: Provide a high-level overview of the client's laptop inventory status

**Data Points**:
- Total laptops owned by the company
- Laptops by status (Available, In Transit, At Warehouse, Delivered, Retired)
- Laptops by brand/model breakdown
- Laptops by assignment status (assigned to engineers vs unassigned)
- Average time laptops spend in each status
- Laptops pending reception report approval

**Visualizations**:
- Donut chart: Status distribution
- Bar chart: Brand/model distribution
- Table: Detailed breakdown with serial numbers

**Use Cases**:
- Monthly inventory audits
- Budget planning for new laptop purchases
- Tracking device utilization rates

---

#### 1.2 Laptop Lifecycle Report
**Purpose**: Track individual laptops through their complete lifecycle

**Data Points**:
- Serial number and device specifications
- Current status and location
- Assigned engineer (if applicable)
- Date added to system
- Date delivered (if delivered)
- Time spent in each status
- Associated shipments and JIRA tickets
- Reception report status and dates

**Visualizations**:
- Timeline view: Visual representation of laptop journey
- Status transition history
- Table: Sortable list with all lifecycle events

**Use Cases**:
- Tracking specific device history
- Audit trails for compliance
- Identifying bottlenecks in the delivery process

---

#### 1.3 Inventory Aging Report
**Purpose**: Identify laptops that have been in specific statuses for extended periods

**Data Points**:
- Laptops in transit to warehouse > X days
- Laptops at warehouse > X days
- Laptops in transit to engineer > X days
- Average time per status
- Laptops with overdue status transitions

**Visualizations**:
- Bar chart: Days in current status
- Alert indicators for overdue items
- Table: Detailed list with days in status

**Use Cases**:
- Identifying stuck shipments
- Proactive follow-up with logistics team
- Performance metrics tracking

---

### 2. Shipment Reports

#### 2.1 Shipment Status Dashboard
**Purpose**: Real-time overview of all active shipments

**Data Points**:
- Total shipments (all time, current month, current year)
- Shipments by status breakdown
- Shipments by type (Single Full Journey, Bulk to Warehouse, Warehouse to Engineer)
- Average delivery time by shipment type
- Pending pickups count
- In-transit shipments count
- Delivered shipments count

**Visualizations**:
- Summary cards: Key metrics
- Donut chart: Status distribution
- Line chart: Shipments over time (monthly/quarterly)
- Table: Active shipments list

**Use Cases**:
- Daily operations monitoring
- Executive dashboards
- Performance tracking

---

#### 2.2 Shipment Timeline Report
**Purpose**: Detailed tracking of individual shipments with all milestones

**Data Points**:
- Shipment ID and JIRA ticket number
- Shipment type and laptop count
- Status progression with timestamps:
  - Pickup scheduled date
  - Picked up at
  - Arrived warehouse at
  - Released warehouse at
  - ETA to engineer
  - Delivered at
- Courier information and tracking numbers
- Associated laptops (serial numbers)
- Time between status transitions
- Total shipment duration

**Visualizations**:
- Gantt chart: Visual timeline of shipment progress
- Status flow diagram
- Table: Detailed milestone log

**Use Cases**:
- Tracking specific shipments
- Identifying delays in the process
- Performance analysis

---

#### 2.3 Shipment Performance Report
**Purpose**: Analyze shipment performance metrics and trends

**Data Points**:
- Average delivery time by shipment type
- On-time delivery rate (%)
- Average time per status transition
- Fastest/slowest shipments
- Shipments by courier performance
- Monthly/quarterly trends
- Status transition bottlenecks

**Visualizations**:
- Line chart: Average delivery time trends
- Bar chart: Performance by courier
- Heatmap: Status transition times
- Table: Top performers and outliers

**Use Cases**:
- Performance reviews
- Vendor evaluation (courier selection)
- Process optimization identification

---

#### 2.4 Pickup Form Summary Report
**Purpose**: Overview of all pickup forms submitted by the client

**Data Points**:
- Total pickup forms submitted
- Forms by status (pending, scheduled, completed)
- Pickup locations breakdown
- Average number of laptops per pickup
- Pickup date distribution
- Time slots preferred
- Special instructions frequency
- Forms with accessories included

**Visualizations**:
- Bar chart: Forms by status
- Map view: Pickup locations (if coordinates available)
- Pie chart: Time slot preferences
- Table: Detailed form list

**Use Cases**:
- Understanding pickup patterns
- Planning future pickups
- Identifying recurring issues

---

### 3. Combined Reports

#### 3.1 Company Overview Report
**Purpose**: Comprehensive view of all company-related data

**Data Points**:
- Total laptops in inventory
- Active shipments count
- Delivered shipments (current month/year)
- Pending pickups
- Laptops at warehouse
- Laptops in transit
- Total engineers assigned laptops
- Average delivery time

**Visualizations**:
- Executive summary cards
- Key metrics dashboard
- Recent activity feed

**Use Cases**:
- Monthly/quarterly business reviews
- Executive presentations
- Stakeholder updates

---

#### 3.2 Delivery Performance Report
**Purpose**: Track delivery performance and identify improvement areas

**Data Points**:
- Delivery success rate
- Average time from pickup to delivery
- Average time from warehouse release to delivery
- On-time delivery percentage
- Delayed shipments analysis
- Courier performance comparison
- Engineer assignment time

**Visualizations**:
- Performance scorecards
- Trend analysis charts
- Comparison charts (month-over-month, year-over-year)
- Table: Detailed delivery metrics

**Use Cases**:
- Service level agreement (SLA) tracking
- Vendor performance evaluation
- Process improvement initiatives

---

### 4. Export and Scheduling Features

#### 4.1 Report Export Options
- **PDF Export**: Formatted reports for sharing
- **CSV Export**: Raw data for analysis in Excel/Google Sheets
- **Excel Export**: Formatted spreadsheets with charts
- **Email Reports**: Scheduled email delivery

#### 4.2 Report Scheduling
- Daily summary reports
- Weekly performance reports
- Monthly inventory reports
- Custom date range reports

---

## Implementation Considerations

### Data Access
- All reports must filter by `client_company_id` to ensure Clients only see their own data
- Use existing role-based filtering mechanisms from `internal/handlers/shipments.go` and `internal/handlers/inventory.go`

### Performance
- Consider database indexes on frequently queried fields:
  - `shipments.client_company_id`
  - `shipments.status`
  - `shipments.created_at`
  - `laptops.client_company_id`
  - `laptops.status`
- Implement pagination for large datasets
- Cache frequently accessed reports

### UI/UX
- Use existing Tailwind CSS styling for consistency
- Leverage Chart.js (already in use) for visualizations
- Follow existing navigation patterns
- Add report filters (date range, status, type)
- Include "Export" buttons on all reports

### Security
- Ensure all reports respect role-based access control
- Validate `client_company_id` matches user's company
- Sanitize all user inputs for date ranges and filters
- Implement rate limiting for report generation

---

## Priority Recommendations

### Phase 1 (High Priority)
1. **Shipment Status Dashboard** - Most requested by clients
2. **Inventory Summary Report** - Essential for inventory management
3. **Shipment Timeline Report** - Critical for tracking individual shipments

### Phase 2 (Medium Priority)
4. **Shipment Performance Report** - Useful for process improvement
5. **Laptop Lifecycle Report** - Important for audit trails
6. **Pickup Form Summary Report** - Helps understand pickup patterns

### Phase 3 (Nice to Have)
7. **Company Overview Report** - Executive-level insights
8. **Inventory Aging Report** - Proactive issue identification
9. **Delivery Performance Report** - Advanced analytics

---

## Technical Implementation Notes

### New Routes Needed
```go
// Reports routes (Client role only)
protected.HandleFunc("/reports", reportsHandler.ReportsIndex).Methods("GET")
protected.HandleFunc("/reports/inventory-summary", reportsHandler.InventorySummary).Methods("GET")
protected.HandleFunc("/reports/shipment-status", reportsHandler.ShipmentStatus).Methods("GET")
protected.HandleFunc("/reports/shipment-timeline/{id}", reportsHandler.ShipmentTimeline).Methods("GET")
// ... additional report routes
```

### New Handler Structure
```go
// internal/handlers/reports.go
type ReportsHandler struct {
    DB        *sql.DB
    Templates *template.Template
}

// Methods for each report type
func (h *ReportsHandler) InventorySummary(w http.ResponseWriter, r *http.Request)
func (h *ReportsHandler) ShipmentStatus(w http.ResponseWriter, r *http.Request)
// ... additional report methods
```

### New Model Queries
- Add query methods to `internal/models/` for report-specific data aggregation
- Consider creating `internal/models/reports.go` for report-specific queries
- Leverage existing query patterns from `internal/models/dashboard.go` and `internal/models/charts.go`

---

## Conclusion

These reports will provide Clients with comprehensive insights into their inventory and shipments, enabling better decision-making and improved visibility into the laptop tracking process. The reports are designed to leverage existing data structures and respect the current role-based access control system.

