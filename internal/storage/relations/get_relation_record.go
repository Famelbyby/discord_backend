package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) GetRelationRecord(ctx context.Context, userId string) (models.RelationRecord, error) {
	s.log.Info("[GetRelationRecord] storage started")

	filter := bson.M{"user_id": userId}

	var userRelation models.RelationRecord

	err := s.collection.FindOne(ctx, filter).Decode(&userRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.RelationRecord{}, storage.ErrUserNotFound
		}

		slog.Error("[GetRelationRecord] storage error: " + err.Error())
		return models.RelationRecord{}, errors.New("[GetRelationRecord] storage error: " + err.Error())
	}

	return userRelation, nil
}
