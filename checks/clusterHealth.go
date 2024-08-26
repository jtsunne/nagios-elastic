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

type ClusterHealthResponse struct {
	Status           string `json:"status"`
	ActiveShards     int    `json:"active_shards"`
	RelocatingShards int    `json:"relocating_shards"`
	UnassignedShards int    `json:"unassigned_shards"`
}

func CheckClusterHealth(c *config.Config) *nagios.Plugin {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()

	resp, err := http.Get(fmt.Sprintf("%s/_cluster/health", c.ElasticsearchURL))
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

	var health ClusterHealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	switch health.Status {
	case "green":
		plugin.ServiceOutput = "OK: Cluster health is green"
		plugin.ExitStatusCode = nagios.StateOKExitCode
	case "yellow":
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Cluster health is yellow, relocating shards: %d", health.RelocatingShards)
		pd := []nagios.PerformanceData{
			{Label: "relocating_shards", Value: fmt.Sprintf("%d", health.RelocatingShards)},
			{Label: "unassigned_shards", Value: fmt.Sprintf("%d", health.UnassignedShards)},
			{Label: "active_shards", Value: fmt.Sprintf("%d", health.ActiveShards)},
		}

		if err := plugin.AddPerfData(false, pd...); err != nil {
			log.Printf("failed to add performance data metrics: %v\n", err)
			plugin.Errors = append(plugin.Errors, err)
		}

		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
	case "red":
		plugin.ServiceOutput = "CRITICAL: Cluster health is red"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
	default:
		plugin.ServiceOutput = fmt.Sprintf("UNKNOWN: Cluster health is %s", health.Status)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode
	}

	return plugin
}
