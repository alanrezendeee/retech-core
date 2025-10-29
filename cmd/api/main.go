package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/bootstrap"
	"github.com/theretech/retech-core/internal/cache"
	"github.com/theretech/retech-core/internal/config"
	nethttp "github.com/theretech/retech-core/internal/http"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/observability"
	"github.com/theretech/retech-core/internal/storage"
)

func main() {
	// Validar ENVs obrigatórias (falha rápido se não configuradas)
	config.ValidateExternalAPIsConfig()
	
	cfg := config.Load()
	log := observability.NewLogger(cfg.Env)

	// Mongo
	m, err := storage.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal().Err(err).Msg("mongo_connect_error")
	}

	// Redis (opcional - graceful degradation se não configurado)
	var redisClient interface{} // interface{} para permitir nil
	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		client, err := cache.NewRedisClient(redisURL, "", 0, log)
		if err != nil {
			log.Warn().Err(err).Msg("⚠️  Redis não disponível, usando apenas MongoDB (performance reduzida)")
			redisClient = nil // Continua sem Redis (fallback gracioso)
		} else {
			log.Info().Msg("✅ Redis conectado - cache ultra-rápido habilitado!")
			redisClient = client
		}
	} else {
		log.Warn().Msg("⚠️  REDIS_URL não configurado, usando apenas MongoDB")
		redisClient = nil
	}

	// Executar migrations/seeds
	log.Info().Msg("Executando migrations e seeds...")
	migrationManager := bootstrap.NewMigrationManager(m.DB, log)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := migrationManager.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("migration_error")
	}
	log.Info().Msg("Migrations concluídas com sucesso")

	// Criar índices
	if err := bootstrap.CreateIndexes(ctx, m.DB, log); err != nil {
		log.Warn().Err(err).Msg("index_creation_warning")
	}

	// Migrar configurações (adicionar campos novos)
	if err := bootstrap.MigrateSettings(ctx, m.DB, log); err != nil {
		log.Warn().Err(err).Msg("settings_migration_warning")
	}

	// Repos
	tenants := storage.NewTenantsRepo(m.DB)
	apikeys := storage.NewAPIKeysRepo(m.DB)
	users := storage.NewUsersRepo(m.DB)
	estados := storage.NewEstadosRepo(m.DB)
	municipios := storage.NewMunicipiosRepo(m.DB)

	// Settings
	settings := storage.NewSettingsRepo(m.DB)

	// Garantir que configurações padrão existam
	if err := settings.Ensure(context.Background()); err != nil {
		log.Warn().Err(err).Msg("failed to ensure default settings")
	}

	// Activity Logs
	activityLogs := storage.NewActivityLogsRepo(m.DB)

	// Criar índices para activity logs
	if err := activityLogs.EnsureIndexes(context.Background()); err != nil {
		log.Warn().Err(err).Msg("failed to create activity logs indexes")
	}

	// 🎯 DESABILITADO: API Key Demo agora é criada manualmente via admin/settings
	// Usar botão "Gerar Nova" ou "Rotacionar" em /admin/settings
	// if err := bootstrap.EnsureDemoAPIKey(context.Background(), apikeys, tenants, settings, m.DB); err != nil {
	// 	log.Warn().Err(err).Msg("failed to ensure demo API key")
	// }

	// JWT Service
	jwtService := auth.NewJWTService(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)

	// Router
	health := handlers.NewHealthHandler(m.Client, redisClient) // ✅ Passar Redis também
	router := nethttp.NewRouter(log, m, redisClient, health, apikeys, tenants, users, estados, municipios, settings, activityLogs, jwtService)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info().Msgf("listening on :%s (env=%s)", cfg.HTTPPort, cfg.Env)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server_error")
	}
	fmt.Println()
}
