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

func NewUser(email string, passwordRaw string, name string) User {
	userID := NewUserID()
	password := NewPassword(passwordRaw)
	return User{
		UserID:   userID,
		Email:    email,
		Password: password,
		Name:     name,
	}
}

func NewUserID() UserID {
	return UserID(uuid.New().String())
}

type UserID string

func NewPassword(passwordRaw string) Password {
	return Password{
		PasswordRaw:  passwordRaw,
		PasswordSalt: uuid.New().String(),
	}
}

type Password struct {
	PasswordRaw  string
	PasswordSalt string
}

func (password *Password) Encrypt() string {
	crypto := sha512.New()
	crypto.Write([]byte(password.PasswordSalt + password.PasswordRaw))

	return base64.URLEncoding.EncodeToString(crypto.Sum(nil))
}

type UserRepository interface {
	IsEmailExists(email string) bool
	Store(user User) error
	FindByEmail(email string) *User
}
