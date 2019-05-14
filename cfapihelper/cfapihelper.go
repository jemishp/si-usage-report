package cfapihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/krujos/cfcurl"
	"strconv"
	"time"
)

type CFAPIHelper interface {
	GetOrgs() ([]plugin_models.GetOrgs_Model, error)
	GetServices() ([]plugin_models.GetServices_Model, error)
	GetServicePlans() ([]plugin_models.GetService_ServicePlan, error)
	GetServiceInstances() ([]plugin_models.GetSpace_ServiceInstance, error)
	GetServiceInstancesWithDetails() ([]ServiceInstance_Details, error)
	GetServiceInstancePlanDetails(servicePlanURL string) (string, error)
	GetServiceInstanceServiceDetails(serviceURL string) (string, error)
	IsLoggedIn() (bool, error)
}

type APIHelper struct {
	cliConnection plugin.CliConnection
}

type ServiceInstance_Details struct {
	Guid      string
	Name      string
	Org       string
	Space     string
	Plan      string
	Service   string
	Type      string
	CreatedAt time.Time
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

func (a *APIHelper) GetServiceInstancesWithDetails() ([]ServiceInstance_Details, error) {
	queryPath := "/v2/service_instances?q=&inline-relations-depth=2"
	serviceInstanceDetailsJSON, err := cfcurl.Curl(a.cliConnection, queryPath)
	if err != nil {
		return nil, err
	}
	pages := int(serviceInstanceDetailsJSON["total_pages"].(float64))
	layout := "2006-01-02T15:04:05Z"
	var serviceInstances []ServiceInstance_Details
	for i := 1; i <= pages; i++ {
		if i != 1 {
			serviceInstanceDetailsJSON, err = cfcurl.Curl(a.cliConnection, queryPath+"&page="+strconv.Itoa(i))
		}
		for _, o := range serviceInstanceDetailsJSON["resources"].([]interface{}) {
			theServiceInstance := o.(map[string]interface{})
			entity := theServiceInstance["entity"].(map[string]interface{})
			metadata := theServiceInstance["metadata"].(map[string]interface{})
			dateString := metadata["created_at"].(string)
			createdTime, err := time.Parse(layout, dateString)
			if err != nil {
				return nil, err
			}

			//getting associated space & org name
			theSpace := entity["space"].(map[string]interface{})
			spaceEntity := theSpace["entity"].(map[string]interface{})
			spaceName := spaceEntity["name"].(string)
			theOrg := spaceEntity["organization"].(map[string]interface{})
			orgEntity := theOrg["entity"].(map[string]interface{})
			orgName := orgEntity["name"].(string)

			//getting associated plan name
			theServicePlan := entity["service_plan"].(map[string]interface{})
			servicePlanEntity := theServicePlan["entity"].(map[string]interface{})
			servicePlanName := servicePlanEntity["name"].(string)

			//getting associated service name
			theService := servicePlanEntity["service"].(map[string]interface{})
			serviceEntity := theService["entity"].(map[string]interface{})
			serviceName := serviceEntity["label"].(string)

			serviceInstances = append(serviceInstances,
				ServiceInstance_Details{
					Guid:      metadata["guid"].(string),
					Name:      entity["name"].(string),
					Org:       orgName,
					Space:     spaceName,
					Plan:      servicePlanName,
					Service:   serviceName,
					Type:      entity["type"].(string),
					CreatedAt: createdTime,
				})
		}
	}
	return serviceInstances, nil
}

func (a *APIHelper) GetServiceInstancePlanDetails(servicePlanURL string) (string, error) {
	serviceInstancePlanDetailsJSON, err := cfcurl.Curl(a.cliConnection, servicePlanURL)
	if err != nil {
		return "", err
	}
	var servicePlans plugin_models.GetService_ServicePlan
	entity := serviceInstancePlanDetailsJSON["entity"].(map[string]interface{})
	metadata := serviceInstancePlanDetailsJSON["metadata"].(map[string]interface{})
	servicePlans = plugin_models.GetService_ServicePlan{
		Name: entity["name"].(string),
		Guid: metadata["guid"].(string),
	}

	return servicePlans.Name, nil
}

func (a *APIHelper) GetServiceInstanceServiceDetails(serviceURL string) (string, error) {
	serviceInstanceServiceDetailsJSON, err := cfcurl.Curl(a.cliConnection, serviceURL)
	if err != nil {
		return "", err
	}
	var service plugin_models.GetService_Model
	entity := serviceInstanceServiceDetailsJSON["entity"].(map[string]interface{})
	metadata := serviceInstanceServiceDetailsJSON["metadata"].(map[string]interface{})
	service = plugin_models.GetService_Model{
		Name: entity["label"].(string),
		Guid: metadata["guid"].(string),
	}

	return service.Name, nil
}

func (a *APIHelper) IsLoggedIn() (bool, error) {
	return a.cliConnection.IsLoggedIn()
}

func fillinPlanServiceDetails(serviceInstances []ServiceInstance_Details) ([]ServiceInstance_Details, error) {
	//for i, si := range serviceInstances.([]ServiceInstance_Details) {
	//
	//}
	return nil, nil
}
