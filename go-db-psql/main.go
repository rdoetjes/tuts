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

func main() {
	const api_port = ":3000"

	fmt.Println("Connecting to the database")

	sql := db.Connect()
	defer sql.Close()

	r := chi.NewRouter()

	user := models.User{}
	restUser := rest.NewCrudAPI[models.User](sql, user.SetupQueries())
	restUser.RegisterRoutes(r, "/users")

	fmt.Println("Starting REST API on port ", api_port)
	err := http.ListenAndServe(api_port, r)
	if err != nil {
		log.Println(err)
	}
}
