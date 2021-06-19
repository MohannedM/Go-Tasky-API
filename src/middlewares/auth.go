package middlewares

import (
	"TaskyBE/src/controllers"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		bearerToken := req.Header.Get("Authorization")
		if bearerToken == "" {
			controllers.ErrorThrower(http.StatusUnauthorized, "No bearer token", res)
			return
		}
		strArr := strings.Split(bearerToken, " ")
		tokenString := strArr[1]
		if tokenString == "" {
			controllers.ErrorThrower(http.StatusUnauthorized, "Errors with token", res)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey, _ := os.LookupEnv("JWT_SECRET")
			return []byte(secretKey), nil
		})
		if err != nil {
			controllers.ErrorThrower(http.StatusUnauthorized, "Expired token", res)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			controllers.ErrorThrower(http.StatusUnauthorized, "Error with token data", res)
			return
		}
		context.Set(req, "user_id", claims["user_id"])
		next.ServeHTTP(res, req)
	})
}
