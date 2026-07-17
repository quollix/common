package api_client

import (
	"encoding/json"

	api "github.com/quollix/common/quollix/api"
)

type AuthClient struct {
	quollix *QuollixClient
}

func (c *AuthClient) SignIn(username, password string) error {
	signInCredentials := api.Credentials{Username: username, Password: password}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendSignIn, signInCredentials)
	return err
}

func (c *AuthClient) SignOut() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersSignOut, nil)
	return err
}

func (c *AuthClient) GetCurrentUser() (*api.User, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendCheckAuth, nil)
	if err != nil {
		return nil, err
	}
	var authResponse api.User
	err = json.Unmarshal(responseBody, &authResponse)
	if err != nil {
		return nil, err
	}
	return &authResponse, nil
}

func (c *AuthClient) ChangePassword(oldPassword, newPassword string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersChangeOwnPassword, api.ChangeOwnPasswordRequest{
		CurrentPassword: oldPassword,
		NewPassword:     newPassword,
	})
	return err
}
