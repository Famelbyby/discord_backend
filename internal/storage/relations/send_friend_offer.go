package relations

import (
	"context"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) SendFriendOffer(ctx context.Context, senderId, recieverId string) (string, error) {
	s.log.Info("[SendFriendRequest] storage started, sender id=" + senderId + " recieverId=" + recieverId)

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", storage.ErrUserNotFound
		}

		slog.Error("[SendFriendRequest] storage error: " + err.Error())
		return "", errors.New("[SendFriendRequest] storage error: " + err.Error())
	}
	senderRelation.OutgoingIds = append(senderRelation.OutgoingIds, recieverId)

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[SendFriendRequest] storage error: " + err.Error())
		return "", errors.New("[SendFriendRequest] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	var recieverRelation dtoRelations

	err = s.collection.FindOne(ctx, filter).Decode(&recieverRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", storage.ErrUserNotFound
		}

		slog.Error("[SendFriendRequest] storage error: " + err.Error())
		return "", errors.New("[SendFriendRequest] storage error: " + err.Error())
	}
	recieverRelation.IncomingIds = append(recieverRelation.IncomingIds, senderId)

	update = bson.M{
		"$set": recieverRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[SendFriendRequest] storage error: " + err.Error())
		return "", errors.New("[SendFriendRequest] storage error: " + err.Error())
	}

	return recieverId, nil
}
