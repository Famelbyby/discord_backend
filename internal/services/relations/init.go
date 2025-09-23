package relations

import (
	"context"
	"discord_backend/internal/domain/models"
	"log/slog"
)

type RelationsService struct {
	log       *slog.Logger
	relations Relations
}

// DeclineFriendOffer implements relations.Relations.
func (r *RelationsService) DeclineFriendOffer(ctx context.Context, senderId string, revieverId string) error {
	panic("unimplemented")
}

type Relations interface {
	CreateNewRelation(ctx context.Context, userId string) error
	SendFriendOffer(ctx context.Context, senderId, recievedId string) (string, error)
	CancelFriendOffer(ctx context.Context, senderId, recieverId string) (string, error)
	AcceptFriendRequest(ctx context.Context, senderId, recievedId string) error
	DeclineFriendRequest(ctx context.Context, senderId, recievedId string) error
	RemoveFriend(ctx context.Context, senderId, toRemoveId string) error
	BlockUser(ctx context.Context, senderId, toBlockId string) error
	UnblockUser(ctx context.Context, senderId, toUnblockId string) error
	GetUserRelation(ctx context.Context, senderId, targetId string) (models.UserRelation, error)
	GetAllFriends(ctx context.Context, senderId string, page, limit int64) ([]string, error)
	GetAllIncomingOffers(ctx context.Context, senderId string) ([]string, error)
	GetAllOutgoingOffers(ctx context.Context, senderId string) ([]string, error)
	GetAllBlockedUsers(ctx context.Context, senderId string) ([]string, error)
}

func New(log *slog.Logger, relations Relations) *RelationsService {
	return &RelationsService{
		relations: relations,
		log:       log,
	}
}
