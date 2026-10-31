package wallet

import (
	"context"
	"fmt"
	"math"
	"strings"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/money"
	"expense.com/m/shared/tx"
)

var allowedWalletTypes = map[string]bool{
	"cash":        true,
	"bank":        true,
	"ewallet":     true,
	"credit_card": true,
	"paylater":    true,
	"investment":  true,
}

type WalletService struct {
	repo *WalletRepository
	txm  *tx.Manager
}

func NewWalletService(repo *WalletRepository, txm *tx.Manager) *WalletService {
	return &WalletService{repo: repo, txm: txm}
}

func (s *WalletService) CreateWallet(ctx context.Context, dto *entity.WalletDTO) (*entity.WalletDTO, error) {
	normalizeWallet(dto)
	if err := validateWallet(dto); err != nil {
		return nil, err
	}

	// Dukung payload lama yang hanya mengirim "balance" sebagai saldo awal.
	if dto.OpeningBalance == 0 && dto.Balance != 0 {
		dto.OpeningBalance = dto.Balance
	}
	dto.Balance = dto.OpeningBalance

	created, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	applyWalletMetrics(created)
	return created, nil
}

func (s *WalletService) UpdateWallet(ctx context.Context, dto *entity.WalletDTO) (*entity.WalletDTO, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("%w: id akun wajib diisi", apperror.ErrValidation)
	}
	normalizeWallet(dto)
	if err := validateWallet(dto); err != nil {
		return nil, err
	}

	var out *entity.WalletDTO
	err := s.txm.Run(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindByIDForUpdate(ctx, dto.ID, dto.UserID)
		if err != nil {
			return err
		}

		// Saldo awal adalah satu-satunya input yang boleh menggeser saldo berjalan.
		delta := money.Round2(dto.OpeningBalance - existing.OpeningBalance)

		if err := s.repo.UpdateFields(ctx, dto); err != nil {
			return err
		}
		if delta != 0 {
			if err := s.repo.AdjustBalance(ctx, dto.UserID, dto.ID, delta); err != nil {
				return err
			}
		}

		out, err = s.repo.FindByID(ctx, dto.ID, dto.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}

	applyWalletMetrics(out)
	return out, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id, userID string) (*entity.WalletDTO, error) {
	wallet, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	applyWalletMetrics(wallet)
	return wallet, nil
}

func (s *WalletService) FindWalletByUserID(ctx context.Context, userID string, includeArchived bool) ([]*entity.WalletDTO, error) {
	wallets, err := s.repo.FindByUser(ctx, userID, includeArchived)
	if err != nil {
		return nil, err
	}
	for _, wallet := range wallets {
		applyWalletMetrics(wallet)
	}
	return wallets, nil
}

func (s *WalletService) DeleteWallet(ctx context.Context, id, userID string, archive bool) error {
	if id == "" {
		return fmt.Errorf("%w: id akun wajib diisi", apperror.ErrValidation)
	}

	if archive {
		return s.repo.Archive(ctx, id, userID)
	}

	count, err := s.repo.CountTransactions(ctx, userID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%w: akun masih memiliki %d transaksi, gunakan arsip", apperror.ErrConflict, count)
	}

	return s.repo.Delete(ctx, id, userID)
}

func (s *WalletService) GetSummary(ctx context.Context, userID string) (*entity.WalletSummaryDTO, error) {
	wallets, err := s.repo.FindByUser(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	return summarizeWallets(wallets), nil
}

func normalizeWallet(dto *entity.WalletDTO) {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Type = strings.ToLower(strings.TrimSpace(dto.Type))
	dto.Currency = strings.ToUpper(strings.TrimSpace(dto.Currency))
	if dto.Currency == "" {
		dto.Currency = "IDR"
	}
	if dto.Color == "" {
		dto.Color = "#0F766E"
	}
	if dto.Icon == "" {
		dto.Icon = "wallet"
	}
	dto.AccountNumber = strings.TrimSpace(dto.AccountNumber)
	dto.Balance = money.Round2(dto.Balance)
	dto.OpeningBalance = money.Round2(dto.OpeningBalance)
	if dto.InterestRate != nil {
		value := money.Round2(*dto.InterestRate)
		dto.InterestRate = &value
	}
}

func validateWallet(dto *entity.WalletDTO) error {
	if dto.Name == "" {
		return fmt.Errorf("%w: nama akun wajib diisi", apperror.ErrValidation)
	}
	if !allowedWalletTypes[dto.Type] {
		return fmt.Errorf("%w: tipe akun tidak dikenal", apperror.ErrValidation)
	}
	if dto.Type == "credit_card" && (dto.CreditLimit == nil || *dto.CreditLimit <= 0) {
		return fmt.Errorf("%w: creditLimit wajib diisi untuk kartu kredit", apperror.ErrValidation)
	}
	if dto.CreditLimit != nil && *dto.CreditLimit < 0 {
		return fmt.Errorf("%w: creditLimit tidak boleh negatif", apperror.ErrValidation)
	}
	if dto.StatementDay != nil && (*dto.StatementDay < 1 || *dto.StatementDay > 28) {
		return fmt.Errorf("%w: statementDay harus antara 1 dan 28", apperror.ErrValidation)
	}
	if dto.DueDay != nil && (*dto.DueDay < 1 || *dto.DueDay > 28) {
		return fmt.Errorf("%w: dueDay harus antara 1 dan 28", apperror.ErrValidation)
	}
	return nil
}

func applyWalletMetrics(dto *entity.WalletDTO) {
	if dto == nil {
		return
	}

	if !money.IsLiability(dto.Type) {
		dto.Used, dto.Available, dto.Utilization = 0, 0, 0
		return
	}

	dto.Used = money.Round2(dto.Balance)
	if dto.CreditLimit == nil || *dto.CreditLimit <= 0 {
		dto.Available, dto.Utilization = 0, 0
		return
	}

	dto.Available = money.Round2(*dto.CreditLimit - dto.Used)
	dto.Utilization = math.Round(dto.Used / *dto.CreditLimit * 1000) / 10
}

func summarizeWallets(wallets []*entity.WalletDTO) *entity.WalletSummaryDTO {
	summary := &entity.WalletSummaryDTO{ByType: make([]entity.WalletTypeSummaryDTO, 0)}
	totals := map[string]float64{}
	order := make([]string, 0)

	for _, wallet := range wallets {
		if wallet == nil || wallet.IsArchived {
			continue
		}

		if _, seen := totals[wallet.Type]; !seen {
			order = append(order, wallet.Type)
		}
		totals[wallet.Type] = money.Round2(totals[wallet.Type] + wallet.Balance)

		if wallet.IsExclude {
			continue
		}

		if money.IsLiability(wallet.Type) {
			summary.CreditUsed = money.Round2(summary.CreditUsed + wallet.Balance)
			if wallet.CreditLimit != nil {
				summary.CreditLimit = money.Round2(summary.CreditLimit + *wallet.CreditLimit)
			}
			continue
		}
		summary.TotalBalance = money.Round2(summary.TotalBalance + wallet.Balance)
	}

	summary.Available = money.Round2(summary.CreditLimit - summary.CreditUsed)
	for _, walletType := range order {
		summary.ByType = append(summary.ByType, entity.WalletTypeSummaryDTO{Type: walletType, Total: totals[walletType]})
	}

	return summary
}
