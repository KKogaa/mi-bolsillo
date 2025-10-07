package repositories

import (
	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/jmoiron/sqlx"
)

type ExpenseRepositoryImpl struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) *ExpenseRepositoryImpl {
	return &ExpenseRepositoryImpl{db: db}
}

func (r *ExpenseRepositoryImpl) Create(expense *entities.Expense) error {
	query := `
		INSERT INTO expenses (expense_id, amount_pen, amount_usd, exchange_rate, currency, description, category, date, bill_id, user_id, created_at, updated_at)
		VALUES (:expense_id, :amount_pen, :amount_usd, :exchange_rate, :currency, :description, :category, :date, :bill_id, :user_id, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, expense)
	return err
}

func (r *ExpenseRepositoryImpl) CreateBatch(expenses []*entities.Expense) error {
	query := `
		INSERT INTO expenses (expense_id, amount_pen, amount_usd, exchange_rate, currency, description, category, date, bill_id, user_id, created_at, updated_at)
		VALUES (:expense_id, :amount_pen, :amount_usd, :exchange_rate, :currency, :description, :category, :date, :bill_id, :user_id, :created_at, :updated_at)
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, expense := range expenses {
		if _, err := tx.NamedExec(query, expense); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepositoryImpl) FindByBillID(billID string) ([]*entities.Expense, error) {
	var expenses []*entities.Expense
	query := `SELECT * FROM expenses WHERE bill_id = ?`
	err := r.db.Select(&expenses, query, billID)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepositoryImpl) DeleteByBillID(billID string) error {
	query := `DELETE FROM expenses WHERE bill_id = ?`
	_, err := r.db.Exec(query, billID)
	return err
}
