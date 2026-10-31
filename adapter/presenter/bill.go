package presenter

import (
	"time"

	"expense.com/m/entity"
)

type BillPresenter struct {
	ID               string                        `json:"id"`
	LenderId         string                        `json:"lenderId"`
	Lender           *LenderPresenter              `json:"lender,omitempty"`
	Amount           float64                       `json:"amount"`
	PaidAmount       float64                       `json:"paidAmount"`
	RemainingAmount  float64                       `json:"remainingAmount"`
	Status           string                        `json:"status"`
	IsOverdue        bool                          `json:"isOverdue"`
	DaysUntilDue     int                           `json:"daysUntilDue"`
	DueDate          string                        `json:"dueDate"`
	Deskripsi        string                        `json:"deskripsi"`
	WalletId         *string                       `json:"walletId,omitempty"`
	CategoryId       *string                       `json:"categoryId,omitempty"`
	InstallmentTotal *int                          `json:"installmentTotal,omitempty"`
	InstallmentIndex *int                          `json:"installmentIndex,omitempty"`
	InterestRate     *float64                      `json:"interestRate,omitempty"`
	RemindH3         bool                          `json:"remindH3"`
	RemindH1         bool                          `json:"remindH1"`
	RemindDue        bool                          `json:"remindDue"`
	Currency         string                        `json:"currency"`
	Wallet           *TransactionWalletPresenter   `json:"wallet,omitempty"`
	Category         *TransactionCategoryPresenter `json:"category,omitempty"`
	CreatedAt        time.Time                     `json:"createdAt"`
	UpdatedAt        time.Time                     `json:"updatedAt"`
}

func (v BillPresenter) FromDTO(dto *entity.BillDTO) *BillPresenter {
	if dto == nil {
		return nil
	}

	return &BillPresenter{
		ID:               dto.ID,
		LenderId:         dto.LenderId,
		Lender:           new(LenderPresenter).FromDTO(dto.Lender),
		Amount:           dto.Amount,
		PaidAmount:       dto.PaidAmount,
		RemainingAmount:  dto.RemainingAmount,
		Status:           dto.Status,
		IsOverdue:        dto.IsOverdue,
		DaysUntilDue:     dto.DaysUntilDue,
		DueDate:          dto.DueDate,
		Deskripsi:        dto.Deskripsi,
		WalletId:         dto.WalletId,
		CategoryId:       dto.CategoryId,
		InstallmentTotal: dto.InstallmentTotal,
		InstallmentIndex: dto.InstallmentIndex,
		InterestRate:     dto.InterestRate,
		RemindH3:         dto.RemindH3,
		RemindH1:         dto.RemindH1,
		RemindDue:        dto.RemindDue,
		Currency:         dto.Currency,
		Wallet:           walletBrief(dto.Wallet),
		Category:         categoryBrief(dto.Category),
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}

type BillPaymentPresenter struct {
	ID            string    `json:"id"`
	BillId        string    `json:"billId"`
	Amount        float64   `json:"amount"`
	PaidAt        string    `json:"paidAt"`
	WalletId      *string   `json:"walletId,omitempty"`
	Note          string    `json:"note"`
	TransactionId *string   `json:"transactionId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (v BillPaymentPresenter) FromDTO(dto *entity.BillPaymentDTO) *BillPaymentPresenter {
	if dto == nil {
		return nil
	}

	return &BillPaymentPresenter{
		ID:            dto.ID,
		BillId:        dto.BillId,
		Amount:        dto.Amount,
		PaidAt:        dto.PaidAt,
		WalletId:      dto.WalletId,
		Note:          dto.Note,
		TransactionId: dto.TransactionId,
		CreatedAt:     dto.CreatedAt,
	}
}

type BillLenderSummaryPresenter struct {
	LenderId  string           `json:"lenderId"`
	Lender    *LenderPresenter `json:"lender"`
	Count     int64            `json:"count"`
	Total     float64          `json:"total"`
	Paid      float64          `json:"paid"`
	Remaining float64          `json:"remaining"`
}

type BillStatusSummaryPresenter struct {
	Status    string  `json:"status"`
	Count     int64   `json:"count"`
	Total     float64 `json:"total"`
	Paid      float64 `json:"paid"`
	Remaining float64 `json:"remaining"`
}

type BillWalletSummaryPresenter struct {
	WalletId   string  `json:"walletId,omitempty"`
	WalletName string  `json:"walletName,omitempty"`
	Count      int64   `json:"count"`
	Total      float64 `json:"total"`
	Paid       float64 `json:"paid"`
	Remaining  float64 `json:"remaining"`
}

type BillSummaryPresenter struct {
	Total        float64                      `json:"total"`
	Paid         float64                      `json:"paid"`
	Remaining    float64                      `json:"remaining"`
	RemainingAll float64                      `json:"remainingAll"`
	OverdueTotal float64                      `json:"overdueTotal"`
	DueThisWeek  float64                      `json:"dueThisWeek"`
	PerLender    []BillLenderSummaryPresenter `json:"perLender"`
	ByStatus     []BillStatusSummaryPresenter `json:"byStatus"`
	ByWallet     []BillWalletSummaryPresenter `json:"byWallet"`
}

func (v BillSummaryPresenter) FromDTO(dto *entity.BillSummaryDTO) *BillSummaryPresenter {
	if dto == nil {
		return nil
	}

	perLender := make([]BillLenderSummaryPresenter, len(dto.PerLender))
	for i, item := range dto.PerLender {
		perLender[i] = BillLenderSummaryPresenter{
			LenderId:  item.LenderId,
			Lender:    new(LenderPresenter).FromDTO(item.Lender),
			Count:     item.Count,
			Total:     item.Total,
			Paid:      item.Paid,
			Remaining: item.Remaining,
		}
	}

	byStatus := make([]BillStatusSummaryPresenter, len(dto.ByStatus))
	for i, item := range dto.ByStatus {
		byStatus[i] = BillStatusSummaryPresenter{
			Status:    item.Status,
			Count:     item.Count,
			Total:     item.Total,
			Paid:      item.Paid,
			Remaining: item.Remaining,
		}
	}

	byWallet := make([]BillWalletSummaryPresenter, len(dto.ByWallet))
	for i, item := range dto.ByWallet {
		byWallet[i] = BillWalletSummaryPresenter{
			WalletId:   item.WalletId,
			WalletName: item.WalletName,
			Count:      item.Count,
			Total:      item.Total,
			Paid:       item.Paid,
			Remaining:  item.Remaining,
		}
	}

	return &BillSummaryPresenter{
		Total:        dto.Total,
		Paid:         dto.Paid,
		Remaining:    dto.Remaining,
		RemainingAll: dto.RemainingAll,
		OverdueTotal: dto.OverdueTotal,
		DueThisWeek:  dto.DueThisWeek,
		PerLender:    perLender,
		ByStatus:     byStatus,
		ByWallet:     byWallet,
	}
}
