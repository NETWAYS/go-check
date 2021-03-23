package checkhttpmock_test

import (
	"fmt"
	"github.com/NETWAYS/go-check/http/mock"
	"github.com/jarcoal/httpmock"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func ExampleActivateRecorder() {
	// Activate the normal httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Activate recorder
	_ = os.Remove(checkhttpmock.RecordFile) // Remove any prior recording
	checkhttpmock.ActivateRecorder()

	// We don't set any mock examples here
	//httpmock.RegisterResponder("GET", "http://localhost:8080/test",
	//	func(request *http.Request) (*http.Response, error) {
	//		return httpmock.NewStringResponse(200, "Hello World"), nil
	//	})

	// Start a simple HTTP server
	runHTTP()

	resp, err := http.Get("http://localhost:8080/test") // nolint:noctx
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print response body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", data)

	// Print recording
	data, _ = ioutil.ReadFile(checkhttpmock.RecordFile)
	fmt.Printf("%s\n", data)

	_ = resp.Body.Close()

	// Output:
	// Hello World
	// ---
	// url: http://localhost:8080/test
	// method: GET
	// query: ""
	// status: 200 OK
	// body: Hello World
}

func runHTTP() {
	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, `Hello World`)
	})

	go http.ListenAndServe(":8080", nil) //nolint:errcheck
}
