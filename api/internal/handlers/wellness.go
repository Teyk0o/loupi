package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/middleware"
	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/services"
)

// WellnessHandler handles wellness-related HTTP requests.
type WellnessHandler struct {
	wellnessService *services.WellnessService
}

// NewWellnessHandler creates a new WellnessHandler instance.
func NewWellnessHandler(wellnessService *services.WellnessService) *WellnessHandler {
	return &WellnessHandler{wellnessService: wellnessService}
}

// Upsert handles wellness entry creation or update (POST /v1/wellness).
func (h *WellnessHandler) Upsert(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	var req models.CreateWellnessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	entry, err := h.wellnessService.Upsert(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// GetByDate handles wellness entry retrieval (GET /v1/wellness?date=YYYY-MM-DD).
func (h *WellnessHandler) GetByDate(c *gin.Context) {
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

	entry, err := h.wellnessService.GetByDate(c.Request.Context(), userID, date)
	if err != nil {
		if errors.Is(err, services.ErrWellnessNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "No wellness entry for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to get wellness entry"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// Update handles wellness entry modification (PUT /v1/wellness/:id).
func (h *WellnessHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid wellness entry ID"})
		return
	}

	var req models.CreateWellnessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	entry, err := h.wellnessService.Update(c.Request.Context(), userID, entryID, req)
	if err != nil {
		if errors.Is(err, services.ErrWellnessNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Wellness entry not found"})
			return
		}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}
