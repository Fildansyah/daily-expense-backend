package model

import "expense.com/m/model/concern"

type TransactionTypeEnum string

const (
	Income   TransactionTypeEnum = "INCOME"
	Expense  TransactionTypeEnum = "EXPENSE"
	Transfer TransactionTypeEnum = "TRANSFER"
)

var TransactionTypeValues = []TransactionTypeEnum{Income, Expense, Transfer}

type TransactionType struct {
	concern.BaseFields

	Type TransactionTypeEnum
}
