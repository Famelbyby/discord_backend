package relations

import (
	"context"
	relationsv1 "discord_backend/gen/go/relations"
	"discord_backend/internal/domain/models"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Relations interface {
	CreateNewRelation(ctx context.Context, userId string) error
	SendFriendOffer(ctx context.Context, senderId, recievedId string) error
	CancelFriendOffer(ctx context.Context, senderId, recieverId string) error
	AcceptFriendOffer(ctx context.Context, senderId, recievedId string) error
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

type serverAPI struct {
	relationsv1.UnimplementedRelationsServer
	relations Relations
	log       *slog.Logger
}

func Register(gRPC *grpc.Server, relations Relations, log *slog.Logger) {
	relationsv1.RegisterRelationsServer(gRPC, &serverAPI{relations: relations, log: log})
}

func (s *serverAPI) CreateNewRelation(ctx context.Context, req *relationsv1.CreateNewRelationRequest) (*emptypb.Empty, error) {
	s.log.Info("[CreateNewRelation] grpc started")

	err := s.relations.CreateNewRelation(ctx, req.UserId)

	if err != nil {
		s.log.Error("[CreateNewRelation] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[CreateNewRelation] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) SendFriendOffer(ctx context.Context, req *relationsv1.SendFriendOfferRequest) (*emptypb.Empty, error) {
	s.log.Info("[SendFriendOffer] grpc started")

	err := s.relations.SendFriendOffer(ctx, req.SenderId, req.RecieverId)

	if err != nil {
		s.log.Error("[SendFriendOffer] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[SendFriendOffer] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) AcceptFriendOffer(ctx context.Context, req *relationsv1.AcceptFriendOfferRequest) (*emptypb.Empty, error) {
	s.log.Info("[AcceptFriendOffer] grpc started")

	err := s.relations.AcceptFriendOffer(ctx, req.SenderId, req.RecieverId)

	if err != nil {
		s.log.Error("[AcceptFriendOffer] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[AcceptFriendOffer] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) DeclineFriendOffer(ctx context.Context, req *relationsv1.DeclineFriendOfferRequest) (*emptypb.Empty, error) {
	s.log.Info("[DeclineFriendOffer] grpc started")

	err := s.relations.DeclineFriendOffer(ctx, req.SenderId, req.RecieverId)

	if err != nil {
		s.log.Error("[DeclineFriendOffer] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[DeclineFriendOffer] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) CancelFriendOffer(ctx context.Context, req *relationsv1.CancelFriendOfferRequest) (*emptypb.Empty, error) {
	s.log.Info("[CancelFriendOffer] grpc started")

	err := s.relations.CancelFriendOffer(ctx, req.SenderId, req.RecieverId)

	if err != nil {
		s.log.Error("[CancelFriendOffer] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[CancelFriendOffer] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) RemoveFriend(ctx context.Context, req *relationsv1.RemoveFriendRequest) (*emptypb.Empty, error) {
	s.log.Info("[RemoveFriend] grpc started")

	err := s.relations.RemoveFriend(ctx, req.SenderId, req.ToRemoveId)

	if err != nil {
		s.log.Error("[RemoveFriend] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[RemoveFriend] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) BlockUser(ctx context.Context, req *relationsv1.BlockUserRequest) (*emptypb.Empty, error) {
	s.log.Info("[BlockUser] grpc started")

	err := s.relations.BlockUser(ctx, req.SenderId, req.ToBlockId)

	if err != nil {
		s.log.Error("[BlockUser] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[BlockUser] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) UnblockUser(ctx context.Context, req *relationsv1.UnblockUserRequest) (*emptypb.Empty, error) {
	s.log.Info("[UnblockUser] grpc started")

	err := s.relations.UnblockUser(ctx, req.SenderId, req.ToUnblockId)

	if err != nil {
		s.log.Error("[UnblockUser] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[UnblockUser] grpc error: "+err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) GetUserRelation(ctx context.Context, req *relationsv1.GetUserRelationRequest) (*relationsv1.GetUserRelationResponse, error) {
	s.log.Info("[GetUserRelation] grpc started")

	relation, err := s.relations.GetUserRelation(ctx, req.SenderId, req.TargetId)

	if err != nil {
		s.log.Error("[GetUserRelation] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[GetUserRelation] grpc error: "+err.Error())
	}

	return &relationsv1.GetUserRelationResponse{
		IsFriend:   relation.IsFriend,
		IsOutgoing: relation.IsOutcoming,
		IsIncoming: relation.IsIncoming,
		IsBlocked:  relation.IsBlocked,
	}, nil
}

func (s *serverAPI) GetAllBlockedUsers(ctx context.Context, req *relationsv1.GetAllBlockedUsersRequest) (*relationsv1.GetAllBlockedUsersResponse, error) {
	s.log.Info("[GetAllBlockedUsers] grpc started")

	users, err := s.relations.GetAllBlockedUsers(ctx, req.SenderId, req.Page, req.Limit)

	if err != nil {
		s.log.Error("[GetAllBlockedUsers] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[GetAllBlockedUsers] grpc error: "+err.Error())
	}

	return &relationsv1.GetAllBlockedUsersResponse{
		Ids: users,
	}, nil
}

func (s *serverAPI) GetAllFriends(ctx context.Context, req *relationsv1.GetAllFriendsRequest) (*relationsv1.GetAllFriendsResponse, error) {
	s.log.Info("[GetAllFriends] grpc started")

	users, err := s.relations.GetAllFriends(ctx, req.SenderId, req.Page, req.Limit)

	if err != nil {
		s.log.Error("[GetAllFriends] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[GetAllFriends] grpc error: "+err.Error())
	}

	return &relationsv1.GetAllFriendsResponse{
		Friends: users,
	}, nil
}

func (s *serverAPI) GetAllOutgoingOffers(ctx context.Context, req *relationsv1.GetAllOutgoingOffersRequest) (*relationsv1.GetAllOutgoingOffersResponse, error) {
	s.log.Info("[GetAllOutgoingOffers] grpc started")

	users, err := s.relations.GetAllOutgoingOffers(ctx, req.SenderId, req.Page, req.Limit)

	if err != nil {
		s.log.Error("[GetAllOutgoingOffers] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[GetAllOutgoingOffers] grpc error: "+err.Error())
	}

	return &relationsv1.GetAllOutgoingOffersResponse{
		Ids: users,
	}, nil
}

func (s *serverAPI) GetAllIncomingOffers(ctx context.Context, req *relationsv1.GetAllIncomingOffersRequest) (*relationsv1.GetAllIncomingOffersResponse, error) {
	s.log.Info("[GetAllIncomingOffers] grpc started")

	users, err := s.relations.GetAllIncomingOffers(ctx, req.SenderId, req.Page, req.Limit)

	if err != nil {
		s.log.Error("[GetAllIncomingOffers] grpc error: " + err.Error())
		return nil, fmt.Errorf("%s", "[GetAllIncomingOffers] grpc error: "+err.Error())
	}

	return &relationsv1.GetAllIncomingOffersResponse{
		Ids: users,
	}, nil
}
