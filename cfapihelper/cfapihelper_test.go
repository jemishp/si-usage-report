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
						Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
						Name:      "test-si-1",
						Org:       "test-org",
						Space:     "test-space",
						Plan:      "spark",
						Service:   "cleardb",
						Type:      "managed_service_instance",
						CreatedAt: time.Date(2019, 05, 06, 21, 18, 47, 0, time.UTC),
					},
					{
						Guid:      "31e3efc1-1898-4156-b936-a666131483e1",
						Name:      "test-si-2",
						Org:       "test-org",
						Space:     "test-space",
						Plan:      "spark",
						Service:   "cleardb",
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
