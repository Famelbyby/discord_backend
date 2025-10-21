package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"errors"
	"log/slog"
	"slices"
)

func (s *RelationsStorage) GetUserRelation(ctx context.Context, senderId string, targetId string) (models.UserRelation, error) {
	s.log.Info("[GetUserRelation] storage started")

	senderRelationRecord, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[GetUserRelation] storage error: " + err.Error())
		return models.UserRelation{}, errors.New("[GetUserRelation] storage error: " + err.Error())
	}

	relation := models.UserRelation{}

	if slices.Contains(senderRelationRecord.FriendsIds, targetId) {
		relation.IsFriend = true
	}

	if slices.Contains(senderRelationRecord.BlockedIds, targetId) {
		relation.IsBlocked = true
	}

	if slices.Contains(senderRelationRecord.IncomingIds, targetId) {
		relation.IsIncoming = true
	}

	if slices.Contains(senderRelationRecord.OutgoingIds, targetId) {
		relation.IsOutcoming = true
	}

	return relation, nil
}
