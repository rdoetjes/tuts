package rest

import (
	"net/http"

	"phonax.com/db/metrics"
)

// MetricsHandler returns query metrics as JSON
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	api := &RestAPI[interface{}]{}
	api.RespondJSON(w, http.StatusOK, metrics.GetMetrics())
}
