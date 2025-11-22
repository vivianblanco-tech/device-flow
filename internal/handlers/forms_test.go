package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestFormsHandler_FormsPage_RequiresLogisticsRole(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	tests := []struct {
		name           string
		userRole       models.UserRole
		expectedStatus int
	}{
		{
			name:           "logistics user can access forms page",
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "client user cannot access forms page",
			userRole:       models.RoleClient,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "warehouse user cannot access forms page",
			userRole:       models.RoleWarehouse,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "project manager cannot access forms page",
			userRole:       models.RoleProjectManager,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test user
			user := &models.User{
				ID:    1,
				Email: "test@example.com",
				Role:  tt.userRole,
			}

			// Create request
			req := httptest.NewRequest("GET", "/forms", nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.FormsPage(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestClientCompanyEditPage_FormatsJSONContactInfo tests that JSON contact info is converted to plain text in the form
func TestClientCompanyEditPage_FormatsJSONContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create test company with JSON contact info
	var companyID int64
	jsonContactInfo := `{"email":"contact@example.com","phone":"+1-555-0100","address":"123 Main St"}`
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"Test Company", jsonContactInfo, time.Now(), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	req := httptest.NewRequest("GET", "/forms/client-companies/"+strconv.FormatInt(companyID, 10)+"/edit", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(companyID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyEditPage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check that the response contains formatted contact info (not raw JSON)
	body := rr.Body.String()
	if strings.Contains(body, jsonContactInfo) {
		t.Error("Form should not display raw JSON contact info")
	}
	// The formatted version should have "Email:" and "Phone:" labels
	if !strings.Contains(body, "Email:") || !strings.Contains(body, "Phone:") {
		t.Error("Form should display formatted contact info with labels")
	}
}

// TestClientCompanyAddSubmit_AcceptsPlainTextContactInfo tests that plain text contact info is stored correctly
func TestClientCompanyAddSubmit_AcceptsPlainTextContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with plain text contact info
	formData := url.Values{}
	formData.Set("name", "New Company")
	formData.Set("contact_info", "Email: contact@example.com\nPhone: +1-555-0100")

	req := httptest.NewRequest("POST", "/forms/client-companies/add", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyAddSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the company was created with plain text contact info
	var contactInfo string
	err = db.QueryRowContext(context.Background(),
		`SELECT contact_info FROM client_companies WHERE name = $1`,
		"New Company",
	).Scan(&contactInfo)
	if err != nil {
		t.Fatalf("Failed to retrieve created company: %v", err)
	}

	expected := "Email: contact@example.com\nPhone: +1-555-0100"
	if contactInfo != expected {
		t.Errorf("expected contact_info %q, got %q", expected, contactInfo)
	}
}

// TestClientCompanyEditSubmit_AcceptsPlainTextContactInfo tests that plain text contact info updates are stored correctly
func TestClientCompanyEditSubmit_AcceptsPlainTextContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create test company with JSON contact info (simulating old data)
	var companyID int64
	jsonContactInfo := `{"email":"old@example.com","phone":"+1-555-0000"}`
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"Test Company", jsonContactInfo, time.Now(), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with plain text contact info
	formData := url.Values{}
	formData.Set("name", "Test Company")
	formData.Set("contact_info", "Email: new@example.com\nPhone: +1-555-9999\nAddress: 456 New St")

	req := httptest.NewRequest("POST", "/forms/client-companies/"+strconv.FormatInt(companyID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(companyID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the company was updated with plain text contact info
	var contactInfo string
	err = db.QueryRowContext(ctx,
		`SELECT contact_info FROM client_companies WHERE id = $1`,
		companyID,
	).Scan(&contactInfo)
	if err != nil {
		t.Fatalf("Failed to retrieve updated company: %v", err)
	}

	expected := "Email: new@example.com\nPhone: +1-555-9999\nAddress: 456 New St"
	if contactInfo != expected {
		t.Errorf("expected contact_info %q, got %q", expected, contactInfo)
	}
}

// # RED: Test that address confirmation timestamp is set when checkbox is checked in Add form
func TestSoftwareEngineerAddSubmit_SetsAddressConfirmationTimestamp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with address_confirmed checked
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("email", "john@bairesdev.com")
	formData.Set("phone", "+1-555-0100")
	formData.Set("address", "123 Main St")
	formData.Set("address_confirmed", "on")

	req := httptest.NewRequest("POST", "/forms/software-engineers/add", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerAddSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer was created with address_confirmation_at set
	var addressConfirmationAt sql.NullTime
	err = db.QueryRowContext(ctx,
		`SELECT address_confirmation_at FROM software_engineers WHERE email = $1`,
		"john@bairesdev.com",
	).Scan(&addressConfirmationAt)
	if err != nil {
		t.Fatalf("Failed to retrieve created engineer: %v", err)
	}

	if !addressConfirmationAt.Valid {
		t.Error("Expected address_confirmation_at to be set when address_confirmed checkbox is checked")
	}
	if addressConfirmationAt.Time.IsZero() {
		t.Error("Expected address_confirmation_at to be a valid timestamp")
	}
}

// # RED: Test that address confirmation timestamp is set when checkbox changes from unchecked to checked in Edit form
func TestSoftwareEngineerEditSubmit_SetsAddressConfirmationTimestampWhenChecked(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create engineer with address_confirmed = false
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, phone, address, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		"Jane Smith", "jane@bairesdev.com", "+1-555-0200", "456 Oak St", false, nil, time.Now(), time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create test engineer: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with address_confirmed checked
	formData := url.Values{}
	formData.Set("name", "Jane Smith")
	formData.Set("email", "jane@bairesdev.com")
	formData.Set("phone", "+1-555-0200")
	formData.Set("address", "456 Oak St")
	formData.Set("address_confirmed", "on")

	req := httptest.NewRequest("POST", "/forms/software-engineers/"+strconv.FormatInt(engineerID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(engineerID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer was updated with address_confirmation_at set
	var addressConfirmationAt sql.NullTime
	err = db.QueryRowContext(ctx,
		`SELECT address_confirmation_at FROM software_engineers WHERE id = $1`,
		engineerID,
	).Scan(&addressConfirmationAt)
	if err != nil {
		t.Fatalf("Failed to retrieve updated engineer: %v", err)
	}

	if !addressConfirmationAt.Valid {
		t.Error("Expected address_confirmation_at to be set when address_confirmed checkbox changes from unchecked to checked")
	}
	if addressConfirmationAt.Time.IsZero() {
		t.Error("Expected address_confirmation_at to be a valid timestamp")
	}
}

// # RED: Test that address confirmation timestamp is cleared when checkbox is unchecked in Edit form
func TestSoftwareEngineerEditSubmit_ClearsAddressConfirmationTimestampWhenUnchecked(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create engineer with address_confirmed = true and a timestamp
	confirmationTime := time.Now().Add(-24 * time.Hour)
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, phone, address, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		"Bob Johnson", "bob@bairesdev.com", "+1-555-0300", "789 Pine St", true, confirmationTime, time.Now(), time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create test engineer: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with address_confirmed unchecked (not included in form)
	formData := url.Values{}
	formData.Set("name", "Bob Johnson")
	formData.Set("email", "bob@bairesdev.com")
	formData.Set("phone", "+1-555-0300")
	formData.Set("address", "789 Pine St")
	// address_confirmed is not set (unchecked)

	req := httptest.NewRequest("POST", "/forms/software-engineers/"+strconv.FormatInt(engineerID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(engineerID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer was updated with address_confirmation_at cleared
	var addressConfirmationAt sql.NullTime
	err = db.QueryRowContext(ctx,
		`SELECT address_confirmation_at FROM software_engineers WHERE id = $1`,
		engineerID,
	).Scan(&addressConfirmationAt)
	if err != nil {
		t.Fatalf("Failed to retrieve updated engineer: %v", err)
	}

	if addressConfirmationAt.Valid {
		t.Error("Expected address_confirmation_at to be cleared (NULL) when address_confirmed checkbox is unchecked")
	}
}

// # RED: Test that address confirmation timestamp is preserved when checkbox remains checked in Edit form
func TestSoftwareEngineerEditSubmit_PreservesAddressConfirmationTimestampWhenStillChecked(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create engineer with address_confirmed = true and a specific timestamp
	originalConfirmationTime := time.Now().Add(-48 * time.Hour)
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, phone, address, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		"Alice Brown", "alice@bairesdev.com", "+1-555-0400", "321 Elm St", true, originalConfirmationTime, time.Now(), time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create test engineer: %v", err)
	}

	// Reload the engineer to get the actual stored timestamp (accounting for any database timezone conversions)
	var storedConfirmationTime sql.NullTime
	err = db.QueryRowContext(ctx,
		`SELECT address_confirmation_at FROM software_engineers WHERE id = $1`,
		engineerID,
	).Scan(&storedConfirmationTime)
	if err != nil {
		t.Fatalf("Failed to retrieve stored timestamp: %v", err)
	}
	if !storedConfirmationTime.Valid {
		t.Fatalf("Expected stored timestamp to be valid")
	}
	originalConfirmationTime = storedConfirmationTime.Time

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with address_confirmed still checked
	formData := url.Values{}
	formData.Set("name", "Alice Brown")
	formData.Set("email", "alice@bairesdev.com")
	formData.Set("phone", "+1-555-0400")
	formData.Set("address", "321 Elm St")
	formData.Set("address_confirmed", "on")

	req := httptest.NewRequest("POST", "/forms/software-engineers/"+strconv.FormatInt(engineerID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(engineerID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer still has the original address_confirmation_at timestamp
	var addressConfirmationAt sql.NullTime
	err = db.QueryRowContext(ctx,
		`SELECT address_confirmation_at FROM software_engineers WHERE id = $1`,
		engineerID,
	).Scan(&addressConfirmationAt)
	if err != nil {
		t.Fatalf("Failed to retrieve updated engineer: %v", err)
	}

	if !addressConfirmationAt.Valid {
		t.Error("Expected address_confirmation_at to still be set when address_confirmed checkbox remains checked")
	}

	// The timestamp should be preserved (within 1 second tolerance for database precision)
	timeDiff := addressConfirmationAt.Time.Sub(originalConfirmationTime)
	if timeDiff < -time.Second || timeDiff > time.Second {
		t.Errorf("Expected address_confirmation_at to be preserved, but got difference of %v", timeDiff)
	}
}

// # RED: Test that software engineer can be created with international address format
func TestSoftwareEngineerAddSubmit_WithInternationalAddressFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with international address format
	formData := url.Values{}
	formData.Set("name", "Jane Smith")
	formData.Set("email", "jane@bairesdev.com")
	formData.Set("phone", "+1-555-0100")
	formData.Set("address_street", "456 Tech Avenue, Apt 12B")
	formData.Set("address_city", "Buenos Aires")
	formData.Set("address_country", "Argentina")
	formData.Set("address_state", "Buenos Aires")
	formData.Set("address_postal_code", "C1000ABC")

	req := httptest.NewRequest("POST", "/forms/software-engineers/add", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerAddSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer was created with international address fields
	var addressStreet, addressCity, addressCountry, addressState, addressPostalCode sql.NullString
	err = db.QueryRowContext(ctx,
		`SELECT address_street, address_city, address_country, address_state, address_postal_code 
		FROM software_engineers WHERE email = $1`,
		"jane@bairesdev.com",
	).Scan(&addressStreet, &addressCity, &addressCountry, &addressState, &addressPostalCode)
	if err != nil {
		t.Fatalf("Failed to retrieve created engineer: %v", err)
	}

	if !addressStreet.Valid || addressStreet.String != "456 Tech Avenue, Apt 12B" {
		t.Errorf("Expected address_street to be '456 Tech Avenue, Apt 12B', got %v", addressStreet)
	}
	if !addressCity.Valid || addressCity.String != "Buenos Aires" {
		t.Errorf("Expected address_city to be 'Buenos Aires', got %v", addressCity)
	}
	if !addressCountry.Valid || addressCountry.String != "Argentina" {
		t.Errorf("Expected address_country to be 'Argentina', got %v", addressCountry)
	}
	if !addressState.Valid || addressState.String != "Buenos Aires" {
		t.Errorf("Expected address_state to be 'Buenos Aires', got %v", addressState)
	}
	if !addressPostalCode.Valid || addressPostalCode.String != "C1000ABC" {
		t.Errorf("Expected address_postal_code to be 'C1000ABC', got %v", addressPostalCode)
	}
}

// # RED: Test that software engineer can be updated with international address format
func TestSoftwareEngineerEditSubmit_WithInternationalAddressFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create engineer with old address format
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, phone, address, address_confirmed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		"Bob Johnson", "bob@bairesdev.com", "+1-555-0200", "789 Old St", false, time.Now(), time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create test engineer: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with international address format
	formData := url.Values{}
	formData.Set("name", "Bob Johnson")
	formData.Set("email", "bob@bairesdev.com")
	formData.Set("phone", "+1-555-0200")
	formData.Set("address_street", "123 Main Street")
	formData.Set("address_city", "London")
	formData.Set("address_country", "United Kingdom")
	formData.Set("address_state", "")
	formData.Set("address_postal_code", "SW1A 1AA")

	req := httptest.NewRequest("POST", "/forms/software-engineers/"+strconv.FormatInt(engineerID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(engineerID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.SoftwareEngineerEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the engineer was updated with international address fields
	var addressStreet, addressCity, addressCountry, addressState, addressPostalCode sql.NullString
	err = db.QueryRowContext(ctx,
		`SELECT address_street, address_city, address_country, address_state, address_postal_code 
		FROM software_engineers WHERE id = $1`,
		engineerID,
	).Scan(&addressStreet, &addressCity, &addressCountry, &addressState, &addressPostalCode)
	if err != nil {
		t.Fatalf("Failed to retrieve updated engineer: %v", err)
	}

	if !addressStreet.Valid || addressStreet.String != "123 Main Street" {
		t.Errorf("Expected address_street to be '123 Main Street', got %v", addressStreet)
	}
	if !addressCity.Valid || addressCity.String != "London" {
		t.Errorf("Expected address_city to be 'London', got %v", addressCity)
	}
	if !addressCountry.Valid || addressCountry.String != "United Kingdom" {
		t.Errorf("Expected address_country to be 'United Kingdom', got %v", addressCountry)
	}
	if addressState.Valid && addressState.String != "" {
		t.Errorf("Expected address_state to be empty, got %v", addressState)
	}
	if !addressPostalCode.Valid || addressPostalCode.String != "SW1A 1AA" {
		t.Errorf("Expected address_postal_code to be 'SW1A 1AA', got %v", addressPostalCode)
	}
}
