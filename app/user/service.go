package user

import (
	"errors"
	"fmt"
	"go-api-base/model"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var ErrDuplicateEmail = errors.New("this email already exists")
var ErrAccountNotExists = errors.New("this account is not exists")
var ErrNotAuthorize = errors.New("not authorize")

type Service interface {
	Register(user model.User) (string, error)
	Login(email string, password string) (string, error)
}

func NewService(
	jwtSecretKey []byte,
	userRepository model.UserRepository,
) Service {
	return &service{
		jwtSecretKey:   jwtSecretKey,
		userRepository: userRepository,
	}
}

type service struct {
	jwtSecretKey   []byte
	userRepository model.UserRepository
}

func (service *service) Register(
	user model.User,
) (string, error) {
	if service.userRepository.IsEmailExists(user.Email) {
		return "", ErrDuplicateEmail
	}

	if err := service.userRepository.Store(user); err != nil {
		return "", err
	}

	tokenString, err := service.newJWTTokenString(user.UserID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *service) Login(email string, password string) (string, error) {
	user := service.userRepository.FindByEmail(email)

	if user.UserID == "" {
		return "", ErrAccountNotExists
	}

	authPassword := user.Password
	authPassword.Password = password
	if authPassword.Encrypt() != user.PasswordEncrypted {
		return "", ErrAccountNotExists
	}

	tokenString, err := service.newJWTTokenString(user.UserID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type jwtClaims struct {
	UserID model.UserID `json:"user_id"`
	jwt.StandardClaims
}

func (service *service) newJWTTokenString(userID model.UserID) (string, error) {
	issueAt := time.Now()

	lifeTimeDuration, _ := time.ParseDuration("20m")
	expiresAt := issueAt.Add(lifeTimeDuration)

	claims := jwtClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  issueAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(service.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *service) parseJWTTokenString(
	tokenString string,
) (*jwtClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return service.jwtSecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtClaims)
	if ok && token.Valid {
		return &claims, nil
	}

	return nil, ErrNotAuthorize
}
