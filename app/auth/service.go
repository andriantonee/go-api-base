package auth

import (
	"fmt"
	"go-api-base/model"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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
	UserID   model.UserID `json:"user_id"`
	IssuedAt int64        `json:"iat"`
}

func (jwtC jwtClaims) Valid() error {
	return nil
}

func (service *service) NewIdentifier(userID model.UserID) (string, error) {
	issuedAt := time.Now()

	claims := jwtClaims{
		UserID:   userID,
		IssuedAt: issuedAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(service.jwtSecretKey)
}

func (service *service) Authorize(tokenString string) (*model.UserID, error) {
	var claims jwtClaims

	_, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
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

	return &claims.UserID, nil
}
