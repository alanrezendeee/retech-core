package storage

import (
	"context"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SystemSettingsID = "system-settings-singleton" // ID único para singleton
)

type SettingsRepo struct {
	col *mongo.Collection
}

func NewSettingsRepo(db *mongo.Database) *SettingsRepo {
	return &SettingsRepo{col: db.Collection("system_settings")}
}

// Get retorna as configurações do sistema (ou padrões se não existir)
func (r *SettingsRepo) Get(ctx context.Context) (*domain.SystemSettings, error) {
	var settings domain.SystemSettings
	err := r.col.FindOne(ctx, bson.M{"_id": SystemSettingsID}).Decode(&settings)

	if err == mongo.ErrNoDocuments {
		// Retornar configurações padrão se não existir
		return domain.GetDefaultSettings(), nil
	}

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// Update atualiza as configurações do sistema
func (r *SettingsRepo) Update(ctx context.Context, settings *domain.SystemSettings) error {
	settings.UpdatedAt = time.Now().UTC()

	// Upsert: cria se não existir, atualiza se existir
	upsert := true
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": SystemSettingsID},
		bson.M{
			"$set": settings,
			"$setOnInsert": bson.M{
				"_id":       SystemSettingsID,
				"createdAt": time.Now().UTC(),
			},
		},
		&options.UpdateOptions{
			Upsert: &upsert,
		},
	)

	return err
}

// Ensure cria as configurações padrão se não existirem
func (r *SettingsRepo) Ensure(ctx context.Context) error {
	count, err := r.col.CountDocuments(ctx, bson.M{"_id": SystemSettingsID})
	if err != nil {
		return err
	}

	if count == 0 {
		defaultSettings := domain.GetDefaultSettings()
		defaultSettings.ID = SystemSettingsID
		_, err = r.col.InsertOne(ctx, defaultSettings)
		return err
	}

	return nil
}
