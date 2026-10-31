package budget

import (
	"testing"

	"expense.com/m/entity"
)

func TestNormalizeMonth(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{"valid", "2026-10", "2026-10", false},
		{"spasi", " 2026-01 ", "2026-01", false},
		{"format salah", "2026/10", "", true},
		{"tanggal lengkap", "2026-10-01", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeMonth(tc.value)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr && got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestMonthBounds(t *testing.T) {
	cases := []struct {
		month string
		from  string
		to    string
	}{
		{"2026-10", "2026-10-01", "2026-10-31"},
		{"2026-02", "2026-02-01", "2026-02-28"},
		{"2024-02", "2024-02-01", "2024-02-29"},
	}

	for _, tc := range cases {
		from, to := monthBounds(tc.month)
		if from != tc.from || to != tc.to {
			t.Errorf("monthBounds(%s) = %s..%s, want %s..%s", tc.month, from, to, tc.from, tc.to)
		}
	}
}

func TestBudgetStatus(t *testing.T) {
	cases := []struct {
		percent float64
		want    string
	}{
		{0, StatusSafe},
		{79.9, StatusSafe},
		{80, StatusWarning},
		{100, StatusWarning},
		{100.1, StatusExceeded},
	}

	for _, tc := range cases {
		if got := budgetStatus(tc.percent); got != tc.want {
			t.Errorf("budgetStatus(%v) = %s, want %s", tc.percent, got, tc.want)
		}
	}
}

func TestPercentage(t *testing.T) {
	if got := percentage(7300000, 10000000); got != 73 {
		t.Errorf("percentage = %v, want 73", got)
	}
	if got := percentage(100, 0); got != 0 {
		t.Errorf("percentage with zero amount = %v, want 0", got)
	}
}

func TestSpentIncludingChildren(t *testing.T) {
	spent := map[string]float64{
		"makan": 100000,
		"kopi":  50000,
		"luar":  25000,
	}
	children := map[string][]string{
		"makan": {"kopi", "luar"},
	}

	if got := spentIncludingChildren("makan", spent, children); got != 175000 {
		t.Fatalf("spent = %v, want 175000", got)
	}
	if got := spentIncludingChildren("kopi", spent, children); got != 50000 {
		t.Fatalf("spent child = %v, want 50000", got)
	}
}

func TestChildrenAndSumIndex(t *testing.T) {
	categories := []*entity.CategoryDTO{
		{CommonDTOFields: entity.CommonDTOFields{ID: "makan"}, Name: "Makan"},
		{CommonDTOFields: entity.CommonDTOFields{ID: "kopi"}, Name: "Kopi", ParentId: ptr("makan")},
	}

	children := childrenIndex(categories)
	if len(children["makan"]) != 1 || children["makan"][0] != "kopi" {
		t.Fatalf("children index salah: %v", children)
	}

	sums := sumIndex([]entity.CategorySumDTO{{CategoryId: "kopi", Total: 50000}})
	if sums["kopi"] != 50000 {
		t.Fatalf("sum index salah: %v", sums)
	}
}

func ptr(value string) *string {
	return &value
}
