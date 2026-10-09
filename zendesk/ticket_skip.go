package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TicketSkip is a record of an agent skipping a ticket
// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/
type TicketSkip struct {
	ID        int64      `json:"id,omitempty"`
	TicketID  int64      `json:"ticket_id,omitempty"`
	UserID    int64      `json:"user_id,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	Ticket    *Ticket    `json:"ticket,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// TicketSkipOptions is the payload for recording a ticket skip
type TicketSkipOptions struct {
	TicketID int64  `json:"ticket_id"`
	Reason   string `json:"reason,omitempty"`
}

// TicketSkipListOptions is options for listing ticket skips
type TicketSkipListOptions struct {
	PageOptions
	SortOrder string `url:"sort_order,omitempty"`
}

// TicketSkipAPI an interface containing all ticket skip related methods
type TicketSkipAPI interface {
	GetTicketSkips(ctx context.Context, opts *TicketSkipListOptions) ([]TicketSkip, Page, error)
	GetTicketSkipsByTicket(ctx context.Context, ticketID int64, opts *TicketSkipListOptions) ([]TicketSkip, Page, error)
	GetTicketSkipsByUser(ctx context.Context, userID int64, opts *TicketSkipListOptions) ([]TicketSkip, Page, error)
	CreateTicketSkip(ctx context.Context, opts TicketSkipOptions) (TicketSkip, error)
}

// GetTicketSkips fetches all ticket skips
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/#list-all-skips
func (z *Client) GetTicketSkips(ctx context.Context, opts *TicketSkipListOptions) ([]TicketSkip, Page, error) {
	return z.getTicketSkips(ctx, "/skips.json", opts)
}

// GetTicketSkipsByTicket fetches the skips for a ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/#list-ticket-skips-by-ticket
func (z *Client) GetTicketSkipsByTicket(ctx context.Context, ticketID int64, opts *TicketSkipListOptions) ([]TicketSkip, Page, error) {
	return z.getTicketSkips(ctx, fmt.Sprintf("/tickets/%d/skips.json", ticketID), opts)
}

// GetTicketSkipsByUser fetches the skips for a user
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/#list-ticket-skips-by-user
func (z *Client) GetTicketSkipsByUser(ctx context.Context, userID int64, opts *TicketSkipListOptions) ([]TicketSkip, Page, error) {
	return z.getTicketSkips(ctx, fmt.Sprintf("/users/%d/skips.json", userID), opts)
}

func (z *Client) getTicketSkips(ctx context.Context, path string, opts *TicketSkipListOptions) ([]TicketSkip, Page, error) {
	var data struct {
		Skips []TicketSkip `json:"skips"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketSkipListOptions{}
	}

	u, err := addOptions(path, tmp)
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
	return data.Skips, data.Page, nil
}

// CreateTicketSkip records a new ticket skip for the current user
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/#record-a-new-skip-for-the-current-user
func (z *Client) CreateTicketSkip(ctx context.Context, opts TicketSkipOptions) (TicketSkip, error) {
	var data struct {
		Skip TicketSkipOptions `json:"skip"`
	}
	var result struct {
		Skip TicketSkip `json:"skip"`
	}
	data.Skip = opts

	body, err := z.post(ctx, "/skips.json", data)
	if err != nil {
		return TicketSkip{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return TicketSkip{}, err
	}
	return result.Skip, nil
}
