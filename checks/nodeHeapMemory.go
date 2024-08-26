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

func CheckNodeHeapMemory(c *config.Config) *nagios.Plugin {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()

	resp, err := http.Get(fmt.Sprintf("%s/_nodes/stats/jvm", c.ElasticsearchURL))
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

	var health ClusterNodesStatsResponse
	if err := json.Unmarshal(body, &health); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	var pd []nagios.PerformanceData
	var maxHeap int
	for _, node := range health.Nodes {
		if node.IP == c.NodeIP {
			nodeHeapPercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.JVM.Mem.HeapUsedPercent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}

			if err := plugin.AddPerfData(false, nodeHeapPercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.JVM.Mem.HeapUsedPercent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Heap size on node %s is %d%%", node.IP, node.JVM.Mem.HeapUsedPercent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.JVM.Mem.HeapUsedPercent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: Heap size on node %s is %d%%", node.IP, node.JVM.Mem.HeapUsedPercent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: Heap size on node %s less then %d", node.IP, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}

			return plugin
		}

		if node.Name == c.NodeName {
			nodeHeapPercent := nagios.PerformanceData{
				Label:             node.Name,
				Value:             fmt.Sprintf("%d", node.JVM.Mem.HeapUsedPercent),
				Warn:              fmt.Sprintf("%d", c.WarningThreshold),
				Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
				Min:               "0",
				Max:               "100",
				UnitOfMeasurement: "%",
			}

			if err := plugin.AddPerfData(false, nodeHeapPercent); err != nil {
				log.Printf("failed to add performance data metrics: %v\n", err)
				plugin.Errors = append(plugin.Errors, err)
			}

			switch {
			case node.JVM.Mem.HeapUsedPercent > c.CriticalThreshold:
				plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Heap size on node %s is %d%%", node.Name, node.JVM.Mem.HeapUsedPercent)
				plugin.ExitStatusCode = nagios.StateCRITICALExitCode
			case node.JVM.Mem.HeapUsedPercent > c.WarningThreshold:
				plugin.ServiceOutput = fmt.Sprintf("WARNING: Heap size on node %s is %d%%", node.Name, node.JVM.Mem.HeapUsedPercent)
				plugin.ExitStatusCode = nagios.StateWARNINGExitCode
			default:
				plugin.ServiceOutput = fmt.Sprintf("OK: Heap size on node %s less then %d", node.Name, c.WarningThreshold)
				plugin.ExitStatusCode = nagios.StateOKExitCode
			}

			return plugin
		}

		if node.JVM.Mem.HeapUsedPercent > maxHeap {
			maxHeap = node.JVM.Mem.HeapUsedPercent
		}

		nodeHeapPercent := nagios.PerformanceData{
			Label:             node.Name,
			Value:             fmt.Sprintf("%d", node.JVM.Mem.HeapUsedPercent),
			Warn:              fmt.Sprintf("%d", c.WarningThreshold),
			Crit:              fmt.Sprintf("%d", c.CriticalThreshold),
			Min:               "0",
			Max:               "100",
			UnitOfMeasurement: "%",
		}
		pd = append(pd, nodeHeapPercent)
	}

	if err := plugin.AddPerfData(false, pd...); err != nil {
		log.Printf("failed to add performance data metrics: %v\n", err)
		plugin.Errors = append(plugin.Errors, err)
	}

	switch {
	case maxHeap > c.CriticalThreshold:
		plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Max(Heap size) on cluster is %d%%", maxHeap)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
	case maxHeap > c.WarningThreshold:
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Max(Heap size) on cluster is %d%%", maxHeap)
		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
	default:
		plugin.ServiceOutput = fmt.Sprintf("OK: Max(Heap size) on cluster less then %d", c.WarningThreshold)
		plugin.ExitStatusCode = nagios.StateOKExitCode
	}

	return plugin
}
