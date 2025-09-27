package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) GetUserRelation(ctx context.Context, senderId string, targetId string) (models.UserRelation, error) {
	s.log.Info("[GetUserRelation] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.UserRelation{}, storage.ErrUserNotFound
		}

		slog.Error("[GetUserRelation] storage error: " + err.Error())
		return models.UserRelation{}, errors.New("[GetUserRelation] storage error: " + err.Error())
	}

	relation := models.UserRelation{}

	if slices.Contains(senderRelation.FriendsIds, targetId) {
		relation.IsFriend = true
	}

	if slices.Contains(senderRelation.BlockedIds, targetId) {
		relation.IsBlocked = true
	}

	if slices.Contains(senderRelation.IncomingIds, targetId) {
		relation.IsIncoming = true
	}

	if slices.Contains(senderRelation.OutgoingIds, targetId) {
		relation.IsOutcoming = true
	}

	return relation, nil
}
