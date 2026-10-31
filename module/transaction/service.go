package transaction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/money"
	"expense.com/m/shared/tx"
)

const (
	TypeIncome   = "INCOME"
	TypeExpense  = "EXPENSE"
	TypeTransfer = "TRANSFER"
)

// WalletGateway adalah kontrak data akun yang dibutuhkan service transaksi,
// sehingga service tidak bergantung pada implementasi repository wallet.
type WalletGateway interface {
	FindByIDs(ctx context.Context, userID string, ids []string) ([]*entity.WalletDTO, error)
	AdjustBalance(ctx context.Context, userID, walletID string, delta float64) error
}

type TransactionService struct {
	repo    *TransactionRepository
	wallets WalletGateway
	txm     *tx.Manager
}

func NewTransactionService(repo *TransactionRepository, wallets WalletGateway, txm *tx.Manager) *TransactionService {
	return &TransactionService{repo: repo, wallets: wallets, txm: txm}
}

func (s *TransactionService) Create(ctx context.Context, dto *entity.TransactionDTO) (*entity.TransactionDTO, error) {
	if err := prepareTransaction(dto); err != nil {
		return nil, err
	}

	var out *entity.TransactionDTO
	err := s.txm.Run(ctx, func(ctx context.Context) error {
		if err := s.resolveReferences(ctx, dto); err != nil {
			return err
		}

		deltas, err := s.deltasFor(ctx, dto.UserID, dto.Type, dto.WalletId, dto.TransferWalletId, dto.Amount)
		if err != nil {
			return err
		}
		if err := s.applyDeltas(ctx, dto.UserID, deltas); err != nil {
			return err
		}

		created, err := s.repo.Create(ctx, dto)
		if err != nil {
			return err
		}

		out, err = s.repo.FindByID(ctx, created.ID, dto.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *TransactionService) Update(ctx context.Context, dto *entity.TransactionDTO) (*entity.TransactionDTO, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("%w: id transaksi wajib diisi", apperror.ErrValidation)
	}
	if err := prepareTransaction(dto); err != nil {
		return nil, err
	}

	var out *entity.TransactionDTO
	err := s.txm.Run(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindByIDForUpdate(ctx, dto.ID, dto.UserID)
		if err != nil {
			return err
		}

		if err := s.resolveReferences(ctx, dto); err != nil {
			return err
		}

		oldDeltas, err := s.deltasFor(ctx, dto.UserID, existing.Type, existing.WalletId, existing.TransferWalletId, existing.Amount)
		if err != nil {
			return err
		}
		newDeltas, err := s.deltasFor(ctx, dto.UserID, dto.Type, dto.WalletId, dto.TransferWalletId, dto.Amount)
		if err != nil {
			return err
		}

		net := negateDeltas(oldDeltas)
		for walletID, delta := range newDeltas {
			net[walletID] = money.Round2(net[walletID] + delta)
		}
		if err := s.applyDeltas(ctx, dto.UserID, net); err != nil {
			return err
		}

		if err := s.repo.Update(ctx, dto); err != nil {
			return err
		}

		out, err = s.repo.FindByID(ctx, dto.ID, dto.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *TransactionService) Delete(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("%w: id transaksi wajib diisi", apperror.ErrValidation)
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindByIDForUpdate(ctx, id, userID)
		if err != nil {
			return err
		}

		deltas, err := s.deltasFor(ctx, userID, existing.Type, existing.WalletId, existing.TransferWalletId, existing.Amount)
		if err != nil {
			return err
		}
		if err := s.applyDeltas(ctx, userID, negateDeltas(deltas)); err != nil {
			return err
		}

		return s.repo.Delete(ctx, id, userID)
	})
}

func (s *TransactionService) Get(ctx context.Context, id, userID string) (*entity.TransactionDTO, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *TransactionService) List(ctx context.Context, userID string, filter *entity.TransactionFilter) ([]*entity.TransactionDTO, int64, error) {
	normalizeFilter(filter)
	return s.repo.FindAll(ctx, userID, *filter)
}

func (s *TransactionService) Summary(ctx context.Context, userID, from, to string) (*entity.TransactionSummaryDTO, error) {
	rows, err := s.repo.SummaryRows(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	summary := &entity.TransactionSummaryDTO{}
	for _, row := range rows {
		switch row.Type {
		case TypeIncome:
			summary.Income = money.Round2(row.Total)
		case TypeExpense:
			summary.Expense = money.Round2(row.Total)
		}
		summary.Count += row.Count
	}
	summary.Net = money.Round2(summary.Income - summary.Expense)

	return summary, nil
}

func (s *TransactionService) GetRange(ctx context.Context, userID, from, to string) ([]*entity.TransactionDTO, error) {
	if _, err := parseRangeDate(from); err != nil {
		return nil, err
	}
	if _, err := parseRangeDate(to); err != nil {
		return nil, err
	}
	return s.repo.FindRange(ctx, userID, from, to)
}

func (s *TransactionService) Recent(ctx context.Context, userID string, limit int) ([]*entity.TransactionDTO, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}
	return s.repo.Recent(ctx, userID, limit)
}

func (s *TransactionService) SumByCategory(ctx context.Context, userID, from, to, transactionType string) ([]entity.CategorySumDTO, error) {
	if _, err := parseRangeDate(from); err != nil {
		return nil, err
	}
	if _, err := parseRangeDate(to); err != nil {
		return nil, err
	}
	return s.repo.SumByCategory(ctx, userID, from, to, transactionType)
}

func (s *TransactionService) Series(ctx context.Context, userID, from, to, granularity string) ([]entity.TransactionSeriesDTO, error) {
	if _, err := parseRangeDate(from); err != nil {
		return nil, err
	}
	if _, err := parseRangeDate(to); err != nil {
		return nil, err
	}

	switch granularity {
	case "", "day", "week", "month":
	default:
		return nil, fmt.Errorf("%w: granularity harus day, week, atau month", apperror.ErrValidation)
	}

	return s.repo.Series(ctx, userID, from, to, granularity)
}

func parseRangeDate(value string) (time.Time, error) {
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: tanggal harus berformat yyyy-mm-dd", apperror.ErrValidation)
	}
	return parsed, nil
}

func (s *TransactionService) Bulk(ctx context.Context, userID string, bulk *entity.BulkTransactionDTO) error {
	if err := validateBulk(bulk); err != nil {
		return err
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		switch bulk.Action {
		case "delete":
			return s.bulkDelete(ctx, userID, bulk.IDs)
		case "categorize":
			return s.bulkCategorize(ctx, userID, bulk)
		case "wallet":
			return s.bulkMoveWallet(ctx, userID, bulk)
		default:
			return fmt.Errorf("%w: action tidak dikenal", apperror.ErrValidation)
		}
	})
}

func (s *TransactionService) Transfer(ctx context.Context, dto *entity.TransactionDTO) (*entity.TransactionDTO, error) {
	dto.Type = TypeTransfer
	dto.CategoryId = nil
	return s.Create(ctx, dto)
}

func (s *TransactionService) bulkDelete(ctx context.Context, userID string, ids []string) error {
	for _, id := range uniqueIDs(ids) {
		existing, err := s.repo.FindByIDForUpdate(ctx, id, userID)
		if err != nil {
			return err
		}

		deltas, err := s.deltasFor(ctx, userID, existing.Type, existing.WalletId, existing.TransferWalletId, existing.Amount)
		if err != nil {
			return err
		}
		if err := s.applyDeltas(ctx, userID, negateDeltas(deltas)); err != nil {
			return err
		}
		if err := s.repo.Delete(ctx, id, userID); err != nil {
			return err
		}
	}
	return nil
}

func (s *TransactionService) bulkCategorize(ctx context.Context, userID string, bulk *entity.BulkTransactionDTO) error {
	categoryType, err := s.repo.FindCategoryType(ctx, userID, bulk.Value)
	if err != nil {
		return err
	}

	ids := uniqueIDs(bulk.IDs)
	rows, err := s.repo.FindByIDs(ctx, userID, ids)
	if err != nil {
		return err
	}
	if len(rows) != len(ids) {
		return fmt.Errorf("%w: sebagian transaksi tidak ditemukan", apperror.ErrNotFound)
	}

	for _, row := range rows {
		if row.Type != categoryType {
			return fmt.Errorf("%w: kategori tidak sesuai tipe transaksi %s", apperror.ErrValidation, row.Type)
		}
	}

	return s.repo.BulkUpdateCategory(ctx, userID, ids, bulk.Value)
}

func (s *TransactionService) bulkMoveWallet(ctx context.Context, userID string, bulk *entity.BulkTransactionDTO) error {
	target, err := s.wallets.FindByIDs(ctx, userID, []string{bulk.Value})
	if err != nil {
		return err
	}
	if len(target) == 0 {
		return fmt.Errorf("%w: akun tujuan tidak ditemukan", apperror.ErrNotFound)
	}

	for _, id := range uniqueIDs(bulk.IDs) {
		existing, err := s.repo.FindByIDForUpdate(ctx, id, userID)
		if err != nil {
			return err
		}
		if existing.Type == TypeTransfer {
			return fmt.Errorf("%w: transaksi transfer tidak bisa dipindah akunnya", apperror.ErrValidation)
		}

		oldDeltas, err := s.deltasFor(ctx, userID, existing.Type, existing.WalletId, existing.TransferWalletId, existing.Amount)
		if err != nil {
			return err
		}
		if err := s.applyDeltas(ctx, userID, negateDeltas(oldDeltas)); err != nil {
			return err
		}

		newDeltas, err := s.deltasFor(ctx, userID, existing.Type, bulk.Value, nil, existing.Amount)
		if err != nil {
			return err
		}
		if err := s.applyDeltas(ctx, userID, newDeltas); err != nil {
			return err
		}

		if err := s.repo.UpdateWalletID(ctx, userID, id, bulk.Value); err != nil {
			return err
		}
	}

	return nil
}

func (s *TransactionService) resolveReferences(ctx context.Context, dto *entity.TransactionDTO) error {
	typeID, err := s.repo.FindTypeID(ctx, dto.UserID, dto.Type)
	if err != nil {
		return err
	}
	dto.TransactionTypeId = typeID

	if dto.Type == TypeTransfer {
		dto.CategoryId = nil
		return nil
	}

	categoryType, err := s.repo.FindCategoryType(ctx, dto.UserID, *dto.CategoryId)
	if err != nil {
		return err
	}
	if categoryType != dto.Type {
		return fmt.Errorf("%w: kategori bukan untuk tipe transaksi %s", apperror.ErrValidation, dto.Type)
	}

	return nil
}

func (s *TransactionService) deltasFor(ctx context.Context, userID, transactionType, walletID string, transferWalletID *string, amount float64) (map[string]float64, error) {
	ids := []string{walletID}
	if transferWalletID != nil {
		ids = append(ids, *transferWalletID)
	}

	wallets, err := s.wallets.FindByIDs(ctx, userID, ids)
	if err != nil {
		return nil, err
	}

	types := make(map[string]string, len(wallets))
	for _, wallet := range wallets {
		types[wallet.ID] = wallet.Type
	}
	for _, id := range uniqueIDs(ids) {
		if _, ok := types[id]; !ok {
			return nil, fmt.Errorf("%w: akun %s tidak ditemukan", apperror.ErrNotFound, id)
		}
	}

	return buildDeltas(transactionType, walletID, transferWalletID, amount, types)
}

func (s *TransactionService) applyDeltas(ctx context.Context, userID string, deltas map[string]float64) error {
	for walletID, delta := range deltas {
		if delta == 0 {
			continue
		}
		if err := s.wallets.AdjustBalance(ctx, userID, walletID, delta); err != nil {
			return err
		}
	}
	return nil
}

func buildDeltas(transactionType, walletID string, transferWalletID *string, amount float64, walletTypes map[string]string) (map[string]float64, error) {
	fromType, ok := walletTypes[walletID]
	if !ok {
		return nil, fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
	}

	deltas := map[string]float64{}

	switch transactionType {
	case TypeIncome:
		deltas[walletID] = money.SignedDelta(fromType, false, amount)
	case TypeExpense:
		deltas[walletID] = money.SignedDelta(fromType, true, amount)
	case TypeTransfer:
		if transferWalletID == nil {
			return nil, fmt.Errorf("%w: transferWalletId wajib diisi", apperror.ErrValidation)
		}
		toType, ok := walletTypes[*transferWalletID]
		if !ok {
			return nil, fmt.Errorf("%w: akun tujuan tidak ditemukan", apperror.ErrNotFound)
		}
		deltas[walletID] = money.SignedDelta(fromType, true, amount)
		deltas[*transferWalletID] = money.SignedDelta(toType, false, amount)
	default:
		return nil, fmt.Errorf("%w: tipe transaksi tidak dikenal", apperror.ErrValidation)
	}

	return deltas, nil
}

func prepareTransaction(dto *entity.TransactionDTO) error {
	dto.Type = strings.ToUpper(strings.TrimSpace(dto.Type))
	dto.Memo = strings.TrimSpace(dto.Memo)
	dto.Note = strings.TrimSpace(dto.Note)
	dto.Amount = money.Round2(dto.Amount)
	dto.Date = strings.TrimSpace(dto.Date)
	dto.TransferWalletId = normalizeOptionalID(dto.TransferWalletId)
	dto.CategoryId = normalizeOptionalID(dto.CategoryId)

	if dto.Date == "" {
		dto.Date = time.Now().UTC().Format(dateLayout)
	}

	return validateTransaction(dto)
}

func validateTransaction(dto *entity.TransactionDTO) error {
	switch dto.Type {
	case TypeIncome, TypeExpense, TypeTransfer:
	default:
		return fmt.Errorf("%w: type harus INCOME, EXPENSE, atau TRANSFER", apperror.ErrValidation)
	}

	if dto.Amount <= 0 {
		return fmt.Errorf("%w: amount harus lebih dari 0", apperror.ErrValidation)
	}
	if dto.WalletId == "" {
		return fmt.Errorf("%w: walletId wajib diisi", apperror.ErrValidation)
	}

	if dto.Type == TypeTransfer {
		if dto.TransferWalletId == nil {
			return fmt.Errorf("%w: transferWalletId wajib diisi untuk transfer", apperror.ErrValidation)
		}
		if *dto.TransferWalletId == dto.WalletId {
			return fmt.Errorf("%w: akun asal dan tujuan tidak boleh sama", apperror.ErrValidation)
		}
		dto.CategoryId = nil
	} else {
		if dto.TransferWalletId != nil {
			return fmt.Errorf("%w: transferWalletId hanya untuk tipe TRANSFER", apperror.ErrValidation)
		}
		if dto.CategoryId == nil {
			return fmt.Errorf("%w: categoryId wajib diisi", apperror.ErrValidation)
		}
	}

	parsed, err := time.Parse(dateLayout, dto.Date)
	if err != nil {
		return fmt.Errorf("%w: tanggal harus berformat yyyy-mm-dd", apperror.ErrValidation)
	}
	if parsed.After(time.Now().UTC().AddDate(1, 0, 0)) {
		return fmt.Errorf("%w: tanggal tidak boleh lebih dari 1 tahun ke depan", apperror.ErrValidation)
	}

	return nil
}

func validateBulk(bulk *entity.BulkTransactionDTO) error {
	if len(bulk.IDs) == 0 {
		return fmt.Errorf("%w: ids tidak boleh kosong", apperror.ErrValidation)
	}

	switch bulk.Action {
	case "delete":
		return nil
	case "categorize", "wallet":
		if strings.TrimSpace(bulk.Value) == "" {
			return fmt.Errorf("%w: value wajib diisi untuk aksi %s", apperror.ErrValidation, bulk.Action)
		}
		return nil
	default:
		return fmt.Errorf("%w: action harus delete, categorize, atau wallet", apperror.ErrValidation)
	}
}

func normalizeFilter(filter *entity.TransactionFilter) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	filter.Type = strings.ToUpper(strings.TrimSpace(filter.Type))
}

func normalizeOptionalID(id *string) *string {
	if id == nil {
		return nil
	}
	value := strings.TrimSpace(*id)
	if value == "" {
		return nil
	}
	return &value
}

func negateDeltas(deltas map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(deltas))
	for id, delta := range deltas {
		out[id] = -delta
	}
	return out
}

func uniqueIDs(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
