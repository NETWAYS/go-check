package checkhttpmock

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io"
	"net/http"
	"os"
	"path"
)

var (
	CurrentRecorder *Recorder
	// Default storage location for the recorded data
	RecordFile = TestData + "/httpmock-record.yml"
)

// Helper to record http.Response and http.Request for httpmock when no responser was found
type Recorder struct {
	RecordFile string
	Writer     io.Writer
}

// Activate a Recorder and set it as noResponder with httpmock
//
// Usually you would use this during developing unit tests, and not for finished tests.
//
// By default records to RecordFile
func ActivateRecorder() (rec *Recorder) {
	rec = &Recorder{RecordFile: RecordFile}

	httpmock.RegisterNoResponder(rec.Respond)

	CurrentRecorder = rec

	return
}

// Open the Writer when needed
func (rec *Recorder) writer() io.Writer {
	if rec.Writer != nil {
		return rec.Writer
	} else if rec.RecordFile != "-" && rec.RecordFile != "" {
		// Ensure directory is writable
		dir := path.Dir(rec.RecordFile)
		_ = os.MkdirAll(dir, 0755)

		// Open file in append mode
		f, err := os.OpenFile(RecordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			rec.Writer = f
		}
	}

	if rec.Writer == nil {
		rec.Writer = os.Stdout
	}

	return rec.Writer
}

// Handle http.request for httpmock, and execute a real HTTP connection using httpmock.InitialTransport
//
// Recording and returning the real http.Response
func (rec *Recorder) Respond(request *http.Request) (response *http.Response, err error) {
	r := NewRecord(request)

	// Do a real request bypassing mock
	response, err = httpmock.InitialTransport.RoundTrip(request)
	if err != nil {
		err = fmt.Errorf("could not execute HTTP request: %w", err)
		return
	}

	r.Complete(response)

	err = r.EmitYAML(rec.writer())

	return
}
