package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
)

type SIUsageReport struct{}

func (c *SIUsageReport) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	switch args[0] {
	case "si-usage-report":
		fmt.Println("Running the si-usage-report")
	}
}

func (c *SIUsageReport) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "SIUsageReport",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "si-usage-report",
				HelpText: "Shows Service Instance Usage Report",
				UsageDetails: plugin.Usage{
					Usage: "si-usage-report\n   cf si-usage-report",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(SIUsageReport))
}
