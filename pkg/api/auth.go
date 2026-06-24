package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("go-final-pablo-secret")

func hashPassword(pass string) string {
	h := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", h)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	claims := jwt.MapClaims{
		"hash": hashPassword(pass),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка формирования токена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
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
