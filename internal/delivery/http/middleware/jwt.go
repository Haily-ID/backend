package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haily-id/engine/internal/pkg/response"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.ErrUnauthorized
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
			}

			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int64(userID))
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if status, ok := claims["status"].(string); ok {
				c.Set("status", status)
			}

			return next(c)
		}
	}
}

