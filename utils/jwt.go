package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"os"
)

func GenerateJWT(nid string) (string, error) {
	claims := jwt.MapClaims{
		"nid":  nid,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
