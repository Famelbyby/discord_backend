package relations

import (
	"context"
	relationsv1 "discord_backend/gen/go/relations"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Relations interface {
	CreateNewRelation(ctx context.Context, userId string) error
	SendFriendOffer(ctx context.Context, senderId, recieverId string) error
	AcceptFriendOffer(ctx context.Context, senderId, recieverId string) error
	DeclineFriendOffer(ctx context.Context, senderId, revieverId string) error
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
