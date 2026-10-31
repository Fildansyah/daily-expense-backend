package transaction

import (
	"testing"
	"time"

	"expense.com/m/entity"
	"expense.com/m/shared/money"
)

func ptr(value string) *string {
	return &value
}

func validExpense() *entity.TransactionDTO {
	return &entity.TransactionDTO{
		Type:       TypeExpense,
		WalletId:   "w-bank",
		CategoryId: ptr("c-food"),
		Amount:     50000,
		Memo:       " Makan siang ",
		Date:       time.Now().UTC().Format(dateLayout),
	}
}

func TestPrepareTransaction(t *testing.T) {
	dto := &entity.TransactionDTO{
		Type:       " expense ",
		WalletId:   "w-bank",
		CategoryId: ptr("c-food"),
		Amount:     1000.005,
	}

	if err := prepareTransaction(dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Type != TypeExpense {
		t.Fatalf("type = %s, want EXPENSE", dto.Type)
	}
	if dto.Amount != 1000.01 {
		t.Fatalf("amount = %v, want 1000.01", dto.Amount)
	}
	if dto.Date != time.Now().UTC().Format(dateLayout) {
		t.Fatalf("date default = %s", dto.Date)
	}
}

func TestValidateTransaction(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*entity.TransactionDTO)
		wantErr bool
	}{
		{"pengeluaran valid", func(*entity.TransactionDTO) {}, false},
		{"pemasukan valid", func(d *entity.TransactionDTO) { d.Type = TypeIncome }, false},
		{"transfer valid", func(d *entity.TransactionDTO) {
			d.Type = TypeTransfer
			d.CategoryId = nil
			d.TransferWalletId = ptr("w-cash")
		}, false},
		{"amount nol", func(d *entity.TransactionDTO) { d.Amount = 0 }, true},
		{"amount negatif", func(d *entity.TransactionDTO) { d.Amount = -1 }, true},
		{"type tidak dikenal", func(d *entity.TransactionDTO) { d.Type = "UNKNOWN" }, true},
		{"tanpa wallet", func(d *entity.TransactionDTO) { d.WalletId = "" }, true},
		{"pengeluaran tanpa kategori", func(d *entity.TransactionDTO) { d.CategoryId = nil }, true},
		{"transfer tanpa tujuan", func(d *entity.TransactionDTO) { d.Type = TypeTransfer; d.CategoryId = nil }, true},
		{"transfer ke akun yang sama", func(d *entity.TransactionDTO) {
			d.Type = TypeTransfer
			d.CategoryId = nil
			d.TransferWalletId = ptr("w-bank")
		}, true},
		{"tanggal salah format", func(d *entity.TransactionDTO) { d.Date = "28-10-2026" }, true},
		{"tanggal terlalu jauh", func(d *entity.TransactionDTO) {
			d.Date = time.Now().UTC().AddDate(2, 0, 0).Format(dateLayout)
		}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dto := validExpense()
			tc.mutate(dto)
			err := validateTransaction(dto)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBuildDeltas(t *testing.T) {
	walletTypes := map[string]string{
		"w-bank": "bank",
		"w-cash": "cash",
		"w-cc":   "credit_card",
	}

	cases := []struct {
		name            string
		transactionType string
		walletID        string
		transferWallet  *string
		amount          float64
		expect          map[string]float64
	}{
		{
			name:            "pengeluaran dari bank mengurangi saldo",
			transactionType: TypeExpense,
			walletID:        "w-bank",
			amount:          50000,
			expect:          map[string]float64{"w-bank": -50000},
		},
		{
			name:            "pemasukan ke bank menambah saldo",
			transactionType: TypeIncome,
			walletID:        "w-bank",
			amount:          12500000,
			expect:          map[string]float64{"w-bank": 12500000},
		},
		{
			name:            "pengeluaran kartu kredit menambah utang",
			transactionType: TypeExpense,
			walletID:        "w-cc",
			amount:          300000,
			expect:          map[string]float64{"w-cc": 300000},
		},
		{
			name:            "transfer mengurangi asal dan menambah tujuan",
			transactionType: TypeTransfer,
			walletID:        "w-bank",
			transferWallet:  ptr("w-cash"),
			amount:          200000,
			expect:          map[string]float64{"w-bank": -200000, "w-cash": 200000},
		},
		{
			name:            "transfer ke kartu kredit mengurangi utang",
			transactionType: TypeTransfer,
			walletID:        "w-bank",
			transferWallet:  ptr("w-cc"),
			amount:          300000,
			expect:          map[string]float64{"w-bank": -300000, "w-cc": -300000},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildDeltas(tc.transactionType, tc.walletID, tc.transferWallet, tc.amount, walletTypes)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.expect) {
				t.Fatalf("got %d deltas, want %d (%v)", len(got), len(tc.expect), got)
			}
			for walletID, want := range tc.expect {
				if got[walletID] != want {
					t.Errorf("delta %s = %v, want %v", walletID, got[walletID], want)
				}
			}
		})
	}
}

func TestBuildDeltasUnknownWallet(t *testing.T) {
	if _, err := buildDeltas(TypeExpense, "w-missing", nil, 1000, map[string]string{}); err == nil {
		t.Fatal("expected error for unknown wallet")
	}
}

func TestValidateBulk(t *testing.T) {
	cases := []struct {
		name    string
		bulk    *entity.BulkTransactionDTO
		wantErr bool
	}{
		{"delete valid", &entity.BulkTransactionDTO{IDs: []string{"a"}, Action: "delete"}, false},
		{"categorize tanpa value", &entity.BulkTransactionDTO{IDs: []string{"a"}, Action: "categorize"}, true},
		{"tanpa ids", &entity.BulkTransactionDTO{Action: "delete"}, true},
		{"aksi tidak dikenal", &entity.BulkTransactionDTO{IDs: []string{"a"}, Action: "archive"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBulk(tc.bulk)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// ledger mensimulasikan urutan operasi CRUD memakai buildDeltas yang sama
// dengan yang dipakai service.
type ledger struct {
	balances map[string]float64
	types    map[string]string
}

func newLedger() *ledger {
	return &ledger{
		balances: map[string]float64{"w-bank": 1000000, "w-cash": 0, "w-cc": 0},
		types:    map[string]string{"w-bank": "bank", "w-cash": "cash", "w-cc": "credit_card"},
	}
}

func (l *ledger) apply(transactionType, walletID string, transferWallet *string, amount float64, reverse bool) error {
	deltas, err := buildDeltas(transactionType, walletID, transferWallet, amount, l.types)
	if err != nil {
		return err
	}
	for walletID, delta := range deltas {
		if reverse {
			delta = -delta
		}
		l.balances[walletID] = money.Round2(l.balances[walletID] + delta)
	}
	return nil
}

func TestLedgerCreateUpdateDeleteIsConsistent(t *testing.T) {
	l := newLedger()

	// Create expense 50.000 dari bank.
	if err := l.apply(TypeExpense, "w-bank", nil, 50000, false); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-bank"] != 950000 {
		t.Fatalf("setelah create bank = %v, want 950000", l.balances["w-bank"])
	}

	// Update nominal menjadi 75.000: reverse lama lalu apply baru.
	if err := l.apply(TypeExpense, "w-bank", nil, 50000, true); err != nil {
		t.Fatal(err)
	}
	if err := l.apply(TypeExpense, "w-bank", nil, 75000, false); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-bank"] != 925000 {
		t.Fatalf("setelah update bank = %v, want 925000", l.balances["w-bank"])
	}

	// Pindah akun ke cash: reverse bank, apply cash.
	if err := l.apply(TypeExpense, "w-bank", nil, 75000, true); err != nil {
		t.Fatal(err)
	}
	if err := l.apply(TypeExpense, "w-cash", nil, 75000, false); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-bank"] != 1000000 {
		t.Fatalf("setelah pindah akun bank = %v, want 1000000", l.balances["w-bank"])
	}
	if l.balances["w-cash"] != -75000 {
		t.Fatalf("setelah pindah akun cash = %v, want -75000", l.balances["w-cash"])
	}

	// Delete: kembalikan efek terakhir.
	if err := l.apply(TypeExpense, "w-cash", nil, 75000, true); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-bank"] != 1000000 || l.balances["w-cash"] != 0 {
		t.Fatalf("saldo tidak kembali setelah delete: %v", l.balances)
	}
}

func TestLedgerTransferKeepsTotal(t *testing.T) {
	l := newLedger()

	before := l.balances["w-bank"] + l.balances["w-cash"]
	if err := l.apply(TypeTransfer, "w-bank", ptr("w-cash"), 400000, false); err != nil {
		t.Fatal(err)
	}
	after := l.balances["w-bank"] + l.balances["w-cash"]

	if before != after {
		t.Fatalf("total aset berubah dari %v menjadi %v", before, after)
	}
	if l.balances["w-bank"] != 600000 || l.balances["w-cash"] != 400000 {
		t.Fatalf("distribusi saldo salah: %v", l.balances)
	}

	if err := l.apply(TypeTransfer, "w-bank", ptr("w-cash"), 400000, true); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-bank"] != 1000000 || l.balances["w-cash"] != 0 {
		t.Fatalf("saldo tidak kembali setelah delete transfer: %v", l.balances)
	}
}

func TestLedgerCreditCardLifecycle(t *testing.T) {
	l := newLedger()

	if err := l.apply(TypeExpense, "w-cc", nil, 300000, false); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-cc"] != 300000 {
		t.Fatalf("utang kartu kredit = %v, want 300000", l.balances["w-cc"])
	}

	if err := l.apply(TypeTransfer, "w-bank", ptr("w-cc"), 300000, false); err != nil {
		t.Fatal(err)
	}
	if l.balances["w-cc"] != 0 {
		t.Fatalf("utang setelah dibayar = %v, want 0", l.balances["w-cc"])
	}
	if l.balances["w-bank"] != 700000 {
		t.Fatalf("saldo bank setelah bayar = %v, want 700000", l.balances["w-bank"])
	}
}
