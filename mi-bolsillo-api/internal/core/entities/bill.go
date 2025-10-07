package entities

import "time"

type Bill struct {
	BillId      string    `json:"billId" db:"bill_id"`
	AmountPen   float64   `json:"amountPen" db:"amount_pen"`
	AmountUsd   float64   `json:"amountUsd" db:"amount_usd"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Currency    string    `json:"currency" db:"currency"`
	UserID      string    `json:"userId" db:"user_id"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
