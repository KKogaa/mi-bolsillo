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
