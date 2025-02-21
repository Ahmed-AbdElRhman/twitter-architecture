package users_api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/users_services"
	"github.com/labstack/echo/v4"
)

type UserManager interface {
	GetUser(UsrLogin *users_services.UsrLoginparam) (*users_services.UserInfo, error)
	CreateUser(UsrRegister *users_services.UsrRegisterparam) (int, error)
}

type JWTMiddleware interface {
	JWTMiddleware() echo.MiddlewareFunc
	GenerateJWT(userID int, username string, authList []string) (string, error)
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

// ----------------- Register -----------------
// TODO: Create a duplicate key validation for the Email
func (obj *UsersRouter) Register(c echo.Context) error {
	//------ Get the request body -------------
	var UsrRegisterparam users_services.UsrRegisterparam
	err := c.Bind(&UsrRegisterparam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, printerror("Failed to bind the request body", err.Error()))
	}
	// ------ Validate required fields ------
	if err := validateStruct(UsrRegisterparam); err != nil {
		return c.JSON(http.StatusBadRequest, printerror("Failed to validate the request body", err.Error()))
	}
	//------ Users Register Services Logic --
	id, err := obj.userMngr.CreateUser(&UsrRegisterparam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, printerror("Failed to create the user", err.Error()))
	}
	return c.JSON(http.StatusOK, "Register:"+fmt.Sprint(id))
}

// ----------------- Login -----------------
func (obj *UsersRouter) Login(c echo.Context) error {
	//------ Get the request body -------------
	var UsrLoginparam users_services.UsrLoginparam
	err := c.Bind(&UsrLoginparam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, printerror("Failed to bind the request body", err.Error()))
	}

	//------ Validate required fields ------
	if err := validateStruct(UsrLoginparam); err != nil {
		return c.JSON(http.StatusBadRequest, printerror("Failed to validate the request body", err.Error()))
	}
	//------ Users Login Services Logic ------
	user, err := obj.userMngr.GetUser(&UsrLoginparam)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, printerror("Failed to get the user", err.Error()))
	}
	//------ Generate JWT -------------
	token, err := obj.jwtMiddleware.GenerateJWT(user.ID, user.Username, user.AuthList)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, printerror("Failed to generate the JWT", err.Error()))
	}
	//------ Return the token ----------
	return c.JSON(http.StatusOK, token)
}

// ----------------- GetUserTweets -----------------
func (obj *UsersRouter) GetUserTweets(c echo.Context) error {
	fmt.Println("GetUserTweets")
	return c.JSON(http.StatusOK, "GetUserTweets")
}

// ### Print Error
func printerror(msg string, errormsg string) map[string]string {
	// Simulate an error
	return map[string]string{
		"error": fmt.Sprintf("%s: %s", msg, errormsg),
	}
}

// ### Validate Struct
func validateStruct(input interface{}) error {
	val := reflect.ValueOf(input)

	// Ensure we are working with a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input is not a struct")
	}

	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is a string and empty
		if field.Kind() == reflect.String && field.String() == "" {
			// fieldName := reflect.TypeOf(input).Field(i).Name
			fieldName := reflect.TypeOf(input).Field(i).Tag.Get("json")
			return fmt.Errorf("field '%s' is required", fieldName)
		}
	}
	return nil
}
