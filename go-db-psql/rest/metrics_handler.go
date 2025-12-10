package rest

import (
	"net/http"

	"phonax.com/db/metrics"
)

// MetricsHandler returns query metrics as JSON
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Let RespondJSON set the Content-Type and status code to avoid
	api := &RestAPI[interface{}]{}
	api.RespondJSON(w, http.StatusOK, metrics.GetMetrics())
}
