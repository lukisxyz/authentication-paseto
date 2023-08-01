package account

import (
	"flukis/login-system/src/utils/helper"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	Id               ulid.ULID
	Email            string
	Password         string
	VerificationCode string
	VerifiedAt       null.Time
	CreatedAt        time.Time
	UpdatedAt        null.Time
}

func NewAccount(email, pwd string) (Account, error) {
	vcode := helper.RandString(10)
	hashedPassword, err := hashParams.generate(pwd)
	if err != nil {
		return Account{}, err
	}
	account := Account{
		Id:               ulid.Make(),
		Email:            email,
		Password:         hashedPassword,
		VerificationCode: vcode,
		CreatedAt:        time.Time{},
	}
	return account, nil
}
