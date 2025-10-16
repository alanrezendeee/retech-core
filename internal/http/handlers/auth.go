package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

type AuthHandler struct {
	Users  *storage.UsersRepo
	Tokens *storage.TokensRepo
	JWT    *auth.JWTService
}

func NewAuthHandler(u *storage.UsersRepo, t *storage.TokensRepo, j *auth.JWTService) *AuthHandler {
	return &AuthHandler{Users: u, Tokens: t, JWT: j}
}

type registerDTO struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginDTO struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	// opcional: habilitar apenas em dev
	if os.Getenv("ENV") == "production" {
		c.JSON(http.StatusForbidden, gin.H{"error":"disabled in production"}); return
	}
	var in registerDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	exists, _ := h.Users.ByEmail(c, in.Email)
	if exists != nil { c.JSON(http.StatusConflict, gin.H{"error":"email exists"}); return }
	hash, _ := auth.HashPassword(in.Password)
	u := &domain.User{
		ID: uuid.NewString(),
		Email: in.Email,
		Password: hash,
		Roles: []string{"user"},
	}
	if err := h.Users.Insert(c, u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"user insert"}); return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in loginDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	u, _ := h.Users.ByEmail(c, in.Email)
	if u == nil || !auth.CheckPassword(u.Password, in.Password) || !u.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return
	}
	at, atJTI, err := h.JWT.SignAccess(u.ID, u.Roles)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"sign access"}); return }
	rt, rtJTI, exp, err := h.JWT.SignRefresh(u.ID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"sign refresh"}); return }

	_ = h.Tokens.Insert(c, &domain.RefreshToken{
		ID: rtJTI, UserID: u.ID, ExpiresAt: exp, Revoked: false,
		CreatedAt: time.Now().UTC(), ParentJTI: "",
	})
	c.JSON(http.StatusOK, gin.H{
		"access_token": at,
		"access_jti":   atJTI,
		"refresh_token": rt,
		"refresh_jti":   rtJTI,
		"expires_in":    int(h.JWT.AccessTTL().Seconds()),
	})
}


type refreshDTO struct{ RefreshToken string `json:"refresh_token" binding:"required"` }

func (h *AuthHandler) Refresh(c *gin.Context) {
	var in refreshDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	_, claims, err := h.JWT.Parse(in.RefreshToken)
	if err != nil || claims["typ"] != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid refresh"}); return
	}
	jti := claims["jti"].(string)
	rt, _ := h.Tokens.FindActive(c, jti)
	if rt == nil || rt.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"revoked/expired refresh"}); return
	}
	// rotação segura: revoga o atual e emite novo
	_ = h.Tokens.Revoke(c, jti)

	userID := claims["sub"].(string)
	// (roles por user)
	u, _ := h.Users.ByID(c, userID)
	roles := []string{"user"}
	if u != nil { roles = u.Roles }

	at, atJTI, err := h.JWT.SignAccess(userID, roles)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"sign access"}); return }
	rtNew, rtJTI, exp, err := h.JWT.SignRefresh(userID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"sign refresh"}); return }

	_ = h.Tokens.Insert(c, &domain.RefreshToken{
		ID: rtJTI, UserID: userID, ExpiresAt: exp, Revoked: false,
		CreatedAt: time.Now().UTC(), ParentJTI: jti,
	})
	c.JSON(http.StatusOK, gin.H{
		"access_token": at,
		"access_jti":   atJTI,
		"refresh_token": rtNew,
		"refresh_jti":   rtJTI,
	})
}

type logoutDTO struct{ RefreshToken string `json:"refresh_token" binding:"required"` }

func (h *AuthHandler) Logout(c *gin.Context) {
	var in logoutDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	_, claims, err := h.JWT.Parse(in.RefreshToken)
	if err != nil || claims["typ"] != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid refresh"}); return
	}
	_ = h.Tokens.Revoke(c, claims["jti"].(string))
	c.Status(http.StatusNoContent)
}

