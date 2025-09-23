package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) SendFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[SendFriendOffer] service started")

	//проверка на существование

	_, err := r.relations.SendFriendOffer(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[SendFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[SendFriendOffer] service error: " + err.Error())
	}

	return nil
}
