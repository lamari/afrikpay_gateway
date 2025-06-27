package services

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
)

type TransactionService struct {
    Repo repositories.TransactionRepository
}

func (s *TransactionService) Create(ctx context.Context, t *models.Transaction) error {
    if err := t.Validate(); err != nil {
        return err
    }
    return s.Repo.Create(ctx, t)
}

func (s *TransactionService) Get(ctx context.Context, id string) (*models.Transaction, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *TransactionService) Update(ctx context.Context, t *models.Transaction) error {
    if err := t.Validate(); err != nil {
        return err
    }
    return s.Repo.Update(ctx, t)
}

func (s *TransactionService) Delete(ctx context.Context, id string) error {
    return s.Repo.Delete(ctx, id)
}
