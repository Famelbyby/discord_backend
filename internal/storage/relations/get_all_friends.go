package relations

import (
	"context"
	"errors"
	"log/slog"
)

func (s *RelationsStorage) GetAllFriends(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllFriends] storage started")

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[GetAllFriends] storage error: " + err.Error())
		return nil, errors.New("[GetAllFriends] storage error: " + err.Error())
	}

	if (page-1)*limit > int64(len(senderRelation.FriendsIds)) {
		return []string{}, nil
	}
	return senderRelation.FriendsIds[(page-1)*limit : min((page)*limit, int64(len(senderRelation.FriendsIds)))], nil
}
