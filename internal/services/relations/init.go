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

type Relations interface {
	CreateNewRelation(ctx context.Context, userId string) error
	SendFriendOffer(ctx context.Context, senderId, recievedId string) error
	CancelFriendOffer(ctx context.Context, senderId, recieverId string) error
	AcceptFriendRequest(ctx context.Context, senderId, recievedId string) error
	DeclineFriendOffer(ctx context.Context, senderId, recievedId string) error
	RemoveFriend(ctx context.Context, senderId, toRemoveId string) error
	BlockUser(ctx context.Context, senderId, toBlockId string) error
	UnblockUser(ctx context.Context, senderId, toUnblockId string) error
	GetUserRelation(ctx context.Context, senderId, targetId string) (models.UserRelation, error)
	GetAllFriends(ctx context.Context, senderId string, page, limit int64) ([]string, error)
	GetAllIncomingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error)
	GetAllOutgoingOffers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error)
	GetAllBlockedUsers(ctx context.Context, senderId string, page int64, limit int64) ([]string, error)
}

func New(log *slog.Logger, relations Relations) *RelationsService {
	return &RelationsService{
		relations: relations,
		log:       log,
	}
}
