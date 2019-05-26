package auth

import (
	"errors"
	"fmt"
	"go-api-base/model"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var ErrNotAuthorize = errors.New("not authorize")

type Service interface {
	NewIdentifier(userID model.UserID) (string, error)
	Authorize(tokenString string) (*model.UserID, error)
}

func NewService(jwtSecretKey []byte) Service {
	return &service{
		jwtSecretKey: jwtSecretKey,
	}
}

type service struct {
	jwtSecretKey []byte
}

type jwtClaims struct {
	UserID model.UserID `json:"user_id"`
	jwt.StandardClaims
}

func (service *service) NewIdentifier(userID model.UserID) (string, error) {
	issuedAt := time.Now()

	claims := jwtClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: issuedAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(service.jwtSecretKey)
}

func (service *service) Authorize(tokenString string) (*model.UserID, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return service.jwtSecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtClaims)
	if ok {
		return &claims.UserID, nil
	}

	return nil, ErrNotAuthorize
}
