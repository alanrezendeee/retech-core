package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/theretech/retech-core/internal/domain"
)

type UsersRepo struct{ col *mongo.Collection }

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{col: db.Collection("users")}
}

func (r *UsersRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &u, err
}

func (r *UsersRepo) ByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &u, err
}

func (r *UsersRepo) Insert(ctx context.Context, u *domain.User) error {
	u.CreatedAt = time.Now().UTC()
	u.Active = true
	_, err := r.col.InsertOne(ctx, u)
	return err
}

