package api

type ObservabilitySettings struct {
	LogsRetentionDays    string `json:"logs_retention_days" validate:"number"`
	LogsMaxDiskSpaceGB   string `json:"logs_max_disk_space_gb" validate:"number"`
	MetricsRetentionDays string `json:"metrics_retention_days" validate:"number"`
}
