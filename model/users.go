package model

import "expense.com/m/model/concern"

type Users struct {
	concern.BaseFields
	FirstName string
	LastName  string
	Email     string
	Password  string
}
