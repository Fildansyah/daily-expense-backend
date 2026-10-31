package model

import (
	"time"

	"expense.com/m/model/concern"
)

type Budget struct {
	concern.BaseFields

	UserID      string
	CategoryId  string
	PeriodMonth time.Time
	Amount      float64
	UpdatedAt   time.Time
}
