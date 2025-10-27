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
	if err != nil { return err }

	// playground_rate_limits: índice composto para segurança
	_, err = db.Collection("playground_rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "ipAddress", Value: 1},
			{Key: "apiKey", Value: 1},
			{Key: "date", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil { return err }

	// playground_rate_limits: índice por IP para consultas rápidas
	_, err = db.Collection("playground_rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "ipAddress", Value: 1},
			{Key: "date", Value: 1},
		},
	})
	if err != nil { return err }

	// playground_rate_limits: TTL para limpeza automática
	_, err = db.Collection("playground_rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "updatedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(7 * 24 * 60 * 60), // 7 dias
	})
	if err != nil { return err }

	// playground_global_limits: índice composto
	_, err = db.Collection("playground_global_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "apiKey", Value: 1},
			{Key: "date", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil { return err }

	// playground_global_limits: TTL para limpeza automática
	_, err = db.Collection("playground_global_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "updatedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(7 * 24 * 60 * 60), // 7 dias
	})
	if err != nil { return err }

	return nil
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

	// ✅ PERFORMANCE: Índice único para CEP cache (hot path)
	_, err = db.Collection("cep_cache").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cep", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// ✅ PERFORMANCE: Índice único para CNPJ cache (hot path)
	_, err = db.Collection("cnpj_cache").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cnpj", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// ✅ PERFORMANCE: Índice para tenant_id (hot path - rate limiting)
	_, err = db.Collection("rate_limits").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenantId", Value: 1}, {Key: "resetAt", Value: 1}},
	})
	if err != nil {
		return err
	}

	// ✅ PERFORMANCE: Índice composto para rate_limits_minute
	_, err = db.Collection("rate_limits_minute").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenantId", Value: 1}, {Key: "resetAt", Value: 1}},
	})
	if err != nil {
		return err
	}

	log.Info().Msg("Índices criados com sucesso (incluindo cache performance)")
	return nil
}

// MigrateSettings adiciona campos faltantes nas configurações
func MigrateSettings(ctx context.Context, db *mongo.Database, log zerolog.Logger) error {
	log.Info().Msg("Migrando configurações...")

	// Verificar se settings tem o campo contact
	var settings bson.M
	err := db.Collection("system_settings").FindOne(ctx, bson.M{
		"_id": "system-settings-singleton",
	}).Decode(&settings)

	if err == nil {
		// Settings existe, verificar se tem contact
		if _, hasContact := settings["contact"]; !hasContact {
			log.Info().Msg("Adicionando campo contact nas configurações...")
			
			_, err = db.Collection("system_settings").UpdateOne(
				ctx,
				bson.M{"_id": "system-settings-singleton"},
				bson.M{
					"$set": bson.M{
						"contact": bson.M{
							"whatsapp": "48999616679",
							"email":    "suporte@theretech.com.br",
							"phone":    "+55 48 99961-6679",
						},
					},
				},
			)
			
			if err != nil {
				log.Error().Err(err).Msg("Erro ao migrar campo contact")
				return err
			}
			
			log.Info().Msg("Campo contact adicionado com sucesso!")
		}
		
		// Verificar se settings tem o campo cache
		if _, hasCache := settings["cache"]; !hasCache {
			log.Info().Msg("Adicionando campo cache nas configurações...")
			
			_, err = db.Collection("system_settings").UpdateOne(
				ctx,
				bson.M{"_id": "system-settings-singleton"},
				bson.M{
					"$set": bson.M{
						"cache": bson.M{
							"enabled":      true,
							"cepTtlDays":   7,
							"cnpjTtlDays":  30, // ✅ CNPJ: 30 dias
							"maxSizeMb":    100,
							"autoCleanup":  true,
						},
					},
				},
			)
			
			if err != nil {
				log.Error().Err(err).Msg("Erro ao migrar campo cache")
				return err
			}
			
			log.Info().Msg("Campo cache adicionado com sucesso!")
		}
	}

	log.Info().Msg("Migração de configurações concluída")
	return nil
}

