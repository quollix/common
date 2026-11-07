package api

type AlertSettings struct {
	CPUUsageEnabled          bool `json:"cpu_usage_enabled"`
	MemoryUsageEnabled       bool `json:"memory_usage_enabled"`
	DiskUsageEnabled         bool `json:"disk_usage_enabled"`
	CertificateExpiryEnabled bool `json:"certificate_expiry_enabled"`
	QuollixErrorsEnabled     bool `json:"quollix_errors_enabled"`
}
