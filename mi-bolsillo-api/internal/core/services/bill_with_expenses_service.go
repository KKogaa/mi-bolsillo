package services

import (
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/ports"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services/dtos"
	"github.com/google/uuid"
)

type BillWithExpensesService struct {
	billRepo    ports.BillRepository
	expenseRepo ports.ExpenseRepository
}

func NewBillWithExpensesService(billRepo ports.BillRepository, expenseRepo ports.ExpenseRepository) *BillWithExpensesService {
	return &BillWithExpensesService{
		billRepo:    billRepo,
		expenseRepo: expenseRepo,
	}
}

func (s *BillWithExpensesService) CreateBillWithExpenses(dto dtos.CreateBillWithExpensesDTO) (*entities.Bill, []*entities.Expense, error) {
	now := time.Now()
	billID := uuid.New().String()

	// Create bill entity
	bill := &entities.Bill{
		BillId:      billID,
		AmountPen:   dto.AmountPen,
		AmountUsd:   dto.AmountUsd,
		Description: dto.Description,
		Category:    dto.Category,
		UserID:      dto.UserID,
		Date:        dto.Date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save bill
	if err := s.billRepo.Create(bill); err != nil {
		return nil, nil, err
	}

	// Create expense entities
	expenses := make([]*entities.Expense, 0, len(dto.Expenses))
	for _, expenseDTO := range dto.Expenses {
		expense := &entities.Expense{
			ExpenseId:    uuid.New().String(),
			AmountPen:    expenseDTO.AmountPen,
			AmountUsd:    expenseDTO.AmountUsd,
			ExchangeRate: expenseDTO.ExchangeRate,
			Currency:     expenseDTO.Currency,
			Description:  expenseDTO.Description,
			Category:     expenseDTO.Category,
			Date:         expenseDTO.Date,
			BillID:       billID,
			UserID:       dto.UserID,
			CreatedAt:    now.Format(time.RFC3339),
			UpdatedAt:    now.Format(time.RFC3339),
		}
		expenses = append(expenses, expense)
	}

	// Save expenses in batch
	if len(expenses) > 0 {
		if err := s.expenseRepo.CreateBatch(expenses); err != nil {
			return nil, nil, err
		}
	}

	return bill, expenses, nil
}
