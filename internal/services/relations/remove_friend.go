package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) RemoveFriend(ctx context.Context, senderId string, toRemoveId string) error {
	r.log.Info("[RemoveFriend] service started")

	err := r.relations.RemoveFriend(ctx, senderId, toRemoveId)
	if err != nil {
		r.log.Error("[RemoveFriend] service error: " + err.Error())
		return fmt.Errorf("[RemoveFriend] service error: " + err.Error())
	}
	return nil
}
