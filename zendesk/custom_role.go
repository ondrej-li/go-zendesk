package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Configuration is a dictionary of custom configuration fields
type Configuration map[string]interface{}

// CustomRole is zendesk CustomRole JSON payload format
// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
type CustomRole struct {
	Description     string        `json:"description,omitempty"`
	ID              int64         `json:"id,omitempty"`
	TeamMemberCount int64         `json:"team_member_count"`
	Name            string        `json:"name"`
	Configuration   Configuration `json:"configuration"`
	RoleType        int64         `json:"role_type"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
}

// CustomRoleAPI an interface containing all CustomRole related methods
type CustomRoleAPI interface {
	GetCustomRoles(ctx context.Context) ([]CustomRole, error)
	GetCustomRole(ctx context.Context, customRoleID int64) (CustomRole, error)
	CreateCustomRole(ctx context.Context, customRole CustomRole) (CustomRole, error)
	UpdateCustomRole(ctx context.Context, customRoleID int64, customRole CustomRole) (CustomRole, error)
	DeleteCustomRole(ctx context.Context, customRoleID int64) error
}

// GetRoles fetch CustomRoles list
func (z *Client) GetCustomRoles(ctx context.Context) ([]CustomRole, error) {
	var data struct {
		CustomRoles []CustomRole `json:"custom_roles"`
		Page
	}

	u := "/custom_roles.json"

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.CustomRoles, nil
}

// GetCustomRole gets a specified custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#show-custom-role
func (z *Client) GetCustomRole(ctx context.Context, customRoleID int64) (CustomRole, error) {
	var result struct {
		CustomRole CustomRole `json:"custom_role"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/custom_roles/%d.json", customRoleID))
	if err != nil {
		return CustomRole{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return CustomRole{}, err
	}
	return result.CustomRole, nil
}

// CreateCustomRole creates a custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#create-custom-role
func (z *Client) CreateCustomRole(ctx context.Context, customRole CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	data.CustomRole = customRole

	body, err := z.post(ctx, "/custom_roles.json", data)
	if err != nil {
		return CustomRole{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return CustomRole{}, err
	}
	return result.CustomRole, nil
}

// UpdateCustomRole updates a specified custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#update-custom-role
func (z *Client) UpdateCustomRole(ctx context.Context, customRoleID int64, customRole CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	data.CustomRole = customRole

	body, err := z.put(ctx, fmt.Sprintf("/custom_roles/%d.json", customRoleID), data)
	if err != nil {
		return CustomRole{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return CustomRole{}, err
	}
	return result.CustomRole, nil
}

// DeleteCustomRole deletes a specified custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#delete-custom-role
func (z *Client) DeleteCustomRole(ctx context.Context, customRoleID int64) error {
	return z.delete(ctx, fmt.Sprintf("/custom_roles/%d.json", customRoleID))
}
