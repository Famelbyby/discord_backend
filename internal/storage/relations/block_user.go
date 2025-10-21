package relations

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) BlockUser(ctx context.Context, senderId string, toBlockId string) error {
	s.log.Info("[BlockUser] storage started, sender id=" + senderId + " toBlockId=" + toBlockId)

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[BlockUser] storage error: " + err.Error())
		return errors.New("[BlockUser] storage error: " + err.Error())
	}

	if !slices.Contains(senderRelation.BlockedIds, toBlockId) {
		senderRelation.BlockedIds = append(senderRelation.BlockedIds, toBlockId)
	}

	filter := bson.M{"user_id": senderId}
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
