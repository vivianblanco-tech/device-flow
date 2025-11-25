package views

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestNavbarComponentRendering tests that the navbar component renders correctly
func TestNavbarComponentRendering(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Define template functions needed for navbar
	funcMap := template.FuncMap{
		"title": func(v interface{}) string {
			// Convert interface{} to string
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			default:
				s = ""
			}
			return strings.Title(s)
		},
	}

	tests := []struct {
		name                string
		userRole            models.UserRole
		currentPage         string
		expectedLinks       []string
		unexpectedLinks     []string
		expectedStickyClass bool
	}{
		{
			name:        "logistics user sees all navigation links",
			userRole:    models.RoleLogistics,
			currentPage: "dashboard",
			expectedLinks: []string{
				`href="/dashboard"`,
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/pickup-forms"`,
				`href="/reception-reports"`,
			},
			unexpectedLinks:     []string{},
			expectedStickyClass: true,
		},
		{
			name:        "project manager sees dashboard and reports",
			userRole:    models.RoleProjectManager,
			currentPage: "dashboard",
			expectedLinks: []string{
				`href="/dashboard"`,
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/reports"`, // Project Manager should have access to Reports
			},
			unexpectedLinks: []string{
				`href="/pickup-forms"`,
				`href="/reception-reports"`,
			},
			expectedStickyClass: true,
		},
		{
			name:        "warehouse user sees inventory and reception",
			userRole:    models.RoleWarehouse,
			currentPage: "inventory",
			expectedLinks: []string{
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/reception-reports"`,
			},
			unexpectedLinks: []string{
				`href="/dashboard"`,
				`href="/pickup-forms"`,
			},
			expectedStickyClass: true,
		},
		{
			name:        "client user has limited access",
			userRole:    models.RoleClient,
			currentPage: "shipments",
			expectedLinks: []string{
				`href="/shipments"`,
				`href="/inventory"`, // Client users can now view their company's inventory
				`href="/calendar"`,
				`href="/pickup-forms"`,
			},
			unexpectedLinks: []string{
				`href="/dashboard"`,
				`href="/reception-reports"`,
			},
			expectedStickyClass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the navbar template
			tmpl, err := template.New("navbar.html").Funcs(funcMap).ParseFiles("../../templates/components/navbar.html")
			if err != nil {
				t.Fatalf("failed to parse navbar template: %v", err)
			}

			// Prepare template data
			user := &models.User{
				ID:    1,
				Email: "test@bairesdev.com",
				Role:  tt.userRole,
			}

			nav := GetNavigationLinks(tt.userRole)

			data := map[string]interface{}{
				"User":        user,
				"Nav":         nav,
				"CurrentPage": tt.currentPage,
			}

			// Render the template
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				t.Fatalf("failed to execute navbar template: %v", err)
			}

			html := buf.String()

			// Check for sticky positioning classes
			if tt.expectedStickyClass {
				if !strings.Contains(html, "sticky") && !strings.Contains(html, "fixed") {
					t.Errorf("expected sticky positioning class in navbar, but not found")
				}
			}

			// Check expected links are present
			for _, link := range tt.expectedLinks {
				if !strings.Contains(html, link) {
					t.Errorf("expected link %s in navbar HTML for %s role, but it was not found", link, tt.userRole)
				}
			}

			// Check unexpected links are NOT present
			for _, link := range tt.unexpectedLinks {
				if strings.Contains(html, link) {
					t.Errorf("did not expect link %s in navbar HTML for %s role, but it was found", link, tt.userRole)
				}
			}

			// Check user info is displayed
			if !strings.Contains(html, user.Email) {
				t.Errorf("expected user email %s in navbar, but not found", user.Email)
			}

			// Check logout link is present
			if !strings.Contains(html, `href="/logout"`) {
				t.Error("expected logout link in navbar, but not found")
			}

			// Check logo/title is present
			if !strings.Contains(html, "Laptop Tracking System") {
				t.Error("expected application title in navbar, but not found")
			}
		})
	}
}

// TestNavbarActiveLink tests that the current page is highlighted correctly
func TestNavbarActiveLink(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	funcMap := template.FuncMap{
		"title": func(v interface{}) string {
			// Convert interface{} to string
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			default:
				s = ""
			}
			return strings.Title(s)
		},
	}

	tests := []struct {
		name            string
		currentPage     string
		expectedActive  string
		expectedClasses []string
	}{
		{
			name:           "dashboard page is active",
			currentPage:    "dashboard",
			expectedActive: "dashboard",
			expectedClasses: []string{
				"text-blue-600",
				"font-medium",
			},
		},
		{
			name:           "inventory page is active",
			currentPage:    "inventory",
			expectedActive: "inventory",
			expectedClasses: []string{
				"text-blue-600",
				"font-medium",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the navbar template
			tmpl, err := template.New("navbar.html").Funcs(funcMap).ParseFiles("../../templates/components/navbar.html")
			if err != nil {
				t.Fatalf("failed to parse navbar template: %v", err)
			}

			// Prepare template data with logistics user (has all links)
			user := &models.User{
				ID:    1,
				Email: "test@bairesdev.com",
				Role:  models.RoleLogistics,
			}

			nav := GetNavigationLinks(models.RoleLogistics)

			data := map[string]interface{}{
				"User":        user,
				"Nav":         nav,
				"CurrentPage": tt.currentPage,
			}

			// Render the template
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				t.Fatalf("failed to execute navbar template: %v", err)
			}

			html := buf.String()

			// Check that active classes are present
			for _, class := range tt.expectedClasses {
				if !strings.Contains(html, class) {
					t.Errorf("expected active class %s in navbar for page %s, but not found", class, tt.currentPage)
				}
			}
		})
	}
}

// TestNavbarClientUserCompanyName tests that client users with company_id display company name instead of email
func TestNavbarClientUserCompanyName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Define template functions needed for navbar
	funcMap := template.FuncMap{
		"title": func(v interface{}) string {
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			default:
				s = ""
			}
			return strings.Title(s)
		},
	}

	// Load the navbar template
	tmpl, err := template.New("navbar.html").Funcs(funcMap).ParseFiles("../../templates/components/navbar.html")
	if err != nil {
		t.Fatalf("failed to parse navbar template: %v", err)
	}

	companyID := int64(1)
	companyName := "Apex Digital Group"
	userEmail := "procurement@apexdigital.com"

	tests := []struct {
		name                string
		user                *models.User
		expectedDisplayText string
		shouldNotContain    string
	}{
		{
			name: "client user with company_id displays company name",
			user: &models.User{
				ID:                1,
				Email:             userEmail,
				Role:              models.RoleClient,
				ClientCompanyID:   &companyID,
				ClientCompanyName: companyName,
			},
			expectedDisplayText: companyName,
			shouldNotContain:    userEmail,
		},
		{
			name: "client user without company_id displays email",
			user: &models.User{
				ID:              2,
				Email:           userEmail,
				Role:            models.RoleClient,
				ClientCompanyID: nil,
			},
			expectedDisplayText: userEmail,
			shouldNotContain:    "",
		},
		{
			name: "non-client user always displays email",
			user: &models.User{
				ID:              3,
				Email:           "logistics@bairesdev.com",
				Role:            models.RoleLogistics,
				ClientCompanyID: nil,
			},
			expectedDisplayText: "logistics@bairesdev.com",
			shouldNotContain:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := GetNavigationLinks(tt.user.Role)

			data := map[string]interface{}{
				"User":        tt.user,
				"Nav":         nav,
				"CurrentPage": "shipments",
			}

			// Render the template
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				t.Fatalf("failed to execute navbar template: %v", err)
			}

			html := buf.String()

			// Check that expected text is displayed
			if !strings.Contains(html, tt.expectedDisplayText) {
				t.Errorf("expected navbar to contain %s, but it was not found. HTML: %s", tt.expectedDisplayText, html)
			}

			// Check that email is NOT displayed for client users with company name
			if tt.shouldNotContain != "" && strings.Contains(html, tt.shouldNotContain) {
				t.Errorf("expected navbar NOT to contain %s for client user with company, but it was found. HTML: %s", tt.shouldNotContain, html)
			}

			// Check logout link is present
			if !strings.Contains(html, `href="/logout"`) {
				t.Error("expected logout link in navbar, but not found")
			}
		})
	}
}
