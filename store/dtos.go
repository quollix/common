package store

import (
	"time"
)

type CreatedVersionResponse struct {
	AppId     int `json:"app_id"`
	VersionId int `json:"version_id"`
}

type NumberString struct {
	Value string `json:"value" validate:"number"`
}

type UserNameString struct {
	Value string `json:"value" validate:"number"`
}

type MaintainerNameString struct {
	Value string `json:"value" validate:"default"`
}

type AppNameString struct {
	Value string `json:"value" validate:"default"`
}

type SecretString struct {
	Value string `json:"value" validate:"secret"`
}

type EmailString struct {
	Value string `json:"value" validate:"email"`
}

type LeanVersionDto struct {
	VersionId             int       `json:"version_id"`
	Name                  string    `json:"name"`
	CreationTimestamp     time.Time `json:"creation_timestamp"`
	SizeInBytes           int64     `json:"size_in_bytes"`
	IsMigrationCheckpoint bool      `json:"is_migration_checkpoint"`
	DownloadCount         int64     `json:"download_count"`
}

type ChangePasswordForm struct {
	OldPassword string `validate:"password"`
	NewPassword string `validate:"password"`
}

type Version struct {
	VersionId                        int `json:"version_id"`
	Maintainer, AppName, VersionName string
	Content, Signature               []byte
	MaintainerPublicKeyRaw           []byte
	VersionCreationTimestamp         time.Time
	IsMigrationCheckpoint            bool
}

type SetupInitialPasswordForm struct {
	Token    string `validate:"secret"`
	Password string `validate:"password"`
}

type AdminMaintainerCreateForm struct {
	Name               string `validate:"default"`
	Email              string `validate:"email"`
	PublicKeyRaw       []byte
	PublicKeySignature []byte
}

type AdminMaintainerCreateResponse struct {
	SetupToken     string
	ExpirationDate time.Time
}

type MaintainerPublicKeyRecord struct {
	Maintainer         string
	PublicKeyRaw       []byte
	PublicKeySignature []byte
}

type AdminEmailRequest struct {
	Subject string `validate:"email_subject"`
	Body    string `validate:"email_body"`
}

type AdminEmailSendResult struct {
	SentCount   int
	FailedCount int
	Failures    []AdminEmailSendFailure
}

type AdminEmailSendFailure struct {
	Recipient string
	Error     string
}

type LoginCredentials struct {
	User     string `validate:"default"`
	Password string `validate:"password"`
}

type AppWithLatestVersion struct {
	Maintainer                     string
	AppName                        string
	LatestVersionId                int
	LatestVersionName              string
	LatestVersionCreationTimestamp time.Time
}

type SearchRequest struct {
	MaintainerSearchTerm string `validate:"search_term"`
	AppSearchTerm        string `validate:"search_term"`
	ShowUnofficialApps   bool
}

type AccountDetails struct {
	Name                 string
	Email                string
	PublicKeyRaw         []byte
	CookieExpirationDate time.Time
	UsedSpaceInBytes     int64
	StorageLimitInBytes  int64
	IsAdmin              bool
}

type VersionUploadDto struct {
	AppName           string `validate:"default"`
	Version           string `validate:"version_name"`
	CreationTimestamp time.Time
	Content           []byte
	Signature         []byte
}

type VersionTree struct {
	Maintainer  string `validate:"default"`
	AppName     string `validate:"default"`
	VersionName string `validate:"version_name"`
}

type VersionID struct {
	VersionId int `json:"version_id"`
}

type NextVersionForUpdateRequest struct {
	Maintainer                      string `validate:"default"`
	AppName                         string `validate:"default"`
	CurrentVersionCreationTimestamp time.Time
}

type NextVersionForUpdateResponse struct {
	UpdateAvailable bool
	Version         *Version
}

type VersionMigrationCheckpointRequest struct {
	VersionId             int `json:"version_id"`
	IsMigrationCheckpoint bool
}

type AppTree struct {
	Maintainer string `validate:"default"`
	AppName    string `validate:"default"`
}
