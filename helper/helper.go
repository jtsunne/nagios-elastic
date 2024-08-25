package helper

import (
	"fmt"

	"github.com/atc0005/go-nagios"
)

func ErrorUnknown(error string) {
	plugin := nagios.NewPlugin()
	defer plugin.ReturnCheckResults()
	plugin.ServiceOutput = fmt.Sprintf("UNKNOWN: %s", error)
	plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode
	plugin.ReturnCheckResults()
}

func CalculateDiskUsagePercentage(total, free int64) int {
	return int(100 * (total - free) / total)
}
