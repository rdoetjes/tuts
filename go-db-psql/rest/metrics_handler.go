package rest

import (
	"encoding/json"
	"net/http"

	"phonax.com/db/metrics"
)

// MetricsHandler returns query metrics as JSON
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Let RespondJSON set the Content-Type and status code to avoid
	api := &RestAPI[interface{}]{}
	api.RespondJSON(w, http.StatusOK, metrics.GetMetrics())
}

func MetricsHandlerPrometheus(w http.ResponseWriter, r *http.Request) {
	// 1. Get metrics
	data, err := json.Marshal(metrics.GetMetrics())
	if err != nil {
		http.Error(w, "Failed to marshal metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Convert to Prometheus format
	prom, err := metrics.ConvertMetricsToPrometheus(data)
	if err != nil {
		http.Error(w, "Failed to convert to prometheus: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Send raw text response
	// Prometheus expects text/plain, conventionally version 0.0.4
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	// Convert string to byte array and Write the data directly.
	w.Write([]byte(prom))
}
