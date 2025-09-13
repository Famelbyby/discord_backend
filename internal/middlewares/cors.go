package middlewares

import (
	"idea-store-auth/internal/utils"
	"net/http"
)

func setupCORS(w http.ResponseWriter, req *http.Request, sessionLifetime string) {
	allowedOrigins := map[string]struct{}{
		"http://0.0.0.0:3000":     {},
		"http://0.0.0.0:8000":     {},
		"http://81.177.223.234":   {},
		"https://ideastore.space": {},
	}
	origin := req.Header.Get("Origin")
	if _, ok := allowedOrigins[origin]; ok {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods",
		http.MethodPost+", "+http.MethodGet+", "+http.MethodOptions+", "+http.MethodPut+", "+http.MethodDelete+", "+http.MethodPatch)
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, Authorization, Access-Control-Allow-Origin, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", sessionLifetime)
	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		setupCORS(w, req, utils.DefaultSessionLifetimeString)

		// Если метод OPTIONS, то возвращаем пустой ответ с нужными заголовками
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, req)
	})
}
