package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SatisfactionRating is a customer satisfaction (CSAT) rating on a ticket
// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/
type SatisfactionRating struct {
	ID          int64      `json:"id,omitempty"`
	URL         string     `json:"url,omitempty"`
	Score       string     `json:"score"`
	Comment     string     `json:"comment,omitempty"`
	Reason      string     `json:"reason,omitempty"`
	ReasonID    int64      `json:"reason_id,omitempty"`
	ReasonCode  int64      `json:"reason_code,omitempty"`
	TicketID    int64      `json:"ticket_id,omitempty"`
	AssigneeID  int64      `json:"assignee_id,omitempty"`
	GroupID     int64      `json:"group_id,omitempty"`
	RequesterID int64      `json:"requester_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// SatisfactionRatingCount is an approximate count of satisfaction ratings
type SatisfactionRatingCount struct {
	RefreshedAt *time.Time `json:"refreshed_at,omitempty"`
	Value       int64      `json:"value"`
}

// SatisfactionRatingListOptions is options for listing or counting satisfaction ratings
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#list-satisfaction-ratings
type SatisfactionRatingListOptions struct {
	PageOptions
	Score     string `url:"score,omitempty"`
	StartTime int64  `url:"start_time,omitempty"`
	EndTime   int64  `url:"end_time,omitempty"`
	Sort      string `url:"sort,omitempty"`
}

// SatisfactionRatingAPI an interface containing all satisfaction rating related methods
type SatisfactionRatingAPI interface {
	GetSatisfactionRatings(ctx context.Context, opts *SatisfactionRatingListOptions) ([]SatisfactionRating, Page, error)
	GetSatisfactionRatingsCount(ctx context.Context, opts *SatisfactionRatingListOptions) (SatisfactionRatingCount, error)
	GetSatisfactionRating(ctx context.Context, satisfactionRatingID int64) (SatisfactionRating, error)
	CreateTicketSatisfactionRating(ctx context.Context, ticketID int64, rating SatisfactionRating) (SatisfactionRating, error)
}

// GetSatisfactionRatings fetches satisfaction rating list
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#list-satisfaction-ratings
func (z *Client) GetSatisfactionRatings(ctx context.Context, opts *SatisfactionRatingListOptions) ([]SatisfactionRating, Page, error) {
	var data struct {
		SatisfactionRatings []SatisfactionRating `json:"satisfaction_ratings"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &SatisfactionRatingListOptions{}
	}

	u, err := addOptions("/satisfaction_ratings.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, Page{}, err
	}
	return data.SatisfactionRatings, data.Page, nil
}

// GetSatisfactionRatingsCount returns an approximate count of satisfaction ratings
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#count-satisfaction-ratings
func (z *Client) GetSatisfactionRatingsCount(ctx context.Context, opts *SatisfactionRatingListOptions) (SatisfactionRatingCount, error) {
	var result struct {
		Count SatisfactionRatingCount `json:"count"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &SatisfactionRatingListOptions{}
	}

	u, err := addOptions("/satisfaction_ratings/count.json", tmp)
	if err != nil {
		return SatisfactionRatingCount{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return SatisfactionRatingCount{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SatisfactionRatingCount{}, err
	}
	return result.Count, nil
}

// GetSatisfactionRating gets a specified satisfaction rating
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#show-satisfaction-rating
func (z *Client) GetSatisfactionRating(ctx context.Context, satisfactionRatingID int64) (SatisfactionRating, error) {
	var result struct {
		SatisfactionRating json.RawMessage `json:"satisfaction_rating"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/satisfaction_ratings/%d.json", satisfactionRatingID))
	if err != nil {
		return SatisfactionRating{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SatisfactionRating{}, err
	}
	return decodeSatisfactionRating(result.SatisfactionRating)
}

// CreateTicketSatisfactionRating creates a satisfaction rating for a ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#create-a-satisfaction-rating
func (z *Client) CreateTicketSatisfactionRating(ctx context.Context, ticketID int64, rating SatisfactionRating) (SatisfactionRating, error) {
	var data struct {
		SatisfactionRating SatisfactionRating `json:"satisfaction_rating"`
	}
	var result struct {
		SatisfactionRating json.RawMessage `json:"satisfaction_rating"`
	}
	data.SatisfactionRating = rating

	body, err := z.post(ctx, fmt.Sprintf("/tickets/%d/satisfaction_rating.json", ticketID), data)
	if err != nil {
		return SatisfactionRating{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SatisfactionRating{}, err
	}
	return decodeSatisfactionRating(result.SatisfactionRating)
}

// decodeSatisfactionRating decodes a satisfaction_rating payload. The API documents
// the show and create responses as a one-element array, so both shapes are accepted.
func decodeSatisfactionRating(raw json.RawMessage) (SatisfactionRating, error) {
	if len(raw) == 0 {
		return SatisfactionRating{}, nil
	}

	if raw[0] == '[' {
		var ratings []SatisfactionRating
		if err := json.Unmarshal(raw, &ratings); err != nil {
			return SatisfactionRating{}, err
		}
		if len(ratings) == 0 {
			return SatisfactionRating{}, nil
		}
		return ratings[0], nil
	}

	var rating SatisfactionRating
	if err := json.Unmarshal(raw, &rating); err != nil {
		return SatisfactionRating{}, err
	}
	return rating, nil
}
