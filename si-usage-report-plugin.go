package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"io"
	"os"
)

type SIUsageReport struct {
	CliConnection plugin.CliConnection
	APIHelper     cfapihelper.CFAPIHelper
	OutBuf        io.Writer
}

type Report struct {
	Orgs []Org
}

type Org struct {
	OrgName              string
	SpaceName            string
	ProductName          string
	PlanName             string
	ServiceInstanceCount int
}

func (s *SIUsageReport) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command si-usage-report
	switch args[0] {
	case "si-usage-report":
		s.CliConnection = cliConnection
		s.APIHelper = cfapihelper.New(s.CliConnection)
		s.OutBuf = os.Stdout
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
			report := s.GenerateReport(sis)
			sisJSON, err := json.Marshal(report)
			if err != nil {
				fmt.Fprintf(s.OutBuf, "error converting to json: %s", err.Error())
			}
			fmt.Fprint(s.OutBuf, string(sisJSON))
		}
	} else {
		fmt.Fprint(s.OutBuf, "error: not logged in.\n run cf login")
	}
}

func main() {
	plugin.Start(new(SIUsageReport))
}

func (s *SIUsageReport) GenerateReport(serviceInstances []cfapihelper.ServiceInstance_Details) []Org {
	var report []Org
	for _, si := range serviceInstances {
		switch si.Service {
		case "p.mysql", "p.redis", "p.pcc", "p.rabbit":
			newOrg := Org{
				OrgName:              si.Org,
				SpaceName:            si.Space,
				ProductName:          si.Service,
				PlanName:             si.Plan,
				ServiceInstanceCount: 1,
			}
			report = append(report, newOrg)
		}

	}
	return report
}