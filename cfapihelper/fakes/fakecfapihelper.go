package fakes

import (
	"bufio"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"os"
)

type FakeAPIHelper struct {
	CliConnection       *pluginfakes.FakeCliConnection
	GetResponseFromFile bool
}

type CFAPIHelper interface {
	GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error)
	IsLoggedIn() (bool, error)
	GetResponse(filename string) []string
}

var _ cfapihelper.CFAPIHelper = new(FakeAPIHelper)

func (f *FakeAPIHelper) GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error) {
	return []cfapihelper.ServiceInstance_Details{
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-1",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "10mb",
			Service:   "p.mysql",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-2",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "small",
			Service:   "p.pcc",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-3",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "medium",
			Service:   "p.redis",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-4",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "lemur",
			Service:   "p.rabbit",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-5",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "lemur",
			Service:   "p.rabbit",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-6",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "10mb",
			Service:   "p.mysql",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-7",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "10mb",
			Service:   "p.mysql",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-8",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "100mb",
			Service:   "p.mysql",
			Type:      "managed_service_instance",
			CreatedAt: "2018-02-21 23:30:26",
		},
	}, nil
}

func (f *FakeAPIHelper) IsLoggedIn() (bool, error) {
	return true, nil
}

func (f *FakeAPIHelper) GetResponse(filename string) []string {
	var b []string
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b = append(b, scanner.Text())
	}
	return b
}
