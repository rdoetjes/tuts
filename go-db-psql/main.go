package main

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "192.168.178.92"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func connectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

// Converts interface{} to string safely
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

func printRows(rows *sqlx.Rows) {
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			panic(err)
		}

		// MapScan returns []byte for text columns, so convert
		id := row["id"]
		fn := asString(row["first_name"])
		ln := asString(row["last_name"])

		fmt.Printf("name: %d %s %s\n", id, fn, ln)
	}
}

func connect() *sqlx.DB {
	db := sqlx.MustConnect("postgres", connectionString())
	return db
}

func query(db *sqlx.DB, sql string, wg *sync.WaitGroup) {
	defer wg.Done()

	rows, err := db.Queryx(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	printRows(rows)
}

func main() {
	fmt.Println("Connecting to:", connectionString())

	db := connect()
	defer db.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go query(db, "SELECT * FROM test WHERE first_name LIKE '%Keith%'", &wg)
	go query(db, "SELECT * FROM test WHERE last_name LIKE 'Doetjes'", &wg)

	wg.Wait()
}
