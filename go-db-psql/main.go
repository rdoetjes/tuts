package main

import (
	"fmt"

	"phonax.com/db/db"
)

const (
	host     = "192.168.178.92"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

// Converts and prints result set rows
func printRows(rows []map[string]interface{}) {
	for _, row := range rows {
		id := row["id"]
		fn := row["first_name"]
		ln := row["last_name"]

		fmt.Printf("name: %v %v %v\n", id, fn, ln)
	}
}

func main() {
	fmt.Println("Connecting to the database")

	sql := db.Connect()
	defer sql.Close()

	res, err := db.Query(sql, "SELECT * FROM test WHERE first_name LIKE '%Keith%'")
	if err != nil {
		panic(err)
	}

	res2, err := db.Query(sql, "SELECT * FROM test WHERE last_name LIKE 'Doetjes'")
	if err != nil {
		panic(err)
	}

	printRows(res)
	printRows(res2)
}
