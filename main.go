package main

import (
	"nagios-es/checks"
	"nagios-es/config"
	"nagios-es/helper"

	"github.com/atc0005/go-nagios"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		helper.ErrorUnknown("Can't load config")
	}

	if cfg.ElasticsearchURL == "" {
		helper.ErrorUnknown("Elasticsearch URL is required")
	}

	if cfg.Check == "" {
		helper.ErrorUnknown("Check name is required")
	}

	var plugin *nagios.Plugin

	switch cfg.Check {
	case "health":
		plugin = checks.CheckClusterHealth(cfg)
	case "node_count":
		plugin = checks.CheckClusterNodeCount(cfg)
	case "data_node_count":
		plugin = checks.CheckClusterDataNodeCount(cfg)
	case "cpu_usage":
		plugin = checks.CheckNodeCPUUsage(cfg)
	case "heap_size":
		plugin = checks.CheckNodeHeapMemory(cfg)
	case "disk_usage":
		plugin = checks.CheckNodeDiskUsage(cfg)
	default:
		helper.ErrorUnknown(cfg.Check)
	}

	if plugin != nil {
		plugin.ReturnCheckResults()
	}
}
