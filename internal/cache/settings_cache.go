package cache

import (
	"context"
	"sync"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

// SettingsCache mantém settings em memória para evitar consultas MongoDB a cada request
type SettingsCache struct {
	mu       sync.RWMutex
	settings *domain.SystemSettings
	lastLoad time.Time
	ttl      time.Duration
	repo     *storage.SettingsRepo
}

// NewSettingsCache cria um novo cache de settings
func NewSettingsCache(repo *storage.SettingsRepo) *SettingsCache {
	return &SettingsCache{
		repo: repo,
		ttl:  30 * time.Second, // Cache de 30 segundos (balance entre freshness e performance)
	}
}

// Get retorna settings do cache ou carrega do MongoDB se expirado
func (sc *SettingsCache) Get(ctx context.Context) (*domain.SystemSettings, error) {
	sc.mu.RLock()
	
	// Se cache ainda é válido, retorna imediatamente
	if sc.settings != nil && time.Since(sc.lastLoad) < sc.ttl {
		sc.mu.RUnlock()
		return sc.settings, nil
	}
	
	sc.mu.RUnlock()
	
	// Cache expirado ou vazio, precisa carregar
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	// Double-check: outro goroutine pode ter carregado enquanto esperávamos o lock
	if sc.settings != nil && time.Since(sc.lastLoad) < sc.ttl {
		return sc.settings, nil
	}
	
	// Carregar do MongoDB
	settings, err := sc.repo.Get(ctx)
	if err != nil {
		// Fallback para defaults se falhar
		defaultSettings := domain.GetDefaultSettings()
		return defaultSettings, nil
	}
	
	// Atualizar cache
	sc.settings = settings
	sc.lastLoad = time.Now()
	
	return sc.settings, nil
}

// Invalidate força recarregar do MongoDB na próxima requisição
func (sc *SettingsCache) Invalidate() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.lastLoad = time.Time{} // Zero time = expirado
}

// Refresh recarrega imediatamente do MongoDB
func (sc *SettingsCache) Refresh(ctx context.Context) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	settings, err := sc.repo.Get(ctx)
	if err != nil {
		return err
	}
	
	sc.settings = settings
	sc.lastLoad = time.Now()
	
	return nil
}

