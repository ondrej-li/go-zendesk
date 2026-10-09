package zendesk

import "time"

// Count is a cached record count together with the time it was last refreshed.
// Count endpoints that return a nested count object
// ("count": {"value": ..., "refreshed_at": ...}) decode into this type.
type Count struct {
	RefreshedAt *time.Time `json:"refreshed_at,omitempty"`
	Value       int64      `json:"value"`
}

// Page is base struct for resource pagination
type Page struct {
	PreviousPage *string `json:"previous_page"`
	NextPage     *string `json:"next_page"`
	Count        int64   `json:"count"`
}

// PageOptions is options for list method of paginatable resources.
// It's used to create query string.
//
// ref: https://developer.zendesk.com/rest_api/docs/support/introduction#pagination
type PageOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}

// HasPrev checks if the Page has previous page
func (p Page) HasPrev() bool {
	return (p.PreviousPage != nil)
}

// HasNext checks if the Page has next page
func (p Page) HasNext() bool {
	return (p.NextPage != nil)
}
