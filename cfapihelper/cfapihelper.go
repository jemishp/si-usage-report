package cfapihelper

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
	"strconv"
	"sync"
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
	CreatedAt string
}

func New(cliConnection plugin.CliConnection) CFAPIHelper {
	return &APIHelper{cliConnection: cliConnection}
}

func (a *APIHelper) GetServiceInstancesWithDetails() ([]ServiceInstance_Details, error) {
	var (
		serviceInstanceDetailsJSON map[string]interface{}
		err                        error
		wg                         sync.WaitGroup
		serviceInstances           []ServiceInstance_Details
	)
	queryPath := "/v2/service_instances?q=&inline-relations-depth=2&" +
		"exclude-relations=developers,managers,auditors,domains,security_groups,staging_security_groups," +
		"apps,routes,service_keys,service_bindings,app,service_binding"
	serviceInstanceDetailsJSON, err = cfcurl.Curl(a.cliConnection, queryPath)
	if err != nil {
		return nil, err
	}
	pages := int(serviceInstanceDetailsJSON["total_pages"].(float64)) // *0 + 100
	results := make(chan ServiceInstance_Details, 50)

	wg.Add(pages)

	//read from last page to page #1
	for i := pages; i > 0; i-- {
		pathUrl := queryPath + "&page=" + strconv.Itoa(i)
		go a.processUrl(&wg, pathUrl, results)
	}

	go func() {
		for instance := range results {
			//	//fmt.Println("trying to read results")
			//	//fmt.Print("adding results to the service instance: ")
			//	//fmt.Println(instance)
			serviceInstances = append(serviceInstances, instance)
		}
	}()
	wg.Wait()
	return serviceInstances, nil
}

func (a *APIHelper) IsLoggedIn() (bool, error) {
	return a.cliConnection.IsLoggedIn()
}

func (a *APIHelper) parseResource(theServiceInstance map[string]interface{}) ServiceInstance_Details {
	entity := theServiceInstance["entity"].(map[string]interface{})
	metadata := theServiceInstance["metadata"].(map[string]interface{})
	dateString := metadata["created_at"].(string)
	//getting associated plan name
	theServicePlan := entity["service_plan"].(map[string]interface{})
	servicePlanEntity := theServicePlan["entity"].(map[string]interface{})
	servicePlanName := servicePlanEntity["name"].(string)

	//getting associated service name
	theService := servicePlanEntity["service"].(map[string]interface{})
	serviceEntity := theService["entity"].(map[string]interface{})
	serviceName := serviceEntity["label"].(string)

	switch serviceName {
	case "p.mysql", "p.redis", "p.pcc", "p.rabbit":
		//getting associated space & org name
		theSpace := entity["space"].(map[string]interface{})
		spaceEntity := theSpace["entity"].(map[string]interface{})
		spaceName := spaceEntity["name"].(string)
		theOrg := spaceEntity["organization"].(map[string]interface{})
		orgEntity := theOrg["entity"].(map[string]interface{})
		orgName := orgEntity["name"].(string)

		resource := ServiceInstance_Details{
			Guid:      metadata["guid"].(string),
			Name:      entity["name"].(string),
			Org:       orgName,
			Space:     spaceName,
			Plan:      servicePlanName,
			Service:   serviceName,
			Type:      entity["type"].(string),
			CreatedAt: dateString,
		}
		return resource
	}
	return ServiceInstance_Details{}

}

func (a *APIHelper) processUrl(waitGroup *sync.WaitGroup, pathUrl string, results chan<- ServiceInstance_Details) {
	//fmt.Println("worker: ", id, " waking up")
	//fmt.Println("send message to increment counter")
	defer waitGroup.Done()
	serviceInstanceDetailsJSON, err := cfcurl.Curl(a.cliConnection, pathUrl)
	if err != nil {
		fmt.Println("error ", err)
	}
	//fmt.Println("procesing: ", pathUrl)

	for _, resource := range serviceInstanceDetailsJSON["resources"].([]interface{}) {
		castedResource := resource.(map[string]interface{})
		parsedResource := a.parseResource(castedResource)
		if parsedResource != (ServiceInstance_Details{}) {
			results <- parsedResource
			time.Sleep(1 * time.Millisecond)
		}
	}
}
