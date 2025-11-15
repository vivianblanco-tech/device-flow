package models

import (
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetAllSoftwareEngineers tests retrieving all software engineers from the database
func TestGetAllSoftwareEngineers(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test software engineers
	engineer1 := &SoftwareEngineer{
		Name:  "Alice Johnson",
		Email: "alice.johnson@example.com",
		Phone: "+1-555-0001",
	}
	engineer2 := &SoftwareEngineer{
		Name:  "Bob Smith",
		Email: "bob.smith@example.com",
		Phone: "+1-555-0002",
	}
	engineer3 := &SoftwareEngineer{
		Name:  "Charlie Brown",
		Email: "charlie.brown@example.com",
		Phone: "+1-555-0003",
	}

	// Insert engineers
	if err := CreateSoftwareEngineer(db, engineer1); err != nil {
		t.Fatalf("Failed to create engineer1: %v", err)
	}
	if err := CreateSoftwareEngineer(db, engineer2); err != nil {
		t.Fatalf("Failed to create engineer2: %v", err)
	}
	if err := CreateSoftwareEngineer(db, engineer3); err != nil {
		t.Fatalf("Failed to create engineer3: %v", err)
	}

	// Call GetAllSoftwareEngineers
	engineers, err := GetAllSoftwareEngineers(db)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers failed: %v", err)
	}

	// Verify results
	if len(engineers) != 3 {
		t.Errorf("Expected 3 engineers, got %d", len(engineers))
	}

	// Verify engineers are returned in alphabetical order by name
	if len(engineers) >= 3 {
		if engineers[0].Name != "Alice Johnson" {
			t.Errorf("Expected first engineer to be 'Alice Johnson', got %s", engineers[0].Name)
		}
		if engineers[1].Name != "Bob Smith" {
			t.Errorf("Expected second engineer to be 'Bob Smith', got %s", engineers[1].Name)
		}
		if engineers[2].Name != "Charlie Brown" {
			t.Errorf("Expected third engineer to be 'Charlie Brown', got %s", engineers[2].Name)
		}
	}

	// Verify all fields are populated
	if engineers[0].Email == "" {
		t.Error("Engineer email should be populated")
	}
	if engineers[0].ID == 0 {
		t.Error("Engineer ID should be populated")
	}
}
