package relations

import (
	"context"
	"errors"
	"log/slog"
)

func (s *RelationsStorage) GetAllOutgoingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllOutgoingOffers] storage started")

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[GetAllOutgoingOffers] storage error: " + err.Error())
		return nil, errors.New("[GetAllOutgoingOffers] storage error: " + err.Error())
	}

	if (page-1)*limit > int64(len(senderRelation.OutgoingIds)) {
		return []string{}, nil
	}
	return senderRelation.OutgoingIds[(page-1)*limit : min((page)*limit, int64(len(senderRelation.OutgoingIds)))], nil
}
