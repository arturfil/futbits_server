package middlewares

import (
	"chi_soccer/helpers"
	"chi_soccer/models"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func IsAuthorized(next http.Handler) http.Handler {
	godotenv.Load(".env")
	var myKey = []byte(os.Getenv("SECRET_KEY"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return myKey, nil
			})
			if err != nil {
				payload := models.JsonResponse{
					Error:   true,
					Message: err.Error(),
				}
				_ = helpers.WriteJSON(w, http.StatusUnauthorized, payload)
				return
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			helpers.ErrorJSON(w, errors.New("authorization headers missing"))
		}
	})
}
