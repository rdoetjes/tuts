package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"phonax.com/db/metrics"
)

// MetricsHandler returns query metrics as JSON
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Let RespondJSON set the Content-Type and status code to avoid
	api := &RestAPI[interface{}]{}
	api.RespondJSON(w, http.StatusOK, metrics.GetMetrics())
}

// MetricsHandlerPrometheus returns the metrics JSON into prometheus format
func MetricsHandlerPrometheus(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(metrics.GetMetrics())
	if err != nil {
		fmt.Println("ERROR ", err)
	}

	prom, err := metrics.ConvertMetricsToPrometheus(data)
	if err != nil {
		fmt.Println("ERROR ", err)
	}

	// Let RespondJSON set the Content-Type and status code to avoid
	api := &RestAPI[interface{}]{}
	api.RespondJSON(w, http.StatusOK, prom)
}
