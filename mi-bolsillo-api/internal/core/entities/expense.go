package entities

// Expense represents an expense entity
type Expense struct {
	ExpenseId    string  `json:"expenseId" db:"expense_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	AmountPen    float64 `json:"amountPen" db:"amount_pen" example:"95.75"`
	AmountUsd    float64 `json:"amountUsd" db:"amount_usd" example:"25.50"`
	ExchangeRate float64 `json:"exchangeRate" db:"exchange_rate" example:"3.75"`
	Currency     string  `json:"currency" db:"currency" example:"USD"`
	Description  string  `json:"description" db:"description" example:"Apples"`
	Category     string  `json:"category" db:"category" example:"Fruits"`
	Date         string  `json:"date" db:"date" example:"2025-10-10"`
	BillID       string  `json:"billId" db:"bill_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID       string  `json:"userId" db:"user_id" example:"user_123456789"`
	CreatedAt    string  `json:"createdAt" db:"created_at" example:"2025-10-10T10:00:00Z"`
	UpdatedAt    string  `json:"updatedAt" db:"updated_at" example:"2025-10-10T10:00:00Z"`
}
