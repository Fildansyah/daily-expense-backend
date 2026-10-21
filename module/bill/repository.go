package bill

import (
	"context"
	"time"

	"expense.com/m/entity"
	"expense.com/m/model"
	baserepository "expense.com/m/shared/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillRepository struct {
	*baserepository.BaseRepository[model.Bill]
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{
		BaseRepository: baserepository.NewBaseRepository[model.Bill](db),
		db:             db,
	}
}

func (r *BillRepository) CreateBill(input *entity.BillDTO) (*entity.BillDTO, error) {
	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

// findByUser mengambil bill milik user tertentu agar user tidak bisa mengakses bill user lain
func (r *BillRepository) findByUser(ctx context.Context, id, userID string) (*model.Bill, error) {
	var modelEntity model.Bill
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&modelEntity).Error
	if err != nil {
		return nil, err
	}
	return &modelEntity, nil
}

func (r *BillRepository) FindBillByID(ctx context.Context, id, userID string) (*entity.BillDTO, error) {
	modelEntity, err := r.findByUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return new(entity.BillDTO).FromModel(modelEntity), nil
}

func (r *BillRepository) UpdateBill(input *entity.BillDTO) (*entity.BillDTO, error) {
	existing, err := r.findByUser(context.Background(), input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	dueDate, _ := time.Parse("2006-01-02", input.DueDate)
	existing.Lender = input.Lender
	existing.Amount = input.Amount
	existing.PaidAmount = input.PaidAmount
	existing.DueDate = dueDate
	existing.Deskripsi = input.Deskripsi
	// UserID & CreatedAt dipertahankan; UpdatedAt di-set otomatis oleh GORM saat Save

	if err := r.Update(existing); err != nil {
		return nil, err
	}
	return input.FromModel(existing), nil
}

func (r *BillRepository) DeleteBill(id, userID string) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bill{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BillRepository) FindBills(ctx context.Context, userID string, f entity.BillFilter) ([]*entity.BillDTO, error) {
	query := r.db.WithContext(ctx).Model(&model.Bill{}).Where("user_id = ?", userID)

	if f.From != "" {
		if from, err := time.Parse("2006-01-02", f.From); err == nil {
			query = query.Where("due_date >= ?", from)
		}
	}
	if f.To != "" {
		if to, err := time.Parse("2006-01-02", f.To); err == nil {
			query = query.Where("due_date <= ?", to)
		}
	}
	if f.Lender != "" {
		query = query.Where("lender = ?", f.Lender)
	}
	if f.Status != "" {
		// HARUS sinkron dengan entity.ComputeBillStatus
		query = query.Where(`CASE
			WHEN paid_amount >= amount THEN 'lunas'
			WHEN paid_amount > 0 THEN 'sebagian'
			ELSE 'belum_bayar'
		END = ?`, f.Status)
	}
	if f.Q != "" {
		like := "%" + f.Q + "%"
		query = query.Where("lender ILIKE ? OR deskripsi ILIKE ?", like, like)
	}

	var entities []model.Bill
	if err := query.Order("due_date ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	dtos := make([]*entity.BillDTO, len(entities))
	for i, modelEntity := range entities {
		dtos[i] = new(entity.BillDTO).FromModel(&modelEntity)
	}

	return dtos, nil
}

func (r *BillRepository) GetBillSummary(ctx context.Context, userID string, f entity.BillFilter) (*entity.BillSummaryDTO, error) {
	type aggRow struct {
		Lender    string
		Count     int64
		Total     float64
		Paid      float64
		Remaining float64
	}

	// Agregasi per lender dengan filter periode
	query := r.db.WithContext(ctx).Model(&model.Bill{}).Where("user_id = ?", userID)
	if f.From != "" {
		if from, err := time.Parse("2006-01-02", f.From); err == nil {
			query = query.Where("due_date >= ?", from)
		}
	}
	if f.To != "" {
		if to, err := time.Parse("2006-01-02", f.To); err == nil {
			query = query.Where("due_date <= ?", to)
		}
	}

	var rows []aggRow
	err := query.
		Select("lender, COUNT(*) as count, " +
			"COALESCE(SUM(amount), 0) as total, " +
			"COALESCE(SUM(paid_amount), 0) as paid, " +
			"COALESCE(SUM(amount - paid_amount), 0) as remaining").
		Group("lender").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	summary := &entity.BillSummaryDTO{PerLender: make([]*entity.BillSummaryLenderDTO, len(rows))}
	for i, row := range rows {
		summary.Total += row.Total
		summary.Paid += row.Paid
		summary.Remaining += row.Remaining
		summary.PerLender[i] = &entity.BillSummaryLenderDTO{
			Lender:    row.Lender,
			Count:     row.Count,
			Total:     row.Total,
			Paid:      row.Paid,
			Remaining: row.Remaining,
		}
	}

	// remainingAll: total remaining dari SEMUA bills (tanpa filter periode)
	var all struct{ RemainingAll float64 }
	if err := r.db.WithContext(ctx).Model(&model.Bill{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount - paid_amount), 0) as remaining_all").
		Scan(&all).Error; err != nil {
		return nil, err
	}
	summary.RemainingAll = all.RemainingAll

	return summary, nil
}
