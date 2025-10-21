package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
)

// AuthJWT middleware para validar JWT e extrair claims
func AuthJWT(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": http.StatusUnauthorized,
				"detail": "Token de autenticação não fornecido",
			})
			c.Abort()
			return
		}

		// Formato esperado: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": http.StatusUnauthorized,
				"detail": "Formato de token inválido. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			detail := "Token inválido"

			if err == ErrExpiredToken {
				detail = "Token expirado"
			}

			c.JSON(status, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": status,
				"detail": detail,
			})
			c.Abort()
			return
		}

		// Armazenar claims no contexto
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("tenantID", claims.TenantID)

		c.Next()
	}
}

// RequireSuperAdmin middleware que requer role SUPER_ADMIN
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": http.StatusUnauthorized,
				"detail": "Autenticação necessária",
			})
			c.Abort()
			return
		}

		if role != domain.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"type":   "https://retech-core/errors/forbidden",
				"title":  "Forbidden",
				"status": http.StatusForbidden,
				"detail": "Acesso restrito a super administradores",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireTenantUser middleware que requer role TENANT_USER
func RequireTenantUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": http.StatusUnauthorized,
				"detail": "Autenticação necessária",
			})
			c.Abort()
			return
		}

		if role != domain.RoleTenantUser {
			c.JSON(http.StatusForbidden, gin.H{
				"type":   "https://retech-core/errors/forbidden",
				"title":  "Forbidden",
				"status": http.StatusForbidden,
				"detail": "Acesso restrito a usuários de tenant",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID helper para extrair userID do contexto
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return userID.(string)
	}
	return ""
}

// GetTenantID helper para extrair tenantID do contexto
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get("tenantID"); exists {
		return tenantID.(string)
	}
	return ""
}

// GetRole helper para extrair role do contexto
func GetRole(c *gin.Context) domain.UserRole {
	if role, exists := c.Get("role"); exists {
		return role.(domain.UserRole)
	}
	return ""
}

// IsSuperAdmin helper para verificar se é super admin
func IsSuperAdmin(c *gin.Context) bool {
	return GetRole(c) == domain.RoleSuperAdmin
}

