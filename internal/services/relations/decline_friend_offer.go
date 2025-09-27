package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) DeclineFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[DeclineFriendOffer] service started")

	err := r.relations.DeclineFriendOffer(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[DeclineFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[DeclineFriendOffer] service error: " + err.Error())
	}
	return nil
}
