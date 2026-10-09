package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// UserField is struct for user_field payload
type UserField struct {
	ID                  int64               `json:"id,omitempty"`
	URL                 string              `json:"url,omitempty"`
	Key                 string              `json:"key,omitempty"`
	Type                string              `json:"type"`
	Title               string              `json:"title"`
	RawTitle            string              `json:"raw_title,omitempty"`
	Description         string              `json:"description,omitempty"`
	RawDescription      string              `json:"raw_description,omitempty"`
	Position            int64               `json:"position,omitempty"`
	Active              bool                `json:"active,omitempty"`
	System              bool                `json:"system,omitempty"`
	RegexpForValidation string              `json:"regexp_for_validation,omitempty"`
	Tag                 string              `json:"tag,omitempty"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	CreatedAt           time.Time           `json:"created_at,omitempty"`
	UpdatedAt           time.Time           `json:"updated_at,omitempty"`
}

type UserFieldListOptions struct {
	PageOptions
}

type UserFieldAPI interface {
	GetUserFields(ctx context.Context, opts *UserFieldListOptions) ([]UserField, Page, error)
	CreateUserField(ctx context.Context, userField UserField) (UserField, error)
	GetUserField(ctx context.Context, userFieldID int64) (UserField, error)
	UpdateUserField(ctx context.Context, userFieldID int64, userField UserField) (UserField, error)
	DeleteUserField(ctx context.Context, userFieldID int64) error
	GetUserFieldsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[UserField]
	GetUserFieldsOBP(ctx context.Context, opts *OBPOptions) ([]UserField, Page, error)
	GetUserFieldsCBP(ctx context.Context, opts *CBPOptions) ([]UserField, CursorPaginationMeta, error)
}

// GetUserFields fetch trigger list
//
// https://developer.zendesk.com/rest_api/docs/support/user_fields#list-user-fields
func (z *Client) GetUserFields(ctx context.Context, opts *UserFieldListOptions) ([]UserField, Page, error) {
	var data struct {
		UserFields []UserField `json:"user_fields"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &UserFieldListOptions{}
	}

	u, err := addOptions("/user_fields.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.UserFields, data.Page, nil
}

// CreateUserField creates new user field
// ref: https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#create-user-field
func (z *Client) CreateUserField(ctx context.Context, userField UserField) (UserField, error) {
	var data, result struct {
		UserField UserField `json:"user_field"`
	}
	data.UserField = userField

	body, err := z.post(ctx, "/user_fields.json", data)
	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return UserField{}, err
	}
	return result.UserField, nil
}

// GetUserField gets a specified user field
// ref: https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#show-user-field
func (z *Client) GetUserField(ctx context.Context, userFieldID int64) (UserField, error) {
	var result struct {
		UserField UserField `json:"user_field"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/user_fields/%d.json", userFieldID))
	if err != nil {
		return UserField{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return UserField{}, err
	}
	return result.UserField, nil
}

// UpdateUserField updates a specified user field
// ref: https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#update-user-field
func (z *Client) UpdateUserField(ctx context.Context, userFieldID int64, userField UserField) (UserField, error) {
	var data, result struct {
		UserField UserField `json:"user_field"`
	}
	data.UserField = userField

	body, err := z.put(ctx, fmt.Sprintf("/user_fields/%d.json", userFieldID), data)
	if err != nil {
		return UserField{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return UserField{}, err
	}
	return result.UserField, nil
}

// DeleteUserField deletes a specified user field
// ref: https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#delete-user-field
func (z *Client) DeleteUserField(ctx context.Context, userFieldID int64) error {
	return z.delete(ctx, fmt.Sprintf("/user_fields/%d.json", userFieldID))
}
