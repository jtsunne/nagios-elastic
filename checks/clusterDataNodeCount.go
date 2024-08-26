package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"nagios-es/config"
	"net/http"

	"github.com/atc0005/go-nagios"
)

func CheckClusterDataNodeCount(c *config.Config) *nagios.Plugin {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()

	resp, err := http.Get(fmt.Sprintf("%s/_cluster/health", c.ElasticsearchURL))
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

	var health ClusterNodeCountResponse
	if err := json.Unmarshal(body, &health); err != nil {
		plugin.ServiceOutput = "CRITICAL: Failed to parse JSON response from Elasticsearch"
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
		return plugin
	}

	switch {
	case health.NumberOfDataNodes < c.CriticalThreshold:
		plugin.ServiceOutput = fmt.Sprintf("CRITICAL: Number of nodes is %d", health.NumberOfDataNodes)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode
	case health.NumberOfDataNodes < c.WarningThreshold:
		plugin.ServiceOutput = fmt.Sprintf("WARNING: Number of nodes is %d", health.NumberOfDataNodes)
		plugin.ExitStatusCode = nagios.StateWARNINGExitCode
	default:
		plugin.ServiceOutput = fmt.Sprintf("OK: Number of nodes is %d", health.NumberOfDataNodes)
		plugin.ExitStatusCode = nagios.StateOKExitCode
	}

	return plugin
}
