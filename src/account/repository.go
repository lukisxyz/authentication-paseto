package account

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func save(ctx context.Context, tx pgx.Tx, a Account) error {
	query := `
		INSERT INTO accounts
			(id, email, password, verification_code, created_at)
		VALUES
			($1, $2, $3, $4, $5)
	`
	_, err := tx.Exec(
		ctx,
		query,
		a.Id,
		a.Email,
		a.Password,
		a.VerificationCode,
		a.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func getByEmail(ctx context.Context, email string) (Account, error) {
	query := `
		SELECT
			id, email, password, verification_code, verified_at, created_at, updated_at
		FROM
			accounts
		WHERE email = $1
	`
	row := pool.QueryRow(
		ctx,
		query,
		&email,
	)
	account := Account{}
	if err := row.Scan(
		&account.Id,
		&account.Email,
		&account.Password,
		&account.VerificationCode,
		&account.VerifiedAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any account")
			return Account{}, err
		}
		return Account{}, err
	}
	return account, nil
}

func getById(ctx context.Context, id ulid.ULID) (Account, error) {
	query := `
		SELECT
			id, email, password, verification_code, verified_at, created_at, updated_at
		FROM
			accounts
		WHERE id = $1
	`
	row := pool.QueryRow(
		ctx,
		query,
		&id,
	)
	account := Account{}
	if err := row.Scan(
		&account.Id,
		&account.Email,
		&account.Password,
		&account.VerificationCode,
		&account.VerifiedAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any account")
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}
	return account, nil
}
