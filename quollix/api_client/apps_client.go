package api_client

import (
	"encoding/json"

	api "github.com/quollix/common/quollix/api"

	"github.com/quollix/common/store"
)

type QuollixAppsClient struct {
	quollix *QuollixClient
}

func (c *QuollixAppsClient) InstallFromStore(maintainer, appName, version string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendStoreVersionsInstall, store.VersionTree{
		Maintainer:  maintainer,
		AppName:     appName,
		VersionName: version,
	})
	return err
}

func (c *QuollixAppsClient) SearchStore(maintainerName, appName string, searchForUnofficialApps bool) ([]store.AppWithLatestVersion, error) {
	appSearchRequest := store.SearchRequest{
		MaintainerSearchTerm: maintainerName,
		AppSearchTerm:        appName,
		ShowUnofficialApps:   searchForUnofficialApps,
	}
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendStoreSearch, appSearchRequest)
	if err != nil {
		return nil, err
	}
	var apps []store.AppWithLatestVersion
	err = json.Unmarshal(responseBody, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *QuollixAppsClient) ListVersions(userName, appName string) ([]store.LeanVersionDto, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendStoreVersionsList, store.AppTree{
		Maintainer: userName,
		AppName:    appName,
	})
	if err != nil {
		return nil, err
	}
	var versions []store.LeanVersionDto
	if err = json.Unmarshal(responseBody, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (c *QuollixAppsClient) ListInstalled() ([]api.AppDto, error) {
	responseBody, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsList, nil)
	if err != nil {
		return nil, err
	}
	var apps []api.AppDto
	err = json.Unmarshal(responseBody, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *QuollixAppsClient) Delete(appId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsDelete, api.NumberString{Value: appId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixAppsClient) Update(appId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsUpdate, api.NumberString{Value: appId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixAppsClient) RegenerateSecret(appId, name string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppSecretRegenerate, api.AppSecretRegenerationRequest{
		AppId: appId,
		Name:  name,
	})
	return err
}

func (c *QuollixAppsClient) Start(appId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsStart, api.NumberString{Value: appId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixAppsClient) Stop(appId string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsStop, api.NumberString{Value: appId})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixAppsClient) SetAccessPolicy(appId, accessPolicy string) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppsChangeAccessPolicy, api.ChangeAccessPolicyRequest{
		AppId:        appId,
		AccessPolicy: accessPolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *QuollixAppsClient) UpdateMaintenanceSettings(appId string, autoUpdatesEnabled, autoBackupsEnabled bool) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppAutomaticMaintenanceSettings, api.AutoMaintenanceSettingsResponse{
		AppId:                   appId,
		AutomaticUpdatesEnabled: autoUpdatesEnabled,
		AutomaticBackupsEnabled: autoBackupsEnabled,
	})
	return err
}

func (c *QuollixAppsClient) UploadVersionFile(file api.BinaryFile) error {
	return c.quollix.uploadBinaryFile(api.Paths.BackendAppUploadToApplication, file)
}

func (c *QuollixAppsClient) DownloadVersionFile(appId string) (*api.BinaryFile, error) {
	return c.quollix.downloadBinaryFile(api.Paths.BackendAppDownloadFromApplication, api.NumberString{Value: appId})
}

func (c *QuollixAppsClient) GetCurrentOperations() ([]string, bool, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendAppOperationInfo, nil)
	if err != nil {
		return nil, false, err
	}
	var out api.AppOperationInfoResponse
	if err = json.Unmarshal(body, &out); err != nil {
		return nil, false, err
	}
	return out.Operations, out.IsOngoing, nil
}

func (c *QuollixAppsClient) DownloadVersion(maintainer, appName, versionName string) (*api.BinaryFile, error) {
	return c.quollix.downloadBinaryFile(api.Paths.BackendStoreVersionsDownload, store.VersionTree{
		Maintainer:  maintainer,
		AppName:     appName,
		VersionName: versionName,
	})
}
