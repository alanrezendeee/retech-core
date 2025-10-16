package main

import (
	"fmt"
	"net/http"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/config"
	"github.com/theretech/retech-core/internal/http/handlers"
	nethttp "github.com/theretech/retech-core/internal/http"
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

	// Repos
	users := storage.NewUsersRepo(m.DB)
	tokens := storage.NewTokensRepo(m.DB)
	apikeys := storage.NewAPIKeysRepo(m.DB)

	// JWT
	jwtSvc, err := auth.NewJWTService()
	if err != nil {
		log.Fatal().Err(err).Msg("jwt")
	}

	// Router
	health := handlers.NewHealthHandler(m.Client)
	router := nethttp.NewRouter(log, health, jwtSvc, users, tokens, apikeys)

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

