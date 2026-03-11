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

// CustomOptionHandler handles custom option HTTP requests.
type CustomOptionHandler struct {
	optionService *services.CustomOptionService
}

// NewCustomOptionHandler creates a new CustomOptionHandler instance.
func NewCustomOptionHandler(optionService *services.CustomOptionService) *CustomOptionHandler {
	return &CustomOptionHandler{optionService: optionService}
}

// List handles retrieval of custom options by category (GET /v1/options/:category).
func (h *CustomOptionHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	category := c.Param("category")
	options, err := h.optionService.ListByCategory(c.Request.Context(), userID, category)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid category. Must be one of: symptom_type, meal_category, sport_type"})
			return
		}
		log.Printf("list options error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to list options"})
		return
	}

	if options == nil {
		options = []models.CustomOption{}
	}

	c.JSON(http.StatusOK, options)
}

// Create handles creation of a new custom option (POST /v1/options/:category).
func (h *CustomOptionHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	category := c.Param("category")

	var req models.CreateCustomOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	option, err := h.optionService.Create(c.Request.Context(), userID, category, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid category"})
			return
		}
		if errors.Is(err, services.ErrOptionAlreadyExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{Error: "conflict", Message: "An option with this value already exists"})
			return
		}
		log.Printf("create option error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to create option"})
		return
	}

	c.JSON(http.StatusCreated, option)
}

// Update handles modification of a custom option (PUT /v1/options/:id).
func (h *CustomOptionHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	optionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid option ID"})
		return
	}

	var req models.UpdateCustomOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	option, err := h.optionService.Update(c.Request.Context(), userID, optionID, req)
	if err != nil {
		if errors.Is(err, services.ErrOptionNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Option not found"})
			return
		}
		log.Printf("update option error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to update option"})
		return
	}

	c.JSON(http.StatusOK, option)
}

// Delete handles deletion of a custom option (DELETE /v1/options/:id).
func (h *CustomOptionHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	optionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid option ID"})
		return
	}

	if err := h.optionService.Delete(c.Request.Context(), userID, optionID); err != nil {
		if errors.Is(err, services.ErrOptionNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Option not found"})
			return
		}
		log.Printf("delete option error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to delete option"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Reorder handles reordering of custom options (PUT /v1/options/:category/reorder).
func (h *CustomOptionHandler) Reorder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	category := c.Param("category")

	var req models.ReorderCustomOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	if err := h.optionService.Reorder(c.Request.Context(), userID, category, req); err != nil {
		if errors.Is(err, services.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid category"})
			return
		}
		if errors.Is(err, services.ErrOptionNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "One or more option IDs not found"})
			return
		}
		log.Printf("reorder options error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to reorder options"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Options reordered successfully"})
}
