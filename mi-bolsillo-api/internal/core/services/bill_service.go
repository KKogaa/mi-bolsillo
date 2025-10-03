package services

import "github.com/KKogaa/mi-bolsillo-api/internal/core/entities"

type BillService struct{}

func NewBillService() *BillService {
	return &BillService{}
}

func (s *BillService) CreateBill() (entities.Bill, error) {

	// create the list of expenses 
	
	// for each expense in the list, create the expense and associate it with the bill



	return entities.Bill{}, nil
}
