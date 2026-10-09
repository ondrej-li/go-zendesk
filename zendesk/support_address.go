package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SupportAddress is an email support address (recipient address)
// https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/
type SupportAddress struct {
	ID                       int64      `json:"id,omitempty"`
	BrandID                  int64      `json:"brand_id,omitempty"`
	Name                     string     `json:"name,omitempty"`
	Email                    string     `json:"email"`
	Default                  bool       `json:"default,omitempty"`
	ForwardingStatus         string     `json:"forwarding_status,omitempty"`
	SPFStatus                string     `json:"spf_status,omitempty"`
	CNAMEStatus              string     `json:"cname_status,omitempty"`
	DNSResults               string     `json:"dns_results,omitempty"`
	DomainVerificationCode   string     `json:"domain_verification_code,omitempty"`
	DomainVerificationStatus string     `json:"domain_verification_status,omitempty"`
	CreatedAt                *time.Time `json:"created_at,omitempty"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`
}

// SupportAddressAPI an interface containing all support address related methods
type SupportAddressAPI interface {
	GetSupportAddresses(ctx context.Context) ([]SupportAddress, error)
	GetSupportAddress(ctx context.Context, supportAddressID int64) (SupportAddress, error)
	CreateSupportAddress(ctx context.Context, address SupportAddress) (SupportAddress, error)
	UpdateSupportAddress(ctx context.Context, supportAddressID int64, address SupportAddress) (SupportAddress, error)
	DeleteSupportAddress(ctx context.Context, supportAddressID int64) error
	VerifySupportAddressForwarding(ctx context.Context, supportAddressID int64, verificationType string) error
}

// GetSupportAddresses fetches the support address list
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#list-support-addresses
func (z *Client) GetSupportAddresses(ctx context.Context) ([]SupportAddress, error) {
	var result struct {
		RecipientAddresses []SupportAddress `json:"recipient_addresses"`
	}

	body, err := z.get(ctx, "/recipient_addresses.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.RecipientAddresses, nil
}

// GetSupportAddress gets a specified support address
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#show-support-address
func (z *Client) GetSupportAddress(ctx context.Context, supportAddressID int64) (SupportAddress, error) {
	var result struct {
		RecipientAddress SupportAddress `json:"recipient_address"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/recipient_addresses/%d.json", supportAddressID))
	if err != nil {
		return SupportAddress{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SupportAddress{}, err
	}
	return result.RecipientAddress, nil
}

// CreateSupportAddress adds a Zendesk or external support address to the account
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#create-support-address
func (z *Client) CreateSupportAddress(ctx context.Context, address SupportAddress) (SupportAddress, error) {
	var data, result struct {
		RecipientAddress SupportAddress `json:"recipient_address"`
	}
	data.RecipientAddress = address

	body, err := z.post(ctx, "/recipient_addresses.json", data)
	if err != nil {
		return SupportAddress{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SupportAddress{}, err
	}
	return result.RecipientAddress, nil
}

// UpdateSupportAddress updates a specified support address. The email address of an
// existing address cannot be changed.
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#update-support-address
func (z *Client) UpdateSupportAddress(ctx context.Context, supportAddressID int64, address SupportAddress) (SupportAddress, error) {
	var data, result struct {
		RecipientAddress SupportAddress `json:"recipient_address"`
	}
	data.RecipientAddress = address

	body, err := z.put(ctx, fmt.Sprintf("/recipient_addresses/%d.json", supportAddressID), data)
	if err != nil {
		return SupportAddress{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return SupportAddress{}, err
	}
	return result.RecipientAddress, nil
}

// DeleteSupportAddress deletes a specified support address
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#delete-support-address
func (z *Client) DeleteSupportAddress(ctx context.Context, supportAddressID int64) error {
	return z.delete(ctx, fmt.Sprintf("/recipient_addresses/%d.json", supportAddressID))
}

// VerifySupportAddressForwarding requests a verification check for a support address.
// An empty verificationType defaults to "forwarding"; "spf" and "dns" are also accepted.
// The endpoint returns no results: re-fetch the address to read its status.
// ref: https://developer.zendesk.com/api-reference/ticketing/channels/support_addresses/#verify-support-address-forwarding
func (z *Client) VerifySupportAddressForwarding(ctx context.Context, supportAddressID int64, verificationType string) error {
	if verificationType == "" {
		verificationType = "forwarding"
	}

	data := struct {
		Type string `json:"type"`
	}{Type: verificationType}

	_, err := z.put(ctx, fmt.Sprintf("/recipient_addresses/%d/verify.json", supportAddressID), data)
	return err
}
