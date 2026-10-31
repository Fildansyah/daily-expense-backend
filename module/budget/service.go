package budget

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/money"
	"expense.com/m/shared/tx"
)

const (
	StatusSafe     = "aman"
	StatusWarning  = "waspada"
	StatusExceeded = "lebih"
)

type CategoryReader interface {
	FindVisible(ctx context.Context, userID, transactionType string, includeArchived bool) ([]*entity.CategoryDTO, error)
	GetCategoryType(ctx context.Context, id, userID string) (string, error)
}

type SpendingReader interface {
	SumByCategory(ctx context.Context, userID, from, to, transactionType string) ([]entity.CategorySumDTO, error)
}

type BudgetService struct {
	repo       *BudgetRepository
	txm        *tx.Manager
	categories CategoryReader
	spending   SpendingReader
}

func NewBudgetService(repo *BudgetRepository, txm *tx.Manager, categories CategoryReader, spending SpendingReader) *BudgetService {
	return &BudgetService{
		repo:       repo,
		txm:        txm,
		categories: categories,
		spending:   spending,
	}
}

func (s *BudgetService) GetBudgets(ctx context.Context, userID, month string) ([]*entity.BudgetDTO, error) {
	normalized, err := normalizeMonth(month)
	if err != nil {
		return nil, err
	}

	budgets, err := s.repo.ListByMonth(ctx, userID, normalized)
	if err != nil {
		return nil, err
	}
	return s.enrich(ctx, userID, normalized, budgets)
}

func (s *BudgetService) GetSummary(ctx context.Context, userID, month string) (*entity.BudgetSummaryDTO, error) {
	normalized, err := normalizeMonth(month)
	if err != nil {
		return nil, err
	}

	budgets, err := s.GetBudgets(ctx, userID, normalized)
	if err != nil {
		return nil, err
	}

	summary := &entity.BudgetSummaryDTO{PeriodMonth: normalized, Unbudgeted: []entity.BudgetUnbudgetedDTO{}}
	budgeted := make(map[string]struct{}, len(budgets))

	for _, item := range budgets {
		summary.TotalBudget = money.Round2(summary.TotalBudget + item.Amount)
		summary.TotalSpent = money.Round2(summary.TotalSpent + item.Spent)
		budgeted[item.CategoryId] = struct{}{}
	}
	summary.Remaining = money.Round2(summary.TotalBudget - summary.TotalSpent)
	summary.Percentage = percentage(summary.TotalSpent, summary.TotalBudget)

	from, to := monthBounds(normalized)
	categories, err := s.categories.FindVisible(ctx, userID, "EXPENSE", false)
	if err != nil {
		return nil, err
	}
	sums, err := s.spending.SumByCategory(ctx, userID, from, to, "EXPENSE")
	if err != nil {
		return nil, err
	}

	children := childrenIndex(categories)
	spentByCategory := sumIndex(sums)

	for _, category := range categories {
		if _, ok := budgeted[category.ID]; ok {
			continue
		}
		spent := spentIncludingChildren(category.ID, spentByCategory, children)
		if spent <= 0 {
			continue
		}
		summary.Unbudgeted = append(summary.Unbudgeted, entity.BudgetUnbudgetedDTO{
			CategoryId:   category.ID,
			CategoryName: category.Name,
			CategoryIcon: category.Icon,
			Spent:        spent,
		})
	}

	return summary, nil
}

func (s *BudgetService) CreateBudget(ctx context.Context, userID string, payload *entity.BudgetPayload) (*entity.BudgetDTO, error) {
	month, err := normalizeMonth(payload.PeriodMonth)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload.CategoryId) == "" {
		return nil, fmt.Errorf("%w: categoryId wajib diisi", apperror.ErrValidation)
	}
	amount := money.Round2(payload.Amount)
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount harus lebih dari 0", apperror.ErrValidation)
	}

	categoryType, err := s.categories.GetCategoryType(ctx, payload.CategoryId, userID)
	if err != nil {
		return nil, err
	}
	if categoryType != "EXPENSE" {
		return nil, fmt.Errorf("%w: anggaran hanya untuk kategori pengeluaran", apperror.ErrValidation)
	}

	exists, err := s.repo.Exists(ctx, userID, payload.CategoryId, month, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: anggaran kategori ini sudah ada untuk bulan tersebut", apperror.ErrConflict)
	}

	created, err := s.repo.Create(ctx, &entity.BudgetDTO{
		UserID:      userID,
		CategoryId:  payload.CategoryId,
		PeriodMonth: month,
		Amount:      amount,
	})
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrich(ctx, userID, month, []*entity.BudgetDTO{created})
	if err != nil {
		return nil, err
	}
	return enriched[0], nil
}

func (s *BudgetService) UpdateBudget(ctx context.Context, userID, id string, amount float64) (*entity.BudgetDTO, error) {
	amount = money.Round2(amount)
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount harus lebih dari 0", apperror.ErrValidation)
	}

	existing, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAmount(ctx, id, userID, amount); err != nil {
		return nil, err
	}

	existing.Amount = amount
	enriched, err := s.enrich(ctx, userID, existing.PeriodMonth, []*entity.BudgetDTO{existing})
	if err != nil {
		return nil, err
	}
	return enriched[0], nil
}

func (s *BudgetService) DeleteBudget(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *BudgetService) Copy(ctx context.Context, userID, fromMonth, toMonth string, skipExisting bool) (int, error) {
	source, err := normalizeMonth(fromMonth)
	if err != nil {
		return 0, err
	}
	target, err := normalizeMonth(toMonth)
	if err != nil {
		return 0, err
	}

	created := 0
	err = s.txm.Run(ctx, func(ctx context.Context) error {
		budgets, err := s.repo.ListByMonth(ctx, userID, source)
		if err != nil {
			return err
		}

		targetBudgets, err := s.repo.ListByMonth(ctx, userID, target)
		if err != nil {
			return err
		}
		targetByCategory := make(map[string]*entity.BudgetDTO, len(targetBudgets))
		for _, budget := range targetBudgets {
			targetByCategory[budget.CategoryId] = budget
		}

		for _, item := range budgets {
			if existing, ok := targetByCategory[item.CategoryId]; ok {
				if skipExisting {
					continue
				}
				if err := s.repo.UpdateAmount(ctx, existing.ID, userID, item.Amount); err != nil {
					return err
				}
				created++
				continue
			}

			if _, err := s.repo.Create(ctx, &entity.BudgetDTO{
				UserID:      userID,
				CategoryId:  item.CategoryId,
				PeriodMonth: target,
				Amount:      item.Amount,
			}); err != nil {
				return err
			}
			created++
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return created, nil
}

func (s *BudgetService) enrich(ctx context.Context, userID, month string, budgets []*entity.BudgetDTO) ([]*entity.BudgetDTO, error) {
	from, to := monthBounds(month)

	categories, err := s.categories.FindVisible(ctx, userID, "EXPENSE", false)
	if err != nil {
		return nil, err
	}
	sums, err := s.spending.SumByCategory(ctx, userID, from, to, "EXPENSE")
	if err != nil {
		return nil, err
	}

	byID := categoryIndex(categories)
	children := childrenIndex(categories)
	spentByCategory := sumIndex(sums)

	for _, item := range budgets {
		if category, ok := byID[item.CategoryId]; ok {
			item.CategoryName = category.Name
			item.CategoryIcon = category.Icon
			item.CategoryColor = category.BgColor
		}
		item.Spent = spentIncludingChildren(item.CategoryId, spentByCategory, children)
		item.Remaining = money.Round2(item.Amount - item.Spent)
		item.Percentage = percentage(item.Spent, item.Amount)
		item.Status = budgetStatus(item.Percentage)
	}

	return budgets, nil
}

func categoryIndex(categories []*entity.CategoryDTO) map[string]*entity.CategoryDTO {
	index := make(map[string]*entity.CategoryDTO, len(categories))
	for _, category := range categories {
		index[category.ID] = category
	}
	return index
}

func childrenIndex(categories []*entity.CategoryDTO) map[string][]string {
	children := make(map[string][]string)
	for _, category := range categories {
		if category.ParentId == nil {
			continue
		}
		children[*category.ParentId] = append(children[*category.ParentId], category.ID)
	}
	return children
}

func sumIndex(sums []entity.CategorySumDTO) map[string]float64 {
	index := make(map[string]float64, len(sums))
	for _, row := range sums {
		index[row.CategoryId] = row.Total
	}
	return index
}

func spentIncludingChildren(categoryID string, spent map[string]float64, children map[string][]string) float64 {
	total := spent[categoryID]
	for _, childID := range children[categoryID] {
		total += spent[childID]
	}
	return money.Round2(total)
}

func percentage(spent, amount float64) float64 {
	if amount <= 0 {
		return 0
	}
	return math.Round(spent/amount*1000) / 10
}

func budgetStatus(percent float64) string {
	switch {
	case percent > 100:
		return StatusExceeded
	case percent >= 80:
		return StatusWarning
	default:
		return StatusSafe
	}
}

func normalizeMonth(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC().Format(monthLayout), nil
	}

	parsed, err := time.Parse(monthLayout, value)
	if err != nil {
		return "", fmt.Errorf("%w: periodMonth harus berformat yyyy-mm", apperror.ErrValidation)
	}
	return parsed.Format(monthLayout), nil
}

func monthBounds(month string) (string, string) {
	start, _ := time.Parse(monthLayout, month)
	end := start.AddDate(0, 1, -1)
	return start.Format(dateLayout), end.Format(dateLayout)
}
