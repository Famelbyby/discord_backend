package relations

import (
	"context"
	"discord_backend/internal/storage"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) RemoveFriend(ctx context.Context, senderId string, toRemoveId string) error {
	s.log.Info("[RemoveFriend] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[RemoveFriend] storage error: " + err.Error())
		return errors.New("[RemoveFriend] storage error: " + err.Error())
	}

	if !slices.Contains(senderRelation.FriendsIds, toRemoveId) {
		return storage.ErrUserIsntFriend
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

	var recieverRelation dtoRelations

	err = s.collection.FindOne(ctx, filter).Decode(&recieverRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

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
