package main

import (
	relationsv1 "discord_backend/gen/go/relations"
	"discord_backend/internal/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type RelationsClient struct {
	relationsAPi relationsv1.RelationsClient
}

func (c *RelationsClient) CreateNewRelation(w http.ResponseWriter, r *http.Request) {
	type createNewRelationRequest struct {
		UserId string `json:"userId"`
	}

	var req createNewRelationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("[CreateNewRelation] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &relationsv1.CreateNewRelationRequest{
		UserId: req.UserId,
	}

	_, err = c.relationsAPi.CreateNewRelation(r.Context(), request)
	if err != nil {
		slog.Error("[CreateNewRelation] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *RelationsClient) SendFriendOffer(w http.ResponseWriter, r *http.Request) {
	type sendFriendOfferRequest struct {
		FriendId string `json:"friendId"`
	}

	type sendFriendOfferResponse struct {
		UserId string `json:"id"`
	}

	var req sendFriendOfferRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("[SendFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[SendFriendOffer] senderId=" + senderId)

	request := &relationsv1.SendFriendOfferRequest{
		RecieverId: req.FriendId,
		SenderId:   senderId,
	}

	_, err = c.relationsAPi.SendFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[SendFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	resp := sendFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[SendFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)
}

func NewRelationsClient(addr string, timeout time.Duration, retriesCount int) (*RelationsClient, error) {
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

	return &RelationsClient{
		relationsAPi: relationsv1.NewRelationsClient(cc),
	}, nil
}
