package cfapihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/krujos/cfcurl"
	"strconv"
)

type CFAPIHelper interface {
	GetOrgs() ([]plugin_models.GetOrgs_Model, error)
}

type APIHelper struct {
	cliConnection plugin.CliConnection
}

func New(cliConnection plugin.CliConnection) CFAPIHelper {
	return &APIHelper{cliConnection: cliConnection}
}

func (a *APIHelper) GetOrgs() ([]plugin_models.GetOrgs_Model, error) {
	orgsJSON, err := cfcurl.Curl(a.cliConnection, "/v2/organizations")
	if err != nil {
		return nil, err
	}

	pages := int(orgsJSON["total_pages"].(float64))
	var orgs []plugin_models.GetOrgs_Model
	for i := 1; i <= pages; i++ {
		if 1 != i {
			orgsJSON, err = cfcurl.Curl(a.cliConnection, "/v2/organizations?page="+strconv.Itoa(i))
		}
		for _, o := range orgsJSON["resources"].([]interface{}) {
			theOrg := o.(map[string]interface{})
			entity := theOrg["entity"].(map[string]interface{})
			metadata := theOrg["metadata"].(map[string]interface{})
			orgs = append(orgs,
				plugin_models.GetOrgs_Model{
					Guid: metadata["guid"].(string),
					Name: entity["name"].(string),
				})
				//Organization{
				//	Name:      entity["name"].(string),
				//	URL:       metadata["url"].(string),
				//	QuotaURL:  entity["quota_definition_url"].(string),
				//	SpacesURL: entity["spaces_url"].(string),
				//})
		}
	}
	return orgs, nil
}
