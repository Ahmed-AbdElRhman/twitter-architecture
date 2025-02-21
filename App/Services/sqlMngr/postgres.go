package sqlmngr

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/users/users_services"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connStr string, schemaFile string) (*Postgres, error) {
	//TODO:SingleTone Pattern

	// Open the connection to the database
	dbconn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Initialize the database
	if schemaFile != "" {
		if err := initializeDatabase(dbconn, schemaFile); err != nil {
			return nil, err
		}
	}
	return &Postgres{
		db: dbconn,
	}, nil
}

// ******* InitializeDatabase initializes the database with the schema in the given file *******
func initializeDatabase(db *sql.DB, schemaFile string) error {
	// Open the SQL file
	file, err := os.Open(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split the content into individual SQL statements
	queries := strings.Split(string(content), ";")

	// Execute each query
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue // Skip empty queries
		}
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query (%s): %w", query, err)
		}
	}

	return nil
}

// ******* CreateUser creates a new user in the database *******
func (obj *Postgres) CreateUser(userRegisterparam *users_services.UsrRegisterparam) (int, error) {
	var userid int
	// Insert the user into the database
	query := "INSERT INTO users(username, email, password) VALUES($1, $2, $3) RETURNING id"
	err := obj.db.QueryRow(query, userRegisterparam.Username, userRegisterparam.Email, userRegisterparam.Password).Scan(&userid)
	if err != nil {
		return -1, fmt.Errorf("failed to insert user into database: %w", err)
	}
	return userid, nil
}

// ******* GetUser retrieves a user from the database by username *******
func (obj *Postgres) GetUser(username string) (*users_services.UserInfo, error) {

	user := &users_services.UserInfo{}
	if obj.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT id, username, email, password FROM users WHERE username=$1"

	err := obj.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// ******* GetUserAuthlist retrieves the list of endpoints that a user has access to *******
func (obj *Postgres) GetUserAuthlist(user *users_services.UserInfo, userid int) error {
	queryEndpoints := "SELECT authendpoint FROM user_data_with_endpoints WHERE userid = $1"
	rows, err := obj.db.Query(queryEndpoints, user.ID)
	if err != nil {
		return fmt.Errorf("unable to select endpoint for User:%d- %w", user.ID, err)
	}
	defer rows.Close()

	var endpoints []string
	for rows.Next() {
		var endpoint string
		err := rows.Scan(&endpoint)
		if err == sql.ErrNoRows {
			// return fmt.Errorf("user Dont has any group")
		} else if err != nil {
			return fmt.Errorf("unable to fetch endpoint - %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	user.AuthList = endpoints
	return nil
}

// ******* Insert inserts a new record into the specified table *******
func (obj *Postgres) Insert(tableName string, argsKeys []string, argsVals []string) error {
	sql := fmt.Sprintf(`INSERT INTO %s(%s) VALUES ('%s')`, tableName, strings.Join(argsKeys, ","), strings.Join(argsVals, `','`))
	fmt.Println(sql)
	_, err := obj.db.Query(sql)
	if err != nil {
		log.Println("Error inserting into Database:", err)
		return err
	}
	return nil
}

// ******* CloseDB closes the database connection *******
func (obj *Postgres) CloseDB() error {
	err := obj.db.Close()
	if err != nil {
		log.Println("Error While Clossing Database:", err)
		return err
	}
	fmt.Println("Closed the database connection")
	return nil
}
