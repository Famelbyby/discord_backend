package main

import (
	relationsv1 "discord_backend/gen/go/relations"
	"discord_backend/internal/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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
		UserId string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
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
		UserId string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
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

func (c *RelationsClient) AcceptFriendOffer(w http.ResponseWriter, r *http.Request) {
	type acceptFriendOfferRequest struct {
		FriendId string `json:"friendId"`
	}

	type acceptFriendOfferResponse struct {
		UserId string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var req acceptFriendOfferRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("[AcceptFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[AcceptFriendOffer] senderId=" + senderId)

	request := &relationsv1.AcceptFriendOfferRequest{
		RecieverId: req.FriendId,
		SenderId:   senderId,
	}

	_, err = c.relationsAPi.AcceptFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[AcceptFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	resp := acceptFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[AcceptFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)
}

func (c *RelationsClient) DeclineFriendOffer(w http.ResponseWriter, r *http.Request) {
	type declineFriendOfferRequest struct {
		FriendId string `json:"friendId"`
	}

	type declineFriendOfferResponse struct {
		UserId string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var req declineFriendOfferRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("[DeclineFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[DeclineFriendOffer] senderId=" + senderId)

	request := &relationsv1.DeclineFriendOfferRequest{
		RecieverId: req.FriendId,
		SenderId:   senderId,
	}

	_, err = c.relationsAPi.DeclineFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[DeclineFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	resp := declineFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[DeclineFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)
}

func (c *RelationsClient) CancelFriendOffer(w http.ResponseWriter, r *http.Request) {

	type cancelFriendOfferResponse struct {
		UserId string `json:"id,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[CancelFriendOffer] senderId=" + senderId)

	friendId := params.Get("friendId")
	if friendId == "" {
		resp := cancelFriendOfferResponse{Error: "no friendId given"}
		respJson, err := json.Marshal(resp)
		if err != nil {
			slog.Error("[CancelFriendOffer] client error: " + err.Error())
			utils.WriteError(w, "Internal error: "+err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(respJson)
	}

	request := &relationsv1.CancelFriendOfferRequest{
		RecieverId: friendId,
		SenderId:   senderId,
	}

	_, err := c.relationsAPi.CancelFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[CancelFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error: "+err.Error())
		return
	}

	resp := cancelFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[CancelFriendOffer] client error: " + err.Error())
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
