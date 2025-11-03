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

	// Set content type before any writes
	w.Header().Set("Content-Type", "application/json")

	// Get chart data
	data, err := models.GetShipmentsOverTime(h.DB, days)
	if err != nil {
		log.Printf("Error getting shipments over time: %v", err)
		// Return empty array instead of error to keep frontend happy
		json.NewEncoder(w).Encode([]models.ChartDataPoint{})
		return
	}

	// Ensure we always return an array, even if empty
	if data == nil {
		data = []models.ChartDataPoint{}
	}

	// Return JSON response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		json.NewEncoder(w).Encode([]models.ChartDataPoint{})
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

	// Set content type before any writes
	w.Header().Set("Content-Type", "application/json")

	// Get distribution data
	data, err := models.GetShipmentStatusDistribution(h.DB)
	if err != nil {
		log.Printf("Error getting status distribution: %v", err)
		// Return empty array instead of error to keep frontend happy
		json.NewEncoder(w).Encode([]models.StatusDistribution{})
		return
	}

	// Ensure we always return an array, even if empty
	if data == nil {
		data = []models.StatusDistribution{}
	}

	// Return JSON response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		json.NewEncoder(w).Encode([]models.StatusDistribution{})
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

	// Set content type before any writes
	w.Header().Set("Content-Type", "application/json")

	// Get trends data
	data, err := models.GetDeliveryTimeTrends(h.DB, weeks)
	if err != nil {
		log.Printf("Error getting delivery time trends: %v", err)
		// Return empty array instead of error to keep frontend happy
		json.NewEncoder(w).Encode([]models.DeliveryTimeTrend{})
		return
	}

	// Ensure we always return an array, even if empty
	if data == nil {
		data = []models.DeliveryTimeTrend{}
	}

	// Return JSON response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		json.NewEncoder(w).Encode([]models.DeliveryTimeTrend{})
		return
	}
}

