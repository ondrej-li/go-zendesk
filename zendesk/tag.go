package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Tag is an alias for string
type Tag string

// TagAPI an interface containing all tag related methods
type TagAPI interface {
	GetTicketTags(ctx context.Context, ticketID int64) ([]Tag, error)
	GetOrganizationTags(ctx context.Context, organizationID int64) ([]Tag, error)
	GetUserTags(ctx context.Context, userID int64) ([]Tag, error)
	AddTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error)
	AddOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error)
	AddUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error)
	SetTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error)
	SetOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error)
	SetUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error)
	DeleteTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error)
	DeleteOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error)
	DeleteUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error)
}

// GetTicketTags get ticket tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetTicketTags(ctx context.Context, ticketID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/tickets/%d/tags.json", ticketID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

// GetOrganizationTags get organization tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetOrganizationTags(ctx context.Context, organizationID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/organizations/%d/tags.json", organizationID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

// GetUserTags get user tag list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#show-tags
func (z *Client) GetUserTags(ctx context.Context, userID int64) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/users/%d/tags.json", userID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, err
}

// AddTicketTags add tags to ticket
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/tickets/%d/tags", ticketID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// AddOrganizationTags add tags to organization
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/organizations/%d/tags", organizationID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// AddUserTags add tags to user
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tags#add-tags
func (z *Client) AddUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.put(ctx, fmt.Sprintf("/users/%d/tags", userID), data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// SetTicketTags replaces all tags on a ticket
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (z *Client) SetTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error) {
	return z.setTags(ctx, fmt.Sprintf("/tickets/%d/tags", ticketID), tags)
}

// SetOrganizationTags replaces all tags on an organization
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (z *Client) SetOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error) {
	return z.setTags(ctx, fmt.Sprintf("/organizations/%d/tags", organizationID), tags)
}

// SetUserTags replaces all tags on a user
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (z *Client) SetUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error) {
	return z.setTags(ctx, fmt.Sprintf("/users/%d/tags", userID), tags)
}

func (z *Client) setTags(ctx context.Context, path string, tags []Tag) ([]Tag, error) {
	var data, result struct {
		Tags []Tag `json:"tags"`
	}
	data.Tags = tags

	body, err := z.post(ctx, path, data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// DeleteTicketTags removes the given tags from a ticket and returns the remaining tags
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (z *Client) DeleteTicketTags(ctx context.Context, ticketID int64, tags []Tag) ([]Tag, error) {
	path := fmt.Sprintf("/tickets/%d/tags?tags=%s", ticketID, encodeTags(tags))
	return z.removeTags(ctx, path, nil)
}

// DeleteOrganizationTags removes the given tags from an organization and returns the remaining tags
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (z *Client) DeleteOrganizationTags(ctx context.Context, organizationID int64, tags []Tag) ([]Tag, error) {
	return z.removeTags(ctx, fmt.Sprintf("/organizations/%d/tags", organizationID), tagsBody(tags))
}

// DeleteUserTags removes the given tags from a user and returns the remaining tags
//
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (z *Client) DeleteUserTags(ctx context.Context, userID int64, tags []Tag) ([]Tag, error) {
	return z.removeTags(ctx, fmt.Sprintf("/users/%d/tags", userID), tagsBody(tags))
}

// removeTags deletes tags on the resource at path and returns the remaining tags.
// Tickets take the tags as a query parameter (body is nil); organizations and users
// take them in the request body.
func (z *Client) removeTags(ctx context.Context, path string, body interface{}) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	response, err := z.deleteWithBody(ctx, path, body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// tagsBody wraps tags in the request body shape expected by the tag endpoints
func tagsBody(tags []Tag) interface{} {
	return struct {
		Tags []Tag `json:"tags"`
	}{Tags: tags}
}

// encodeTags encodes tags as a comma-separated query parameter value
func encodeTags(tags []Tag) string {
	encoded := make([]string, len(tags))
	for i, tag := range tags {
		encoded[i] = url.QueryEscape(string(tag))
	}
	return strings.Join(encoded, ",")
}
