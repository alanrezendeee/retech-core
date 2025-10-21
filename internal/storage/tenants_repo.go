package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/theretech/retech-core/internal/domain"
)

type TenantsRepo struct {
	col *mongo.Collection
}

func NewTenantsRepo(db *mongo.Database) *TenantsRepo {
	return &TenantsRepo{col: db.Collection("tenants")}
}

func (r *TenantsRepo) Insert(ctx context.Context, t *domain.Tenant) error {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, t)
	return err
}

func (r *TenantsRepo) ByTenantID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.col.FindOne(ctx, bson.M{"tenantId": tenantID}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &t, err
}

func (r *TenantsRepo) List(ctx context.Context) ([]*domain.Tenant, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tenants []*domain.Tenant
	if err := cursor.All(ctx, &tenants); err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantsRepo) Update(ctx context.Context, tenantID string, updates map[string]interface{}) error {
	updates["updatedAt"] = time.Now().UTC()
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"tenantId": tenantID},
		bson.M{"$set": updates},
	)
	return err
}

func (r *TenantsRepo) Delete(ctx context.Context, idOrTenantID string) error {
	// Tentar deletar por tenantId primeiro
	result, err := r.col.DeleteOne(ctx, bson.M{"tenantId": idOrTenantID})
	if err != nil {
		return err
	}

	// Se não deletou nada (tenantId vazio ou não encontrado), tentar por _id
	if result.DeletedCount == 0 {
		// Tentar converter para ObjectID
		if objectID, err := primitive.ObjectIDFromHex(idOrTenantID); err == nil {
			_, err = r.col.DeleteOne(ctx, bson.M{"_id": objectID})
			return err
		}
	}

	return nil
}
