package middlewares

import (
	"net/http"
	"strings"
	"warungjwt_postgre2/config"
	"warungjwt_postgre2/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Middleware untuk otorisasi berdasarkan role pengguna
func IsAuthorized(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Ambil header Authorization dari request
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing authorization header"})
			}

			// Pisahkan kata 'Bearer' dari token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse token JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unexpected signing method"})
				}
				return config.JWTSecret, nil
			})

			// Jika terjadi error saat parsing token
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			// Periksa klaim token JWT
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userRole := claims["role"].(string)

				// Otorisasi berdasarkan role
				if role == "admin" && userRole != "admin" {
					return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden for non-admin"})
				}

				if role == "staff" && userRole != "staff" {
					return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden for non-staff"})
				}

				// Batasi akses staff hanya untuk metode GET
				if role == "staff" && c.Request().Method != http.MethodGet {
					return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden for staff"})
				}

				// Set token di konteks untuk digunakan di handler
				c.Set("user", token)
				return next(c)
			} else {
				// return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
				return utils.HandleError(c, utils.NewUnauthorizedError("Invalid token claims"))
			}
		}
	}
}
