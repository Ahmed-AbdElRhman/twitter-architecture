package users_services

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type DbMngr interface {
	CloseDB() error
	Insert(tableName string, argsKeys []string, argsVals []string) error
	//TODO: Remove SQL Package Make it puplic
	Selectrow(sql string) (*sql.Rows, error)
}
type UserService struct {
	dbMngr DbMngr
}

func NewUsersService(dbMngr DbMngr) *UserService {
	return &UserService{
		dbMngr: dbMngr,
	}
}
func (obj *UserService) GetUser(UsrLogin UsrLoginparam) (*User, error) {
	usertags := getJSONTags(User{})
	querySt := fmt.Sprintf(`SELECT %s from %s WHERE userId='%s' and password='%s' `, strings.Join(usertags, ","), "users",
		UsrLogin.UserId, UsrLogin.Password)
	fmt.Println("GetUser Query " + querySt)
	row, err := obj.dbMngr.Selectrow(querySt)
	if err != nil {
		return nil, err
	}
	fmt.Println("GetUser Query Done")
	user := &User{}
	// Check if there's a row to scan
	if row.Next() {
		// Scan the values from the row into the Product struct
		err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Printf("GetUSer:Error scanning row values: %s", err)
			return nil, err
		}
		return user, nil
	} else {
		// No rows found
		log.Printf("GetUSer:No rows found")
		return nil, nil
	}
}

/************************************************************************************************************/
// Global functions to get the JSON tags of a struct
func getJSONTags(s interface{}) []string {
	var tags []string
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
