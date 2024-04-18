package token_service

import (
	"errors"
	token_model "paper_back/models/token"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	*jwt.StandardClaims
	Username string
	Role     string
	Id       int64
}

type Payload struct {
	Username string
	Role     string
	Id       int64
}

var SecretAccess = []byte("fg54gSDv45njg45bgDFSDFn4u!@435SDFn45")

func CreateToken(payload Payload) (token_model.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		payload.Username,
		payload.Role,
		payload.Id,
	}

	tokenString, err := token.SignedString(SecretAccess)
	if err != nil {
		return token_model.Token{}, errors.New("error create token")
	}

	tokenFromDatabase, err := token_model.AddToken(tokenString, payload.Id)
	if err != nil {
		return token_model.Token{}, err
	}

	return tokenFromDatabase, nil
}
