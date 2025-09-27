package relations

import (
	"context"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) BlockUser(ctx context.Context, senderId string, toBlockId string) error {
	s.log.Info("[BlockUser] storage started, sender id=" + senderId + " toBlockId=" + toBlockId)

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[BlockUser] storage error: " + err.Error())
		return errors.New("[BlockUser] storage error: " + err.Error())
	}
	if !slices.Contains(senderRelation.BlockedIds, toBlockId) {
		senderRelation.BlockedIds = append(senderRelation.BlockedIds, toBlockId)
	}

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[BlockUser] storage error: " + err.Error())
		return errors.New("[BlockUser] storage error: " + err.Error())
	}

	return nil
}
