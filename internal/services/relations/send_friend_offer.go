package relations

import (
	"context"
	"discord_backend/internal/storage"
	"fmt"
)

func (r *RelationsService) SendFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[SendFriendOffer] service started")

	relation, err := r.relations.GetUserRelation(ctx, recieverId, senderId)
	if err != nil {
		r.log.Error("[SendFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[SendFriendOffer] service error: " + err.Error())
	}
	if relation.IsBlocked {
		return storage.ErrCantAddBlockedUser
	}
	if relation.IsFriend {
		return storage.ErrCantAddAlreadyFriend
	}
	if relation.IsIncoming {
		return storage.ErrFriendOfferAlreadySent
	}

	err = r.relations.SendFriendOffer(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[SendFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[SendFriendOffer] service error: " + err.Error())
	}

	return nil
}
