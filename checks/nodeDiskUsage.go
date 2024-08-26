package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nagios-es/config"
	"nagios-es/helper"
	"net/http"

	"github.com/atc0005/go-nagios"
)

func CheckNodeDiskUsage(c *config.Config) *nagios.Plugin {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()

	resp, err := http.Get(fmt.Sprintf("%s/_nodes/stats/fs", c.ElasticsearchURL))
	if err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to connect to Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to read response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	var nodeStats ClusterNodesStatsResponse
	if err := json.Unmarshal(body, &nodeStats); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	var pd []nagios.PerformanceData
	var maxDiskUsage int
	for _, node := range nodeStats.Nodes {
		node.FS.Total.UsedPercent = helper.CalculateDiskUsagePercentage(node.FS.Total.TotalInBytes, node.FS.Total.FreeInBytes)
		if node.IP == c.NodeIP {
			nodeDiskUsagePercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.FS.Total.UsedPercent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}

			if err := plugin.AddPerfData(false, nodeDiskUsagePercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.FS.Total.UsedPercent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Disk usage on node %s is %d%%", node.IP, node.FS.Total.UsedPercent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.FS.Total.UsedPercent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: Disk usage on node %s is %d%%", node.IP, node.FS.Total.UsedPercent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: Disk usage on node %s less then %d", node.IP, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}

			return plugin
		}

		if node.Name == c.NodeName {
			nodeDiskUsagePercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.FS.Total.UsedPercent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}

			if err := plugin.AddPerfData(false, nodeDiskUsagePercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.FS.Total.UsedPercent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Disk usage on node %s is %d%%", node.Name, node.FS.Total.UsedPercent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.FS.Total.UsedPercent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: Disk usage on node %s is %d%%", node.Name, node.FS.Total.UsedPercent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: Disk uage on node %s less then %d", node.Name, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}

			return plugin
		}

		if node.FS.Total.UsedPercent > maxDiskUsage {
			maxDiskUsage = node.FS.Total.UsedPercent
		}

		nodeDiskUsagePercent := nagios.PerformanceData{
			Label:             node.Name,
			Value:             fmt.Sprintf("%d", node.FS.Total.UsedPercent),
			Warn:              fmt.Sprintf("%d", c.WarningThreshold),
			Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
			Min:               "0",
			Max:               "100",
			UnitOfMeasurement: "%",
		}
		pd = append(pd, nodeDiskUsagePercent)
	}

	if err := plugin.AddPerfData(false, pd...); err != nil {
		log.Printf("failed to add performance data metrics: %v\n", err)
		plugin.Errors = append(plugin.Errors, err)
	}

	switch {
	case maxDiskUsage > c.CriticalThreshold:
		plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Max(Disk usage) on cluster is %d%%", maxDiskUsage)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
	case maxDiskUsage > c.WarningThreshold:
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Max(Disk usage) on cluster is %d%%", maxDiskUsage)
		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
	default:
		plugin.ServiceOutput = fmt.Sprintf("OK: Max(Disk usage) on cluster less then %d", c.WarningThreshold)
		plugin.ExitStatusCode = nagios.StateOKExitCode
	}

	return plugin
}
