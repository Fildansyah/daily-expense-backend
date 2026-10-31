package tx

import (
	"context"

	"gorm.io/gorm"
)

type contextKey struct{}

// Manager menjalankan beberapa operasi repository dalam satu transaksi database
// tanpa membocorkan *gorm.DB ke lapisan service.
type Manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

// Run mengeksekusi fn di dalam transaksi. Pemanggilan bersarang memakai
// transaksi yang sama (savepoint tidak diperlukan untuk pola ini).
func (m *Manager) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	if existing, ok := ctx.Value(contextKey{}).(*gorm.DB); ok && existing != nil {
		return fn(ctx)
	}

	return m.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		return fn(context.WithValue(ctx, contextKey{}, transaction))
	})
}

// DB mengembalikan handle transaksi dari context bila ada, atau koneksi default.
func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if transaction, ok := ctx.Value(contextKey{}).(*gorm.DB); ok && transaction != nil {
		return transaction.WithContext(ctx)
	}
	return fallback.WithContext(ctx)
}
