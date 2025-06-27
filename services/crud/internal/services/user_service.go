package services

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
)

type UserService struct {
    Repo repositories.UserRepository
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
    if err := user.Validate(); err != nil {
        return err
    }
    return s.Repo.Create(ctx, user)
}

func (s *UserService) Get(ctx context.Context, id string) (*models.User, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
    if err := user.Validate(); err != nil {
        return err
    }
    return s.Repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
    return s.Repo.Delete(ctx, id)
}
