package custommiddleware

import (
	"context"
	"flukis/login-system/src/utils/cookie"
	"flukis/login-system/src/utils/token"
	"net/http"
	"time"
)

var payload *token.Payload

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyUserId = contextKey("user")
)

func PasetoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ck, err := cookie.ReadCookie(w, req)
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusUnauthorized)
			return
		}
		payload, err = token.Verify(ck)
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusUnauthorized)
			return
		}
		if payload.ExpiredAt.Before(time.Now()) {
			http.Redirect(w, req, "/login", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), ContextKeyUserId, payload.Id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
