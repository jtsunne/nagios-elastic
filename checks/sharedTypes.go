package checks

type ClusterNodesStatsResponse struct {
	Nodes map[string]NodeStats `json:"nodes"`
}

type NodeStats struct {
	Name string   `json:"name"`
	IP   string   `json:"host"`
	OS   OSStats  `json:"os"`
	JVM  JVMStats `json:"jvm"`
	FS   FSStats  `json:"fs"`
}

// JVMStats represents the JVM-related statistics for a node.
type JVMStats struct {
	Mem MemStats `json:"mem"`
}

// MemStats represents the memory-related statistics for JVM.
type MemStats struct {
	HeapUsedPercent int `json:"heap_used_percent"`
}

type OSStats struct {
	CPU CPUStats `json:"cpu"`
}

// CPUStats represents the CPU-related statistics for a node.
type CPUStats struct {
	Percent int `json:"percent"`
}

// NodeFSStats represents the filesystem statistics for a specific node in the Elasticsearch cluster.
type NodeFSStats struct {
	FS FSStats `json:"fs"`
}

// FSStats represents the filesystem-related statistics for a node.
type FSStats struct {
	Total TotalStats `json:"total"`
}

// TotalStats represents the total, free, and available disk space.
type TotalStats struct {
	TotalInBytes     int64 `json:"total_in_bytes"`
	FreeInBytes      int64 `json:"free_in_bytes"`
	AvailableInBytes int64 `json:"available_in_bytes"`
	UsedPercent      int
}
