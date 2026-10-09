package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDynamicContentItems(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "dynamic_content/items.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	items, page, err := client.GetDynamicContentItems(ctx)
	if err != nil {
		t.Fatalf("Failed to get dynamic content items: %s", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected length of dynamic content items is 2, but got %d", len(items))
	}

	if len(items[0].Variants) != 3 {
		t.Fatalf("expected length of items[0].Variants is 3, but got %d", len(items[0].Variants))
	}

	if page.HasPrev() || page.HasNext() {
		t.Fatalf("page fields are wrong: %v", page)
	}
}

func TestCreateDynamicContentItem(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPost, "dynamic_content/items.json", http.StatusCreated)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	item, err := client.CreateDynamicContentItem(ctx, DynamicContentItem{})
	if err != nil {
		t.Fatalf("Failed to get valid response: %s", err)
	}
	if item.ID == 0 {
		t.Fatal("Failed to create dynamic content item")
	}
}

func TestGetDynamicContentItem(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if want := "/dynamic_content/items/360000346774.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.Write([]byte(`{"item":{"id":360000346774,"name":"title","placeholder":"{{dc.title}}","default_locale_id":1,
			"variants":[{"id":360001100153,"content":"hello","locale_id":1}]}}`))
	}))
	defer mockAPI.Close()

	item, err := newTestClient(mockAPI).GetDynamicContentItem(ctx, 360000346774)
	if err != nil {
		t.Fatalf("Failed to get dynamic content item: %s", err)
	}

	if item.ID != 360000346774 {
		t.Fatalf("expected item id 360000346774, but got %d", item.ID)
	}
	if item.Name != "title" || item.Placeholder != "{{dc.title}}" || item.DefaultLocaleID != 1 {
		t.Fatalf("Unexpected dynamic content item: %+v", item)
	}
	if len(item.Variants) != 1 || item.Variants[0].Content != "hello" {
		t.Fatalf("Unexpected variants: %+v", item.Variants)
	}
}

func TestUpdateDynamicContentItem(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		if want := "/dynamic_content/items/360000346774.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}

		var body struct {
			Item DynamicContentItem `json:"item"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		if body.Item.Name != "updated" {
			t.Errorf(`Expected item to be wrapped in "item", got %+v`, body)
		}

		w.Write([]byte(`{"item":{"id":360000346774,"name":"updated","default_locale_id":1}}`))
	}))
	defer mockAPI.Close()

	item, err := newTestClient(mockAPI).UpdateDynamicContentItem(ctx, 360000346774, DynamicContentItem{Name: "updated"})
	if err != nil {
		t.Fatalf("Failed to update dynamic content item: %s", err)
	}
	if item.Name != "updated" {
		t.Fatalf("Unexpected dynamic content item name %q", item.Name)
	}
}

func TestDeleteDynamicContentItem(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if want := "/dynamic_content/items/360000346774.json"; r.URL.Path != want {
			t.Errorf("Expected path %q, got %q", want, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockAPI.Close()

	if err := newTestClient(mockAPI).DeleteDynamicContentItem(ctx, 360000346774); err != nil {
		t.Fatalf("Failed to delete dynamic content item: %s", err)
	}
}

func TestDynamicContentItemFailures(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) error
	}{
		{"show", func(c *Client) error { _, err := c.GetDynamicContentItem(ctx, 1); return err }},
		{"update", func(c *Client) error { _, err := c.UpdateDynamicContentItem(ctx, 1, DynamicContentItem{}); return err }},
		{"delete", func(c *Client) error { return c.DeleteDynamicContentItem(ctx, 1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer mockAPI.Close()

			if err := tt.call(newTestClient(mockAPI)); err == nil {
				t.Fatal("Expected error, got nil")
			}
		})
	}
}
