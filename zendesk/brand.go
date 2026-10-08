package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Brand is struct for brand payload
// https://developer.zendesk.com/rest_api/docs/support/brands
type Brand struct {
	ID                int64      `json:"id,omitempty"`
	URL               string     `json:"url,omitempty"`
	Name              string     `json:"name"`
	BrandURL          string     `json:"brand_url,omitempty"`
	HasHelpCenter     bool       `json:"has_help_center,omitempty"`
	HelpCenterState   string     `json:"help_center_state,omitempty"`
	Active            bool       `json:"active,omitempty"`
	Default           bool       `json:"default,omitempty"`
	Logo              Attachment `json:"logo,omitempty"`
	TicketFormIDs     []int64    `json:"ticket_form_ids,omitempty"`
	Subdomain         string     `json:"subdomain"`
	HostMapping       string     `json:"host_mapping,omitempty"`
	SignatureTemplate string     `json:"signature_template"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"`
}

// BrandListOptions is options for GetBrands
//
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#list-brands
type BrandListOptions struct {
	PageOptions
}

// BrandAPI an interface containing all methods associated with zendesk brands
type BrandAPI interface {
	GetBrands(ctx context.Context, opts *BrandListOptions) ([]Brand, Page, error)
	GetBrandsOBP(ctx context.Context, opts *OBPOptions) ([]Brand, Page, error)
	GetBrandsCBP(ctx context.Context, opts *CBPOptions) ([]Brand, CursorPaginationMeta, error)
	GetBrandsIterator(ctx context.Context, opts *PaginationOptions) *Iterator[Brand]
	CreateBrand(ctx context.Context, brand Brand) (Brand, error)
	GetBrand(ctx context.Context, brandID int64) (Brand, error)
	UpdateBrand(ctx context.Context, brandID int64, brand Brand) (Brand, error)
	DeleteBrand(ctx context.Context, brandID int64) error
}

// GetBrands fetches brand list
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/#list-brands
func (z *Client) GetBrands(ctx context.Context, opts *BrandListOptions) ([]Brand, Page, error) {
	var data struct {
		Brands []Brand `json:"brands"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &BrandListOptions{}
	}

	u, err := addOptions("/brands.json", tmp)
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
	return data.Brands, data.Page, nil
}

// CreateBrand creates new brand
// https://developer.zendesk.com/rest_api/docs/support/brands#create-brand
func (z *Client) CreateBrand(ctx context.Context, brand Brand) (Brand, error) {
	var data, result struct {
		Brand Brand `json:"brand"`
	}
	data.Brand = brand

	body, err := z.post(ctx, "/brands.json", data)
	if err != nil {
		return Brand{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Brand{}, err
	}
	return result.Brand, nil
}

// GetBrand gets a specified brand
// ref: https://developer.zendesk.com/rest_api/docs/support/brands#show-brand
func (z *Client) GetBrand(ctx context.Context, brandID int64) (Brand, error) {
	var result struct {
		Brand Brand `json:"brand"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/brands/%d.json", brandID))

	if err != nil {
		return Brand{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Brand{}, err
	}

	return result.Brand, err
}

// UpdateBrand updates a brand with the specified brand
// ref: https://developer.zendesk.com/rest_api/docs/support/brands#update-brand
func (z *Client) UpdateBrand(ctx context.Context, brandID int64, brand Brand) (Brand, error) {
	var result, data struct {
		Brand Brand `json:"brand"`
	}

	data.Brand = brand

	body, err := z.put(ctx, fmt.Sprintf("/brands/%d.json", brandID), data)

	if err != nil {
		return Brand{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Brand{}, err
	}

	return result.Brand, err
}

// DeleteBrand deletes the specified brand
// ref: https://developer.zendesk.com/rest_api/docs/support/brands#delete-brand
func (z *Client) DeleteBrand(ctx context.Context, brandID int64) error {
	err := z.delete(ctx, fmt.Sprintf("/brands/%d.json", brandID))

	if err != nil {
		return err
	}

	return nil
}
