package dtos

import (
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
)

// BillWithExpensesResponse represents a bill with its associated expenses
type BillWithExpensesResponse struct {
	BillId      string              `json:"billId" example:"123e4567-e89b-12d3-a456-426614174000"`
	AmountPen   float64             `json:"amountPen" example:"95.75"`
	AmountUsd   float64             `json:"amountUsd" example:"25.50"`
	Description string              `json:"description" example:"Grocery shopping"`
	Category    string              `json:"category" example:"Food"`
	Currency    string              `json:"currency" example:"USD"`
	UserID      string              `json:"userId" example:"user_123456789"`
	Date        time.Time           `json:"date" example:"2025-10-10T10:00:00Z"`
	CreatedAt   time.Time           `json:"createdAt" example:"2025-10-10T10:00:00Z"`
	UpdatedAt   time.Time           `json:"updatedAt" example:"2025-10-10T10:00:00Z"`
	Expenses    []*entities.Expense `json:"expenses"`
}
