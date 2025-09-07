package service

import (
	"crypto/rand"
	"db-rest/config"
	"github.com/golang-jwt/jwt"
	"time"
)

var tokenCache = map[string]string{}

var defaultSec = rand.Text()

var secret = config.GetEnvValue[string](config.VIPER_KEY_JWT_SECRET)

func IsValidToken(token string) bool {

	contain := false
	for _, v := range tokenCache {
		if v == token {
			contain = true
			break
		}
	}

	if contain {

		token, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {

			if secret != "" {
				return []byte(secret), nil
			}
			return []byte(defaultSec), nil
		})

		if err == nil && token.Valid {
			return true
		}

	}

	return false
}

func GenToken(user string) (tokenStr string, err error) {

	oldToken := tokenCache[user]

	if oldToken != "" && IsValidToken(oldToken) {
		return oldToken, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    "db-rest",
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user,
		Id:        user,
	})

	if secret != "" {
		tokenStr, err = token.SignedString([]byte(secret))
	} else {
		tokenStr, err = token.SignedString([]byte(defaultSec))
	}

	if err == nil {
		tokenCache[user] = tokenStr
	}

	return tokenStr, err
}
