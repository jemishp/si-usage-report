package cfapihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/krujos/cfcurl"
	"strconv"
)

type CFAPIHelper interface {
	GetOrgs() ([]plugin_models.GetOrgs_Model, error)
	GetServices() ([]plugin_models.GetServices_Model, error)
	GetServicePlans() ([]plugin_models.GetService_ServicePlan, error)
	GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error)
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
		}
	}
	return orgs, nil
}

func (a *APIHelper) GetServices() ([]plugin_models.GetServices_Model, error) {
	servicesJSON, err := cfcurl.Curl(a.cliConnection, "/v2/services")

	if err != nil {
		return nil, err
	}
	pages := int(servicesJSON["total_pages"].(float64))
	var services []plugin_models.GetServices_Model
	for i := 1; i <= pages; i++ {
		if i != 1 {
			servicesJSON, err = cfcurl.Curl(a.cliConnection, "/v2/services?page="+strconv.Itoa(i))
		}
		for _, o := range servicesJSON["resources"].([]interface{}) {
			theService := o.(map[string]interface{})
			entity := theService["entity"].(map[string]interface{})
			metadata := theService["metadata"].(map[string]interface{})
			services = append(services,
				plugin_models.GetServices_Model{
					Guid:             metadata["guid"].(string),
					Name:             entity["label"].(string),
					ServicePlan:      plugin_models.GetServices_ServicePlan{},
					Service:          plugin_models.GetServices_ServiceFields{},
					LastOperation:    plugin_models.GetServices_LastOperation{},
					ApplicationNames: nil,
					IsUserProvided:   false,
				})
		}

	}
	return services, nil
}

func (a *APIHelper) GetServicePlans() ([]plugin_models.GetService_ServicePlan, error) {
	servicePlansJSON, err := cfcurl.Curl(a.cliConnection, "/v2/service_plans")
	if err != nil {
		return nil, err
	}
	pages := int(servicePlansJSON["total_pages"].(float64))
	var servicePlans []plugin_models.GetService_ServicePlan
	for i := 1; i <= pages; i++ {
		if i != 1 {
			servicePlansJSON, err = cfcurl.Curl(a.cliConnection, "/v2/service_plans?page="+strconv.Itoa(i))
		}
		for _, o := range servicePlansJSON["resources"].([]interface{}) {
			theServicePlan := o.(map[string]interface{})
			entity := theServicePlan["entity"].(map[string]interface{})
			metadata := theServicePlan["metadata"].(map[string]interface{})
			servicePlans = append(servicePlans,
				plugin_models.GetService_ServicePlan{
					Name: entity["name"].(string),
					Guid: metadata["guid"].(string),
				})
		}

	}
	return servicePlans, nil
}

func (a *APIHelper) GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error) {
	serviceInstancesJSON, err := cfcurl.Curl(a.cliConnection, "/v2/service_instances")
	if err != nil {
		return nil, err
	}
	pages := int(serviceInstancesJSON["total_pages"].(float64))
	var serviceInstances []plugin_models.GetSpace_ServiceInstance
	for i := 1; i <= pages; i++ {
		if i != 1 {
			serviceInstancesJSON, err = cfcurl.Curl(a.cliConnection, "/v2/service_instances?page="+strconv.Itoa(i))
		}
		for _, o := range serviceInstancesJSON["resources"].([]interface{}) {
			theServiceInstance := o.(map[string]interface{})
			entity := theServiceInstance["entity"].(map[string]interface{})
			metadata := theServiceInstance["metadata"].(map[string]interface{})
			serviceInstances = append(serviceInstances,
				plugin_models.GetSpace_ServiceInstance{
					Name: entity["name"].(string),
					Guid: metadata["guid"].(string),
				})
		}

	}
	return serviceInstances, nil
}
