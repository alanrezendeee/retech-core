package storage

import (
	"context"
	"time"

	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepo struct {
	coll *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{coll: db.Collection("users")}
}

// Create cria um novo usuário (com hash de senha)
func (r *UsersRepo) Create(ctx context.Context, user *domain.User, password string) error {
	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user.Password = string(hashedPassword)
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Active = true

	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Converter ObjectID para string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	
	return nil
}

// FindByEmail busca usuário por email
func (r *UsersRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID busca usuário por ID
func (r *UsersRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyPassword verifica se a senha está correta
func (r *UsersRepo) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// UpdateLastLogin atualiza o último login do usuário
func (r *UsersRepo) UpdateLastLogin(ctx context.Context, userID string) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"lastLogin": now, "updatedAt": now}},
	)
	return err
}

// Update atualiza um usuário
func (r *UsersRepo) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// UpdatePassword atualiza a senha do usuário
func (r *UsersRepo) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = r.coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"password":  string(hashedPassword),
			"updatedAt": time.Now(),
		}},
	)
	return err
}

// Delete deleta um usuário
func (r *UsersRepo) Delete(ctx context.Context, id string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListByTenant lista usuários de um tenant
func (r *UsersRepo) ListByTenant(ctx context.Context, tenantID string) ([]*domain.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{"tenantId": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ListAll lista todos os usuários (admin only)
func (r *UsersRepo) ListAll(ctx context.Context) ([]*domain.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Count retorna total de usuários
func (r *UsersRepo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

// CountByTenant retorna total de usuários de um tenant
func (r *UsersRepo) CountByTenant(ctx context.Context, tenantID string) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{"tenantId": tenantID})
}
