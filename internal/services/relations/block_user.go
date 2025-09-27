package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) BlockUser(ctx context.Context, senderId string, toBlockId string) error {
	relation, err := r.relations.GetUserRelation(ctx, senderId, toBlockId)
	if err != nil {
		r.log.Error("[BlockUser] service error: " + err.Error())
		return fmt.Errorf("[BlockUser] service error: " + err.Error())
	}
	if relation.IsFriend {
		err = r.relations.RemoveFriend(ctx, senderId, toBlockId)
		if err != nil {
			r.log.Error("[BlockUser] service error: " + err.Error())
			return fmt.Errorf("[BlockUser] service error: " + err.Error())
		}
	}

	if relation.IsOutcoming {
		err = r.relations.CancelFriendOffer(ctx, senderId, toBlockId)
		if err != nil {
			r.log.Error("[BlockUser] service error: " + err.Error())
			return fmt.Errorf("[BlockUser] service error: " + err.Error())
		}
	}

	if relation.IsIncoming {
		err = r.relations.DeclineFriendOffer(ctx, toBlockId, senderId)
		if err != nil {
			r.log.Error("[BlockUser] service error: " + err.Error())
			return fmt.Errorf("[BlockUser] service error: " + err.Error())
		}
	}

	err = r.relations.BlockUser(ctx, senderId, toBlockId)
	if err != nil {
		r.log.Error("[BlockUser] service error: " + err.Error())
		return fmt.Errorf("[BlockUser] service error: " + err.Error())
	}

	return nil
}
