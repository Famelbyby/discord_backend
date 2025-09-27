package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) GetAllFriends(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	r.log.Info("[GetAllFriends] service started")

	ids, err := r.relations.GetAllFriends(ctx, senderId, page, limit)
	if err != nil {
		r.log.Error("[GetAllFriends] service error: " + err.Error())
		return nil, fmt.Errorf("[GetAllFriends] service error: " + err.Error())
	}
	return ids, nil
}
