package api_client

import (
	"encoding/json"
	"strconv"

	api "github.com/quollix/common/quollix/api"
)

type OidcRelyingPartiesClient struct {
	quollix *QuollixClient
}

func (c *OidcRelyingPartiesClient) Create(client *api.OidcRelyingPartyDto) error {
	request := api.OidcRelyingPartyRequest{
		Id:     "0",
		Name:   client.Name,
		Domain: client.Domain,
	}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcRelyingPartiesCreate, request)
	return err
}

func (c *OidcRelyingPartiesClient) Update(client *api.OidcRelyingPartyDto) error {
	request := api.OidcRelyingPartyRequest{
		Id:     strconv.Itoa(client.Id),
		Name:   client.Name,
		Domain: client.Domain,
	}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcRelyingPartiesUpdate, request)
	return err
}

func (c *OidcRelyingPartiesClient) Regenerate(clientId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcRelyingPartiesRegenerate, api.NumberString{Value: strconv.Itoa(clientId)})
	return err
}

func (c *OidcRelyingPartiesClient) List() ([]api.OidcRelyingPartyDto, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcRelyingPartiesList, nil)
	if err != nil {
		return nil, err
	}

	var clients []api.OidcRelyingPartyDto
	err = json.Unmarshal(body, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *OidcRelyingPartiesClient) Delete(clientId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcRelyingPartiesDelete, api.NumberString{Value: strconv.Itoa(clientId)})
	return err
}
