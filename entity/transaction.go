package entity

import "time"

type CategoryBriefDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Icon       string  `json:"icon"`
	BgColor    string  `json:"bgColor"`
	ParentId   *string `json:"parentId,omitempty"`
	ParentName string  `json:"parentName,omitempty"`
}

type WalletBriefDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type TransactionDTO struct {
	CommonDTOFields

	UpdatedAt         time.Time `json:"updatedAt"`
	UserID            string    `json:"-"`
	Type              string    `json:"type"`
	TransactionTypeId string    `json:"transactionTypeId"`
	WalletId          string    `json:"walletId"`
	TransferWalletId  *string   `json:"transferWalletId,omitempty"`
	CategoryId        *string   `json:"categoryId,omitempty"`
	Amount            float64   `json:"amount"`
	Memo              string    `json:"memo"`
	Note              string    `json:"note"`
	Date              string    `json:"date"`

	Category       *CategoryBriefDTO `json:"category,omitempty"`
	Wallet         *WalletBriefDTO   `json:"wallet,omitempty"`
	TransferWallet *WalletBriefDTO   `json:"transferWallet,omitempty"`
}

type TransactionFilter struct {
	From       string
	To         string
	Type       string
	CategoryId string
	WalletId   string
	Q          string
	MinAmount  *float64
	MaxAmount  *float64
	Page       int
	Limit      int
	Sort       string
}

type TransactionSummaryDTO struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
	Count   int64   `json:"count"`
}

type BulkTransactionDTO struct {
	IDs    []string `json:"ids"`
	Action string   `json:"action"`
	Value  string   `json:"value,omitempty"`
}
