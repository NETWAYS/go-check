package checkhttpmock

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const contentTypeHeader = "Content-Type"
const contentTypeUrlencoded = "application/x-www-form-urlencoded"

// Read all data from a io.ReadCloser, return the data as string and return a new io.ReadCloser to pass on
//
// This can be quite tricky and is only used for mocking and testing here.
func dumpAndBuffer(r io.ReadCloser) (string, io.ReadCloser) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	_ = r.Close()

	return string(data), ioutil.NopCloser(bytes.NewReader(data))
}

// Extract a URL query from the request body, when the Content-Type is set to be urlencoded
func extractFormQuery(request *http.Request) string {
	if strings.Contains(request.Header.Get(contentTypeHeader), contentTypeUrlencoded) {
		query, newReader := dumpAndBuffer(request.Body)
		request.Body = newReader

		return query
	}

	return ""
}
