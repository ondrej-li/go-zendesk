package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketTags(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "tags.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := client.GetTicketTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get ticket tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestGetOrganizationTags(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "tags.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := client.GetOrganizationTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get organization tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestGetUserTags(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "tags.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tags, err := client.GetUserTags(ctx, int64(2))
	if err != nil {
		t.Fatalf("Failed to get user tags: %s", err)
	}

	expectedLength := 2
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
}

func TestAddTicketTags(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPut, "tags.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := client.AddTicketTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}

func TestAddOrganizationTags(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPut, "tags.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := client.AddOrganizationTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}

func TestAddUserTags(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPut, "tags.json", http.StatusOK)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	tag := Tag("example")

	tags, err := client.AddUserTags(ctx, 2, []Tag{tag})
	if err != nil {
		t.Fatalf("Failed to add ticket tags: %s", err)
	}

	expectedLength := 3
	if len(tags) != expectedLength {
		t.Fatalf("Returned tags does not have the expexted length %d. Tags length is %d", expectedLength, len(tags))
	}
	if tags[2] != tag {
		t.Fatalf("Returned tags does not have the expexted tag %s. %s given", "important", tags[0])
	}
}

func TestSetTicketTags(t *testing.T) {
	var received struct {
		Tags []Tag `json:"tags"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"tags":["important","customer","urgent"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	tags, err := client.SetTicketTags(ctx, 2, []Tag{"important", "customer"})
	if err != nil {
		t.Fatalf("Failed to set ticket tags: %s", err)
	}
	if len(received.Tags) != 2 {
		t.Fatalf("Expected 2 tags in the request, got %d", len(received.Tags))
	}
	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags in the response, got %d", len(tags))
	}
}

func TestSetOrganizationTags(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"tags":["important"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.SetOrganizationTags(ctx, 16, []Tag{"important"}); err != nil {
		t.Fatalf("Failed to set organization tags: %s", err)
	}
}

func TestSetUserTags(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"tags":["important"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.SetUserTags(ctx, 42, []Tag{"important"}); err != nil {
		t.Fatalf("Failed to set user tags: %s", err)
	}
}

func TestDeleteTicketTags(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if got := r.URL.Query().Get("tags"); got != "important,customer" {
			t.Errorf("Expected tags query parameter %q, got %q", "important,customer", got)
		}
		w.Write([]byte(`{"tags":["urgent"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	tags, err := client.DeleteTicketTags(ctx, 2, []Tag{"important", "customer"})
	if err != nil {
		t.Fatalf("Failed to delete ticket tags: %s", err)
	}
	if len(tags) != 1 || tags[0] != "urgent" {
		t.Fatalf("Unexpected remaining tags %v", tags)
	}
}

func TestDeleteOrganizationTags(t *testing.T) {
	var received struct {
		Tags []Tag `json:"tags"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(`{"tags":["urgent"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.DeleteOrganizationTags(ctx, 16, []Tag{"customer"}); err != nil {
		t.Fatalf("Failed to delete organization tags: %s", err)
	}
	if len(received.Tags) != 1 || received.Tags[0] != "customer" {
		t.Fatalf("Unexpected request tags %v", received.Tags)
	}
}

func TestDeleteUserTags(t *testing.T) {
	var received struct {
		Tags []Tag `json:"tags"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(`{"tags":["urgent"]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.DeleteUserTags(ctx, 42, []Tag{"customer"}); err != nil {
		t.Fatalf("Failed to delete user tags: %s", err)
	}
	if len(received.Tags) != 1 || received.Tags[0] != "customer" {
		t.Fatalf("Unexpected request tags %v", received.Tags)
	}
}
