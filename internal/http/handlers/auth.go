package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"github.com/theretech/retech-core/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	users        *storage.UsersRepo
	tenants      *storage.TenantsRepo
	apikeys      *storage.APIKeysRepo
	activityRepo *storage.ActivityLogsRepo
	settings     *storage.SettingsRepo
	jwt          *auth.JWTService
}

func NewAuthHandler(users *storage.UsersRepo, tenants *storage.TenantsRepo, apikeys *storage.APIKeysRepo, activityRepo *storage.ActivityLogsRepo, settings *storage.SettingsRepo, jwt *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		users:        users,
		tenants:      tenants,
		apikeys:      apikeys,
		activityRepo: activityRepo,
		settings:     settings,
		jwt:          jwt,
	}
}

// Login autentica um usuário e retorna JWT
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	// Buscar usuário
	user, err := h.users.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "Unauthorized",
				"status": http.StatusUnauthorized,
				"detail": "Email ou senha incorretos",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar usuário",
		})
		return
	}

	// Verificar se usuário está ativo
	if !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Usuário inativo",
		})
		return
	}

	// Verificar senha
	if !h.users.VerifyPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Email ou senha incorretos",
		})
		return
	}

	// ✅ Buscar TTLs dinâmicos do settings
	ctx := c.Request.Context()
	sysSettings, err := h.settings.Get(ctx)
	if err != nil {
		// Fallback para defaults se settings não encontrado
		sysSettings = domain.GetDefaultSettings()
	}
	
	// ✅ Atualizar TTLs do JWT Service dinamicamente
	h.jwt.SetAccessTTL(time.Duration(sysSettings.JWT.AccessTokenTTL) * time.Second)
	h.jwt.SetRefreshTTL(time.Duration(sysSettings.JWT.RefreshTokenTTL) * time.Second)
	
	fmt.Printf("🔑 [Login] JWT TTLs: Access=%ds, Refresh=%ds\n", 
		sysSettings.JWT.AccessTokenTTL, sysSettings.JWT.RefreshTokenTTL)
	
	// Gerar tokens com TTLs dinâmicos
	accessToken, err := h.jwt.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao gerar token",
		})
		return
	}

	refreshToken, err := h.jwt.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao gerar refresh token",
		})
		return
	}

	// Atualizar último login
	_ = h.users.UpdateLastLogin(c.Request.Context(), user.ID)

	// Log da atividade
	utils.LogActivity(
		c,
		h.activityRepo,
		domain.ActivityTypeUserLogin,
		domain.ActionLogin,
		domain.Actor{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Role:   string(user.Role),
		},
		domain.Resource{
			Type: domain.ResourceTypeUser,
			ID:   user.ID,
			Name: user.Name,
		},
		map[string]interface{}{
			"role":     user.Role,
			"tenantId": user.TenantID,
		},
	)

	// Remover senha da resposta
	user.Password = ""

	c.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.jwt.GetAccessTTL().Seconds()),
		User:         user,
	})
}

// Register cria um novo tenant com primeiro usuário (self-service)
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Verificar se email do usuário já existe
	existingUser, _ := h.users.FindByEmail(ctx, req.UserEmail)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"type":   "https://retech-core/errors/conflict",
			"title":  "Email já em uso",
			"status": http.StatusConflict,
			"detail": fmt.Sprintf("O email '%s' já está cadastrado. Por favor, use outro email ou faça login.", req.UserEmail),
		})
		return
	}

	// Verificar se email do tenant já existe
	existingTenant, _ := h.tenants.FindByEmail(ctx, req.TenantEmail)
	if existingTenant != nil {
		c.JSON(http.StatusConflict, gin.H{
			"type":   "https://retech-core/errors/conflict",
			"title":  "Email da empresa já em uso",
			"status": http.StatusConflict,
			"detail": fmt.Sprintf("O email da empresa '%s' já está cadastrado.", req.TenantEmail),
		})
		return
	}

	// Buscar rate limit padrão das configurações
	settings, err := h.settings.Get(ctx)
	var defaultRateLimit *domain.RateLimitConfig
	if err == nil && settings != nil {
		defaultRateLimit = &domain.RateLimitConfig{
			RequestsPerDay:    settings.DefaultRateLimit.RequestsPerDay,
			RequestsPerMinute: settings.DefaultRateLimit.RequestsPerMinute,
		}
		fmt.Printf("✅ [Register] Aplicando rate limit padrão: %d/dia, %d/min\n",
			defaultRateLimit.RequestsPerDay, defaultRateLimit.RequestsPerMinute)
	} else {
		// Fallback se não encontrar settings
		defaultRateLimit = &domain.RateLimitConfig{
			RequestsPerDay:    1000,
			RequestsPerMinute: 60,
		}
		fmt.Printf("⚠️ [Register] Usando fallback: 1000/dia, 60/min\n")
	}

	// Criar tenant com rate limit padrão
	tenant := &domain.Tenant{
		TenantID:  generateTenantID(req.TenantName),
		Name:      req.TenantName,
		Email:     req.TenantEmail,
		Active:    true,
		RateLimit: defaultRateLimit, // ✅ Sempre salvar rate limit!
	}

	if err := h.tenants.Insert(ctx, tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao criar tenant",
		})
		return
	}

	// Criar primeiro usuário
	user := &domain.User{
		Email:    req.UserEmail,
		Name:     req.UserName,
		Role:     domain.RoleTenantUser,
		TenantID: tenant.TenantID, // Usar TenantID, não ID
		Active:   true,
	}

	if err := h.users.Create(ctx, user, req.UserPassword); err != nil {
		// Rollback: deletar tenant se falhar (TODO: implementar Delete)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao criar usuário",
		})
		return
	}

	// Criar primeira API Key automaticamente
	// TODO: Implementar criação de API Key
	// Por enquanto, tenant pode criar depois no painel
	var generatedKey string

	// ✅ Buscar TTLs dinâmicos do settings
	sysSettings, err := h.settings.Get(ctx)
	if err != nil {
		sysSettings = domain.GetDefaultSettings()
	}
	
	// ✅ Atualizar TTLs do JWT Service
	h.jwt.SetAccessTTL(time.Duration(sysSettings.JWT.AccessTokenTTL) * time.Second)
	h.jwt.SetRefreshTTL(time.Duration(sysSettings.JWT.RefreshTokenTTL) * time.Second)
	
	// Gerar tokens JWT com TTLs dinâmicos
	accessToken, err := h.jwt.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao gerar token",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao gerar token de acesso",
		})
		return
	}

	refreshToken, err := h.jwt.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao gerar token",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao gerar refresh token",
		})
		return
	}

	// Log da atividade de criação de tenant
	utils.LogActivity(
		c,
		h.activityRepo,
		domain.ActivityTypeTenantCreated,
		domain.ActionCreate,
		domain.Actor{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Role:   string(user.Role),
		},
		domain.Resource{
			Type: domain.ResourceTypeTenant,
			ID:   tenant.TenantID,
			Name: tenant.Name,
		},
		map[string]interface{}{
			"email":   tenant.Email,
			"company": tenant.Company,
			"purpose": tenant.Purpose,
			"via":     "self-register",
		},
	)

	// Log da atividade de criação de usuário
	utils.LogActivity(
		c,
		h.activityRepo,
		domain.ActivityTypeUserCreated,
		domain.ActionCreate,
		domain.Actor{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Role:   string(user.Role),
		},
		domain.Resource{
			Type: domain.ResourceTypeUser,
			ID:   user.ID,
			Name: user.Name,
		},
		map[string]interface{}{
			"email":    user.Email,
			"role":     user.Role,
			"tenantId": user.TenantID,
			"via":      "self-register",
		},
	)

	// Remover senha
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"tenant":       tenant,
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"apiKey":       generatedKey, // Key só é mostrada uma vez
	})
}

// RefreshToken renova o access token usando refresh token
// POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	// Validar refresh token
	claims, err := h.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Refresh token inválido ou expirado",
		})
		return
	}

	// Buscar usuário
	user, err := h.users.FindByID(c.Request.Context(), claims.UserID)
	if err != nil || !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Usuário não encontrado ou inativo",
		})
		return
	}

	// ✅ Buscar TTLs dinâmicos do settings (para refresh também)
	ctx := c.Request.Context()
	sysSettings, err := h.settings.Get(ctx)
	if err != nil {
		sysSettings = domain.GetDefaultSettings()
	}
	
	// ✅ Atualizar TTLs
	h.jwt.SetAccessTTL(time.Duration(sysSettings.JWT.AccessTokenTTL) * time.Second)
	h.jwt.SetRefreshTTL(time.Duration(sysSettings.JWT.RefreshTokenTTL) * time.Second)

	// Gerar novo access token com TTL dinâmico
	newAccessToken, err := h.jwt.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao gerar novo token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": newAccessToken,
		"expiresIn":   int64(h.jwt.GetAccessTTL().Seconds()),
	})
}

// Me retorna dados do usuário logado
// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Autenticação necessária",
		})
		return
	}

	user, err := h.users.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://retech-core/errors/not-found",
			"title":  "Not Found",
			"status": http.StatusNotFound,
			"detail": "Usuário não encontrado",
		})
		return
	}

	// Remover senha
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// generateTenantID gera um ID único para o tenant baseado no nome
func generateTenantID(name string) string {
	// Simplificado - em produção use algo mais robusto
	return "tenant-" + time.Now().Format("20060102150405")
}
