package api_client

import (
	"encoding/json"
	"strconv"

	api "github.com/quollix/common/quollix/api"
)

type UsersClient struct {
	quollix *QuollixClient
}

func (c *UsersClient) List() ([]api.User, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersList, nil)
	if err != nil {
		return nil, err
	}
	var users []api.User
	err = json.Unmarshal(responseBody, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *UsersClient) GetByUsername(username string) (api.User, bool, error) {
	users, err := c.List()
	if err != nil {
		return api.User{}, false, err
	}
	for _, user := range users {
		if user.Username == username {
			return user, true, nil
		}
	}
	return api.User{}, false, nil
}

func (c *UsersClient) Delete(userId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersDelete, api.NumberString{Value: strconv.Itoa(userId)})
	return err
}

func (c *UsersClient) Invite(username, email string) error {
	request := api.InviteUserRequest{Username: username, Email: email}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersInviteUser, request)
	return err
}

func (c *UsersClient) SetPasswordViaToken(password, token string) error {
	request := api.AcceptNewPasswordViaTokenRequest{Password: password, Token: token}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersSetPassword, request)
	return err
}

func (c *UsersClient) SetOwnPassword(password string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersSetOwnPassword, api.SetOwnPasswordRequest{
		NewPassword: password,
	})
	return err
}

func (c *UsersClient) ResetPassword(userId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersResetPassword, api.NumberString{Value: strconv.Itoa(userId)})
	return err
}

func (c *UsersClient) ChangeUsername(userId int, newUsername string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersChangeUsername, api.ChangeUsernameRequest{
		UserId:   strconv.Itoa(userId),
		Username: newUsername,
	})
	return err
}

func (c *UsersClient) ChangeEmail(userId int, newEmail string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersChangeEmail, api.ChangeEmailRequest{
		UserId:   strconv.Itoa(userId),
		NewEmail: newEmail,
	})
	return err
}

func (c *UsersClient) SetEnabled(userId int, isEnabled bool) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendUsersSetEnabled, api.SetUserEnabledRequest{
		UserId:    strconv.Itoa(userId),
		IsEnabled: isEnabled,
	})
	return err
}
