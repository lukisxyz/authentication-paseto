package account

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flukis/login-system/src/email"
	"flukis/login-system/src/utils/token"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func requestUpdatePassword(ctx context.Context, id ulid.ULID) error {
	account, err := getById(ctx, id)
	if err != nil {
		return err
	}

	addr := os.Getenv("ADDR")
	bin := make([]byte, 64)
	rand.Read(bin)
	state := base64.URLEncoding.EncodeToString(bin)
	rand.Read(bin)
	state2 := base64.URLEncoding.EncodeToString(bin)
	buildLink := fmt.Sprintf("http://%s/c/%s", addr, state)
	emailBody := email.EmailLinkForgotPasswordBodyRequest{
		SUBJECT: "Forgot Password Link",
		EMAIL:   account.Email,
		CODE:    buildLink,
	}
	fmt.Println(state2)
	go email.SendLinkForgotPassword(account.Email, emailBody)

	return nil
}

func createAccount(ctx context.Context, email, pwd string) error {
	account, err := NewAccount(email, pwd)
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = getByEmail(ctx, email)
	if err == nil {
		tx.Rollback(ctx)
		return nil
	}
	fmt.Println(err)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		tx.Rollback(ctx)
		return err
	}

	err = save(ctx, tx, account)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	return nil
}

var (
	ErrWrongPassword = errors.New("wrong password")
)

func loginAccount(ctx context.Context, email, pwd string) (string, error) {
	account, err := getByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAccountNotFound
		}
		return "", err
	}
	ok, err := hashParams.compare(pwd, account.Password)
	if err != nil {
		return "", err
	}
	tokenPayload := token.New(account.Id, account.Email)
	tokenStr, err := tokenPayload.CreateToken()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrWrongPassword
	}
	return tokenStr, nil
}
