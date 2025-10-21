package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/bootstrap"
	"github.com/theretech/retech-core/internal/config"
	nethttp "github.com/theretech/retech-core/internal/http"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/observability"
	"github.com/theretech/retech-core/internal/storage"
)

func main() {
	cfg := config.Load()
	log := observability.NewLogger(cfg.Env)

	// Mongo
	m, err := storage.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal().Err(err).Msg("mongo_connect_error")
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

	// Repos
	tenants := storage.NewTenantsRepo(m.DB)
	apikeys := storage.NewAPIKeysRepo(m.DB)
	users := storage.NewUsersRepo(m.DB)
	estados := storage.NewEstadosRepo(m.DB)
	municipios := storage.NewMunicipiosRepo(m.DB)

	// JWT Service
	jwtService := auth.NewJWTService(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)

	// Router
	health := handlers.NewHealthHandler(m.Client)
	router := nethttp.NewRouter(log, m, health, apikeys, tenants, users, estados, municipios, jwtService)

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
