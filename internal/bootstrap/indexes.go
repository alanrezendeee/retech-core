package bootstrap

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// users: email único
	_, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil { return err }

	// refresh_tokens: TTL em expiresAt
	_, err = db.Collection("refresh_tokens").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil { return err }

	// api_keys: TTL
	_, err = db.Collection("api_keys").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	return err
}

