package model

import "expense.com/m/model/concern"

type Lender struct {
	concern.BaseFields

	UserID string
	Name   string
}
