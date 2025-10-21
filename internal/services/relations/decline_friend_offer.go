package relations

import (
	"context"
	"discord_backend/internal/storage"
	"fmt"
)

func (r *RelationsService) DeclineFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[DeclineFriendOffer] service started")

	userRelation, err := r.relations.GetUserRelation(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[DeclineFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[DeclineFriendOffer] service error: " + err.Error())
	}

	if !userRelation.IsOutcoming {
		return storage.ErrNoOutgoingOffer
	}

	err = r.relations.DeclineFriendOffer(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[DeclineFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[DeclineFriendOffer] service error: " + err.Error())
	}
	return nil
}
