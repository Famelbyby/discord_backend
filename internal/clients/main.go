package main

import (
	"fmt"
	common "idea-store-auth/cmd"
	"idea-store-auth/internal/config"
	"idea-store-auth/internal/middlewares"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var (
	AppPortEnv      = "APP_PORT"
	AppSecretEnv    = "APP_SECRET"
	TokenCookieName = "token"
	EmptyValue      = int64(-1)
)

func main() {

	cfg := config.MustLoad()
	authClient, _ := NewAuthClient(common.GrpcAuthAddress(cfg), cfg.Clients.Auth.Timeout, cfg.Clients.Auth.RetriesCount)
	router := mux.NewRouter()
	router.HandleFunc("/api/register", authClient.Regsiter).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/login", authClient.Login).Methods(http.MethodPost, http.MethodOptions)

	/*router.HandleFunc("/app/files/{name}", GetImages).Methods(http.MethodGet)

	router.HandleFunc("/api/idea", ideasClient.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/idea", ideasClient.GetIdea).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/ideas", ideasClient.GetAllIdeas).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/ideas", ideasClient.GetIdeas).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/search-ideas", ideasClient.GetIdeasFromSearch).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/board", boardsClient.CreateBoard).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/board", boardsClient.GetBoard).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/board", boardsClient.DeleteBoard).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/my-boards", boardsClient.GetCurrentUsersBoards).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/boards", boardsClient.GetBoards).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/ideas-in-board", boardsClient.GetIdeasInBoard).Methods(http.MethodGet, http.MethodOptions)


	router.HandleFunc("/api/profile", profilesClient.CreateProfile).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/profile", profilesClient.UpdateProfile).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/profile", profilesClient.GetProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/my-profile", profilesClient.GetCurrentProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/toggle-save-idea", profilesClient.ToggleSaveIdea).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/is-idea-saved", profilesClient.IsIdeaSaved).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/get-saved-ideas", profilesClient.GetSavedIdeas).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/search-profiles", profilesClient.GetProfilesFromSearch).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/toggle-like-idea", profilesClient.ToggleLikeIdea).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/is-idea-liked", profilesClient.IsIdeaLiked).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/comments", commentsClient.GetComments).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/comment", commentsClient.CreateComment).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/chats", chatsClient.GetChats).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/chat_ws", chatsClient.HandleChatWebSocket).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/messages", chatsClient.GetMessages).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/message", chatsClient.SendMessage).Methods(http.MethodPost, http.MethodOptions)*/
	handler := middlewares.CorsMiddleware(router)
	fmt.Println("Server is listening...")

	log.Fatal(http.ListenAndServe(os.Getenv(AppPortEnv), handler))
}

func GetUserIdByRequestWithCookie(r *http.Request) (int64, error) {
	tokenCookie, err := r.Cookie(TokenCookieName)
	if err != nil {
		slog.Error("GetUserIdByRequestWithCookie error: " + err.Error())
		return EmptyValue, fmt.Errorf("GetUserIdByRequestWithCookie error: " + err.Error())
	}

	claims := jwt.MapClaims{}
	tokenStr := tokenCookie.String()[6:]
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(AppSecretEnv)), nil
	})
	if err != nil {
		slog.Error("GetUserIdByRequestWithCookie error: " + err.Error())
		return EmptyValue, fmt.Errorf("GetUserIdByRequestWithCookie error: " + err.Error())
	}

	userId := claims["uid"].(float64)
	userIdStr := fmt.Sprint(userId)
	userIdInt, _ := strconv.ParseInt(userIdStr, 10, 64)
	return userIdInt, nil
}
