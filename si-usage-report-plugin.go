package main

import (
	"code.cloudfoundry.org/cli/plugin"
		)

type SIUsageReport struct{
	CliConnection plugin.CliConnection
}

func (c *SIUsageReport) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command si-usage-report
	switch args[0] {
	case "si-usage-report":
		c.GetSIUsageReport(args)
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

func (c *SIUsageReport) GetSIUsageReport(args []string) {
	c.CliConnection.CliCommand("cf -h")
}

func main() {
	plugin.Start(new(SIUsageReport))
}
