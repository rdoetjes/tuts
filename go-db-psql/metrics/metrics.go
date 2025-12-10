package metrics

import (
	"sync"
	"time"
)

// QueryMetrics tracks query statistics
type QueryMetrics struct {
	mu               sync.RWMutex
	totalQueries     int64
	totalDuration    time.Duration
	queryDurations   []time.Duration
	operationMetrics map[string]*OperationMetric
}

// OperationMetric tracks metrics for a specific operation type
type OperationMetric struct {
	Count           int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
}

var globalMetrics = &QueryMetrics{
	operationMetrics: make(map[string]*OperationMetric),
}

// RecordQuery records a query execution
func RecordQuery(operation string, duration time.Duration) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.totalQueries++
	globalMetrics.totalDuration += duration
	globalMetrics.queryDurations = append(globalMetrics.queryDurations, duration)

	// Record operation-specific metrics
	if opMetric, exists := globalMetrics.operationMetrics[operation]; exists {
		opMetric.Count++
		opMetric.TotalDuration += duration
		opMetric.AverageDuration = opMetric.TotalDuration / time.Duration(opMetric.Count)

		if duration < opMetric.MinDuration || opMetric.MinDuration == 0 {
			opMetric.MinDuration = duration
		}
		if duration > opMetric.MaxDuration {
			opMetric.MaxDuration = duration
		}
	} else {
		globalMetrics.operationMetrics[operation] = &OperationMetric{
			Count:           1,
			TotalDuration:   duration,
			AverageDuration: duration,
			MinDuration:     duration,
			MaxDuration:     duration,
		}
	}
}

// GetMetrics returns current metrics
func GetMetrics() map[string]interface{} {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	var avgDuration time.Duration
	if globalMetrics.totalQueries > 0 {
		avgDuration = globalMetrics.totalDuration / time.Duration(globalMetrics.totalQueries)
	}

	// Convert operation metrics to map
	opMetrics := make(map[string]interface{})
	for op, metric := range globalMetrics.operationMetrics {
		opMetrics[op] = map[string]interface{}{
			"count":            metric.Count,
			"total_duration":   metric.TotalDuration.String(),
			"average_duration": metric.AverageDuration.String(),
			"min_duration":     metric.MinDuration.String(),
			"max_duration":     metric.MaxDuration.String(),
		}
	}

	return map[string]interface{}{
		"total_queries":     globalMetrics.totalQueries,
		"total_duration":    globalMetrics.totalDuration.String(),
		"average_duration":  avgDuration.String(),
		"operation_metrics": opMetrics,
	}
}

// ResetMetrics clears all metrics
func ResetMetrics() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.totalQueries = 0
	globalMetrics.totalDuration = 0
	globalMetrics.queryDurations = []time.Duration{}
	globalMetrics.operationMetrics = make(map[string]*OperationMetric)
}

// GetTotalQueries returns total number of queries executed
func GetTotalQueries() int64 {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()
	return globalMetrics.totalQueries
}

// GetAverageDuration returns average query duration
func GetAverageDuration() time.Duration {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	if globalMetrics.totalQueries == 0 {
		return 0
	}
	return globalMetrics.totalDuration / time.Duration(globalMetrics.totalQueries)
}

// GetOperationMetrics returns metrics for a specific operation
func GetOperationMetrics(operation string) *OperationMetric {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	if metric, exists := globalMetrics.operationMetrics[operation]; exists {
		// Return a copy to avoid race conditions
		m := *metric
		return &m
	}
	return nil
}
