package auth

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v4"
)

const defaultExpiry = 24 * time.Hour

// JWTManager handles JWT operations.
type JWTManager interface {
	Generate(userID int64) (string, error)
	TokenAuth() *jwtauth.JWTAuth
}

type jwtManager struct {
	secret []byte
	expiry time.Duration
	auth   *jwtauth.JWTAuth
}

func NewJWTManager(secret string) JWTManager {
	a := jwtauth.New("HS256", []byte(secret), nil)
	return &jwtManager{secret: []byte(secret), expiry: defaultExpiry, auth: a}
}

func (j *jwtManager) Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(j.expiry).Unix()}
	_, tokenString, err := j.auth.Encode(claims)
	return tokenString, err
}

// JWTAuth returns the underlying jwtauth instance for middleware.
func (j *jwtManager) TokenAuth() *jwtauth.JWTAuth {
	return j.auth
}
