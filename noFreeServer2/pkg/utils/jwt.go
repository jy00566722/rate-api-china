package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTUtil interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	GenerateTemporaryToken(userID uint) (string, error)
	ValidateTemporaryToken(tokenString string) (*JWTClaims, error)
}

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

type jwtUtil struct {
	secretKey         string
	tokenDuration     time.Duration
	tempTokenDuration time.Duration
}

func NewJWTUtil(secretKey string, tokenDuration time.Duration, tempTokenDuration time.Duration) JWTUtil {
	return &jwtUtil{
		secretKey:         secretKey,
		tokenDuration:     tokenDuration,
		tempTokenDuration: tempTokenDuration,
	}
}

func (j *jwtUtil) GenerateToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func (j *jwtUtil) GenerateTemporaryToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.tempTokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtUtil) ValidateTemporaryToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
