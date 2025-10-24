package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// RedisClient encapsula o cliente Redis
type RedisClient struct {
	client *redis.Client
	log    zerolog.Logger
}

// NewRedisClient cria um novo cliente Redis
// Aceita tanto URL completa (redis://...) quanto addr simples (host:port)
func NewRedisClient(urlOrAddr, password string, db int, log zerolog.Logger) (*RedisClient, error) {
	var client *redis.Client
	
	// Se começa com "redis://", parsear como URL
	if len(urlOrAddr) > 8 && urlOrAddr[:8] == "redis://" {
		opt, err := redis.ParseURL(urlOrAddr)
		if err != nil {
			log.Error().Err(err).Str("url", urlOrAddr).Msg("❌ Erro ao parsear REDIS_URL")
			return nil, fmt.Errorf("erro ao parsear REDIS_URL: %w", err)
		}
		
		// Configurar pool e timeouts
		opt.PoolSize = 50
		opt.MinIdleConns = 10
		opt.DialTimeout = 5 * time.Second
		opt.ReadTimeout = 3 * time.Second
		opt.WriteTimeout = 3 * time.Second
		
		client = redis.NewClient(opt)
	} else {
		// Formato antigo: addr + password separados
		client = redis.NewClient(&redis.Options{
			Addr:         urlOrAddr,
			Password:     password,
			DB:           db,
			PoolSize:     50,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})
	}

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Str("url", urlOrAddr).Msg("❌ Falha ao conectar no Redis")
		return nil, err
	}

	log.Info().Msg("⚡ Redis conectado - cache ultra-rápido habilitado!")

	return &RedisClient{
		client: client,
		log:    log,
	}, nil
}

// Get busca um valor no Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Chave não existe
	}
	return val, err
}

// Set salva um valor no Redis com TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Serializar para JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// SetString salva uma string no Redis com TTL
func (r *RedisClient) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete remove uma chave do Redis
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// FlushPattern remove todas as chaves que começam com pattern
func (r *RedisClient) FlushPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	
	keys := []string{}
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	
	if err := iter.Err(); err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	
	return nil
}

// Exists verifica se uma chave existe
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// TTL retorna o tempo de vida restante de uma chave
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Close fecha a conexão com Redis
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Ping testa a conexão
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

