package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
)

type QuollixMaintenanceClient struct {
	quollix *QuollixClient
}

func (c *QuollixMaintenanceClient) SaveConfigs(request *api.MaintenanceConfigDto) error {
	_, err := c.quollix.Parent.DoRequestWithFullResponse(api.Paths.BackendMaintenanceConfigsSave, request)
	return err
}

func (c *QuollixMaintenanceClient) ReadConfigs() (*api.MaintenanceConfig, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendMaintenanceConfigsRead, nil)
	if err != nil {
		return nil, err
	}
	var config *api.MaintenanceConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *QuollixMaintenanceClient) ReadRetentionPolicy() (*api.RetentionPolicy, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendMaintenanceRetentionPolicyRead, nil)
	if err != nil {
		return nil, err
	}

	var policy api.RetentionPolicy
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (c *QuollixMaintenanceClient) SaveRetentionPolicy(policy *api.RetentionPolicy) error {
	_, err := c.quollix.Parent.DoRequestWithFullResponse(api.Paths.BackendMaintenanceRetentionPolicySave, policy)
	return err
}

func (c *QuollixMaintenanceClient) ExecuteJob() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendMaintenanceTriggerMaintenanceJob, nil)
	return err
}
