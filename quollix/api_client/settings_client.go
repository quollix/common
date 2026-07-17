package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
)

type QuollixSettingsClient struct {
	quollix *QuollixClient
}

func (c *QuollixSettingsClient) SetBaseDomainValue(baseDomain string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsBaseDomainSave, api.NumberString{Value: baseDomain})
	return err
}

func (c *QuollixSettingsClient) GetBaseDomainValue() (string, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsBaseDomainRead, nil)
	if err != nil {
		return "", err
	}
	var baseDomain string
	err = json.Unmarshal(responseBody, &baseDomain)
	if err != nil {
		return "", err
	}
	return baseDomain, nil
}

func (c *QuollixSettingsClient) ReadSshConfigs() (*api.BackupServerConfigs, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsSshRead, nil)
	if err != nil {
		return nil, err
	}
	var backupServerConfig *api.BackupServerConfigs
	err = json.Unmarshal(responseBody, &backupServerConfig)
	if err != nil {
		return nil, err
	}
	return backupServerConfig, nil
}

func (c *QuollixSettingsClient) SaveSshConfigs(repo *api.BackupServerConfigs) error {
	err := c.TestSshAccess(repo)
	if err != nil {
		return err
	}
	_, err = c.quollix.Parent.DoRequest(api.Paths.BackendSettingsSshSave, repo)
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixSettingsClient) TestSshAccess(repo *api.BackupServerConfigs) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsSshTestAccess, repo.ConvertToSshConnectionTestRequest())
	return err
}

func (c *QuollixSettingsClient) GetKnownHosts(repo *api.BackupServerConfigs) (string, error) {
	knownHostsRequest := api.KnownHostsRequest{
		Host: repo.Host,
		Port: repo.SshPort,
	}
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsGetSshKnownHosts, knownHostsRequest)
	if err != nil {
		return "", err
	}
	var knownHostsStruct api.SingleString
	err = json.Unmarshal(responseBody, &knownHostsStruct)
	if err != nil {
		return "", err
	}
	return knownHostsStruct.Value, nil
}
