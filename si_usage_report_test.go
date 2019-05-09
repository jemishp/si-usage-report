package main_test

import (
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/jpatel-pivotal/si-usage-report"
	"github.com/jpatel-pivotal/si-usage-report/cfapihelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
	"os/exec"
)

var _ = Describe("SiUsageReport", func() {
	Describe("happy path", func() {
		var (
			subject               *SIUsageReport
			expectedPluginVersion plugin.VersionType
			expectedCLIVersion    plugin.VersionType
			expectedCommand       plugin.Command
			outBuffer             io.Writer
			errBuffer             io.Writer
			fakeCLIConnection     *pluginfakes.FakeCliConnection
			apiHelper             cfapihelper.CFAPIHelper
		)

		BeforeEach(func() {
			subject = new(SIUsageReport)
			fakeCLIConnection = &pluginfakes.FakeCliConnection{}
			apiHelper = cfapihelper.New(fakeCLIConnection)
			subject.CliConnection = fakeCLIConnection
			subject.APIHelper = apiHelper
			expectedPluginVersion = plugin.VersionType{
				Major: 1,
				Minor: 0,
				Build: 0,
			}
			expectedCLIVersion = plugin.VersionType{
				Major: 6,
				Minor: 7,
				Build: 0,
			}
			expectedCommand = plugin.Command{
				Name:     "si-usage-report",
				HelpText: "Shows Service Instance Usage Report",
				UsageDetails: plugin.Usage{
					Usage: "si-usage-report\n   cf si-usage-report",
				},
			}
			outBuffer = gbytes.NewBuffer()
			errBuffer = gbytes.NewBuffer()
		})

		When("GetMetaData() is called", func() {
			It("returns the correct name for the plugin", func() {
				Expect(subject.GetMetadata().Name).To(Equal("SIUsageReport"))
			})
			It("returns the correct version of the plugin", func() {
				Expect(subject.GetMetadata().Version).To(Equal(expectedPluginVersion))
			})
			It("returns the correct min CLI version of the plugin", func() {
				Expect(subject.GetMetadata().MinCliVersion).To(Equal(expectedCLIVersion))
			})
			It("returns the correct command", func() {
				Expect(len(subject.GetMetadata().Commands)).To(Equal(1))
				Expect(subject.GetMetadata().Commands[0]).To(Equal(expectedCommand))
			})
		})
		//TODO Move this to integration test
		When("cf si-usage-report is run without installing the plugin", func() {
			subject = new(SIUsageReport)
			fakeCLIConnection = &pluginfakes.FakeCliConnection{}
			apiHelper = cfapihelper.New(fakeCLIConnection)
			subject.CliConnection = fakeCLIConnection
			subject.APIHelper = apiHelper
			subject.Run(fakeCLIConnection, []string{"si-usage-report"})
			//subject.GetSIUsageReport([]string{"test"})
			It("prints the usage message", func() {
				args := []string{"si-usage-report"}
				session, err := gexec.Start(exec.Command("cf", args...), outBuffer, errBuffer)
				session.Wait()
				Expect(err).NotTo(HaveOccurred())
				Expect(outBuffer).To(gbytes.Say("'si-usage-report' is not a registered command. See 'cf help -a'"))
				Expect(errBuffer).To(gbytes.Say(""))
			})
		})
		When("GetSIUsageReport is called", func() {
			It("prints completed", func() {

			})
		})

		//When("GetSIUsageReport is called", func() {
		//	var out []string
		//	out, err = subject.ApiHelper.CliCommand("")
		//	It("prints a message", func() {
		//		Expect(err).To(BeNil())
		//		Expect(out).To(Equal("Running the si-usage-report with args:"))
		//	})
		//})
		//When("GetSIUsageReport is called", func() {
		//	var out []string
		//	out, err = subject.ApiHelper.CliCommand("")
		//	It("prints a message", func() {
		//		Expect(err).To(BeNil())
		//		Expect(out).To(Equal("Running the si-usage-report with args:"))
		//	})
		//})
	})

})
