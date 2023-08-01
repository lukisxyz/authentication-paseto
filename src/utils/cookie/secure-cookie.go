package cookie

import (
	"flukis/login-system/src/utils/token"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	s            = &securecookie.SecureCookie{}
	PasetoCookie = "state"
)

func SetCookieKey(hashKeyStr, blockKeyStr string) {
	hashKey := []byte(hashKeyStr)
	blockKey := []byte(blockKeyStr)
	s = securecookie.New(hashKey, blockKey)
}

type CookieValue struct {
	Token string `json:"token"`
}

func SetCookie(w http.ResponseWriter, pasetoToken string) {
	var c = CookieValue{
		Token: pasetoToken,
	}
	if encoded, err := s.Encode(PasetoCookie, c); err == nil {
		cookie := &http.Cookie{
			Name:     PasetoCookie,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(token.EndDuration),
		}
		http.SetCookie(w, cookie)
	}
}

func ReadCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	ck, err := r.Cookie(PasetoCookie)
	if err != nil {
		return "", err
	}
	var c CookieValue
	if err = s.Decode(PasetoCookie, ck.Value, &c); err != nil {
		return "", err
	}
	return c.Token, nil
}
