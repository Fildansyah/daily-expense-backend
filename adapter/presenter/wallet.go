package presenter

import (
	"time"

	"expense.com/m/entity"
)

type WalletPresenter struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Name           string    `json:"name"`
	Balance        float64   `json:"balance"`
	OpeningBalance float64   `json:"openingBalance"`
	Currency       string    `json:"currency"`
	Type           string    `json:"type"`
	Color          string    `json:"color"`
	Icon           string    `json:"icon"`
	IsExclude      bool      `json:"isExclude"`
	IsArchived     bool      `json:"isArchived"`
	CreditLimit    *float64  `json:"creditLimit,omitempty"`
	Used           float64   `json:"used"`
	Available      float64   `json:"available"`
	Utilization    float64   `json:"utilization"`
	StatementDay   *int      `json:"statementDay,omitempty"`
	DueDay         *int      `json:"dueDay,omitempty"`
	InterestRate   *float64  `json:"interestRate,omitempty"`
	AccountNumber  string    `json:"accountNumber,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (v WalletPresenter) FromDTO(dto *entity.WalletDTO) *WalletPresenter {
	if dto == nil {
		return nil
	}

	return &WalletPresenter{
		ID:             dto.ID,
		UserID:         dto.UserID,
		Name:           dto.Name,
		Balance:        dto.Balance,
		OpeningBalance: dto.OpeningBalance,
		Currency:       dto.Currency,
		Type:           dto.Type,
		Color:          dto.Color,
		Icon:           dto.Icon,
		IsExclude:      dto.IsExclude,
		IsArchived:     dto.IsArchived,
		CreditLimit:    dto.CreditLimit,
		Used:           dto.Used,
		Available:      dto.Available,
		Utilization:    dto.Utilization,
		StatementDay:   dto.StatementDay,
		DueDay:         dto.DueDay,
		InterestRate:   dto.InterestRate,
		AccountNumber:  dto.AccountNumber,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

type WalletTypeSummaryPresenter struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
}

type WalletSummaryPresenter struct {
	TotalBalance float64                      `json:"totalBalance"`
	CreditUsed   float64                      `json:"creditUsed"`
	CreditLimit  float64                      `json:"creditLimit"`
	Available    float64                      `json:"available"`
	ByType       []WalletTypeSummaryPresenter `json:"byType"`
}

func (v WalletSummaryPresenter) FromDTO(dto *entity.WalletSummaryDTO) *WalletSummaryPresenter {
	if dto == nil {
		return nil
	}

	byType := make([]WalletTypeSummaryPresenter, len(dto.ByType))
	for i, item := range dto.ByType {
		byType[i] = WalletTypeSummaryPresenter{Type: item.Type, Total: item.Total}
	}

	return &WalletSummaryPresenter{
		TotalBalance: dto.TotalBalance,
		CreditUsed:   dto.CreditUsed,
		CreditLimit:  dto.CreditLimit,
		Available:    dto.Available,
		ByType:       byType,
	}
}
