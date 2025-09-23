package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) AcceptFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[AcceptFriendOffer] service started")

	//проверка на блокировку

	err := r.relations.AcceptFriendRequest(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[AcceptFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[AcceptFriendOffer] service error: " + err.Error())
	}
	return nil
}
