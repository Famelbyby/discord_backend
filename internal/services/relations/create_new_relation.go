package relations

import (
	"context"
	"fmt"
)

func (s *RelationsService) CreateNewRelation(ctx context.Context, userId string) error {
	s.log.Info("[CreateNewRelation] service started")

	//проверка на существование

	err := s.relations.CreateNewRelation(ctx, userId)
	if err != nil {
		s.log.Error("[CreateNewRelation] service error: " + err.Error())
		return fmt.Errorf("[CreateNewRelation] service error: " + err.Error())
	}

	return nil
}
