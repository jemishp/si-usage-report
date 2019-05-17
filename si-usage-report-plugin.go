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
	prodMap := make(map[string]interface{})
	planMap := make(map[string][]Plan)
	spaceMap := make(map[string]interface{})
	orgMap := make(map[string]interface{})

	for _, si := range serviceInstances {
		switch si.Service {
		case "p.mysql", "p.redis", "p.pcc", "p.rabbit":
			if _, ok := planMap[si.Plan]; !ok {
				planMap[si.Plan] = []Plan{
					{
						PlanName:      si.Plan,
						InstanceCount: +1,
					},
				}
			} else {
				for _, plan := range planMap[si.Plan] {
					if plan.PlanName == si.Plan {
						plan.InstanceCount += 1
						planMap[si.Plan] = []Plan{plan}
					}
				}

			}
			if _, ok := prodMap[si.Service]; !ok {
				prodMap[si.Service] = []Product{
					{
						Name:  si.Service,
						Plans: planMap[si.Plan],
					},
				}
			}

			if _, ok := spaceMap[si.Space]; !ok {
				spaceMap[si.Space] = []Space{
					{
						Name:     si.Space,
						Products: prodMap[si.Service].([]Product),
					},
				}
			}
			if _, ok := orgMap[si.Org]; !ok {
				orgMap[si.Org] = []Org{
					{
						OrgName: si.Org,
						Spaces:  spaceMap[si.Space].([]Space),
					},
				}
			}

			prodMap[si.Service] = planMap[si.Plan]

		}
	}
	productKeys := make([]string, 0, len(prodMap))
	for a,_ := range prodMap{
		productKeys = append(productKeys, a)
	}
	sort.Strings(productKeys)
	for _,productKey := range productKeys {
		newProduct := Product{
			Name:  productKey,
			Plans: prodMap[productKey].([]Plan),
		}
		report.Products = append(report.Products, newProduct)
	}

	return report
}
