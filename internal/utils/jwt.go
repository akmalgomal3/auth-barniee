package utils

import (
	"auth-barniee/internal/config"
	"auth-barniee/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID  `json:"user_id"`
	Email    string     `json:"email"`
	Role     string     `json:"role"`
	SchoolID *uuid.UUID `json:"school_id,omitempty"` // Add SchoolID to claims
	jwt.StandardClaims
}

func GenerateToken(user *models.User, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	if user.SchoolID != uuid.Nil { // Only add if user is associated with a school
		claims.SchoolID = &user.SchoolID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
