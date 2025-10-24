package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	now := time.Now().UTC()
	
	fmt.Printf("📝 Atualizando settings: %+v\n", settings)
	
	// Verificar se documento já existe
	var existing domain.SystemSettings
	err := r.col.FindOne(ctx, bson.M{"_id": SystemSettingsID}).Decode(&existing)
	
	if err == mongo.ErrNoDocuments {
		// Não existe - criar novo
		fmt.Println("ℹ️ Documento não existe, criando novo...")
		newSettings := &domain.SystemSettings{
			ID:               SystemSettingsID,
			DefaultRateLimit: settings.DefaultRateLimit,
			CORS:             settings.CORS,
			JWT:              settings.JWT,
			API:              settings.API,
			Contact:          settings.Contact,   // ✅ ADICIONADO
			Cache:            settings.Cache,     // ✅ ADICIONADO
			Playground:       settings.Playground, // ✅ ADICIONADO
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		_, err = r.col.InsertOne(ctx, newSettings)
		if err != nil {
			fmt.Printf("❌ Erro ao inserir: %v\n", err)
		} else {
			fmt.Println("✅ Documento criado com sucesso!")
		}
		return err
	}
	
	if err != nil {
		fmt.Printf("❌ Erro ao buscar documento existente: %v\n", err)
		return err
	}
	
	// Existe - atualizar apenas os campos necessários
	fmt.Println("ℹ️ Documento existe, atualizando...")
	_, err = r.col.UpdateOne(
		ctx,
		bson.M{"_id": SystemSettingsID},
		bson.M{
			"$set": bson.M{
				"defaultRateLimit": settings.DefaultRateLimit,
				"cors":             settings.CORS,
				"jwt":              settings.JWT,
				"api":              settings.API,
				"contact":          settings.Contact,   // ✅ ADICIONADO
				"cache":            settings.Cache,     // ✅ ADICIONADO
				"playground":       settings.Playground, // ✅ ADICIONADO
				"updatedAt":        now,
			},
		},
	)
	
	if err != nil {
		fmt.Printf("❌ Erro ao atualizar: %v\n", err)
	} else {
		fmt.Println("✅ Documento atualizado com sucesso!")
	}

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
