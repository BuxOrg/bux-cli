package cmd

import (
	"fmt"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/ryanuber/columnize"
)

// displayTracingResults displays the tracing results into the terminal per request
func displayTracingResults(tracing resty.TraceInfo, statusCode int) {

	// Add the network time columns
	output := []string{
		fmt.Sprintf(`DNSLookup | %s | TTFB | %s`, tracing.DNSLookup.String(), tracing.ServerTime.String()),
		fmt.Sprintf(`TLSHandshake | %s | ConnTime | %s`, tracing.TLSHandshake.String(), tracing.ConnTime.String()),
		fmt.Sprintf(`FB to Close  | %s | %s | %s [%d]`, tracing.ResponseTime.String(), "TotalTime", tracing.TotalTime.String(), statusCode),
	}

	// Connection was idle?
	if tracing.IsConnWasIdle {
		output = append(output,
			fmt.Sprintf(`IsConnWasIdle | %s | ConnIdleTime | %s`,
				color.MagentaString(fmt.Sprintf("%v", tracing.IsConnWasIdle)),
				color.MagentaString(tracing.ConnIdleTime.String()),
			))
	}

	// Connection reused?
	if tracing.IsConnReused {
		output = append(output, fmt.Sprintf(`IsConnReused | %s `, color.MagentaString(fmt.Sprintf("%v", tracing.IsConnReused))))
	}

	// Render the data
	chalker.Log(chalker.DIM, columnize.SimpleFormat(output))
}

// displayHeader will display a standard header
func displayHeader(level, text string) {
	chalker.Log(level, "\n==========| "+text)
}
