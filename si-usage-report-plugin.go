package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"io"
	"os"
	"sort"
)

type SIUsageReport struct {
	CliConnection plugin.CliConnection
	APIHelper     cfapihelper.CFAPIHelper
	OutBuf        io.Writer
}

type Report struct {
	Products []Product
}

type Org struct {
	OrgName string
	Spaces  []Space
}

type Space struct {
	Name     string
	Products []Product
}

type Product struct {
	Name  string
	Plans []Plan
}

type Plan struct {
	PlanName      string
	InstanceCount int
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
		return
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

func (s *SIUsageReport) GenerateReport(serviceInstances []cfapihelper.ServiceInstance_Details) Report {
	var report Report
	report = Report{}
	planMap := make(map[string][]Plan, 0)

	for _, si := range serviceInstances {
		if _, ok := planMap[si.Service]; !ok {
			planMap[si.Service] = []Plan{
				{
					PlanName:      si.Plan,
					InstanceCount: +1,
				},
			}
		} else {
			var planModified bool
			planModified = false
			for k, plan := range planMap[si.Service] {
				if plan.PlanName == si.Plan {
					plan.InstanceCount += 1
					planMap[si.Service][k] = plan
					planModified = true
				}
			}
			if !planModified {

				plan := Plan{
					PlanName:      si.Plan,
					InstanceCount: +1,
				}
				planMap[si.Service] = append(planMap[si.Service], plan)
			}

		}


	}
	productKeys := make([]string, 0, len(planMap))
	for a, _ := range planMap {
		productKeys = append(productKeys, a)
	}
	sort.Strings(productKeys)
	for _, productKey := range productKeys {
		newProduct := Product{
			Name:  productKey,
			Plans: planMap[productKey],
		}
		report.Products = append(report.Products, newProduct)
	}

	return report
}
