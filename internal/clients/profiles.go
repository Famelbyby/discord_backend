package main

import (
	"bytes"
	relationsv1 "discord_backend/gen/go/relations"
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

type RelatedProfile struct {
	ID          string `json:"id"`
	ShortLink   string `json:"short_link"`
	Mail        string `json:"mail"`
	Username    string `json:"username"`
	CreatedAt   int64  `json:"created_at"`
	AvatarURL   string `json:"avatar_url"`
	Status      string `json:"status"`
	IsFriend    bool   `json:"isFriend"`
	IsIncoming  bool   `json:"isIncoming"`
	IsOutcoming bool   `json:"isOutcoming"`
	IsBlocked   bool   `json:"isBlocked"`
}

type ProfilesClient struct {
	relationsApi relationsv1.RelationsClient
}

func (c *ProfilesClient) GetProfiles(profileIds []string) []Profile {
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

func (c *ProfilesClient) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	requestUrl := profileUrl + r.URL.RequestURI()

	resp, err := http.Get(requestUrl)
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

	type ProfilesClientResponse struct {
		Profiles []Profile `json:"profiles"`
	}

	var profiles ProfilesClientResponse
	if err := json.Unmarshal(body, &profiles); err != nil {
		slog.Error("client profiles error: " + err.Error())
		return
	}

	userId := r.Context().Value(middlewares.UserKey).(string)

	targetIds := []string{}
	for _, p := range profiles.Profiles {
		targetIds = append(targetIds, p.ID)
	}

	userRelationsResponse, err := c.relationsApi.GetUserRelations(r.Context(), &relationsv1.GetUserRelationsRequest{
		SenderId:  userId,
		TargetIds: targetIds,
	})
	if err != nil {
		slog.Error("client HandleGetProfile error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	type RelatedProfiles struct {
		Profiles []RelatedProfile `json:"profiles"`
	}
	var relatedProfiles RelatedProfiles

	for i, p := range profiles.Profiles {
		relatedProfiles.Profiles = append(relatedProfiles.Profiles, RelatedProfile{
			ID:          p.ID,
			AvatarURL:   p.AvatarURL,
			Mail:        p.Mail,
			ShortLink:   p.ShortLink,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
			Username:    p.Username,
			IsFriend:    userRelationsResponse.Datas[i].IsFriend,
			IsIncoming:  userRelationsResponse.Datas[i].IsIncoming,
			IsOutcoming: userRelationsResponse.Datas[i].IsOutgoing,
			IsBlocked:   userRelationsResponse.Datas[i].IsBlocked,
		})
	}
	respJson, err := json.Marshal(relatedProfiles)
	if err != nil {
		slog.Error("[GetRelation] client error: " + err.Error())
		utils.WriteError(w, "Internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (c *ProfilesClient) HandleGetProfileById(w http.ResponseWriter, r *http.Request) {
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

func (c *ProfilesClient) HandleDeleteProfileById(w http.ResponseWriter, r *http.Request) {
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

func (c *ProfilesClient) HandleSaveProfile(w http.ResponseWriter, r *http.Request) {
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

func (c *ProfilesClient) HandleEditProfileById(w http.ResponseWriter, r *http.Request) {
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
