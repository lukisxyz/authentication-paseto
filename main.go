package main

import (
	"context"
	"flukis/login-system/src/account"
	"flukis/login-system/src/utils/cookie"
	"flukis/login-system/src/utils/token"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load file .env")
	}
	sKey := SymmetricKeyString()
	dbString := DbConnString()
	cookieHash, cookieBlock := CookieHashString()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}
	m, i, sl, kl, p, err := HashParams()
	if err != nil {
		log.Error().Err(err).Msg("unable to get env for hashing")
	}
	account.SetPool(pool)
	account.SetHashParam(m, i, sl, kl, p)
	token.SetSymetricKey(sKey)
	cookie.SetCookieKey(cookieHash, cookieBlock)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Mount("/account", account.Router())
	r.Mount("/login", account.RouterAuth())

	log.Info().Msg("starting up server...")
	if err := http.ListenAndServe(AddrString(), r); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}

	log.Info().Msg("server Stopped")
}
