package handlers

import "github.com/labstack/echo/v4"

type ExpenseHandler struct{}

func NewExpenseHandler() *ExpenseHandler {
	return &ExpenseHandler{}
}

func (h *ExpenseHandler) CreateExpense(c echo.Context) error {
	return c.JSON(201, map[string]string{
		"message": "Expense created",
	})
}
