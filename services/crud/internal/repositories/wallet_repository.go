package repositories

import (
    "context"
    "errors"
    "sync"

    "github.com/afrikpay/gateway/services/crud/internal/models"
)

// WalletRepository abstracts wallet storage operations.
type WalletRepository interface {
    Create(ctx context.Context, w *models.Wallet) error
    GetByID(ctx context.Context, id string) (*models.Wallet, error)
    GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Wallet, error)
    Update(ctx context.Context, w *models.Wallet) error
    Delete(ctx context.Context, id string) error
}

type InMemoryWalletRepository struct {
    mu    sync.RWMutex
    store map[string]*models.Wallet
}

func NewInMemoryWalletRepository() *InMemoryWalletRepository {
    return &InMemoryWalletRepository{store: make(map[string]*models.Wallet)}
}

func (r *InMemoryWalletRepository) Create(_ context.Context, w *models.Wallet) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[w.ID]; ok {
        return errors.New("wallet already exists")
    }
    cp := *w
    r.store[w.ID] = &cp
    return nil
}

func (r *InMemoryWalletRepository) GetByID(_ context.Context, id string) (*models.Wallet, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    w, ok := r.store[id]
    if !ok {
        return nil, ErrNotFound
    }
    cp := *w
    return &cp, nil
}

func (r *InMemoryWalletRepository) Update(_ context.Context, w *models.Wallet) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[w.ID]; !ok {
        return ErrNotFound
    }
    cp := *w
    r.store[w.ID] = &cp
    return nil
}

func (r *InMemoryWalletRepository) Delete(_ context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.store[id]; !ok {
        return ErrNotFound
    }
    delete(r.store, id)
    return nil
}

func (r *InMemoryWalletRepository) GetByUserIDAndCurrency(_ context.Context, userID string, currency string) (*models.Wallet, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    for _, wallet := range r.store {
        if wallet.UserID == userID && wallet.Currency == currency {
            cp := *wallet
            return &cp, nil
        }
    }
    
    return nil, ErrNotFound
}
