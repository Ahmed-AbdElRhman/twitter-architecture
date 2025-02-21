package users_services

import (
	"fmt"
	"reflect"
	"strings"
)

type DbMngr interface {
	CloseDB() error
	GetUser(username string) (*UserInfo, error)
	CreateUser(UsrRegister *UsrRegisterparam) (int, error)
	Insert(tableName string, argsKeys []string, argsVals []string) error
	GetUserAuthlist(user *UserInfo, userid int) error
}
type HashMngr interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashedPassword string) error
}
type UserService struct {
	dbMngr   DbMngr
	hashMngr HashMngr
}

func NewUsersService(dbMngr DbMngr, hashMngr HashMngr) *UserService {
	return &UserService{
		dbMngr:   dbMngr,
		hashMngr: hashMngr,
	}
}

// ----------------- CreateUser -----------------
func (obj *UserService) CreateUser(usrRegisterparam *UsrRegisterparam) (int, error) {
	// ----- Hash the password -----
	hashedPassword, err := obj.hashMngr.HashPassword(usrRegisterparam.Password)
	if err != nil {
		return -1, err
	}
	usrRegisterparam.Password = hashedPassword
	// ----- Insert the user into the database -----
	id, err := obj.dbMngr.CreateUser(usrRegisterparam)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// ----------------- GetUser -----------------
func (obj *UserService) GetUser(UsrLogin *UsrLoginparam) (*UserInfo, error) {
	// ----- Check if the user exists in the database -----
	user, err := obj.dbMngr.GetUser(UsrLogin.Username)
	if err != nil {
		return nil, err
	}
	// ----- Check if the password is correct -----
	if err_ := obj.hashMngr.CheckPassword(UsrLogin.Password, user.Password); err_ != nil {
		return nil, err_

	}
	// ----- Get the user's auth list -----
	if err_ := obj.dbMngr.GetUserAuthlist(user, user.ID); err_ != nil {
		return nil, err_

	}
	if err_ := len(user.AuthList) == 0; err_ {
		fmt.Println("user Dont has any group")

	}
	return user, nil
}

/************************************************************************************************************/
// Global functions to get the JSON tags of a struct
func getJSONTagsandvalues(s interface{}) []string {
	var tags []string
	// var value []string
	// Get the type of the struct
	t := reflect.TypeOf(s)

	// Iterate over all fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Get the "json" tag value
		tag := field.Tag.Get("json")
		// Split to handle cases like `json:"name,omitempty"`
		parts := strings.Split(tag, ",")
		// Take the first part (the actual JSON key)
		if len(parts) > 0 && parts[0] != "" {
			tags = append(tags, parts[0])
		}
	}
	return tags
}

// func getStructKeys(v interface{}) []string {
// 	val := reflect.ValueOf(v)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	if val.Kind() != reflect.Struct {
// 		return nil
// 	}

// 	var keys []string
// 	for i := 0; i < val.NumField(); i++ {
// 		keys = append(keys, val.Type().Field(i).Name)
// 	}
// 	return keys
// }
