package entity

import "time"

type BudgetDTO struct {
	CommonDTOFields

	UpdatedAt   time.Time `json:"updatedAt"`
	UserID      string    `json:"-"`
	CategoryId  string    `json:"categoryId"`
	PeriodMonth string    `json:"periodMonth"`
	Amount      float64   `json:"amount"`

	CategoryName  string  `json:"categoryName"`
	CategoryIcon  string  `json:"categoryIcon"`
	CategoryColor string  `json:"categoryColor"`
	Spent         float64 `json:"spent"`
	Remaining     float64 `json:"remaining"`
	Percentage    float64 `json:"percentage"`
	Status        string  `json:"status"`
}

type BudgetPayload struct {
	CategoryId  string  `json:"categoryId"`
	PeriodMonth string  `json:"periodMonth"`
	Amount      float64 `json:"amount"`
}

type BudgetCopyPayload struct {
	FromMonth    string `json:"fromMonth"`
	ToMonth      string `json:"toMonth"`
	SkipExisting *bool  `json:"skipExisting"`
}

type BudgetUnbudgetedDTO struct {
	CategoryId   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	CategoryIcon string  `json:"categoryIcon"`
	Spent        float64 `json:"spent"`
}

type BudgetSummaryDTO struct {
	PeriodMonth string                `json:"periodMonth"`
	TotalBudget float64               `json:"totalBudget"`
	TotalSpent  float64               `json:"totalSpent"`
	Remaining   float64               `json:"remaining"`
	Percentage  float64               `json:"percentage"`
	Unbudgeted  []BudgetUnbudgetedDTO `json:"unbudgeted"`
}
