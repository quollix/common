package store

import (
	"fmt"
	"time"

	u "github.com/quollix/common/utils"
	"github.com/quollix/common/validation"
)

var (
	AppStoreOfficialMaintainerPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKlo+3ABMOUtj4JSvbYVnIDhCqrxkiMvleXudu9Lwc87"

	ApiPrefix    = "/api"
	WipeDataPath = ApiPrefix + "/wipe-data"

	userPath                = ApiPrefix + "/account"
	RegistrationPath        = userPath + "/registration"
	ConfirmRegistrationPath = userPath + "/validate"
	LoginPath               = userPath + "/login"
	LogoutPath              = userPath + "/logout"
	DeleteUserPath          = userPath + "/delete"
	ChangePasswordPath      = userPath + "/change-password"
	AccountDetailsPath      = userPath + "/details"

	EmailChangePath        = userPath + "/email-change"
	RequestEmailChangePath = EmailChangePath + "/request"
	ConfirmEmailChangePath = EmailChangePath + "/confirm"

	VersionPath       = ApiPrefix + "/versions"
	VersionUploadPath = VersionPath + "/upload"
	VersionDeletePath = VersionPath + "/delete"
	GetVersionsPath   = VersionPath + "/list"
	DownloadPath      = VersionPath + "/download"

	AppPath         = ApiPrefix + "/apps"
	AppCreationPath = AppPath + "/create"
	AppGetListPath  = AppPath + "/get-list"
	AppDeletePath   = AppPath + "/delete"
	SearchAppsPath  = AppPath + "/search"

	EmailConfigPath      = ApiPrefix + "/email-config"
	EmailConfigReadPath  = EmailConfigPath + "/read"
	EmailConfigWritePath = EmailConfigPath + "/write"
)

type AppStoreClient interface {
	Register(maintainer, password, email string, maintainerPublicKeyRaw []byte) error
	ConfirmationRegistration(code string) error
	Login(username, password string) error
	DeleteOwnAccount() error
	CreateApp(appName string) error
	SearchForApps(maintainerSearchTerm, appSearchTerm string, showUnofficialApps bool) ([]AppWithLatestVersion, error)
	ListOwnApps() ([]string, error)
	UploadVersion(appName, versionName string, creationTimestamp time.Time, content, signature []byte) error
	DownloadVersion(userName, appName, versionName string) (*Version, error)
	ListVersions(userName, appName string) ([]LeanVersionDto, error)
	DeleteVersion(appName, versionName string) error
	DeleteApp(appName string) error
	ChangePassword(oldPassword, newPassword string) error
	Logout() error
	GetAccountDetails() (*AccountDetails, error)

	GetEmailConfig() (*u.EmailConfig, error)
	SetEmailConfig(emailConfig *u.EmailConfig) error

	RequestEmailChange(email string) error
	ConfirmEmailChange(emailChangeConfirmationCode string) error
}

type AppStoreClientImpl struct {
	Parent    u.ComponentClient
	Validator validation.VersionValidator
}

func (h *AppStoreClientImpl) UploadVersion(appName, versionName string, creationTimestamp time.Time, content, signature []byte) error {
	versionUploads := VersionUploadDto{
		AppName:           appName,
		Version:           versionName,
		CreationTimestamp: creationTimestamp,
		Content:           content,
		Signature:         signature,
	}
	result, err := h.Parent.DoRequest(VersionUploadPath, versionUploads)
	if err != nil {
		return err
	}
	_ = result
	return nil
}

func (h *AppStoreClientImpl) Register(maintainer, password, email string, maintainerPublicKeyRaw []byte) error {
	form := RegistrationForm{
		User:                   maintainer,
		Password:               password,
		Email:                  email,
		MaintainerPublicKeyRaw: maintainerPublicKeyRaw,
	}
	_, err := h.Parent.DoRequest(RegistrationPath, form)
	return err
}

func (h *AppStoreClientImpl) ConfirmationRegistration(registrationCode string) error {
	_, err := h.Parent.DoRequest(ConfirmRegistrationPath, SecretString{Value: registrationCode})
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
		return fmt.Errorf("expected 1 cookie, got %d", len(cookies))
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

func (h *AppStoreClientImpl) DownloadVersion(userName, appName, versionName string) (*Version, error) {
	result, err := h.Parent.DoRequest(DownloadPath, VersionTree{
		Maintainer:  userName,
		AppName:     appName,
		VersionName: versionName,
	})
	if err != nil {
		return nil, err
	}

	version, err := u.UnpackResponse[Version](result)
	if err != nil {
		return nil, err
	}

	err = h.Validator.Validate(version.Content, version.Maintainer, version.AppName)
	if err != nil {
		return nil, fmt.Errorf("version validation failed: %w", err)
	}

	return version, nil
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

func (h *AppStoreClientImpl) DeleteVersion(appName, versionName string) error {
	_, err := h.Parent.DoRequest(VersionDeletePath, AppAndVersion{
		AppName:     appName,
		VersionName: versionName,
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

func (h *AppStoreClientImpl) RequestEmailChange(email string) error {
	_, err := h.Parent.DoRequest(RequestEmailChangePath, EmailString{Value: email})
	return err
}

func (h *AppStoreClientImpl) ConfirmEmailChange(emailChangeConfirmationCode string) error {
	_, err := h.Parent.DoRequest(ConfirmEmailChangePath, SecretString{Value: emailChangeConfirmationCode})
	return err
}

func (s *AppStoreClientImpl) WipeData() {
	_, err := s.Parent.DoRequest(WipeDataPath, nil)
	if err != nil {
		panic("failed to wipe data: " + err.Error())
	}
}
