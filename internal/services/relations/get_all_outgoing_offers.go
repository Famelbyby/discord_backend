package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) GetAllOutgoingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	r.log.Info("[GetAllOutgoingOffers] service started")

	ids, err := r.relations.GetAllOutgoingOffers(ctx, senderId, page, limit)
	if err != nil {
		r.log.Error("[GetAllOutgoingOffers] service error: " + err.Error())
		return nil, fmt.Errorf("[GetAllOutgoingOffers] service error: " + err.Error())
	}
	return ids, nil
}
