package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AuthConfig struct {
	JWTSecret      string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	CookieSecure   bool
	CookieSameSite string
}

func LoadAuthConfig() (*AuthConfig, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi minimal 32 karakter (generate: openssl rand -hex 32)")
	}

	accessTTL, err := parseDurationEnv("JWT_ACCESS_TTL", "15m")
	if err != nil {
		return nil, err
	}

	refreshTTL, err := parseDurationEnv("JWT_REFRESH_TTL", "168h")
	if err != nil {
		return nil, err
	}

	if refreshTTL <= accessTTL {
		return nil, fmt.Errorf("JWT_REFRESH_TTL harus lebih besar dari JWT_ACCESS_TTL")
	}

	secure, err := strconv.ParseBool(envOr("COOKIE_SECURE", "false"))
	if err != nil {
		return nil, fmt.Errorf("COOKIE_SECURE harus bernilai true atau false: %w", err)
	}

	sameSite, err := parseSameSite(envOr("COOKIE_SAMESITE", "Lax"))
	if err != nil {
		return nil, err
	}

	return &AuthConfig{
		JWTSecret:      secret,
		AccessTTL:      accessTTL,
		RefreshTTL:     refreshTTL,
		CookieSecure:   secure,
		CookieSameSite: sameSite,
	}, nil
}

func parseDurationEnv(key, fallback string) (time.Duration, error) {
	value := envOr(key, fallback)
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s tidak valid (%q): %w", key, value, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%s harus lebih besar dari 0", key)
	}
	return d, nil
}

func parseSameSite(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "lax":
		return "Lax", nil
	case "strict":
		return "Strict", nil
	case "none":
		return "None", nil
	default:
		return "", fmt.Errorf("COOKIE_SAMESITE harus salah satu dari Lax, Strict, None")
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
