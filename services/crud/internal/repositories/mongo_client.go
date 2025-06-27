package repositories

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MustNewMongoClient creates a client or panics (handy for tests).
func MustNewMongoClient(uri string) *mongo.Client {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Panicf("cannot connect mongo: %v", err)
    }
    return cli
}
