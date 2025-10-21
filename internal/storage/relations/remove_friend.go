package relations

import (
	"context"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) RemoveFriend(ctx context.Context, senderId string, toRemoveId string) error {
	s.log.Info("[RemoveFriend] storage started")

	filter := bson.M{"user_id": senderId}

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[RemoveFriend] storage error: " + err.Error())
		return errors.New("[RemoveFriend] storage error: " + err.Error())
	}
	senderRelation.FriendsIds = utils.RemoveByValue(senderRelation.FriendsIds, toRemoveId)

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[RemoveFriend] storage error: " + err.Error())
		return errors.New("[RemoveFriend] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": toRemoveId}

	recieverRelation, err := s.GetRelationRecord(ctx, toRemoveId)
	if err != nil {
		slog.Error("[RemoveFriend] storage error: " + err.Error())
		return errors.New("[RemoveFriend] storage error: " + err.Error())
	}

	recieverRelation.FriendsIds = utils.RemoveByValue(recieverRelation.FriendsIds, senderId)

	update = bson.M{
		"$set": recieverRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[RemoveFriend] storage error: " + err.Error())
		return errors.New("[RemoveFriend] storage error: " + err.Error())
	}

	return nil
}
