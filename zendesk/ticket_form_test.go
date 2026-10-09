package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketForms(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_forms.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketForms, _, err := client.GetTicketForms(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket forms: %s", err)
	}

	if len(ticketForms) != 1 {
		t.Fatalf("expected length of ticket forms is , but got %d", len(ticketForms))
	}
}

func TestCreateTicketForm(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPost, "ticket_form.json", http.StatusCreated)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.CreateTicketForm(ctx, TicketForm{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket form: %s", err)
	}
}

func TestDeleteTicketForm(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	err := c.DeleteTicketForm(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete ticket field: %s", err)
	}
}

func TestDeleteTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	err := c.DeleteTicketForm(ctx, 1234)
	if err == nil {
		t.Fatal("Client did not return error when api failed")
	}
}

func TestGetTicketForm(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_form.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	f, err := client.GetTicketForm(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(47)
	if f.ID != expectedID {
		t.Fatalf("Returned ticket form does not have the expected ID %d. Ticket id is %d", expectedID, f.ID)
	}
}

func TestGetTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	_, err := c.GetTicketForm(ctx, 1234)
	if err == nil {
		t.Fatal("Client did not return error when api failed")
	}
}

func TestUpdateTicketForm(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPut, "ticket_form.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	f, err := client.UpdateTicketForm(ctx, 123, TicketForm{})
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(47)
	if f.ID != expectedID {
		t.Fatalf("Returned ticket form does not have the expected ID %d. Ticket id is %d", expectedID, f.ID)
	}
}

func TestUpdateTicketFormFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	_, err := c.UpdateTicketForm(ctx, 1234, TicketForm{})
	if err == nil {
		t.Fatal("Client did not return error when api failed")
	}
}

func TestGetTicketFormsShowMany(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2" {
			t.Errorf("Expected ids query parameter %q, got %q", "1,2", got)
		}
		w.Write([]byte(`{"ticket_forms":[{"id":1},{"id":2}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	forms, err := client.GetTicketFormsShowMany(ctx, []int64{1, 2})
	if err != nil {
		t.Fatalf("Failed to get ticket forms: %s", err)
	}
	if len(forms) != 2 {
		t.Fatalf("Expected 2 ticket forms, got %d", len(forms))
	}
}

func TestReorderTicketForms(t *testing.T) {
	var received struct {
		TicketFormIDs []int64 `json:"ticket_form_ids"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(`{"ticket_forms":[{"id":2},{"id":23}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	forms, err := client.ReorderTicketForms(ctx, []int64{2, 23})
	if err != nil {
		t.Fatalf("Failed to reorder ticket forms: %s", err)
	}
	if len(received.TicketFormIDs) != 2 || received.TicketFormIDs[0] != 2 {
		t.Fatalf("Unexpected request ids %v", received.TicketFormIDs)
	}
	if len(forms) != 2 {
		t.Fatalf("Expected 2 ticket forms in the response, got %d", len(forms))
	}
}
