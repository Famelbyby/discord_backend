package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"fmt"
)

func (r *RelationsService) GetUserRelations(ctx context.Context, senderId string, targetIds []string) ([]models.UserRelation, error) {
	r.log.Info("[GetUserRelation] service started")
	var relations []models.UserRelation
	for i := 0; i < len(targetIds); i++ {
		relation, err := r.relations.GetUserRelation(ctx, senderId, targetIds[i])
		if err != nil {
			r.log.Error("[GetUserRelation] service error: " + err.Error())
			return nil, fmt.Errorf("[GetUserRelation] service error: " + err.Error())
		}
		relations = append(relations, relation)
	}

	return relations, nil
}
