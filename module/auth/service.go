package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"expense.com/m/entity"
	"expense.com/m/module/user"
	"expense.com/m/shared/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Grace period untuk menoleransi dua tab yang me-refresh hampir bersamaan.
const refreshGracePeriod = 30 * time.Second

// Hash dummy agar waktu respons saat email tidak ditemukan mirip dengan
// saat email ditemukan (mengurangi user enumeration lewat timing).
var dummyPasswordHash = func() []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte("timing-protection-only"), bcrypt.DefaultCost)
	return hash
}()

type AuthResult struct {
	AccessToken  string
	ExpiresIn    int64
	RefreshToken string
	User         *entity.UserDTO
}

type Service struct {
	users      *user.UserRepository
	refresh    *RefreshTokenRepository
	tokens     *token.Manager
	refreshTTL time.Duration
}

func NewService(users *user.UserRepository, refresh *RefreshTokenRepository, tokens *token.Manager, refreshTTL time.Duration) *Service {
	return &Service{
		users:      users,
		refresh:    refresh,
		tokens:     tokens,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	account, err := s.users.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return nil, ErrInvalidCredentials
	}

	if account.Password == "" || bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	_ = s.refresh.DeleteExpired(ctx, account.ID)

	return s.issue(ctx, account, uuid.NewString())
}

func (s *Service) Refresh(ctx context.Context, rawToken string) (*AuthResult, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	stored, err := s.refresh.FindByHash(ctx, hashToken(rawToken))
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrInvalidRefreshToken
	}

	now := time.Now().UTC()
	if !stored.ExpiresAt.After(now) {
		return nil, ErrInvalidRefreshToken
	}

	account, err := s.users.FindUserByID(stored.UserID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		_ = s.refresh.RevokeFamily(ctx, stored.FamilyID)
		return nil, ErrInvalidRefreshToken
	}

	accessToken, expiresIn, err := s.tokens.GenerateAccess(account.ID, account.Email)
	if err != nil {
		return nil, err
	}

	// Token sudah dirotasi: beri grace period tanpa memutar cookie lagi.
	// Grace hanya berlaku untuk token hasil rotasi (replaced_by terisi);
	// token yang dicabut karena logout harus langsung ditolak.
	if stored.RevokedAt != nil {
		if stored.ReplacedBy == nil || now.Sub(*stored.RevokedAt) > refreshGracePeriod {
			return nil, ErrInvalidRefreshToken
		}
		return &AuthResult{
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
			User:        account,
		}, nil
	}

	newRaw, newHash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	newStored, err := s.refresh.Create(ctx, account.ID, stored.FamilyID, newHash, now.Add(s.refreshTTL))
	if err != nil {
		return nil, err
	}

	if err := s.refresh.Revoke(ctx, stored.ID, &newStored.ID); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  accessToken,
		ExpiresIn:    expiresIn,
		RefreshToken: newRaw,
		User:         account,
	}, nil
}

func (s *Service) Logout(ctx context.Context, rawToken string) error {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil
	}

	stored, err := s.refresh.FindByHash(ctx, hashToken(rawToken))
	if err != nil {
		return err
	}
	if stored == nil || stored.RevokedAt != nil {
		return nil
	}

	return s.refresh.Revoke(ctx, stored.ID, nil)
}

func (s *Service) issue(ctx context.Context, account *entity.UserDTO, familyID string) (*AuthResult, error) {
	accessToken, expiresIn, err := s.tokens.GenerateAccess(account.ID, account.Email)
	if err != nil {
		return nil, err
	}

	rawToken, tokenHash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	if _, err := s.refresh.Create(ctx, account.ID, familyID, tokenHash, time.Now().UTC().Add(s.refreshTTL)); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  accessToken,
		ExpiresIn:    expiresIn,
		RefreshToken: rawToken,
		User:         account,
	}, nil
}

func newRefreshToken() (raw string, tokenHash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(buf)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
