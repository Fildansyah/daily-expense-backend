package wallet

import (
	"context"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/module/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository struct {
	*base_repository.BaseRepository[model.Wallet]
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		BaseRepository: base_repository.NewBaseRepository[model.Wallet](db),
		db:             db,
	}
}

func (w *WalletRepository) CreateWallet(input *entity.WalletDTO) (*entity.WalletDTO, error) {
	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := w.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (w *WalletRepository) FindWalletByUserID(ctx context.Context, userID string) ([]*entity.WalletDTO, error) {
	var wallet []model.Wallet
	err := w.db.WithContext(ctx).Model(&model.Wallet{}).
		Preload("Users").
		Where("user_id = ?", userID).Find(&wallet).Error
	if err != nil {
		return nil, err
	}

	dtos := make([]*entity.WalletDTO, len(wallet))
	for i, modelEntity := range wallet {
		dto := new(entity.WalletDTO)
		dtos[i] = dto.FromModel(&modelEntity)
	}

	return dtos, nil
}
