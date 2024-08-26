package helper

import (
	"github.com/atc0005/go-nagios"
)

func ErrorUnknown(errString string) {
	plugin := nagios.NewPlugin()

	defer plugin.ReturnCheckResults()

	plugin.ServiceOutput = "UNKNOWN: " + errString
	plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode
	plugin.ReturnCheckResults()
}

func CalculateDiskUsagePercentage(total, free int64) int {
	return int(100 * (total - free) / total)
}
