package entities

import "time"

// Bill represents a bill entity
type Bill struct {
	BillId      string    `json:"billId" db:"bill_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AmountPen   float64   `json:"amountPen" db:"amount_pen" example:"95.75"`
	AmountUsd   float64   `json:"amountUsd" db:"amount_usd" example:"25.50"`
	Description string    `json:"description" db:"description" example:"Grocery shopping"`
	Category    string    `json:"category" db:"category" example:"Food"`
	Currency    string    `json:"currency" db:"currency" example:"USD"`
	UserID      string    `json:"userId" db:"user_id" example:"user_123456789"`
	Date        time.Time `json:"date" db:"date" example:"2025-10-10T10:00:00Z"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at" example:"2025-10-10T10:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at" example:"2025-10-10T10:00:00Z"`
}
