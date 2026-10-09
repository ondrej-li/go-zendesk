package zendesk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	status := http.StatusOK
	resp := &http.Response{
		StatusCode: status,
	}
	body := []byte("foo")
	err := Error{
		body: body,
		resp: resp,
	}

	expected := fmt.Sprintf("%d: %s", status, body)
	if v := err.Error(); v != expected {
		t.Fatalf("Error %s did not have expected value %s", v, expected)
	}
}

func TestError_ErrorEmptyBody(t *testing.T) {
	status := http.StatusOK
	resp := &http.Response{
		StatusCode: status,
	}
	err := Error{
		resp: resp,
	}

	expected := fmt.Sprintf("%d: %s", status, http.StatusText(status))
	if v := err.Error(); v != expected {
		t.Fatalf("Error %s did not have expected value %s", v, expected)
	}
}

func TestError_Headers(t *testing.T) {
	retryAfter := "Retry-After"
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			retryAfter: []string{"92"},
		},
	}

	err := Error{
		resp: resp,
	}

	if _, ok := err.Headers()[retryAfter]; !ok {
		t.Fatal("Could not get header values from zendesk error")
	}
}

func TestError_Status(t *testing.T) {
	retryAfter := "Retry-After"
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			retryAfter: []string{"92"},
		},
	}

	err := Error{
		resp: resp,
	}

	if status := err.Status(); status != http.StatusTooManyRequests {
		t.Fatal("Status returned from error was not the correct status code")
	}
}

func TestNewError(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusNotFound}
	body := []byte(`{"error":"RecordNotFound"}`)

	err := NewError(body, resp)

	if err.Status() != http.StatusNotFound {
		t.Fatalf("Unexpected status %d", err.Status())
	}

	expected := fmt.Sprintf("%d: %s", http.StatusNotFound, body)
	if v := err.Error(); v != expected {
		t.Fatalf("Error %s did not have expected value %s", v, expected)
	}
}

func TestOptionsError(t *testing.T) {
	err := &OptionsError{opts: "search options"}

	expected := "invalid options: search options"
	if v := err.Error(); v != expected {
		t.Fatalf("Error %s did not have expected value %s", v, expected)
	}
}

func TestErrorWithoutResponse(t *testing.T) {
	body := []byte(`{"error":"RecordNotFound"}`)
	err := NewError(body, nil)

	if v := err.Error(); v == "" {
		t.Fatal("Expected a non-empty error message")
	} else if v != string(body) {
		t.Fatalf("Error %s did not have expected value %s", v, body)
	}

	if v := err.Status(); v != 0 {
		t.Fatalf("Expected status 0 for a response-less error, got %d", v)
	}

	if v := err.Headers(); v != nil {
		t.Fatalf("Expected no headers for a response-less error, got %v", v)
	}

	read, err2 := ioutil.ReadAll(err.Body())
	if err2 != nil {
		t.Fatalf("Failed to read body: %s", err2)
	}
	if string(read) != string(body) {
		t.Fatalf("Body %s did not have expected value %s", read, body)
	}
}

func TestErrorWithoutResponseAndBody(t *testing.T) {
	err := NewError(nil, nil)

	if v := err.Error(); v == "" {
		t.Fatal("Expected a non-empty error message")
	}
}
