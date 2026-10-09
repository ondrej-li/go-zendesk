package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrganizationFields(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "organization_fields.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	ticketFields, _, err := client.GetOrganizationFields(ctx)
	if err != nil {
		t.Fatalf("Failed to get organization fields: %s", err)
	}

	if len(ticketFields) != 2 {
		t.Fatalf("expected length of organization fields is , but got %d", len(ticketFields))
	}
}

func TestOrganizationField(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPost, "organization_fields.json", http.StatusCreated)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := client.CreateOrganizationField(ctx, OrganizationField{})
	if err != nil {
		t.Fatalf("Failed to send request to create organization field: %s", err)
	}
}

func TestGetOrganizationField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"organization_field":{"id":1,"title":"Plan"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	field, err := client.GetOrganizationField(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get organization field: %s", err)
	}
	if field.Title != "Plan" {
		t.Fatalf("Unexpected organization field title %q", field.Title)
	}
}

func TestUpdateOrganizationField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"organization_field":{"id":1,"title":"Updated"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	field, err := client.UpdateOrganizationField(ctx, 1, OrganizationField{Title: "Updated"})
	if err != nil {
		t.Fatalf("Failed to update organization field: %s", err)
	}
	if field.Title != "Updated" {
		t.Fatalf("Unexpected organization field title %q", field.Title)
	}
}

func TestDeleteOrganizationField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.DeleteOrganizationField(ctx, 1); err != nil {
		t.Fatalf("Failed to delete organization field: %s", err)
	}
}
