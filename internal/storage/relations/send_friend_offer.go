package relations

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *RelationsStorage) SendFriendOffer(ctx context.Context, senderId, recieverId string) error {
	s.log.Info("[SendFriendOffer] storage started, sender id=" + senderId + " recieverId=" + recieverId)

	filter := bson.M{"user_id": senderId}

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[SendFriendOffer] storage error: " + err.Error())
		return errors.New("[SendFriendOffer] storage error: " + err.Error())
	}
	if slices.Contains(senderRelation.OutgoingIds, recieverId) {
		return errors.New("friend offer is already sent")
	}
	senderRelation.OutgoingIds = append(senderRelation.OutgoingIds, recieverId)

	update := bson.M{
		"$set": senderRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[SendFriendOffer] storage error: " + err.Error())
		return errors.New("[SendFriendOffer] storage error: " + err.Error())
	}

	filter = bson.M{"user_id": recieverId}

	recieverRelation, err := s.GetRelationRecord(ctx, recieverId)
	if err != nil {
		slog.Error("[SendFriendOffer] storage error: " + err.Error())
		return errors.New("[SendFriendOffer] storage error: " + err.Error())
	}

	recieverRelation.IncomingIds = append(recieverRelation.IncomingIds, senderId)

	update = bson.M{
		"$set": recieverRelation,
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("[SendFriendOffer] storage error: " + err.Error())
		return errors.New("[SendFriendOffer] storage error: " + err.Error())
	}

	return nil
}
