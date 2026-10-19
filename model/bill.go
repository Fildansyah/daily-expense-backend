package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Bill struct {
	concern.BaseFields

	UserID     string
	Lender     string
	Amount     float64
	PaidAmount float64
	DueDate    time.Time
	Deskripsi  string
	UpdatedAt  time.Time
}
