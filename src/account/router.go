package account

import (
	"encoding/json"
	"flukis/ecommerce/src/utils/response"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Post("/", createAccountHandler)

	return r
}

func RouterAuth() *chi.Mux {
	r := chi.NewMux()
	r.Post("/", loginAccountHandler)

	return r
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
	account, err := loginAccount(ctx, input.Email, input.Password)
	if err != nil {
		if err == ErrAccountNotFound {
			response.Error(w, http.StatusNotFound, err)
			return
		}
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	response.Response(w, "success login", http.StatusOK, account, nil)
}
