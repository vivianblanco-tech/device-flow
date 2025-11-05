# Phase 6 Manual Testing Plan
**Date**: November 3, 2025  
**Phase**: Dashboard & Visualization  
**Tester**: _______________  
**Duration**: ~2 hours

---

## 🎯 Testing Objectives

Verify that all Phase 6 features work correctly in the browser:
- ✅ Dashboard displays statistics
- ✅ Charts render with real data
- ✅ Calendar shows events
- ✅ Inventory CRUD operations work
- ✅ Role-based access control
- ✅ Responsive design

---

## 📋 Pre-Testing Setup

### 1. Start Application
```powershell
# Ensure Docker PostgreSQL is running
docker-compose up -d postgres

# Start the application
go run cmd/web/main.go
```

**Expected**: Server starts on `http://localhost:8080`

### 2. Create Test Data (if needed)

**Option A: Use existing data** (if database has data from previous testing)

**Option B: Create fresh test data**
```powershell
# Connect to database
docker-compose exec postgres psql -U postgres -d laptop_tracking_dev

# Run test data scripts (if they exist)
\i scripts/create-test-data.sql
```

### 3. Prepare Test Users

You'll need access to users with different roles:
- **Logistics User**: Can access dashboard, calendar, inventory, all features
- **Client User**: Cannot access dashboard
- **Warehouse User**: Cannot access dashboard, can access inventory
- **Project Manager User**: Can access dashboard (read-only)

**Default Test Users** (from previous phases):
- Email: `logistics@bairesdev.com` | Role: Logistics
- Email: `client@bairesdev.com` | Role: Client
- Email: `warehouse@bairesdev.com` | Role: Warehouse
- Email: `pm@bairesdev.com` | Role: Project Manager

---

## 🧪 Test Cases

### Test Suite 1: Dashboard Access Control (15 min)

#### Test 1.1: Logistics User Access ✅
**Steps**:
1. Open browser: `http://localhost:8080/login`
2. Login as logistics user
3. Click "Dashboard" in navigation

**Expected**:
- ✅ Dashboard loads successfully
- ✅ Statistics cards display numbers
- ✅ "Dashboard" link visible in nav
- ✅ No error messages

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 1.2: Project Manager Access ✅
**Steps**:
1. Logout
2. Login as project manager user
3. Click "Dashboard" in navigation

**Expected**:
- ✅ Dashboard loads successfully
- ✅ Can view all statistics (read-only)
- ✅ "Dashboard" link visible in nav

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 1.3: Client User Denied ✅
**Steps**:
1. Logout
2. Login as client user
3. Try to access dashboard

**Expected**:
- ✅ "Dashboard" link NOT visible in nav
- ✅ Direct access to `/dashboard` shows "Forbidden" error

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 1.4: Warehouse User Denied ✅
**Steps**:
1. Logout
2. Login as warehouse user
3. Try to access dashboard

**Expected**:
- ✅ "Dashboard" link NOT visible in nav
- ✅ Direct access to `/dashboard` shows "Forbidden" error

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 2: Dashboard Statistics (20 min)

**Prerequisites**: Login as logistics user, navigate to dashboard

#### Test 2.1: Statistics Cards Display ✅
**Steps**:
1. Observe the 4 statistics cards at top of dashboard

**Expected**: All 4 cards show:
- ✅ **Total Shipments**: Shows number (can be 0)
- ✅ **Pending Pickups**: Shows number (can be 0)
- ✅ **In Transit**: Shows number (can be 0)
- ✅ **Delivered**: Shows number (can be 0)
- ✅ Cards have icons and colors
- ✅ Numbers are readable

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 2.2: Average Delivery Time ✅
**Steps**:
1. Look for "Average Delivery Time" section

**Expected**:
- ✅ Shows number of days (or "N/A" if no deliveries)
- ✅ Label is clear
- ✅ Formatting is correct

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 2.3: Inventory Statistics ✅
**Steps**:
1. Look for inventory breakdown section

**Expected**:
- ✅ Shows available laptops count
- ✅ Shows laptops by status
- ✅ Status labels are clear
- ✅ Colors match laptop statuses

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 2.4: Shipment Status Breakdown ✅
**Steps**:
1. Look for shipment status section

**Expected**:
- ✅ Lists all shipment statuses
- ✅ Shows count for each status
- ✅ Status labels are formatted nicely
- ✅ Colors are distinct

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 3: Charts and Visualization (25 min)

**Prerequisites**: Dashboard page loaded as logistics user

#### Test 3.1: Charts Load ✅
**Steps**:
1. Wait for page to fully load (2-3 seconds)
2. Observe the 3 chart sections

**Expected**:
- ✅ "Shipments Over Time" line chart appears
- ✅ "Status Distribution" donut chart appears
- ✅ "Delivery Time Trends" bar chart appears
- ✅ No JavaScript errors in console (press F12 to check)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 3.2: Line Chart - Shipments Over Time ✅
**Steps**:
1. Examine the line chart at top

**Expected**:
- ✅ Chart has X-axis (dates)
- ✅ Chart has Y-axis (counts)
- ✅ Blue line connects data points
- ✅ Shows last 30 days of data
- ✅ Hover shows tooltips with values
- ✅ Chart is responsive (resize browser to test)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 3.3: Donut Chart - Status Distribution ✅
**Steps**:
1. Examine the donut chart

**Expected**:
- ✅ Different colored segments for each status
- ✅ Legend shows status names and colors
- ✅ Hover shows percentage/count
- ✅ Center shows total or main value
- ✅ Colors match status colors elsewhere

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 3.4: Bar Chart - Delivery Time Trends ✅
**Steps**:
1. Examine the bar chart at bottom

**Expected**:
- ✅ Bars show average delivery time per week/month
- ✅ X-axis shows time periods
- ✅ Y-axis shows days
- ✅ Hover shows exact values
- ✅ Bars have consistent color

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 3.5: Charts with Empty Data ✅
**Steps**:
1. If database is empty, charts should handle gracefully

**Expected** (if no data):
- ✅ Charts show empty state or zero values
- ✅ No JavaScript errors
- ✅ No broken images
- ✅ Message like "No data available" (optional)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 4: Calendar View (20 min)

**Prerequisites**: Login as any authenticated user

#### Test 4.1: Calendar Access ✅
**Steps**:
1. Click "Calendar" in navigation menu

**Expected**:
- ✅ Calendar page loads
- ✅ Shows current month by default
- ✅ Month/year displayed at top
- ✅ No errors

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 4.2: Calendar Display ✅
**Steps**:
1. Observe calendar layout

**Expected**:
- ✅ Days of week headers (Sun-Sat)
- ✅ Date numbers visible
- ✅ Current day highlighted (if current month)
- ✅ Responsive grid layout
- ✅ Professional appearance

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 4.3: Calendar Events ✅
**Steps**:
1. Look for events on calendar dates

**Expected** (if events exist):
- ✅ Events appear on correct dates
- ✅ Color-coded by type:
  - 🔵 Pickup scheduled
  - 🟡 In transit to warehouse
  - 🟢 In transit to engineer
  - 🟣 Delivery scheduled
- ✅ Event titles are readable
- ✅ Multiple events per day stack properly

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 4.4: Calendar Navigation ✅
**Steps**:
1. Click "Previous Month" or "←" button
2. Click "Next Month" or "→" button
3. Click "Today" button (if exists)

**Expected**:
- ✅ Previous month loads with correct dates
- ✅ Next month loads with correct dates
- ✅ Events update for selected month
- ✅ Month/year label updates
- ✅ Navigation is smooth (no flicker)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 4.5: Event Details ✅
**Steps**:
1. Click on an event (if interactive)

**Expected**:
- ✅ Shows event details (shipment info, type, date)
- ✅ Link to shipment (if applicable)
- ✅ Can close detail view
- ✅ Details are accurate

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 5: Inventory Management (30 min)

**Prerequisites**: Login as logistics or warehouse user

#### Test 5.1: Inventory List Access ✅
**Steps**:
1. Click "Inventory" in navigation menu

**Expected**:
- ✅ Inventory list page loads
- ✅ Shows all laptops (or message if empty)
- ✅ Table/grid layout is clean
- ✅ "Add Laptop" button visible (for logistics/warehouse)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.2: Inventory List Display ✅
**Steps**:
1. Observe laptop list

**Expected**: Each laptop shows:
- ✅ Serial Number
- ✅ Brand
- ✅ Model
- ✅ Status (with color badge)
- ✅ Actions (View, Edit, Delete buttons)
- ✅ Table is sortable/organized
- ✅ Status colors match dashboard

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.3: Search/Filter Laptops ✅
**Steps**:
1. Use search box (if present)
2. Type serial number or brand
3. Use status filter dropdown (if present)

**Expected**:
- ✅ Search filters list in real-time
- ✅ Results match search term
- ✅ Clear search returns full list
- ✅ Status filter shows only matching laptops

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.4: Add New Laptop ✅
**Steps**:
1. Click "Add Laptop" button
2. Fill form:
   - Serial Number: TEST-001
   - Brand: Dell
   - Model: XPS 15
   - Specs: i7, 16GB RAM, 512GB SSD
   - Status: Available
3. Click "Save" or "Add"

**Expected**:
- ✅ Form validation works (required fields)
- ✅ Submission succeeds
- ✅ Redirected to inventory list
- ✅ New laptop appears in list
- ✅ Success message shown (optional)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.5: View Laptop Details ✅
**Steps**:
1. Click on a laptop or "View" button

**Expected**:
- ✅ Detail page loads
- ✅ Shows all laptop information
- ✅ Status badge visible
- ✅ Edit/Delete buttons present (for authorized users)
- ✅ Shipment history (if applicable)

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.6: Edit Laptop ✅
**Steps**:
1. From laptop detail or list, click "Edit"
2. Modify Brand to "HP"
3. Modify Status to "At Warehouse"
4. Click "Update" or "Save"

**Expected**:
- ✅ Edit form loads with current values
- ✅ Can modify all fields
- ✅ Update succeeds
- ✅ Redirected to detail or list
- ✅ Changes are saved and visible

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.7: Delete Laptop ✅
**Steps**:
1. From laptop detail or list, click "Delete"
2. Confirm deletion (if confirmation prompt)

**Expected**:
- ✅ Confirmation dialog appears (good UX)
- ✅ Deletion succeeds
- ✅ Redirected to inventory list
- ✅ Laptop removed from list
- ✅ Success message shown

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 5.8: Inventory Access Control ✅
**Steps**:
1. Logout
2. Login as client user
3. Try to access inventory

**Expected**:
- ✅ Can VIEW inventory (read-only)
- ✅ Cannot see "Add Laptop" button
- ✅ Cannot Edit or Delete

**Alternative** (if client has no access):
- ✅ "Inventory" link not visible
- ✅ Direct access shows forbidden error

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 6: Mobile Responsiveness (15 min)

**Prerequisites**: Any user logged in

#### Test 6.1: Dashboard on Mobile ✅
**Steps**:
1. Resize browser to mobile size (375x667px)
2. Or use browser DevTools mobile emulation
3. Navigate to dashboard

**Expected**:
- ✅ Statistics cards stack vertically
- ✅ Charts resize and remain readable
- ✅ Navigation menu adapts (hamburger menu)
- ✅ All text is readable
- ✅ No horizontal scroll

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 6.2: Calendar on Mobile ✅
**Steps**:
1. Keep mobile viewport
2. Navigate to calendar

**Expected**:
- ✅ Calendar grid adapts to small screen
- ✅ Events are readable
- ✅ Navigation buttons accessible
- ✅ No layout breaks

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

#### Test 6.3: Inventory on Mobile ✅
**Steps**:
1. Keep mobile viewport
2. Navigate to inventory

**Expected**:
- ✅ Table converts to cards or stacks
- ✅ All laptop info visible
- ✅ Action buttons accessible
- ✅ Add button visible
- ✅ Search/filter works

**Result**: ☐ Pass  ☐ Fail  
**Notes**: _______________________

---

### Test Suite 7: Browser Compatibility (Optional, 15 min)

#### Test 7.1: Chrome ✅
**Result**: ☐ Pass  ☐ Fail  
**Version**: _______  

#### Test 7.2: Firefox ✅
**Result**: ☐ Pass  ☐ Fail  
**Version**: _______  

#### Test 7.3: Edge ✅
**Result**: ☐ Pass  ☐ Fail  
**Version**: _______  

#### Test 7.4: Safari (Mac only) ✅
**Result**: ☐ Pass  ☐ Fail  
**Version**: _______  

---

## 🐛 Bugs Found

| # | Component | Severity | Description | Steps to Reproduce |
|---|-----------|----------|-------------|--------------------|
| 1 |           |          |             |                    |
| 2 |           |          |             |                    |
| 3 |           |          |             |                    |

**Severity Levels**:
- 🔴 Critical: Blocks functionality
- 🟡 Major: Impacts usability
- 🟢 Minor: Cosmetic or edge case

---

## ✅ Testing Summary

### Test Results
- **Total Test Cases**: 33
- **Passed**: _____ / 33
- **Failed**: _____ / 33
- **Pass Rate**: _____%

### Feature Status
| Feature | Status | Notes |
|---------|--------|-------|
| Dashboard Statistics | ☐ ✅ ☐ ⚠️ ☐ ❌ | |
| Charts | ☐ ✅ ☐ ⚠️ ☐ ❌ | |
| Calendar | ☐ ✅ ☐ ⚠️ ☐ ❌ | |
| Inventory CRUD | ☐ ✅ ☐ ⚠️ ☐ ❌ | |
| Access Control | ☐ ✅ ☐ ⚠️ ☐ ❌ | |
| Responsive Design | ☐ ✅ ☐ ⚠️ ☐ ❌ | |

### Overall Assessment
☐ **Ready for Production**  
☐ **Minor Issues - Can Deploy**  
☐ **Major Issues - Needs Fixes**  
☐ **Critical Issues - Cannot Deploy**

---

## 📝 Recommendations

### Immediate Actions
1. _______________________________________
2. _______________________________________
3. _______________________________________

### Future Improvements
1. _______________________________________
2. _______________________________________
3. _______________________________________

---

## ✅ Sign-Off

**Tested By**: _______________  
**Date**: _______________  
**Time Spent**: _____ hours  
**Phase 6 Status**: ☐ **APPROVED** ☐ **NEEDS WORK**

---

**Signature**: _______________

---

## 📎 Attachments

- [ ] Screenshots of dashboard
- [ ] Screenshots of charts
- [ ] Screenshots of calendar
- [ ] Screenshots of inventory
- [ ] Browser console logs (if errors)
- [ ] Network tab (if API issues)

---

**Last Updated**: November 3, 2025  
**Next Review**: After fixes (if any)

