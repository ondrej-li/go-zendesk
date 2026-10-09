package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCustomRoles(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"custom_roles":[{"id":1,"name":"Agent"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	roles, err := client.GetCustomRoles(ctx)
	if err != nil {
		t.Fatalf("Failed to get custom roles: %s", err)
	}
	if len(roles) != 1 {
		t.Fatalf("Expected 1 custom role, got %d", len(roles))
	}
}

func TestGetCustomRole(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"custom_role":{"id":1,"name":"Agent"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	role, err := client.GetCustomRole(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get custom role: %s", err)
	}
	if role.Name != "Agent" {
		t.Fatalf("Unexpected custom role name %q", role.Name)
	}
}

func TestCreateCustomRole(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Write([]byte(`{"custom_role":{"id":2,"name":"Created"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	role, err := client.CreateCustomRole(ctx, CustomRole{Name: "Created"})
	if err != nil {
		t.Fatalf("Failed to create custom role: %s", err)
	}
	if role.ID != 2 {
		t.Fatalf("Unexpected custom role id %d", role.ID)
	}
}

func TestUpdateCustomRole(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"custom_role":{"id":2,"name":"Updated"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	role, err := client.UpdateCustomRole(ctx, 2, CustomRole{Name: "Updated"})
	if err != nil {
		t.Fatalf("Failed to update custom role: %s", err)
	}
	if role.Name != "Updated" {
		t.Fatalf("Unexpected custom role name %q", role.Name)
	}
}

func TestDeleteCustomRole(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.DeleteCustomRole(ctx, 2); err != nil {
		t.Fatalf("Failed to delete custom role: %s", err)
	}
}
