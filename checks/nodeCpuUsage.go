package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nagios-es/config"
	"net/http"

	"github.com/atc0005/go-nagios"
)

func CheckNodeCPUUsage(c *config.Config) *nagios.Plugin {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()

	resp, err := http.Get(fmt.Sprintf("%s/_nodes/stats/os", c.ElasticsearchURL))
	if err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to connect to Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			plugin.ServiceOutput = "CRITICAL: Failed to read response from Elasticsearch"
			plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			plugin.Errors = append(plugin.Errors, err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to read response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	var health ClusterNodesStatsResponse
	if err := json.Unmarshal(body, &health); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	var pd []nagios.PerformanceData
	var maxCPU int
	for _, node := range health.Nodes {
		if node.IP == c.NodeIP {
			nodeCpuPercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.OS.CPU.Percent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}

			if err := plugin.AddPerfData(false, nodeCpuPercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.OS.CPU.Percent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: CPU usage on node %s is %d%%", node.IP, node.OS.CPU.Percent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.OS.CPU.Percent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: CPU usage on node %s is %d%%", node.IP, node.OS.CPU.Percent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: CPU usage on node %s less then %d", node.IP, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}
			return plugin
		}

		if node.Name == c.NodeName {
			nodeCpuPercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.OS.CPU.Percent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}
			if err := plugin.AddPerfData(false, nodeCpuPercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.OS.CPU.Percent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: CPU usage on node %s is %d%%", node.Name, node.OS.CPU.Percent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.OS.CPU.Percent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: CPU usage on node %s is %d%%", node.Name, node.OS.CPU.Percent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: CPU usage on node %s less then %d", node.Name, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}

			return plugin
		}

		if node.OS.CPU.Percent > maxCPU {
			maxCPU = node.OS.CPU.Percent
		}

		nodeCpuPercent := nagios.PerformanceData{
			Label:             node.Name,
			Value:             fmt.Sprintf("%d", node.OS.CPU.Percent),
			Warn:              fmt.Sprintf("%d", c.WarningThreshold),
			Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
			Min:               "0",
			Max:               "100",
			UnitOfMeasurement: "%",
		}
		pd = append(pd, nodeCpuPercent)
	}

	if err := plugin.AddPerfData(false, pd...); err != nil {
		log.Printf("failed to add performance data metrics: %v\n", err)
		plugin.Errors = append(plugin.Errors, err)
	}

	switch {
	case maxCPU > c.CriticalThreshold:
		plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Max(CPU usage) on cluster is %d%%", maxCPU)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
	case maxCPU > c.WarningThreshold:
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Max(CPU usage) on cluster is %d%%", maxCPU)
		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
	default:
		plugin.ServiceOutput = fmt.Sprintf("OK: Max(CPU usage) on cluster less then %d", c.WarningThreshold)
		plugin.ExitStatusCode = nagios.StateOKExitCode
	}

	return plugin
}
