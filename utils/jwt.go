package utils

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

var signingKey []byte

func init() {
	str, _ := strconv.Unquote(secretKey)
	key, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Panicf("cannot get the signing key: %v", err.Error())
	}
	signingKey = key
}

func GenerateToken(email string, userId int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(signingKey)
}

func VerifyToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, isValidClaimType := parsedToken.Claims.(jwt.MapClaims)

	if !isValidClaimType {
		return nil, errors.New("unexpected token claims")
	}

	return claims, nil
}
