package model

import (
	"github.com/google/uuid"

	"crypto/sha512"
	"encoding/base64"
)

type User struct {
	UserID            UserID
	Email             string
	PasswordEncrypted string
	Password          Password
	Name              string
}

func NewUserID() UserID {
	return UserID(uuid.New().String())
}

type UserID string

func NewPassword(password string) Password {
	return Password{
		Password:     password,
		PasswordSalt: uuid.New().String(),
	}
}

type Password struct {
	Password     string
	PasswordSalt string
}

func (password *Password) Encrypt() string {
	crypto := sha512.New()
	crypto.Write([]byte(password.PasswordSalt + password.Password))

	return base64.URLEncoding.EncodeToString(crypto.Sum(nil))
}

type UserRepository interface {
	IsEmailExists(email string) bool
	Store(user User) error
	FindByEmail(email string) User
}
