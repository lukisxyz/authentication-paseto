package account

import (
	"encoding/json"
	"flukis/login-system/src/utils/cookie"
	custommiddleware "flukis/login-system/src/utils/custom-middleware"
	"flukis/login-system/src/utils/response"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Post("/", createAccountHandler)
	r.Group(func(rt chi.Router) {
		rt.Use(custommiddleware.PasetoMiddleware)
		rt.Post("/change-password", reqUpdatePassword)
	})

	return r
}

func reqUpdatePassword(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := ctx.Value(custommiddleware.ContextKeyUserId).(ulid.ULID)
	err := requestUpdatePassword(ctx, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Response(w, "please check your email", http.StatusOK, nil, nil)
}

func createAccountHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	_, err = mail.ParseAddress(input.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	ctx := req.Context()
	err = createAccount(ctx, input.Email, input.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Response(w, "success create an account, please check your email for verification", http.StatusCreated, nil, nil)
}

func RouterAuth() *chi.Mux {
	r := chi.NewMux()
	r.Post("/", loginAccountHandler)

	return r
}

func loginAccountHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	_, err = mail.ParseAddress(input.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	ctx := req.Context()
	token, err := loginAccount(ctx, input.Email, input.Password)
	if err != nil {
		if err == ErrAccountNotFound {
			response.Error(w, http.StatusNotFound, err)
			return
		}
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	cookie.SetCookie(w, token)

	response.Response(w, "success login", http.StatusOK, nil, nil)
}
