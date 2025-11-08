package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(db *sql.DB, templates *template.Template) *DashboardHandler {
	return &DashboardHandler{
		DB:        db,
		Templates: templates,
	}
}

// Dashboard displays the main dashboard with statistics
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics and project manager users can access the dashboard
	if user.Role != models.RoleLogistics && user.Role != models.RoleProjectManager {
		http.Error(w, "Forbidden: Only logistics and project manager users can access this page", http.StatusForbidden)
		return
	}

	// Get dashboard statistics
	stats, err := models.GetDashboardStats(h.DB)
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		http.Error(w, "Failed to load dashboard statistics", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":        user,
		"Stats":       stats,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "dashboard",
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "dashboard-with-charts.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
		return
	}
}

