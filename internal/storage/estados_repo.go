package storage

import (
	"context"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EstadosRepo struct {
	coll *mongo.Collection
}

func NewEstadosRepo(db *mongo.Database) *EstadosRepo {
	return &EstadosRepo{coll: db.Collection("estados")}
}

// FindAll retorna todos os estados
func (r *EstadosRepo) FindAll(ctx context.Context) ([]domain.Estado, error) {
	opts := options.Find().SetSort(bson.D{{Key: "nome", Value: 1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var estados []domain.Estado
	if err := cursor.All(ctx, &estados); err != nil {
		return nil, err
	}
	return estados, nil
}

// FindBySigla retorna um estado pela sigla
func (r *EstadosRepo) FindBySigla(ctx context.Context, sigla string) (*domain.Estado, error) {
	var estado domain.Estado
	err := r.coll.FindOne(ctx, bson.M{"sigla": sigla}).Decode(&estado)
	if err != nil {
		return nil, err
	}
	return &estado, nil
}

// FindByID retorna um estado pelo ID
func (r *EstadosRepo) FindByID(ctx context.Context, id int) (*domain.Estado, error) {
	var estado domain.Estado
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&estado)
	if err != nil {
		return nil, err
	}
	return &estado, nil
}

// InsertMany insere múltiplos estados
func (r *EstadosRepo) InsertMany(ctx context.Context, estados []domain.Estado) error {
	if len(estados) == 0 {
		return nil
	}

	docs := make([]interface{}, len(estados))
	now := time.Now()
	for i, e := range estados {
		e.CreatedAt = now
		e.UpdatedAt = now
		docs[i] = e
	}

	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

// Count retorna a quantidade de estados
func (r *EstadosRepo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

// DeleteAll remove todos os estados (usado para re-seed)
func (r *EstadosRepo) DeleteAll(ctx context.Context) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{})
	return err
}

