package model

import "expense.com/m/model/concern"

type CategoryNameEnum string

const (
	Bills          CategoryNameEnum = "BILLS"
	Entertainment  CategoryNameEnum = "ENTERTAINMENT"
	Education      CategoryNameEnum = "EDUCATION"
	Health         CategoryNameEnum = "HEALTH"
	Travel         CategoryNameEnum = "TRAVEL"
	Food           CategoryNameEnum = "FOOD"
	Shopping       CategoryNameEnum = "SHOPPING"
	Transportation CategoryNameEnum = "TRANSPORTATION"
	Other          CategoryNameEnum = "OTHER"
)

var CategoryNameEnumValues = []CategoryNameEnum{Bills, Entertainment, Education, Health, Travel, Food, Shopping, Transportation, Other}

type Category struct {
	concern.BaseFields

	Name              CategoryNameEnum
	TransactionTypeId string
	TransactionTypes  TransactionType `gorm:"-"`
	Icon              string
}
