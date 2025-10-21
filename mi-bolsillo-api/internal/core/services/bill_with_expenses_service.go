package services

import (
	"errors"
	"log"
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

	// log the incoming DTO for debugging
	log.Printf("Creating bill with DTO: %+v", dto)

	// Create expense entities and calculate totals
	expenses := make([]*entities.Expense, 0, len(dto.Expenses))
	var totalAmountPen float64
	var totalAmountUsd float64

	for _, expenseDTO := range dto.Expenses {
		var amountPen, amountUsd float64

		// Convert amount based on bill's currency
		if dto.Currency == "PEN" {
			amountPen = expenseDTO.Amount
			amountUsd = expenseDTO.Amount / dto.ExchangeRate
		} else if dto.Currency == "USD" {
			amountUsd = expenseDTO.Amount
			amountPen = expenseDTO.Amount * dto.ExchangeRate
		}

		totalAmountPen += amountPen
		totalAmountUsd += amountUsd

		expense := &entities.Expense{
			ExpenseId:    uuid.New().String(),
			AmountPen:    amountPen,
			AmountUsd:    amountUsd,
			ExchangeRate: dto.ExchangeRate,
			Currency:     dto.Currency,
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

	// Create bill entity with calculated totals
	bill := &entities.Bill{
		BillId:      billID,
		AmountPen:   totalAmountPen,
		AmountUsd:   totalAmountUsd,
		Description: dto.Description,
		Category:    dto.Category,
		Currency:    dto.Currency,
		UserID:      dto.UserID,
		Date:        dto.Date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save bill
	if err := s.billRepo.Create(bill); err != nil {
		return nil, nil, err
	}

	// Save expenses in batch
	if len(expenses) > 0 {
		if err := s.expenseRepo.CreateBatch(expenses); err != nil {
			return nil, nil, err
		}
	}

	return bill, expenses, nil
}

// TODO: move this to another service that only lists the bills
func (s *BillWithExpensesService) ListBillsByUserID(userID string) ([]*dtos.BillWithExpensesResponse, error) {
	// Get all bills for the user
	bills, err := s.billRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// For each bill, get its expenses and create response
	result := make([]*dtos.BillWithExpensesResponse, 0, len(bills))
	for _, bill := range bills {
		expenses, err := s.expenseRepo.FindByBillID(bill.BillId)
		if err != nil {
			return nil, err
		}

		result = append(result, &dtos.BillWithExpensesResponse{
			BillId:      bill.BillId,
			AmountPen:   bill.AmountPen,
			AmountUsd:   bill.AmountUsd,
			Description: bill.Description,
			Category:    bill.Category,
			Currency:    bill.Currency,
			UserID:      bill.UserID,
			Date:        bill.Date,
			CreatedAt:   bill.CreatedAt,
			UpdatedAt:   bill.UpdatedAt,
			Expenses:    expenses,
		})
	}

	return result, nil
}

func (s *BillWithExpensesService) GetBillWithExpenses(billID string, userID string) (*entities.Bill, []*entities.Expense, error) {
	// Get the bill
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return nil, nil, err
	}

	// Verify the bill belongs to the user
	if bill.UserID != userID {
		return nil, nil, ErrUnauthorized
	}

	// Get expenses for the bill
	expenses, err := s.expenseRepo.FindByBillID(billID)
	if err != nil {
		return nil, nil, err
	}

	return bill, expenses, nil
}

func (s *BillWithExpensesService) DeleteBillWithExpenses(billID string, userID string) error {
	// Get the bill
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return err
	}

	// Verify the bill belongs to the user
	if bill.UserID != userID {
		return ErrUnauthorized
	}

	// Delete all expenses associated with the bill
	if err := s.expenseRepo.DeleteByBillID(billID); err != nil {
		return err
	}

	// Delete the bill
	if err := s.billRepo.Delete(billID); err != nil {
		return err
	}

	return nil
}

var ErrUnauthorized = errors.New("unauthorized access to bill")
