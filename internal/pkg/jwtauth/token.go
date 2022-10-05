package jwtauth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JwtData struct {
	ID string
}

var hmacSignKey = []byte("qwerty123123")

func parseToken(tokenString string) (*JwtData, error) {
	token, err := jwt.Parse(tokenString, jwtParseHmacKey(hmacSignKey))

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtData := JwtData{}

		jwtData.ID = claim["id"].(string)

		return &jwtData, nil
	}

	return nil, errors.New("token invalid")
}

func jwtParseHmacKey(hmacSignKey []byte) func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSignKey, nil
	}
}

func CreateJwtToken(jwtData *JwtData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"id": jwtData.ID,
	})

	return token.SignedString(hmacSignKey)
}
