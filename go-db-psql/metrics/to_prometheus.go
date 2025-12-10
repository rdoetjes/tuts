package metrics

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func durToMilliSec(d time.Duration) float64 {
	return float64(d) / 1000
}

func ConvertMetricsToPrometheus(metrics []byte) (string, error) {
	var m QueryMetrics
	if err := json.Unmarshal(metrics, &m); err != nil {
		return "", fmt.Errorf("unmarshal metrics: %w", err)
	}
	var b strings.Builder

	// Global metrics
	fmt.Fprintf(&b, "api_total_duration_seconds %f\n", durToMilliSec(m.TotalDuration))
	fmt.Fprintf(&b, "api_total_errors %d\n", m.TotalErrors)
	fmt.Fprintf(&b, "api_total_queries %d\n", m.TotalQueries)
	fmt.Fprintf(&b, "api_average_duration_seconds %f\n", durToMilliSec(m.TotalDuration/time.Duration(m.TotalQueries)))
	fmt.Fprintf(&b, "api_uptime %f\n", time.Since(startTime).Seconds())
	fmt.Fprintf(&b, "api_qps %d\n", m.Qps)

	// Per-operation metrics
	for op, om := range m.OperationMetrics {
		opl := strings.ToLower(op)

		fmt.Fprintf(&b, "api_operation_duration_seconds_avg{operation=\"%s\"} %f\n", opl, durToMilliSec(om.AverageDuration))
		fmt.Fprintf(&b, "api_operation_duration_seconds_min{operation=\"%s\"} %f\n", opl, durToMilliSec(om.MinDuration))
		fmt.Fprintf(&b, "api_operation_duration_seconds_max{operation=\"%s\"} %f\n", opl, durToMilliSec(om.MaxDuration))
		fmt.Fprintf(&b, "api_operation_duration_seconds_total{operation=\"%s\"} %f\n", opl, durToMilliSec(om.TotalDuration))
		fmt.Fprintf(&b, "api_operation_count{operation=\"%s\"} %d\n", opl, om.Count)
		fmt.Fprintf(&b, "api_operation_errors{operation=\"%s\"} %d\n", opl, om.NrErrors)
	}

	return b.String(), nil
}
