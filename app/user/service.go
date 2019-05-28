package user

import (
	"errors"
	"go-api-base/model"
)

var (
	ErrDuplicateEmail   = errors.New("this email already exists")
	ErrAccountNotExists = errors.New("this account is not exists")
)

type Service interface {
	Register(user model.User) error
	FindUserID(email string, passwordRaw string) (*model.UserID, error)
}

func NewService(userRepository model.UserRepository) Service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository model.UserRepository
}

func (service *service) Register(user model.User) error {
	if service.userRepository.IsEmailExists(user.Email) {
		return ErrDuplicateEmail
	}

	if err := service.userRepository.Store(user); err != nil {
		return ErrDuplicateEmail
	}

	return nil
}

func (service *service) FindUserID(email string, passwordRaw string) (*model.UserID, error) {
	user := service.userRepository.FindByEmail(email)

	if user == nil {
		return nil, ErrAccountNotExists
	}

	password := user.Password
	password.PasswordRaw = passwordRaw
	if password.Encrypt() != user.PasswordEncrypted {
		return nil, ErrAccountNotExists
	}

	return &user.UserID, nil
}
