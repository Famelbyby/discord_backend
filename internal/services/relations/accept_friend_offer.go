package relations

import (
	"context"
	"discord_backend/internal/storage"
	"fmt"
)

func (r *RelationsService) AcceptFriendOffer(ctx context.Context, senderId string, recieverId string) error {
	r.log.Info("[AcceptFriendOffer] service started")

	relation, err := r.relations.GetUserRelation(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[AcceptFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[AcceptFriendOffer] service error: " + err.Error())
	}

	if relation.IsFriend {
		return storage.ErrCantAddAlreadyFriend
	}

	if !relation.IsOutcoming {
		return storage.ErrNoOutgoingOffer
	}

	err = r.relations.AcceptFriendRequest(ctx, senderId, recieverId)
	if err != nil {
		r.log.Error("[AcceptFriendOffer] service error: " + err.Error())
		return fmt.Errorf("[AcceptFriendOffer] service error: " + err.Error())
	}
	return nil
}
