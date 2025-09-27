package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) GetAllBlockedUsers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	r.log.Info("[GetAllBlockedUsers] service started")

	ids, err := r.relations.GetAllBlockedUsers(ctx, senderId, page, limit)
	if err != nil {
		r.log.Error("[GetAllBlockedUsers] service error: " + err.Error())
		return nil, fmt.Errorf("[GetAllBlockedUsers] service error: " + err.Error())
	}
	return ids, nil
}
