package zendesk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Error an error type containing the http response from zendesk
type Error struct {
	body []byte
	resp *http.Response
}

// NewError is a function to initialize the Error type. This function will be useful
// for unit testing and mocking purposes in the client side
// to test their behavior by the API response.
//
// resp may be nil, in which case Error, Status and Headers report the body,
// no status code and no headers rather than panicking.
func NewError(body []byte, resp *http.Response) Error {
	return Error{
		body: body,
		resp: resp,
	}
}

// Error the error string for this error
func (e Error) Error() string {
	msg := string(e.body)

	// An Error built without a response, e.g. by NewError for mocking, has no
	// status code to report.
	if e.resp == nil {
		if msg == "" {
			return "unknown zendesk API error"
		}
		return msg
	}

	if msg == "" {
		msg = http.StatusText(e.resp.StatusCode)
	}

	return fmt.Sprintf("%d: %s", e.resp.StatusCode, msg)
}

// Body is the Body of the HTTP response
func (e Error) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(e.body))
}

// Headers the HTTP headers returned from zendesk
func (e Error) Headers() http.Header {
	if e.resp == nil {
		return nil
	}

	return e.resp.Header
}

// Status the HTTP status code returned from zendesk
func (e Error) Status() int {
	if e.resp == nil {
		return 0
	}

	return e.resp.StatusCode
}

// OptionsError is an error type for invalid option argument.
type OptionsError struct {
	opts interface{}
}

func (e *OptionsError) Error() string {
	return fmt.Sprintf("invalid options: %v", e.opts)
}
