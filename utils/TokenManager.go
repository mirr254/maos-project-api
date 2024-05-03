package utils

import (
	"maos-cloud-project-api/models"

	"github.com/dgrijalva/jwt-go"
	"crypto/rand"
	"encoding/base64"
)

/*
 GenerateToken function generates a token for email verification
return: string, error
*/
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

/*
 ParseToken function receives a cookie as tokenString and gets the details stored in that cookie for logged
in users
Params: tokenString string
return: claims *models.Claims, err error
*/
func ParseToken(tokenString string) (claims *models.Claims, err error) {

	secretKey, err := GetSecretKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
