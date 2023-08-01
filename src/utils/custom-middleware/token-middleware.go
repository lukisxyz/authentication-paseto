package custommiddleware

import (
	"context"
	"flukis/login-system/src/utils/cookie"
	"flukis/login-system/src/utils/token"
	"net/http"
)

var payload *token.Payload

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
		ctx := context.WithValue(req.Context(), "user", payload.Id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
