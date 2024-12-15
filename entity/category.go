package entity

import "expense.com/m/model"

type CategoryDto struct {
	CommonDTOFields

	Name              model.CategoryNameEnum
	TransactionTypeId string
	TransactionTypes  model.TransactionTypeEnum
	Icon              string
}
