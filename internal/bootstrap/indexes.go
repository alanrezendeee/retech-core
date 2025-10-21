package bootstrap

import (
	"context"
	"time"

	"github.com/rs/zerolog"
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
	if err != nil { return err }

	// rate_limits: índice por API key
	_, err = db.Collection("rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "apiKey", Value: 1}},
	})
	if err != nil { return err }

	// rate_limits: TTL em resetAt
	_, err = db.Collection("rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "resetAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	return err
}

// CreateIndexes cria todos os índices necessários
func CreateIndexes(ctx context.Context, db *mongo.Database, log zerolog.Logger) error {
	log.Info().Msg("Criando índices...")

	// Estados: índice único por ID e por sigla
	_, err := db.Collection("estados").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("estados").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "sigla", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Municípios: índice único por ID
	_, err = db.Collection("municipios").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Municípios: índice para busca por UF
	_, err = db.Collection("municipios").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "microrregiao.mesorregiao.UF.sigla", Value: 1}},
	})
	if err != nil {
		return err
	}

	// Municípios: índice para busca por nome
	_, err = db.Collection("municipios").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "nome", Value: 1}},
	})
	if err != nil {
		return err
	}

	log.Info().Msg("Índices criados com sucesso")
	return nil
}

