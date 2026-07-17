package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
	"strconv"

	u "github.com/quollix/common/utils"
)

type GroupsClient struct {
	quollix *QuollixClient
}

func (c *GroupsClient) ListAppsAccessByGroup(groupId int) (*api.AppsAccessByGroup, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsListAppsAccessByGroup, api.NumberString{Value: strconv.Itoa(groupId)})
	if err != nil {
		return nil, err
	}
	var out api.AppsAccessByGroup
	return &out, json.Unmarshal(body, &out)
}

func (c *GroupsClient) CreateGroup(name string) (*api.Group, error) {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsCreate, api.DefaultString{Value: name})
	if err != nil {
		return nil, err
	}

	allGroups, err := c.ListAllGroups()
	if err != nil {
		return nil, err
	}
	for _, group := range allGroups {
		if group.Name == name {
			return &group, nil
		}
	}
	return nil, u.Logger.NewError("group not found")
}

func (c *GroupsClient) DeleteGroup(groupId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsDelete, api.NumberString{Value: strconv.Itoa(groupId)})
	return err
}

func (c *GroupsClient) ListAllGroups() ([]api.Group, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendListAllGroups, nil)
	if err != nil {
		return nil, err
	}
	var out []api.Group
	return out, json.Unmarshal(body, &out)
}

func (c *GroupsClient) AddUserToGroup(userId, groupId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsAddUsers, api.GroupIdAndUserIds{
		GroupId: strconv.Itoa(groupId),
		UserIds: []string{strconv.Itoa(userId)},
	})
	return err
}

func (c *GroupsClient) RemoveUserFromGroup(userId, groupId int) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsRemoveUsers, api.GroupIdAndUserIds{
		GroupId: strconv.Itoa(groupId),
		UserIds: []string{strconv.Itoa(userId)},
	})
	return err
}

func (c *GroupsClient) ListUsersByGroupMembership(groupId int) (*api.UsersByGroupMembership, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsListUsersByMembership, api.NumberString{Value: strconv.Itoa(groupId)})
	if err != nil {
		return nil, err
	}
	var out api.UsersByGroupMembership
	return &out, json.Unmarshal(body, &out)
}

func (c *GroupsClient) GrantAppAccess(groupId int, appName string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsGrantGroupAccessToApps, api.GroupIdAndAppNames{
		GroupId:  strconv.Itoa(groupId),
		AppNames: []string{appName},
	})
	return err
}

func (c *GroupsClient) RevokeAppAccess(groupId int, appName string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendGroupsRevokeGroupAccessToApps, api.GroupIdAndAppNames{
		GroupId:  strconv.Itoa(groupId),
		AppNames: []string{appName},
	})
	return err
}
