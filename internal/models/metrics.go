package models

type ContainerMetrics struct {
	ContainerID string  `json:"container_id"`
	Name        string  `json:"name"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemUsageMB  float64 `json:"memory_mb"`
	MemPercent  float64 `json:"memory_percent"`
	NetRxKB     float64 `json:"network_rx_kb"`
	NetTxKB     float64 `json:"network_tx_kb"`
}

type DockerStats struct {
	Name string `json:"name"`

	CPUStats struct {
		CPUUsage struct {
			TotalUsage  uint64   `json:"total_usage"`
			PercpuUsage []uint64 `json:"percpu_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`

	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`

	MemoryStats struct {
		Usage    uint64 `json:"usage"`
		Limit    uint64 `json:"limit"`
		MaxUsage uint64 `json:"max_usage"`
		Failcnt  uint64 `json:"failcnt"`
	} `json:"memory_stats"`

	Networks map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}
