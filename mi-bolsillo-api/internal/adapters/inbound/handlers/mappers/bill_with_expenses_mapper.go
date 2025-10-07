package mappers

import (
	handlerdtos "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/dtos"
	servicedtos "github.com/KKogaa/mi-bolsillo-api/internal/core/services/dtos"
)

func ToCreateBillWithExpensesServiceDTO(handlerDTO handlerdtos.CreateBillWithExpensesRequest) servicedtos.CreateBillWithExpensesDTO {
	serviceExpenses := make([]servicedtos.CreateExpenseForBill, len(handlerDTO.Expenses))
	for i, expense := range handlerDTO.Expenses {
		serviceExpenses[i] = servicedtos.CreateExpenseForBill{
			Amount:      expense.Amount,
			Description: expense.Description,
			Category:    expense.Category,
			Date:        expense.Date,
		}
	}

	return servicedtos.CreateBillWithExpensesDTO{
		Description:  handlerDTO.Description,
		Category:     handlerDTO.Category,
		UserID:       handlerDTO.UserID,
		Date:         handlerDTO.Date,
		Currency:     handlerDTO.Currency,
		ExchangeRate: handlerDTO.ExchangeRate,
		Expenses:     serviceExpenses,
	}
}
