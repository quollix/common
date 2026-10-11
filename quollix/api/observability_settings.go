package api

type ObservabilitySettings struct {
	LogsRetentionDays               string `json:"logs_retention_days" validate:"number"`
	LogsMaxDiskSpaceGB              string `json:"logs_max_disk_space_gb" validate:"number"`
	MetricsRetentionDays            string `json:"metrics_retention_days" validate:"number"`
	HostMetricsIntervalSeconds      string `json:"host_metrics_interval_seconds" validate:"number"`
	ContainerMetricsIntervalSeconds string `json:"container_metrics_interval_seconds" validate:"number"`
	VolumeMetricsIntervalSeconds    string `json:"volume_metrics_interval_seconds" validate:"number"`
}
