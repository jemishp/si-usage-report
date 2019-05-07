package cfapihelper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfapihelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfapihelper Suite")
}
