package relations

import (
	"context"
	"discord_backend/internal/storage"
	"errors"
	"fmt"
)

func (s *RelationsService) CreateNewRelation(ctx context.Context, userId string) error {
	s.log.Info("[CreateNewRelation] service started")

	_, err := s.relations.GetRelationRecord(ctx, userId)

	if err == nil {
		return storage.ErrRecordExists
	} else {
		if !errors.Is(err, storage.ErrUserNotFound) {
			s.log.Error("[CreateNewRelation] service error: " + err.Error())
			return fmt.Errorf("[CreateNewRelation] service error: " + err.Error())
		}
	}

	err = s.relations.CreateNewRelation(ctx, userId)
	if err != nil {
		s.log.Error("[CreateNewRelation] service error: " + err.Error())
		return fmt.Errorf("[CreateNewRelation] service error: " + err.Error())
	}

	return nil
}
