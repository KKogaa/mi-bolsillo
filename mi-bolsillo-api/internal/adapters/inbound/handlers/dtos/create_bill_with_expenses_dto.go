package dtos

import "time"

// CreateBillWithExpensesRequest represents the request to create a bill with expenses
type CreateBillWithExpensesRequest struct {
	Description  string                 `json:"description" example:"Grocery shopping"`
	Category     string                 `json:"category" example:"Food"`
	Date         time.Time              `json:"date" example:"2025-10-10T10:00:00Z"`
	Currency     string                 `json:"currency" example:"USD"`
	ExchangeRate float64                `json:"exchangeRate" example:"3.75"`
	Expenses     []CreateExpenseForBill `json:"expenses"`
	// UserID is set from JWT token in the handler, not from request body
	UserID string `json:"-" swaggerignore:"true"`
	// Source is set by the handler (web or telegram), not from request body
	Source string `json:"-" swaggerignore:"true"`
}

// CreateExpenseForBill represents an expense item within a bill
type CreateExpenseForBill struct {
	Amount      float64 `json:"amount" example:"25.50"`
	Description string  `json:"description" example:"Apples"`
	Category    string  `json:"category" example:"Fruits"`
	Date        string  `json:"date" example:"2025-10-10"`
}
