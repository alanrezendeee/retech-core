package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

// LogActivity registra uma atividade no sistema (de forma assíncrona)
func LogActivity(
	c *gin.Context,
	repo *storage.ActivityLogsRepo,
	activityType string,
	action string,
	actor domain.Actor,
	resource domain.Resource,
	metadata map[string]interface{},
) {
	// Capturar IP e User-Agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Criar log
	log := &domain.ActivityLog{
		Timestamp: domain.TimeNow(),
		Type:      activityType,
		Action:    action,
		Actor:     actor,
		Resource:  resource,
		Metadata:  metadata,
		IP:        ip,
		UserAgent: userAgent,
	}

	// Salvar de forma assíncrona (não bloqueia a request)
	go func() {
		ctx := context.Background()
		if err := repo.Log(ctx, log); err != nil {
			// Log interno (não falhar a request por causa disso)
			// TODO: adicionar log estruturado aqui se necessário
		}
	}()
}

// BuildActorFromContext extrai informações do usuário do contexto Gin
func BuildActorFromContext(c *gin.Context) domain.Actor {
	// O middleware AuthJWT já deve ter setado essas informações
	userID, _ := c.Get("userId")
	email, _ := c.Get("userEmail")
	name, _ := c.Get("userName")
	role, _ := c.Get("userRole")

	return domain.Actor{
		UserID: toString(userID),
		Email:  toString(email),
		Name:   toString(name),
		Role:   toString(role),
	}
}

// toString converte interface{} para string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

