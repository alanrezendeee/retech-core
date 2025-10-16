package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongo(uri, db string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Mongo{Client: cl, DB: cl.Database(db)}, nil
}

func (m *Mongo) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

