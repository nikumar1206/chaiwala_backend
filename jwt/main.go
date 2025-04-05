package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SIGNING_KEY = []byte(os.Getenv("SIGNING_KEY"))
	issuer      = "chaiwala"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateTokens(username string) (string, string, time.Time, error) {
	accessExp := time.Now().Add(4 * time.Hour)
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			// todo(nick): pull from app config
			Issuer: issuer,
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := at.SignedString(SIGNING_KEY)
	if err != nil {
		return "", "", accessExp, err
	}

	refreshClaims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    issuer,
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString(SIGNING_KEY)
	if err != nil {
		return "", "", accessExp, err
	}

	return signedToken, refreshToken, accessExp, nil
}

func ValidateToken(token string) (Claims, error) {
	claims := new(Claims)

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return SIGNING_KEY, nil
	})

	if err != nil || !t.Valid {
		fmt.Println(err.Error())
		return *claims, errors.New("Invalid or expired token")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return *claims, errors.New("Token has expired")
	}

	return *claims, nil
}
