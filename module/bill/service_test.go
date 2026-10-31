package bill

import (
	"testing"
	"time"

	"expense.com/m/entity"
)

func TestBillStatus(t *testing.T) {
	cases := []struct {
		amount float64
		paid   float64
		want   string
	}{
		{100000, 0, StatusUnpaid},
		{100000, -10, StatusUnpaid},
		{100000, 50000, StatusPartial},
		{100000, 100000, StatusPaid},
		{100000, 120000, StatusPaid},
	}

	for _, tc := range cases {
		if got := billStatus(tc.amount, tc.paid); got != tc.want {
			t.Errorf("billStatus(%v, %v) = %s, want %s", tc.amount, tc.paid, got, tc.want)
		}
	}
}

func TestEnrichBill(t *testing.T) {
	today := time.Now().UTC()

	cases := []struct {
		name        string
		dueDate     string
		amount      float64
		paid        float64
		wantStatus  string
		wantOverdue bool
		wantDays    int
		wantRemains float64
	}{
		{
			name:        "jatuh tempo kemarin dan belum dibayar",
			dueDate:     today.AddDate(0, 0, -1).Format(dateLayout),
			amount:      100000,
			paid:        0,
			wantStatus:  StatusUnpaid,
			wantOverdue: true,
			wantDays:    -1,
			wantRemains: 100000,
		},
		{
			name:        "jatuh tempo hari ini sebagian",
			dueDate:     today.Format(dateLayout),
			amount:      100000,
			paid:        40000,
			wantStatus:  StatusPartial,
			wantOverdue: false,
			wantDays:    0,
			wantRemains: 60000,
		},
		{
			name:        "lunas walau sudah lewat",
			dueDate:     today.AddDate(0, 0, -3).Format(dateLayout),
			amount:      100000,
			paid:        100000,
			wantStatus:  StatusPaid,
			wantOverdue: false,
			wantDays:    -3,
			wantRemains: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dto := &entity.BillDTO{
				Amount:     tc.amount,
				PaidAmount: tc.paid,
				DueDate:    tc.dueDate,
			}
			enrichBill(dto)

			if dto.Status != tc.wantStatus {
				t.Errorf("status = %s, want %s", dto.Status, tc.wantStatus)
			}
			if dto.IsOverdue != tc.wantOverdue {
				t.Errorf("isOverdue = %v, want %v", dto.IsOverdue, tc.wantOverdue)
			}
			if dto.DaysUntilDue != tc.wantDays {
				t.Errorf("daysUntilDue = %d, want %d", dto.DaysUntilDue, tc.wantDays)
			}
			if dto.RemainingAmount != tc.wantRemains {
				t.Errorf("remaining = %v, want %v", dto.RemainingAmount, tc.wantRemains)
			}
		})
	}
}

func TestPrepareBill(t *testing.T) {
	today := time.Now().UTC().Format(dateLayout)
	installmentTotal := 12
	installmentIndex := 13
	interest := -1.0

	cases := []struct {
		name    string
		dto     *entity.BillDTO
		wantErr bool
	}{
		{"valid", &entity.BillDTO{LenderId: "Kredivo", Amount: 100000, DueDate: today}, false},
		{"lender kosong", &entity.BillDTO{LenderId: "  ", Amount: 100000, DueDate: today}, true},
		{"amount nol", &entity.BillDTO{LenderId: "Kredivo", Amount: 0, DueDate: today}, true},
		{"paid negatif", &entity.BillDTO{LenderId: "Kredivo", Amount: 100000, PaidAmount: -1, DueDate: today}, true},
		{"paid melebihi amount", &entity.BillDTO{LenderId: "Kredivo", Amount: 100000, PaidAmount: 120000, DueDate: today}, true},
		{"tanggal salah", &entity.BillDTO{LenderId: "Kredivo", Amount: 100000, DueDate: "31-12-2026"}, true},
		{"index cicilan melebihi total", &entity.BillDTO{
			LenderId: "Kredivo", Amount: 100000, DueDate: today,
			InstallmentTotal: &installmentTotal, InstallmentIndex: &installmentIndex,
		}, true},
		{"bunga negatif", &entity.BillDTO{LenderId: "Kredivo", Amount: 100000, DueDate: today, InterestRate: &interest}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := prepareBill(tc.dto)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePaymentAmount(t *testing.T) {
	cases := []struct {
		name    string
		amount  float64
		paid    float64
		total   float64
		wantErr bool
	}{
		{"pembayaran sebagian", 50000, 0, 100000, false},
		{"pelunasan", 60000, 40000, 100000, false},
		{"melebihi sisa", 70000, 40000, 100000, true},
		{"nol", 0, 0, 100000, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePaymentAmount(tc.amount, tc.paid, tc.total)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPreparePayment(t *testing.T) {
	walletID := "w-1"
	cases := []struct {
		name    string
		dto     *entity.BillPaymentDTO
		wantErr bool
	}{
		{"valid", &entity.BillPaymentDTO{Amount: 50000}, false},
		{"amount nol", &entity.BillPaymentDTO{Amount: 0}, true},
		{"tanggal salah", &entity.BillPaymentDTO{Amount: 50000, PaidAt: "2026/10/31"}, true},
		{"buat transaksi tanpa wallet", &entity.BillPaymentDTO{Amount: 50000, CreateTransaction: true}, true},
		{"buat transaksi dengan wallet", &entity.BillPaymentDTO{Amount: 50000, CreateTransaction: true, WalletId: &walletID}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := preparePayment(tc.dto)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
