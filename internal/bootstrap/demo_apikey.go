package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DemoTenantID   = "demo-playground-tenant"
	DemoTenantName = "Playground Demo"
)

// EnsureDemoAPIKey cria/atualiza a API Key especial para o playground público
// Agora usa configurações do admin/settings (Playground.APIKey e Playground.RateLimit)
func EnsureDemoAPIKey(ctx context.Context, apikeys *storage.APIKeysRepo, tenants *storage.TenantsRepo, settings *storage.SettingsRepo, db *mongo.Database) error {
	fmt.Println("🎯 Verificando API Key Demo para Playground...")

	// Buscar configurações do playground
	sysSettings, err := settings.Get(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar settings: %w", err)
	}

	if !sysSettings.Playground.Enabled {
		fmt.Println("⚠️ Playground desabilitado nas configurações, pulando...")
		return nil
	}

	demoAPIKey := sysSettings.Playground.APIKey
	if demoAPIKey == "" {
		demoAPIKey = "rtc_demo_playground_2024" // Fallback
	}

	// 1. Verificar se tenant demo existe
	tenantsCollection := db.Collection("tenants")
	var existingTenant domain.Tenant
	errTenant := tenantsCollection.FindOne(ctx, bson.M{"tenantId": DemoTenantID}).Decode(&existingTenant)
	
	if errTenant == mongo.ErrNoDocuments {
		// 2. Criar tenant demo (usando rate limit do settings)
		fmt.Println("📦 Criando tenant demo...")
		demoTenant := &domain.Tenant{
			TenantID:  DemoTenantID,
			Name:      DemoTenantName,
			Email:     "playground@theretech.com.br",
			Company:   "Playground Demo",
			Purpose:   "API Testing & Demo",
			RateLimit: &sysSettings.Playground.RateLimit, // Do settings!
			Active:    true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		
		_, errInsert := tenantsCollection.InsertOne(ctx, demoTenant)
		if errInsert != nil {
			return fmt.Errorf("erro ao criar tenant demo: %w", errInsert)
		}
		fmt.Println("✅ Tenant demo criado!")
	} else if errTenant != nil {
		return fmt.Errorf("erro ao verificar tenant demo: %w", errTenant)
	} else {
		fmt.Println("✅ Tenant demo já existe")
	}

	// 3. Verificar se API Key demo existe (usando keyHash)
	apikeysCollection := db.Collection("apikeys")
	var existingKey domain.APIKey
	err = apikeysCollection.FindOne(ctx, bson.M{"keyId": "rtc_demo_playground"}).Decode(&existingKey)
	
	if err == mongo.ErrNoDocuments {
		// Criar nova API Key demo (usando settings.Playground.APIKey)
		fmt.Printf("🔑 Criando API Key demo: %s\n", demoAPIKey)
		
		demoKey := &domain.APIKey{
			KeyID:     "rtc_demo_playground",
			KeyHash:   demoAPIKey, // Usar valor do settings
			Scopes:    sysSettings.Playground.AllowedAPIs, // APIs permitidas do settings
			OwnerID:   DemoTenantID,
			ExpiresAt: time.Now().UTC().AddDate(10, 0, 0), // Expira em 10 anos
			Revoked:   false,
			CreatedAt: time.Now().UTC(),
		}
		
		_, err = apikeysCollection.InsertOne(ctx, demoKey)
		if err != nil {
			return fmt.Errorf("erro ao criar API key demo: %w", err)
		}
		
		fmt.Printf("✅ API Key demo criada: %s\n", demoAPIKey)
		return nil
	}
	
	if err != nil {
		return fmt.Errorf("erro ao verificar API key demo: %w", err)
	}

	// API Key já existe - atualizar com valores do settings
	fmt.Println("✅ API Key demo já existe, atualizando com settings...")
	
	// Atualizar keyHash e scopes conforme settings
	update := bson.M{
		"$set": bson.M{
			"keyHash": demoAPIKey,                       // Atualizar chave do settings
			"revoked": false,
			"scopes":  sysSettings.Playground.AllowedAPIs, // Atualizar scopes do settings
		},
	}
	
	_, err = apikeysCollection.UpdateOne(ctx, bson.M{"keyId": "rtc_demo_playground"}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar API key demo: %w", err)
	}

	fmt.Printf("✅ API Key demo atualizada: %s (scopes: %v)\n", demoAPIKey, sysSettings.Playground.AllowedAPIs)
	return nil
}

