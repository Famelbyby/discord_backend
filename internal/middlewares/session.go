package middlewares

import (
	"context"
	authv1 "discord_backend/gen/go/auth"
	"discord_backend/internal/utils"
	"log/slog"
	"net/http"
	"strings"

	"google.golang.org/grpc"
)

const UserKey string = "userId"

type SessionClient interface {
	IsRegistered(ctx context.Context, r *authv1.IsRegisteredRequest, options ...grpc.CallOption) (*authv1.IsRegisteredResponse, error)
}

func SessionMiddleware(auth SessionClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			r, err := req.Cookie("session_id")

			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			request := &authv1.IsRegisteredRequest{
				SessionId: r.Value,
			}

			isRegistered, err := auth.IsRegistered(context.Background(), request)

			if err != nil {
				if strings.Contains(err.Error(), "no session") {
					slog.Error("client login error: no session")
					w.WriteHeader(http.StatusForbidden)
					utils.WriteError(w, "forbidden")

					return
				}

				slog.Error("client Login error: " + err.Error())
				utils.WriteError(w, "Internal error")
				return
			}

			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), UserKey, isRegistered.UserId)))
		})
	}
}
