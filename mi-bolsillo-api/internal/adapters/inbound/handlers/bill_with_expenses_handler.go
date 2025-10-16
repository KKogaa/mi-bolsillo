package handlers

import (
	"net/http"

	handlerdtos "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/dtos"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/mappers"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/labstack/echo/v4"
)

type BillWithExpensesHandler struct {
	service *services.BillWithExpensesService
}

func NewBillWithExpensesHandler(service *services.BillWithExpensesService) *BillWithExpensesHandler {
	return &BillWithExpensesHandler{service: service}
}

// CreateBillWithExpenses godoc
// @Summary Create a new bill with expenses
// @Description Creates a new bill with associated expenses for the authenticated user
// @Tags bills
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dtos.CreateBillWithExpensesRequest true "Bill and expenses data"
// @Success 201 {object} map[string]interface{} "bill and expenses created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "User ID not found in context"
// @Failure 500 {object} map[string]string "Failed to create bill with expenses"
// @Security BearerAuth
// @Router /bills [post]
func (h *BillWithExpensesHandler) CreateBillWithExpenses(c echo.Context) error {
	var handlerDTO handlerdtos.CreateBillWithExpensesRequest

	if err := c.Bind(&handlerDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (set by Clerk auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in context",
		})
	}

	// Set the user ID from the authenticated token
	handlerDTO.UserID = userID

	// Map handler DTO to service DTO using mapper
	serviceDTO := mappers.ToCreateBillWithExpensesServiceDTO(handlerDTO)

	bill, expenses, err := h.service.CreateBillWithExpenses(serviceDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create bill with expenses",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"bill":     bill,
		"expenses": expenses,
	})
}

// ListBills godoc
// @Summary List all bills for the authenticated user
// @Description Retrieves all bills with their associated expenses for the authenticated user
// @Tags bills
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} dtos.BillWithExpensesResponse "List of bills with expenses"
// @Failure 401 {object} map[string]string "User ID not found in context"
// @Failure 500 {object} map[string]string "Failed to retrieve bills"
// @Security BearerAuth
// @Router /bills [get]
func (h *BillWithExpensesHandler) ListBills(c echo.Context) error {
	// Get user ID from context (set by Clerk auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in context",
		})
	}

	billsWithExpenses, err := h.service.ListBillsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve bills",
		})
	}

	return c.JSON(http.StatusOK, billsWithExpenses)
}

// GetBillByID godoc
// @Summary Get a specific bill by ID
// @Description Retrieves a bill with its associated expenses by ID for the authenticated user
// @Tags bills
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Bill ID"
// @Success 200 {object} map[string]interface{} "Bill with expenses"
// @Failure 400 {object} map[string]string "Bill ID is required"
// @Failure 401 {object} map[string]string "User ID not found in context"
// @Failure 403 {object} map[string]string "Access denied to this bill"
// @Failure 500 {object} map[string]string "Failed to retrieve bill"
// @Security BearerAuth
// @Router /bills/{id} [get]
func (h *BillWithExpensesHandler) GetBillByID(c echo.Context) error {
	// Get user ID from context (set by Clerk auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in context",
		})
	}

	// Get bill ID from URL parameter
	billID := c.Param("id")
	if billID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Bill ID is required",
		})
	}

	bill, expenses, err := h.service.GetBillWithExpenses(billID, userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied to this bill",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve bill",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"bill":     bill,
		"expenses": expenses,
	})
}

// DeleteBillByID godoc
// @Summary Delete a bill by ID
// @Description Deletes a bill and its associated expenses by ID for the authenticated user
// @Tags bills
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Bill ID"
// @Success 200 {object} map[string]string "Bill and associated expenses deleted successfully"
// @Failure 400 {object} map[string]string "Bill ID is required"
// @Failure 401 {object} map[string]string "User ID not found in context"
// @Failure 403 {object} map[string]string "Access denied to this bill"
// @Failure 500 {object} map[string]string "Failed to delete bill"
// @Security BearerAuth
// @Router /bills/{id} [delete]
func (h *BillWithExpensesHandler) DeleteBillByID(c echo.Context) error {
	// Get user ID from context (set by Clerk auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in context",
		})
	}

	// Get bill ID from URL parameter
	billID := c.Param("id")
	if billID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Bill ID is required",
		})
	}

	err := h.service.DeleteBillWithExpenses(billID, userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied to this bill",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete bill",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Bill and associated expenses deleted successfully",
	})
}
