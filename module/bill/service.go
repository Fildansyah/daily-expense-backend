package bill

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
	StatusUnpaid  = "belum_bayar"
	StatusPartial = "sebagian"
	StatusPaid    = "lunas"
)

type WalletReader interface {
	FindByIDs(ctx context.Context, userID string, ids []string) ([]*entity.WalletDTO, error)
}

type CategoryReader interface {
	FindByID(ctx context.Context, id, userID string) (*entity.CategoryDTO, error)
}

type ExpenseRecorder interface {
	Create(ctx context.Context, dto *entity.TransactionDTO) (*entity.TransactionDTO, error)
	Delete(ctx context.Context, id, userID string) error
}

type BillService struct {
	repo       *BillRepository
	txm        *tx.Manager
	wallets    WalletReader
	categories CategoryReader
	expenses   ExpenseRecorder
}

func NewBillService(repo *BillRepository, txm *tx.Manager, wallets WalletReader, categories CategoryReader, expenses ExpenseRecorder) *BillService {
	return &BillService{
		repo:       repo,
		txm:        txm,
		wallets:    wallets,
		categories: categories,
		expenses:   expenses,
	}
}

func (s *BillService) CreateBill(ctx context.Context, dto *entity.BillDTO) (*entity.BillDTO, error) {
	if err := prepareBill(dto); err != nil {
		return nil, err
	}
	if err := s.resolveReferences(ctx, dto); err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	enrichBill(created)
	return created, nil
}

func (s *BillService) UpdateBill(ctx context.Context, dto *entity.BillDTO) (*entity.BillDTO, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("%w: id tagihan wajib diisi", apperror.ErrValidation)
	}
	if err := prepareBill(dto); err != nil {
		return nil, err
	}

	var out *entity.BillDTO
	err := s.txm.Run(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindByIDForUpdate(ctx, dto.ID, dto.UserID)
		if err != nil {
			return err
		}

		if dto.PaidAmount < existing.PaidAmount {
			return fmt.Errorf("%w: paidAmount tidak boleh dikurangi lewat edit, hapus pembayaran terkait", apperror.ErrValidation)
		}
		if dto.PaidAmount != existing.PaidAmount {
			payments, err := s.repo.CountPayments(ctx, dto.ID, dto.UserID)
			if err != nil {
				return err
			}
			if payments > 0 {
				return fmt.Errorf("%w: perubahan pembayaran harus lewat endpoint pembayaran", apperror.ErrConflict)
			}
		}

		if err := s.resolveReferences(ctx, dto); err != nil {
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

	enrichBill(out)
	return out, nil
}

func (s *BillService) DeleteBill(ctx context.Context, id, userID string) error {
	if id == "" {
		return fmt.Errorf("%w: id tagihan wajib diisi", apperror.ErrValidation)
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		payments, err := s.repo.ListPayments(ctx, id, userID)
		if err != nil {
			return err
		}

		for _, payment := range payments {
			if payment.TransactionId != nil {
				if err := s.expenses.Delete(ctx, *payment.TransactionId, userID); err != nil {
					return err
				}
			}
		}

		if err := s.repo.DeletePaymentsByBill(ctx, id, userID); err != nil {
			return err
		}

		return s.repo.Delete(ctx, id, userID)
	})
}

func (s *BillService) GetBillByID(ctx context.Context, id, userID string) (*entity.BillDTO, error) {
	bill, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	enrichBill(bill)
	return bill, nil
}

func (s *BillService) GetBills(ctx context.Context, userID string, filter *entity.BillFilter) ([]*entity.BillDTO, int64, error) {
	if err := normalizeBillFilter(filter); err != nil {
		return nil, 0, err
	}

	items, total, err := s.repo.FindAll(ctx, userID, *filter)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range items {
		enrichBill(item)
	}
	return items, total, nil
}

func (s *BillService) GetRange(ctx context.Context, userID, from, to string) ([]*entity.BillDTO, error) {
	if _, err := parseDateStrict(from); err != nil {
		return nil, err
	}
	if _, err := parseDateStrict(to); err != nil {
		return nil, err
	}

	items, err := s.repo.FindRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		enrichBill(item)
	}
	return items, nil
}

func (s *BillService) GetBillSummary(ctx context.Context, userID, from, to string) (*entity.BillSummaryDTO, error) {
	if from != "" {
		if _, err := parseDateStrict(from); err != nil {
			return nil, err
		}
	}
	if to != "" {
		if _, err := parseDateStrict(to); err != nil {
			return nil, err
		}
	}

	filter := entity.BillFilter{From: from, To: to}

	totals, err := s.repo.Totals(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	lenders, err := s.repo.AggregateByLender(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	statuses, err := s.repo.AggregateByStatus(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	wallets, err := s.repo.AggregateByWallet(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	remainingAll, err := s.repo.OutstandingAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC()
	todayStr := today.Format(dateLayout)
	weekAhead := today.AddDate(0, 0, 7).Format(dateLayout)

	overdueTotal, err := s.repo.OverdueTotal(ctx, userID, todayStr)
	if err != nil {
		return nil, err
	}
	dueThisWeek, err := s.repo.DueTotalBetween(ctx, userID, todayStr, weekAhead)
	if err != nil {
		return nil, err
	}

	summary := &entity.BillSummaryDTO{
		Total:        money.Round2(totals.Total),
		Paid:         money.Round2(totals.Paid),
		Remaining:    money.Round2(totals.Remaining),
		RemainingAll: money.Round2(remainingAll),
		OverdueTotal: money.Round2(overdueTotal),
		DueThisWeek:  money.Round2(dueThisWeek),
		PerLender:    make([]entity.BillLenderSummaryDTO, len(lenders)),
		ByStatus:     make([]entity.BillStatusSummaryDTO, len(statuses)),
		ByWallet:     make([]entity.BillWalletSummaryDTO, len(wallets)),
	}

	for i, row := range lenders {
		summary.PerLender[i] = entity.BillLenderSummaryDTO{
			Lender:    row.Lender,
			Count:     row.Count,
			Total:     money.Round2(row.Total),
			Paid:      money.Round2(row.Paid),
			Remaining: money.Round2(row.Remaining),
		}
	}
	for i, row := range statuses {
		summary.ByStatus[i] = entity.BillStatusSummaryDTO{
			Status:    row.Status,
			Count:     row.Count,
			Total:     money.Round2(row.Total),
			Paid:      money.Round2(row.Paid),
			Remaining: money.Round2(row.Remaining),
		}
	}
	for i, row := range wallets {
		summary.ByWallet[i] = entity.BillWalletSummaryDTO{
			WalletId:   row.WalletId,
			WalletName: row.WalletName,
			Count:      row.Count,
			Total:      money.Round2(row.Total),
			Paid:       money.Round2(row.Paid),
			Remaining:  money.Round2(row.Remaining),
		}
	}

	return summary, nil
}

func (s *BillService) GetPayments(ctx context.Context, billID, userID string) ([]*entity.BillPaymentDTO, error) {
	if _, err := s.repo.FindByID(ctx, billID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListPayments(ctx, billID, userID)
}

func (s *BillService) AddPayment(ctx context.Context, billID, userID string, dto *entity.BillPaymentDTO) (*entity.BillPaymentDTO, *entity.BillDTO, *entity.TransactionDTO, error) {
	if billID == "" {
		return nil, nil, nil, fmt.Errorf("%w: id tagihan wajib diisi", apperror.ErrValidation)
	}
	if err := preparePayment(dto); err != nil {
		return nil, nil, nil, err
	}

	var payment *entity.BillPaymentDTO
	var updatedBill *entity.BillDTO
	var createdTransaction *entity.TransactionDTO

	err := s.txm.Run(ctx, func(ctx context.Context) error {
		bill, err := s.repo.FindByIDForUpdate(ctx, billID, userID)
		if err != nil {
			return err
		}

		if err := validatePaymentAmount(dto.Amount, bill.PaidAmount, bill.Amount); err != nil {
			return err
		}

		dto.BillId = billID
		dto.UserID = userID

		payment, err = s.repo.CreatePayment(ctx, dto)
		if err != nil {
			return err
		}

		paidAmount := money.Round2(bill.PaidAmount + dto.Amount)
		if err := s.repo.UpdatePaidAmount(ctx, billID, userID, paidAmount); err != nil {
			return err
		}

		if dto.CreateTransaction {
			if dto.WalletId == nil {
				return fmt.Errorf("%w: walletId wajib diisi untuk membuat transaksi", apperror.ErrValidation)
			}

			categoryID := dto.CategoryId
			if categoryID == nil {
				categoryID = bill.CategoryId
			}
			if categoryID == nil {
				return fmt.Errorf("%w: categoryId wajib diisi untuk membuat transaksi", apperror.ErrValidation)
			}

			createdTransaction, err = s.expenses.Create(ctx, &entity.TransactionDTO{
				UserID:     userID,
				Type:       "EXPENSE",
				WalletId:   *dto.WalletId,
				CategoryId: categoryID,
				Amount:     dto.Amount,
				Memo:       "Bayar " + bill.Lender,
				Date:       dto.PaidAt,
			})
			if err != nil {
				return err
			}

			if err := s.repo.SetPaymentTransaction(ctx, payment.ID, createdTransaction.ID); err != nil {
				return err
			}
			payment.TransactionId = &createdTransaction.ID
		}

		updatedBill, err = s.repo.FindByID(ctx, billID, userID)
		if err != nil {
			return err
		}
		enrichBill(updatedBill)

		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return payment, updatedBill, createdTransaction, nil
}

func (s *BillService) DeletePayment(ctx context.Context, billID, paymentID, userID string) error {
	if billID == "" || paymentID == "" {
		return fmt.Errorf("%w: id tagihan dan pembayaran wajib diisi", apperror.ErrValidation)
	}

	return s.txm.Run(ctx, func(ctx context.Context) error {
		bill, err := s.repo.FindByIDForUpdate(ctx, billID, userID)
		if err != nil {
			return err
		}

		payment, err := s.repo.FindPayment(ctx, paymentID, billID, userID)
		if err != nil {
			return err
		}

		if err := s.repo.DeletePayment(ctx, paymentID, billID, userID); err != nil {
			return err
		}

		paidAmount := money.Round2(bill.PaidAmount - payment.Amount)
		if paidAmount < 0 {
			paidAmount = 0
		}
		if err := s.repo.UpdatePaidAmount(ctx, billID, userID, paidAmount); err != nil {
			return err
		}

		if payment.TransactionId != nil {
			if err := s.expenses.Delete(ctx, *payment.TransactionId, userID); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *BillService) resolveReferences(ctx context.Context, dto *entity.BillDTO) error {
	if dto.WalletId != nil {
		wallets, err := s.wallets.FindByIDs(ctx, dto.UserID, []string{*dto.WalletId})
		if err != nil {
			return err
		}
		if len(wallets) == 0 {
			return fmt.Errorf("%w: akun tidak ditemukan", apperror.ErrNotFound)
		}
	}

	if dto.CategoryId != nil {
		if _, err := s.categories.FindByID(ctx, *dto.CategoryId, dto.UserID); err != nil {
			return err
		}
	}

	return nil
}

func prepareBill(dto *entity.BillDTO) error {
	dto.Lender = strings.TrimSpace(dto.Lender)
	dto.Deskripsi = strings.TrimSpace(dto.Deskripsi)
	dto.Currency = strings.ToUpper(strings.TrimSpace(dto.Currency))
	if dto.Currency == "" {
		dto.Currency = "IDR"
	}
	dto.WalletId = normalizeOptionalID(dto.WalletId)
	dto.CategoryId = normalizeOptionalID(dto.CategoryId)
	dto.Amount = money.Round2(dto.Amount)
	dto.PaidAmount = money.Round2(dto.PaidAmount)

	if dto.Lender == "" {
		return fmt.Errorf("%w: lender wajib diisi", apperror.ErrValidation)
	}
	if dto.Amount <= 0 {
		return fmt.Errorf("%w: amount harus lebih dari 0", apperror.ErrValidation)
	}
	if dto.PaidAmount < 0 {
		return fmt.Errorf("%w: paidAmount tidak boleh negatif", apperror.ErrValidation)
	}
	if dto.PaidAmount > dto.Amount {
		return fmt.Errorf("%w: paidAmount tidak boleh melebihi amount", apperror.ErrValidation)
	}
	if _, err := parseDateStrict(dto.DueDate); err != nil {
		return err
	}
	if dto.InstallmentTotal != nil && *dto.InstallmentTotal <= 0 {
		return fmt.Errorf("%w: installmentTotal harus lebih dari 0", apperror.ErrValidation)
	}
	if dto.InstallmentIndex != nil {
		if *dto.InstallmentIndex <= 0 {
			return fmt.Errorf("%w: installmentIndex harus lebih dari 0", apperror.ErrValidation)
		}
		if dto.InstallmentTotal != nil && *dto.InstallmentIndex > *dto.InstallmentTotal {
			return fmt.Errorf("%w: installmentIndex tidak boleh melebihi installmentTotal", apperror.ErrValidation)
		}
	}
	if dto.InterestRate != nil && *dto.InterestRate < 0 {
		return fmt.Errorf("%w: interestRate tidak boleh negatif", apperror.ErrValidation)
	}

	return nil
}

func preparePayment(dto *entity.BillPaymentDTO) error {
	dto.Note = strings.TrimSpace(dto.Note)
	dto.WalletId = normalizeOptionalID(dto.WalletId)
	dto.CategoryId = normalizeOptionalID(dto.CategoryId)
	dto.Amount = money.Round2(dto.Amount)

	if dto.PaidAt == "" {
		dto.PaidAt = time.Now().UTC().Format(dateLayout)
	}
	if _, err := parseDateStrict(dto.PaidAt); err != nil {
		return err
	}
	if dto.Amount <= 0 {
		return fmt.Errorf("%w: jumlah pembayaran harus lebih dari 0", apperror.ErrValidation)
	}
	if dto.CreateTransaction && dto.WalletId == nil {
		return fmt.Errorf("%w: walletId wajib diisi untuk membuat transaksi", apperror.ErrValidation)
	}

	return nil
}

func enrichBill(dto *entity.BillDTO) {
	if dto == nil {
		return
	}

	dto.RemainingAmount = money.Round2(dto.Amount - dto.PaidAmount)
	dto.Status = billStatus(dto.Amount, dto.PaidAmount)

	dueDate, err := time.Parse(dateLayout, dto.DueDate)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dto.DaysUntilDue = int(dueDate.Sub(today).Hours() / 24)
	dto.IsOverdue = dto.Status != StatusPaid && dto.DaysUntilDue < 0
}

func billStatus(amount, paid float64) string {
	switch {
	case paid <= 0:
		return StatusUnpaid
	case paid >= amount:
		return StatusPaid
	default:
		return StatusPartial
	}
}

func validatePaymentAmount(amount, paid, total float64) error {
	if money.Round2(amount) <= 0 {
		return fmt.Errorf("%w: jumlah pembayaran harus lebih dari 0", apperror.ErrValidation)
	}
	if money.Round2(paid+amount) > total {
		return fmt.Errorf("%w: pembayaran melebihi sisa tagihan", apperror.ErrValidation)
	}
	return nil
}

func normalizeBillFilter(filter *entity.BillFilter) error {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	filter.Lender = strings.TrimSpace(filter.Lender)
	filter.Status = strings.ToLower(strings.TrimSpace(filter.Status))

	switch filter.Status {
	case "", StatusUnpaid, StatusPartial, StatusPaid:
	default:
		return fmt.Errorf("%w: status harus belum_bayar, sebagian, atau lunas", apperror.ErrValidation)
	}

	if filter.From != "" {
		if _, err := parseDateStrict(filter.From); err != nil {
			return err
		}
	}
	if filter.To != "" {
		if _, err := parseDateStrict(filter.To); err != nil {
			return err
		}
	}

	return nil
}

func parseDateStrict(value string) (time.Time, error) {
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: tanggal harus berformat yyyy-mm-dd", apperror.ErrValidation)
	}
	return parsed, nil
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
