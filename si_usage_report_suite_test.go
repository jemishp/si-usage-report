package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSiUsageReport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SiUsageReport Suite")
}
