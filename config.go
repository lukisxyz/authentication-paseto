package main

import (
	"os"
	"strconv"
)

func DbConnString() string {
	conStr := os.Getenv("PSQL")
	return conStr
}

func CookieHashString() (string, string) {
	hk := os.Getenv("HASH_COOKIE_KEY")
	bk := os.Getenv("BLOCK_COOKIE_KEY")
	return hk, bk
}

func AddrString() string {
	conStr := os.Getenv("SERVER_ADDR")
	return conStr
}

func SymmetricKeyString() string {
	conStr := os.Getenv("SYMMETRIC_KEY")
	return conStr
}

func HashParams() (memory, iterations, saltLength, keyLength uint32, parallelism uint8, err error) {
	mStr := os.Getenv("HASH_MEMORY")
	mInt, err := strconv.Atoi(mStr)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	memory = uint32(mInt)

	iStr := os.Getenv("HASH_ITERATIONS")
	iInt, err := strconv.Atoi(iStr)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	iterations = uint32(iInt)

	slStr := os.Getenv("HASH_SALT_LENGTH")
	slInt, err := strconv.Atoi(slStr)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	saltLength = uint32(slInt)

	klStr := os.Getenv("HASH_SALT_LENGTH")
	klInt, err := strconv.Atoi(klStr)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	keyLength = uint32(klInt)

	tStr := os.Getenv("HASH_TIMES")
	tInt, err := strconv.Atoi(tStr)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	parallelism = uint8(tInt)

	return memory, iterations, saltLength, keyLength, parallelism, nil
}
