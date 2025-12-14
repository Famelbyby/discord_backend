package main

import (
	common "discord_backend/cmd"
	"discord_backend/internal/config"
	"discord_backend/internal/middlewares"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

var (
	AppPortEnv      = "APP_PORT"
	AppSecretEnv    = "APP_SECRET"
	TokenCookieName = "token"
	GOMAXPROCS      = 5
	EmptyValue      = int64(-1)
)
var profilesClient ProfilesClient

func main() {

	cfg := config.MustLoad()
	router := mux.NewRouter()
	authedRouter := router.NewRoute().Subrouter()
	runtime.GOMAXPROCS(GOMAXPROCS)

	InitKafka()
	defer kafkaWriter.Close()
	defer kafkaReader.Close()

	// 2. Запуск Consumer'а для чтения chat-out в фоне
	go StartKafkaConsumer()

	relationsClient, _ := NewRelationsClient(common.GrpcRelationsAddress(cfg), cfg.Clients.Relations.Timeout, cfg.Clients.Relations.RetriesCount)
	profilesClient = ProfilesClient{
		relationsApi: relationsClient.relationsAPi,
	}
	authClient, _ := NewAuthClient(common.GrpcAuthAddress(cfg), cfg.Clients.Auth.Timeout, cfg.Clients.Auth.RetriesCount, relationsClient.relationsAPi)

	authedRouter.Use(middlewares.SessionMiddleware(authClient.authApi))

	router.HandleFunc("/api/register", authClient.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/login", authClient.Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/is-registered", authClient.IsRegistered).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/logout", authClient.Logout).Methods(http.MethodPost, http.MethodOptions)

	authedRouter.HandleFunc("/api/relation", relationsClient.CreateNewRelation).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/friends/{id}", relationsClient.SendFriendOffer).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/friends/{id}", relationsClient.RemoveFriend).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/outcoming/{id}/delete", relationsClient.CancelFriendOffer).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/incoming/{id}/accept", relationsClient.AcceptFriendOffer).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/incoming/{id}/decline", relationsClient.DeclineFriendOffer).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/blocks/{id}", relationsClient.BlockUser).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/blocks/{id}", relationsClient.UnblockUser).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/relation/{id}", relationsClient.GetRelation).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/relations", relationsClient.GetUserRelations).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/friends/{id}", relationsClient.GetAllFriends).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/incoming/{id}", relationsClient.GetAllIncomingOffers).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/outcoming/{id}", relationsClient.GetAllOutgoingOffers).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/blocks/{id}", relationsClient.GetAllBlockedUsers).Methods(http.MethodGet, http.MethodOptions)

	authedRouter.HandleFunc("/api/profile", profilesClient.HandleGetProfile).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/profile", profilesClient.HandleSaveProfile).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/profile/{id}", profilesClient.HandleGetProfileById).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/profile/{id}", profilesClient.HandleDeleteProfileById).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/profile/{id}", profilesClient.HandleEditProfileById).Methods(http.MethodPut, http.MethodOptions)

	authedRouter.HandleFunc("/api/chat", HandleGetUserChats).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/chat", HandleCreateChat).Methods(http.MethodPost, http.MethodOptions)
	authedRouter.HandleFunc("/api/chat/{id}", HandleGetChatById).Methods(http.MethodGet, http.MethodOptions)
	authedRouter.HandleFunc("/api/chat/{id}", HandleDeleteChatById).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/chat/{id}/delete_user", HandleDeleteUserFromChat).Methods(http.MethodDelete, http.MethodOptions)
	authedRouter.HandleFunc("/api/chat/{id}", HandleUpdateChat).Methods(http.MethodPut, http.MethodOptions)
	authedRouter.HandleFunc("/api/ws/chat", ChatWsHandler).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)

	handler := middlewares.CorsMiddleware(router)
	fmt.Println("Server is listening...")

	log.Fatal(http.ListenAndServe(os.Getenv(AppPortEnv), handler))
}
