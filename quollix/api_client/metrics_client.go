package api_client

import (
	"encoding/json"

	"github.com/quollix/common/quollix/api"
)

type MetricsClient struct{ quollix *QuollixClient }

func (c *MetricsClient) Read() ([]api.MetricSeries, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendMetricsRead, nil)
	if err != nil {
		return nil, err
	}
	var records []api.MetricSeries
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (c *MetricsClient) Collect() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendMetricsCollect, nil)
	return err
}

func (c *MetricsClient) Wipe() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendMetricsWipe, nil)
	return err
}
