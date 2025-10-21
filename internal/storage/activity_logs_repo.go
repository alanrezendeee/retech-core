package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ActivityLogsRepo gerencia logs de atividades no MongoDB
type ActivityLogsRepo struct {
	col *mongo.Collection
}

// NewActivityLogsRepo cria um novo repositório de activity logs
func NewActivityLogsRepo(db *mongo.Database) *ActivityLogsRepo {
	return &ActivityLogsRepo{
		col: db.Collection("activity_logs"),
	}
}

// EnsureIndexes cria índices necessários para performance
func (r *ActivityLogsRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "timestamp", Value: -1}, // Ordenação DESC por timestamp
			},
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "actor.userId", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "resource.type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "resource.id", Value: 1},
			},
		},
	}

	_, err := r.col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("erro ao criar índices de activity_logs: %w", err)
	}

	fmt.Println("✅ Índices de activity_logs criados com sucesso")
	return nil
}

// Log cria um novo registro de atividade (assíncrono recomendado)
func (r *ActivityLogsRepo) Log(ctx context.Context, log *domain.ActivityLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now().UTC()
	}

	_, err := r.col.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("erro ao salvar activity log: %w", err)
	}

	return nil
}

// Recent retorna as atividades mais recentes (limite padrão: 20)
func (r *ActivityLogsRepo) Recent(ctx context.Context, limit int) ([]*domain.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Máximo de 100 por segurança
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}). // DESC
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar activity logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*domain.ActivityLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("erro ao decodificar activity logs: %w", err)
	}

	return logs, nil
}

// ByUser retorna atividades de um usuário específico
func (r *ActivityLogsRepo) ByUser(ctx context.Context, userID string, limit int) ([]*domain.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, bson.M{"actor.userId": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar activity logs por usuário: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*domain.ActivityLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("erro ao decodificar activity logs: %w", err)
	}

	return logs, nil
}

// ByType retorna atividades de um tipo específico
func (r *ActivityLogsRepo) ByType(ctx context.Context, eventType string, limit int) ([]*domain.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, bson.M{"type": eventType}, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar activity logs por tipo: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*domain.ActivityLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("erro ao decodificar activity logs: %w", err)
	}

	return logs, nil
}

// ByResource retorna atividades relacionadas a um recurso específico
func (r *ActivityLogsRepo) ByResource(ctx context.Context, resourceType, resourceID string, limit int) ([]*domain.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	filter := bson.M{
		"resource.type": resourceType,
		"resource.id":   resourceID,
	}

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar activity logs por recurso: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*domain.ActivityLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("erro ao decodificar activity logs: %w", err)
	}

	return logs, nil
}

// DeleteOlderThan deleta logs mais antigos que a duração especificada (para limpeza)
func (r *ActivityLogsRepo) DeleteOlderThan(ctx context.Context, duration time.Duration) (int64, error) {
	cutoffTime := time.Now().UTC().Add(-duration)

	result, err := r.col.DeleteMany(ctx, bson.M{
		"timestamp": bson.M{"$lt": cutoffTime},
	})
	if err != nil {
		return 0, fmt.Errorf("erro ao deletar activity logs antigos: %w", err)
	}

	return result.DeletedCount, nil
}

// Count retorna o total de logs (com filtro opcional)
func (r *ActivityLogsRepo) Count(ctx context.Context, filter bson.M) (int64, error) {
	if filter == nil {
		filter = bson.M{}
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar activity logs: %w", err)
	}

	return count, nil
}

