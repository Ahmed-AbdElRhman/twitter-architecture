package sqlmngr

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

// ******* Selectrow executes a SELECT query and returns the resulting rows *******
func (obj *Postgres) Selectrow(sql string) (*sql.Rows, error) {
	row, err := obj.db.Query(sql)
	fmt.Println("Selectrow Query " + sql)
	if err != nil {
		return nil, err
	}
	return row, nil
}

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

func (obj *Postgres) CloseDB() error {
	err := obj.db.Close()
	if err != nil {
		log.Println("Error While Clossing Database:", err)
		return err
	}
	fmt.Println("Closed the database connection")
	return nil
}
