package cfapihelper_test

import (
	"bufio"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"fmt"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
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
