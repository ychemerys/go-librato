package librato

import (
	"fmt"
	"net/http"
)

// MetricsService handles communication with the Librato API methods related to
// metrics.
type MetricsService struct {
	client *Client
}

// Metric represents a Librato Metric.
type Metric struct {
	Name        *string           `json:"name"`
	Period      *uint             `json:"period,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
	Attributes  *MetricAttributes `json:"attributes,omitempty"`
}

type MetricAttributes struct {
	Color *string `json:"color"`
	// These are interface{} because sometimes the Librato API
	// returns strings, and sometimes it returns integers
	DisplayMax        interface{} `json:"display_max"`
	DisplayMin        interface{} `json:"display_min"`
	DisplayUnitsShort string      `json:"display_units_short"`
	DisplayStacked    bool        `json:"display_stacked"`
	DisplayTransform  string      `json:"display_transform"`
}

type ListMetricsOptions struct {
	*PaginationMeta
	Name string `url:"name,omitempty"`
}

// Advance to the specified page in result set, while retaining
// the filtering options.
func (l *ListMetricsOptions) AdvancePage(next *PaginationMeta) ListMetricsOptions {
	return ListMetricsOptions{
		PaginationMeta: next,
		Name:           l.Name,
	}
}

type ListMetricsResponse struct {
	ThisPage *PaginationResponseMeta
	NextPage *PaginationMeta
}

// List metrics using the provided options.
//
// Librato API docs: https://www.librato.com/docs/api/#retrieve-metrics
func (m *MetricsService) List(opts *ListMetricsOptions) ([]Metric, *ListMetricsResponse, error) {
	u, err := urlWithOptions("metrics", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var metricsResponse struct {
		Query   PaginationResponseMeta
		Metrics []Metric
	}

	_, err = m.client.Do(req, &metricsResponse)
	if err != nil {
		return nil, nil, err
	}

	return metricsResponse.Metrics,
		&ListMetricsResponse{
			ThisPage: &metricsResponse.Query,
			NextPage: metricsResponse.Query.nextPage(opts.PaginationMeta),
		},
		nil
}

// Update a metric.
//
// Librato API docs: https://www.librato.com/docs/api/#update-metric-by-name
func (m *MetricsService) Update(metric *Metric) (*http.Response, error) {
	u := fmt.Sprintf("metrics/%s", *metric.Name)

	req, err := m.client.NewRequest("PUT", u, metric)
	if err != nil {
		return nil, err
	}

	return m.client.Do(req, nil)
}
