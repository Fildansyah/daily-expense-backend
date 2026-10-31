package entity

import "time"

type BillDTO struct {
	CommonDTOFields

	UpdatedAt  time.Time  `json:"updatedAt"`
	UserID     string     `json:"-"`
	LenderId   string     `json:"lenderId"`
	Lender     *LenderDTO `json:"lender,omitempty"`
	Amount     float64    `json:"amount"`
	PaidAmount float64    `json:"paidAmount"`
	DueDate    string     `json:"dueDate"`
	Deskripsi  string     `json:"deskripsi"`

	WalletId   *string `json:"walletId,omitempty"`
	CategoryId *string `json:"categoryId,omitempty"`

	InstallmentTotal *int     `json:"installmentTotal,omitempty"`
	InstallmentIndex *int     `json:"installmentIndex,omitempty"`
	InterestRate     *float64 `json:"interestRate,omitempty"`

	RemindH3  bool   `json:"remindH3"`
	RemindH1  bool   `json:"remindH1"`
	RemindDue bool   `json:"remindDue"`
	Currency  string `json:"currency"`

	// Field turunan yang dihitung service.
	RemainingAmount float64           `json:"remainingAmount"`
	Status          string            `json:"status"`
	IsOverdue       bool              `json:"isOverdue"`
	DaysUntilDue    int               `json:"daysUntilDue"`
	Wallet          *WalletBriefDTO   `json:"wallet,omitempty"`
	Category        *CategoryBriefDTO `json:"category,omitempty"`
}

type BillFilter struct {
	From     string
	To       string
	LenderId string
	Status   string
	WalletId string
	Q        string
	Page     int
	Limit    int
	Sort     string
}

type BillPaymentDTO struct {
	CommonDTOFields

	BillId        string    `json:"billId"`
	UserID        string    `json:"-"`
	Amount        float64   `json:"amount"`
	PaidAt        string    `json:"paidAt"`
	WalletId      *string   `json:"walletId,omitempty"`
	Note          string    `json:"note"`
	TransactionId *string   `json:"transactionId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`

	// Field khusus request: buat transaksi pengeluaran otomatis.
	CreateTransaction bool    `json:"createTransaction"`
	CategoryId        *string `json:"categoryId,omitempty"`
}

type BillLenderSummaryDTO struct {
	LenderId  string     `json:"lenderId"`
	Lender    *LenderDTO `json:"lender"`
	Count     int64      `json:"count"`
	Total     float64    `json:"total"`
	Paid      float64    `json:"paid"`
	Remaining float64    `json:"remaining"`
}

type BillStatusSummaryDTO struct {
	Status    string  `json:"status"`
	Count     int64   `json:"count"`
	Total     float64 `json:"total"`
	Paid      float64 `json:"paid"`
	Remaining float64 `json:"remaining"`
}

type BillWalletSummaryDTO struct {
	WalletId   string  `json:"walletId,omitempty"`
	WalletName string  `json:"walletName,omitempty"`
	Count      int64   `json:"count"`
	Total      float64 `json:"total"`
	Paid       float64 `json:"paid"`
	Remaining  float64 `json:"remaining"`
}

type BillSummaryDTO struct {
	Total        float64                `json:"total"`
	Paid         float64                `json:"paid"`
	Remaining    float64                `json:"remaining"`
	RemainingAll float64                `json:"remainingAll"`
	OverdueTotal float64                `json:"overdueTotal"`
	DueThisWeek  float64                `json:"dueThisWeek"`
	PerLender    []BillLenderSummaryDTO `json:"perLender"`
	ByStatus     []BillStatusSummaryDTO `json:"byStatus"`
	ByWallet     []BillWalletSummaryDTO `json:"byWallet"`
}
