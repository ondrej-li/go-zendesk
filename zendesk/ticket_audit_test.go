package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTicketAudits(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_audits.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketAudits, _, err := client.GetAllTicketAudits(ctx, CursorOption{})
	if err != nil {
		t.Fatalf("Failed to get ticket audits: %s", err)
	}

	if len(ticketAudits) != 1 {
		t.Fatalf("expected length of ticket audit is %d, but got %d", 1, len(ticketAudits))
	}
}

func TestGetTicketAudits(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_audits.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketAudits, _, err := client.GetTicketAudits(ctx, 666, PageOptions{})
	if err != nil {
		t.Fatalf("Failed to get ticket audits: %s", err)
	}

	if len(ticketAudits) != 1 {
		t.Fatalf("expected length of ticket audit is %d, but got %d", 1, len(ticketAudits))
	}
}

func TestGetTicketAudit(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "ticket_audit.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketAudit, err := client.GetTicketAudit(ctx, 666, 2127301143)
	if err != nil {
		t.Fatalf("Failed to get ticket audit: %s", err)
	}

	expectedID := int64(2127301143)
	if ticketAudit.ID != expectedID {
		t.Fatalf("Returned ticket audit does not have the expected ID %d. Ticket audit id is %d", expectedID, ticketAudit.ID)
	}
}

func TestGetTicketAuditsCount(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"count":{"refreshed_at":"2020-04-06T02:18:17Z","value":5}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	count, err := client.GetTicketAuditsCount(ctx, 2)
	if err != nil {
		t.Fatalf("Failed to get ticket audits count: %s", err)
	}
	if count.Value != 5 {
		t.Fatalf("Unexpected count value %d", count.Value)
	}
}

func TestMakeTicketAuditPrivate(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(``))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.MakeTicketAuditPrivate(ctx, 2, 2127301143); err != nil {
		t.Fatalf("Failed to make ticket audit private: %s", err)
	}
}
