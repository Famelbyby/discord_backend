package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	authv1 "idea-store-auth/gen/go/auth"
	"idea-store-auth/internal/utils"
	"log/slog"
	"net/http"
	"time"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	authAPi authv1.AuthClient
}

func (c *AuthClient) Login(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginResponse, err := c.authAPi.Login(context.Background(), request)
	if err != nil {
		if errors.Is(err, status.Error(codes.InvalidArgument, "Invalid credentials")) {
			utils.WriteError(w, "Invalid credentials")
			return
		}

		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	responseJson, err := json.Marshal(loginResponse)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthClient) Regsiter(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &authv1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	registerResponse, err := c.authAPi.Register(r.Context(), request)
	if err != nil {
		if errors.Is(err, status.Error(codes.AlreadyExists, "User already exists")) {
			utils.WriteError(w, "User already exists")
			return
		}

		slog.Error("client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	result, err := json.Marshal(registerResponse)
	if err != nil {
		slog.Error("client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
func NewAuthClient(addr string, timeout time.Duration, retriesCount int) (*AuthClient, error) {
	retryOptions := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		grpcretry.UnaryClientInterceptor(retryOptions...),
	))
	if err != nil {
		slog.Error("client Regsiter error: " + err.Error())
		return nil, fmt.Errorf("client Regsiter error: " + err.Error())
	}

	return &AuthClient{
		authAPi: authv1.NewAuthClient(cc),
	}, nil
}
