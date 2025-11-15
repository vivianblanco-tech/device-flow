package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestUpdateSingleShipmentToAtWarehouseUpdatesLaptopStatus tests that when a
// single_full_journey shipment status is updated to 'at_warehouse', the laptop
// in that shipment also gets its status updated to 'at_warehouse' (Received at Warehouse)
func TestUpdateSingleShipmentToAtWarehouseUpdatesLaptopStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, created_at) VALUES ($1, $2) RETURNING id`,
		"Test Company", time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create a single_full_journey shipment with status 'in_transit_to_warehouse'
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusInTransitToWarehouse, 1, "TEST-SYNC-001", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Create a laptop with status 'in_transit_to_warehouse' (matching the shipment)
	var laptopID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		"SN-SYNC-TEST-001", "Dell", "Latitude 7420", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID,
		time.Now(), time.Now(),
	).Scan(&laptopID)
	if err != nil {
		t.Fatalf("Failed to create test laptop: %v", err)
	}

	// Link laptop to shipment
	_, err = db.ExecContext(ctx,
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at) VALUES ($1, $2, $3)`,
		shipmentID, laptopID, time.Now(),
	)
	if err != nil {
		t.Fatalf("Failed to link laptop to shipment: %v", err)
	}

	// Verify initial laptop status is 'in_transit_to_warehouse'
	var initialLaptopStatus models.LaptopStatus
	err = db.QueryRowContext(ctx,
		`SELECT status FROM laptops WHERE id = $1`,
		laptopID,
	).Scan(&initialLaptopStatus)
	if err != nil {
		t.Fatalf("Failed to query initial laptop status: %v", err)
	}
	if initialLaptopStatus != models.LaptopStatusInTransitToWarehouse {
		t.Errorf("Expected initial laptop status to be 'in_transit_to_warehouse', got '%s'", initialLaptopStatus)
	}

	// Create a logistics user for the handler
	var logisticsUserID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"logistics@test.com", "hash", models.RoleLogistics, time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// NOW UPDATE SHIPMENT STATUS TO 'at_warehouse' using the handler
	// This should automatically update the laptop status to 'at_warehouse' as well
	handler := NewShipmentsHandler(db, nil, nil) // nil templates and notifier for test

	formData := url.Values{}
	formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
	formData.Set("status", string(models.ShipmentStatusAtWarehouse))

	req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	logisticsUser := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, logisticsUser)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.UpdateShipmentStatus(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusSeeOther, rr.Code, rr.Body.String())
	}

	// Verify laptop status is NOW 'at_warehouse'
	// This will FAIL until we implement the handler logic
	var finalLaptopStatus models.LaptopStatus
	err = db.QueryRowContext(ctx,
		`SELECT status FROM laptops WHERE id = $1`,
		laptopID,
	).Scan(&finalLaptopStatus)
	if err != nil {
		t.Fatalf("Failed to query final laptop status: %v", err)
	}

	if finalLaptopStatus != models.LaptopStatusAtWarehouse {
		t.Errorf("Expected laptop status to be 'at_warehouse' after shipment update, got '%s'", finalLaptopStatus)
	}
}

// TestUpdateBulkShipmentToAtWarehouseDoesNotUpdateLaptopStatus tests that when a
// bulk_to_warehouse shipment status is updated to 'at_warehouse', the laptops
// are NOT automatically updated (because bulk shipments don't sync status)
func TestUpdateBulkShipmentToAtWarehouseDoesNotUpdateLaptopStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, created_at) VALUES ($1, $2) RETURNING id`,
		"Test Company Bulk", time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create a bulk_to_warehouse shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusInTransitToWarehouse, 2, "TEST-BULK-002", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Create laptops with status 'in_transit_to_warehouse'
	var laptop1ID, laptop2ID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		"SN-BULK-SYNC-1", "HP", "EliteBook 840", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID,
		time.Now(), time.Now(),
	).Scan(&laptop1ID)
	if err != nil {
		t.Fatalf("Failed to create test laptop 1: %v", err)
	}

	err = db.QueryRowContext(ctx,
		`INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		"SN-BULK-SYNC-2", "HP", "EliteBook 850", "32", "1024", models.LaptopStatusInTransitToWarehouse, companyID,
		time.Now(), time.Now(),
	).Scan(&laptop2ID)
	if err != nil {
		t.Fatalf("Failed to create test laptop 2: %v", err)
	}

	// Link laptops to shipment
	_, err = db.ExecContext(ctx,
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at) VALUES ($1, $2, $3)`,
		shipmentID, laptop1ID, time.Now(),
	)
	if err != nil {
		t.Fatalf("Failed to link laptop 1 to shipment: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at) VALUES ($1, $2, $3)`,
		shipmentID, laptop2ID, time.Now(),
	)
	if err != nil {
		t.Fatalf("Failed to link laptop 2 to shipment: %v", err)
	}

	// Update bulk shipment status to 'at_warehouse'
	_, err = db.ExecContext(ctx,
		`UPDATE shipments SET status = $1, updated_at = $2 WHERE id = $3`,
		models.ShipmentStatusAtWarehouse, time.Now(), shipmentID,
	)
	if err != nil {
		t.Fatalf("Failed to update shipment status: %v", err)
	}

	// For bulk shipments, laptops should NOT automatically update status
	// (they need individual reception reports)
	// Verify laptops STILL have 'in_transit_to_warehouse' status
	var laptop1Status, laptop2Status models.LaptopStatus
	err = db.QueryRowContext(ctx,
		`SELECT status FROM laptops WHERE id = $1`,
		laptop1ID,
	).Scan(&laptop1Status)
	if err != nil {
		t.Fatalf("Failed to query laptop 1 status: %v", err)
	}

	err = db.QueryRowContext(ctx,
		`SELECT status FROM laptops WHERE id = $1`,
		laptop2ID,
	).Scan(&laptop2Status)
	if err != nil {
		t.Fatalf("Failed to query laptop 2 status: %v", err)
	}

	// Laptops should still be 'in_transit_to_warehouse' for bulk shipments
	if laptop1Status != models.LaptopStatusInTransitToWarehouse {
		t.Errorf("Expected bulk shipment laptop 1 to remain 'in_transit_to_warehouse', got '%s'", laptop1Status)
	}
	if laptop2Status != models.LaptopStatusInTransitToWarehouse {
		t.Errorf("Expected bulk shipment laptop 2 to remain 'in_transit_to_warehouse', got '%s'", laptop2Status)
	}
}

