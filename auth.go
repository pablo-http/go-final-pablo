package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret — секрет для подписи жвт
var jwtSecret = []byte("go-final-pablo-secret")

// hashPassword возвращает SHA-256 хэш пароля
func hashPassword(pass string) string {
	h := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", h)
}

// signinHandler обрабатывает POST /api/signin.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "ошибка десериализации JSON")
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeError(w, "неверный пароль")
		return
	}

	// создание жвт с хэшем пароля в качестве claim
	claims := jwt.MapClaims{
		"hash": hashPassword(pass),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeError(w, "ошибка формирования токена")
		return
	}

	writeJSON(w, map[string]string{"token": signed})
}

// auth  middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			// если пароль не задан, то авторизация не требуется
			next(w, r)
			return
		}

		var tokenStr string
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenStr = cookie.Value
		}

		valid := false
		if tokenStr != "" {
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неожиданный метод подписи")
				}
				return jwtSecret, nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					// проверка на хэш в токене совпадает с хэшем текущего пароля
					if claims["hash"] == hashPassword(pass) {
						valid = true
					}
				}
			}
		}

		if !valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
