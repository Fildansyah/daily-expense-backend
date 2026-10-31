package entity

type CategorySumDTO struct {
	CategoryId string  `json:"categoryId"`
	Total      float64 `json:"total"`
	Count      int64   `json:"count"`
}

type TransactionSeriesDTO struct {
	Period string  `json:"period"`
	Type   string  `json:"type"`
	Total  float64 `json:"total"`
}

type ReportPeriodDTO struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ReportTotalsDTO struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
	Count   int64   `json:"count"`
}

type ReportKPI struct {
	TotalBalance float64 `json:"totalBalance"`
	Income       float64 `json:"income"`
	Expense      float64 `json:"expense"`
	Net          float64 `json:"net"`
	Debt         float64 `json:"debt"`
	IncomeDelta  float64 `json:"incomeDelta"`
	ExpenseDelta float64 `json:"expenseDelta"`
	NetDelta     float64 `json:"netDelta"`
	DebtDelta    float64 `json:"debtDelta"`
}

type CashflowPointDTO struct {
	Period  string  `json:"period"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type ReportCategoryDTO struct {
	CategoryId string              `json:"categoryId"`
	Name       string              `json:"name"`
	Icon       string              `json:"icon"`
	BgColor    string              `json:"bgColor"`
	ParentId   *string             `json:"parentId,omitempty"`
	Amount     float64             `json:"amount"`
	Percentage float64             `json:"percentage"`
	Count      int64               `json:"count"`
	Children   []ReportCategoryDTO `json:"children,omitempty"`
}

type TrendPointDTO struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
}

type TrendDTO struct {
	Current  []TrendPointDTO `json:"current"`
	Previous []TrendPointDTO `json:"previous,omitempty"`
}

type ReportSummaryDTO struct {
	Period       ReportPeriodDTO `json:"period"`
	Income       float64         `json:"income"`
	Expense      float64         `json:"expense"`
	Net          float64         `json:"net"`
	Count        int64           `json:"count"`
	Previous     ReportTotalsDTO `json:"previous"`
	IncomeDelta  float64         `json:"incomeDelta"`
	ExpenseDelta float64         `json:"expenseDelta"`
	NetDelta     float64         `json:"netDelta"`
}

type DebtPointDTO struct {
	Period    string  `json:"period"`
	Remaining float64 `json:"remaining"`
}

type DebtReportDTO struct {
	Outstanding float64                `json:"outstanding"`
	Overdue     float64                `json:"overdue"`
	DueThisWeek float64                `json:"dueThisWeek"`
	ByLender    []BillLenderSummaryDTO `json:"byLender"`
	ByStatus    []BillStatusSummaryDTO `json:"byStatus"`
	ByWallet    []BillWalletSummaryDTO `json:"byWallet"`
	Series      []DebtPointDTO         `json:"series"`
}

type NetWorthPointDTO struct {
	Period    string  `json:"period"`
	NetChange float64 `json:"netChange"`
}

type NetWorthReportDTO struct {
	Assets      float64                `json:"assets"`
	Liabilities float64                `json:"liabilities"`
	Net         float64                `json:"net"`
	ByType      []WalletTypeSummaryDTO `json:"byType"`
	Series      []NetWorthPointDTO     `json:"series"`
}

type WalletCreditDTO struct {
	WalletId    string  `json:"walletId"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Limit       float64 `json:"limit"`
	Used        float64 `json:"used"`
	Available   float64 `json:"available"`
	Utilization float64 `json:"utilization"`
}

type CalendarDayDTO struct {
	Date             string            `json:"date"`
	Income           float64           `json:"income"`
	Expense          float64           `json:"expense"`
	Net              float64           `json:"net"`
	TransactionCount int               `json:"transactionCount"`
	Transactions     []*TransactionDTO `json:"transactions"`
	Bills            []*BillDTO        `json:"bills"`
}

type CalendarDTO struct {
	From     string           `json:"from"`
	To       string           `json:"to"`
	Days     []CalendarDayDTO `json:"days"`
	Summary  ReportTotalsDTO  `json:"summary"`
	BillsDue float64          `json:"billsDue"`
}

type DashboardDTO struct {
	Period             ReportPeriodDTO     `json:"period"`
	KPI                ReportKPI           `json:"kpi"`
	Cashflow           []CashflowPointDTO  `json:"cashflow"`
	CategoryBreakdown  []ReportCategoryDTO `json:"categoryBreakdown"`
	UpcomingBills      []*BillDTO          `json:"upcomingBills"`
	RecentTransactions []*TransactionDTO   `json:"recentTransactions"`
	BudgetProgress     []*BudgetDTO        `json:"budgetProgress"`
	CreditUtilization  []WalletCreditDTO   `json:"creditUtilization"`
}
