package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlaygroundAPIKeyHandler struct {
	apikeys  *storage.APIKeysRepo
	settings *storage.SettingsRepo
	db       *mongo.Database
}

func NewPlaygroundAPIKeyHandler(apikeys *storage.APIKeysRepo, settings *storage.SettingsRepo, db *mongo.Database) *PlaygroundAPIKeyHandler {
	return &PlaygroundAPIKeyHandler{
		apikeys:  apikeys,
		settings: settings,
		db:       db,
	}
}

// GenerateAPIKey cria uma nova API Key demo para o playground
func (h *PlaygroundAPIKeyHandler) GenerateAPIKey(c *gin.Context) {
	ctx := context.Background()

	// 1. Buscar configurações atuais
	settings, err := h.settings.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao buscar configurações",
			"status": http.StatusInternalServerError,
			"detail": "Não foi possível carregar as configurações do playground",
		})
		return
	}

	// 2. Gerar nova API Key (formato: keyId.keySecret)
	now := time.Now()
	keyID := fmt.Sprintf("rtc_demo_playground_%s", now.Format("20060102150405"))
	keySecret := generateRandomString(32)
	newAPIKey := fmt.Sprintf("%s.%s", keyID, keySecret)

	// 3. Determinar scopes baseado nos allowedAPIs configurados
	scopes := settings.Playground.AllowedAPIs
	if len(scopes) == 0 {
		scopes = []string{"cep", "cnpj", "geo"} // Default
	}

	// 4. Gerar hash da API Key (HMAC-SHA256)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	if secret == "" {
		secret = "default-secret-key" // Fallback (não usar em produção)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(newAPIKey))
	keyHash := hex.EncodeToString(mac.Sum(nil))

	// 5. Criar registro da API Key no MongoDB
	apiKeyDoc := domain.APIKey{
		KeyID:     keyID,
		KeyHash:   keyHash,
		Scopes:    scopes,
		OwnerID:   "playground-public",
		ExpiresAt: now.AddDate(10, 0, 0), // Expira em 10 anos
		Revoked:   false,
		CreatedAt: now,
	}

	// 5. Salvar no MongoDB
	_, err = h.db.Collection("api_keys").InsertOne(ctx, apiKeyDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao criar API Key",
			"status": http.StatusInternalServerError,
			"detail": "Não foi possível salvar a API Key no banco de dados",
		})
		return
	}

	// 6. Atualizar settings com a nova API Key
	_, err = h.db.Collection("system_settings").UpdateOne(
		ctx,
		bson.M{"_id": "system-settings-singleton"},
		bson.M{
			"$set": bson.M{
				"playground.apiKey": newAPIKey,
				"updatedAt":         now,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao atualizar configurações",
			"status": http.StatusInternalServerError,
			"detail": "API Key criada mas não foi possível atualizar as configurações",
		})
		return
	}

	// 7. Retornar sucesso
	c.JSON(http.StatusOK, gin.H{
		"apiKey":  newAPIKey,
		"scopes":  scopes,
		"message": "API Key demo gerada com sucesso",
	})
}

// RotateAPIKey desativa a API Key atual e gera uma nova
func (h *PlaygroundAPIKeyHandler) RotateAPIKey(c *gin.Context) {
	ctx := context.Background()

	// 1. Buscar configurações atuais
	settings, err := h.settings.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao buscar configurações",
			"status": http.StatusInternalServerError,
			"detail": "Não foi possível carregar as configurações do playground",
		})
		return
	}

	oldAPIKey := settings.Playground.APIKey

	// 2. Desativar API Key antiga (se existir e não for a default inválida)
	if oldAPIKey != "" && oldAPIKey != "rtc_demo_playground_2024" {
		_, err = h.db.Collection("api_keys").UpdateOne(
			ctx,
			bson.M{"key": oldAPIKey},
			bson.M{
				"$set": bson.M{
					"active":    false,
					"updatedAt": time.Now(),
				},
			},
		)
		if err != nil {
			fmt.Printf("⚠️  [PLAYGROUND] Erro ao desativar API Key antiga: %v\n", err)
			// Não retorna erro, continua com a geração
		}
	}

	// 3. Gerar nova API Key (formato: keyId.keySecret)
	now := time.Now()
	keyID := fmt.Sprintf("rtc_demo_playground_%s", now.Format("20060102150405"))
	keySecret := generateRandomString(32)
	newAPIKey := fmt.Sprintf("%s.%s", keyID, keySecret)

	// 4. Determinar scopes baseado nos allowedAPIs configurados
	scopes := settings.Playground.AllowedAPIs
	if len(scopes) == 0 {
		scopes = []string{"cep", "cnpj", "geo"} // Default
	}

	// 5. Gerar hash da API Key (HMAC-SHA256)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	if secret == "" {
		secret = "default-secret-key" // Fallback
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(newAPIKey))
	keyHash := hex.EncodeToString(mac.Sum(nil))

	// 6. Criar registro da nova API Key no MongoDB
	apiKeyDoc := domain.APIKey{
		KeyID:     keyID,
		KeyHash:   keyHash,
		Scopes:    scopes,
		OwnerID:   "playground-public",
		ExpiresAt: now.AddDate(10, 0, 0), // Expira em 10 anos
		Revoked:   false,
		CreatedAt: now,
	}

	// 6. Salvar no MongoDB
	_, err = h.db.Collection("api_keys").InsertOne(ctx, apiKeyDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao criar nova API Key",
			"status": http.StatusInternalServerError,
			"detail": "Não foi possível salvar a nova API Key no banco de dados",
		})
		return
	}

	// 7. Atualizar settings com a nova API Key
	_, err = h.db.Collection("system_settings").UpdateOne(
		ctx,
		bson.M{"_id": "system-settings-singleton"},
		bson.M{
			"$set": bson.M{
				"playground.apiKey": newAPIKey,
				"updatedAt":         now,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao atualizar configurações",
			"status": http.StatusInternalServerError,
			"detail": "Nova API Key criada mas não foi possível atualizar as configurações",
		})
		return
	}

	// 8. Retornar sucesso
	c.JSON(http.StatusOK, gin.H{
		"oldKey":  oldAPIKey,
		"newKey":  newAPIKey,
		"scopes":  scopes,
		"message": "API Key rotacionada com sucesso",
	})
}

// generateRandomString gera uma string aleatória segura
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback para timestamp se houver erro
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
