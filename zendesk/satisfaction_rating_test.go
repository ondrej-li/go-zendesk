package zendesk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSatisfactionRatings(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"satisfaction_ratings":[{"id":1,"score":"good","ticket_id":208}],"count":1}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	ratings, _, err := client.GetSatisfactionRatings(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get satisfaction ratings: %s", err)
	}
	if len(ratings) != 1 || ratings[0].Score != "good" {
		t.Fatalf("Unexpected satisfaction ratings %+v", ratings)
	}
}

func TestGetSatisfactionRatingsCount(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"count":{"refreshed_at":"2020-04-06T02:18:17Z","value":102}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	count, err := client.GetSatisfactionRatingsCount(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get satisfaction ratings count: %s", err)
	}
	if count.Value != 102 {
		t.Fatalf("Unexpected count value %d", count.Value)
	}
}

func TestGetSatisfactionRating(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"satisfaction_rating":[{"id":35436,"score":"good"}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	rating, err := client.GetSatisfactionRating(ctx, 35436)
	if err != nil {
		t.Fatalf("Failed to get satisfaction rating: %s", err)
	}
	if rating.ID != 35436 {
		t.Fatalf("Unexpected satisfaction rating id %d", rating.ID)
	}
}

func TestGetSatisfactionRatingObjectShape(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"satisfaction_rating":{"id":35436,"score":"bad"}}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	rating, err := client.GetSatisfactionRating(ctx, 35436)
	if err != nil {
		t.Fatalf("Failed to get satisfaction rating: %s", err)
	}
	if rating.Score != "bad" {
		t.Fatalf("Unexpected satisfaction rating score %q", rating.Score)
	}
}

func TestCreateTicketSatisfactionRating(t *testing.T) {
	var received struct {
		SatisfactionRating SatisfactionRating `json:"satisfaction_rating"`
	}

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode request body: %s", err)
		}
		w.Write([]byte(`{"satisfaction_rating":[{"id":35436,"score":"good","ticket_id":208}]}`))
	}))
	defer mockAPI.Close()

	client := newTestClient(mockAPI)
	rating, err := client.CreateTicketSatisfactionRating(ctx, 208, SatisfactionRating{Score: "good", Comment: "Awesome"})
	if err != nil {
		t.Fatalf("Failed to create ticket satisfaction rating: %s", err)
	}
	if received.SatisfactionRating.Score != "good" || received.SatisfactionRating.Comment != "Awesome" {
		t.Fatalf("Unexpected request payload %+v", received.SatisfactionRating)
	}
	if rating.ID != 35436 {
		t.Fatalf("Unexpected satisfaction rating id %d", rating.ID)
	}
}
