package repositories

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoTransactionRepository struct {
    col *mongo.Collection
}

func NewMongoTransactionRepository(col *mongo.Collection) *MongoTransactionRepository {
    return &MongoTransactionRepository{col: col}
}

func (r *MongoTransactionRepository) Create(ctx context.Context, t *models.Transaction) error {
    _, err := r.col.InsertOne(ctx, t)
    return err
}

func (r *MongoTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
    var result models.Transaction
    err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
    if err == mongo.ErrNoDocuments {
        return nil, ErrNotFound
    }
    return &result, err
}

func (r *MongoTransactionRepository) Update(ctx context.Context, t *models.Transaction) error {
    _, err := r.col.ReplaceOne(ctx, bson.M{"_id": t.ID}, t)
    return err
}

func (r *MongoTransactionRepository) Delete(ctx context.Context, id string) error {
    res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return ErrNotFound
    }
    return nil
}
