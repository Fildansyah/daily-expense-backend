package transaction

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const dateLayout = "2006-01-02"

type SummaryRow struct {
	Type  string
	Total float64
	Count int64
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, input *entity.TransactionDTO) (*entity.TransactionDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()
	modelEntity.TransactionDate = parseDate(input.Date)

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}

	return fromModel(modelEntity), nil
}

func (r *TransactionRepository) Update(ctx context.Context, input *entity.TransactionDTO) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Updates(map[string]interface{}{
			"wallet_id":           input.WalletId,
			"transfer_wallet_id":  input.TransferWalletId,
			"category_id":         input.CategoryId,
			"transaction_type_id": input.TransactionTypeId,
			"amount":              input.Amount,
			"memo":                input.Memo,
			"note":                input.Note,
			"transaction_date":    parseDate(input.Date),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: transaksi tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Transaction{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: transaksi tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id, userID string) (*entity.TransactionDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID)
}

func (r *TransactionRepository) FindByIDForUpdate(ctx context.Context, id, userID string) (*entity.TransactionDTO, error) {
	return r.findByID(tx.DB(ctx, r.db).Clauses(clause.Locking{Strength: "UPDATE"}), id, userID)
}

func (r *TransactionRepository) FindAll(ctx context.Context, userID string, filter entity.TransactionFilter) ([]*entity.TransactionDTO, int64, error) {
	base := func() *gorm.DB {
		query := tx.DB(ctx, r.db).
			Model(&model.Transaction{}).
			Joins("JOIN transaction_types AS tt ON tt.id = transactions.transaction_type_id").
			Where("transactions.user_id = ?", userID)

		if filter.From != "" {
			query = query.Where("transactions.transaction_date >= ?", parseDate(filter.From))
		}
		if filter.To != "" {
			query = query.Where("transactions.transaction_date <= ?", parseDate(filter.To))
		}
		if filter.Type != "" {
			query = query.Where("tt.type = ?", strings.ToUpper(filter.Type))
		}
		if filter.CategoryId != "" {
			query = query.Where("transactions.category_id = ?", filter.CategoryId)
		}
		if filter.WalletId != "" {
			query = query.Where("(transactions.wallet_id = ? OR transactions.transfer_wallet_id = ?)", filter.WalletId, filter.WalletId)
		}
		if filter.Q != "" {
			like := "%" + filter.Q + "%"
			query = query.Where("(transactions.memo ILIKE ? OR transactions.note ILIKE ?)", like, like)
		}
		if filter.MinAmount != nil {
			query = query.Where("transactions.amount >= ?", *filter.MinAmount)
		}
		if filter.MaxAmount != nil {
			query = query.Where("transactions.amount <= ?", *filter.MaxAmount)
		}
		return query
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.Transaction
	err := base().
		Preload("Category").
		Preload("Wallet").
		Preload("TransferWallet").
		Preload("TransactionType").
		Order(orderClause(filter.Sort)).
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]*entity.TransactionDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}

	return items, total, nil
}

func (r *TransactionRepository) SummaryRows(ctx context.Context, userID, from, to string) ([]SummaryRow, error) {
	query := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Select("tt.type AS type, COALESCE(SUM(transactions.amount), 0) AS total, COUNT(*) AS count").
		Joins("JOIN transaction_types AS tt ON tt.id = transactions.transaction_type_id").
		Where("transactions.user_id = ?", userID)

	if from != "" {
		query = query.Where("transactions.transaction_date >= ?", parseDate(from))
	}
	if to != "" {
		query = query.Where("transactions.transaction_date <= ?", parseDate(to))
	}

	var rows []SummaryRow
	if err := query.Group("tt.type").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *TransactionRepository) FindTypeID(ctx context.Context, userID, transactionType string) (string, error) {
	var row struct {
		ID string
	}
	err := tx.DB(ctx, r.db).
		Model(&model.TransactionType{}).
		Select("id").
		Where("type = ? AND (user_id IS NULL OR user_id = ?)", transactionType, userID).
		Order("CASE WHEN user_id IS NULL THEN 1 ELSE 0 END").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return "", err
	}
	if row.ID == "" {
		return "", fmt.Errorf("%w: transaction type %s tidak tersedia", apperror.ErrNotFound, transactionType)
	}
	return row.ID, nil
}

func (r *TransactionRepository) FindCategoryType(ctx context.Context, userID, categoryID string) (string, error) {
	var row struct {
		Type string
	}
	err := tx.DB(ctx, r.db).
		Table("categories AS c").
		Select("tt.type AS type").
		Joins("JOIN transaction_types AS tt ON tt.id = c.transaction_type_id").
		Where("c.id = ? AND c.deleted_at IS NULL AND (c.user_id IS NULL OR c.user_id = ?)", categoryID, userID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return "", err
	}
	if row.Type == "" {
		return "", fmt.Errorf("%w: kategori tidak ditemukan", apperror.ErrNotFound)
	}
	return row.Type, nil
}

func (r *TransactionRepository) FindByIDs(ctx context.Context, userID string, ids []string) ([]*entity.TransactionDTO, error) {
	var rows []model.Transaction
	err := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Preload("Category").
		Preload("Wallet").
		Preload("TransferWallet").
		Preload("TransactionType").
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.TransactionDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *TransactionRepository) BulkUpdateCategory(ctx context.Context, userID string, ids []string, categoryID string) error {
	return tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Update("category_id", categoryID).Error
}

func (r *TransactionRepository) UpdateWalletID(ctx context.Context, userID, id, walletID string) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("wallet_id", walletID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: transaksi tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *TransactionRepository) findByID(db *gorm.DB, id, userID string) (*entity.TransactionDTO, error) {
	var modelEntity model.Transaction
	err := db.Preload("Category").
		Preload("Wallet").
		Preload("TransferWallet").
		Preload("TransactionType").
		Where("id = ? AND user_id = ?", id, userID).
		First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: transaksi tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func (r *TransactionRepository) FindRange(ctx context.Context, userID, from, to string) ([]*entity.TransactionDTO, error) {
	var rows []model.Transaction
	err := tx.DB(ctx, r.db).
		Preload("Category").
		Preload("Wallet").
		Preload("TransferWallet").
		Preload("TransactionType").
		Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userID, parseDate(from), parseDate(to)).
		Order("transaction_date ASC, created_at ASC").
		Limit(5000).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.TransactionDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *TransactionRepository) Recent(ctx context.Context, userID string, limit int) ([]*entity.TransactionDTO, error) {
	var rows []model.Transaction
	err := tx.DB(ctx, r.db).
		Preload("Category").
		Preload("Wallet").
		Preload("TransferWallet").
		Preload("TransactionType").
		Where("user_id = ?", userID).
		Order("transaction_date DESC, created_at DESC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.TransactionDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *TransactionRepository) SumByCategory(ctx context.Context, userID, from, to, transactionType string) ([]entity.CategorySumDTO, error) {
	query := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Select("transactions.category_id AS category_id, COALESCE(SUM(transactions.amount), 0) AS total, COUNT(*) AS count").
		Joins("JOIN transaction_types AS tt ON tt.id = transactions.transaction_type_id").
		Where("transactions.user_id = ? AND transactions.category_id IS NOT NULL AND transactions.transaction_date BETWEEN ? AND ?",
			userID, parseDate(from), parseDate(to))

	if transactionType != "" {
		query = query.Where("tt.type = ?", strings.ToUpper(transactionType))
	}

	var rows []entity.CategorySumDTO
	if err := query.Group("transactions.category_id").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *TransactionRepository) Series(ctx context.Context, userID, from, to, granularity string) ([]entity.TransactionSeriesDTO, error) {
	trunc := "month"
	switch granularity {
	case "day", "week", "month":
		trunc = granularity
	}

	query := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Select("to_char(date_trunc(?, transactions.transaction_date), 'YYYY-MM-DD') AS period, tt.type AS type, COALESCE(SUM(transactions.amount), 0) AS total", trunc).
		Joins("JOIN transaction_types AS tt ON tt.id = transactions.transaction_type_id").
		Where("transactions.user_id = ? AND transactions.transaction_date BETWEEN ? AND ?", userID, parseDate(from), parseDate(to)).
		Group("period, tt.type").
		Order("period ASC")

	var rows []entity.TransactionSeriesDTO
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func toModel(input *entity.TransactionDTO) *model.Transaction {
	return &model.Transaction{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID:            input.UserID,
		WalletId:          input.WalletId,
		TransferWalletId:  input.TransferWalletId,
		CategoryId:        input.CategoryId,
		TransactionTypeId: input.TransactionTypeId,
		Amount:            input.Amount,
		Memo:              input.Memo,
		Note:              input.Note,
	}
}

func fromModel(modelEntity *model.Transaction) *entity.TransactionDTO {
	if modelEntity == nil {
		return nil
	}

	dto := &entity.TransactionDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UpdatedAt:         modelEntity.UpdatedAt,
		UserID:            modelEntity.UserID,
		Type:              string(modelEntity.TransactionType.Type),
		TransactionTypeId: modelEntity.TransactionTypeId,
		WalletId:          modelEntity.WalletId,
		TransferWalletId:  modelEntity.TransferWalletId,
		CategoryId:        modelEntity.CategoryId,
		Amount:            modelEntity.Amount,
		Memo:              modelEntity.Memo,
		Note:              modelEntity.Note,
		Date:              modelEntity.TransactionDate.Format(dateLayout),
	}

	if modelEntity.Wallet.ID != "" {
		dto.Wallet = &entity.WalletBriefDTO{
			ID:    modelEntity.Wallet.ID,
			Name:  modelEntity.Wallet.Name,
			Type:  modelEntity.Wallet.Type,
			Color: modelEntity.Wallet.Color,
			Icon:  modelEntity.Wallet.Icon,
		}
	}
	if modelEntity.TransferWallet != nil && modelEntity.TransferWallet.ID != "" {
		dto.TransferWallet = &entity.WalletBriefDTO{
			ID:    modelEntity.TransferWallet.ID,
			Name:  modelEntity.TransferWallet.Name,
			Type:  modelEntity.TransferWallet.Type,
			Color: modelEntity.TransferWallet.Color,
			Icon:  modelEntity.TransferWallet.Icon,
		}
	}
	if modelEntity.Category != nil && modelEntity.Category.ID != "" {
		dto.Category = &entity.CategoryBriefDTO{
			ID:       modelEntity.Category.ID,
			Name:     string(modelEntity.Category.Name),
			Icon:     modelEntity.Category.Icon,
			BgColor:  modelEntity.Category.BgColor,
			ParentId: modelEntity.Category.ParentId,
		}
	}

	return dto
}

func parseDate(value string) time.Time {
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func orderClause(sort string) string {
	desc := strings.HasPrefix(sort, "-")
	key := strings.TrimPrefix(sort, "-")

	columns := map[string]string{
		"date":       "transactions.transaction_date",
		"amount":     "transactions.amount",
		"created_at": "transactions.created_at",
		"memo":       "transactions.memo",
	}

	column, ok := columns[key]
	if !ok {
		return "transactions.transaction_date DESC, transactions.created_at DESC"
	}

	direction := "ASC"
	if desc {
		direction = "DESC"
	}
	if column == "transactions.transaction_date" {
		return column + " " + direction + ", transactions.created_at DESC"
	}
	return column + " " + direction
}
