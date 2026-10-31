package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Bill struct {
	concern.BaseFields

	UserID     string
	LenderId   string
	Lender     *Lender `gorm:"foreignKey:LenderId;references:ID"`
	Amount     float64
	PaidAmount float64
	DueDate    time.Time
	Deskripsi  string

	WalletId   *string
	Wallet     *Wallet `gorm:"foreignKey:WalletId;references:ID"`
	CategoryId *string
	Category   *Category `gorm:"foreignKey:CategoryId;references:ID"`

	InstallmentTotal *int
	InstallmentIndex *int
	InterestRate     *float64

	RemindH3  bool
	RemindH1  bool
	RemindDue bool

	Currency  string
	UpdatedAt time.Time
}
