package main

import (
	"bytes"
	"context"
	authv1 "discord_backend/gen/go/auth"
	"discord_backend/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxAge = time.Hour * 24 * 7
)

type AuthClient struct {
	authApi authv1.AuthClient
}

type loginRequest struct {
	Email    string `json:"mail"`
	Password string `json:"password"`
}

func (c *AuthClient) Login(w http.ResponseWriter, r *http.Request) {

	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
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
	type ProfilesResponse struct {
		Profiles []Profile `json:"profiles"`
	}

	resp, err := http.Get("http://profile_py:9999/api/profile?email=" + req.Email)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}
	// Обязательно закрываем тело ответа после завершения работы функции.
	defer resp.Body.Close()

	// Читаем тело ответа.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	// Создаем переменную для хранения результата разбора.
	var profilesResponse ProfilesResponse

	// Разбираем (анмаршалим) JSON из тела ответа в нашу структуру.
	if err := json.Unmarshal(body, &profilesResponse); err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		Id:       profilesResponse.Profiles[0].ID,
	}

	loginResponse, err := c.authApi.Login(context.Background(), request)

	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			utils.WriteError(w, "Invalid credentials")
			return
		}
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    string(loginResponse.SessionId),
		MaxAge:   int(maxAge / 1_000_000_000),
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	responseJson, err := json.Marshal(profilesResponse.Profiles[0])
	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

type registerRequest struct {
	Email    string `form:"mail"`
	Password string `form:"password"`
}

func (c *AuthClient) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Не удалось разобрать multipart-форму", http.StatusBadRequest)
		return
	}

	// 2. Создание нового тела запроса.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Копирование текстовых полей (mail, username, status).
	for key, values := range r.MultipartForm.Value {
		for _, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				slog.Error("client Regsiter error: " + err.Error())
				utils.WriteError(w, "Internal error")
				return
			}
		}
	}

	// Копирование файла (avatar).
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		part, err := writer.CreateFormFile("avatar", header.Filename)
		if err != nil {
			slog.Error("client Regsiter error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			slog.Error("client Regsiter error: " + err.Error())
			utils.WriteError(w, "Internal error")
			return
		}
	}

	// Завершаем формирование multipart-тела.
	if err := writer.Close(); err != nil {
		slog.Error("client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	// 3. Создание и отправка нового запроса.
	targetURL := "http://profile_py:9999/api/profile"
	postReq, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		slog.Error("http.NewRequest client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	// Устанавливаем Content-Type с правильной границей (boundary).
	postReq.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос.
	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		slog.Error(" client.Do client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}
	defer resp.Body.Close()

	req := registerRequest{
		Email:    r.FormValue("mail"),
		Password: r.FormValue("password"),
	}

	request := &authv1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	registerResponse, err := c.authApi.Register(r.Context(), request)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			utils.WriteError(w, "User already exists")
			return
		}

		slog.Error(" c.authAPi.Register client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	_, err = json.Marshal(registerResponse)
	if err != nil {
		slog.Error(" json.Marshal client Regsiter error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	// 4. (Опционально) Пересылка ответа от целевого сервера обратно клиенту.
	// Копируем заголовки ответа.

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	// Устанавливаем статус-код ответа.
	w.WriteHeader(resp.StatusCode)
	// Копируем тело ответа.
	io.Copy(w, resp.Body)
}

func (c *AuthClient) Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "unauthorized")
		return
	}

	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &authv1.LogoutRequest{
		SessionId: sessionCookie.Value,
	}

	_, err = c.authApi.Logout(context.Background(), request)

	if err != nil {
		if strings.Contains(err.Error(), "no session") {
			slog.Error("client login error: no session")
			utils.WriteError(w, "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthClient) IsRegistered(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "unauthorized")
		return
	}

	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	request := &authv1.IsRegisteredRequest{
		SessionId: sessionCookie.Value,
	}

	isRegisteredResponse, err := c.authApi.IsRegistered(context.Background(), request)

	if err != nil {
		if strings.Contains(err.Error(), "no session") {
			slog.Error("client login error: no session")
			utils.WriteError(w, "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	var ids = []string{isRegisteredResponse.UserId}

	profiles := GetProfiles(ids)
	profileJson, err := json.Marshal(profiles[0])

	if err != nil {
		slog.Error("client Login error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(profileJson)
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
		authApi: authv1.NewAuthClient(cc),
	}, nil
}
