package dtos

import "time"

type CreateBillWithExpensesDTO struct {
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	UserID      string                 `json:"userId"`
	Date        time.Time              `json:"date"`
	Expenses    []CreateExpenseForBill `json:"expenses"`
}

type CreateExpenseForBill struct {
	AmountPen    float64 `json:"amountPen"`
	AmountUsd    float64 `json:"amountUsd"`
	ExchangeRate float64 `json:"exchangeRate"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Date         string  `json:"date"`
}
