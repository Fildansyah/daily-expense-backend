package model

import (
	"time"

	"expense.com/m/model/concern"
)

type BillPayment struct {
	concern.BaseFields

	BillId        string
	UserId        string
	Amount        float64
	PaidAt        time.Time
	WalletId      *string
	Note          string
	TransactionId *string
}

func (BillPayment) TableName() string {
	return "bill_payments"
}
