package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/theretech/retech-core/internal/domain"
)

type TokensRepo struct{ col *mongo.Collection }

func NewTokensRepo(db *mongo.Database) *TokensRepo {
	return &TokensRepo{col: db.Collection("refresh_tokens")}
}

func (r *TokensRepo) Insert(ctx context.Context, t *domain.RefreshToken) error {
	t.CreatedAt = time.Now().UTC()
	_, err := r.col.InsertOne(ctx, t); return err
}

func (r *TokensRepo) Revoke(ctx context.Context, jti string) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": jti}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}

func (r *TokensRepo) FindActive(ctx context.Context, jti string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.col.FindOne(ctx, bson.M{"_id": jti, "revoked": false}).Decode(&t)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &t, err
}

