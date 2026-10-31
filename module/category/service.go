package category

import (
	"context"
	"fmt"
	"strings"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
)

type CategoryService struct {
	repo *CategoryRepository
	txm  *tx.Manager
}

func NewCategoryService(repo *CategoryRepository, txm *tx.Manager) *CategoryService {
	return &CategoryService{repo: repo, txm: txm}
}

func (s *CategoryService) FindVisible(ctx context.Context, userID, transactionType string, includeArchived bool) ([]*entity.CategoryDTO, error) {
	return s.repo.FindVisible(ctx, userID, transactionType, includeArchived)
}

func (s *CategoryService) GetTree(ctx context.Context, userID, transactionType string, includeArchived bool) ([]*entity.CategoryDTO, error) {
	flat, err := s.repo.FindVisible(ctx, userID, transactionType, includeArchived)
	if err != nil {
		return nil, err
	}

	counts, err := s.repo.CountByCategory(ctx, userID)
	if err != nil {
		return nil, err
	}

	typeIDs := make([]string, 0, len(flat))
	for _, category := range flat {
		typeIDs = append(typeIDs, category.TransactionTypeId)
	}
	names, err := s.repo.FindTransactionTypeNames(ctx, userID, uniqueStrings(typeIDs))
	if err != nil {
		return nil, err
	}

	for _, category := range flat {
		category.TransactionCount = counts[category.ID]
		category.TransactionTypeName = names[category.TransactionTypeId]
	}

	return buildTree(flat), nil
}

func (s *CategoryService) FindCategoryByTransactionTypeID(ctx context.Context, transactionTypeID, userID string) ([]*entity.CategoryDTO, error) {
	flat, err := s.repo.FindVisible(ctx, userID, "", false)
	if err != nil {
		return nil, err
	}

	filtered := make([]*entity.CategoryDTO, 0, len(flat))
	for _, category := range flat {
		if category.TransactionTypeId == transactionTypeID {
			filtered = append(filtered, category)
		}
	}
	return filtered, nil
}

func (s *CategoryService) GetCategoryType(ctx context.Context, id, userID string) (string, error) {
	category, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return "", err
	}
	transactionType, err := s.repo.FindTransactionType(ctx, category.TransactionTypeId, userID)
	if err != nil {
		return "", err
	}
	return transactionType.Type, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, dto *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	normalizeCategory(dto)
	if err := validateCategory(dto); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindTransactionType(ctx, dto.TransactionTypeId, dto.UserID); err != nil {
		return nil, err
	}

	if dto.ParentId != nil {
		if err := s.validateParent(ctx, dto.UserID, dto.TransactionTypeId, *dto.ParentId, ""); err != nil {
			return nil, err
		}
	}

	return s.repo.Create(ctx, dto)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, dto *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	existing, err := s.repo.FindOwned(ctx, dto.ID, dto.UserID)
	if err != nil {
		return nil, err
	}

	normalizeCategory(dto)
	if err := validateCategory(dto); err != nil {
		return nil, err
	}

	// Tipe transaksi kategori tidak boleh berubah setelah dibuat.
	dto.TransactionTypeId = existing.TransactionTypeId

	if dto.ParentId != nil {
		if err := s.validateParent(ctx, dto.UserID, existing.TransactionTypeId, *dto.ParentId, dto.ID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateFields(ctx, dto); err != nil {
		return nil, err
	}

	return s.repo.FindOwned(ctx, dto.ID, dto.UserID)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id, userID string, archive bool) error {
	if id == "" {
		return fmt.Errorf("%w: id kategori wajib diisi", apperror.ErrValidation)
	}

	if archive {
		return s.repo.Archive(ctx, id, userID)
	}

	children, err := s.repo.CountChildren(ctx, userID, id)
	if err != nil {
		return err
	}
	if children > 0 {
		return fmt.Errorf("%w: kategori masih memiliki %d sub-kategori", apperror.ErrConflict, children)
	}

	count, err := s.repo.CountTransactions(ctx, userID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%w: kategori masih dipakai oleh %d transaksi, gunakan arsip atau merge", apperror.ErrConflict, count)
	}

	return s.repo.Delete(ctx, id, userID)
}

func (s *CategoryService) Reorder(ctx context.Context, userID string, items []entity.CategoryReorderItem) error {
	if len(items) == 0 {
		return fmt.Errorf("%w: items tidak boleh kosong", apperror.ErrValidation)
	}

	for _, item := range items {
		category, err := s.repo.FindOwned(ctx, item.ID, userID)
		if err != nil {
			return err
		}
		if item.ParentId == nil {
			continue
		}
		if *item.ParentId == item.ID {
			return fmt.Errorf("%w: kategori tidak boleh menjadi induk dirinya sendiri", apperror.ErrValidation)
		}
		if err := s.validateParent(ctx, userID, category.TransactionTypeId, *item.ParentId, item.ID); err != nil {
			return err
		}
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		return s.repo.Reorder(ctx, userID, items)
	})
}

func (s *CategoryService) Merge(ctx context.Context, userID, sourceID, targetID string) error {
	if sourceID == targetID {
		return fmt.Errorf("%w: kategori sumber dan tujuan tidak boleh sama", apperror.ErrValidation)
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		source, err := s.repo.FindOwned(ctx, sourceID, userID)
		if err != nil {
			return err
		}

		target, err := s.repo.FindByID(ctx, targetID, userID)
		if err != nil {
			return err
		}

		sourceType, err := s.repo.FindTransactionType(ctx, source.TransactionTypeId, userID)
		if err != nil {
			return err
		}
		targetType, err := s.repo.FindTransactionType(ctx, target.TransactionTypeId, userID)
		if err != nil {
			return err
		}
		if sourceType.Type != targetType.Type {
			return fmt.Errorf("%w: kategori harus memiliki tipe transaksi yang sama", apperror.ErrValidation)
		}

		children, err := s.repo.CountChildren(ctx, userID, sourceID)
		if err != nil {
			return err
		}
		if children > 0 {
			return fmt.Errorf("%w: pindahkan sub-kategori terlebih dahulu", apperror.ErrConflict)
		}

		if err := s.repo.ReassignTransactions(ctx, userID, sourceID, targetID); err != nil {
			return err
		}

		return s.repo.Delete(ctx, sourceID, userID)
	})
}

func (s *CategoryService) validateParent(ctx context.Context, userID, transactionTypeID, parentID, categoryID string) error {
	if categoryID != "" && parentID == categoryID {
		return fmt.Errorf("%w: kategori tidak boleh menjadi induk dirinya sendiri", apperror.ErrValidation)
	}

	parent, err := s.repo.FindByID(ctx, parentID, userID)
	if err != nil {
		return err
	}
	if parent.ParentId != nil {
		return fmt.Errorf("%w: sub-kategori tidak boleh memiliki sub-kategori", apperror.ErrValidation)
	}

	parentType, err := s.repo.FindTransactionType(ctx, parent.TransactionTypeId, userID)
	if err != nil {
		return err
	}
	childType, err := s.repo.FindTransactionType(ctx, transactionTypeID, userID)
	if err != nil {
		return err
	}
	if parentType.Type != childType.Type {
		return fmt.Errorf("%w: tipe kategori harus sama dengan induk", apperror.ErrValidation)
	}

	return nil
}

func normalizeCategory(dto *entity.CategoryDTO) {
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.BgColor == "" {
		dto.BgColor = "#64748B"
	}
	if dto.Icon == "" {
		dto.Icon = "tag"
	}
}

func validateCategory(dto *entity.CategoryDTO) error {
	if dto.Name == "" {
		return fmt.Errorf("%w: nama kategori wajib diisi", apperror.ErrValidation)
	}
	if dto.TransactionTypeId == "" {
		return fmt.Errorf("%w: transactionTypeId wajib diisi", apperror.ErrValidation)
	}
	return nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
