package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func FetchCollection(ctx context.Context, client *mongo.Client, dbName, dbCollection string, limit ...int64) ([]bson.M, error) {
	collection := client.Database(dbName).Collection(dbCollection)

	opts := options.Find()
	if len(limit) > 0 {
		opts.SetLimit(limit[0])
	}

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch documents: %w", err)
	}

	return results, nil
}
