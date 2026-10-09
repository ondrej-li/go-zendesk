package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestGetUserFields(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "user_fields.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	fields, page, err := client.GetUserFields(ctx, nil)
	if err != nil {
		t.Fatalf("Received error calling API: %v", err)
	}

	if page.Count != 1 {
		t.Fatalf("Did not receive the correct count in the page field. Was %d expected 1", page.Count)
	}

	n := len(fields)
	if n != 1 {
		t.Fatalf("Expected 1 entry in fields list. Got %d", n)
	}

	id := fields[0].ID
	if id != 7 {
		t.Fatalf("Field did not have the expected id. Was %d", id)
	}
}

func TestUserFieldQueryParamsSet(t *testing.T) {
	opts := UserFieldListOptions{
		PageOptions{
			Page: 2,
		},
	}
	expected := fmt.Sprintf("page=%d", opts.Page)
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryString := r.URL.Query().Encode()
		if queryString != expected {
			t.Fatalf(`Did not get the expect query string: "%s". Was: "%s"`, expected, queryString)
		}
		w.Write(readFixture(filepath.Join(http.MethodGet, "user_fields.json")))
	}))

	defer mockAPI.Close()
	client := newTestClient(mockAPI)
	_, _, err := client.GetUserFields(ctx, &opts)
	if err != nil {
		t.Fatalf("Received error calling API: %v", err)
	}
}

func TestGetUserField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"user_field":{"id":7,"title":"Role"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	field, err := client.GetUserField(ctx, 7)
	if err != nil {
		t.Fatalf("Failed to get user field: %s", err)
	}
	if field.Title != "Role" {
		t.Fatalf("Unexpected user field title %q", field.Title)
	}
}

func TestUpdateUserField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"user_field":{"id":7,"title":"Updated"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	field, err := client.UpdateUserField(ctx, 7, UserField{Title: "Updated"})
	if err != nil {
		t.Fatalf("Failed to update user field: %s", err)
	}
	if field.Title != "Updated" {
		t.Fatalf("Unexpected user field title %q", field.Title)
	}
}

func TestDeleteUserField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if err := client.DeleteUserField(ctx, 7); err != nil {
		t.Fatalf("Failed to delete user field: %s", err)
	}
}

func TestCreateUserField(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if want := "/user_fields.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}

		var body struct {
			UserField UserField `json:"user_field"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		if body.UserField.Title != "Role" {
			t.Errorf(`Expected user field to be wrapped in "user_field", got %+v`, body)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"user_field":{"id":7,"title":"Role"}}`))
	}))
	defer mockAPI.Close()

	field, err := newTestClient(mockAPI).CreateUserField(ctx, UserField{Title: "Role"})
	if err != nil {
		t.Fatalf("Failed to create user field: %s", err)
	}
	if field.ID != 7 || field.Title != "Role" {
		t.Fatalf("Unexpected user field: %+v", field)
	}
}

func TestCreateUserFieldFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}))
	defer mockAPI.Close()

	if _, err := newTestClient(mockAPI).CreateUserField(ctx, UserField{}); err == nil {
		t.Fatal("Expected error, got nil")
	}
}
