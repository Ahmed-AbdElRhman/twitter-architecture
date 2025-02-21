package authmiddleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type LocalMiddlewareMngr struct {
	JWTSecret string
	//User Caims
}

func NewLocalMiddlewareMngr(JWTSecret string) *LocalMiddlewareMngr {
	return &LocalMiddlewareMngr{JWTSecret}
}
func (obj *LocalMiddlewareMngr) JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(obj.JWTSecret),
		// TokenLookup: "header:Authorization:Bearer ",
		TokenLookup: "header:Authorization",
	})
}

func (obj *LocalMiddlewareMngr) GenerateJWT(userID int, username string, authList []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"username":  username,
		"auth_list": authList,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(obj.JWTSecret))
}
func (obj *LocalMiddlewareMngr) GroupAuthorization(groups []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ---------- Get the JWT token from the Authorization header ----------
			authHeader := c.Request().Header.Get("Authorization")
			fmt.Println(authHeader)
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}
			// ---------- Extract token part ----------
			tokenString := authHeader
			// For Bearer Schema
			// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			// if tokenString == authHeader {
			// 	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
			// }

			// ---------- Parse the token ----------
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(obj.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// ---------- Extract claims ----------
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
			}

			// ---------- Extract auth_list from claims ----------
			authListInterface, exists := claims["auth_list"]
			if !exists {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "no authorization list found"})
			}

			// ---------- Convert auth_list to []string ----------
			authList, ok := authListInterface.([]interface{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "invalid auth_list format"})
			}

			// ---------- Convert []interface{} to []string ----------
			var userGroups []string
			for _, v := range authList {
				if str, ok := v.(string); ok {
					userGroups = append(userGroups, str)
				}
			}
			// Check if the user belongs to any of the required groups
			for _, requiredGroup := range groups {
				for _, userGroup := range userGroups {
					if userGroup == requiredGroup {
						return next(c) // Authorized, continue request
					}
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
		}
	}
}

// func (obj *LocalMiddlewareMngr) GroupAuthorization(groups []string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user := c.Get("user").(*jwt.Token)
// 			claims := user.Claims.(jwt.MapClaims)
// 			authList := claims["auth_list"].([]interface{})

// 			for _, group := range groups {
// 				for _, auth := range authList {
// 					if group == auth {
// 						return next(c)
// 					}
// 				}
// 			}
// 			return echo.ErrUnauthorized
// 		}
// 	}
// }
