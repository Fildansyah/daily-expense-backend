package entity

import "time"

type CategoryDTO struct {
	CommonDTOFields

	UpdatedAt         time.Time `json:"updatedAt"`
	UserID            string    `json:"-"`
	Name              string    `json:"name"`
	TransactionTypeId string    `json:"transactionTypeId"`
	BgColor           string    `json:"bgColor"`
	Icon              string    `json:"icon"`
	ParentId          *string   `json:"parentId,omitempty"`
	SortOrder         int       `json:"sortOrder"`
	IsArchived        bool      `json:"isArchived"`

	// Field turunan yang diisi service.
	TransactionTypeName string         `json:"transactionTypeName,omitempty"`
	TransactionCount    int64          `json:"transactionCount"`
	Children            []*CategoryDTO `json:"children,omitempty"`
}

type CategoryReorderItem struct {
	ID        string  `json:"id"`
	ParentId  *string `json:"parentId"`
	SortOrder int     `json:"sortOrder"`
}

type CategoryReorderRequest struct {
	Items []CategoryReorderItem `json:"items"`
}
