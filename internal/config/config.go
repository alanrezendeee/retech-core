package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string
	HTTPPort       string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MongoURI       string
	MongoDB        string
	EnableCORS     bool
	CORSOrigins    []string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	c := &Config{
		Env:          getenv("ENV", "development"),
		HTTPPort:     getenv("PORT", "8080"),
		MongoURI:     getenv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDB:      getenv("MONGO_DB", "retech_core"),
		EnableCORS:   getenv("CORS_ENABLE", "true") == "true",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("[config] ENV=%s PORT=%s MONGO_URI=%s DB=%s",
		c.Env, c.HTTPPort, c.MongoURI, c.MongoDB)
	return c
}

