package handlers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// ReportsHandler handles report-related requests
type ReportsHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewReportsHandler creates a new ReportsHandler
func NewReportsHandler(db *sql.DB, templates *template.Template) *ReportsHandler {
	return &ReportsHandler{
		DB:        db,
		Templates: templates,
	}
}

// requireClientRole checks if the user is a client user
func (h *ReportsHandler) requireClientRole(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, false
	}
	if user.Role != models.RoleClient {
		http.Error(w, "Forbidden: Only client users can access reports", http.StatusForbidden)
		return nil, false
	}
	if user.ClientCompanyID == nil {
		http.Error(w, "Forbidden: Client user must be associated with a company", http.StatusForbidden)
		return nil, false
	}
	return user, true
}

// ReportsIndex displays the reports index page
func (h *ReportsHandler) ReportsIndex(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireClientRole(w, r)
	if !ok {
		return
	}

	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "reports",
	}

	if err := h.Templates.ExecuteTemplate(w, "reports-index.html", data); err != nil {
		log.Printf("Error executing reports index template: %v", err)
		http.Error(w, "Failed to render reports index", http.StatusInternalServerError)
		return
	}
}

// ShipmentStatusDashboard displays the shipment status dashboard report
func (h *ReportsHandler) ShipmentStatusDashboard(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireClientRole(w, r)
	if !ok {
		return
	}

	format := r.URL.Query().Get("format")

	// Get report data
	reportData, err := h.getShipmentStatusData(user.ClientCompanyID)
	if err != nil {
		log.Printf("Error getting shipment status data: %v", err)
		http.Error(w, "Failed to load report data", http.StatusInternalServerError)
		return
	}

	// Handle exports
	switch format {
	case "csv":
		h.exportShipmentStatusCSV(w, reportData)
		return
	case "xlsx":
		h.exportShipmentStatusExcel(w, reportData)
		return
	case "pdf":
		h.exportShipmentStatusPDF(w, reportData)
		return
	}

	// HTML view
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "reports",
		"ReportData":  reportData,
		"ReportType":  "shipment-status",
	}

	if err := h.Templates.ExecuteTemplate(w, "report-shipment-status.html", data); err != nil {
		log.Printf("Error executing shipment status report template: %v", err)
		http.Error(w, "Failed to render report", http.StatusInternalServerError)
		return
	}
}

// InventorySummaryReport displays the inventory summary report
func (h *ReportsHandler) InventorySummaryReport(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireClientRole(w, r)
	if !ok {
		return
	}

	format := r.URL.Query().Get("format")

	// Get report data
	reportData, err := h.getInventorySummaryData(user.ClientCompanyID)
	if err != nil {
		log.Printf("Error getting inventory summary data: %v", err)
		http.Error(w, "Failed to load report data", http.StatusInternalServerError)
		return
	}

	// Handle exports
	switch format {
	case "csv":
		h.exportInventorySummaryCSV(w, reportData)
		return
	case "xlsx":
		h.exportInventorySummaryExcel(w, reportData)
		return
	case "pdf":
		h.exportInventorySummaryPDF(w, reportData)
		return
	}

	// HTML view
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "reports",
		"ReportData":  reportData,
		"ReportType":  "inventory-summary",
	}

	if err := h.Templates.ExecuteTemplate(w, "report-inventory-summary.html", data); err != nil {
		log.Printf("Error executing inventory summary report template: %v", err)
		http.Error(w, "Failed to render report", http.StatusInternalServerError)
		return
	}
}

// ShipmentTimelineReport displays the shipment timeline report
func (h *ReportsHandler) ShipmentTimelineReport(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireClientRole(w, r)
	if !ok {
		return
	}

	format := r.URL.Query().Get("format")

	// Get report data
	reportData, err := h.getShipmentTimelineData(user.ClientCompanyID)
	if err != nil {
		log.Printf("Error getting shipment timeline data: %v", err)
		http.Error(w, "Failed to load report data", http.StatusInternalServerError)
		return
	}

	// Handle exports
	switch format {
	case "csv":
		h.exportShipmentTimelineCSV(w, reportData)
		return
	case "xlsx":
		h.exportShipmentTimelineExcel(w, reportData)
		return
	case "pdf":
		h.exportShipmentTimelinePDF(w, reportData)
		return
	}

	// HTML view
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "reports",
		"ReportData":  reportData,
		"ReportType":  "shipment-timeline",
	}

	if err := h.Templates.ExecuteTemplate(w, "report-shipment-timeline.html", data); err != nil {
		log.Printf("Error executing shipment timeline report template: %v", err)
		http.Error(w, "Failed to render report", http.StatusInternalServerError)
		return
	}
}

// ShipmentStatusData represents data for shipment status dashboard
type ShipmentStatusData struct {
	TotalShipments      int
	TotalThisMonth      int
	TotalThisYear       int
	ByStatus            map[string]int
	ByType              map[string]int
	PendingPickups      int
	InTransitCount      int
	DeliveredCount      int
	AverageDeliveryTime float64
	Shipments           []ShipmentStatusRow
}

// ShipmentStatusRow represents a row in the shipment status table
type ShipmentStatusRow struct {
	ID                int64
	JiraTicket        string
	Type              string
	Status            string
	LaptopCount       int
	CourierName       string
	TrackingNumber    string
	CreatedAt         time.Time
	DeliveredAt       *time.Time
	DaysSinceCreated  int
	DaysSinceDelivery *int
}

// getShipmentStatusData retrieves shipment status data for a client company
func (h *ReportsHandler) getShipmentStatusData(companyID *int64) (*ShipmentStatusData, error) {
	if companyID == nil {
		return nil, fmt.Errorf("company ID is required")
	}

	data := &ShipmentStatusData{
		ByStatus: make(map[string]int),
		ByType:   make(map[string]int),
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// Get total shipments
	err := h.DB.QueryRow(
		`SELECT COUNT(*) FROM shipments WHERE client_company_id = $1`,
		*companyID,
	).Scan(&data.TotalShipments)
	if err != nil {
		return nil, fmt.Errorf("failed to get total shipments: %w", err)
	}

	// Get this month's shipments
	err = h.DB.QueryRow(
		`SELECT COUNT(*) FROM shipments WHERE client_company_id = $1 AND created_at >= $2`,
		*companyID, startOfMonth,
	).Scan(&data.TotalThisMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly shipments: %w", err)
	}

	// Get this year's shipments
	err = h.DB.QueryRow(
		`SELECT COUNT(*) FROM shipments WHERE client_company_id = $1 AND created_at >= $2`,
		*companyID, startOfYear,
	).Scan(&data.TotalThisYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly shipments: %w", err)
	}

	// Get shipments by status
	rows, err := h.DB.Query(
		`SELECT status, COUNT(*) FROM shipments WHERE client_company_id = $1 GROUP BY status`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipments by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		data.ByStatus[status] = count
	}

	// Get shipments by type
	rows, err = h.DB.Query(
		`SELECT shipment_type, COUNT(*) FROM shipments WHERE client_company_id = $1 GROUP BY shipment_type`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipments by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var shipmentType string
		var count int
		if err := rows.Scan(&shipmentType, &count); err != nil {
			continue
		}
		data.ByType[shipmentType] = count
	}

	// Calculate pending pickups
	data.PendingPickups = data.ByStatus[string(models.ShipmentStatusPendingPickup)] +
		data.ByStatus[string(models.ShipmentStatusPickupScheduled)]

	// Calculate in-transit count
	data.InTransitCount = data.ByStatus[string(models.ShipmentStatusInTransitToWarehouse)] +
		data.ByStatus[string(models.ShipmentStatusInTransitToEngineer)]

	data.DeliveredCount = data.ByStatus[string(models.ShipmentStatusDelivered)]

	// Get average delivery time for delivered shipments
	var avgDays sql.NullFloat64
	err = h.DB.QueryRow(
		`SELECT AVG(EXTRACT(EPOCH FROM (delivered_at - created_at)) / 86400)
		 FROM shipments 
		 WHERE client_company_id = $1 AND delivered_at IS NOT NULL`,
		*companyID,
	).Scan(&avgDays)
	if err == nil && avgDays.Valid {
		data.AverageDeliveryTime = avgDays.Float64
	}

	// Get detailed shipment list
	rows, err = h.DB.Query(
		`SELECT id, jira_ticket_number, shipment_type, status, laptop_count, 
		 courier_name, tracking_number, created_at, delivered_at
		 FROM shipments 
		 WHERE client_company_id = $1 
		 ORDER BY created_at DESC`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment list: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row ShipmentStatusRow
		var deliveredAt sql.NullTime
		err := rows.Scan(
			&row.ID, &row.JiraTicket, &row.Type, &row.Status, &row.LaptopCount,
			&row.CourierName, &row.TrackingNumber, &row.CreatedAt, &deliveredAt,
		)
		if err != nil {
			continue
		}

		if deliveredAt.Valid {
			row.DeliveredAt = &deliveredAt.Time
			days := int(now.Sub(deliveredAt.Time).Hours() / 24)
			row.DaysSinceDelivery = &days
		}

		row.DaysSinceCreated = int(now.Sub(row.CreatedAt).Hours() / 24)
		data.Shipments = append(data.Shipments, row)
	}

	return data, nil
}

// InventorySummaryData represents data for inventory summary report
type InventorySummaryData struct {
	TotalLaptops       int
	ByStatus           map[string]int
	ByBrand            map[string]int
	AssignedCount      int
	UnassignedCount    int
	Laptops            []InventorySummaryRow
}

// InventorySummaryRow represents a row in the inventory summary table
type InventorySummaryRow struct {
	ID                 int64
	SerialNumber       string
	Brand              string
	Model              string
	CPU                string
	RAMGB              string
	SSDGB              string
	Status             string
	SoftwareEngineerID *int64
	SoftwareEngineerName string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// getInventorySummaryData retrieves inventory summary data for a client company
func (h *ReportsHandler) getInventorySummaryData(companyID *int64) (*InventorySummaryData, error) {
	if companyID == nil {
		return nil, fmt.Errorf("company ID is required")
	}

	data := &InventorySummaryData{
		ByStatus: make(map[string]int),
		ByBrand:  make(map[string]int),
	}

	// Get total laptops
	err := h.DB.QueryRow(
		`SELECT COUNT(*) FROM laptops WHERE client_company_id = $1`,
		*companyID,
	).Scan(&data.TotalLaptops)
	if err != nil {
		return nil, fmt.Errorf("failed to get total laptops: %w", err)
	}

	// Get laptops by status
	rows, err := h.DB.Query(
		`SELECT status, COUNT(*) FROM laptops WHERE client_company_id = $1 GROUP BY status`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get laptops by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		data.ByStatus[status] = count
	}

	// Get laptops by brand
	rows, err = h.DB.Query(
		`SELECT brand, COUNT(*) FROM laptops WHERE client_company_id = $1 AND brand IS NOT NULL GROUP BY brand`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get laptops by brand: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var brand string
		var count int
		if err := rows.Scan(&brand, &count); err != nil {
			continue
		}
		data.ByBrand[brand] = count
	}

	// Get assigned vs unassigned count
	err = h.DB.QueryRow(
		`SELECT COUNT(*) FROM laptops WHERE client_company_id = $1 AND software_engineer_id IS NOT NULL`,
		*companyID,
	).Scan(&data.AssignedCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned count: %w", err)
	}

	data.UnassignedCount = data.TotalLaptops - data.AssignedCount

	// Get detailed laptop list
	rows, err = h.DB.Query(
		`SELECT l.id, l.serial_number, l.brand, l.model, l.cpu, l.ram_gb, l.ssd_gb, 
		 l.status, l.software_engineer_id, se.name as engineer_name,
		 l.created_at, l.updated_at
		 FROM laptops l
		 LEFT JOIN software_engineers se ON se.id = l.software_engineer_id
		 WHERE l.client_company_id = $1 
		 ORDER BY l.created_at DESC`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get laptop list: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row InventorySummaryRow
		var engineerID sql.NullInt64
		var engineerName sql.NullString
		err := rows.Scan(
			&row.ID, &row.SerialNumber, &row.Brand, &row.Model, &row.CPU, &row.RAMGB, &row.SSDGB,
			&row.Status, &engineerID, &engineerName,
			&row.CreatedAt, &row.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if engineerID.Valid {
			row.SoftwareEngineerID = &engineerID.Int64
		}
		if engineerName.Valid {
			row.SoftwareEngineerName = engineerName.String
		}

		data.Laptops = append(data.Laptops, row)
	}

	return data, nil
}

// ShipmentTimelineData represents data for shipment timeline report
type ShipmentTimelineData struct {
	Shipments []ShipmentTimelineRow
}

// ShipmentTimelineRow represents a row in the shipment timeline table
type ShipmentTimelineRow struct {
	ID                  int64
	JiraTicket          string
	Type                string
	Status              string
	LaptopCount         int
	CourierName         string
	TrackingNumber      string
	PickupScheduledDate *time.Time
	PickedUpAt          *time.Time
	ArrivedWarehouseAt  *time.Time
	ReleasedWarehouseAt *time.Time
	ETAToEngineer       *time.Time
	DeliveredAt         *time.Time
	CreatedAt           time.Time
	TotalDurationDays   *int
	TimeBetweenStages   map[string]int // Days between stages
}

// getShipmentTimelineData retrieves shipment timeline data for a client company
func (h *ReportsHandler) getShipmentTimelineData(companyID *int64) (*ShipmentTimelineData, error) {
	if companyID == nil {
		return nil, fmt.Errorf("company ID is required")
	}

	data := &ShipmentTimelineData{}

	rows, err := h.DB.Query(
		`SELECT id, jira_ticket_number, shipment_type, status, laptop_count,
		 courier_name, tracking_number, pickup_scheduled_date, picked_up_at,
		 arrived_warehouse_at, released_warehouse_at, eta_to_engineer, delivered_at,
		 created_at
		 FROM shipments 
		 WHERE client_company_id = $1 
		 ORDER BY created_at DESC`,
		*companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment timeline data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row ShipmentTimelineRow
		var pickupScheduled sql.NullTime
		var pickedUp sql.NullTime
		var arrivedWarehouse sql.NullTime
		var releasedWarehouse sql.NullTime
		var etaToEngineer sql.NullTime
		var deliveredAt sql.NullTime

		err := rows.Scan(
			&row.ID, &row.JiraTicket, &row.Type, &row.Status, &row.LaptopCount,
			&row.CourierName, &row.TrackingNumber, &pickupScheduled, &pickedUp,
			&arrivedWarehouse, &releasedWarehouse, &etaToEngineer, &deliveredAt,
			&row.CreatedAt,
		)
		if err != nil {
			continue
		}

		if pickupScheduled.Valid {
			row.PickupScheduledDate = &pickupScheduled.Time
		}
		if pickedUp.Valid {
			row.PickedUpAt = &pickedUp.Time
		}
		if arrivedWarehouse.Valid {
			row.ArrivedWarehouseAt = &arrivedWarehouse.Time
		}
		if releasedWarehouse.Valid {
			row.ReleasedWarehouseAt = &releasedWarehouse.Time
		}
		if etaToEngineer.Valid {
			row.ETAToEngineer = &etaToEngineer.Time
		}
		if deliveredAt.Valid {
			row.DeliveredAt = &deliveredAt.Time
			duration := int(deliveredAt.Time.Sub(row.CreatedAt).Hours() / 24)
			row.TotalDurationDays = &duration
		}

		// Calculate time between stages
		row.TimeBetweenStages = make(map[string]int)
		if row.PickupScheduledDate != nil && row.PickedUpAt != nil {
			days := int(row.PickedUpAt.Sub(*row.PickupScheduledDate).Hours() / 24)
			row.TimeBetweenStages["pickup_to_picked"] = days
		}
		if row.PickedUpAt != nil && row.ArrivedWarehouseAt != nil {
			days := int(row.ArrivedWarehouseAt.Sub(*row.PickedUpAt).Hours() / 24)
			row.TimeBetweenStages["picked_to_warehouse"] = days
		}
		if row.ArrivedWarehouseAt != nil && row.ReleasedWarehouseAt != nil {
			days := int(row.ReleasedWarehouseAt.Sub(*row.ArrivedWarehouseAt).Hours() / 24)
			row.TimeBetweenStages["warehouse_processing"] = days
		}
		if row.ReleasedWarehouseAt != nil && row.DeliveredAt != nil {
			days := int(row.DeliveredAt.Sub(*row.ReleasedWarehouseAt).Hours() / 24)
			row.TimeBetweenStages["warehouse_to_delivered"] = days
		}

		data.Shipments = append(data.Shipments, row)
	}

	return data, nil
}

// Export functions will be implemented next...
// Placeholder functions to satisfy compilation
func (h *ReportsHandler) exportShipmentStatusCSV(w http.ResponseWriter, data *ShipmentStatusData) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-status-%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "JIRA Ticket", "Type", "Status", "Laptop Count", "Courier", "Tracking Number", "Created At", "Delivered At", "Days Since Created"})

	// Write data
	for _, shipment := range data.Shipments {
		deliveredAt := ""
		if shipment.DeliveredAt != nil {
			deliveredAt = shipment.DeliveredAt.Format("2006-01-02 15:04:05")
		}
		writer.Write([]string{
			strconv.FormatInt(shipment.ID, 10),
			shipment.JiraTicket,
			shipment.Type,
			shipment.Status,
			strconv.Itoa(shipment.LaptopCount),
			shipment.CourierName,
			shipment.TrackingNumber,
			shipment.CreatedAt.Format("2006-01-02 15:04:05"),
			deliveredAt,
			strconv.Itoa(shipment.DaysSinceCreated),
		})
	}
}

func (h *ReportsHandler) exportShipmentStatusExcel(w http.ResponseWriter, data *ShipmentStatusData) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Shipment Status"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := []string{"ID", "JIRA Ticket", "Type", "Status", "Laptop Count", "Courier", "Tracking Number", "Created At", "Delivered At", "Days Since Created"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, shipment := range data.Shipments {
		row := i + 2
		deliveredAt := ""
		if shipment.DeliveredAt != nil {
			deliveredAt = shipment.DeliveredAt.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), shipment.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), shipment.JiraTicket)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), shipment.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), shipment.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), shipment.LaptopCount)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), shipment.CourierName)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), shipment.TrackingNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), shipment.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), deliveredAt)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), shipment.DaysSinceCreated)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-status-%s.xlsx", time.Now().Format("2006-01-02")))
	f.Write(w)
}

func (h *ReportsHandler) exportShipmentStatusPDF(w http.ResponseWriter, data *ShipmentStatusData) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Shipment Status Report")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, fmt.Sprintf("Total Shipments: %d", data.TotalShipments))
	pdf.Ln(10)

	// Table headers
	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "JIRA", "Type", "Status", "Count", "Courier", "Tracking", "Created"}
	colWidth := 35.0
	for i, header := range headers {
		pdf.Cell(colWidth, 7, header)
		if i < len(headers)-1 {
			pdf.Cell(colWidth, 7, "")
		}
	}
	pdf.Ln(7)

	// Table data
	pdf.SetFont("Arial", "", 8)
	for _, shipment := range data.Shipments {
		pdf.Cell(colWidth, 6, strconv.FormatInt(shipment.ID, 10))
		pdf.Cell(colWidth, 6, shipment.JiraTicket)
		pdf.Cell(colWidth, 6, shipment.Type)
		pdf.Cell(colWidth, 6, shipment.Status)
		pdf.Cell(colWidth, 6, strconv.Itoa(shipment.LaptopCount))
		pdf.Cell(colWidth, 6, shipment.CourierName)
		pdf.Cell(colWidth, 6, shipment.TrackingNumber)
		pdf.Cell(colWidth, 6, shipment.CreatedAt.Format("2006-01-02"))
		pdf.Ln(6)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-status-%s.pdf", time.Now().Format("2006-01-02")))
	pdf.Output(w)
}

func (h *ReportsHandler) exportInventorySummaryCSV(w http.ResponseWriter, data *InventorySummaryData) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=inventory-summary-%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Serial Number", "Brand", "Model", "CPU", "RAM", "SSD", "Status", "Engineer", "Created At"})

	for _, laptop := range data.Laptops {
		writer.Write([]string{
			strconv.FormatInt(laptop.ID, 10),
			laptop.SerialNumber,
			laptop.Brand,
			laptop.Model,
			laptop.CPU,
			laptop.RAMGB,
			laptop.SSDGB,
			laptop.Status,
			laptop.SoftwareEngineerName,
			laptop.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

func (h *ReportsHandler) exportInventorySummaryExcel(w http.ResponseWriter, data *InventorySummaryData) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Inventory Summary"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Serial Number", "Brand", "Model", "CPU", "RAM", "SSD", "Status", "Engineer", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, laptop := range data.Laptops {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), laptop.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), laptop.SerialNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), laptop.Brand)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), laptop.Model)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), laptop.CPU)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), laptop.RAMGB)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), laptop.SSDGB)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), laptop.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), laptop.SoftwareEngineerName)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), laptop.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=inventory-summary-%s.xlsx", time.Now().Format("2006-01-02")))
	f.Write(w)
}

func (h *ReportsHandler) exportInventorySummaryPDF(w http.ResponseWriter, data *InventorySummaryData) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Inventory Summary Report")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, fmt.Sprintf("Total Laptops: %d", data.TotalLaptops))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "Serial", "Brand", "Model", "Status", "Engineer"}
	colWidth := 40.0
	for i, header := range headers {
		pdf.Cell(colWidth, 7, header)
		if i < len(headers)-1 {
			pdf.Cell(colWidth, 7, "")
		}
	}
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 8)
	for _, laptop := range data.Laptops {
		pdf.Cell(colWidth, 6, strconv.FormatInt(laptop.ID, 10))
		pdf.Cell(colWidth, 6, laptop.SerialNumber)
		pdf.Cell(colWidth, 6, laptop.Brand)
		pdf.Cell(colWidth, 6, laptop.Model)
		pdf.Cell(colWidth, 6, laptop.Status)
		pdf.Cell(colWidth, 6, laptop.SoftwareEngineerName)
		pdf.Ln(6)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=inventory-summary-%s.pdf", time.Now().Format("2006-01-02")))
	pdf.Output(w)
}

func (h *ReportsHandler) exportShipmentTimelineCSV(w http.ResponseWriter, data *ShipmentTimelineData) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-timeline-%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "JIRA Ticket", "Type", "Status", "Pickup Scheduled", "Picked Up", "Arrived Warehouse", "Released", "Delivered", "Total Days"})

	for _, shipment := range data.Shipments {
		formatTime := func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("2006-01-02 15:04:05")
		}
		totalDays := ""
		if shipment.TotalDurationDays != nil {
			totalDays = strconv.Itoa(*shipment.TotalDurationDays)
		}
		writer.Write([]string{
			strconv.FormatInt(shipment.ID, 10),
			shipment.JiraTicket,
			shipment.Type,
			shipment.Status,
			formatTime(shipment.PickupScheduledDate),
			formatTime(shipment.PickedUpAt),
			formatTime(shipment.ArrivedWarehouseAt),
			formatTime(shipment.ReleasedWarehouseAt),
			formatTime(shipment.DeliveredAt),
			totalDays,
		})
	}
}

func (h *ReportsHandler) exportShipmentTimelineExcel(w http.ResponseWriter, data *ShipmentTimelineData) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Shipment Timeline"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "JIRA Ticket", "Type", "Status", "Pickup Scheduled", "Picked Up", "Arrived Warehouse", "Released", "Delivered", "Total Days"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	formatTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	}

	for i, shipment := range data.Shipments {
		row := i + 2
		totalDays := ""
		if shipment.TotalDurationDays != nil {
			totalDays = strconv.Itoa(*shipment.TotalDurationDays)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), shipment.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), shipment.JiraTicket)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), shipment.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), shipment.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), formatTime(shipment.PickupScheduledDate))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), formatTime(shipment.PickedUpAt))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), formatTime(shipment.ArrivedWarehouseAt))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), formatTime(shipment.ReleasedWarehouseAt))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), formatTime(shipment.DeliveredAt))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), totalDays)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-timeline-%s.xlsx", time.Now().Format("2006-01-02")))
	f.Write(w)
}

func (h *ReportsHandler) exportShipmentTimelinePDF(w http.ResponseWriter, data *ShipmentTimelineData) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Shipment Timeline Report")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 9)
	headers := []string{"ID", "JIRA", "Status", "Pickup", "Picked", "Warehouse", "Released", "Delivered"}
	colWidth := 30.0
	for i, header := range headers {
		pdf.Cell(colWidth, 7, header)
		if i < len(headers)-1 {
			pdf.Cell(colWidth, 7, "")
		}
	}
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 7)
	formatTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02")
	}

	for _, shipment := range data.Shipments {
		pdf.Cell(colWidth, 5, strconv.FormatInt(shipment.ID, 10))
		pdf.Cell(colWidth, 5, shipment.JiraTicket)
		pdf.Cell(colWidth, 5, shipment.Status)
		pdf.Cell(colWidth, 5, formatTime(shipment.PickupScheduledDate))
		pdf.Cell(colWidth, 5, formatTime(shipment.PickedUpAt))
		pdf.Cell(colWidth, 5, formatTime(shipment.ArrivedWarehouseAt))
		pdf.Cell(colWidth, 5, formatTime(shipment.ReleasedWarehouseAt))
		pdf.Cell(colWidth, 5, formatTime(shipment.DeliveredAt))
		pdf.Ln(5)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=shipment-timeline-%s.pdf", time.Now().Format("2006-01-02")))
	pdf.Output(w)
}

