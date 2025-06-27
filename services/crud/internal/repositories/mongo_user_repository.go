package repositories

import (
    "context"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
    col *mongo.Collection
}

func NewMongoUserRepository(col *mongo.Collection) *MongoUserRepository {
    return &MongoUserRepository{col: col}
}

func (r *MongoUserRepository) Create(ctx context.Context, u *models.User) error {
    _, err := r.col.InsertOne(ctx, u)
    return err
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    var result models.User
    err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
    if err == mongo.ErrNoDocuments {
        return nil, ErrNotFound
    }
    return &result, err
}

func (r *MongoUserRepository) Update(ctx context.Context, u *models.User) error {
    _, err := r.col.ReplaceOne(ctx, bson.M{"_id": u.ID}, u)
    return err
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
    res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return ErrNotFound
    }
    return nil
}
