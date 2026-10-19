package bill

import (
	"context"

	"expense.com/m/entity"
)

type BillService struct {
	repo *BillRepository
}

func NewBillService(repo *BillRepository) *BillService {
	return &BillService{repo: repo}
}

func (s *BillService) CreateBill(input *entity.BillDTO) (*entity.BillDTO, error) {
	return s.repo.CreateBill(input)
}

func (s *BillService) GetBillByID(ctx context.Context, id, userID string) (*entity.BillDTO, error) {
	return s.repo.FindBillByID(ctx, id, userID)
}

func (s *BillService) UpdateBill(input *entity.BillDTO) (*entity.BillDTO, error) {
	return s.repo.UpdateBill(input)
}

func (s *BillService) DeleteBill(id, userID string) error {
	return s.repo.DeleteBill(id, userID)
}

func (s *BillService) GetBills(ctx context.Context, userID string, filter entity.BillFilter) ([]*entity.BillDTO, error) {
	return s.repo.FindBills(ctx, userID, filter)
}

func (s *BillService) GetBillSummary(ctx context.Context, userID string, filter entity.BillFilter) (*entity.BillSummaryDTO, error) {
	return s.repo.GetBillSummary(ctx, userID, filter)
}
