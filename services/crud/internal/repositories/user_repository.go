package repositories

import (
    "context"
    "errors"
    "sync"

    "github.com/afrikpay/gateway/services/crud/internal/models"
)

// ErrNotFound is returned when an entity is not found.
var ErrNotFound = errors.New("entity not found")

// UserRepository abstracts access to user storage.
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id string) error
}

// InMemoryUserRepository is an in-memory implementation useful for unit tests.
type InMemoryUserRepository struct {
    mu    sync.RWMutex
    store map[string]*models.User
}

// NewInMemoryUserRepository creates a new in-memory repo.
func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{store: make(map[string]*models.User)}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.store[user.ID]; exists {
        return errors.New("user already exists")
    }
    // copy to avoid external mutation
    cp := *user
    r.store[user.ID] = &cp
    return nil
}

func (r *InMemoryUserRepository) GetByID(_ context.Context, id string) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    u, ok := r.store[id]
    if !ok {
        return nil, ErrNotFound
    }
    cp := *u
    return &cp, nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[user.ID]; !ok {
        return ErrNotFound
    }
    cp := *user
    r.store[user.ID] = &cp
    return nil
}

func (r *InMemoryUserRepository) Delete(_ context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[id]; !ok {
        return ErrNotFound
    }
    delete(r.store, id)
    return nil
}
