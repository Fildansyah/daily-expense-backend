package lender

import (
	"context"
	"fmt"
	"strings"

	"expense.com/m/entity"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
)

type LenderService struct {
	repo *LenderRepository
	txm  *tx.Manager
}

func NewLenderService(repo *LenderRepository, txm *tx.Manager) *LenderService {
	return &LenderService{repo: repo, txm: txm}
}

func (s *LenderService) GetLenders(ctx context.Context, userID string) ([]*entity.LenderDTO, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *LenderService) CreateLender(ctx context.Context, userID, name string) (*entity.LenderDTO, error) {
	name = strings.TrimSpace(name)
	if err := validateLenderName(name); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsName(ctx, userID, name, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: lender dengan nama tersebut sudah ada", apperror.ErrConflict)
	}

	return s.repo.Create(ctx, &entity.LenderDTO{UserID: userID, Name: name})
}

func (s *LenderService) UpdateLender(ctx context.Context, userID, id, name string) (*entity.LenderDTO, error) {
	name = strings.TrimSpace(name)
	if err := validateLenderName(name); err != nil {
		return nil, err
	}

	var out *entity.LenderDTO
	err := s.txm.Run(ctx, func(ctx context.Context) error {
		existing, err := s.repo.FindOwned(ctx, id, userID)
		if err != nil {
			return err
		}

		if existing.Name == name {
			out = existing
			return nil
		}

		exists, err := s.repo.ExistsName(ctx, userID, name, id)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w: lender dengan nama tersebut sudah ada", apperror.ErrConflict)
		}

		// Relasi bill -> lender memakai lender_id, jadi rename cukup di tabel lenders.
		if err := s.repo.Update(ctx, &entity.LenderDTO{CommonDTOFields: entity.CommonDTOFields{ID: id}, UserID: userID, Name: name}); err != nil {
			return err
		}

		out, err = s.repo.FindOwned(ctx, id, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *LenderService) DeleteLender(ctx context.Context, userID, id string) error {
	existing, err := s.repo.FindOwned(ctx, id, userID)
	if err != nil {
		return err
	}

	count, err := s.repo.CountBills(ctx, userID, existing.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%w: lender masih dipakai oleh %d tagihan", apperror.ErrConflict, count)
	}

	return s.repo.Delete(ctx, id, userID)
}

func validateLenderName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: nama lender wajib diisi", apperror.ErrValidation)
	}
	if len([]rune(name)) > 255 {
		return fmt.Errorf("%w: nama lender maksimal 255 karakter", apperror.ErrValidation)
	}
	return nil
}
