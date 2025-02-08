package authmiddleware

import (
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
		SigningKey:  []byte(obj.JWTSecret),
		TokenLookup: "header:Authorization:Bearer ",
	})
}

func (obj *LocalMiddlewareMngr) GenerateJWT(userID int, AuthList []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"auth_list": AuthList,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(obj.JWTSecret))
}

func (obj *LocalMiddlewareMngr) GroupAuthorization(groups []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from the token
			// userID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64)

			// // Check if user belongs to the required group
			// var count int
			// err := db.QueryRow(`
			// 	SELECT COUNT(*)
			// 	FROM user_groups ug
			// 	JOIN groups g ON ug.group_id = g.id
			// 	WHERE ug.user_id = $1 AND g.name = $2`, int(userID), requiredGroup).Scan(&count)
			// if err != nil || count == 0 {
			// 	return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
			// }

			return next(c)
		}
	}
}
