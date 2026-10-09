package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketFields(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_fields.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketFields, _, err := client.GetTicketFields(ctx)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	if len(ticketFields) != 15 {
		t.Fatalf("expected length of ticket fields is , but got %d", len(ticketFields))
	}
}

func TestGetTicketField(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_field.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketField, err := client.GetTicketField(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}

	expectedID := int64(360011737434)
	if ticketField.ID != expectedID {
		t.Fatalf("Returned ticket field does not have the expected ID %d. Ticket id is %d", expectedID, ticketField.ID)
	}
}

func TestCreateTicketField(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPost, "ticket_fields.json", http.StatusCreated)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.CreateTicketField(ctx, TicketField{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket field: %s", err)
	}
}

func TestUpdateTicketField(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPut, "ticket_field.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	updatedField, err := client.UpdateTicketField(ctx, int64(1234), TicketField{})
	if err != nil {
		t.Fatalf("Failed to send request to create ticket field: %s", err)
	}

	expectedID := int64(360011737434)
	if updatedField.ID != expectedID {
		t.Fatalf("Updated field %v did not have expected id %d", updatedField, expectedID)
	}
}

func TestDeleteTicketField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	err := c.DeleteTicketField(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete ticket field: %s", err)
	}
}

func TestGetTicketFieldsShowMany(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if got := r.URL.Query().Get("ids"); got != "1,2,3" {
			t.Errorf("Expected ids query parameter %q, got %q", "1,2,3", got)
		}
		w.Write([]byte(`{"ticket_fields":[{"id":1},{"id":2}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	fields, err := client.GetTicketFieldsShowMany(ctx, []int64{1, 2, 3})
	if err != nil {
		t.Fatalf("Failed to get ticket fields: %s", err)
	}
	if len(fields) != 2 {
		t.Fatalf("Expected 2 ticket fields, got %d", len(fields))
	}
}

func TestReorderTicketFields(t *testing.T) {
	var received struct {
		TicketFieldIDs []int64 `json:"ticket_field_ids"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(``))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.ReorderTicketFields(ctx, []int64{2, 23, 46}); err != nil {
		t.Fatalf("Failed to reorder ticket fields: %s", err)
	}
	if len(received.TicketFieldIDs) != 3 || received.TicketFieldIDs[0] != 2 {
		t.Fatalf("Unexpected request ids %v", received.TicketFieldIDs)
	}
}
