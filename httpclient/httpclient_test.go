package httpclient_test

import (
	. "github.com/jpatel-pivotal/si-usage-report/httpclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

var (
	subject HttpClient
)
var _ = Describe("Httpclient", func() {

	BeforeEach(func() {
		subject = New(5 * time.Second)
	})
	It("has the correct timeout", func() {
		Expect(subject.Client.Timeout).To(Equal(5 * time.Second))
	})
	It("generates requests with headers", func() {
		req := subject.GetRequestWithHeader("http://someURL.com")
		Expect(req.Method).To(Equal("GET"))
		Expect(req.Header.Get("Content-type")).To(Equal("application/json"))
		Expect(req.Header.Get("Host")).To(Equal(""))
		Expect(req.Header.Get("Accept")).To(Equal("application/json"))
		Expect(req.Header.Get("Authorization")).To(Equal("bearer"))
		Expect(req.Header.Get("User-Agent")).To(Equal("go-cli / si-usage-report"))
		Expect(req.Body).To(Equal(http.NoBody))
	})
})
