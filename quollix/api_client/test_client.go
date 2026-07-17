package api_client

import (
	api "github.com/quollix/common/quollix/api"
)

type QuollixTestClient struct {
	quollix *QuollixClient
}

func (c *QuollixTestClient) ResetTestState() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendResetTestState, nil)
	return err
}
