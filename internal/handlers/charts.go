package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// ChartsHandler handles chart data API requests
type ChartsHandler struct {
	DB *sql.DB
}

// NewChartsHandler creates a new ChartsHandler
func NewChartsHandler(db *sql.DB) *ChartsHandler {
	return &ChartsHandler{
		DB: db,
	}
}

// ShipmentsOverTimeAPI returns shipment count data over time for line charts
func (h *ChartsHandler) ShipmentsOverTimeAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get days parameter from query string (default: 30 days)
	days := 30
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	// Get chart data
	data, err := models.GetShipmentsOverTime(h.DB, days)
	if err != nil {
		log.Printf("Error getting shipments over time: %v", err)
		http.Error(w, "Failed to fetch chart data", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// StatusDistributionAPI returns shipment status distribution for pie charts
func (h *ChartsHandler) StatusDistributionAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get distribution data
	data, err := models.GetShipmentStatusDistribution(h.DB)
	if err != nil {
		log.Printf("Error getting status distribution: %v", err)
		http.Error(w, "Failed to fetch chart data", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeliveryTimeTrendsAPI returns delivery time trends for bar charts
func (h *ChartsHandler) DeliveryTimeTrendsAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get weeks parameter from query string (default: 8 weeks)
	weeks := 8
	if weeksParam := r.URL.Query().Get("weeks"); weeksParam != "" {
		if parsedWeeks, err := strconv.Atoi(weeksParam); err == nil && parsedWeeks > 0 {
			weeks = parsedWeeks
		}
	}

	// Get trends data
	data, err := models.GetDeliveryTimeTrends(h.DB, weeks)
	if err != nil {
		log.Printf("Error getting delivery time trends: %v", err)
		http.Error(w, "Failed to fetch chart data", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

