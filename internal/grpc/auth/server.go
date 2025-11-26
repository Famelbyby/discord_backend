package auth

import (
	"context"
	authv1 "discord_backend/gen/go/auth"
	"discord_backend/internal/services/auth"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, id string) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string, id string) (sessionId string, err error)
	Logout(ctx context.Context, sessionId string) error
	IsRegistered(ctx context.Context, sessionId string) (userId string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	slog.Info("started to login")
	if err := validateLogin(req); err != nil {
		slog.Error("grpc Login error: " + err.Error())
		return nil, fmt.Errorf("grpc Login error: " + err.Error())
	}

	sessionId, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.Id)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}

		slog.Error("grpc Login error: " + err.Error())
		return nil, fmt.Errorf("grpc Login error: " + err.Error())
	}

	resp := &authv1.LoginResponse{SessionId: sessionId}
	return resp, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	id := req.GetId()

	if err := validateRegister(req); err != nil {
		slog.Error("grpc Register error: " + err.Error())
		return nil, fmt.Errorf("grpc Register error: " + err.Error())
	}

	sessionId, err := s.auth.RegisterNewUser(ctx, email, password, id)

	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		slog.Error("grpc Register error: " + err.Error())
		return nil, fmt.Errorf("grpc Register error: " + err.Error())
	}
	return &authv1.RegisterResponse{
		SessionId: sessionId,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	slog.Info("started to logout")

	err := s.auth.Logout(ctx, req.SessionId)
	if err != nil {
		slog.Error("grpc Login error: " + err.Error())
		return &authv1.LogoutResponse{}, fmt.Errorf("grpc Login error: " + err.Error())
	}

	return &authv1.LogoutResponse{}, nil
}

func (s *serverAPI) IsRegistered(ctx context.Context, req *authv1.IsRegisteredRequest) (*authv1.IsRegisteredResponse, error) {
	slog.Info("started is registered")

	userId, err := s.auth.IsRegistered(ctx, req.SessionId)

	if err != nil {
		slog.Error("grpc Login error: " + err.Error())
		return &authv1.IsRegisteredResponse{}, fmt.Errorf("grpc Login error: " + err.Error())
	}

	return &authv1.IsRegisteredResponse{
		UserId: userId,
	}, nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email must not be empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password must not be empty")
	}
	return nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email must not be empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password must not be empty")
	}
	return nil
}
