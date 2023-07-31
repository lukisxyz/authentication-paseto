package account

import (
	"context"
	"errors"
)

func createAccount(ctx context.Context, email, pwd string) error {
	account, err := NewAccount(email, pwd)
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
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

func loginAccount(ctx context.Context, email, pwd string) (Account, error) {
	account, err := getByEmail(ctx, email)
	if err != nil {
		return Account{}, err
	}
	ok, err := hashParams.compare(pwd, account.Password)
	if err != nil {
		return Account{}, err
	}
	if !ok {
		return Account{}, ErrWrongPassword
	}
	return account, nil
}
