package metrics

import (
	"sync"
	"time"
)

// QueryMetrics tracks query statistics
type QueryMetrics struct {
	mu               sync.RWMutex
	TotalQueries     uint64        `json:"total_queries"`
	TotalDuration    time.Duration `json:"total_duration"`
	TotalErrors      uint64        `json:"total_errors"`
	queryDurations   []time.Duration
	OperationMetrics map[string]*OperationMetric `json:"operation_metrics"`
}

// OperationMetric tracks metrics for a specific operation type
type OperationMetric struct {
	Count           uint64        `json:"count"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`
	NrErrors        uint64        `json:"nr_errors"`
}

var globalMetrics = &QueryMetrics{
	OperationMetrics: make(map[string]*OperationMetric),
}

// RecordQuery records a query execution
func RecordQuery(operation string, duration time.Duration, err bool) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.TotalQueries++
	globalMetrics.TotalDuration += duration
	if err {
		globalMetrics.TotalErrors++
	}

	globalMetrics.queryDurations = append(globalMetrics.queryDurations, duration)

	// Record operation-specific metrics
	if opMetric, exists := globalMetrics.OperationMetrics[operation]; exists {
		opMetric.Count++
		opMetric.TotalDuration += duration
		opMetric.AverageDuration = opMetric.TotalDuration / time.Duration(opMetric.Count)
		if err {
			opMetric.NrErrors++
		}

		if duration < opMetric.MinDuration || opMetric.MinDuration == 0 {
			opMetric.MinDuration = duration
		}
		if duration > opMetric.MaxDuration {
			opMetric.MaxDuration = duration
		}
	} else {
		globalMetrics.OperationMetrics[operation] = &OperationMetric{
			Count:           1,
			TotalDuration:   duration,
			AverageDuration: duration,
			MinDuration:     duration,
			MaxDuration:     duration,
			NrErrors:        0,
		}
	}
}

// GetMetrics returns current metrics
func GetMetrics() map[string]interface{} {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	var avgDuration time.Duration
	if globalMetrics.TotalQueries > 0 {
		avgDuration = globalMetrics.TotalDuration / time.Duration(globalMetrics.TotalQueries)
	}

	// Convert operation metrics to map
	opMetrics := make(map[string]interface{})
	for op, metric := range globalMetrics.OperationMetrics {
		opMetrics[op] = map[string]interface{}{
			"count": metric.Count,
			// CHANGE: Remove .String() to keep them as time.Duration (int64 nanoseconds)
			"total_duration":   metric.TotalDuration,
			"average_duration": metric.AverageDuration,
			"min_duration":     metric.MinDuration,
			"max_duration":     metric.MaxDuration,
			"nr_errors":        metric.NrErrors,
		}
	}

	return map[string]interface{}{
		"total_queries": globalMetrics.TotalQueries,
		// CHANGE: Remove .String() here as well
		"total_duration":    globalMetrics.TotalDuration,
		"average_duration":  avgDuration,
		"total_errors":      globalMetrics.TotalErrors,
		"operation_metrics": opMetrics,
	}
}

// ResetMetrics clears all metrics
func ResetMetrics() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.TotalQueries = 0
	globalMetrics.TotalDuration = 0
	globalMetrics.queryDurations = []time.Duration{}
	globalMetrics.OperationMetrics = make(map[string]*OperationMetric)
}

// GetTotalQueries returns total number of queries executed
func GetTotalQueries() uint64 {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()
	return globalMetrics.TotalQueries
}

// GetAverageDuration returns average query duration
func GetAverageDuration() time.Duration {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	if globalMetrics.TotalQueries == 0 {
		return 0
	}
	return globalMetrics.TotalDuration / time.Duration(globalMetrics.TotalQueries)
}

// GetOperationMetrics returns metrics for a specific operation
func GetOperationMetrics(operation string) *OperationMetric {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	if metric, exists := globalMetrics.OperationMetrics[operation]; exists {
		// Return a copy to avoid race conditions
		m := *metric
		return &m
	}
	return nil
}
