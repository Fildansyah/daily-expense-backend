package model

import "expense.com/m/model/concern"

type Wallet struct {
	concern.BaseFields
	UserID    string
	Users     Users `gorm:"foreignKey:UserID;references:ID"`
	Name      string
	Balance   float64
	Currency  string
	Type      string
	Color     string
	Icon      string
	IsExclude bool
}
