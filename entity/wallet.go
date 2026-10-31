package entity

import "time"

type WalletDTO struct {
	CommonDTOFields

	UserID         string    `json:"-"`
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
	StatementDay   *int      `json:"statementDay,omitempty"`
	DueDay         *int      `json:"dueDay,omitempty"`
	InterestRate   *float64  `json:"interestRate,omitempty"`
	AccountNumber  string    `json:"accountNumber,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Field turunan yang dihitung service untuk kebutuhan presentasi.
	Used        float64 `json:"used"`
	Available   float64 `json:"available"`
	Utilization float64 `json:"utilization"`
}

type WalletTypeSummaryDTO struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
}

type WalletSummaryDTO struct {
	TotalBalance float64                `json:"totalBalance"`
	CreditUsed   float64                `json:"creditUsed"`
	CreditLimit  float64                `json:"creditLimit"`
	Available    float64                `json:"available"`
	ByType       []WalletTypeSummaryDTO `json:"byType"`
}
