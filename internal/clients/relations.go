package main

import (
	relationsv1 "discord_backend/gen/go/relations"
	"discord_backend/internal/storage"
	"discord_backend/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const recordsOnPage = 30

type RelationsClient struct {
	relationsAPi relationsv1.RelationsClient
}

type Profile struct {
	ID        string `json:"id"`
	ShortLink string `json:"short_link"`
	Mail      string `json:"mail"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	AvatarURL string `json:"avatar_url"`
	Status    string `json:"status"`
}

func GetProfileById(id string) (Profile, error) {
	type ProfileClientResponse struct {
		Profile Profile `json:"profile"`
	}

	resp, err := http.Get("http://profile_py:9999/api/profile/" + id)

	if err != nil {
		slog.Error("client profile error: " + err.Error())
		return Profile{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		return Profile{}, err
	}

	var profileResponse ProfileClientResponse

	if err := json.Unmarshal(body, &profileResponse); err != nil {
		slog.Error("client Login error: " + err.Error())

		return Profile{}, err
	}

	return profileResponse.Profile, nil
}

func GetProfiles(profileIds []string) []Profile {
	var profiles = []Profile{}

	for _, id := range profileIds {
		profile, err := GetProfileById(id)

		if err != nil {
			slog.Error("client profile error: " + err.Error())
		} else {
			profiles = append(profiles, profile)
		}
	}

	return profiles
}

func (c *RelationsClient) CreateNewRelation(w http.ResponseWriter, r *http.Request) {
	type createNewRelationRequest struct {
		UserId string `json:"id"`
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
		if strings.Contains(err.Error(), storage.ErrRecordExists.Error()) {
			utils.WriteError(w, "Record already exists")
			return
		}
		slog.Error("[CreateNewRelation] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
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
		outputError := "Internal error"
		if strings.Contains(err.Error(), storage.ErrUserNotFound.Error()) {
			outputError = storage.ErrUserNotFound.Error()
		}
		if strings.Contains(err.Error(), storage.ErrFriendOfferAlreadySent.Error()) {
			outputError = storage.ErrFriendOfferAlreadySent.Error()
		}
		if strings.Contains(err.Error(), storage.ErrCantAddBlockedUser.Error()) {
			outputError = storage.ErrCantAddBlockedUser.Error()
		}
		if strings.Contains(err.Error(), storage.ErrCantAddAlreadyFriend.Error()) {
			outputError = storage.ErrCantAddAlreadyFriend.Error()
		}
		utils.WriteError(w, outputError)
		return
	}

	resp := sendFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[SendFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
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
	recieverId := vars["id"]
	slog.Info("[AcceptFriendOffer] recieverId=" + recieverId)

	request := &relationsv1.AcceptFriendOfferRequest{
		RecieverId: recieverId,
		SenderId:   req.FriendId,
	}

	_, err = c.relationsAPi.AcceptFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[AcceptFriendOffer] client error: " + err.Error())
		outputError := "Internal error"
		if strings.Contains(err.Error(), storage.ErrCantAddAlreadyFriend.Error()) {
			outputError = storage.ErrCantAddAlreadyFriend.Error()
		}
		if strings.Contains(err.Error(), storage.ErrUserNotFound.Error()) {
			outputError = storage.ErrUserNotFound.Error()
		}
		if strings.Contains(err.Error(), storage.ErrNoOutgoingOffer.Error()) {
			outputError = storage.ErrNoOutgoingOffer.Error()
		}
		utils.WriteError(w, outputError)
		return
	}

	resp := acceptFriendOfferResponse{UserId: recieverId}
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
	recieverId := vars["id"]
	slog.Info("[DeclineFriendOffer] recieverId=" + recieverId)

	request := &relationsv1.DeclineFriendOfferRequest{
		RecieverId: recieverId,
		SenderId:   req.FriendId,
	}

	_, err = c.relationsAPi.DeclineFriendOffer(r.Context(), request)
	if err != nil {
		slog.Error("[DeclineFriendOffer] client error: " + err.Error())
		if strings.Contains(err.Error(), storage.ErrNoOutgoingOffer.Error()) {
			utils.WriteError(w, storage.ErrNoOutgoingOffer.Error())
			return
		}
		utils.WriteError(w, "Internal error")
		return
	}

	resp := declineFriendOfferResponse{UserId: recieverId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[DeclineFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
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
		utils.WriteError(w, "no friendId parameter given")
		return
	}

	request := &relationsv1.CancelFriendOfferRequest{
		RecieverId: friendId,
		SenderId:   senderId,
	}

	_, err := c.relationsAPi.CancelFriendOffer(r.Context(), request)
	if err != nil {
		if strings.Contains(err.Error(), storage.ErrNoOutgoingOffer.Error()) {
			utils.WriteError(w, storage.ErrNoOutgoingOffer.Error())
			return
		}
		slog.Error("[CancelFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	resp := cancelFriendOfferResponse{UserId: senderId}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[CancelFriendOffer] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)
}

func (c *RelationsClient) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[RemoveFriend] senderId=" + senderId)

	friendId := params.Get("friendId")
	if friendId == "" {
		utils.WriteError(w, "no friendId parameter given")
		return
	}

	request := &relationsv1.RemoveFriendRequest{
		ToRemoveId: friendId,
		SenderId:   senderId,
	}

	_, err := c.relationsAPi.RemoveFriend(r.Context(), request)
	if err != nil {
		slog.Error("[RemoveFriend] client error: " + err.Error())
		if strings.Contains(err.Error(), storage.ErrUserIsntFriend.Error()) {
			utils.WriteError(w, storage.ErrUserIsntFriend.Error())
			return
		}
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *RelationsClient) BlockUser(w http.ResponseWriter, r *http.Request) {
	type blockUserRequest struct {
		ProfileId string `json:"profileId"`
	}

	var req blockUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("[BlockUser] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[BlockUser] senderId=" + senderId)

	request := &relationsv1.BlockUserRequest{
		ToBlockId: req.ProfileId,
		SenderId:  senderId,
	}

	_, err = c.relationsAPi.BlockUser(r.Context(), request)
	if err != nil {
		slog.Error("[BlockUser] client error: " + err.Error())
		outputError := "Internal error"
		if strings.Contains(err.Error(), storage.ErrUserAlreadyBlocked.Error()) {
			outputError = storage.ErrUserAlreadyBlocked.Error()
		}
		utils.WriteError(w, outputError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *RelationsClient) UnblockUser(w http.ResponseWriter, r *http.Request) {
	type unblockUserResponse struct {
		Error string `json:"error,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[UnblockUser] senderId=" + senderId)

	profileId := params.Get("profileId")
	if profileId == "" {
		resp := unblockUserResponse{Error: "no profileId given"}
		respJson, err := json.Marshal(resp)
		if err != nil {
			slog.Error("[UnblockUser] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
	}

	request := &relationsv1.UnblockUserRequest{
		ToUnblockId: profileId,
		SenderId:    senderId,
	}

	_, err := c.relationsAPi.UnblockUser(r.Context(), request)
	if err != nil {
		slog.Error("[UnblockUser] client error: " + err.Error())
		outputError := "Internal error"
		if strings.Contains(err.Error(), storage.ErrUserAlreadyUnblocked.Error()) {
			outputError = storage.ErrUserAlreadyUnblocked.Error()
		}
		utils.WriteError(w, outputError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *RelationsClient) GetRelation(w http.ResponseWriter, r *http.Request) {
	type getRelationResponse struct {
		Error       string `json:"error,omitempty"`
		IsFriend    bool   `json:"isFriend,omitempty"`
		IsIncoming  bool   `json:"isIncoming,omitempty"`
		IsOutcoming bool   `json:"isOutcoming,omitempty"`
		IsBlocked   bool   `json:"isBlocked,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[GetRelation] senderId=" + senderId)

	profileId := params.Get("profileId")
	if profileId == "" {
		resp := getRelationResponse{Error: "no profileId given"}
		respJson, err := json.Marshal(resp)
		if err != nil {
			slog.Error("[GetRelation] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
	}

	request := &relationsv1.GetUserRelationRequest{
		SenderId: senderId,
		TargetId: profileId,
	}

	result, err := c.relationsAPi.GetUserRelation(r.Context(), request)
	if err != nil {
		slog.Error("[GetRelation] client error: " + err.Error())
		outputError := "Internal error"
		utils.WriteError(w, outputError)
		return
	}

	resp := getRelationResponse{IsFriend: result.IsFriend,
		IsIncoming:  result.IsIncoming,
		IsOutcoming: result.IsOutgoing,
		IsBlocked:   result.IsBlocked,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[GetRelation] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (c *RelationsClient) GetAllFriends(w http.ResponseWriter, r *http.Request) {
	type getAllFriendsResponse struct {
		Error   string    `json:"error,omitempty"`
		Friends []Profile `json:"friends,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[GetAllFriends] senderId=" + senderId)

	var page int64
	page = 1
	var limit int64
	limit = recordsOnPage
	var err error
	if params.Has("page") {
		page, err = strconv.ParseInt(params.Get("page"), 10, 64)
		if err != nil {
			slog.Error("[GetAllFriends] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}
	if params.Has("limit") {
		limit, err = strconv.ParseInt(params.Get("limit"), 10, 64)
		if err != nil {
			slog.Error("[GetAllFriends] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}

	request := &relationsv1.GetAllFriendsRequest{
		SenderId: senderId,
		Page:     page,
		Limit:    limit,
	}

	result, err := c.relationsAPi.GetAllFriends(r.Context(), request)
	if err != nil {
		slog.Error("[GetAllFriends] client error: " + err.Error())
		outputError := "Internal error"
		utils.WriteError(w, outputError)
		return
	}

	profiles := GetProfiles(result.Friends)

	resp := getAllFriendsResponse{
		Friends: profiles,
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[GetAllFriends] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (c *RelationsClient) GetAllIncomingOffers(w http.ResponseWriter, r *http.Request) {
	type getAllIncomingssResponse struct {
		Error     string    `json:"error,omitempty"`
		Incomings []Profile `json:"incomings,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[GetAllIncomingOffers] senderId=" + senderId)

	var page int64
	page = 1
	var limit int64
	limit = recordsOnPage
	var err error
	if params.Has("page") {
		page, err = strconv.ParseInt(params.Get("page"), 10, 64)
		if err != nil {
			slog.Error("[GetAllIncomingOffers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}
	if params.Has("limit") {
		limit, err = strconv.ParseInt(params.Get("limit"), 10, 64)
		if err != nil {
			slog.Error("[GetAllIncomingOffers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}

	request := &relationsv1.GetAllIncomingOffersRequest{
		SenderId: senderId,
		Page:     page,
		Limit:    limit,
	}

	result, err := c.relationsAPi.GetAllIncomingOffers(r.Context(), request)
	if err != nil {
		slog.Error("[GetAllIncomingOffers] client error: " + err.Error())
		outputError := "Internal error"
		utils.WriteError(w, outputError)
		return
	}

	profiles := GetProfiles(result.Ids)

	resp := getAllIncomingssResponse{
		Incomings: profiles,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[GetAllIncomingOffers] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (c *RelationsClient) GetAllOutgoingOffers(w http.ResponseWriter, r *http.Request) {
	type getAllOutgoingsResponse struct {
		Error      string    `json:"error,omitempty"`
		Outcomings []Profile `json:"outcomings,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[GetAllOutgoingOffers] senderId=" + senderId)

	var page int64
	page = 1
	var limit int64
	limit = recordsOnPage
	var err error
	if params.Has("page") {
		page, err = strconv.ParseInt(params.Get("page"), 10, 64)
		if err != nil {
			slog.Error("[GetAllOutgoingOffers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}
	if params.Has("limit") {
		limit, err = strconv.ParseInt(params.Get("limit"), 10, 64)
		if err != nil {
			slog.Error("[GetAllOutgoingOffers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}

	request := &relationsv1.GetAllOutgoingOffersRequest{
		SenderId: senderId,
		Page:     page,
		Limit:    limit,
	}

	result, err := c.relationsAPi.GetAllOutgoingOffers(r.Context(), request)
	if err != nil {
		slog.Error("[GetAllOutgoingOffers] client error: " + err.Error())
		outputError := "Internal error"
		utils.WriteError(w, outputError)
		return
	}

	profiles := GetProfiles(result.Ids)

	resp := getAllOutgoingsResponse{
		Outcomings: profiles,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[GetAllOutgoingOffers] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (c *RelationsClient) GetAllBlockedUsers(w http.ResponseWriter, r *http.Request) {
	type getAlBlocksResponse struct {
		Error  string    `json:"error,omitempty"`
		Blocks []Profile `json:"blocks,omitempty"`
	}

	myUrl, _ := url.Parse(r.RequestURI)
	params, _ := url.ParseQuery(myUrl.RawQuery)

	vars := mux.Vars(r)
	senderId := vars["id"]
	slog.Info("[GetAllBlockedUsers] senderId=" + senderId)

	var page int64
	page = 1
	var limit int64
	limit = recordsOnPage
	var err error
	if params.Has("page") {
		page, err = strconv.ParseInt(params.Get("page"), 10, 64)
		if err != nil {
			slog.Error("[GetAllBlockedUsers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}
	if params.Has("limit") {
		limit, err = strconv.ParseInt(params.Get("limit"), 10, 64)
		if err != nil {
			slog.Error("[GetAllBlockedUsers] client error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}

	request := &relationsv1.GetAllBlockedUsersRequest{
		SenderId: senderId,
		Page:     page,
		Limit:    limit,
	}

	result, err := c.relationsAPi.GetAllBlockedUsers(r.Context(), request)
	if err != nil {
		slog.Error("[GetAllBlockedUsers] client error: " + err.Error())
		outputError := "Internal error"
		utils.WriteError(w, outputError)
		return
	}

	profiles := GetProfiles(result.Ids)

	resp := getAlBlocksResponse{
		Blocks: profiles,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		slog.Error("[GetAllBlockedUsers] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
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
