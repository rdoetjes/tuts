package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database configuration constants
const (
	host     = "192.168.178.92"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

// connectionString constructs the PostgreSQL connection string
func connectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

// Connect establishes a connection to the PostgreSQL database
func Connect() *sqlx.DB {
	db := sqlx.MustConnect("postgres", connectionString())
	return db
}

// Query executes a SQL query and processes the result set
func Query(db *sqlx.DB, sql string) ([]map[string]interface{}, error) {
	rows, err := db.Queryx(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			return nil, err
		}

		// MapScan returns []byte for text columns, convert them to string
		for k, v := range row {
			row[k] = asString(v)
		}
		results = append(results, row)
	}

	return results, nil
}

// asString safely converts an interface to string
func asString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
