package users_api

import (
	"fmt"
	"net/http"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/users_services"
	"github.com/labstack/echo/v4"
)

type UserManager interface {
	GetUser(UsrLogin users_services.UsrLoginparam) (*users_services.User, error)
}

type JWTMiddleware interface {
	JWTMiddleware() echo.MiddlewareFunc
	GenerateJWT(userID int, AuthList []string) (string, error)
	GroupAuthorization(groups []string) echo.MiddlewareFunc
}

type UsersRouter struct {
	userMngr      UserManager
	jwtMiddleware JWTMiddleware
}

func NewUsersRouter(userMngr UserManager, jwtMiddleware JWTMiddleware) *UsersRouter {
	return &UsersRouter{
		userMngr:      userMngr,
		jwtMiddleware: jwtMiddleware,
	}
}

func (obj *UsersRouter) Login(c echo.Context) error {
	// Get the request body
	var UsrLoginparam users_services.UsrLoginparam
	err := c.Bind(&UsrLoginparam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	//------ Users Login Services Logic -------------
	user, err := obj.userMngr.GetUser(UsrLoginparam)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	fmt.Println("User", user)
	// Generate JWT
	// ToDo : Get the user groups from the database
	authList := []string{"admin"}
	token, err := obj.jwtMiddleware.GenerateJWT(user.ID, authList)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed GenerateJWTe: %w", err))
	}
	// Return the token
	fmt.Println("Token", token)
	return c.JSON(http.StatusOK, token)
}
func (obj *UsersRouter) GetUserTweets(c echo.Context) error {
	fmt.Println("GetUserTweets")
	return c.JSON(http.StatusOK, "GetUserTweets")
}
