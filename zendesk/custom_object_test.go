package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCustomObjectRecord(t *testing.T) {
	var received struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if want := "/custom_objects/car/records.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"custom_object_record":{"id":"1","name":"Tesla","custom_object_key":"car"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	record, err := client.CreateCustomObjectRecord(ctx, CustomObjectRecord{Name: "Tesla"}, "car")
	if err != nil {
		t.Fatalf("Failed to create custom object record: %s", err)
	}
	if received.CustomObjectRecord.Name != "Tesla" {
		t.Fatalf("Unexpected request payload %+v", received.CustomObjectRecord)
	}
	if record.ID != "1" {
		t.Fatalf("Unexpected record id %q", record.ID)
	}
}

func TestCreateCustomObjectRecordFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"RecordInvalid"}`, http.StatusUnprocessableEntity)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.CreateCustomObjectRecord(ctx, CustomObjectRecord{}, "car"); err == nil {
		t.Fatal("Expected an error when the API fails")
	}
}

func TestListCustomObjectRecords(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/custom_objects/car/records"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if got := r.URL.Query().Get("filter[ids]"); got != "1,2" {
			t.Errorf("Expected filter[ids] query parameter %q, got %q", "1,2", got)
		}
		w.Write([]byte(`{"custom_object_records":[{"id":"1","name":"Tesla"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	records, _, err := client.ListCustomObjectRecords(ctx, "car", &CustomObjectListOptions{Ids: "1,2"})
	if err != nil {
		t.Fatalf("Failed to list custom object records: %s", err)
	}
	if len(records) != 1 || records[0].Name != "Tesla" {
		t.Fatalf("Unexpected records %+v", records)
	}
}

func TestListCustomObjectRecordsFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, _, err := client.ListCustomObjectRecords(ctx, "car", nil); err == nil {
		t.Fatal("Expected an error when the API fails")
	}
}

func TestAutocompleteSearchCustomObjectRecords(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/custom_objects/car/records/autocomplete"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if got := r.URL.Query().Get("name"); got != "Tes" {
			t.Errorf("Expected name query parameter %q, got %q", "Tes", got)
		}
		w.Write([]byte(`{"custom_object_records":[{"id":"1","name":"Tesla"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	records, _, err := client.AutocompleteSearchCustomObjectRecords(ctx, "car", &CustomObjectAutocompleteOptions{Name: "Tes"})
	if err != nil {
		t.Fatalf("Failed to autocomplete custom object records: %s", err)
	}
	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}
}

func TestAutocompleteSearchCustomObjectRecordsFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, _, err := client.AutocompleteSearchCustomObjectRecords(ctx, "car", nil); err == nil {
		t.Fatal("Expected an error when the API fails")
	}
}

func TestSearchCustomObjectRecords(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/custom_objects/car/records/search"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if got := r.URL.Query().Get("query"); got != "name:Tesla" {
			t.Errorf("Expected query parameter %q, got %q", "name:Tesla", got)
		}
		if got := r.URL.Query().Get("sort"); got != "-name" {
			t.Errorf("Expected sort parameter %q, got %q", "-name", got)
		}
		w.Write([]byte(`{"custom_object_records":[{"id":"1","name":"Tesla"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	records, _, err := client.SearchCustomObjectRecords(ctx, "car", &SearchCustomObjectRecordsOptions{
		Query: "name:Tesla",
		Sort:  "-name",
	})
	if err != nil {
		t.Fatalf("Failed to search custom object records: %s", err)
	}
	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}
}

func TestSearchCustomObjectRecordsFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, _, err := client.SearchCustomObjectRecords(ctx, "car", nil); err == nil {
		t.Fatal("Expected an error when the API fails")
	}
}

func TestShowCustomObjectRecord(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/custom_objects/car/records/1"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"custom_object_record":{"id":"1","name":"Tesla"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	record, err := client.ShowCustomObjectRecord(ctx, "car", "1")
	if err != nil {
		t.Fatalf("Failed to show custom object record: %s", err)
	}
	if record.Name != "Tesla" {
		t.Fatalf("Unexpected record %+v", record)
	}
}

func TestShowCustomObjectRecordFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"RecordNotFound"}`, http.StatusNotFound)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.ShowCustomObjectRecord(ctx, "car", "404"); err == nil {
		t.Fatal("Expected an error when the record is missing")
	}
}

func TestUpdateCustomObjectRecord(t *testing.T) {
	var received struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		if want := "/custom_objects/car/records/1"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(`{"custom_object_record":{"id":"1","name":"Roadster"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	record, err := client.UpdateCustomObjectRecord(ctx, "car", "1", CustomObjectRecord{Name: "Roadster"})
	if err != nil {
		t.Fatalf("Failed to update custom object record: %s", err)
	}
	if received.CustomObjectRecord.Name != "Roadster" {
		t.Fatalf("Unexpected request payload %+v", received.CustomObjectRecord)
	}
	if record.Name != "Roadster" {
		t.Fatalf("Unexpected record %+v", record)
	}
}

func TestUpdateCustomObjectRecordFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	if _, err := client.UpdateCustomObjectRecord(ctx, "car", "1", CustomObjectRecord{}); err == nil {
		t.Fatal("Expected an error when the API fails")
	}
}
