package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"phonax.com/db/auth"
	"phonax.com/db/db"
	"phonax.com/db/db/models"
	"phonax.com/db/rest"
)

func main() {
	const api_port = ":3000"

	fmt.Println("Connecting to the database")

	sql := db.Connect()
	defer sql.Close()

	// Set up database credential validator for JWT login
	validator := auth.NewDBCredentialValidator(sql)
	auth.SetCredentialValidator(validator)

	r := chi.NewRouter()

	// Public routes
	r.Post("/auth/login", auth.LoginHandler)
	r.Get("/metrics", rest.MetricsHandler)

	// Protected routes
	r.Group(func(rt chi.Router) {
		rt.Use(auth.JWTMiddleware)
		user := models.User{}
		restUser := rest.NewCrudAPI[models.User](sql, user.SetupQueries())
		restUser.RegisterRoutes(rt, "/users")
	})

	fmt.Println("Starting REST API on port ", api_port)
	err := http.ListenAndServe(api_port, r)
	if err != nil {
		log.Println(err)
	}
}
