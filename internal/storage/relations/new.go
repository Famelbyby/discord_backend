package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RelationsStorage struct {
	collection *mongo.Collection
	log        *slog.Logger
}

// BlockUser implements relations.Relations.
func (s *RelationsStorage) BlockUser(ctx context.Context, senderId string, toBlockId string) error {
	panic("unimplemented")
}

// CancelFriendOffer implements relations.Relations.
func (s *RelationsStorage) CancelFriendOffer(ctx context.Context, senderId string, recieverId string) (string, error) {
	panic("unimplemented")
}

// DeclineFriendRequest implements relations.Relations.
func (s *RelationsStorage) DeclineFriendRequest(ctx context.Context, senderId string, recievedId string) error {
	panic("unimplemented")
}

// GetAllBlockedUsers implements relations.Relations.
func (s *RelationsStorage) GetAllBlockedUsers(ctx context.Context, senderId string) ([]string, error) {
	panic("unimplemented")
}

// GetAllFriends implements relations.Relations.
func (s *RelationsStorage) GetAllFriends(ctx context.Context, senderId string, page int64, limit int64) ([]string, error) {
	panic("unimplemented")
}

// GetAllIncomingOffers implements relations.Relations.
func (s *RelationsStorage) GetAllIncomingOffers(ctx context.Context, senderId string) ([]string, error) {
	panic("unimplemented")
}

// GetAllOutgoingOffers implements relations.Relations.
func (s *RelationsStorage) GetAllOutgoingOffers(ctx context.Context, senderId string) ([]string, error) {
	panic("unimplemented")
}

// GetUserRelation implements relations.Relations.
func (s *RelationsStorage) GetUserRelation(ctx context.Context, senderId string, targetId string) (models.UserRelation, error) {
	panic("unimplemented")
}

// RemoveFriend implements relations.Relations.
func (s *RelationsStorage) RemoveFriend(ctx context.Context, senderId string, toRemoveId string) error {
	panic("unimplemented")
}

// UnblockUser implements relations.Relations.
func (s *RelationsStorage) UnblockUser(ctx context.Context, senderId string, toUnblockId string) error {
	panic("unimplemented")
}

func NewStorage(db *mongo.Database, collectionName string, log *slog.Logger) *RelationsStorage {
	storage := &RelationsStorage{
		collection: db.Collection(collectionName),
		log:        log,
	}

	return storage
}
