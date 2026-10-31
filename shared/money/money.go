package money

import "math"

// Round2 membulatkan nominal uang ke 2 desimal.
func Round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// IsLiability menandai tipe akun yang saldonya merepresentasikan utang,
// bukan aset. Untuk akun seperti ini saldo positif berarti kewajiban.
func IsLiability(walletType string) bool {
	switch walletType {
	case "credit_card", "paylater":
		return true
	default:
		return false
	}
}

// SignedDelta mengubah arah arus uang menjadi perubahan saldo akun.
// out=true berarti uang keluar dari akun tersebut.
//
// Aset: uang keluar mengurangi saldo, uang masuk menambah saldo.
// Liabilitas: uang keluar menambah utang, uang masuk mengurangi utang.
func SignedDelta(walletType string, out bool, amount float64) float64 {
	sign := 1.0
	if out {
		sign = -1
	}
	if IsLiability(walletType) {
		sign = -sign
	}
	return Round2(sign * amount)
}
