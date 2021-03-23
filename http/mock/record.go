package checkhttpmock

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
)

// Data structure to store information about a http.Request and http.Response in a simplified way
type Record struct {
	URL    string
	Method string
	Query  string
	Status string
	Body   string
}

// Build a new Record from an http.Request
func NewRecord(request *http.Request) (r *Record) {
	r = &Record{
		URL:    request.URL.String(),
		Method: request.Method,
	}

	// read the query from the request
	r.Query = extractFormQuery(request)

	log.WithFields(log.Fields{
		"url":    r.URL,
		"method": r.Method,
	}).Info("recording request")

	return
}

// Update the Record with a http.Response to get Body and Status
func (r *Record) Complete(response *http.Response) {
	body, newReader := dumpAndBuffer(response.Body)
	response.Body = newReader

	r.Status = response.Status
	r.Body = body

	log.WithFields(log.Fields{
		"status": response.Status,
	}).Info("recording response")
}

// Write a YAML representation of the Record to an io.Writer
func (r Record) EmitYAML(w io.Writer) (err error) {
	out := yaml.NewEncoder(w)
	out.SetIndent(2)

	_, _ = fmt.Fprintln(w, "---")

	err = out.Encode(r)

	return
}
