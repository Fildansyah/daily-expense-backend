package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Transaction struct {
	concern.BaseFields

	UserID            string
	WalletId          string
	Wallet            Wallet `gorm:"foreignKey:WalletId;references:ID"`
	TransferWalletId  *string
	TransferWallet    *Wallet `gorm:"foreignKey:TransferWalletId;references:ID"`
	CategoryId        *string
	Category          *Category `gorm:"foreignKey:CategoryId;references:ID"`
	TransactionTypeId string
	TransactionType   TransactionType `gorm:"foreignKey:TransactionTypeId;references:ID"`
	Amount            float64
	Memo              string
	Note              string
	TransactionDate   time.Time
	UpdatedAt         time.Time
	CreatedByRuleId   *string
}
