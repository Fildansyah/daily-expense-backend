package wallet

import (
	"context"

	"expense.com/m/entity"
)

type WalletService struct {
	repo *WalletRepository
}

func NewWalletService(repo *WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) CreateWallet(input *entity.WalletDTO) (*entity.WalletDTO, error) {
	return s.repo.CreateWallet(input)
}

func (s *WalletService) FindWalletByUserID(ctx context.Context, userID string) ([]*entity.WalletDTO, error) {
	return s.repo.FindWalletByUserID(ctx, userID)
}
