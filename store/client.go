package store

import (
	"time"

	u "github.com/quollix/common/utils"
)

var (
	AppStoreOfficialMaintainerPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKlo+3ABMOUtj4JSvbYVnIDhCqrxkiMvleXudu9Lwc87"

	ApiPrefix    = "/api"
	WipeDataPath = ApiPrefix + "/wipe-data"

	userPath                 = ApiPrefix + "/account"
	LoginPath                = userPath + "/login"
	LogoutPath               = userPath + "/logout"
	DeleteUserPath           = userPath + "/delete"
	ChangePasswordPath       = userPath + "/change-password"
	AccountDetailsPath       = userPath + "/details"
	SetupInitialPasswordPath = userPath + "/setup-password"
	EmailChangePath          = userPath + "/email-change"

	VersionPath                      = ApiPrefix + "/versions"
	VersionUploadPath                = VersionPath + "/upload"
	VersionDeletePath                = VersionPath + "/delete"
	GetVersionsPath                  = VersionPath + "/list"
	DownloadPath                     = VersionPath + "/download"
	DownloadByIDPath                 = VersionPath + "/download-by-id"
	DownloadNextVersionForUpdatePath = VersionPath + "/download-next-for-update"
	VersionMigrationCheckpointPath   = VersionPath + "/migration-checkpoint"

	AppPath         = ApiPrefix + "/apps"
	AppCreationPath = AppPath + "/create"
	AppGetListPath  = AppPath + "/get-list"
	AppDeletePath   = AppPath + "/delete"
	SearchAppsPath  = AppPath + "/search"

	EmailConfigPath      = ApiPrefix + "/email-config"
	EmailConfigReadPath  = EmailConfigPath + "/read"
	EmailConfigWritePath = EmailConfigPath + "/write"

	AdminEmailPath            = ApiPrefix + "/admin/email"
	AdminEmailTestPath        = AdminEmailPath + "/test"
	AdminEmailMaintainersPath = AdminEmailPath + "/maintainers"

	MaintainerPublicKeyPath   = ApiPrefix + "/maintainers/public-key"
	AdminMaintainerPath       = ApiPrefix + "/admin/maintainers"
	AdminMaintainerCreatePath = AdminMaintainerPath + "/create"
	AdminMaintainerDeletePath = AdminMaintainerPath + "/delete"
)

type AppStoreClient interface {
	Login(username, password string) error
	DeleteOwnAccount() error
	CreateApp(appName string) error
	SearchForApps(maintainerSearchTerm, appSearchTerm string, showUnofficialApps bool) ([]AppWithLatestVersion, error)
	ListOwnApps() ([]string, error)
	UploadVersionAndReturnCreatedVersion(appName, versionName string, creationTimestamp time.Time, content, signature []byte) (*CreatedVersionResponse, error)
	DownloadVersionByID(versionId int) (*Version, error)
	DownloadNextVersionForUpdate(userName, appName string, currentVersionCreationTimestamp time.Time) (*NextVersionForUpdateResponse, error)
	ListVersions(userName, appName string) ([]LeanVersionDto, error)
	DeleteVersionByID(versionId int) error
	SetVersionMigrationCheckpoint(versionId int, isMigrationCheckpoint bool) error
	DeleteApp(appName string) error
	ChangePassword(oldPassword, newPassword string) error
	Logout() error
	GetAccountDetails() (*AccountDetails, error)

	GetEmailConfig() (*u.EmailConfig, error)
	SetEmailConfig(emailConfig *u.EmailConfig) error

	ChangeEmail(email string) error

	CreateMaintainerByAdmin(name, email string, publicKeyRaw, publicKeySignature []byte) error
	DeleteMaintainerByAdmin(name string) error
	GetMaintainerPublicKeyRecord(maintainer string) (*MaintainerPublicKeyRecord, error)
	SetupInitialPassword(setupToken, password string) error

	SendTestEmailByAdmin(subject, body string) (*AdminEmailSendResult, error)
	SendEmailToMaintainersByAdmin(subject, body string) (*AdminEmailSendResult, error)
}

type AppStoreClientImpl struct {
	Parent u.ComponentClient
}

func (h *AppStoreClientImpl) UploadVersionAndReturnCreatedVersion(appName, versionName string, creationTimestamp time.Time, content, signature []byte) (*CreatedVersionResponse, error) {
	versionUploads := VersionUploadDto{
		AppName:           appName,
		Version:           versionName,
		CreationTimestamp: creationTimestamp,
		Content:           content,
		Signature:         signature,
	}
	result, err := h.Parent.DoRequest(VersionUploadPath, versionUploads)
	if err != nil {
		return nil, err
	}
	return u.UnpackResponse[CreatedVersionResponse](result)
}

func (h *AppStoreClientImpl) SetupInitialPassword(setupToken, password string) error {
	_, err := h.Parent.DoRequest(SetupInitialPasswordPath, SetupInitialPasswordForm{
		Token:    setupToken,
		Password: password,
	})
	return err
}

func (h *AppStoreClientImpl) Login(username, password string) error {
	creds := LoginCredentials{
		User:     username,
		Password: password,
	}

	resp, err := h.Parent.DoRequestWithFullResponse(LoginPath, creds)
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return u.Logger.NewError("unexpected cookie count", "expected", 1, "actual", len(cookies))
	}
	h.Parent.Cookie = cookies[0]
	return nil
}

func (h *AppStoreClientImpl) DeleteOwnAccount() error {
	_, err := h.Parent.DoRequest(DeleteUserPath, nil)
	return err
}

func (h *AppStoreClientImpl) CreateApp(appName string) error {
	_, err := h.Parent.DoRequest(AppCreationPath, AppNameString{Value: appName})
	return err
}

func (h *AppStoreClientImpl) SearchForApps(maintainerSearchTerm, appSearchTerm string, showUnofficialApps bool) ([]AppWithLatestVersion, error) {
	appSearchRequest := SearchRequest{
		MaintainerSearchTerm: maintainerSearchTerm,
		AppSearchTerm:        appSearchTerm,
		ShowUnofficialApps:   showUnofficialApps,
	}
	result, err := h.Parent.DoRequest(SearchAppsPath, appSearchRequest)
	if err != nil {
		return nil, err
	}

	apps, err := u.UnpackResponse[[]AppWithLatestVersion](result)
	if err != nil {
		return nil, err
	}

	return *apps, nil
}

func (h *AppStoreClientImpl) ListOwnApps() ([]string, error) {
	result, err := h.Parent.DoRequest(AppGetListPath, nil)
	if err != nil {
		return nil, err
	}

	apps, err := u.UnpackResponse[[]string](result)
	if err != nil {
		return nil, err
	}

	return *apps, nil
}

func (h *AppStoreClientImpl) DownloadVersionByID(versionId int) (*Version, error) {
	result, err := h.Parent.DoRequest(DownloadByIDPath, VersionID{VersionId: versionId})
	if err != nil {
		return nil, err
	}

	version, err := u.UnpackResponse[Version](result)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (h *AppStoreClientImpl) DownloadNextVersionForUpdate(userName, appName string, currentVersionCreationTimestamp time.Time) (*NextVersionForUpdateResponse, error) {
	result, err := h.Parent.DoRequest(DownloadNextVersionForUpdatePath, NextVersionForUpdateRequest{
		Maintainer:                      userName,
		AppName:                         appName,
		CurrentVersionCreationTimestamp: currentVersionCreationTimestamp.UTC(),
	})
	if err != nil {
		return nil, err
	}

	response, err := u.UnpackResponse[NextVersionForUpdateResponse](result)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *AppStoreClientImpl) ListVersions(userName, appName string) ([]LeanVersionDto, error) {
	result, err := h.Parent.DoRequest(GetVersionsPath, AppTree{
		Maintainer: userName,
		AppName:    appName,
	})
	if err != nil {
		return nil, err
	}

	versions, err := u.UnpackResponse[[]LeanVersionDto](result)
	if err != nil {
		return nil, err
	}

	return *versions, nil
}

func (h *AppStoreClientImpl) DeleteVersionByID(versionId int) error {
	_, err := h.Parent.DoRequest(VersionDeletePath, VersionID{VersionId: versionId})
	return err
}

func (h *AppStoreClientImpl) SetVersionMigrationCheckpoint(versionId int, isMigrationCheckpoint bool) error {
	_, err := h.Parent.DoRequest(VersionMigrationCheckpointPath, VersionMigrationCheckpointRequest{
		VersionId:             versionId,
		IsMigrationCheckpoint: isMigrationCheckpoint,
	})
	return err
}

func (h *AppStoreClientImpl) DeleteApp(appName string) error {
	_, err := h.Parent.DoRequest(AppDeletePath, AppNameString{appName})
	return err
}

func (h *AppStoreClientImpl) ChangePassword(oldPassword, newPassword string) error {
	form := ChangePasswordForm{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	_, err := h.Parent.DoRequest(ChangePasswordPath, form)
	return err
}

func (h *AppStoreClientImpl) Logout() error {
	_, err := h.Parent.DoRequest(LogoutPath, nil)
	return err
}

func (h *AppStoreClientImpl) GetAccountDetails() (*AccountDetails, error) {
	result, err := h.Parent.DoRequest(AccountDetailsPath, nil)
	if err != nil {
		return nil, err
	}
	userDetails, err := u.UnpackResponse[AccountDetails](result)
	if err != nil {
		return nil, err
	}
	return userDetails, nil
}

func (h *AppStoreClientImpl) GetEmailConfig() (*u.EmailConfig, error) {
	body, err := h.Parent.DoRequest(EmailConfigReadPath, nil)
	if err != nil {
		return nil, err
	}
	emailConfig, err := u.UnpackResponse[u.EmailConfig](body)
	if err != nil {
		return nil, err
	}
	return emailConfig, nil
}

func (h *AppStoreClientImpl) SetEmailConfig(emailConfig *u.EmailConfig) error {
	_, err := h.Parent.DoRequest(EmailConfigWritePath, emailConfig)
	return err
}

func (h *AppStoreClientImpl) ChangeEmail(email string) error {
	_, err := h.Parent.DoRequest(EmailChangePath, EmailString{Value: email})
	return err
}

func (h *AppStoreClientImpl) CreateMaintainerByAdmin(name, email string, publicKeyRaw, publicKeySignature []byte) error {
	_, err := h.Parent.DoRequest(AdminMaintainerCreatePath, AdminMaintainerCreateForm{
		Name:               name,
		Email:              email,
		PublicKeyRaw:       publicKeyRaw,
		PublicKeySignature: publicKeySignature,
	})
	return err
}

func (h *AppStoreClientImpl) GetMaintainerPublicKeyRecord(maintainer string) (*MaintainerPublicKeyRecord, error) {
	responseBody, err := h.Parent.DoRequest(MaintainerPublicKeyPath, MaintainerNameString{Value: maintainer})
	if err != nil {
		return nil, err
	}
	return u.UnpackResponse[MaintainerPublicKeyRecord](responseBody)
}

func (h *AppStoreClientImpl) DeleteMaintainerByAdmin(name string) error {
	_, err := h.Parent.DoRequest(AdminMaintainerDeletePath, MaintainerNameString{Value: name})
	return err
}

func (h *AppStoreClientImpl) SendTestEmailByAdmin(subject, body string) (*AdminEmailSendResult, error) {
	responseBody, err := h.Parent.DoRequest(AdminEmailTestPath, AdminEmailRequest{Subject: subject, Body: body})
	if err != nil {
		return nil, err
	}
	return u.UnpackResponse[AdminEmailSendResult](responseBody)
}

func (h *AppStoreClientImpl) SendEmailToMaintainersByAdmin(subject, body string) (*AdminEmailSendResult, error) {
	responseBody, err := h.Parent.DoRequest(AdminEmailMaintainersPath, AdminEmailRequest{Subject: subject, Body: body})
	if err != nil {
		return nil, err
	}
	return u.UnpackResponse[AdminEmailSendResult](responseBody)
}

func (s *AppStoreClientImpl) WipeData() {
	_, err := s.Parent.DoRequest(WipeDataPath, nil)
	if err != nil {
		panic("failed to wipe data: " + err.Error())
	}
}
