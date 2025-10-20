package storage

import (
	"context"
	"strings"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MunicipiosRepo struct {
	coll *mongo.Collection
}

func NewMunicipiosRepo(db *mongo.Database) *MunicipiosRepo {
	return &MunicipiosRepo{coll: db.Collection("municipios")}
}

// FindAll retorna todos os municípios
func (r *MunicipiosRepo) FindAll(ctx context.Context) ([]domain.Municipio, error) {
	opts := options.Find().SetSort(bson.D{{Key: "nome", Value: 1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var municipios []domain.Municipio
	if err := cursor.All(ctx, &municipios); err != nil {
		return nil, err
	}
	return municipios, nil
}

// FindByUF retorna todos os municípios de um estado
func (r *MunicipiosRepo) FindByUF(ctx context.Context, uf string) ([]domain.Municipio, error) {
	uf = strings.ToUpper(uf)
	filter := bson.M{"microrregiao.mesorregiao.UF.sigla": uf}
	opts := options.Find().SetSort(bson.D{{Key: "nome", Value: 1}})
	
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var municipios []domain.Municipio
	if err := cursor.All(ctx, &municipios); err != nil {
		return nil, err
	}
	return municipios, nil
}

// FindByID retorna um município pelo ID do IBGE
func (r *MunicipiosRepo) FindByID(ctx context.Context, id int) (*domain.Municipio, error) {
	var municipio domain.Municipio
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&municipio)
	if err != nil {
		return nil, err
	}
	return &municipio, nil
}

// Search busca municípios por nome (case-insensitive, parcial)
func (r *MunicipiosRepo) Search(ctx context.Context, query string, uf string) ([]domain.Municipio, error) {
	filter := bson.M{
		"nome": bson.M{"$regex": query, "$options": "i"},
	}
	
	if uf != "" {
		filter["microrregiao.mesorregiao.UF.sigla"] = strings.ToUpper(uf)
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "nome", Value: 1}}).SetLimit(100)
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var municipios []domain.Municipio
	if err := cursor.All(ctx, &municipios); err != nil {
		return nil, err
	}
	return municipios, nil
}

// InsertMany insere múltiplos municípios
func (r *MunicipiosRepo) InsertMany(ctx context.Context, municipios []domain.Municipio) error {
	if len(municipios) == 0 {
		return nil
	}

	// Inserir em lotes para não estourar memória
	batchSize := 1000
	now := time.Now()
	
	for i := 0; i < len(municipios); i += batchSize {
		end := i + batchSize
		if end > len(municipios) {
			end = len(municipios)
		}
		
		batch := municipios[i:end]
		docs := make([]interface{}, len(batch))
		for j, m := range batch {
			m.CreatedAt = now
			m.UpdatedAt = now
			docs[j] = m
		}
		
		if _, err := r.coll.InsertMany(ctx, docs); err != nil {
			return err
		}
	}
	
	return nil
}

// Count retorna a quantidade de municípios
func (r *MunicipiosRepo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

// DeleteAll remove todos os municípios (usado para re-seed)
func (r *MunicipiosRepo) DeleteAll(ctx context.Context) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{})
	return err
}

