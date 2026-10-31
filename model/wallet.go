package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Wallet struct {
	concern.BaseFields
	UserID         string
	Users          Users `gorm:"foreignKey:UserID;references:ID"`
	Name           string
	Balance        float64
	OpeningBalance float64
	Currency       string
	Type           string
	Color          string
	Icon           string
	IsExclude      bool
	IsArchived     bool
	CreditLimit    *float64
	StatementDay   *int
	DueDay         *int
	InterestRate   *float64
	AccountNumber  string
	UpdatedAt      time.Time
}
