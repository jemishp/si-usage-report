package fakes

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"errors"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
)

type FakeAPIHelper struct {
	CliConnection *pluginfakes.FakeCliConnection
}

type CFAPIHelper interface {
	GetOrgs() ([]plugin_models.GetOrgs_Model, error)
	GetServices() ([]plugin_models.GetServices_Model, error)
	GetServicePlans() ([]plugin_models.GetService_ServicePlan, error)
	GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error)
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
