package test_environment

import (
	"github.com/quollix/common/quollix/api"
	"github.com/quollix/common/quollix/api_client"
	u "github.com/quollix/common/utils"
)

const (
	OidcProviderDomain  = "quollix.oidc-provider.localhost"
	OidcClientDomain    = "quollix.oidc-client.localhost"
	OidcProviderBaseUrl = "https://" + OidcProviderDomain
	OidcClientBaseUrl   = "https://" + OidcClientDomain
	OidcProviderHost    = "oidc-provider.localhost"
	OidcClientHost      = "oidc-client.localhost"

	OidcProviderAdminUsername = "provider-admin"
	OidcProviderAdminPassword = "localuserpassword"
	OidcProviderName          = "Provider-Quollix"

	defaultAdminUsername = "admin"
	defaultAdminPassword = "password"
	oidcClientName       = "Two-Instance-Client"
)

type OidcTwoInstanceClients struct {
	ProviderAdmin *api_client.QuollixClient
	ClientAdmin   *api_client.QuollixClient
}

func NewOidcTwoInstanceClients() *OidcTwoInstanceClients {
	return &OidcTwoInstanceClients{
		ProviderAdmin: api_client.NewQuollixClientForRootUrl(OidcProviderBaseUrl),
		ClientAdmin:   api_client.NewQuollixClientForRootUrl(OidcClientBaseUrl),
	}
}

func (c *OidcTwoInstanceClients) Reset() error {
	if err := c.ProviderAdmin.Test.ResetTestState(); err != nil {
		return u.Logger.NewError("reset OIDC provider server", "error", err.Error())
	}
	if err := c.ClientAdmin.Test.ResetTestState(); err != nil {
		return u.Logger.NewError("reset OIDC client server", "error", err.Error())
	}
	return nil
}

func (c *OidcTwoInstanceClients) Configure() error {
	if err := c.ProviderAdmin.Auth.SignIn(defaultAdminUsername, defaultAdminPassword); err != nil {
		return u.Logger.NewError("sign in to OIDC provider server", "error", err.Error())
	}
	if err := c.ClientAdmin.Auth.SignIn(defaultAdminUsername, defaultAdminPassword); err != nil {
		return u.Logger.NewError("sign in to OIDC client server", "error", err.Error())
	}
	if err := renameOidcProviderAdmin(c.ProviderAdmin); err != nil {
		return err
	}
	if err := c.ProviderAdmin.Settings.SetBaseDomainValue(OidcProviderHost); err != nil {
		return u.Logger.NewError("set OIDC provider base domain", "error", err.Error())
	}
	if err := c.ClientAdmin.Settings.SetBaseDomainValue(OidcClientHost); err != nil {
		return u.Logger.NewError("set OIDC client base domain", "error", err.Error())
	}

	relyingParty, err := createOidcRelyingParty(c.ProviderAdmin)
	if err != nil {
		return err
	}
	return createOidcAuthProvider(c.ClientAdmin, relyingParty)
}

func renameOidcProviderAdmin(providerAdmin *api_client.QuollixClient) error {
	admin, exists, err := providerAdmin.Users.GetByUsername(defaultAdminUsername)
	if err != nil {
		return u.Logger.NewError("find provider admin user", "error", err.Error())
	}
	if !exists {
		return u.Logger.NewError("find provider admin user: admin user does not exist")
	}
	if err := providerAdmin.Users.ChangeUsername(admin.Id, OidcProviderAdminUsername); err != nil {
		return u.Logger.NewError("rename provider admin user", "error", err.Error())
	}
	return nil
}

func createOidcRelyingParty(providerAdmin *api_client.QuollixClient) (*api.OidcRelyingPartyDto, error) {
	relyingParty := &api.OidcRelyingPartyDto{
		Name:   oidcClientName,
		Domain: OidcClientDomain,
	}
	if err := providerAdmin.OidcClients.Create(relyingParty); err != nil {
		return nil, u.Logger.NewError("create relying party in OIDC provider server", "error", err.Error())
	}
	clients, err := providerAdmin.OidcClients.List()
	if err != nil {
		return nil, u.Logger.NewError("list OIDC relying parties", "error", err.Error())
	}
	if len(clients) != 1 {
		return nil, u.Logger.NewError("expected exactly one OIDC relying party", "actual_count", len(clients))
	}
	return &clients[0], nil
}

func createOidcAuthProvider(clientAdmin *api_client.QuollixClient, relyingParty *api.OidcRelyingPartyDto) error {
	provider := &api.OidcAuthProviderDto{
		Name:             OidcProviderName,
		IssuerDomainPath: OidcProviderDomain,
		ClientId:         relyingParty.ClientId,
		ClientSecret:     relyingParty.ClientSecret,
	}
	if err := clientAdmin.OidcProviders.Create(provider); err != nil {
		return u.Logger.NewError("create OIDC provider in OIDC client server", "error", err.Error())
	}
	return nil
}
