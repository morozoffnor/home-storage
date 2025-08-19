package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/morozoffnor/home-storage/internal/config"
)

type ContextUserEmailKey string

var ContextUserEmail ContextUserEmailKey = "user_email"

type Auth struct {
	cfg *config.Config
	Jwt *JWT
}

type User struct {
	Id            string `json:"userID"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	Authenticated bool
}

func New(cfg *config.Config) *Auth {
	return &Auth{cfg: cfg, Jwt: &JWT{secret: cfg.JWTSecret}}
}

func (a *Auth) HashPassword(p string) string {
	h := sha256.New()
	h.Write([]byte(p))
	return hex.EncodeToString(h.Sum(nil))
}
