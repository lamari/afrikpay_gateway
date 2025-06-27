package config

import (
    "context"
    "log"
    "os"
    "sync"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    mongoOnce sync.Once
    mongoCli  *mongo.Client
)

// MongoClient returns a singleton mongo client.
func MongoClient() *mongo.Client {
    mongoOnce.Do(func() {
        uri := os.Getenv("MONGO_URI")
        if uri == "" {
            uri = "mongodb://localhost:27017"
        }
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
        if err != nil {
            log.Fatalf("cannot connect mongo: %v", err)
        }
        mongoCli = client
    })
    return mongoCli
}
