package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"fmt"
)

func (r *RelationsService) GetUserRelation(ctx context.Context, senderId string, targetId string) (models.UserRelation, error) {
	r.log.Info("[GetUserRelation] service started")

	relation, err := r.relations.GetUserRelation(ctx, senderId, targetId)
	if err != nil {
		r.log.Error("[GetUserRelation] service error: " + err.Error())
		return models.UserRelation{}, fmt.Errorf("[GetUserRelation] service error: " + err.Error())
	}
	return relation, nil
}
