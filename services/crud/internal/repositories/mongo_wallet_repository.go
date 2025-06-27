package repositories

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoWalletRepository struct {
    col *mongo.Collection
}

func NewMongoWalletRepository(col *mongo.Collection) *MongoWalletRepository {
    return &MongoWalletRepository{col: col}
}

func (r *MongoWalletRepository) Create(ctx context.Context, w *models.Wallet) error {
    _, err := r.col.InsertOne(ctx, w)
    return err
}

func (r *MongoWalletRepository) GetByID(ctx context.Context, id string) (*models.Wallet, error) {
    var result models.Wallet
    err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
    if err == mongo.ErrNoDocuments {
        return nil, ErrNotFound
    }
    return &result, err
}

func (r *MongoWalletRepository) Update(ctx context.Context, w *models.Wallet) error {
    _, err := r.col.ReplaceOne(ctx, bson.M{"_id": w.ID}, w)
    return err
}

func (r *MongoWalletRepository) Delete(ctx context.Context, id string) error {
    res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return ErrNotFound
    }
    return nil
}
