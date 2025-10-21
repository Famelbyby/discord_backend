package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) AcceptFriendRequest(ctx context.Context, senderId string, recieverId string) error {
	s.log.Info("[AcceptFriendRequest] storage started")

	var senderRelation models.RelationRecord
	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	if !slices.Contains(senderRelation.OutgoingIds, recieverId) {
		slog.Error("[AcceptFriendRequest] storage: friend offer is outdated")
		return errors.New("friend offer is outdated")
	}
	senderRelation.OutgoingIds = utils.RemoveByValue(senderRelation.OutgoingIds, recieverId)
	senderRelation.FriendsIds = append(senderRelation.FriendsIds, recieverId)

	update := bson.M{
		"$set": senderRelation,
	}

	filter := bson.M{"user_id": senderId}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	var recieverRelation models.RelationRecord
	recieverRelation, err = s.GetRelationRecord(ctx, recieverId)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	recieverRelation.IncomingIds = utils.RemoveByValue(recieverRelation.IncomingIds, senderId)
	recieverRelation.FriendsIds = append(recieverRelation.FriendsIds, senderId)

	update = bson.M{
		"$set": recieverRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	return nil
}
