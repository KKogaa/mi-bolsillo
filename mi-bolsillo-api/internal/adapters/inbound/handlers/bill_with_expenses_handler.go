package handlers

import (
	"net/http"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services/dtos"
	"github.com/labstack/echo/v4"
)

type BillWithExpensesHandler struct {
	service *services.BillWithExpensesService
}

func NewBillWithExpensesHandler(service *services.BillWithExpensesService) *BillWithExpensesHandler {
	return &BillWithExpensesHandler{service: service}
}

func (h *BillWithExpensesHandler) CreateBillWithExpenses(c echo.Context) error {
	var dto dtos.CreateBillWithExpensesDTO

	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	bill, expenses, err := h.service.CreateBillWithExpenses(dto)
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
