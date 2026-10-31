package presenter

import (
	"time"

	"expense.com/m/entity"
)

type CategoryTransactionTypePresenter struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type CategoryPresenter struct {
	ID                string                            `json:"id"`
	CreatedAt         time.Time                         `json:"createdAt"`
	UpdatedAt         time.Time                         `json:"updatedAt"`
	Name              string                            `json:"name"`
	TransactionTypeId string                            `json:"transactionTypeId"`
	TransactionTypes  *CategoryTransactionTypePresenter `json:"transactionTypes,omitempty"`
	Icon              string                            `json:"icon"`
	BgColor           string                            `json:"bgColor"`
	ParentId          *string                           `json:"parentId,omitempty"`
	SortOrder         int                               `json:"sortOrder"`
	IsArchived        bool                              `json:"isArchived"`
	TransactionCount  int64                             `json:"transactionCount"`
	Children          []CategoryPresenter               `json:"children"`
}

func (v CategoryPresenter) FromDTO(dto *entity.CategoryDTO) *CategoryPresenter {
	if dto == nil {
		return nil
	}

	var transactionTypes *CategoryTransactionTypePresenter
	if dto.TransactionTypeId != "" && dto.TransactionTypeName != "" {
		transactionTypes = &CategoryTransactionTypePresenter{
			ID:   dto.TransactionTypeId,
			Type: dto.TransactionTypeName,
		}
	}

	children := make([]CategoryPresenter, 0, len(dto.Children))
	for _, child := range dto.Children {
		if child == nil {
			continue
		}
		children = append(children, *new(CategoryPresenter).FromDTO(child))
	}

	return &CategoryPresenter{
		ID:                dto.ID,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
		Name:              dto.Name,
		TransactionTypeId: dto.TransactionTypeId,
		TransactionTypes:  transactionTypes,
		Icon:              dto.Icon,
		BgColor:           dto.BgColor,
		ParentId:          dto.ParentId,
		SortOrder:         dto.SortOrder,
		IsArchived:        dto.IsArchived,
		TransactionCount:  dto.TransactionCount,
		Children:          children,
	}
}

func CategoryPresenterList(dtos []*entity.CategoryDTO) []CategoryPresenter {
	out := make([]CategoryPresenter, 0, len(dtos))
	for _, item := range dtos {
		if item == nil {
			continue
		}
		out = append(out, *new(CategoryPresenter).FromDTO(item))
	}
	return out
}
