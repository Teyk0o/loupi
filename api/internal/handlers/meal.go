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

// MealHandler handles meal-related HTTP requests.
type MealHandler struct {
	mealService *services.MealService
}

// NewMealHandler creates a new MealHandler instance.
func NewMealHandler(mealService *services.MealService) *MealHandler {
	return &MealHandler{mealService: mealService}
}

// Create handles meal creation (POST /v1/meals).
func (h *MealHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	var req models.CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	meal, err := h.mealService.Create(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("create meal error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to create meal"})
		return
	}

	c.JSON(http.StatusCreated, meal)
}

// GetByID handles meal retrieval (GET /v1/meals/:id).
func (h *MealHandler) GetByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	mealID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid meal ID"})
		return
	}

	meal, err := h.mealService.GetByID(c.Request.Context(), userID, mealID)
	if err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Meal not found"})
			return
		}
		log.Printf("get meal error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to get meal"})
		return
	}

	c.JSON(http.StatusOK, meal)
}

// ListByDate handles listing meals for a date (GET /v1/meals?date=YYYY-MM-DD).
func (h *MealHandler) ListByDate(c *gin.Context) {
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

	meals, err := h.mealService.ListByDate(c.Request.Context(), userID, date)
	if err != nil {
		log.Printf("list meals error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to list meals"})
		return
	}

	c.JSON(http.StatusOK, meals)
}

// Update handles meal modification (PUT /v1/meals/:id).
func (h *MealHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	mealID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid meal ID"})
		return
	}

	var req models.UpdateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	meal, err := h.mealService.Update(c.Request.Context(), userID, mealID, req)
	if err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Meal not found"})
			return
		}
		log.Printf("update meal error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to update meal"})
		return
	}

	c.JSON(http.StatusOK, meal)
}

// Delete handles meal deletion (DELETE /v1/meals/:id).
func (h *MealHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	mealID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid meal ID"})
		return
	}

	if err := h.mealService.Delete(c.Request.Context(), userID, mealID); err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Meal not found"})
			return
		}
		log.Printf("delete meal error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to delete meal"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Meal deleted successfully"})
}

// GetCheckins handles listing check-ins for a meal (GET /v1/meals/:id/check-ins).
func (h *MealHandler) GetCheckins(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	mealID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid meal ID"})
		return
	}

	checkins, err := h.mealService.GetCheckins(c.Request.Context(), userID, mealID)
	if err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Meal not found"})
			return
		}
		log.Printf("get checkins error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to get check-ins"})
		return
	}

	c.JSON(http.StatusOK, checkins)
}

// CreateCheckin handles check-in creation (POST /v1/meals/:id/check-ins).
func (h *MealHandler) CreateCheckin(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized", Message: "Not authenticated"})
		return
	}

	mealID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Invalid meal ID"})
		return
	}

	var req models.CreateCheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Invalid input data"})
		return
	}

	checkin, err := h.mealService.CreateCheckin(c.Request.Context(), userID, mealID, req)
	if err != nil {
		if errors.Is(err, services.ErrMealNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Meal not found"})
			return
		}
		log.Printf("create checkin error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Failed to create check-in"})
		return
	}

	c.JSON(http.StatusCreated, checkin)
}
