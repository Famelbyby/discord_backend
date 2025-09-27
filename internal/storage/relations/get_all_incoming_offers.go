package relations

import (
	"context"
	"discord_backend/internal/storage"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *RelationsStorage) GetAllIncomingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	s.log.Info("[GetAllIncomingOffers] storage started")

	filter := bson.M{"user_id": senderId}

	var senderRelation dtoRelations

	err := s.collection.FindOne(ctx, filter).Decode(&senderRelation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, storage.ErrUserNotFound
		}

		slog.Error("[GetAllIncomingOffers] storage error: " + err.Error())
		return nil, errors.New("[GetAllIncomingOffers] storage error: " + err.Error())
	}

	return senderRelation.IncomingIds[(page-1)*recordsOnPage : (page-1)*recordsOnPage+limit], nil
}
