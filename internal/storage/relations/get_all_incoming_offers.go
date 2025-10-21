package relations

import (
	"context"
	"errors"
	"log/slog"
)

func (s *RelationsStorage) GetAllIncomingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllIncomingOffers] storage started")

	senderRelation, err := s.GetRelationRecord(ctx, senderId)
	if err != nil {
		slog.Error("[GetAllIncomingOffers] storage error: " + err.Error())
		return nil, errors.New("[GetAllIncomingOffers] storage error: " + err.Error())
	}

	if (page-1)*limit > int64(len(senderRelation.IncomingIds)) {
		return []string{}, nil
	}
	return senderRelation.IncomingIds[(page-1)*limit : min((page)*limit, int64(len(senderRelation.IncomingIds)))], nil
}
