package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketMetrics(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"ticket_metrics":[{"id":1,"ticket_id":208,"replies":2}],"count":1}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	metrics, _, err := client.GetTicketMetrics(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get ticket metrics: %s", err)
	}
	if len(metrics) != 1 || metrics[0].TicketID != 208 {
		t.Fatalf("Unexpected ticket metrics %+v", metrics)
	}
}

func TestGetTicketMetric(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/ticket_metrics/1.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"ticket_metric":{"id":1,"ticket_id":208}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	metric, err := client.GetTicketMetric(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get ticket metric: %s", err)
	}
	if metric.ID != 1 {
		t.Fatalf("Unexpected ticket metric id %d", metric.ID)
	}
}

func TestGetTicketMetricByTicket(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/tickets/208/metrics.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"ticket_metric":{"id":1,"ticket_id":208}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	metric, err := client.GetTicketMetricByTicket(ctx, 208)
	if err != nil {
		t.Fatalf("Failed to get ticket metric by ticket: %s", err)
	}
	if metric.TicketID != 208 {
		t.Fatalf("Unexpected ticket id %d", metric.TicketID)
	}
}

// TestTicketMetricsAvailableViaAPI guards the interface embedding: these methods
// were unreachable through zendesk.API before TicketMetricsAPI was embedded.
func TestTicketMetricsAvailableViaAPI(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ticket_metrics":[]}`))
	}))
	defer mockAPI.Close()

	var api API = newTestClient(mockAPI)
	if _, _, err := api.GetTicketMetrics(ctx, nil); err != nil {
		t.Fatalf("Failed to get ticket metrics via API: %s", err)
	}
}
