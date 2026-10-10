package api_client

import (
	"encoding/json"

	"github.com/quollix/common/quollix/api"
)

type LogsClient struct{ quollix *QuollixClient }

func (c *LogsClient) Read() ([]api.LogRecord, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendLogsRead, nil)
	if err != nil {
		return nil, err
	}
	var records []api.LogRecord
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (c *LogsClient) Flush() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendLogsFlush, nil)
	return err
}

func (c *LogsClient) Wipe() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendLogsWipe, nil)
	return err
}
