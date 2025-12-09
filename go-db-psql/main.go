package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"phonax.com/db/db"
	"phonax.com/db/models"
	"phonax.com/db/rest"
)

const (
	host     = "192.168.178.92"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func userQueries() rest.Queries {
	return rest.Queries{
		Create: `INSERT INTO users (firstname, lastname, dob, email) VALUES ($1, $2, $3, $4) RETURNING id`,
		GetOne: `SELECT id, firstname, lastname, dob, email FROM users WHERE id = $1`,
		List:   `SELECT id, firstname, lastname, dob, email FROM users`,
		Update: `UPDATE users SET firstname=$2, lastname=$3, dob=$4, email=$5 WHERE id = $1`,
		Delete: `DELETE FROM users WHERE id = $1`,
	}
}

func main() {
	const api_port = ":3000"

	fmt.Println("Connecting to the database")

	sql := db.Connect()
	defer sql.Close()

	r := chi.NewRouter()

	restUser := rest.NewCrudAPI[models.User](sql, userQueries())
	restUser.RegisterRoutes(r, "/users")

	fmt.Println("Starting REST API on port ", api_port)
	err := http.ListenAndServe(api_port, r)
	if err != nil {
		log.Println(err)
	}
}
