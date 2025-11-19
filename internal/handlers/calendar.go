package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// CalendarHandler handles calendar-related requests
type CalendarHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewCalendarHandler creates a new CalendarHandler
func NewCalendarHandler(db *sql.DB, templates *template.Template) *CalendarHandler {
	return &CalendarHandler{
		DB:        db,
		Templates: templates,
	}
}

// Calendar displays the calendar view with shipment events
func (h *CalendarHandler) Calendar(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse query parameters for date range
	startDateStr := r.URL.Query().Get("start")
	endDateStr := r.URL.Query().Get("end")

	// Default to current month if no dates provided
	var startDate, endDate time.Time
	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		var err error
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			startDate = time.Now().AddDate(0, 0, -30)
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			endDate = time.Now().AddDate(0, 0, 30)
		}
	}

	// Get calendar events
	// For client users, filter by their company ID; for other roles, show all events
	var clientCompanyID *int64
	if user.Role == models.RoleClient && user.ClientCompanyID != nil {
		clientCompanyID = user.ClientCompanyID
	}

	events, err := models.GetCalendarEvents(h.DB, startDate, endDate, clientCompanyID)
	if err != nil {
		log.Printf("Error getting calendar events: %v", err)
		http.Error(w, "Failed to load calendar events", http.StatusInternalServerError)
		return
	}

	// Generate calendar grid with events
	calendarGrid := models.GenerateCalendarGridWithEvents(
		startDate.Year(),
		startDate.Month(),
		events,
	)

	// Prepare template data
	data := map[string]interface{}{
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "calendar",
		"Events":       events,
		"CalendarGrid": calendarGrid,
		"StartDate":    startDate,
		"EndDate":      endDate,
		"CurrentMonth": startDate.Format("January 2006"),
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "calendar.html", data); err != nil {
		log.Printf("Error executing calendar template: %v", err)
		http.Error(w, "Failed to render calendar", http.StatusInternalServerError)
		return
	}
}

