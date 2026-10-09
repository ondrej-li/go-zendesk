package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// OrganizationField represents the Organization Custom field structure
type OrganizationField struct {
	ID                     int64               `json:"id,omitempty"`
	URL                    string              `json:"url,omitempty"`
	Title                  string              `json:"title"`
	Type                   string              `json:"type"`
	RelationshipTargetType string              `json:"relationship_target_type"`
	RelationshipFilter     RelationshipFilter  `json:"relationship_filter"`
	Active                 bool                `json:"active,omitempty"`
	CustomFieldOptions     []CustomFieldOption `json:"custom_field_options,omitempty"`
	Description            string              `json:"description,omitempty"`
	Key                    string              `json:"key"`
	Position               int64               `json:"position,omitempty"`
	RawDescription         string              `json:"raw_description,omitempty"`
	RawTitle               string              `json:"raw_title,omitempty"`
	RegexpForValidation    string              `json:"regexp_for_validation,omitempty"`
	System                 bool                `json:"system,omitempty"`
	Tag                    string              `json:"tag,omitempty"`
	CreatedAt              *time.Time          `json:"created_at,omitempty"`
	UpdatedAt              *time.Time          `json:"updated_at,omitempty"`
}

// OrganizationFieldAPI an interface containing all the organization field related zendesk methods
type OrganizationFieldAPI interface {
	GetOrganizationFields(ctx context.Context) ([]OrganizationField, Page, error)
	CreateOrganizationField(ctx context.Context, organizationField OrganizationField) (OrganizationField, error)
	GetOrganizationField(ctx context.Context, organizationFieldID int64) (OrganizationField, error)
	UpdateOrganizationField(ctx context.Context, organizationFieldID int64, organizationField OrganizationField) (OrganizationField, error)
	DeleteOrganizationField(ctx context.Context, organizationFieldID int64) error
	GetOrganizationFieldsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[OrganizationField]
	GetOrganizationFieldsOBP(ctx context.Context, opts *OBPOptions) ([]OrganizationField, Page, error)
	GetOrganizationFieldsCBP(ctx context.Context, opts *CBPOptions) ([]OrganizationField, CursorPaginationMeta, error)
}

// GetOrganizationFields fetches organization field list
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#list-organization-fields
func (z *Client) GetOrganizationFields(ctx context.Context) ([]OrganizationField, Page, error) {
	var data struct {
		OrganizationFields []OrganizationField `json:"organization_fields"`
		Page
	}

	body, err := z.get(ctx, "/organization_fields.json")
	if err != nil {
		return []OrganizationField{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []OrganizationField{}, Page{}, err
	}
	return data.OrganizationFields, data.Page, nil
}

// CreateOrganizationField creates new organization field
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#create-organization-field
func (z *Client) CreateOrganizationField(ctx context.Context, organizationField OrganizationField) (OrganizationField, error) {
	var data, result struct {
		OrganizationField OrganizationField `json:"organization_field"`
	}
	data.OrganizationField = organizationField

	body, err := z.post(ctx, "/organization_fields.json", data)
	if err != nil {
		return OrganizationField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OrganizationField{}, err
	}
	return result.OrganizationField, nil
}

// GetOrganizationField gets a specified organization field
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#show-organization-field
func (z *Client) GetOrganizationField(ctx context.Context, organizationFieldID int64) (OrganizationField, error) {
	var result struct {
		OrganizationField OrganizationField `json:"organization_field"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/organization_fields/%d.json", organizationFieldID))
	if err != nil {
		return OrganizationField{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return OrganizationField{}, err
	}
	return result.OrganizationField, nil
}

// UpdateOrganizationField updates a specified organization field
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#update-organization-field
func (z *Client) UpdateOrganizationField(ctx context.Context, organizationFieldID int64, organizationField OrganizationField) (OrganizationField, error) {
	var data, result struct {
		OrganizationField OrganizationField `json:"organization_field"`
	}
	data.OrganizationField = organizationField

	body, err := z.put(ctx, fmt.Sprintf("/organization_fields/%d.json", organizationFieldID), data)
	if err != nil {
		return OrganizationField{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return OrganizationField{}, err
	}
	return result.OrganizationField, nil
}

// DeleteOrganizationField deletes a specified organization field
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#delete-organization-field
func (z *Client) DeleteOrganizationField(ctx context.Context, organizationFieldID int64) error {
	return z.delete(ctx, fmt.Sprintf("/organization_fields/%d.json", organizationFieldID))
}
