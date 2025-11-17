package main

import (
	"bytes"
	"discord_backend/internal/utils"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"sync"
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

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, id := range profileIds {
		wg.Add(1)

		go func() {
			defer wg.Done()

			profile, err := GetProfileById(id)

			if err != nil {
				slog.Error("client profile error: " + err.Error())
			} else {
				mutex.Lock()

				profiles = append(profiles, profile)

				mutex.Unlock()
			}
		}()
	}

	wg.Wait()

	return profiles
}

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	url := "http://profile_py:9999" + r.URL.RequestURI()

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
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func HandleGetProfileById(w http.ResponseWriter, r *http.Request) {
	url := "http://profile_py:9999" + r.URL.RequestURI()

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
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func HandleDeleteProfileById(w http.ResponseWriter, r *http.Request) {
	url := "http://profile_py:9999" + r.URL.RequestURI()

	resp, err := http.NewRequest(http.MethodDelete, url, r.Body)
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
	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func HandleEditProfileById(w http.ResponseWriter, r *http.Request) {
	url := "http://profile_py:9999" + r.URL.RequestURI()

	resp, err := http.NewRequest(http.MethodPut, url, r.Body)
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
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
