package bootstrap

import (
	"os"

	u "github.com/quollix/common/utils"
	"github.com/quollix/common/validation"
)

const (
	InitialAdminPasswordEnvVar        = "INITIAL_ADMIN_PASSWORD"
	InitialAdminNameEnvVar            = "INITIAL_ADMIN_NAME"
	GeneratedInitialAdminPasswordSize = 20
)

func GetInitialAdminCredentials(defaultUsername string) (string, string, error) {
	username, err := getInitialAdminName(defaultUsername)
	if err != nil {
		return "", "", err
	}

	password, err := getInitialAdminPassword(username)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func getInitialAdminName(defaultUsername string) (string, error) {
	adminName := os.Getenv(InitialAdminNameEnvVar)
	if adminName == "" {
		adminName = defaultUsername
	}

	if err := validation.Validate("adminName", validation.FieldUsername, adminName); err != nil {
		return "", u.Logger.NewError("env variable is not valid", "env_variable", InitialAdminNameEnvVar)
	}
	return adminName, nil
}

func getInitialAdminPassword(adminName string) (string, error) {
	adminPassword := os.Getenv(InitialAdminPasswordEnvVar)
	isGeneratedPassword := adminPassword == ""
	if isGeneratedPassword {
		var err error
		adminPassword, err = (&u.AuthHelperImpl{}).GenerateSecret()
		if err != nil {
			return "", err
		}
		if len(adminPassword) > GeneratedInitialAdminPasswordSize {
			adminPassword = adminPassword[:GeneratedInitialAdminPasswordSize]
		}
	}

	if err := validation.Validate("adminPassword", validation.FieldPassword, adminPassword); err != nil {
		return "", u.Logger.NewError("env variable is not valid", "env_variable", InitialAdminPasswordEnvVar)
	}

	if isGeneratedPassword {
		u.Logger.Info("INITIAL_ADMIN_PASSWORD environment variable is not set, generated random initial admin password", "username", adminName, "password", adminPassword)
	} else {
		u.Logger.Info("Creating initial admin user from environment configuration", "username", adminName)
	}
	return adminPassword, nil
}
