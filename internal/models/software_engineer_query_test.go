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

	// Call GetAllSoftwareEngineers with nil filter (should return all)
	engineers, err := GetAllSoftwareEngineers(db, nil)
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

// 🟥 RED: Test search functionality in GetAllSoftwareEngineers
func TestGetAllSoftwareEngineers_WithSearchFilter(t *testing.T) {
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
		EmployeeNumber: "EMP001",
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

	// Test search by name
	filter := &SoftwareEngineerFilter{
		Search: "Alice",
	}
	engineers, err := GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with search filter failed: %v", err)
	}

	if len(engineers) != 1 {
		t.Errorf("Expected 1 engineer matching 'Alice', got %d", len(engineers))
	}
	if len(engineers) > 0 && engineers[0].Name != "Alice Johnson" {
		t.Errorf("Expected engineer name 'Alice Johnson', got %s", engineers[0].Name)
	}

	// Test search by email
	filter = &SoftwareEngineerFilter{
		Search: "bob.smith",
	}
	engineers, err = GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with email search failed: %v", err)
	}

	if len(engineers) != 1 {
		t.Errorf("Expected 1 engineer matching 'bob.smith', got %d", len(engineers))
	}
	if len(engineers) > 0 && engineers[0].Email != "bob.smith@example.com" {
		t.Errorf("Expected engineer email 'bob.smith@example.com', got %s", engineers[0].Email)
	}

	// Test search by employee number
	filter = &SoftwareEngineerFilter{
		Search: "EMP001",
	}
	engineers, err = GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with employee number search failed: %v", err)
	}

	if len(engineers) != 1 {
		t.Errorf("Expected 1 engineer matching 'EMP001', got %d", len(engineers))
	}
	if len(engineers) > 0 && engineers[0].EmployeeNumber != "EMP001" {
		t.Errorf("Expected employee number 'EMP001', got %s", engineers[0].EmployeeNumber)
	}

	// Test search with no results
	filter = &SoftwareEngineerFilter{
		Search: "NonExistent",
	}
	engineers, err = GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with no-match search failed: %v", err)
	}

	if len(engineers) != 0 {
		t.Errorf("Expected 0 engineers matching 'NonExistent', got %d", len(engineers))
	}
}

// 🟥 RED: Test sort functionality in GetAllSoftwareEngineers
func TestGetAllSoftwareEngineers_WithSortFilter(t *testing.T) {
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

	// Test sort by name descending
	filter := &SoftwareEngineerFilter{
		SortBy:    "name",
		SortOrder: "desc",
	}
	engineers, err := GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with sort filter failed: %v", err)
	}

	if len(engineers) != 3 {
		t.Errorf("Expected 3 engineers, got %d", len(engineers))
	}
	if len(engineers) >= 3 {
		if engineers[0].Name != "Charlie Brown" {
			t.Errorf("Expected first engineer to be 'Charlie Brown' (desc), got %s", engineers[0].Name)
		}
		if engineers[2].Name != "Alice Johnson" {
			t.Errorf("Expected last engineer to be 'Alice Johnson' (desc), got %s", engineers[2].Name)
		}
	}

	// Test sort by email ascending
	filter = &SoftwareEngineerFilter{
		SortBy:    "email",
		SortOrder: "asc",
	}
	engineers, err = GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with email sort failed: %v", err)
	}

	if len(engineers) != 3 {
		t.Errorf("Expected 3 engineers, got %d", len(engineers))
	}
	if len(engineers) >= 3 {
		if engineers[0].Email != "alice.johnson@example.com" {
			t.Errorf("Expected first engineer email 'alice.johnson@example.com', got %s", engineers[0].Email)
		}
		if engineers[2].Email != "charlie.brown@example.com" {
			t.Errorf("Expected last engineer email 'charlie.brown@example.com', got %s", engineers[2].Email)
		}
	}

	// Test default sort (by name ascending)
	filter = &SoftwareEngineerFilter{}
	engineers, err = GetAllSoftwareEngineers(db, filter)
	if err != nil {
		t.Fatalf("GetAllSoftwareEngineers with default sort failed: %v", err)
	}

	if len(engineers) != 3 {
		t.Errorf("Expected 3 engineers, got %d", len(engineers))
	}
	if len(engineers) >= 3 {
		if engineers[0].Name != "Alice Johnson" {
			t.Errorf("Expected first engineer to be 'Alice Johnson' (default), got %s", engineers[0].Name)
		}
	}
}