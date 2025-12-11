package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"phonax.com/db/auth"
	"phonax.com/db/db"
	"phonax.com/db/db/models"
	"phonax.com/db/rest"
)

func setupAuth(sql *sqlx.DB) {
	// Set up database credential validator for JWT login
	if os.Getenv("JWT_SECRET") != "" {
		log.Println("Using provided JWT secret")
		auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	}
	validator := auth.NewDBCredentialValidator(sql)
	auth.SetCredentialValidator(validator)
}

func setupRoutes(sql *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	// Public routes
	r.Post("/auth/login", auth.LoginHandler)
	r.Get("/v1/metrics", rest.MetricsHandler)
	r.Get("/v1/metrics_prometheus", rest.MetricsHandlerPrometheus)

	// Protected routes
	r.Group(func(rt chi.Router) {
		rt.Use(auth.JWTMiddleware)

		user := models.User{}
		restUser := rest.NewCrudAPI[models.User](sql, user.SetupQueries())
		restUser.RegisterRoutes(rt, "/v1/users")
	})
	return r
}

func main() {
	const api_port = ":3000"

	fmt.Println("Connecting to the database")

	sql := db.Connect()
	defer sql.Close()

	setupAuth(sql)

	r := setupRoutes(sql)

	fmt.Println("Starting REST API on port ", api_port)
	err := http.ListenAndServe(api_port, r)
	if err != nil {
		log.Println(err)
	}
}
