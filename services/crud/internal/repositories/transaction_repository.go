package repositories

import (
    "context"
    "errors"
    "sync"

    "github.com/afrikpay/gateway/services/crud/internal/models"
)

// TransactionRepository abstracts transaction storage.
type TransactionRepository interface {
    Create(ctx context.Context, t *models.Transaction) error
    GetByID(ctx context.Context, id string) (*models.Transaction, error)
    Update(ctx context.Context, t *models.Transaction) error
    Delete(ctx context.Context, id string) error
}

type InMemoryTransactionRepository struct {
    mu    sync.RWMutex
    store map[string]*models.Transaction
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
    return &InMemoryTransactionRepository{store: make(map[string]*models.Transaction)}
}

func (r *InMemoryTransactionRepository) Create(_ context.Context, t *models.Transaction) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[t.ID]; ok {
        return errors.New("transaction already exists")
    }
    cp := *t
    r.store[t.ID] = &cp
    return nil
}

func (r *InMemoryTransactionRepository) GetByID(_ context.Context, id string) (*models.Transaction, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    tr, ok := r.store[id]
    if !ok {
        return nil, ErrNotFound
    }
    cp := *tr
    return &cp, nil
}

func (r *InMemoryTransactionRepository) Update(_ context.Context, t *models.Transaction) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[t.ID]; !ok {
        return ErrNotFound
    }
    cp := *t
    r.store[t.ID] = &cp
    return nil
}

func (r *InMemoryTransactionRepository) Delete(_ context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[id]; !ok {
        return ErrNotFound
    }
    delete(r.store, id)
    return nil
}
