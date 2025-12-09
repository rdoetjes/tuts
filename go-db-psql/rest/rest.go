package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
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

func NewRestAPI[T any](database *sqlx.DB, queries Queries) *RestAPI[T] {
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

	// Create expects you to return ID; so obj must contain ID after creation
	id, err := api.CRUD.Create(r.Context(), api.Queries.Create, structToArgs(obj)...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	// TODO: implement reflect-based struct → []any mapping
	return []any{} // placeholder
}
