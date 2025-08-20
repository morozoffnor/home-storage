package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret string
}

type Claims struct {
	jwt.RegisteredClaims
	UserEmail string `json:"user_email"`
}

func (j *JWT) GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)),
		},
		UserEmail: email,
	})
	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenStr, nil
}

func (j *JWT) ParseToken(token string) (*Claims, error) {
	token = token[7:]
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (j *JWT) CheckToken(r *http.Request) bool {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		return false
	}

	claims, err := j.ParseToken(cookie.Value)
	if err != nil {
		return false
	}
	if claims.UserEmail == "" {
		return false
	}
	return true
}

func (j *JWT) GetUserIdFromToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		return "", err
	}

	claims, err := j.ParseToken(cookie.Value)
	if err != nil {
		return "", err
	}
	if claims.UserEmail == "" {
		return "", err
	}
	return claims.UserEmail, nil
}

func (j *JWT) AddTokenToCookies(w *http.ResponseWriter, r *http.Request, token string) (context.Context, error) {
	http.SetCookie(*w, &http.Cookie{
		Name:    "Authorization",
		Value:   token,
		Expires: time.Now().Add(24 * 30 * time.Hour),
		Path:    "/",
	})

	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Context(), ContextUserEmail, claims.UserEmail)

	return ctx, nil
}
