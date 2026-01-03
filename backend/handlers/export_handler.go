package handlers

import (
	"fmt"
	"net/http"
	"time"

	"daybook-backend/services"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	service services.ExportService
}

func NewExportHandler(service services.ExportService) *ExportHandler {
	return &ExportHandler{service: service}
}

// ExportData handles data export requests
func (h *ExportHandler) ExportData(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dataType := c.Query("type") // transactions, accounts, budgets, goals, categories, assets, all
	format := c.Query("format") // csv, json
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Validate parameters
	if dataType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'type' parameter"})
		return
	}

	if format == "" {
		format = "json" // default to JSON
	}

	if format != "csv" && format != "json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Must be 'csv' or 'json'"})
		return
	}

	var data []byte
	var err error
	var filename string

	ctx := c.Request.Context()

	switch dataType {
	case "transactions":
		// Parse date range
		var startDate, endDate time.Time
		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
				return
			}
		} else {
			// Default to last year
			startDate = time.Now().AddDate(-1, 0, 0)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
				return
			}
		} else {
			endDate = time.Now()
		}

		if format == "csv" {
			data, err = h.service.ExportTransactionsCSV(ctx, userID, startDate, endDate)
			filename = fmt.Sprintf("transactions_%s_%s.csv", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		} else {
			data, err = h.service.ExportTransactionsJSON(ctx, userID, startDate, endDate)
			filename = fmt.Sprintf("transactions_%s_%s.json", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		}

	case "accounts":
		if format == "csv" {
			data, err = h.service.ExportAccountsCSV(ctx, userID)
			filename = "accounts.csv"
		} else {
			data, err = h.service.ExportAccountsJSON(ctx, userID)
			filename = "accounts.json"
		}

	case "budgets":
		if format == "csv" {
			data, err = h.service.ExportBudgetsCSV(ctx, userID)
			filename = "budgets.csv"
		} else {
			data, err = h.service.ExportBudgetsJSON(ctx, userID)
			filename = "budgets.json"
		}

	case "goals":
		if format == "csv" {
			data, err = h.service.ExportGoalsCSV(ctx, userID)
			filename = "goals.csv"
		} else {
			data, err = h.service.ExportGoalsJSON(ctx, userID)
			filename = "goals.json"
		}

	case "categories":
		if format == "csv" {
			data, err = h.service.ExportCategoriesCSV(ctx, userID)
			filename = "categories.csv"
		} else {
			data, err = h.service.ExportCategoriesJSON(ctx, userID)
			filename = "categories.json"
		}

	case "assets":
		if format == "csv" {
			data, err = h.service.ExportAssetsCSV(ctx, userID)
			filename = "assets.csv"
		} else {
			data, err = h.service.ExportAssetsJSON(ctx, userID)
			filename = "assets.json"
		}

	case "all":
		// All data export only supports JSON
		data, err = h.service.ExportAllDataJSON(ctx, userID)
		filename = fmt.Sprintf("daybook_export_%s.json", time.Now().Format("2006-01-02"))

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid type. Must be one of: transactions, accounts, budgets, goals, categories, assets, all",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate content type and headers
	var contentType string
	if format == "csv" {
		contentType = "text/csv"
	} else {
		contentType = "application/json"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Data(http.StatusOK, contentType, data)
}
