package dtos

import "time"

type CreateBillWithExpensesRequest struct {
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	Date         time.Time              `json:"date"`
	Currency     string                 `json:"currency"`
	ExchangeRate float64                `json:"exchangeRate"`
	Expenses     []CreateExpenseForBill `json:"expenses"`
	// UserID is set from JWT token in the handler, not from request body
	UserID string `json:"-"`
}

type CreateExpenseForBill struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}
