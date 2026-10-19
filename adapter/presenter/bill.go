package presenter

import (
	"time"

	"expense.com/m/entity"
)

type BillPresenter struct {
	ID              string    `json:"id"`
	Lender          string    `json:"lender"`
	Amount          float64   `json:"amount"`
	PaidAmount      float64   `json:"paidAmount"`
	RemainingAmount float64   `json:"remainingAmount"`
	Status          string    `json:"status"`
	DueDate         string    `json:"dueDate"`
	Deskripsi       string    `json:"deskripsi"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (v BillPresenter) FromDTO(dto *entity.BillDTO) *BillPresenter {
	if dto == nil {
		return nil
	}

	return &BillPresenter{
		ID:              dto.ID,
		Lender:          dto.Lender,
		Amount:          dto.Amount,
		PaidAmount:      dto.PaidAmount,
		RemainingAmount: dto.Amount - dto.PaidAmount,
		Status:          entity.ComputeBillStatus(dto.Amount, dto.PaidAmount),
		DueDate:         dto.DueDate,
		Deskripsi:       dto.Deskripsi,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

type BillSummaryLenderPresenter struct {
	Lender    string  `json:"lender"`
	Count     int64   `json:"count"`
	Total     float64 `json:"total"`
	Paid      float64 `json:"paid"`
	Remaining float64 `json:"remaining"`
}

type BillSummaryPresenter struct {
	Total        float64                      `json:"total"`
	Paid         float64                      `json:"paid"`
	Remaining    float64                      `json:"remaining"`
	RemainingAll float64                      `json:"remainingAll"`
	PerLender    []BillSummaryLenderPresenter `json:"perLender"`
}

func (v BillSummaryPresenter) FromDTO(dto *entity.BillSummaryDTO) *BillSummaryPresenter {
	if dto == nil {
		return nil
	}

	perLender := make([]BillSummaryLenderPresenter, len(dto.PerLender))
	for i, l := range dto.PerLender {
		perLender[i] = BillSummaryLenderPresenter{
			Lender:    l.Lender,
			Count:     l.Count,
			Total:     l.Total,
			Paid:      l.Paid,
			Remaining: l.Remaining,
		}
	}

	return &BillSummaryPresenter{
		Total:        dto.Total,
		Paid:         dto.Paid,
		Remaining:    dto.Remaining,
		RemainingAll: dto.RemainingAll,
		PerLender:    perLender,
	}
}
