package main

import (
	"bytes"
	"discord_backend/internal/middlewares"
	"discord_backend/internal/utils"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	profileUrl = "http://profile_py:9999"
)

type Profile struct {
	ID        string `json:"id"`
	ShortLink string `json:"short_link"`
	Mail      string `json:"mail"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	AvatarURL string `json:"avatar_url"`
	Status    string `json:"status"`
}

func GetProfiles(profileIds []string) []Profile {
	if len(profileIds) == 0 {
		return []Profile{}
	}

	type ProfilesClientResponse struct {
		Profiles []Profile `json:"profiles"`
	}

	profilesRequestBody, err := json.Marshal(profileIds)

	if err != nil {
		slog.Error("http.NewRequest client profiles error: " + err.Error())
		return []Profile{}
	}

	targetURL := profileUrl + "/api/profile/by-array"
	postReq, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(profilesRequestBody))
	if err != nil {
		slog.Error("http.NewRequest client profiles error: " + err.Error())
		return []Profile{}
	}

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		slog.Error(" client.Do client profiles error: " + err.Error())
		return []Profile{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client profiles error: " + err.Error())
		return []Profile{}
	}

	var profiles ProfilesClientResponse

	if err := json.Unmarshal(body, &profiles); err != nil {
		slog.Error("client profiles error: " + err.Error())
		return []Profile{}
	}

	return profiles.Profiles
}

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	url := profileUrl + r.URL.RequestURI()

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	defer resp.Body.Close()

	// Читаем тело ответа.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleGetProfileById(w http.ResponseWriter, r *http.Request) {
	url := profileUrl + r.URL.RequestURI()

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

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleDeleteProfileById(w http.ResponseWriter, r *http.Request) {
	err := utils.CompareUserIDsInParams(r, middlewares.UserKey, "id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

	url := profileUrl + r.URL.RequestURI()

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

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func HandleSaveProfile(w http.ResponseWriter, r *http.Request) {
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
	targetURL := profileUrl + "api/profile"
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

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func HandleEditProfileById(w http.ResponseWriter, r *http.Request) {
	err := utils.CompareUserIDsInParams(r, middlewares.UserKey, "id")

	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		w.WriteHeader(http.StatusForbidden)
		utils.WriteError(w, "forbidden")

		return
	}

	url := profileUrl + r.URL.RequestURI()

	req, err := http.NewRequest(http.MethodPut, url, r.Body)
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

	// Читаем тело ответа.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("client HandleEditProfileById error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
