package ports

import "github.com/KKogaa/mi-bolsillo-api/internal/core/entities"

type BillRepository interface {
	Create(bill *entities.Bill) error
	FindByID(billID string) (*entities.Bill, error)
}
