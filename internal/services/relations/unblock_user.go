package relations

import (
	"context"
	"fmt"
)

func (r *RelationsService) UnblockUser(ctx context.Context, senderId string, toUnblockId string) error {
	r.log.Info("[UnblockUser] service started")

	err := r.relations.UnblockUser(ctx, senderId, toUnblockId)
	if err != nil {
		r.log.Error("[UnblockUser] service error: " + err.Error())
		return fmt.Errorf("[UnblockUser] service error: " + err.Error())
	}
	return nil
}
