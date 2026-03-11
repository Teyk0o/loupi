package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/middleware"
	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/services"
)

// SymptomHandler handles standalone symptom-related HTTP requests.
type SymptomHandler struct {
	symptomService *services.SymptomService
}

// NewSymptomHandler creates a new SymptomHandler instance.
func NewSymptomHandler(symptomService *services.SymptomService) *SymptomHandler {
	return &SymptomHandler{symptomService: symptomService}
}

// Create handles standalone symptom entry creation (POST /v1/symptoms).
func (h *SymptomHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	var req models.CreateSymptomEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	entry, err := h.symptomService.Create(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("create symptom entry error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to create symptom entry"})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// ListByDate handles listing symptom entries for a date (GET /v1/symptoms?date=YYYY-MM-DD).
func (h *SymptomHandler) ListByDate(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Missing required 'date' query parameter (YYYY-MM-DD)"})
		return
	}

	entries, err := h.symptomService.ListByDate(c.Request.Context(), userID, date)
	if err != nil {
		log.Printf("list symptom entries error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to list symptom entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// Update handles standalone symptom entry modification (PUT /v1/symptoms/:id).
func (h *SymptomHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid symptom entry ID"})
		return
	}

	var req models.CreateSymptomEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	entry, err := h.symptomService.Update(c.Request.Context(), userID, entryID, req)
	if err != nil {
		if errors.Is(err, services.ErrSymptomEntryNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Symptom entry not found"})
			return
		}
		log.Printf("update symptom entry error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to update symptom entry"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// Delete handles standalone symptom entry deletion (DELETE /v1/symptoms/:id).
func (h *SymptomHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid symptom entry ID"})
		return
	}

	if err := h.symptomService.Delete(c.Request.Context(), userID, entryID); err != nil {
		if errors.Is(err, services.ErrSymptomEntryNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Symptom entry not found"})
			return
		}
		log.Printf("delete symptom entry error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to delete symptom entry"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Symptom entry deleted successfully"})
}

// UpdateCheckin handles check-in modification (PUT /v1/check-ins/:id).
func (h *SymptomHandler) UpdateCheckin(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	checkinID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid check-in ID"})
		return
	}

	var req models.CreateCheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	checkin, err := h.symptomService.UpdateCheckin(c.Request.Context(), userID, checkinID, req)
	if err != nil {
		if errors.Is(err, services.ErrCheckinNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Check-in not found"})
			return
		}
		log.Printf("update checkin error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to update check-in"})
		return
	}

	c.JSON(http.StatusOK, checkin)
}

// DeleteCheckin handles check-in deletion (DELETE /v1/check-ins/:id).
func (h *SymptomHandler) DeleteCheckin(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	checkinID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid check-in ID"})
		return
	}

	if err := h.symptomService.DeleteCheckin(c.Request.Context(), userID, checkinID); err != nil {
		if errors.Is(err, services.ErrCheckinNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Check-in not found"})
			return
		}
		log.Printf("delete checkin error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to delete check-in"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Check-in deleted successfully"})
}
