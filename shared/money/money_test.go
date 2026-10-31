package money

import "testing"

func TestRound2(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{1000.004, 1000},
		{1000.005, 1000.01},
		{1250.555, 1250.56},
		{-1250.555, -1250.56},
		{0, 0},
	}

	for _, tc := range cases {
		if got := Round2(tc.in); got != tc.want {
			t.Errorf("Round2(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestSignedDelta(t *testing.T) {
	cases := []struct {
		name     string
		wallet   string
		out      bool
		amount   float64
		expected float64
	}{
		{"aset uang keluar", "bank", true, 50000, -50000},
		{"aset uang masuk", "bank", false, 50000, 50000},
		{"kartu kredit belanja", "credit_card", true, 50000, 50000},
		{"kartu kredit pembayaran", "credit_card", false, 50000, -50000},
		{"paylater belanja", "paylater", true, 50000, 50000},
		{"tunai uang masuk", "cash", false, 50000, 50000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SignedDelta(tc.wallet, tc.out, tc.amount); got != tc.expected {
				t.Errorf("SignedDelta(%s, out=%v, %v) = %v, want %v", tc.wallet, tc.out, tc.amount, got, tc.expected)
			}
		})
	}
}

func TestTransferKeepsTotal(t *testing.T) {
	balances := map[string]float64{
		"bank": 1000000,
		"cash": 0,
	}

	const amount = 250000
	balances["bank"] = Round2(balances["bank"] + SignedDelta("bank", true, amount))
	balances["cash"] = Round2(balances["cash"] + SignedDelta("cash", false, amount))

	if balances["bank"] != 750000 {
		t.Fatalf("saldo bank = %v, want 750000", balances["bank"])
	}
	if balances["cash"] != 250000 {
		t.Fatalf("saldo cash = %v, want 250000", balances["cash"])
	}
	if total := balances["bank"] + balances["cash"]; total != 1000000 {
		t.Fatalf("total saldo berubah menjadi %v, want 1000000", total)
	}
}

func TestCreditCardPaymentReducesDebt(t *testing.T) {
	// Belanja 300rb dengan kartu kredit menambah utang.
	debt := SignedDelta("credit_card", true, 300000)
	if debt != 300000 {
		t.Fatalf("utang setelah belanja = %v, want 300000", debt)
	}

	// Bayar dari rekening bank: saldo bank berkurang, utang berkurang.
	debt = Round2(debt + SignedDelta("credit_card", false, 300000))
	if debt != 0 {
		t.Fatalf("utang setelah dibayar = %v, want 0", debt)
	}
}
