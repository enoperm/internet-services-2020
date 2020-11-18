package model

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/scrypt"
)

// TODO: Move errors and pw related stuff.
var (
	ErrRegisterOk                  = errors.New("success")
	ErrRegisterUsernameUnavailable = errors.New("requested username is not available")
	ErrRegisterInvalidUsername     = errors.New("username must not be empty and may only consists of printable ASCII characters")
	ErrRegisterInvalidPassword     = errors.New("password must be at least <X> characters long")
)

var (
	ErrAuthenticationOk   = errors.New("successful login")
	ErrAuthenticationFail = errors.New("invalid username or password")
)

const (
	HASH_LEN = 32
	SALT_LEN = 8
)

func GenerateRandomSalt() (*[SALT_LEN]byte, error) {
	var salt [SALT_LEN]byte
	n, err := rand.Read(salt[:])

	switch {
	case err != nil:
		return nil, err
	case n != SALT_LEN:
		return nil, errors.New("failed to generate salt")
	}

	return &salt, nil
}

func NewPassword(value []byte, salt []byte) (*Password, error) {
	hash, err := scrypt.Key(value, salt, 32768, 8, 1, HASH_LEN)
	if err != nil {
		return nil, err
	}

	var fsHash [HASH_LEN]byte
	var fsSalt [SALT_LEN]byte

	copy(fsHash[:], hash)
	copy(fsSalt[:], salt)

	p := Password{
		Hash: fsHash,
		Salt: fsSalt,
	}
	return &p, nil
}

type User struct {
	Id       int64     `json:"-"`
	Name     string    `json:"name"`
	Password *Password `json:"-"`
}

type Password struct {
	Hash [HASH_LEN]byte
	Salt [SALT_LEN]byte

	_ struct{} `json:"-"`
}

func (p Password) Check(v []byte) error {
	result := ErrAuthenticationFail

	hashedInput, err := NewPassword(v, p.Salt[:])
	if err != nil {
		return result
	}

	if subtle.ConstantTimeCompare(p.Hash[:], hashedInput.Hash[:]) == 1 {
		result = ErrAuthenticationOk
	}
	return result
}
