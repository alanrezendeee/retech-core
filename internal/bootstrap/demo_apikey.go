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
	DemoAPIKeyValue  = "rtc_demo_playground_2024"
	DemoTenantID     = "demo-playground-tenant"
	DemoTenantName   = "Playground Demo"
)

// EnsureDemoAPIKey cria/atualiza a API Key especial para o playground público
func EnsureDemoAPIKey(ctx context.Context, apikeys *storage.APIKeysRepo, tenants *storage.TenantsRepo, db *mongo.Database) error {
	fmt.Println("🎯 Verificando API Key Demo para Playground...")

	// 1. Verificar se tenant demo existe
	tenantsCollection := db.Collection("tenants")
	var existingTenant domain.Tenant
	err := tenantsCollection.FindOne(ctx, bson.M{"tenantId": DemoTenantID}).Decode(&existingTenant)
	
	if err == mongo.ErrNoDocuments {
		// 2. Criar tenant demo
		fmt.Println("📦 Criando tenant demo...")
		demoTenant := &domain.Tenant{
			TenantID: DemoTenantID,
			Name:     DemoTenantName,
			Email:    "playground@theretech.com.br",
			Company:  "Playground Demo",
			Purpose:  "API Testing & Demo",
			RateLimit: &domain.RateLimitConfig{
				RequestsPerDay:    100,  // Limite BEM BAIXO para playground
				RequestsPerMinute: 10,   // Rate limit agressivo
			},
			Active:    true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		
		_, err = tenantsCollection.InsertOne(ctx, demoTenant)
		if err != nil {
			return fmt.Errorf("erro ao criar tenant demo: %w", err)
		}
		fmt.Println("✅ Tenant demo criado!")
	} else if err != nil {
		return fmt.Errorf("erro ao verificar tenant demo: %w", err)
	} else {
		fmt.Println("✅ Tenant demo já existe")
	}

	// 3. Verificar se API Key demo existe (usando keyHash)
	apikeysCollection := db.Collection("apikeys")
	var existingKey domain.APIKey
	err = apikeysCollection.FindOne(ctx, bson.M{"keyId": "rtc_demo_playground"}).Decode(&existingKey)
	
	if err == mongo.ErrNoDocuments {
		// Criar nova API Key demo
		fmt.Println("🔑 Criando API Key demo...")
		
		demoKey := &domain.APIKey{
			KeyID:     "rtc_demo_playground",
			KeyHash:   DemoAPIKeyValue, // Usar valor completo como hash (simplificado para demo)
			Scopes:    []string{"cep", "cnpj", "geo"},
			OwnerID:   DemoTenantID,
			ExpiresAt: time.Now().UTC().AddDate(10, 0, 0), // Expira em 10 anos
			Revoked:   false,
			CreatedAt: time.Now().UTC(),
		}
		
		_, err = apikeysCollection.InsertOne(ctx, demoKey)
		if err != nil {
			return fmt.Errorf("erro ao criar API key demo: %w", err)
		}
		
		fmt.Println("✅ API Key demo criada: rtc_demo_playground_2024")
		return nil
	}
	
	if err != nil {
		return fmt.Errorf("erro ao verificar API key demo: %w", err)
	}

	// API Key já existe - atualizar se necessário
	fmt.Println("✅ API Key demo já existe")
	
	// Garantir que não está revogada e com scopes corretos
	update := bson.M{
		"$set": bson.M{
			"revoked": false,
			"scopes":  []string{"cep", "cnpj", "geo"},
		},
	}
	
	_, err = apikeysCollection.UpdateOne(ctx, bson.M{"keyId": "rtc_demo_playground"}, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar API key demo: %w", err)
	}

	fmt.Println("✅ API Key demo atualizada e ativa!")
	return nil
}

