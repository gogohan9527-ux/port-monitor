package main

import (
	"crypto/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	hardcodedUser = "admin"
	hardcodedPass = "admin"
	tokenTTL      = 8 * time.Hour
)

type Auth struct {
	secret []byte
}

func NewAuth() (*Auth, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return &Auth{secret: b}, nil
}

func (a *Auth) Issue(username string) (string, time.Time, error) {
	exp := time.Now().Add(tokenTTL)
	claims := jwt.MapClaims{
		"sub": username,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(a.secret)
	return s, exp, err
}

func (a *Auth) Verify(token string) bool {
	if token == "" {
		return false
	}
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return a.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !parsed.Valid {
		return false
	}
	return true
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearer(r.Header.Get("Authorization"))
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if !a.Verify(token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearer(h string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(h, prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	return ""
}
