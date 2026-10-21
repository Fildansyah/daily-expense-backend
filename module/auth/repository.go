package auth

import (
	"context"
	"errors"
	"time"

	"expense.com/m/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID, familyID, tokenHash string, expiresAt time.Time) (*model.RefreshToken, error) {
	entity := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		FamilyID:  familyID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var entity model.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id string, replacedBy *string) error {
	updates := map[string]interface{}{"revoked_at": time.Now().UTC()}
	if replacedBy != nil {
		updates["replaced_by"] = *replacedBy
	}

	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Updates(updates).Error
}

func (r *RefreshTokenRepository) RevokeFamily(ctx context.Context, familyID string) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		Update("revoked_at", time.Now().UTC()).Error
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at < ?", userID, time.Now().UTC()).
		Delete(&model.RefreshToken{}).Error
}
