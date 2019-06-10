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

		Context("4 service instances with details are returned", func() {
			var serviceInstancesJSON []string
			var expectedServiceInstances []cfapihelper.ServiceInstance_Details

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
				expectedServiceInstances = []cfapihelper.ServiceInstance_Details{
					{
						Guid:      "168f3d45-1e11-4966-9fdc-b52632f00fba",
						Name:      "datadog-metrics-work-alerting",
						Org:       "dedicated-mysql-dev",
						Space:     "acceptance",
						Plan:      "db-lf-small",
						Service:   "p.mysql",
						Type:      "managed_service_instance",
						CreatedAt: "2017-09-07T20:55:22Z",
					},
					{
						Guid:      "6ce0a8d2-c23f-47ca-9252-1a98be6543b3",
						Name:      "migrate1",
						Org:       "dedicated-mysql-dev",
						Space:     "acceptance",
						Plan:      "db-small",
						Service:   "p.mysql",
						Type:      "managed_service_instance",
						CreatedAt: "2018-02-21T23:30:26Z",
					},
					{
						Guid:      "6e4f3b2d-503d-4ed7-b0f8-85511ae93a7f",
						Name:      "to",
						Org:       "dedicated-mysql-dev",
						Space:     "acceptance",
						Plan:      "db-small",
						Service:   "p.mysql",
						Type:      "managed_service_instance",
						CreatedAt: "2018-03-12T19:04:26Z",
					},
					{
						Guid:      "d97fa089-e5c8-4e22-b2da-5ceb49db85f2",
						Name:      "backup-test-ha",
						Org:       "dedicated-mysql-dev",
						Space:     "acceptance",
						Plan:      "db-ha-small",
						Service:   "p.mysql",
						Type:      "managed_service_instance",
						CreatedAt: "2019-05-01T18:56:01Z",
					},
				}
			})
			It("should return list of service instances with details", func() {
				serviceInstances, err := apiHelper.GetServiceInstancesWithDetails()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serviceInstances)).To(Equal(4))
				Expect(serviceInstances).To(Equal(expectedServiceInstances))
			})
		})
		Context("41350 service instances with details exist across pages", func() {
			var serviceInstancesJSON []string

			BeforeEach(func() {
				serviceInstancesJSON = getResponse("test-data/lot-of-service-instances.json")
				fakeClientConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
			})
			It("should return list of service instances with details", func() {
				serviceInstances, err := apiHelper.GetServiceInstancesWithDetails()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(serviceInstances)).To(Equal(9097))
			})
			Measure("it should do this efficiently", func(b Benchmarker) {
				runtime := b.Time("runtime", func() {
					serviceInstances, err := apiHelper.GetServiceInstancesWithDetails()
					Expect(err).NotTo(HaveOccurred())
					Expect(len(serviceInstances)).To(Equal(9097))
				})

				Ω(runtime.Seconds()).Should(BeNumerically("<", 7), "GetServiceInstancesWithDetails() shouldn't take too long.")

				//b.RecordValue("disk usage (in MB)", HowMuchDiskSpaceDidYouUse())
			}, 10)

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
