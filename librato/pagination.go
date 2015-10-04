package librato

import (
	"fmt"
	"net/url"
)

// Pagination metadata from Librato API responses
type PaginationResponseMeta struct {
	Offset uint `json:"offset"`
	Length uint `json:"length"`
	Total  uint `json:"total"`
	Found  uint `json:"found"`
}

// Calculate the pagination metadata for the next page of the result set.
// Takes the metadata used to request the current page so that it can use the
// same sort/orderby options
func (p *PaginationResponseMeta) nextPage(originalQuery *PaginationMeta) (next *PaginationMeta) {
	nextOffset := p.Offset + p.Length

	if nextOffset >= p.Found {
		return nil
	}

	next = &PaginationMeta{}
	next.Offset = nextOffset
	next.Length = p.Length

	if originalQuery != nil {
		next.OrderBy = originalQuery.OrderBy
		next.Sort = originalQuery.Sort
	}

	return next
}

// Metadata that the librato API requires for pagination
// http://dev.librato.com/v1/pagination
type PaginationMeta struct {
	Offset  uint   `url:"offset,omitempty"`
	Length  uint   `url:"length,omitempty"`
	OrderBy string `url:"orderby,omitempty"`
	Sort    string `url:"sort,omitempty"`
}

// Custom Encoder for the query string encoder library.
// The encoder allows other structs to embed PaginationMeta, and have it
// appear in the top-level query string fields without nesting.
func (m *PaginationMeta) EncodeValues(name string, values *url.Values) error {
	if m == nil {
		return nil
	}

	values.Set("offset", fmt.Sprintf("%d", m.Offset))
	values.Set("length", fmt.Sprintf("%d", m.Length))

	if m.OrderBy != "" {
		values.Set("orderby", m.OrderBy)
	}
	if m.Sort != "" {
		values.Set("sort", m.Sort)
	}

	return nil
}
