package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateJWT(dni int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	secretKey := os.Getenv("SECRET_KEY")

	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(120 * time.Minute)
	claims["authorized"] = true
	claims["dni"] = dni

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token)(interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					_, err := w.Write([]byte("You're Unauthorized!"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err2 := w.Write([]byte("You're unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}
			}

			if token.Valid {
				endpointHandler(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("You're unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}