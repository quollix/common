package api

import "time"

type LogRecord struct {
	ReadTimestamp time.Time      `json:"_time"`
	LogTimestamp  *time.Time     `json:"log_timestamp,omitempty"`
	Severity      string         `json:"severity"`
	Message       string         `json:"_msg"`
	Attributes    map[string]any `json:"attributes,omitempty"`
	Maintainer    string         `json:"maintainer"`
	App           string         `json:"app"`
	Service       string         `json:"service"`
	Stream        string         `json:"stream,omitempty"`
	Parsed        bool           `json:"parsed"`
}

type MetricSeries struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}
