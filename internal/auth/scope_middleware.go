package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

// RequireScope verifica se a API key tem o scope necessário para acessar a rota
func RequireScope(repo *storage.APIKeysRepo, requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pegar API Key do header
		raw := c.GetHeader("X-API-Key")
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "API Key Required",
				"status": http.StatusUnauthorized,
				"detail": "Header X-API-Key é obrigatório",
			})
			return
		}

		parts := strings.Split(raw, ".")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Invalid API Key Format",
				"status": http.StatusUnauthorized,
				"detail": "API Key deve estar no formato: keyId.keySecret",
			})
			return
		}

		keyId := parts[0]
		k, err := repo.ByKeyID(c, keyId)
		if err != nil || k == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unknown API Key",
				"status": http.StatusUnauthorized,
				"detail": "API Key não encontrada",
			})
			return
		}

		// Verificar se tem scope 'all' (acesso total)
		if hasScope(k.Scopes, "all") {
			c.Next()
			return
		}

		// Verificar se tem o scope específico
		if !hasScope(k.Scopes, requiredScope) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"type":   "https://retech-core/errors/forbidden",
				"title":  "Insufficient Permissions",
				"status": http.StatusForbidden,
				"detail": "API Key não tem permissão para acessar este recurso. Scope necessário: " + requiredScope,
				"meta": gin.H{
					"requiredScope": requiredScope,
					"yourScopes":    k.Scopes,
				},
			})
			return
		}

		c.Next()
	}
}

// hasScope verifica se um slice de scopes contém um scope específico
func hasScope(scopes []string, scope string) bool {
	// Extrair apenas o nome do scope (ex: "geo:read" → "geo")
	scopeName := scope
	if idx := strings.Index(scope, ":"); idx != -1 {
		scopeName = scope[:idx]
	}

	for _, s := range scopes {
		// Comparar sem sufixo :read/:write
		sName := s
		if idx := strings.Index(s, ":"); idx != -1 {
			sName = s[:idx]
		}
		
		if sName == scopeName || s == scope || s == "all" {
			return true
		}
	}
	return false
}

// ValidateAPIKeyScopes valida scopes ao criar API Key
func ValidateAPIKeyScopes(scopes []string) error {
	validScopes := map[string]bool{
		"geo":   true,
		"cep":   true,
		"cnpj":  true,
		"all":   true,
		// Futuros:
		"cpf":    false, // ainda não implementado
		"fipe":   false,
		"moedas": false,
		"bancos": false,
	}

	for _, scope := range scopes {
		// Remover sufixo :read/:write se houver
		scopeName := scope
		if idx := strings.Index(scope, ":"); idx != -1 {
			scopeName = scope[:idx]
		}

		if implemented, exists := validScopes[scopeName]; !exists {
			return &domain.ValidationError{
				Field:   "scopes",
				Message: "Scope '" + scope + "' não reconhecido",
			}
		} else if !implemented {
			return &domain.ValidationError{
				Field:   "scopes",
				Message: "Scope '" + scope + "' ainda não está disponível",
			}
		}
	}

	return nil
}

