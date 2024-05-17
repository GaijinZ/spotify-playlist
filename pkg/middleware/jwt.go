package middleware

import (
	"time"

	"spf-playlist/pkg/config"
	"spf-playlist/users/handler/models"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models.User, cfg config.GlobalEnv) (string, error) {
	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (claims *models.Claims, err error) {
	var cfg config.GlobalEnv

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
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
