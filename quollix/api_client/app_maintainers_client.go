package api_client

import (
	"encoding/json"

	api "github.com/quollix/common/quollix/api"
)

type AppMaintainersClient struct {
	quollix *QuollixClient
}

func (c *AppMaintainersClient) List() ([]api.MaintainerPublicKeyDto, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppMaintainersList, nil)
	if err != nil {
		return nil, err
	}
	var out []api.MaintainerPublicKeyDto
	return out, json.Unmarshal(body, &out)
}

func (c *AppMaintainersClient) Add(name, publicKey string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppMaintainersAdd, api.MaintainerPublicKeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	})
	return err
}

func (c *AppMaintainersClient) Delete(name string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppMaintainersDelete, api.MaintainerPublicKeyDeleteRequest{Name: name})
	return err
}
