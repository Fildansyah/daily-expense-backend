package bill

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

type LenderRow struct {
	Lender    string
	Count     int64
	Total     float64
	Paid      float64
	Remaining float64
}

type StatusRow struct {
	Status    string
	Count     int64
	Total     float64
	Paid      float64
	Remaining float64
}

type WalletRow struct {
	WalletId   string
	WalletName string
	Count      int64
	Total      float64
	Paid       float64
	Remaining  float64
}

type TotalsRow struct {
	Total     float64
	Paid      float64
	Remaining float64
}

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(ctx context.Context, input *entity.BillDTO) (*entity.BillDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()
	modelEntity.DueDate = parseDate(input.DueDate)

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return fromModel(modelEntity), nil
}

func (r *BillRepository) Update(ctx context.Context, input *entity.BillDTO) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Bill{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Updates(map[string]interface{}{
			"lender":            input.Lender,
			"amount":            input.Amount,
			"paid_amount":       input.PaidAmount,
			"due_date":          parseDate(input.DueDate),
			"deskripsi":         input.Deskripsi,
			"wallet_id":         input.WalletId,
			"category_id":       input.CategoryId,
			"installment_total": input.InstallmentTotal,
			"installment_index": input.InstallmentIndex,
			"interest_rate":     input.InterestRate,
			"remind_h3":         input.RemindH3,
			"remind_h1":         input.RemindH1,
			"remind_due":        input.RemindDue,
			"currency":          input.Currency,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: tagihan tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BillRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bill{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: tagihan tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BillRepository) FindByID(ctx context.Context, id, userID string) (*entity.BillDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID)
}

func (r *BillRepository) FindByIDForUpdate(ctx context.Context, id, userID string) (*entity.BillDTO, error) {
	return r.findByID(tx.DB(ctx, r.db).Clauses(clause.Locking{Strength: "UPDATE"}), id, userID)
}

func (r *BillRepository) FindAll(ctx context.Context, userID string, filter entity.BillFilter) ([]*entity.BillDTO, int64, error) {
	base := func() *gorm.DB {
		query := tx.DB(ctx, r.db).Model(&model.Bill{}).Where("user_id = ?", userID)

		if filter.From != "" {
			query = query.Where("due_date >= ?", parseDate(filter.From))
		}
		if filter.To != "" {
			query = query.Where("due_date <= ?", parseDate(filter.To))
		}
		if filter.Lender != "" {
			query = query.Where("lender = ?", filter.Lender)
		}
		if filter.WalletId != "" {
			query = query.Where("wallet_id = ?", filter.WalletId)
		}
		if filter.Q != "" {
			like := "%" + filter.Q + "%"
			query = query.Where("(lender ILIKE ? OR deskripsi ILIKE ?)", like, like)
		}

		switch filter.Status {
		case "lunas":
			query = query.Where("paid_amount >= amount")
		case "sebagian":
			query = query.Where("paid_amount > 0 AND paid_amount < amount")
		case "belum_bayar":
			query = query.Where("paid_amount <= 0")
		}

		return query
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.Bill
	err := base().
		Preload("Wallet").
		Preload("Category").
		Order(billOrder(filter.Sort)).
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]*entity.BillDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, total, nil
}

func (r *BillRepository) UpdatePaidAmount(ctx context.Context, id, userID string, paidAmount float64) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Bill{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("paid_amount", paidAmount)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: tagihan tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BillRepository) CreatePayment(ctx context.Context, input *entity.BillPaymentDTO) (*entity.BillPaymentDTO, error) {
	modelEntity := paymentToModel(input)
	modelEntity.ID = uuid.NewString()
	modelEntity.PaidAt = parseDate(input.PaidAt)

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return paymentFromModel(modelEntity), nil
}

func (r *BillRepository) FindPayment(ctx context.Context, id, billID, userID string) (*entity.BillPaymentDTO, error) {
	var modelEntity model.BillPayment
	err := tx.DB(ctx, r.db).
		Where("id = ? AND bill_id = ? AND user_id = ?", id, billID, userID).
		First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: pembayaran tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return paymentFromModel(&modelEntity), nil
}

func (r *BillRepository) ListPayments(ctx context.Context, billID, userID string) ([]*entity.BillPaymentDTO, error) {
	var rows []model.BillPayment
	if err := tx.DB(ctx, r.db).
		Where("bill_id = ? AND user_id = ?", billID, userID).
		Order("paid_at ASC, created_at ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]*entity.BillPaymentDTO, len(rows))
	for i := range rows {
		items[i] = paymentFromModel(&rows[i])
	}
	return items, nil
}

func (r *BillRepository) DeletePayment(ctx context.Context, id, billID, userID string) error {
	res := tx.DB(ctx, r.db).
		Where("id = ? AND bill_id = ? AND user_id = ?", id, billID, userID).
		Delete(&model.BillPayment{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: pembayaran tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BillRepository) CountPayments(ctx context.Context, billID, userID string) (int64, error) {
	var count int64
	err := tx.DB(ctx, r.db).Model(&model.BillPayment{}).
		Where("bill_id = ? AND user_id = ?", billID, userID).
		Count(&count).Error
	return count, err
}

func (r *BillRepository) DeletePaymentsByBill(ctx context.Context, billID, userID string) error {
	return tx.DB(ctx, r.db).
		Where("bill_id = ? AND user_id = ?", billID, userID).
		Delete(&model.BillPayment{}).Error
}

func (r *BillRepository) SetPaymentTransaction(ctx context.Context, paymentID, transactionID string) error {
	return tx.DB(ctx, r.db).
		Model(&model.BillPayment{}).
		Where("id = ?", paymentID).
		Update("transaction_id", transactionID).Error
}

func (r *BillRepository) Totals(ctx context.Context, userID string, filter entity.BillFilter) (TotalsRow, error) {
	var row TotalsRow
	err := periodQuery(tx.DB(ctx, r.db), userID, filter).
		Select("COALESCE(SUM(amount), 0) AS total, COALESCE(SUM(paid_amount), 0) AS paid, COALESCE(SUM(amount - paid_amount), 0) AS remaining").
		Scan(&row).Error
	return row, err
}

func (r *BillRepository) AggregateByLender(ctx context.Context, userID string, filter entity.BillFilter) ([]LenderRow, error) {
	var rows []LenderRow
	err := periodQuery(tx.DB(ctx, r.db), userID, filter).
		Select("lender, COUNT(*) AS count, COALESCE(SUM(amount), 0) AS total, COALESCE(SUM(paid_amount), 0) AS paid, COALESCE(SUM(amount - paid_amount), 0) AS remaining").
		Group("lender").
		Order("lender ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *BillRepository) AggregateByStatus(ctx context.Context, userID string, filter entity.BillFilter) ([]StatusRow, error) {
	var rows []StatusRow
	err := periodQuery(tx.DB(ctx, r.db), userID, filter).
		Select(`CASE
			WHEN paid_amount >= amount THEN 'lunas'
			WHEN paid_amount > 0 THEN 'sebagian'
			ELSE 'belum_bayar'
		END AS status, COUNT(*) AS count, COALESCE(SUM(amount), 0) AS total, COALESCE(SUM(paid_amount), 0) AS paid, COALESCE(SUM(amount - paid_amount), 0) AS remaining`).
		Group("status").
		Scan(&rows).Error
	return rows, err
}

func (r *BillRepository) AggregateByWallet(ctx context.Context, userID string, filter entity.BillFilter) ([]WalletRow, error) {
	var rows []WalletRow
	err := periodQuery(tx.DB(ctx, r.db), userID, filter).
		Select(`COALESCE(bills.wallet_id::text, '') AS wallet_id,
			COALESCE(wallets.name, '') AS wallet_name,
			COUNT(*) AS count,
			COALESCE(SUM(bills.amount), 0) AS total,
			COALESCE(SUM(bills.paid_amount), 0) AS paid,
			COALESCE(SUM(bills.amount - bills.paid_amount), 0) AS remaining`).
		Joins("LEFT JOIN wallets ON wallets.id = bills.wallet_id").
		Group("bills.wallet_id, wallets.name").
		Order("wallet_name ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *BillRepository) OutstandingAll(ctx context.Context, userID string) (float64, error) {
	var total float64
	err := tx.DB(ctx, r.db).Model(&model.Bill{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount - paid_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *BillRepository) OverdueTotal(ctx context.Context, userID string, today string) (float64, error) {
	var total float64
	err := tx.DB(ctx, r.db).Model(&model.Bill{}).
		Where("user_id = ? AND due_date < ? AND paid_amount < amount", userID, parseDate(today)).
		Select("COALESCE(SUM(amount - paid_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *BillRepository) DueTotalBetween(ctx context.Context, userID, from, to string) (float64, error) {
	var total float64
	err := tx.DB(ctx, r.db).Model(&model.Bill{}).
		Where("user_id = ? AND due_date BETWEEN ? AND ? AND paid_amount < amount", userID, parseDate(from), parseDate(to)).
		Select("COALESCE(SUM(amount - paid_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *BillRepository) findByID(db *gorm.DB, id, userID string) (*entity.BillDTO, error) {
	var modelEntity model.Bill
	err := db.Preload("Wallet").
		Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: tagihan tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func periodQuery(db *gorm.DB, userID string, filter entity.BillFilter) *gorm.DB {
	query := db.Model(&model.Bill{}).Where("user_id = ?", userID)
	if filter.From != "" {
		query = query.Where("due_date >= ?", parseDate(filter.From))
	}
	if filter.To != "" {
		query = query.Where("due_date <= ?", parseDate(filter.To))
	}
	return query
}

func billOrder(sort string) string {
	desc := strings.HasPrefix(sort, "-")
	key := strings.TrimPrefix(sort, "-")

	columns := map[string]string{
		"due_date":   "due_date",
		"amount":     "amount",
		"created_at": "created_at",
		"lender":     "lender",
	}

	column, ok := columns[key]
	if !ok {
		return "due_date ASC"
	}

	direction := "ASC"
	if desc {
		direction = "DESC"
	}
	return column + " " + direction
}

func (r *BillRepository) FindRange(ctx context.Context, userID, from, to string) ([]*entity.BillDTO, error) {
	var rows []model.Bill
	err := tx.DB(ctx, r.db).
		Preload("Wallet").
		Preload("Category").
		Where("user_id = ? AND due_date BETWEEN ? AND ?", userID, parseDate(from), parseDate(to)).
		Order("due_date ASC").
		Limit(5000).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.BillDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func toModel(input *entity.BillDTO) *model.Bill {
	return &model.Bill{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID:           input.UserID,
		Lender:           input.Lender,
		Amount:           input.Amount,
		PaidAmount:       input.PaidAmount,
		Deskripsi:        input.Deskripsi,
		WalletId:         input.WalletId,
		CategoryId:       input.CategoryId,
		InstallmentTotal: input.InstallmentTotal,
		InstallmentIndex: input.InstallmentIndex,
		InterestRate:     input.InterestRate,
		RemindH3:         input.RemindH3,
		RemindH1:         input.RemindH1,
		RemindDue:        input.RemindDue,
		Currency:         input.Currency,
	}
}

func fromModel(modelEntity *model.Bill) *entity.BillDTO {
	if modelEntity == nil {
		return nil
	}

	dto := &entity.BillDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UpdatedAt:        modelEntity.UpdatedAt,
		UserID:           modelEntity.UserID,
		Lender:           modelEntity.Lender,
		Amount:           modelEntity.Amount,
		PaidAmount:       modelEntity.PaidAmount,
		DueDate:          modelEntity.DueDate.Format(dateLayout),
		Deskripsi:        modelEntity.Deskripsi,
		WalletId:         modelEntity.WalletId,
		CategoryId:       modelEntity.CategoryId,
		InstallmentTotal: modelEntity.InstallmentTotal,
		InstallmentIndex: modelEntity.InstallmentIndex,
		InterestRate:     modelEntity.InterestRate,
		RemindH3:         modelEntity.RemindH3,
		RemindH1:         modelEntity.RemindH1,
		RemindDue:        modelEntity.RemindDue,
		Currency:         modelEntity.Currency,
	}

	if modelEntity.Wallet != nil && modelEntity.Wallet.ID != "" {
		dto.Wallet = &entity.WalletBriefDTO{
			ID:    modelEntity.Wallet.ID,
			Name:  modelEntity.Wallet.Name,
			Type:  modelEntity.Wallet.Type,
			Color: modelEntity.Wallet.Color,
			Icon:  modelEntity.Wallet.Icon,
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

func paymentToModel(input *entity.BillPaymentDTO) *model.BillPayment {
	return &model.BillPayment{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		BillId:        input.BillId,
		UserId:        input.UserID,
		Amount:        input.Amount,
		WalletId:      input.WalletId,
		Note:          input.Note,
		TransactionId: input.TransactionId,
	}
}

func paymentFromModel(modelEntity *model.BillPayment) *entity.BillPaymentDTO {
	if modelEntity == nil {
		return nil
	}

	return &entity.BillPaymentDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		BillId:        modelEntity.BillId,
		UserID:        modelEntity.UserId,
		Amount:        modelEntity.Amount,
		PaidAt:        modelEntity.PaidAt.Format(dateLayout),
		WalletId:      modelEntity.WalletId,
		Note:          modelEntity.Note,
		TransactionId: modelEntity.TransactionId,
		CreatedAt:     modelEntity.CreatedAt,
	}
}

func parseDate(value string) time.Time {
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}
