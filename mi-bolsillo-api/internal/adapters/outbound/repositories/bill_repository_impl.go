package repositories

import (
	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/jmoiron/sqlx"
)

type BillRepositoryImpl struct {
	db *sqlx.DB
}

func NewBillRepository(db *sqlx.DB) *BillRepositoryImpl {
	return &BillRepositoryImpl{db: db}
}

func (r *BillRepositoryImpl) Create(bill *entities.Bill) error {
	query := `
		INSERT INTO bills (bill_id, amount_pen, amount_usd, description, category, currency, user_id, date, created_at, updated_at)
		VALUES (:bill_id, :amount_pen, :amount_usd, :description, :category, :currency, :user_id, :date, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, bill)
	return err
}

func (r *BillRepositoryImpl) FindByID(billID string) (*entities.Bill, error) {
	var bill entities.Bill
	query := `SELECT * FROM bills WHERE bill_id = ?`
	err := r.db.Get(&bill, query, billID)
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepositoryImpl) FindByUserID(userID string) ([]*entities.Bill, error) {
	var bills []*entities.Bill
	query := `SELECT * FROM bills WHERE user_id = ? ORDER BY date DESC, created_at DESC`
	err := r.db.Select(&bills, query, userID)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (r *BillRepositoryImpl) Delete(billID string) error {
	query := `DELETE FROM bills WHERE bill_id = ?`
	_, err := r.db.Exec(query, billID)
	return err
}
