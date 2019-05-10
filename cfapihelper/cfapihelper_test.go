package cfapihelper_test

import (
	"bufio"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"errors"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("SiUsageReport", func() {
	Describe("Happy path", func() {
		var (
			apiHelper            cfapihelper.CFAPIHelper
			fakeClientConnection *pluginfakes.FakeCliConnection
		)

		BeforeEach(func() {
			fakeClientConnection = &pluginfakes.FakeCliConnection{}
			apiHelper = cfapihelper.New(fakeClientConnection)
		})

		When("2 orgs exist", func() {
			var orgsJSON []string

			BeforeEach(func() {
				orgsJSON = getResponse("test-data/orgs.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			})
			It("should return 2 orgs", func() {
				orgs, err := apiHelper.GetOrgs()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(orgs)).To(Equal(2))
			})
		})
		When("26 services exist", func() {
			var servicesJSON []string

			BeforeEach(func() {
				servicesJSON = getResponse("test-data/services.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(servicesJSON, nil)
			})
			It("should return list of services", func() {
				services, err := apiHelper.GetServices()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(services)).To(Equal(26))
			})
		})
		When("100 service plans exist", func() {
			var servicePlansJSON []string

			BeforeEach(func() {
				servicePlansJSON = getResponse("test-data/service-plans.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(servicePlansJSON, nil)
			})
			It("should return list of service plans", func() {
				services, err := apiHelper.GetServicePlans()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(services)).To(Equal(100))
			})
		})
		When("2 service instances exist", func() {
			var serviceInstancesJSON []string

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
			})
			It("should return list of service instances", func() {
				serviceInstances, err := apiHelper.GetServiceInstances()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serviceInstances)).To(Equal(2))
			})
		})
		When("2 service instances with details exist", func() {
			var serviceInstancesJSON []string
			var expectedServiceInstances []cfapihelper.ServiceInstance_Details

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
				expectedServiceInstances = []cfapihelper.ServiceInstance_Details{
					{
						Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
						Name:      "test-si-1",
						Plan:      "/v2/service_plans/fd12a21d-6667-4150-8193-884d083b7874",
						Service:   "/v2/services/5e30ff7e-d857-4aa7-9eda-7db9a0d7b19b",
						Type:      "managed_service_instance",
						CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
					},
					{
						Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
						Name:      "test-si-2",
						Plan:      "/v2/service_plans/fd12a21d-6667-4150-8193-884d083b7874",
						Service:   "/v2/services/5e30ff7e-d857-4aa7-9eda-7db9a0d7b19b",
						Type:      "managed_service_instance",
						CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
					},
				}
			})
			It("should return list of service instances with details", func() {
				serviceInstances, err := apiHelper.GetServiceInstancesWithDetails()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serviceInstances)).To(Equal(2))
				Expect(serviceInstances).To(Equal(expectedServiceInstances))
			})
		})
		When("getting the plan name for a service instance", func() {
			var serviceInstancePlanDetailsJSON []string

			BeforeEach(func() {
				serviceInstancePlanDetailsJSON = getResponse("test-data/si-plan-details.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancePlanDetailsJSON, nil)
			})
			It("returns service plan name of the service instance", func() {
				serviceInstancePlanName, err := apiHelper.GetServiceInstancePlanDetails("test")
				Expect(err).NotTo(HaveOccurred())
				Expect(serviceInstancePlanName).To(Equal("spark"))
			})
		})
		When("getting the service name for a service instance", func() {
			var serviceInstanceServiceDetailsJSON []string

			BeforeEach(func() {
				serviceInstanceServiceDetailsJSON = getResponse("test-data/si-service-details.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstanceServiceDetailsJSON, nil)
			})
			It("returns service plan name of the service instance", func() {
				serviceInstanceServiceName, err := apiHelper.GetServiceInstanceServiceDetails("test")
				Expect(err).NotTo(HaveOccurred())
				Expect(serviceInstanceServiceName).To(Equal("cleardb"))
			})
		})
	})
	Describe("testing error cases", func() {
		var (
			apiHelper            cfapihelper.CFAPIHelper
			fakeClientConnection *pluginfakes.FakeCliConnection
		)

		BeforeEach(func() {
			fakeClientConnection = &pluginfakes.FakeCliConnection{}
			apiHelper = cfapihelper.New(fakeClientConnection)
		})
		When("really bad error occurs", func() {
			BeforeEach(func() {
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(nil,
					errors.New("really bad error"))
			})
			It("should return an error", func() {
				returned, err := apiHelper.GetOrgs()
				Expect(err).To(MatchError("really bad error"))
				Expect(len(returned)).To(Equal(0))
			})
		})
		When("no orgs exist", func() {
			var orgsJSON []string

			BeforeEach(func() {
				orgsJSON = getResponse("test-data/no-orgs.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			})
			It("should return an error", func() {
				orgs, err := apiHelper.GetOrgs()
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(len(orgs)).To(Equal(0))
			})
		})
		When("no services exist", func() {
			var servicesJSON []string

			BeforeEach(func() {
				servicesJSON = getResponse("test-data/no-services.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(servicesJSON, nil)
			})
			It("should return an error", func() {
				services, err := apiHelper.GetServices()
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(len(services)).To(Equal(0))
			})
		})
		When("no service plans exist", func() {
			var servicePlansJSON []string

			BeforeEach(func() {
				servicePlansJSON = getResponse("test-data/no-service-plans.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(servicePlansJSON, nil)
			})
			It("should return an error", func() {
				services, err := apiHelper.GetServicePlans()
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(len(services)).To(Equal(0))
			})
		})
		When("no service instances exist", func() {
			var serviceInstancesJSON []string

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/no-service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
			})
			It("should return an error", func() {
				serviceInstances, err := apiHelper.GetServiceInstances()
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(len(serviceInstances)).To(Equal(0))
			})
		})
		When("no plan name for a service instance exists", func() {
			var serviceInstancePlanDetailsJSON []string

			BeforeEach(func() {
				serviceInstancePlanDetailsJSON = getResponse("test-data/no-si-plan-details.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancePlanDetailsJSON, nil)
			})
			It("should return list of service instances", func() {
				serviceInstancePlanName, err := apiHelper.GetServiceInstancePlanDetails("test")
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(serviceInstancePlanName).To(Equal(""))
			})
		})
		When("no service name for a service instance exists", func() {
			var serviceInstanceServiceDetailsJSON []string

			BeforeEach(func() {
				serviceInstanceServiceDetailsJSON = getResponse("test-data/no-si-plan-details.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstanceServiceDetailsJSON, nil)
			})
			It("should return list of service instances", func() {
				serviceInstanceServiceName, err := apiHelper.GetServiceInstanceServiceDetails("test")
				Expect(err).To(MatchError("CF API returned no output"))
				Expect(serviceInstanceServiceName).To(Equal(""))
			})
		})
	})

})

func getResponse(filename string) []string {
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
