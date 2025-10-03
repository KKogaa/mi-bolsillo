package dtos

type CreateBillDTO struct {
	UserID      string             `json:"userId" validate:"required"`
	AmountPen   float64            `json:"amountPen" validate:"required"`
	AmountUsd   float64            `json:"amountUsd" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Category    string             `json:"category" validate:"required"`
	Date        string             `json:"date" validate:"required"`
	Expenses    []CreateExpenseDTO `json:"expenses,omitempty"`
}
