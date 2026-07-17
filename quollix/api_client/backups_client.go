package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
)

type QuollixBackupsClient struct {
	quollix *QuollixClient
}

func (c *QuollixBackupsClient) Create(appId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsCreate, api.NumberString{Value: appId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixBackupsClient) ListByApp(maintainer, appName string) ([]api.BackupInfo, error) {
	backupListRequest := api.MaintainerAndApp{
		Maintainer: maintainer,
		AppName:    appName,
	}
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsList, backupListRequest)
	if err != nil {
		return nil, err
	}
	var backups []api.BackupInfo
	if responseBody == nil {
		return backups, nil
	}
	err = json.Unmarshal(responseBody, &backups)
	if err != nil {
		return nil, err
	}
	return backups, nil
}

func (c *QuollixBackupsClient) Delete(backupIds []string) error {
	deleteBackupRequest := api.BackupsOperationRequest{BackupIds: backupIds}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsDelete, deleteBackupRequest)
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixBackupsClient) Restore(backupId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsRestore, api.BackupOperationRequest{BackupId: backupId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixBackupsClient) ListAppsInRepository() ([]api.MaintainerAndApp, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsListApps, nil)
	if err != nil {
		return nil, err
	}
	var apps []api.MaintainerAndApp
	err = json.Unmarshal(responseBody, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *QuollixBackupsClient) PurgeServer(repo *api.BackupServerConfigs) error {
	purgeRequest := &api.SshConnectionRequest{
		Host:          repo.Host,
		SshPort:       repo.SshPort,
		SshKnownHosts: repo.SshKnownHosts,
		SshUser:       repo.SshUser,
		SshPassword:   repo.SshPassword,
	}
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendBackupsPurgeBackupServer, purgeRequest)
	return err
}
