package lender

import (
	"context"

	"expense.com/m/entity"
)

type LenderService struct {
	repo *LenderRepository
}

func NewLenderService(repo *LenderRepository) *LenderService {
	return &LenderService{repo: repo}
}

func (s *LenderService) GetLenders(ctx context.Context) ([]*entity.LenderDTO, error) {
	return s.repo.FindAllLenders(ctx)
}
