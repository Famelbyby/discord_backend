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

func (s *RelationsStorage) AcceptFriendRequest(ctx context.Context, senderId string, recieverId string) error {
	s.log.Info("[AcceptFriendRequest] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}
	senderRelation.OutgoingIds = utils.RemoveByValue(senderRelation.OutgoingIds, recieverId)
	senderRelation.FriendsIds = append(senderRelation.OutgoingIds, recieverId)

	_, err = s.collection.InsertOne(ctx, senderRelation)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	var recieverRelation dtoRelations

	err = s.collection.FindOne(ctx, filter).Decode(&recieverRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}
	recieverRelation.IncomingIds = utils.RemoveByValue(recieverRelation.IncomingIds, senderId)
	recieverRelation.FriendsIds = append(recieverRelation.IncomingIds, senderId)

	_, err = s.collection.InsertOne(ctx, recieverRelation)
	if err != nil {
		slog.Error("[AcceptFriendRequest] storage error: " + err.Error())
		return errors.New("[AcceptFriendRequest] storage error: " + err.Error())
	}

	return nil
}
