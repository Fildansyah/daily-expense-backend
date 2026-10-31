package wallet

import (
	"context"
	"errors"
	"fmt"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, input *entity.WalletDTO) (*entity.WalletDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return fromModel(modelEntity), nil
}

func (r *WalletRepository) UpdateFields(ctx context.Context, input *entity.WalletDTO) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Wallet{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Updates(map[string]interface{}{
			"name":            input.Name,
			"type":            input.Type,
			"currency":        input.Currency,
			"color":           input.Color,
			"icon":            input.Icon,
			"is_exclude":      input.IsExclude,
			"is_archived":     input.IsArchived,
			"credit_limit":    input.CreditLimit,
			"statement_day":   input.StatementDay,
			"due_day":         input.DueDay,
			"interest_rate":   input.InterestRate,
			"account_number":  input.AccountNumber,
			"opening_balance": input.OpeningBalance,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *WalletRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Wallet{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *WalletRepository) Archive(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Wallet{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_archived", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *WalletRepository) FindByID(ctx context.Context, id, userID string) (*entity.WalletDTO, error) {
	var modelEntity model.Wallet
	err := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func (r *WalletRepository) FindByIDForUpdate(ctx context.Context, id, userID string) (*entity.WalletDTO, error) {
	var modelEntity model.Wallet
	err := tx.DB(ctx, r.db).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func (r *WalletRepository) FindByUser(ctx context.Context, userID string, includeArchived bool) ([]*entity.WalletDTO, error) {
	query := tx.DB(ctx, r.db).Model(&model.Wallet{}).Where("user_id = ?", userID)
	if !includeArchived {
		query = query.Where("is_archived = FALSE")
	}

	var rows []model.Wallet
	if err := query.Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]*entity.WalletDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *WalletRepository) FindByIDs(ctx context.Context, userID string, ids []string) ([]*entity.WalletDTO, error) {
	var rows []model.Wallet
	if err := tx.DB(ctx, r.db).
		Model(&model.Wallet{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]*entity.WalletDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *WalletRepository) AdjustBalance(ctx context.Context, userID, walletID string, delta float64) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Wallet{}).
		Where("id = ? AND user_id = ?", walletID, userID).
		Update("balance", gorm.Expr("balance + ?", delta))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *WalletRepository) CountTransactions(ctx context.Context, userID, walletID string) (int64, error) {
	var count int64
	err := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("user_id = ? AND (wallet_id = ? OR transfer_wallet_id = ?)", userID, walletID, walletID).
		Count(&count).Error
	return count, err
}

func toModel(input *entity.WalletDTO) *model.Wallet {
	return &model.Wallet{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID:         input.UserID,
		Name:           input.Name,
		Balance:        input.Balance,
		OpeningBalance: input.OpeningBalance,
		Currency:       input.Currency,
		Type:           input.Type,
		Color:          input.Color,
		Icon:           input.Icon,
		IsExclude:      input.IsExclude,
		IsArchived:     input.IsArchived,
		CreditLimit:    input.CreditLimit,
		StatementDay:   input.StatementDay,
		DueDay:         input.DueDay,
		InterestRate:   input.InterestRate,
		AccountNumber:  input.AccountNumber,
	}
}

func fromModel(modelEntity *model.Wallet) *entity.WalletDTO {
	if modelEntity == nil {
		return nil
	}

	return &entity.WalletDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UserID:         modelEntity.UserID,
		Name:           modelEntity.Name,
		Balance:        modelEntity.Balance,
		OpeningBalance: modelEntity.OpeningBalance,
		Currency:       modelEntity.Currency,
		Type:           modelEntity.Type,
		Color:          modelEntity.Color,
		Icon:           modelEntity.Icon,
		IsExclude:      modelEntity.IsExclude,
		IsArchived:     modelEntity.IsArchived,
		CreditLimit:    modelEntity.CreditLimit,
		StatementDay:   modelEntity.StatementDay,
		DueDay:         modelEntity.DueDay,
		InterestRate:   modelEntity.InterestRate,
		AccountNumber:  modelEntity.AccountNumber,
		UpdatedAt:      modelEntity.UpdatedAt,
	}
}
