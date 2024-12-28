package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Transaction struct {
	concern.BaseFields

	UserID            string
	Users             Users `gorm:"-"`
	WalletId          string
	Wallet            Wallet `gorm:"-"`
	CategoryId        string
	Category          Category `gorm:"-"`
	TransactionTypeId string
	TransactionType   TransactionType `gorm:"-"`
	Amount            float64
	Memo              string
	TransactionDate   time.Time
}
