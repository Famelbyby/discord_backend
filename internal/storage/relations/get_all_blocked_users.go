package relations

import (
	"context"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) GetAllBlockedUsers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllBlockedUsers] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, storage.ErrUserNotFound
		}

		slog.Error("[GetAllBlockedUsers] storage error: " + err.Error())
		return nil, errors.New("[GetAllBlockedUsers] storage error: " + err.Error())
	}

	if (page-1)*limit > int64(len(senderRelation.BlockedIds)) {
		return []string{}, nil
	}
	return senderRelation.BlockedIds[(page-1)*limit : min((page)*limit, int64(len(senderRelation.BlockedIds)))], nil
}
