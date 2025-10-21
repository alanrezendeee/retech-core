package domain

import "time"

// SystemSettings representa as configurações globais do sistema
type SystemSettings struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`

	// Rate Limiting DEFAULT (para novos tenants)
	DefaultRateLimit RateLimitConfig `bson:"defaultRateLimit" json:"defaultRateLimit"`

	// CORS
	CORS CORSConfig `bson:"cors" json:"cors"`

	// JWT
	JWT JWTConfig `bson:"jwt" json:"jwt"`

	// API Info
	API APIConfig `bson:"api" json:"api"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// CORSConfig define a configuração de CORS
type CORSConfig struct {
	Enabled        bool     `bson:"enabled" json:"enabled"`
	AllowedOrigins []string `bson:"allowedOrigins" json:"allowedOrigins"`
}

// JWTConfig define a configuração de JWT
type JWTConfig struct {
	AccessTokenTTL  int `bson:"accessTokenTTL" json:"accessTokenTTL"`   // em segundos
	RefreshTokenTTL int `bson:"refreshTokenTTL" json:"refreshTokenTTL"` // em segundos
}

// APIConfig define informações da API
type APIConfig struct {
	Version     string `bson:"version" json:"version"`
	Environment string `bson:"environment" json:"environment"` // development, production
	Maintenance bool   `bson:"maintenance" json:"maintenance"`
}

// GetDefaultSettings retorna as configurações padrão do sistema
func GetDefaultSettings() *SystemSettings {
	return &SystemSettings{
		DefaultRateLimit: RateLimitConfig{
			RequestsPerDay:    1000, // 1k requests/dia para plano free
			RequestsPerMinute: 60,   // 60 requests/minuto
		},
		CORS: CORSConfig{
			Enabled: true,
			AllowedOrigins: []string{
				"https://core.theretech.com.br",
				"http://localhost:3000",
				"http://localhost:3001",
			},
		},
		JWT: JWTConfig{
			AccessTokenTTL:  900,    // 15 minutos
			RefreshTokenTTL: 604800, // 7 dias
		},
		API: APIConfig{
			Version:     "1.0.0",
			Environment: "development",
			Maintenance: false,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// UpdateSystemSettingsRequest representa o payload para atualizar configurações
type UpdateSystemSettingsRequest struct {
	DefaultRateLimit *RateLimitConfig `json:"defaultRateLimit,omitempty"`
	CORS             *CORSConfig      `json:"cors,omitempty"`
	JWT              *JWTConfig       `json:"jwt,omitempty"`
	API              *APIConfig       `json:"api,omitempty"`
}
