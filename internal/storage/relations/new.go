package relations

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RelationsStorage struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func NewStorage(db *mongo.Database, collectionName string, log *slog.Logger) *RelationsStorage {
	storage := &RelationsStorage{
		collection: db.Collection(collectionName),
		log:        log,
	}

	return storage
}
