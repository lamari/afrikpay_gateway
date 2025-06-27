package services

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
)

type WalletService struct {
    Repo repositories.WalletRepository
}

func (s *WalletService) Create(ctx context.Context, w *models.Wallet) error {
    if err := w.Validate(); err != nil {
        return err
    }
    return s.Repo.Create(ctx, w)
}

func (s *WalletService) Get(ctx context.Context, id string) (*models.Wallet, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *WalletService) Update(ctx context.Context, w *models.Wallet) error {
    if err := w.Validate(); err != nil {
        return err
    }
    return s.Repo.Update(ctx, w)
}

func (s *WalletService) Delete(ctx context.Context, id string) error {
    return s.Repo.Delete(ctx, id)
}
