package report

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/money"
)

const (
	dateLayout  = "2006-01-02"
	monthLayout = "2006-01"
)

type TransactionReader interface {
	Summary(ctx context.Context, userID, from, to string) (*entity.TransactionSummaryDTO, error)
	GetRange(ctx context.Context, userID, from, to string) ([]*entity.TransactionDTO, error)
	Recent(ctx context.Context, userID string, limit int) ([]*entity.TransactionDTO, error)
	SumByCategory(ctx context.Context, userID, from, to, transactionType string) ([]entity.CategorySumDTO, error)
	Series(ctx context.Context, userID, from, to, granularity string) ([]entity.TransactionSeriesDTO, error)
}

type BillReader interface {
	GetBillSummary(ctx context.Context, userID, from, to string) (*entity.BillSummaryDTO, error)
	GetRange(ctx context.Context, userID, from, to string) ([]*entity.BillDTO, error)
}

type WalletReader interface {
	FindWalletByUserID(ctx context.Context, userID string, includeArchived bool) ([]*entity.WalletDTO, error)
	GetSummary(ctx context.Context, userID string) (*entity.WalletSummaryDTO, error)
}

type CategoryReader interface {
	GetTree(ctx context.Context, userID, transactionType string, includeArchived bool) ([]*entity.CategoryDTO, error)
}

type BudgetReader interface {
	GetBudgets(ctx context.Context, userID, month string) ([]*entity.BudgetDTO, error)
}

type Service struct {
	transactions TransactionReader
	bills        BillReader
	wallets      WalletReader
	categories   CategoryReader
	budgets      BudgetReader
}

func NewService(transactions TransactionReader, bills BillReader, wallets WalletReader, categories CategoryReader, budgets BudgetReader) *Service {
	return &Service{
		transactions: transactions,
		bills:        bills,
		wallets:      wallets,
		categories:   categories,
		budgets:      budgets,
	}
}

func (s *Service) Dashboard(ctx context.Context, userID, from, to string) (*entity.DashboardDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}
	fromStr := fromTime.Format(dateLayout)
	toStr := toTime.Format(dateLayout)

	current, err := s.transactions.Summary(ctx, userID, fromStr, toStr)
	if err != nil {
		return nil, err
	}

	prevFrom, prevTo := previousRange(fromTime, toTime)
	previous, err := s.transactions.Summary(ctx, userID, prevFrom.Format(dateLayout), prevTo.Format(dateLayout))
	if err != nil {
		return nil, err
	}

	walletSummary, err := s.wallets.GetSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	billSummary, err := s.bills.GetBillSummary(ctx, userID, fromStr, toStr)
	if err != nil {
		return nil, err
	}

	cashflow, err := s.cashflowPoints(ctx, userID, addMonths(monthStart(toTime), -5), toTime, "month")
	if err != nil {
		return nil, err
	}

	breakdown, err := s.categoryReport(ctx, userID, fromStr, toStr, "EXPENSE")
	if err != nil {
		return nil, err
	}
	breakdown = collapseCategories(breakdown, 5)

	today := time.Now().UTC()
	todayStr := today.Format(dateLayout)
	weekAhead := today.AddDate(0, 0, 7).Format(dateLayout)

	upcoming, err := s.bills.GetRange(ctx, userID, todayStr, weekAhead)
	if err != nil {
		return nil, err
	}
	upcoming = unpaidBills(upcoming)
	if len(upcoming) > 5 {
		upcoming = upcoming[:5]
	}

	recent, err := s.transactions.Recent(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	budgetProgress, err := s.budgets.GetBudgets(ctx, userID, toTime.Format(monthLayout))
	if err != nil {
		return nil, err
	}
	if len(budgetProgress) > 5 {
		budgetProgress = budgetProgress[:5]
	}

	wallets, err := s.wallets.FindWalletByUserID(ctx, userID, false)
	if err != nil {
		return nil, err
	}

	return &entity.DashboardDTO{
		Period: entity.ReportPeriodDTO{From: fromStr, To: toStr},
		KPI: entity.ReportKPI{
			TotalBalance: walletSummary.TotalBalance,
			Income:       current.Income,
			Expense:      current.Expense,
			Net:          current.Net,
			Debt:         billSummary.RemainingAll,
			IncomeDelta:  percentChange(current.Income, previous.Income),
			ExpenseDelta: percentChange(current.Expense, previous.Expense),
			NetDelta:     percentChange(current.Net, previous.Net),
			DebtDelta:    0,
		},
		Cashflow:           cashflow,
		CategoryBreakdown:  breakdown,
		UpcomingBills:      upcoming,
		RecentTransactions: recent,
		BudgetProgress:     budgetProgress,
		CreditUtilization:  creditUtilization(wallets),
	}, nil
}

func (s *Service) Calendar(ctx context.Context, userID, from, to string) (*entity.CalendarDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	days := int(toTime.Sub(fromTime).Hours()/24) + 1
	if days > 62 {
		return nil, fmt.Errorf("%w: rentang kalender maksimal 62 hari", apperror.ErrValidation)
	}

	fromStr := fromTime.Format(dateLayout)
	toStr := toTime.Format(dateLayout)

	transactions, err := s.transactions.GetRange(ctx, userID, fromStr, toStr)
	if err != nil {
		return nil, err
	}
	bills, err := s.bills.GetRange(ctx, userID, fromStr, toStr)
	if err != nil {
		return nil, err
	}

	dayMap := make(map[string]*entity.CalendarDayDTO, days)
	for cursor := fromTime; !cursor.After(toTime); cursor = cursor.AddDate(0, 0, 1) {
		key := cursor.Format(dateLayout)
		dayMap[key] = &entity.CalendarDayDTO{
			Date:         key,
			Transactions: []*entity.TransactionDTO{},
			Bills:        []*entity.BillDTO{},
		}
	}

	summary := entity.ReportTotalsDTO{From: fromStr, To: toStr}
	billsDue := 0.0

	for _, transaction := range transactions {
		day, ok := dayMap[transaction.Date]
		if !ok {
			continue
		}
		switch transaction.Type {
		case "INCOME":
			day.Income = money.Round2(day.Income + transaction.Amount)
			summary.Income = money.Round2(summary.Income + transaction.Amount)
			summary.Count++
		case "EXPENSE":
			day.Expense = money.Round2(day.Expense + transaction.Amount)
			summary.Expense = money.Round2(summary.Expense + transaction.Amount)
			summary.Count++
		}
		day.TransactionCount++
		day.Transactions = append(day.Transactions, transaction)
	}

	for _, bill := range bills {
		day, ok := dayMap[bill.DueDate]
		if !ok {
			continue
		}
		day.Bills = append(day.Bills, bill)
		if bill.Status != "lunas" {
			billsDue = money.Round2(billsDue + bill.RemainingAmount)
		}
	}

	summary.Net = money.Round2(summary.Income - summary.Expense)

	keys := make([]string, 0, len(dayMap))
	for key := range dayMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	dayList := make([]entity.CalendarDayDTO, 0, len(keys))
	for _, key := range keys {
		day := dayMap[key]
		day.Net = money.Round2(day.Income - day.Expense)
		dayList = append(dayList, *day)
	}

	return &entity.CalendarDTO{
		From:     fromStr,
		To:       toStr,
		Days:     dayList,
		Summary:  summary,
		BillsDue: billsDue,
	}, nil
}

func (s *Service) Summary(ctx context.Context, userID, from, to string, compare bool) (*entity.ReportSummaryDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	current, err := s.transactions.Summary(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout))
	if err != nil {
		return nil, err
	}

	result := &entity.ReportSummaryDTO{
		Period:  entity.ReportPeriodDTO{From: fromTime.Format(dateLayout), To: toTime.Format(dateLayout)},
		Income:  current.Income,
		Expense: current.Expense,
		Net:     current.Net,
		Count:   current.Count,
		Previous: entity.ReportTotalsDTO{
			From: fromTime.Format(dateLayout),
			To:   toTime.Format(dateLayout),
		},
	}

	if compare {
		prevFrom, prevTo := previousRange(fromTime, toTime)
		previous, err := s.transactions.Summary(ctx, userID, prevFrom.Format(dateLayout), prevTo.Format(dateLayout))
		if err != nil {
			return nil, err
		}
		result.Previous = entity.ReportTotalsDTO{
			From:    prevFrom.Format(dateLayout),
			To:      prevTo.Format(dateLayout),
			Income:  previous.Income,
			Expense: previous.Expense,
			Net:     previous.Net,
			Count:   previous.Count,
		}
		result.IncomeDelta = percentChange(current.Income, previous.Income)
		result.ExpenseDelta = percentChange(current.Expense, previous.Expense)
		result.NetDelta = percentChange(current.Net, previous.Net)
	}

	return result, nil
}

func (s *Service) Cashflow(ctx context.Context, userID, from, to, granularity string) ([]entity.CashflowPointDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}
	return s.cashflowPoints(ctx, userID, fromTime, toTime, granularity)
}

func (s *Service) Categories(ctx context.Context, userID, from, to, transactionType string) ([]entity.ReportCategoryDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	transactionType = strings.ToUpper(strings.TrimSpace(transactionType))
	if transactionType == "" {
		transactionType = "EXPENSE"
	}
	if transactionType != "INCOME" && transactionType != "EXPENSE" {
		return nil, fmt.Errorf("%w: type harus INCOME atau EXPENSE", apperror.ErrValidation)
	}

	return s.categoryReport(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout), transactionType)
}

func (s *Service) Trend(ctx context.Context, userID, from, to string, compare bool) (*entity.TrendDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	current, err := s.transactions.Series(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout), "month")
	if err != nil {
		return nil, err
	}

	result := &entity.TrendDTO{Current: trendPoints(current)}

	if compare {
		prevFrom, prevTo := previousRange(fromTime, toTime)
		previous, err := s.transactions.Series(ctx, userID, prevFrom.Format(dateLayout), prevTo.Format(dateLayout), "month")
		if err != nil {
			return nil, err
		}
		result.Previous = trendPoints(previous)
	}

	return result, nil
}

func (s *Service) Debt(ctx context.Context, userID, from, to string) (*entity.DebtReportDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	summary, err := s.bills.GetBillSummary(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout))
	if err != nil {
		return nil, err
	}

	rows, err := s.bills.GetRange(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout))
	if err != nil {
		return nil, err
	}

	remainingByMonth := map[string]float64{}
	for _, bill := range rows {
		if len(bill.DueDate) < 7 {
			continue
		}
		month := bill.DueDate[:7]
		if bill.Status == "lunas" {
			continue
		}
		remainingByMonth[month] = money.Round2(remainingByMonth[month] + bill.RemainingAmount)
	}

	months := make([]string, 0, len(remainingByMonth))
	for month := range remainingByMonth {
		months = append(months, month)
	}
	sort.Strings(months)

	series := make([]entity.DebtPointDTO, 0, len(months))
	for _, month := range months {
		series = append(series, entity.DebtPointDTO{Period: month, Remaining: remainingByMonth[month]})
	}

	return &entity.DebtReportDTO{
		Outstanding: summary.RemainingAll,
		Overdue:     summary.OverdueTotal,
		DueThisWeek: summary.DueThisWeek,
		ByLender:    summary.PerLender,
		ByStatus:    summary.ByStatus,
		ByWallet:    summary.ByWallet,
		Series:      series,
	}, nil
}

func (s *Service) NetWorth(ctx context.Context, userID, from, to string) (*entity.NetWorthReportDTO, error) {
	fromTime, toTime, err := resolveRange(from, to)
	if err != nil {
		return nil, err
	}

	walletSummary, err := s.wallets.GetSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.transactions.Series(ctx, userID, fromTime.Format(dateLayout), toTime.Format(dateLayout), "month")
	if err != nil {
		return nil, err
	}

	netByMonth := map[string]float64{}
	for _, row := range rows {
		switch row.Type {
		case "INCOME":
			netByMonth[row.Period] = money.Round2(netByMonth[row.Period] + row.Total)
		case "EXPENSE":
			netByMonth[row.Period] = money.Round2(netByMonth[row.Period] - row.Total)
		}
	}

	months := make([]string, 0, len(netByMonth))
	for month := range netByMonth {
		months = append(months, month)
	}
	sort.Strings(months)

	series := make([]entity.NetWorthPointDTO, 0, len(months))
	for _, month := range months {
		series = append(series, entity.NetWorthPointDTO{Period: month, NetChange: netByMonth[month]})
	}

	return &entity.NetWorthReportDTO{
		Assets:      walletSummary.TotalBalance,
		Liabilities: walletSummary.CreditUsed,
		Net:         money.Round2(walletSummary.TotalBalance - walletSummary.CreditUsed),
		ByType:      walletSummary.ByType,
		Series:      series,
	}, nil
}

func (s *Service) cashflowPoints(ctx context.Context, userID string, from, to time.Time, granularity string) ([]entity.CashflowPointDTO, error) {
	if granularity == "" {
		granularity = "month"
	}

	rows, err := s.transactions.Series(ctx, userID, from.Format(dateLayout), to.Format(dateLayout), granularity)
	if err != nil {
		return nil, err
	}

	points := map[string]*entity.CashflowPointDTO{}
	order := make([]string, 0)

	ensure := func(period string) *entity.CashflowPointDTO {
		if point, ok := points[period]; ok {
			return point
		}
		point := &entity.CashflowPointDTO{Period: period}
		points[period] = point
		order = append(order, period)
		return point
	}

	for _, row := range rows {
		point := ensure(row.Period)
		switch row.Type {
		case "INCOME":
			point.Income = money.Round2(point.Income + row.Total)
		case "EXPENSE":
			point.Expense = money.Round2(point.Expense + row.Total)
		}
	}

	sort.Strings(order)

	result := make([]entity.CashflowPointDTO, 0, len(order))
	for _, period := range order {
		point := points[period]
		point.Net = money.Round2(point.Income - point.Expense)
		result = append(result, *point)
	}
	return result, nil
}

func (s *Service) categoryReport(ctx context.Context, userID, from, to, transactionType string) ([]entity.ReportCategoryDTO, error) {
	sums, err := s.transactions.SumByCategory(ctx, userID, from, to, transactionType)
	if err != nil {
		return nil, err
	}

	tree, err := s.categories.GetTree(ctx, userID, transactionType, false)
	if err != nil {
		return nil, err
	}

	index := make(map[string]categorySum, len(sums))
	for _, row := range sums {
		index[row.CategoryId] = categorySum{total: row.Total, count: row.Count}
	}

	nodes := buildCategoryReport(tree, index)

	total := 0.0
	for _, node := range nodes {
		total += node.Amount
	}
	applyPercentage(nodes, total)

	return nodes, nil
}

type categorySum struct {
	total float64
	count int64
}

func buildCategoryReport(categories []*entity.CategoryDTO, sums map[string]categorySum) []entity.ReportCategoryDTO {
	nodes := make([]entity.ReportCategoryDTO, 0, len(categories))

	for _, category := range categories {
		node := entity.ReportCategoryDTO{
			CategoryId: category.ID,
			Name:       category.Name,
			Icon:       category.Icon,
			BgColor:    category.BgColor,
			ParentId:   category.ParentId,
		}
		if sum, ok := sums[category.ID]; ok {
			node.Amount = money.Round2(sum.total)
			node.Count = sum.count
		}

		children := buildCategoryReport(category.Children, sums)
		for _, child := range children {
			node.Amount = money.Round2(node.Amount + child.Amount)
			node.Count += child.Count
		}
		node.Children = children

		if node.Amount > 0 || len(children) > 0 {
			nodes = append(nodes, node)
		}
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Amount > nodes[j].Amount
	})

	return nodes
}

func applyPercentage(nodes []entity.ReportCategoryDTO, total float64) {
	for i := range nodes {
		nodes[i].Percentage = percentOf(nodes[i].Amount, total)
		applyPercentage(nodes[i].Children, total)
	}
}

func collapseCategories(nodes []entity.ReportCategoryDTO, limit int) []entity.ReportCategoryDTO {
	if len(nodes) <= limit {
		return nodes
	}

	top := make([]entity.ReportCategoryDTO, limit)
	copy(top, nodes[:limit])

	other := entity.ReportCategoryDTO{CategoryId: "other", Name: "Lainnya"}
	for _, node := range nodes[limit:] {
		other.Amount = money.Round2(other.Amount + node.Amount)
		other.Count += node.Count
		other.Percentage = math.Round((other.Percentage+node.Percentage)*10) / 10
	}
	top = append(top, other)

	return top
}

func trendPoints(rows []entity.TransactionSeriesDTO) []entity.TrendPointDTO {
	values := map[string]float64{}
	order := make([]string, 0)

	for _, row := range rows {
		if row.Type != "EXPENSE" {
			continue
		}
		if _, ok := values[row.Period]; !ok {
			order = append(order, row.Period)
		}
		values[row.Period] = money.Round2(values[row.Period] + row.Total)
	}

	sort.Strings(order)

	points := make([]entity.TrendPointDTO, 0, len(order))
	for _, period := range order {
		points = append(points, entity.TrendPointDTO{Period: period, Value: values[period]})
	}
	return points
}

func creditUtilization(wallets []*entity.WalletDTO) []entity.WalletCreditDTO {
	result := make([]entity.WalletCreditDTO, 0)
	for _, wallet := range wallets {
		if !money.IsLiability(wallet.Type) {
			continue
		}
		limit := 0.0
		if wallet.CreditLimit != nil {
			limit = *wallet.CreditLimit
		}
		result = append(result, entity.WalletCreditDTO{
			WalletId:    wallet.ID,
			Name:        wallet.Name,
			Color:       wallet.Color,
			Limit:       limit,
			Used:        wallet.Used,
			Available:   wallet.Available,
			Utilization: wallet.Utilization,
		})
	}
	return result
}

func unpaidBills(bills []*entity.BillDTO) []*entity.BillDTO {
	result := make([]*entity.BillDTO, 0, len(bills))
	for _, bill := range bills {
		if bill.Status == "lunas" {
			continue
		}
		result = append(result, bill)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].DueDate < result[j].DueDate
	})
	return result
}

func resolveRange(from, to string) (time.Time, time.Time, error) {
	now := time.Now().UTC()

	fromTime, err := parseDate(from)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	toTime, err := parseDate(to)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if from == "" {
		base := now
		if to != "" {
			base = toTime
		}
		fromTime = monthStart(base)
	}
	if to == "" {
		toTime = monthStart(fromTime).AddDate(0, 1, -1)
	}

	if fromTime.After(toTime) {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: from tidak boleh setelah to", apperror.ErrValidation)
	}

	return fromTime, toTime, nil
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: tanggal harus berformat yyyy-mm-dd", apperror.ErrValidation)
	}
	return parsed, nil
}

func previousRange(from, to time.Time) (time.Time, time.Time) {
	days := int(to.Sub(from).Hours()/24) + 1
	prevTo := from.AddDate(0, 0, -1)
	prevFrom := prevTo.AddDate(0, 0, -(days - 1))
	return prevFrom, prevTo
}

func monthStart(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func addMonths(value time.Time, months int) time.Time {
	return value.AddDate(0, months, 0)
}

func percentChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return math.Round((current-previous)/math.Abs(previous)*1000) / 10
}

func percentOf(value, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return math.Round(value/total*1000) / 10
}
