package relations

import (
	"context"
	"discord_backend/internal/storage"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) UnblockUser(ctx context.Context, senderId string, toUnblockId string) error {
	s.log.Info("[UnblockUser] storage started, sender id=" + senderId + " toUnblockId=" + toUnblockId)

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[UnblockUser] storage error: " + err.Error())
		return errors.New("[UnblockUser] storage error: " + err.Error())
	}
	senderRelation.BlockedIds = utils.RemoveByValue(senderRelation.BlockedIds, toUnblockId)

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[UnblockUser] storage error: " + err.Error())
		return errors.New("[UnblockUser] storage error: " + err.Error())
	}

	return nil
}
