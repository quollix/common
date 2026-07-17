package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
	"strconv"
)

type OidcAuthProvidersClient struct {
	quollix *QuollixClient
}

func (c *OidcAuthProvidersClient) Create(provider *api.OidcAuthProviderDto) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcAuthProvidersCreate, provider)
	return err
}

func (c *OidcAuthProvidersClient) Update(provider *api.OidcAuthProviderDto) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcAuthProvidersUpdate, provider)
	return err
}

func (c *OidcAuthProvidersClient) List() ([]api.OidcAuthProviderDto, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcAuthProvidersList, nil)
	if err != nil {
		return nil, err
	}

	var providers []api.OidcAuthProviderDto
	err = json.Unmarshal(body, &providers)
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (c *OidcAuthProvidersClient) Delete(providerId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcAuthProvidersDelete, api.NumberString{Value: strconv.Itoa(providerId)})
	return err
}

func (c *OidcAuthProvidersClient) TestDiscovery(issuerDomainPath string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcAuthProvidersTestDiscovery, api.OidcAuthProviderDiscoveryRequest{IssuerDomainPath: issuerDomainPath})
	return err
}

func (c *OidcAuthProvidersClient) StartLogin(providerId int) (api.OidcStartLoginResponse, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendOidcSignIn, api.NumberString{Value: strconv.Itoa(providerId)})
	if err != nil {
		return api.OidcStartLoginResponse{}, err
	}

	var response api.OidcStartLoginResponse
	err = json.Unmarshal(body, &response)
	return response, err
}
