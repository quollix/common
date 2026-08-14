package api

import (
	"time"

	u "github.com/quollix/common/utils"
)

type Credentials struct {
	Username string `json:"username" validate:"username"`
	Password string `json:"password" validate:"password"`
}

type InviteUserRequest struct {
	Username string `json:"username" validate:"username"`
	Email    string `json:"email" validate:"email_or_empty"`
}

type InviteUserViaEmailRequest struct {
	Username string `json:"username" validate:"username"`
	Email    string `json:"email" validate:"email"`
}

type AcceptNewPasswordViaTokenRequest struct {
	Password string `validate:"password"`
	Token    string `validate:"secret"`
}

type ChangeUsernameRequest struct {
	UserId   string `json:"user_id" validate:"number"`
	Username string `json:"username" validate:"username"`
}

type ChangeEmailRequest struct {
	UserId   string `json:"user_id" validate:"number"`
	NewEmail string `json:"new_email" validate:"email"`
}

type SetUserEnabledRequest struct {
	UserId    string `json:"user_id" validate:"number"`
	IsEnabled bool   `json:"is_enabled"`
}

type AdminAppDto struct {
	AppId, Maintainer, AppName, VersionName, AccessPolicy,
	Port, ClientId, ClientSecret, AppSecret, DocsUrl, VersionCreationTimestampFormatted, VersionCreationTimestampTooltip string
	IsRunning, IsOfficialDatabaseApp, AutomaticBackupsEnabled, AutomaticUpdatesEnabled, IsOfficial bool
	VersionCreationTimestamp                                                                       time.Time
	VersionContent                                                                                 []byte
	Secrets                                                                                        map[string]string
}

type NonAdminAppDto struct {
	Maintainer, AppName string
	IsPublic            bool `json:"is_public"`
}

type AppAccessSecretRequest struct {
	AppName string `json:"app_name" validate:"default"`
}

type AppSecretRequest struct {
	AppId string `json:"app_id" validate:"number"`
	Name  string `json:"name" validate:"compose_secret_name"`
}

type AppSecretUpdateRequest struct {
	AppId string `json:"app_id" validate:"number"`
	Name  string `json:"name" validate:"compose_secret_name"`
	Value string `json:"value" validate:"credential"`
}

type ChangeAccessPolicyRequest struct {
	AppId        string `json:"app_id" validate:"number"`
	AccessPolicy string `json:"access_policy" validate:"access_policy"`
}

type AppOperationInfoResponse struct {
	Operations            []string `json:"operations"`
	IsOngoing             bool     `json:"is_ongoing"`
	AppOperationsFinished []string `json:"app_operations_finished"`
}

type AutoMaintenanceSettingsResponse struct {
	AppId                   string `json:"app_id" validate:"number"`
	AutomaticUpdatesEnabled bool   `json:"automatic_updates_enabled"`
	AutomaticBackupsEnabled bool   `json:"automatic_backups_enabled"`
}

type KnownHostsRequest struct {
	Host string `json:"host" validate:"remote_host"`
	Port string `json:"port" validate:"number"`
}

type Group struct {
	Id   int
	Name string `validate:"default"`
}

type Member struct {
	Id   int
	Name string `validate:"default"`
}

type UsersByGroupMembership struct {
	In    []Member
	NotIn []Member
}

type AppsAccessByGroup struct {
	Granted    []string
	NotGranted []string
}

type GroupIdAndUserIds struct {
	GroupId string   `json:"group_id" validate:"number"`
	UserIds []string `json:"user_ids" validate:"number"`
}

type GroupIdAndAppNames struct {
	GroupId  string   `json:"group_id" validate:"number"`
	AppNames []string `json:"app_names" validate:"default"`
}

type MaintenanceConfigDto struct {
	IanaTimezone               string `json:"iana_timezone" validate:"ignore"`
	MaintenanceWindowStartHour int    `json:"maintenance_window_start_hour"`
}

type MaintenanceConfig struct {
	IanaTimezone               string
	MaintenanceWindowStartHour int
	NextMaintenanceAt          time.Time
}

type RetentionPolicy struct {
	KeepPreUpdate int `json:"keep_pre_update"`
	KeepDaily     int `json:"keep_daily"`
	KeepWeekly    int `json:"keep_weekly"`
	KeepMonthly   int `json:"keep_monthly"`
	KeepYearly    int `json:"keep_yearly"`
}

type Dns01ChallengeInfo struct {
	RecordName      string `json:"record_name"`
	WildcardKeyAuth string `json:"wildcard_key_auth"`
}

type OidcAuthProviderDto struct {
	Id               int    `json:"id"`
	Name             string `json:"name" validate:"loose"`
	IssuerDomainPath string `json:"issuer_domain_path" validate:"domain_path"`
	ClientId         string `json:"client_id" validate:"credential"`
	ClientSecret     string `json:"client_secret" validate:"credential"`
}

type OidcAuthProviderDiscoveryRequest struct {
	IssuerDomainPath string `json:"issuer_domain_path" validate:"domain_path"`
}

type OidcStartLoginResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type OidcRelyingPartyDto struct {
	Id           int    `json:"id"`
	Name         string `json:"name" validate:"loose"`
	Domain       string `json:"domain" validate:"domain"`
	ClientId     string `json:"client_id" validate:"loose"`
	ClientSecret string `json:"client_secret" validate:"loose"`
}

type OidcRelyingPartyRequest struct {
	Id     string `json:"id" validate:"number"`
	Name   string `json:"name" validate:"loose"`
	Domain string `json:"domain" validate:"domain"`
}

type TestEmailRequest struct {
	EmailConfig u.EmailConfig `json:"email_config"`
	ToEmail     string        `json:"to_email" validate:"email"`
}

type InvitationEmailTemplateRequest struct {
	Template string `json:"template" validate:"ignore"`
}
