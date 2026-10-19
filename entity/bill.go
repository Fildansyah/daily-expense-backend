package entity

import (
	"time"

	"expense.com/m/model"
	"expense.com/m/model/concern"
	"gorm.io/gorm"
)

type BillDTO struct {
	CommonDTOFields

	UpdatedAt  time.Time `json:"updatedAt"`
	UserID     string    `json:"user_id"`
	Lender     string    `json:"lender"`
	Amount     float64   `json:"amount"`
	PaidAmount float64   `json:"paidAmount"`
	DueDate    string    `json:"dueDate"`
	Deskripsi  string    `json:"deskripsi"`
}

// ComputeBillStatus adalah single source of truth untuk status tagihan.
// SQL CASE pada filter status di repository HARUS sinkron dengan fungsi ini.
func ComputeBillStatus(amount, paid float64) string {
	switch {
	case paid <= 0:
		return "belum_bayar"
	case paid >= amount:
		return "lunas"
	default:
		return "sebagian"
	}
}

func (d *BillDTO) FromModel(m *model.Bill) *BillDTO {
	var deletedAt *time.Time
	if m.DeletedAt != nil && m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}

	return &BillDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			DeletedAt: deletedAt,
		},
		UpdatedAt:  m.UpdatedAt,
		UserID:     m.UserID,
		Lender:     m.Lender,
		Amount:     m.Amount,
		PaidAmount: m.PaidAmount,
		DueDate:    m.DueDate.Format("2006-01-02"),
		Deskripsi:  m.Deskripsi,
	}
}

func (d *BillDTO) ToModel() *model.Bill {
	dueDate, _ := time.Parse("2006-01-02", d.DueDate)

	m := &model.Bill{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		UserID:     d.UserID,
		Lender:     d.Lender,
		Amount:     d.Amount,
		PaidAmount: d.PaidAmount,
		DueDate:    dueDate,
		Deskripsi:  d.Deskripsi,
	}

	if d.DeletedAt != nil {
		m.DeletedAt = &gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	}

	return m
}

type BillFilter struct {
	From   string
	To     string
	Lender string
	Status string
	Q      string
}

type BillSummaryLenderDTO struct {
	Lender    string  `json:"lender"`
	Count     int64   `json:"count"`
	Total     float64 `json:"total"`
	Paid      float64 `json:"paid"`
	Remaining float64 `json:"remaining"`
}

type BillSummaryDTO struct {
	Total        float64                 `json:"total"`
	Paid         float64                 `json:"paid"`
	Remaining    float64                 `json:"remaining"`
	RemainingAll float64                 `json:"remainingAll"`
	PerLender    []*BillSummaryLenderDTO `json:"perLender"`
}
