package models

import (
	"testing"
	"time"
)

// TestGenerateCalendarGrid tests generating a calendar grid with weeks and days
func TestGenerateCalendarGrid(t *testing.T) {
	// Test with November 2025 which starts on a Saturday (day 6)
	year := 2025
	month := time.November

	grid := GenerateCalendarGrid(year, month)

	// Test: Should have at least 4 weeks, typically 5-6
	if len(grid) < 4 {
		t.Errorf("Expected at least 4 weeks, got %d", len(grid))
	}

	// Test: Each week should have exactly 7 days
	for i, week := range grid {
		if len(week) != 7 {
			t.Errorf("Week %d should have 7 days, got %d", i, len(week))
		}
	}

	// Test: First day of grid should be a Sunday (index 0)
	firstDay := grid[0][0]
	if firstDay.Date.Weekday() != time.Sunday {
		t.Errorf("First day should be Sunday, got %s", firstDay.Date.Weekday())
	}

	// Test: Last day of grid should be a Saturday (index 6)
	lastWeek := grid[len(grid)-1]
	lastDay := lastWeek[6]
	if lastDay.Date.Weekday() != time.Saturday {
		t.Errorf("Last day should be Saturday, got %s", lastDay.Date.Weekday())
	}

	// Test: November 1, 2025 should be marked as in current month
	// November 1, 2025 is a Saturday, so it should be in the first week at index 6
	nov1 := grid[0][6]
	if !nov1.IsCurrentMonth {
		t.Error("November 1 should be marked as in current month")
	}
	if nov1.Date.Day() != 1 {
		t.Errorf("Expected day 1, got %d", nov1.Date.Day())
	}
	if nov1.Date.Month() != time.November {
		t.Errorf("Expected November, got %s", nov1.Date.Month())
	}

	// Test: Days before November 1 should not be in current month
	if grid[0][0].IsCurrentMonth {
		t.Error("Days before November 1 should not be in current month")
	}
}

// TestGenerateCalendarGridWithEvents tests adding events to calendar days
func TestGenerateCalendarGridWithEvents(t *testing.T) {
	year := 2025
	month := time.November

	// Create sample events for November 10 and November 15
	nov10 := time.Date(2025, time.November, 10, 12, 0, 0, 0, time.UTC)
	nov15 := time.Date(2025, time.November, 15, 14, 0, 0, 0, time.UTC)

	events := []CalendarEvent{
		{
			ID:         1,
			Type:       CalendarEventTypePickup,
			Title:      "Pickup from Test Corp",
			Date:       nov10,
			ShipmentID: 100,
		},
		{
			ID:         2,
			Type:       CalendarEventTypeDelivery,
			Title:      "Delivery to Engineer",
			Date:       nov15,
			ShipmentID: 101,
		},
	}

	grid := GenerateCalendarGridWithEvents(year, month, events)

	// Find November 10 in the grid
	var found10 bool
	var day10 *CalendarDay
	for _, week := range grid {
		for i := range week {
			if week[i].Date.Day() == 10 && week[i].Date.Month() == time.November {
				found10 = true
				day10 = &week[i]
				break
			}
		}
		if found10 {
			break
		}
	}

	// Test: November 10 should be found
	if !found10 {
		t.Fatal("November 10 should be in the grid")
	}

	// Test: November 10 should have 1 event
	if len(day10.Events) != 1 {
		t.Errorf("November 10 should have 1 event, got %d", len(day10.Events))
	}

	// Test: Event on November 10 should be the pickup event
	if len(day10.Events) > 0 {
		if day10.Events[0].Type != CalendarEventTypePickup {
			t.Errorf("Expected pickup event type, got %s", day10.Events[0].Type)
		}
		if day10.Events[0].Title != "Pickup from Test Corp" {
			t.Errorf("Expected 'Pickup from Test Corp', got %s", day10.Events[0].Title)
		}
	}

	// Find November 15 in the grid
	var found15 bool
	var day15 *CalendarDay
	for _, week := range grid {
		for i := range week {
			if week[i].Date.Day() == 15 && week[i].Date.Month() == time.November {
				found15 = true
				day15 = &week[i]
				break
			}
		}
		if found15 {
			break
		}
	}

	// Test: November 15 should have the delivery event
	if !found15 {
		t.Fatal("November 15 should be in the grid")
	}

	if len(day15.Events) != 1 {
		t.Errorf("November 15 should have 1 event, got %d", len(day15.Events))
	}

	if len(day15.Events) > 0 && day15.Events[0].Type != CalendarEventTypeDelivery {
		t.Errorf("Expected delivery event type, got %s", day15.Events[0].Type)
	}
}

// TestCalendarDayStructure tests the CalendarDay structure
func TestCalendarDayStructure(t *testing.T) {
	date := time.Date(2025, time.November, 15, 0, 0, 0, 0, time.UTC)
	
	day := CalendarDay{
		Date:           date,
		IsCurrentMonth: true,
		Events:         []CalendarEvent{},
	}

	// Test: Date is set correctly
	if day.Date.Day() != 15 {
		t.Errorf("Expected day 15, got %d", day.Date.Day())
	}

	// Test: IsCurrentMonth is set correctly
	if !day.IsCurrentMonth {
		t.Error("Expected IsCurrentMonth to be true")
	}

	// Test: Events slice is initialized
	if day.Events == nil {
		t.Error("Expected Events slice to be initialized")
	}

	// Test: Can add events to the day
	event := CalendarEvent{
		ID:         1,
		Type:       CalendarEventTypePickup,
		Title:      "Test Event",
		Date:       date,
		ShipmentID: 123,
	}
	day.Events = append(day.Events, event)

	if len(day.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(day.Events))
	}
}

