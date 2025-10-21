package relations

import (
	"context"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) UnblockUser(ctx context.Context, senderId string, toUnblockId string) error {
	s.log.Info("[UnblockUser] storage started, sender id=" + senderId + " toUnblockId=" + toUnblockId)

	filter := bson.M{"user_id": senderId}

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
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
