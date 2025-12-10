package rest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"phonax.com/db/auth"
	"phonax.com/db/db"
)

type Queries struct {
	Create string
	GetOne string
	List   string
	Update string
	Delete string
}

type RestAPI[T any] struct {
	CRUD    *db.CRUD[T]
	Queries Queries
}

func NewCrudAPI[T any](database *sqlx.DB, queries Queries) *RestAPI[T] {
	return &RestAPI[T]{
		CRUD:    &db.CRUD[T]{DB: database},
		Queries: queries,
	}
}

func (api *RestAPI[T]) RegisterRoutes(r chi.Router, basePath string) {
	r.Route(basePath, func(rt chi.Router) {
		rt.Post("/", api.create)
		rt.Get("/", api.list)
		rt.Get("/{id}", api.getOne)
		rt.Put("/{id}", api.update)
		rt.Delete("/{id}", api.delete)
	})
}

// ---- Handlers ----

func (api *RestAPI[T]) create(w http.ResponseWriter, r *http.Request) {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password field if it exists
	hashPasswordField(&obj)

	// Create expects you to return ID; so obj must contain ID after creation
	id, err := api.CRUD.Create(r.Context(), api.Queries.Create, structToArgs(obj)...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with 201 Created for a successful creation
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]any{"id": id})
}

func (api *RestAPI[T]) getOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	obj, err := api.CRUD.GetOne(r.Context(), api.Queries.GetOne, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, obj)
}

func (api *RestAPI[T]) list(w http.ResponseWriter, r *http.Request) {
	obj, err := api.CRUD.List(r.Context(), api.Queries.List)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, obj)
}

func (api *RestAPI[T]) update(w http.ResponseWriter, r *http.Request) {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password field if it exists
	hashPasswordField(&obj)

	id := chi.URLParam(r, "id")

	_, err := api.CRUD.Update(r.Context(), api.Queries.Update, append(structToArgs(obj), id)...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "updated"})
}

func (api *RestAPI[T]) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := api.CRUD.Delete(r.Context(), api.Queries.Delete, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "deleted"})
}

// ---- Helpers ----

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

/*
structToArgs converts struct fields to []any for SQL args.
You can replace this with your own logic or a reflection helper.
*/
func structToArgs[T any](obj T) []any {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	// If the struct is a pointer, dereference it
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return []any{}
	}

	args := []any{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Read db tag
		dbTag := field.Tag.Get("db")

		// Skip db:"-" fields
		if dbTag == "-" {
			continue
		}

		// If no db tag, skip — SQLX requires db tags for mapping
		if dbTag == "" {
			continue
		}

		// Skip "id" field on create when zero
		if strings.EqualFold(dbTag, "id") && isZeroValue(value) {
			continue
		}

		args = append(args, value.Interface())
	}

	return args
}

func isZeroValue(v reflect.Value) bool {
	// Invalid or nil pointer?
	if !v.IsValid() {
		return true
	}

	// For pointer values
	if v.Kind() == reflect.Pointer {
		return v.IsNil()
	}

	// Compare to zero value
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// hashPasswordField finds a "password" field in the struct and hashes it using bcrypt
func hashPasswordField[T any](obj *T) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		return
	}

	v = v.Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return
	}

	// Look for a field with db tag "password" or json tag "password"
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Check if this is a password field
		dbTag := field.Tag.Get("db")
		if dbTag != "password" && field.Name != "Password" {
			continue
		}

		// Only process string fields
		if fieldValue.Kind() != reflect.String {
			continue
		}

		// Skip if password is empty
		password := fieldValue.String()
		if password == "" {
			continue
		}

		// Hash the password
		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			// Log error but continue (could add error handling here if needed)
			continue
		}

		// Update the field with the hashed password
		fieldValue.SetString(hashedPassword)
	}
}
