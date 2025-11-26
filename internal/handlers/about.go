package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// AboutHandler handles about page requests
type AboutHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewAboutHandler creates a new AboutHandler
func NewAboutHandler(db *sql.DB, templates *template.Template) *AboutHandler {
	return &AboutHandler{
		DB:        db,
		Templates: templates,
	}
}

// About displays the about page
func (h *AboutHandler) About(w http.ResponseWriter, r *http.Request) {
	// Get user from context (may be nil if not authenticated)
	user := middleware.GetUserFromContext(r.Context())

	// Prepare template data
	var nav views.NavigationLinks
	if user != nil {
		nav = views.GetNavigationLinks(user.Role)
	}

	data := map[string]interface{}{
		"User":        user,
		"Nav":         nav,
		"CurrentPage": "about",
	}

	// Execute template
	if err := h.Templates.ExecuteTemplate(w, "about.html", data); err != nil {
		http.Error(w, "Failed to render about page", http.StatusInternalServerError)
		return
	}
}

