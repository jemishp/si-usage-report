package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"io"
)

type SIUsageReport struct {
	CliConnection plugin.CliConnection
	APIHelper     cfapihelper.CFAPIHelper
	OutBuf        io.Writer
}

func (s *SIUsageReport) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command si-usage-report
	switch args[0] {
	case "si-usage-report":
		s.CliConnection = cliConnection
		s.APIHelper = cfapihelper.New(s.CliConnection)
		s.GetSIUsageReport(args)
	}
}

func (s *SIUsageReport) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "si-usage-report",
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

func (s *SIUsageReport) GetSIUsageReport(args []string) {
	logedIn, err := s.APIHelper.IsLoggedIn()
	if err != nil {
		fmt.Fprintf(s.OutBuf, "error checking login status: %s", err.Error())
	}
	if logedIn {
		sis, err := s.APIHelper.GetServiceInstancesWithDetails()
		if err != nil {
			fmt.Fprintf(s.OutBuf, "error while getting service instances: %s", err.Error())
		} else {
			sisJSON, err := json.Marshal(sis)
			if err != nil {
				fmt.Fprintf(s.OutBuf, "error converting to json: %s", err.Error())
			}
			fmt.Fprintf(s.OutBuf, "service-instances: %s", sisJSON)
		}
	} else {
		fmt.Fprint(s.OutBuf, "need to log in")
	}
}

func main() {
	plugin.Start(new(SIUsageReport))
}
