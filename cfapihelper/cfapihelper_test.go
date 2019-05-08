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
	Describe("Setup", func() {
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
