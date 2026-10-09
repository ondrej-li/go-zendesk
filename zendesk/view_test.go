package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetView(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "view.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	view, err := client.GetView(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get view: %s", err)
	}

	expectedID := int64(360002440594)
	if view.ID != expectedID {
		t.Fatalf("Returned view does not have the expected ID %d. View ID is %d", expectedID, view.ID)
	}
}

func TestGetViews(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "views.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	views, _, err := client.GetViews(ctx)
	if err != nil {
		t.Fatalf("Failed to get views: %s", err)
	}

	if len(views) != 2 {
		t.Fatalf("expected length of views is 2, but got %d", len(views))
	}
}

func TestGetCountTicketsInViewsTestGetViews(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "views_ticket_count.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()
	ids := []string{"25", "78"}
	viewsCount, err := client.GetCountTicketsInViews(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to get views tickets count: %s", err)
	}

	if len(viewsCount) != 2 {
		t.Fatalf("expected length of views ticket counts is 2, but got %d", len(viewsCount))
	}
}

func TestGetTicketsFromView(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/views/123/tickets.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("Expected page=2, got %q", got)
		}
		w.Write([]byte(`{"tickets":[{"id":1,"subject":"first"},{"id":2,"subject":"second"}],"count":2}`))
	}))
	defer mockAPI.Close()

	tickets, page, err := newTestClient(mockAPI).GetTicketsFromView(ctx, 123, &TicketListOptions{
		PageOptions: PageOptions{Page: 2},
	})
	if err != nil {
		t.Fatalf("Failed to get tickets from view: %s", err)
	}

	if len(tickets) != 2 {
		t.Fatalf("expected 2 tickets, but got %d", len(tickets))
	}
	if tickets[0].Subject != "first" {
		t.Fatalf("Unexpected ticket subject %q", tickets[0].Subject)
	}
	if page.Count != 2 {
		t.Fatalf("expected page count 2, but got %d", page.Count)
	}
}

func TestGetTicketsFromViewWithoutOptions(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if raw := r.URL.RawQuery; raw != "" {
			t.Errorf("Expected no query string, got %q", raw)
		}
		w.Write([]byte(`{"tickets":[]}`))
	}))
	defer mockAPI.Close()

	tickets, _, err := newTestClient(mockAPI).GetTicketsFromView(ctx, 123, nil)
	if err != nil {
		t.Fatalf("Failed to get tickets from view: %s", err)
	}
	if len(tickets) != 0 {
		t.Fatalf("expected no tickets, but got %d", len(tickets))
	}
}

func TestGetTicketsFromViewFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockAPI.Close()

	if _, _, err := newTestClient(mockAPI).GetTicketsFromView(ctx, 123, nil); err == nil {
		t.Fatal("Expected error, got nil")
	}
}
