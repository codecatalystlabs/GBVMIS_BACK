package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))
var RefreshKey = []byte(os.Getenv("REFRESH_SECRET")) // Add this to your .env too

type Claims struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uint, email string, roles []string) (accessToken string, refreshToken string, err error) {
	// Access token - short lived (e.g., 15 min)
	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString(JwtKey)
	if err != nil {
		return
	}

	// Refresh token - longer lived (e.g., 7 days)
	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString(RefreshKey)
	return
}
