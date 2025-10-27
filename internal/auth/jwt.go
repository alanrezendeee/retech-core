package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/theretech/retech-core/internal/domain"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

// JWTClaims representa as claims customizadas do JWT
type JWTClaims struct {
	UserID   string          `json:"userId"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	TenantID string          `json:"tenantId,omitempty"`
	jwt.RegisteredClaims
}

// JWTService gerencia tokens JWT
type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewJWTService cria um novo serviço JWT
func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// GenerateAccessToken gera um access token
func (s *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		TenantID: user.TenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "retech-core",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

// GenerateRefreshToken gera um refresh token
func (s *JWTService) GenerateRefreshToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "retech-core",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

// ValidateAccessToken valida um access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, s.accessSecret)
}

// ValidateRefreshToken valida um refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString, s.refreshSecret)
}

// validateToken valida um token com o secret fornecido
func (s *JWTService) validateToken(tokenString string, secret []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessTTL retorna o TTL do access token
func (s *JWTService) GetAccessTTL() time.Duration {
	return s.accessTTL
}

// GetRefreshTTL retorna o TTL do refresh token
func (s *JWTService) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}

// SetAccessTTL atualiza o TTL do access token dinamicamente
func (s *JWTService) SetAccessTTL(ttl time.Duration) {
	s.accessTTL = ttl
}

// SetRefreshTTL atualiza o TTL do refresh token dinamicamente
func (s *JWTService) SetRefreshTTL(ttl time.Duration) {
	s.refreshTTL = ttl
}
