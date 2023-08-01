package token

import (
	"errors"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"github.com/oklog/ulid/v2"
)

var (
	symetricKey     = []byte{}
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
	EndDuration     = time.Minute * 5
)

type Payload struct {
	Id        ulid.ULID `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func New(id ulid.ULID, email string) Payload {
	payload := Payload{
		Id:        id,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(EndDuration),
	}
	return payload
}

func SetSymetricKey(key string) error {
	if len(key) != chacha20poly1305.KeySize {
		return errors.New("invalid key size")
	}
	symetricKey = []byte(key)
	return nil
}

func (p *Payload) CreateToken() (string, error) {
	token, err := paseto.NewV2().Encrypt(symetricKey, p, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func Verify(token string) (*Payload, error) {
	payload := &Payload{}
	err := paseto.NewV2().Decrypt(token, symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
