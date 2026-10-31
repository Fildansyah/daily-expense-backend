package report

import (
	"testing"
	"time"

	"expense.com/m/entity"
)

func TestResolveRange(t *testing.T) {
	from, to, err := resolveRange("2026-10-01", "2026-10-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if from.Format(dateLayout) != "2026-10-01" || to.Format(dateLayout) != "2026-10-31" {
		t.Fatalf("range salah: %s..%s", from, to)
	}

	if _, _, err := resolveRange("2026-11-01", "2026-10-01"); err == nil {
		t.Fatal("from setelah to harus error")
	}
	if _, _, err := resolveRange("31-10-2026", "2026-10-31"); err == nil {
		t.Fatal("format salah harus error")
	}

	// Default: bulan berjalan.
	from, to, err = resolveRange("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if from.Day() != 1 {
		t.Fatalf("default from harus tanggal 1, got %s", from)
	}
	if to.AddDate(0, 0, 1).Day() != 1 {
		t.Fatalf("default to harus akhir bulan, got %s", to)
	}
}

func TestPreviousRange(t *testing.T) {
	from := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)

	prevFrom, prevTo := previousRange(from, to)
	// Periode pembanding memakai panjang hari yang sama (31 hari).
	if prevFrom.Format(dateLayout) != "2026-08-31" || prevTo.Format(dateLayout) != "2026-09-30" {
		t.Fatalf("previous range salah: %s..%s", prevFrom, prevTo)
	}
}

func TestPercentChange(t *testing.T) {
	cases := []struct {
		current  float64
		previous float64
		want     float64
	}{
		{12500000, 10000000, 25},
		{8000000, 10000000, -20},
		{100, 0, 100},
		{0, 0, 0},
	}
	for _, tc := range cases {
		if got := percentChange(tc.current, tc.previous); got != tc.want {
			t.Errorf("percentChange(%v,%v) = %v, want %v", tc.current, tc.previous, got, tc.want)
		}
	}
}

func TestBuildCategoryReportRollup(t *testing.T) {
	tree := []*entity.CategoryDTO{
		{
			CommonDTOFields: entity.CommonDTOFields{ID: "makan"},
			Name:            "Makan",
			Children: []*entity.CategoryDTO{
				{CommonDTOFields: entity.CommonDTOFields{ID: "kopi"}, Name: "Kopi"},
			},
		},
		{
			CommonDTOFields: entity.CommonDTOFields{ID: "transport"},
			Name:            "Transportasi",
		},
	}

	sums := map[string]categorySum{
		"makan":     {total: 100000, count: 2},
		"kopi":      {total: 50000, count: 3},
		"transport": {total: 25000, count: 1},
	}

	nodes := buildCategoryReport(tree, sums)
	if len(nodes) != 2 {
		t.Fatalf("jumlah root = %d, want 2", len(nodes))
	}

	// Makan = 100.000 + 50.000 (Kopi)
	if nodes[0].CategoryId != "makan" || nodes[0].Amount != 150000 {
		t.Fatalf("rollup salah: %+v", nodes[0])
	}
	if nodes[0].Count != 5 {
		t.Fatalf("count rollup = %d, want 5", nodes[0].Count)
	}
	if len(nodes[0].Children) != 1 || nodes[0].Children[0].Amount != 50000 {
		t.Fatalf("child salah: %+v", nodes[0].Children)
	}

	applyPercentage(nodes, 175000)
	if nodes[0].Percentage != 85.7 {
		t.Fatalf("percentage = %v, want 85.7", nodes[0].Percentage)
	}
}

func TestCollapseCategories(t *testing.T) {
	nodes := []entity.ReportCategoryDTO{
		{CategoryId: "a", Amount: 100, Percentage: 10},
		{CategoryId: "b", Amount: 90, Percentage: 9},
		{CategoryId: "c", Amount: 80, Percentage: 8},
		{CategoryId: "d", Amount: 70, Percentage: 7},
		{CategoryId: "e", Amount: 60, Percentage: 6},
		{CategoryId: "f", Amount: 50, Percentage: 5},
		{CategoryId: "g", Amount: 40, Percentage: 4},
	}

	result := collapseCategories(nodes, 5)
	if len(result) != 6 {
		t.Fatalf("len = %d, want 6 (5 + Lainnya)", len(result))
	}
	if result[5].Name != "Lainnya" || result[5].Amount != 90 {
		t.Fatalf("other salah: %+v", result[5])
	}
}

func TestTrendPointsOnlyExpense(t *testing.T) {
	rows := []entity.TransactionSeriesDTO{
		{Period: "2026-09-01", Type: "INCOME", Total: 5000000},
		{Period: "2026-09-01", Type: "EXPENSE", Total: 100000},
		{Period: "2026-10-01", Type: "EXPENSE", Total: 250000},
	}

	points := trendPoints(rows)
	if len(points) != 2 {
		t.Fatalf("len = %d, want 2", len(points))
	}
	if points[0].Value != 100000 || points[1].Value != 250000 {
		t.Fatalf("points salah: %+v", points)
	}
}

func TestCreditUtilization(t *testing.T) {
	limit := 10000000.0
	wallets := []*entity.WalletDTO{
		{CommonDTOFields: entity.CommonDTOFields{ID: "w1"}, Name: "BCA", Type: "bank", Balance: 1000000},
		{
			CommonDTOFields: entity.CommonDTOFields{ID: "w2"},
			Name:            "Kartu Kredit",
			Type:            "credit_card",
			Balance:         2500000,
			CreditLimit:     &limit,
			Used:            2500000,
			Available:       7500000,
			Utilization:     25,
		},
	}

	result := creditUtilization(wallets)
	if len(result) != 1 || result[0].Name != "Kartu Kredit" || result[0].Utilization != 25 {
		t.Fatalf("credit utilization salah: %+v", result)
	}
}
