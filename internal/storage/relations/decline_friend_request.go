package relations

import (
	"context"
	"discord_backend/internal/utils"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) DeclineFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	s.log.Info("[DeclineFriendOffer] storage started")

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}

	senderRelation.OutgoingIds = utils.RemoveByValue(senderRelation.OutgoingIds, recieverId)

	update := bson.M{
		"$set": senderRelation,
	}
	filter := bson.M{"user_id": senderId}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[DeclineFriendOffer] storage error: " + err.Error())
		return errors.New("[DeclineFriendOffer] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	recieverRelation, err := s.GetRelationRecord(ctx, recieverId)
	if err != nil {
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
