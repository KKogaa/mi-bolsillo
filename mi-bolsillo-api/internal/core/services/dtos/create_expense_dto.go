package dtos

type CreateExpenseDTO struct {
	AmountPen    float64 `json:"amountPen" validate:"required"`
	AmountUsd    float64 `json:"amountUsd" validate:"required"`
	ExchangeRate float64 `json:"exchangeRate" validate:"required"`
	Currency     string  `json:"currency" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Category     string  `json:"category" validate:"required"`
	Date         string  `json:"date" validate:"required"`
}
