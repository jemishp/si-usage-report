package fakes

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"time"
)

type FakeAPIHelper struct {
	CliConnection *pluginfakes.FakeCliConnection
}

type CFAPIHelper interface {
	GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error)
	IsLoggedIn() (bool, error)
}

var _ cfapihelper.CFAPIHelper = new(FakeAPIHelper)

func (f *FakeAPIHelper) GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error) {
	return []cfapihelper.ServiceInstance_Details{
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-1",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "cleardb",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-2",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "p.mysql",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-3",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "p.pcc",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-3",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "p.redis",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		},
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-3",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "p.rabbit",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		},
	}, nil
}

func (f *FakeAPIHelper) IsLoggedIn() (bool, error) {
	return true, nil
}
