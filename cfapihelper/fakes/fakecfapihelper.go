package fakes

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"errors"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"time"
)

type FakeAPIHelper struct {
	CliConnection *pluginfakes.FakeCliConnection
}

type CFAPIHelper interface {
	GetOrgs() ([]plugin_models.GetOrgs_Model, error)
	GetServices() ([]plugin_models.GetServices_Model, error)
	GetServicePlans() ([]plugin_models.GetService_ServicePlan, error)
	GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error)
	GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error)
	GetServiceInstancePlanDetails(servicePlanURL string) (string, error)
	GetServiceInstanceServiceDetails(serviceURL string) (string, error)
	IsLoggedIn() (bool, error)
}

var _ cfapihelper.CFAPIHelper = new(FakeAPIHelper)

func (f *FakeAPIHelper) GetOrgs() ([]plugin_models.GetOrgs_Model, error) {
	return nil, errors.New("blah")
}

func (f *FakeAPIHelper) GetServices() ([]plugin_models.GetServices_Model, error) {
	return nil, errors.New("blah")
}

func (f *FakeAPIHelper) GetServicePlans() ([]plugin_models.GetService_ServicePlan, error) {
	return nil, errors.New("blah")
}

func (f *FakeAPIHelper) GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error) {
	return nil, errors.New("blah")
}

func (f *FakeAPIHelper) GetServiceInstancesWithDetails() ([]cfapihelper.ServiceInstance_Details, error) {
	return []cfapihelper.ServiceInstance_Details{
		{
			Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
			Name:      "test-si-2",
			Org:       "test-org",
			Space:     "test-space",
			Plan:      "spark",
			Service:   "cleardb",
			Type:      "managed_service_instance",
			CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
		}}, nil
}

func (f *FakeAPIHelper) GetServiceInstancePlanDetails(servicePlanURL string) (string, error) {
	return "", errors.New("blah")
}

func (f *FakeAPIHelper) GetServiceInstanceServiceDetails(serviceURL string) (string, error) {
	return "", errors.New("blah")
}

func (f *FakeAPIHelper) IsLoggedIn() (bool, error) {
	return true, nil
}
