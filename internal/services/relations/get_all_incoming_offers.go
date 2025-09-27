package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) GetAllIncomingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	r.log.Info("[GetAllIncomingOffers] service started")

	ids, err := r.relations.GetAllIncomingOffers(ctx, senderId, page, limit)
	if err != nil {
		r.log.Error("[GetAllIncomingOffers] service error: " + err.Error())
		return nil, fmt.Errorf("[GetAllIncomingOffers] service error: " + err.Error())
	}
	return ids, nil
}
