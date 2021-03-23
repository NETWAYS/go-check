package checkhttpmock_test

import (
	"bytes"
	"fmt"
	"github.com/NETWAYS/go-check/http/mock"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
)

func Example() {
	// Activate httpmock as normal
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Use any normal responder
	httpmock.RegisterResponder("GET", "https://example.com/test.json",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `{"allgood":true}`), nil
		})

	req, _ := http.NewRequest("GET", "https://example.com/test.json", nil) //nolint:noctx
	requestAndDump(req)

	// Use additional responders
	checkhttpmock.RegisterQueryMapResponder("POST", "https://exampleapi.com/",
		checkhttpmock.QueryMap{
			"test=1": "test.json",
		})

	req, _ = http.NewRequest("POST", "https://exampleapi.com/", bytes.NewBufferString("test=1")) //nolint:noctx
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requestAndDump(req)

	// Output:
	// {"allgood":true}
	// {"example":true}
}

func requestAndDump(req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}
