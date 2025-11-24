package main

import (
	"bytes"
	"discord_backend/internal/utils"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sort"
)

const (
	chatsUrl = "http://chats_py:9091"
)

type ChatUserProfile struct {
	ID        string   `json:"id"`
	ShortLink string   `json:"short_link"`
	Mail      string   `json:"mail"`
	Username  string   `json:"username"`
	CreatedAt int64    `json:"created_at"`
	AvatarURL string   `json:"avatar_url"`
	Status    string   `json:"status"`
	Rules     []string `json:"rules"`
}

type ChatWithProfiles struct {
	Id     string            `json:"id,omitempty"`
	Name   string            `json:"name"`
	LeadId string            `json:"lead_id"`
	Users  []ChatUserProfile `json:"users"`
}

func GetProfilesOfChatUsers(chatUsers []ChatUser) []ChatUserProfile {
	sort.Slice(chatUsers, func(i, j int) bool { return chatUsers[i].Id < chatUsers[j].Id })

	var usersIds []string

	for _, user := range chatUsers {
		usersIds = append(usersIds, user.Id)
	}

	usersProfiles := GetProfiles(usersIds)

	var chatUsersProfiles []ChatUserProfile

	for index, userProfile := range usersProfiles {
		chatUsersProfiles = append(chatUsersProfiles, ChatUserProfile{
			ID:        userProfile.ID,
			Mail:      userProfile.Mail,
			ShortLink: userProfile.ShortLink,
			Username:  userProfile.Username,
			CreatedAt: userProfile.CreatedAt,
			AvatarURL: userProfile.AvatarURL,
			Status:    userProfile.Status,
			Rules:     chatUsers[index].Rules,
		})
	}

	return chatUsersProfiles
}

func HandleGetUserChats(w http.ResponseWriter, r *http.Request) {
	url := chatsUrl + r.URL.RequestURI()

	resp, err := http.Get(url)

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleGetChatById(w http.ResponseWriter, r *http.Request) {
	url := chatsUrl + r.URL.RequestURI()

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("client HandleGetProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleGetProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	type GetChatResponse struct {
		Chat Chat `json:"chat"`
	}

	var chatResponse GetChatResponse

	err = json.Unmarshal(body, &chatResponse)

	if err != nil {
		slog.Error("client HandleGetProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	chat := ChatWithProfiles{
		Id:     chatResponse.Chat.Id,
		Name:   chatResponse.Chat.Name,
		LeadId: chatResponse.Chat.LeadId,
		Users:  GetProfilesOfChatUsers(chatResponse.Chat.Users),
	}

	jsonChat, err := json.Marshal(chat)

	if err != nil {
		slog.Error("client HandleGetProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(jsonChat)
}

func HandleDeleteChatById(w http.ResponseWriter, r *http.Request) {
	url := chatsUrl + r.URL.RequestURI()

	req, err := http.NewRequest(http.MethodDelete, url, r.Body)
	if err != nil {
		slog.Error("client HandleDeleteProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		slog.Error("client HandleDeleteProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleDeleteProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleDeleteUserFromChat(w http.ResponseWriter, r *http.Request) {
	url := chatsUrl + r.URL.RequestURI()

	req, err := http.NewRequest(http.MethodDelete, url, r.Body)
	if err != nil {
		slog.Error("client HandleDeleteProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleDeleteProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

type ChatUser struct {
	Id    string   `json:"id"`
	Rules []string `json:"rules"`
}

type Chat struct {
	Id     string     `json:"id,omitempty"`
	Name   string     `json:"name"`
	LeadId string     `json:"lead_id"`
	Users  []ChatUser `json:"users"`
}

type UpdateChat struct {
	Name  string     `json:"name,omitempty"`
	Users []ChatUser `json:"users,omitempty"`
}

func HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	bytes1 := &bytes.Buffer{}

	copiedBody := io.TeeReader(r.Body, bytes1)

	reqBody, err := io.ReadAll(copiedBody)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	var chat Chat

	err = json.Unmarshal(reqBody, &chat)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	var usersIds []string

	usersIds = append(usersIds, chat.LeadId)

	for _, user := range chat.Users {
		usersIds = append(usersIds, user.Id)
	}

	usersProfiles := GetProfiles(usersIds)

	if len(usersIds) != len(usersProfiles) {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteError(w, "users don't exist")
		return
	}

	url := chatsUrl + r.URL.RequestURI()

	req, err := http.NewRequest(http.MethodPost, url, bytes1)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleUpdateChat(w http.ResponseWriter, r *http.Request) {
	bytes1 := &bytes.Buffer{}

	copiedBody := io.TeeReader(r.Body, bytes1)
	reqBody, err := io.ReadAll(copiedBody)

	defer r.Body.Close()

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	var chat UpdateChat

	err = json.Unmarshal(reqBody, &chat)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	var usersIds []string

	for _, user := range chat.Users {
		usersIds = append(usersIds, user.Id)
	}

	usersProfiles := GetProfiles(usersIds)

	if len(usersIds) != len(usersProfiles) {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteError(w, "users don't exist")
		return
	}

	url := chatsUrl + r.URL.RequestURI()

	req, err := http.NewRequest(http.MethodPut, url, bytes1)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
