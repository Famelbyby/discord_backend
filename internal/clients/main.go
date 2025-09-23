package main

import (
	common "discord_backend/cmd"
	"discord_backend/internal/config"
	"discord_backend/internal/middlewares"
	"fmt"
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
	router := mux.NewRouter()

	authClient, _ := NewAuthClient(common.GrpcAuthAddress(cfg), cfg.Clients.Auth.Timeout, cfg.Clients.Auth.RetriesCount)
	router.HandleFunc("/api/register", authClient.Regsiter).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/login", authClient.Login).Methods(http.MethodPost, http.MethodOptions)

	relationsClient, _ := NewRelationsClient(common.GrpcRelationsAddress(cfg), cfg.Clients.Relations.Timeout, cfg.Clients.Relations.RetriesCount)
	router.HandleFunc("/api/relation", relationsClient.CreateNewRelation).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/friends/{id}", relationsClient.SendFriendOffer).Methods(http.MethodPost, http.MethodOptions)

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
