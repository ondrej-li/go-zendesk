package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketSkips(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"skips":[{"id":1,"ticket_id":123,"user_id":42}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	skips, _, err := client.GetTicketSkips(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket skips: %s", err)
	}
	if len(skips) != 1 || skips[0].TicketID != 123 {
		t.Fatalf("Unexpected ticket skips %+v", skips)
	}
}

func TestGetTicketSkipsByTicket(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/tickets/123/skips.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"skips":[{"id":1,"ticket_id":123}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	skips, _, err := client.GetTicketSkipsByTicket(ctx, 123, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket skips by ticket: %s", err)
	}
	if len(skips) != 1 {
		t.Fatalf("Expected 1 ticket skip, got %d", len(skips))
	}
}

func TestGetTicketSkipsByUser(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/users/42/skips.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"skips":[{"id":1,"user_id":42}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	skips, _, err := client.GetTicketSkipsByUser(ctx, 42, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket skips by user: %s", err)
	}
	if len(skips) != 1 {
		t.Fatalf("Expected 1 ticket skip, got %d", len(skips))
	}
}

func TestCreateTicketSkip(t *testing.T) {
	var received struct {
		Skip TicketSkipOptions `json:"skip"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"skip":{"id":1,"ticket_id":123,"reason":"I have no idea."}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	skip, err := client.CreateTicketSkip(ctx, TicketSkipOptions{TicketID: 123, Reason: "I have no idea."})
	if err != nil {
		t.Fatalf("Failed to create ticket skip: %s", err)
	}
	if received.Skip.TicketID != 123 || received.Skip.Reason != "I have no idea." {
		t.Fatalf("Unexpected request payload %+v", received.Skip)
	}
	if skip.ID != 1 {
		t.Fatalf("Unexpected ticket skip id %d", skip.ID)
	}
}
