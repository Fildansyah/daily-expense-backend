package model

import "expense.com/m/model/concern"

type User struct {
	concern.BaseFields
	FirstName string
	LastName  string
	Email     string
	Password  string
}
