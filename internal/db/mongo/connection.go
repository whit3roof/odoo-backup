package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/byt3roof/odoo-backup/internal/conf"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Client, error) {
	cfg, err := conf.LoadConfig()
	if err != nil {
		return nil, err
	}

	uri := cfg.MongoUri

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connection error: %w", err)
	}

	return client, nil
}
