package main_test

import (
	"code.cloudfoundry.org/cli/plugin"
	. "github.com/jpatel-pivotal/si-usage-report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SiUsageReport", func() {
	var (
		subject               *BasicPlugin
		expectedPluginVersion plugin.VersionType
		expectedCLIVersion    plugin.VersionType
		expectedCommand       plugin.Command
	)
	BeforeEach(func() {
		subject = new(BasicPlugin)
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
	When("cf si-usage-report is run", func() {
		It("prints a success message", func() {

		})
	})

})
