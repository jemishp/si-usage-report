package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	Client http.Client
}

func New(duration time.Duration) HttpClient {
	client := http.Client{
		Timeout: duration,
	}
	return HttpClient{
		Client: client,
	}
}

func (h HttpClient) GetRequestWithHeader(url string) *http.Request {
	var body []byte
	request, _ := http.NewRequest("GET", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Host", "")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "bearer")
	request.Header.Set("User-Agent", "go-cli / si-usage-report")

	return request
}

func (h HttpClient) GetResponseBody(request *http.Request) []byte {
	response, _ := h.Client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return body

}
