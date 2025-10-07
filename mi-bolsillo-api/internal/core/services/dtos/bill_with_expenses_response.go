package dtos

import (
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
)

type BillWithExpensesResponse struct {
	BillId      string              `json:"billId"`
	AmountPen   float64             `json:"amountPen"`
	AmountUsd   float64             `json:"amountUsd"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
	Currency    string              `json:"currency"`
	UserID      string              `json:"userId"`
	Date        time.Time           `json:"date"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Expenses    []*entities.Expense `json:"expenses"`
}
