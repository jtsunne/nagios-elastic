package checks

import (
	"encoding/json"
	"fmt"
	"github.com/atc0005/go-nagios"
	"io"
	"nagios-es/config"
	"net/http"
)

type ClusterNodeCountResponse struct {
	NumberOfNodes     int `json:"number_of_nodes"`
	NumberOfDataNodes int `json:"number_of_data_nodes"`
}

func CheckClusterNodeCount(c *config.Config) *nagios.Plugin {
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

	var health ClusterNodeCountResponse
	if err := json.Unmarshal(body, &health); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	if health.NumberOfNodes < c.CriticalThreshold {
		plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Number of nodes is %d", health.NumberOfNodes)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}
	if health.NumberOfNodes < c.WarningThreshold {
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Number of nodes is %d", health.NumberOfNodes)
		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
		return plugin
	}

	plugin.ServiceOutput = fmt.Sprintf("OK: Number of nodes is %d", health.NumberOfNodes)
	plugin.ExitStatusCode = nagios.StateOKExitCode
	return plugin
}