package jwt

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var (
	SIGNING_KEY = []byte(os.Getenv("SIGNING_KEY"))
	// todo(nick): pull from app config
	issuer = "chaiwala"

	ErrInvalidToken         = errors.New("Invalid token")
	ErrExpiredToken         = errors.New("Expired token")
	ErrInvalidSigningMethod = errors.New("Unexpected signing method")
)

type Claims struct {
	Username string `json:"username"`
	UserID   int32  `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateTokens(username string, userId int32) (string, string, time.Time, error) {
	accessExp := time.Now().Add(4 * time.Hour)
	claims := Claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			Issuer:    issuer,
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

func ValidateToken(c fiber.Ctx, token string) (Claims, error) {
	claims := new(Claims)

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return SIGNING_KEY, nil
	})

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return *claims, ErrExpiredToken
	}

	if err != nil || !t.Valid {
		slog.ErrorContext(c.Context(), err.Error())
		return *claims, ErrInvalidToken
	}

	return *claims, nil
}
