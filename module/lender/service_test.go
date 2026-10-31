package lender

import (
	"strings"
	"testing"
)

func TestValidateLenderName(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid", "Kredivo", false},
		{"kosong", "", true},
		{"spasi saja", "   ", true},
		{"terlalu panjang", strings.Repeat("a", 256), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := strings.TrimSpace(tc.value)
			err := validateLenderName(value)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
