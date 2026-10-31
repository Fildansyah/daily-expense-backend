package wallet

import (
	"testing"

	"expense.com/m/entity"
)

func newWallet(name, walletType string, balance float64, mutate func(*entity.WalletDTO)) *entity.WalletDTO {
	dto := &entity.WalletDTO{Name: name, Type: walletType, Balance: balance}
	if mutate != nil {
		mutate(dto)
	}
	return dto
}

func TestNormalizeWallet(t *testing.T) {
	dto := &entity.WalletDTO{Name: "  BCA  ", Type: "BANK"}
	normalizeWallet(dto)

	if dto.Name != "BCA" || dto.Type != "bank" {
		t.Fatalf("normalisasi gagal: %+v", dto)
	}
	if dto.Currency != "IDR" || dto.Color == "" || dto.Icon == "" {
		t.Fatalf("default tidak terisi: %+v", dto)
	}
}

func TestValidateWallet(t *testing.T) {
	limit := 25000000.0
	invalidDay := 31

	cases := []struct {
		name    string
		dto     *entity.WalletDTO
		wantErr bool
	}{
		{"valid bank", &entity.WalletDTO{Name: "BCA", Type: "bank"}, false},
		{"valid kartu kredit", &entity.WalletDTO{Name: "CC", Type: "credit_card", CreditLimit: &limit}, false},
		{"tanpa nama", &entity.WalletDTO{Type: "bank"}, true},
		{"tipe tidak dikenal", &entity.WalletDTO{Name: "X", Type: "saham"}, true},
		{"kartu kredit tanpa limit", &entity.WalletDTO{Name: "CC", Type: "credit_card"}, true},
		{"statement day tidak valid", &entity.WalletDTO{Name: "CC", Type: "credit_card", CreditLimit: &limit, StatementDay: &invalidDay}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWallet(tc.dto)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestApplyWalletMetrics(t *testing.T) {
	limit := 25000000.0
	cc := &entity.WalletDTO{Type: "credit_card", Balance: 7850000, CreditLimit: &limit}
	applyWalletMetrics(cc)

	if cc.Used != 7850000 || cc.Available != 17150000 || cc.Utilization != 31.4 {
		t.Fatalf("metrics salah: used=%v available=%v utilization=%v", cc.Used, cc.Available, cc.Utilization)
	}

	bank := &entity.WalletDTO{Type: "bank", Balance: 1000000}
	applyWalletMetrics(bank)
	if bank.Used != 0 || bank.Available != 0 || bank.Utilization != 0 {
		t.Fatal("aset tidak boleh punya metrik kredit")
	}
}

func TestSummarizeWallets(t *testing.T) {
	limit := 25000000.0
	wallets := []*entity.WalletDTO{
		newWallet("BCA", "bank", 12450000, nil),
		newWallet("Dompet", "cash", 850000, nil),
		newWallet("GoPay", "ewallet", 320000, nil),
		newWallet("Kartu Kredit BCA", "credit_card", 7850000, func(w *entity.WalletDTO) {
			w.CreditLimit = &limit
		}),
		newWallet("Akun Lama", "bank", 5000000, func(w *entity.WalletDTO) {
			w.IsArchived = true
		}),
		newWallet("Dana Khusus", "bank", 2000000, func(w *entity.WalletDTO) {
			w.IsExclude = true
		}),
	}

	summary := summarizeWallets(wallets)

	if summary.TotalBalance != 13620000 {
		t.Fatalf("total saldo = %v, want 13620000", summary.TotalBalance)
	}
	if summary.CreditUsed != 7850000 {
		t.Fatalf("credit used = %v, want 7850000", summary.CreditUsed)
	}
	if summary.CreditLimit != 25000000 {
		t.Fatalf("credit limit = %v, want 25000000", summary.CreditLimit)
	}
	if summary.Available != 17150000 {
		t.Fatalf("available = %v, want 17150000", summary.Available)
	}
}
