package relations

import (
	"context"
	"discord_backend/internal/storage"
	"fmt"
)

func (r *RelationsService) UnblockUser(ctx context.Context, senderId string, toUnblockId string) error {
	r.log.Info("[UnblockUser] service started")
	relation, err := r.relations.GetUserRelation(ctx, senderId, toUnblockId)
	if err != nil {
		r.log.Error("[UnblockUser] service error: " + err.Error())
		return fmt.Errorf("[UnblockUser] service error: " + err.Error())
	}

	if !relation.IsBlocked {
		return storage.ErrUserAlreadyUnblocked
	}

	err = r.relations.UnblockUser(ctx, senderId, toUnblockId)
	if err != nil {
		r.log.Error("[UnblockUser] service error: " + err.Error())
		return fmt.Errorf("[UnblockUser] service error: " + err.Error())
	}
	return nil
}
