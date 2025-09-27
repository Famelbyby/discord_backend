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

func (s *RelationsStorage) DeclineFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	s.log.Info("[DeclineFriendOffer] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}

	senderRelation.OutgoingIds = utils.RemoveByValue(senderRelation.OutgoingIds, recieverId)

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	var recieverRelation dtoRelations

	err = s.collection.FindOne(ctx, filter).Decode(&recieverRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}
	recieverRelation.IncomingIds = utils.RemoveByValue(recieverRelation.IncomingIds, senderId)

	update = bson.M{
		"$set": recieverRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}

	return nil
}
