package fakes

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	"errors"
)

type FakeAPIHelper struct {
	CliConnection *pluginfakes.FakeCliConnection
}

var _ cfapihelper.CFAPIHelper = new(FakeAPIHelper)

func (f *FakeAPIHelper) GetOrgs() ([]plugin_models.GetOrgs_Model, error) {
	return nil, errors.New("blah")
}
