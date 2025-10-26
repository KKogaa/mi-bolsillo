package dtos

import "time"

type CreateBillWithExpensesDTO struct {
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	UserID       string                 `json:"userId"`
	Source       string                 `json:"source"`
	Date         time.Time              `json:"date"`
	Currency     string                 `json:"currency"`
	ExchangeRate float64                `json:"exchangeRate"`
	Expenses     []CreateExpenseForBill `json:"expenses"`
}

type CreateExpenseForBill struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}
