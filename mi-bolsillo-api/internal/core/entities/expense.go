package entities

type Expense struct {
	ExpenseId    string  `json:"expenseId" db:"expense_id"`
	AmountPen    float64 `json:"amountPen" db:"amount_pen"`
	AmountUsd    float64 `json:"amountUsd" db:"amount_usd"`
	ExchangeRate float64 `json:"exchangeRate" db:"exchange_rate"`
	Currency     string  `json:"currency" db:"currency"`
	Description  string  `json:"description" db:"description"`
	Category     string  `json:"category" db:"category"`
	Date         string  `json:"date" db:"date"`
	BillID       string  `json:"billId" db:"bill_id"`
	UserID       string  `json:"userId" db:"user_id"`
	CreatedAt    string  `json:"createdAt" db:"created_at"`
	UpdatedAt    string  `json:"updatedAt" db:"updated_at"`
}
