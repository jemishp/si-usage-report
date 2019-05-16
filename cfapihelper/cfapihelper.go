package cfapihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
	"strconv"
	"time"
)

type CFAPIHelper interface {
	GetServiceInstancesWithDetails() ([]ServiceInstance_Details, error)
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

func (a *APIHelper) GetServiceInstancesWithDetails() ([]ServiceInstance_Details, error) {
	queryPath := "/v2/service_instances?q=&inline-relations-depth=2&" +
		"exclude-relations=developers,managers,auditors,domains,security_groups,staging_security_groups," +
		"apps,routes,service_keys,service_bindings,app,service_binding"
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

func (a *APIHelper) IsLoggedIn() (bool, error) {
	return a.cliConnection.IsLoggedIn()
}
