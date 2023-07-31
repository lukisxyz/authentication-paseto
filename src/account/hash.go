package account

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltlength  uint32
	keylength   uint32
}

var (
	hashParams             argon2Params
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func SetHashParam(memory, iterations, saltLength, keyLength uint32, parallelism uint8) {
	hashParams.memory = memory
	hashParams.iterations = iterations
	hashParams.saltlength = saltLength
	hashParams.keylength = keyLength
	hashParams.parallelism = parallelism
}

func generateRandomwBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (a *argon2Params) generate(pwd string) (string, error) {
	salt, err := generateRandomwBytes(a.saltlength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey(
		[]byte(pwd),
		salt,
		a.iterations,
		a.memory,
		a.parallelism,
		a.keylength,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.memory,
		a.iterations,
		a.parallelism,
		b64Salt,
		b64Hash)
	return encodedHash, nil
}

func (a *argon2Params) compare(incoming, stored string) (bool, error) {
	p, salt, hash, err := a.decode(stored)
	if err != nil {
		return false, err
	}
	incomingHashed := argon2.IDKey(
		[]byte(incoming),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keylength,
	)

	if subtle.ConstantTimeCompare(hash, incomingHashed) == 1 {
		return true, nil
	}

	return false, nil
}

func (a *argon2Params) decode(encodedHash string) (p *argon2Params, salt, hash []byte, err error) {
	valStr := strings.Split(encodedHash, "$")
	if len(valStr) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(valStr[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argon2Params{}
	_, err = fmt.Sscanf(valStr[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(valStr[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltlength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(valStr[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keylength = uint32(len(hash))

	return p, salt, hash, nil
}
