package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) CancelFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[CancelFriendOffer] service started")

	err := r.relations.CancelFriendOffer(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[CancelFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[CancelFriendOffer] service error: " + err.Error())
	}
	return nil
}
