package relations

import (
	"context"
	"errors"
	"log/slog"
)

func (s *RelationsStorage) GetAllBlockedUsers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllBlockedUsers] storage started")

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[GetAllBlockedUsers] storage error: " + err.Error())
		return nil, errors.New("[GetAllBlockedUsers] storage error: " + err.Error())
	}

	if (page-1)*limit > int64(len(senderRelation.BlockedIds)) {
		return []string{}, nil
	}
	return senderRelation.BlockedIds[(page-1)*limit : min((page)*limit, int64(len(senderRelation.BlockedIds)))], nil
}
