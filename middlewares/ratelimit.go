package middlewares

import (
	"sync"
	"time"

	"expense.com/m/adapter/handlers/base_response"
	"github.com/gofiber/fiber/v2"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

// RateLimit adalah fixed-window limiter in-memory per IP.
// Cukup untuk satu instance API; tanpa Redis.
func RateLimit(max int, window time.Duration, message string) fiber.Handler {
	var mu sync.Mutex
	entries := make(map[string]*rateLimitEntry)

	return func(c *fiber.Ctx) error {
		key := c.IP()
		now := time.Now()

		mu.Lock()
		entry, ok := entries[key]
		if !ok || now.After(entry.resetAt) {
			entry = &rateLimitEntry{resetAt: now.Add(window)}
			entries[key] = entry
		}
		entry.count++

		allowed := entry.count <= max

		if len(entries) > 10000 {
			for k, v := range entries {
				if now.After(v.resetAt) {
					delete(entries, k)
				}
			}
		}
		mu.Unlock()

		if !allowed {
			return base_response.Error(c, fiber.StatusTooManyRequests, message)
		}

		return c.Next()
	}
}
