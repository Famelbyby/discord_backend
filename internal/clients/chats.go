package main

import (
	"bytes"
	"context"
	"discord_backend/internal/middlewares"
	"discord_backend/internal/utils"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
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

type MType string

const (
	TypeSendMessage   MType = "SEND_MESSAGE"
	TypeDeleteMessage MType = "DELETE_MESSAGE"
	TypeEditMessage   MType = "EDIT_MESSAGE"
)

type Payload struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId,omitempty"` // omitempty, так как для SEND_MESSAGE может быть пустым
}

type WsRequest struct {
	Type    MType   `json:"type"`
	ChatID  string  `json:"chatId"`
	Payload Payload `json:"payload"`
}

type KafkaMessage struct {
	UserID string    `json:"userId"` // ID пользователя
	Body   WsRequest `json:"body"`   // Само сообщение
}

type ConnectionManager struct {
	clients map[string]*websocket.Conn // UserID -> Conn
	mu      sync.RWMutex               // Мьютекс для безопасного доступа из разных горутин
}

func (cm *ConnectionManager) Add(userID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[userID] = conn
}

func (cm *ConnectionManager) Remove(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.clients, userID)
}

func (cm *ConnectionManager) Get(userID string) (*websocket.Conn, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, ok := cm.clients[userID]
	return conn, ok
}

var (
	hub = ConnectionManager{clients: make(map[string]*websocket.Conn)}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Важно: для локальной разработки разрешаем все источники (CORS)
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ChatWsHandler(w http.ResponseWriter, r *http.Request) {
	// Достаем UserID из контекста (положенного middleware)
	userID := r.Context().Value(middlewares.UserKey).(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Регистрируем соединение
	hub.Add(userID, conn)
	log.Printf("User %s connected", userID)

	// При закрытии сокета - удаляем из мапы
	defer func() {
		hub.Remove(userID)
		conn.Close()
		log.Printf("User %s disconnected", userID)
	}()

	// Цикл чтения сообщений от фронта
	for {
		var msg WsRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Обработка входящего сообщения: отправляем в Kafka chat-in
		processIncomingMessage(userID, msg)
	}
}

func processIncomingMessage(userID string, msg WsRequest) {
	// Оборачиваем сообщение, добавляя ID отправителя
	kafkaMsg := KafkaMessage{
		UserID: userID,
		Body:   msg,
	}

	jsonBytes, err := json.Marshal(kafkaMsg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	// Пишем в Kafka (асинхронно или синхронно - зависит от требований)
	// Используем context с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(userID), // Key нужен для партицирования (чтобы сообщения юзера шли по порядку)
		Value: jsonBytes,
	})

	if err != nil {
		log.Println("Error writing to Kafka chat-in:", err)
	} else {
		log.Printf("Sent to Kafka chat-in: %s | %s", msg.Type, msg.Payload.Message)
	}
}

func GetProfilesOfChatUsers(chatUsers []ChatUser) []ChatUserProfile {
	sort.Slice(chatUsers, func(i, j int) bool { return chatUsers[i].Id < chatUsers[j].Id })

	var usersIds []string

	for _, user := range chatUsers {
		usersIds = append(usersIds, user.Id)
	}

	usersProfiles := profilesClient.GetProfiles(usersIds)

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
	err := utils.CompareUserIDsInQuery(r, middlewares.UserKey, "user_id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

	reqUrl := chatsUrl + r.URL.RequestURI()

	resp, err := http.Get(reqUrl)

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
	err := utils.CompareUserIDsInQuery(r, middlewares.UserKey, "user_id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

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
	err := utils.CompareUserIDsInQuery(r, middlewares.UserKey, "lead_id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

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

	usersProfiles := profilesClient.GetProfiles(usersIds)

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
	err := utils.CompareUserIDsInQuery(r, middlewares.UserKey, "user_id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

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

	usersProfiles := profilesClient.GetProfiles(usersIds)

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
