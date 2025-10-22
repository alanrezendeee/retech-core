package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/theretech/retech-core/internal/domain"
)

type APIKeysRepo struct{ col *mongo.Collection }

func NewAPIKeysRepo(db *mongo.Database) *APIKeysRepo {
	return &APIKeysRepo{col: db.Collection("api_keys")}
}

func (r *APIKeysRepo) Insert(ctx context.Context, k *domain.APIKey) error {
	k.CreatedAt = time.Now().UTC()
	_, err := r.col.InsertOne(ctx, k)
	return err
}

func (r *APIKeysRepo) Revoke(ctx context.Context, keyId string) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"keyId": keyId}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}

func (r *APIKeysRepo) ByKeyID(ctx context.Context, keyId string) (*domain.APIKey, error) {
	var a domain.APIKey
	err := r.col.FindOne(ctx, bson.M{"keyId": keyId, "revoked": false}).Decode(&a)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &a, err
}

func (r *APIKeysRepo) ByKeyIDAny(ctx context.Context, keyId string) (*domain.APIKey, error) {
	var a domain.APIKey
	err := r.col.FindOne(ctx, bson.M{"keyId": keyId}).Decode(&a)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &a, err
}

// CountByOwner retorna total de API keys ativas de um tenant
func (r *APIKeysRepo) CountByOwner(ctx context.Context, ownerId string) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{
		"ownerId": ownerId,
		"revoked": false,
	})
}
