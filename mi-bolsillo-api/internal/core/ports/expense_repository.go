package ports

import "github.com/KKogaa/mi-bolsillo-api/internal/core/entities"

type ExpenseRepository interface {
	Create(expense *entities.Expense) error
	CreateBatch(expenses []*entities.Expense) error
	FindByBillID(billID string) ([]*entities.Expense, error)
	DeleteByBillID(billID string) error
	UpdateUserID(oldUserID string, newUserID string) error
}
