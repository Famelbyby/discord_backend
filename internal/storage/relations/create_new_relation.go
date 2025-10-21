package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *RelationsStorage) CreateNewRelation(ctx context.Context, userId string) error {
	s.log.Info("[CreateNewRelation] storage started")

	senderRelation := models.RelationRecord{
		ID:          primitive.NewObjectID(),
		UserId:      userId,
		FriendsIds:  []string{},
		IncomingIds: []string{},
		OutgoingIds: []string{},
		BlockedIds:  []string{},
	}

	_, err := s.collection.InsertOne(ctx, senderRelation)
	if err != nil {
		s.log.Error("[CreateNewRelation] storage error: " + err.Error())
		return errors.New("[CreateNewRelation] storage error: " + err.Error())
	}

	return nil
}
