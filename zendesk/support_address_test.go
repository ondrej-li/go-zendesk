package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSupportAddresses(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"recipient_addresses":[{"id":1,"email":"support@example.zendesk.com"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	addresses, err := client.GetSupportAddresses(ctx)
	if err != nil {
		t.Fatalf("Failed to get support addresses: %s", err)
	}
	if len(addresses) != 1 {
		t.Fatalf("Expected 1 support address, got %d", len(addresses))
	}
}

func TestGetSupportAddress(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"recipient_address":{"id":35436,"email":"support@example.zendesk.com"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	address, err := client.GetSupportAddress(ctx, 35436)
	if err != nil {
		t.Fatalf("Failed to get support address: %s", err)
	}
	if address.Email != "support@example.zendesk.com" {
		t.Fatalf("Unexpected support address email %q", address.Email)
	}
}

func TestCreateSupportAddress(t *testing.T) {
	var received struct {
		RecipientAddress SupportAddress `json:"recipient_address"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"recipient_address":{"id":1,"email":"sales@example.zendesk.com"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	address, err := client.CreateSupportAddress(ctx, SupportAddress{Email: "sales@example.zendesk.com"})
	if err != nil {
		t.Fatalf("Failed to create support address: %s", err)
	}
	if received.RecipientAddress.Email != "sales@example.zendesk.com" {
		t.Fatalf("Unexpected request email %q", received.RecipientAddress.Email)
	}
	if address.ID != 1 {
		t.Fatalf("Unexpected support address id %d", address.ID)
	}
}

func TestUpdateSupportAddress(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"recipient_address":{"id":1,"name":"Updated"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	address, err := client.UpdateSupportAddress(ctx, 1, SupportAddress{Name: "Updated"})
	if err != nil {
		t.Fatalf("Failed to update support address: %s", err)
	}
	if address.Name != "Updated" {
		t.Fatalf("Unexpected support address name %q", address.Name)
	}
}

func TestDeleteSupportAddress(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.DeleteSupportAddress(ctx, 1); err != nil {
		t.Fatalf("Failed to delete support address: %s", err)
	}
}

func TestVerifySupportAddressForwarding(t *testing.T) {
	var received struct {
		Type string `json:"type"`
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
	if err := client.VerifySupportAddressForwarding(ctx, 1, ""); err != nil {
		t.Fatalf("Failed to verify support address: %s", err)
	}
	if received.Type != "forwarding" {
		t.Fatalf("Expected default verification type %q, got %q", "forwarding", received.Type)
	}
}
