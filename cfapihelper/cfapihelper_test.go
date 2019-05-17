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

var _ = Describe("cfapiHelper", func() {
	Describe("Happy path", func() {
		var (
			apiHelper            cfapihelper.CFAPIHelper
			fakeClientConnection *pluginfakes.FakeCliConnection
		)

		BeforeEach(func() {
			fakeClientConnection = &pluginfakes.FakeCliConnection{}
			apiHelper = cfapihelper.New(fakeClientConnection)
		})

		Context("2 service instances with details exist", func() {
			var serviceInstancesJSON []string
			var expectedServiceInstances []cfapihelper.ServiceInstance_Details

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
				expectedServiceInstances = []cfapihelper.ServiceInstance_Details{
					{
						Guid: "2f7f9ddb-0c1b-4198-b164-081f3f05f059",
						Name: "test-upgrade",
						Org: "dedicated-mysql-dev",
						Space: "staging",
						Plan: "db-small",
						Service: "p.mysql-staging",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2017, 04, 18, 18, 44, 16, 0, time.UTC),
					},
					{
						Guid: "168f3d45-1e11-4966-9fdc-b52632f00fba",
						Name: "datadog-metrics-work-alerting",
						Org: "dedicated-mysql-dev",
						Space: "acceptance",
						Plan: "db-lf-small",
						Service: "p.mysql",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2017, 9, 07, 20, 55, 22, 0, time.UTC),
					},
					{
						Guid: "250468d8-be73-475e-bb1e-e62d28b46b08",
						Name: "test-lf-upgrade",
						Org: "dedicated-mysql-dev",
						Space: "staging",
						Plan: "db-lf-small",
						Service: "p.mysql-staging",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2017, 11, 27, 23, 07, 37, 0, time.UTC),
					},
					{
						Guid: "6ce0a8d2-c23f-47ca-9252-1a98be6543b3",
						Name: "migrate1",
						Org: "dedicated-mysql-dev",
						Space: "acceptance",
						Plan: "db-small",
						Service: "p.mysql",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 02, 21, 23, 30, 26, 0, time.UTC),
					},
					{
						Guid: "b5c56b9e-5151-4e6d-a0e5-367df389a132",
						Name: "from",
						Org: "dedicated-mysql-dev",
						Space: "acceptance",
						Plan: "1gb",
						Service: "p-mysql",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 02, 21, 23, 37, 25, 0, time.UTC),
					},
					{
						Guid: "6e4f3b2d-503d-4ed7-b0f8-85511ae93a7f",
						Name: "to",
						Org: "dedicated-mysql-dev",
						Space: "acceptance",
						Plan: "db-small",
						Service: "p.mysql",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 03, 12, 19, 04, 26, 0, time.UTC),
					},
					{
						Guid: "da68e757-390c-41e2-bb5c-c4afeed5a0cc",
						Name: "migrate-test-staging",
						Org: "dedicated-mysql-dev",
						Space: "staging",
						Plan: "db-small",
						Service: "p.mysql-staging",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 10, 03, 17, 04, 56, 0, time.UTC),
					},
					{
						Guid: "02e65477-1e18-4a23-8bda-fabd5a4ac758",
						Name: "test-ha-upgrade",
						Org: "dedicated-mysql-dev",
						Space: "staging",
						Plan: "db-ha-small",
						Service: "p.mysql-staging",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 10, 11, 22, 35, 49, 0, time.UTC),
					},
					{
						Guid: "c0ca147f-c84a-4d11-88b1-fed254943a73",
						Name: "ha-staging-test",
						Org: "dedicated-mysql-dev",
						Space: "staging",
						Plan: "db-ha-small",
						Service: "p.mysql-staging",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2018, 11, 06, 00, 05, 40, 0, time.UTC),
					},
					{
						Guid: "d97fa089-e5c8-4e22-b2da-5ceb49db85f2",
						Name: "backup-test-ha",
						Org: "dedicated-mysql-dev",
						Space: "acceptance",
						Plan: "db-ha-small",
						Service: "p.mysql",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2019, 05, 01, 18, 56, 01, 0, time.UTC),
					},
					{
						Guid: "31e3efc1-1898-4156-b936-a666131483e1",
						Name: "test-si",
						Org: "jpatel-org",
						Space: "development",
						Plan: "spark",
						Service: "cleardb",
						Type: "managed_service_instance",
						CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
					},
				}
			})
			It("should return list of service instances with details", func() {
				serviceInstances, err := apiHelper.GetServiceInstancesWithDetails()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serviceInstances)).To(Equal(11))
				Expect(serviceInstances).To(Equal(expectedServiceInstances))
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
		Context("really bad error occurs", func() {
			BeforeEach(func() {
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(nil,
					errors.New("really bad error"))
			})
			It("should return an error", func() {
				returned, err := apiHelper.GetServiceInstancesWithDetails()
				Expect(err).To(MatchError("really bad error"))
				Expect(len(returned)).To(Equal(0))
			})
		})
		Context("not logged in error occurs", func() {
			BeforeEach(func() {
				fakeClientConnection.IsLoggedInReturns(false, errors.New("really bad login error"))
			})
			It("should return an error", func() {
				returned, err := apiHelper.IsLoggedIn()
				Expect(err).To(MatchError("really bad login error"))
				Expect(returned).To(Equal(false))
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
