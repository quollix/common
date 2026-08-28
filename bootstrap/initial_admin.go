package bootstrap

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"

	u "github.com/quollix/common/utils"
	"github.com/quollix/common/validation"
)

const (
	InitialAdminPasswordEnvVar        = "INITIAL_ADMIN_PASSWORD"
	InitialAdminNameEnvVar            = "INITIAL_ADMIN_NAME"
	GeneratedInitialAdminPasswordSize = 20
)

const (
	generatedInitialAdminPasswordLogMessage = "INITIAL_ADMIN_PASSWORD environment variable is not set, generated random initial admin password"
	initialAdminUsernameLogField            = "username"
	initialAdminPasswordLogField            = "password"
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
		u.Logger.Info(generatedInitialAdminPasswordLogMessage, initialAdminUsernameLogField, adminName, initialAdminPasswordLogField, adminPassword)
	} else {
		u.Logger.Info("Creating initial admin user from environment configuration", "username", adminName)
	}
	return adminPassword, nil
}

func WaitForGeneratedInitialAdminCredentials(containerName string) (string, string, error) {
	var username string
	var password string

	err := u.EventuallyWithTimeout(30*time.Second, 500*time.Millisecond, func() error {
		logs, err := readContainerLogs(containerName)
		if err != nil {
			return err
		}

		username, password, err = ExtractGeneratedInitialAdminCredentials(logs)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

func readContainerLogs(containerName string) (string, error) {
	output, err := exec.Command("docker", "logs", containerName).CombinedOutput() // #nosec G204: fixed docker binary with structured args; containerName is not shell-expanded
	if err != nil {
		return "", u.Logger.NewError("docker logs failed", "error", err.Error(), "container_name", containerName)
	}
	return string(output), nil
}

func ExtractGeneratedInitialAdminCredentials(logs string) (string, string, error) {
	scanner := bufio.NewScanner(strings.NewReader(logs))
	for scanner.Scan() {
		if username, password, found := extractGeneratedInitialAdminCredentialsFromLine(scanner.Text()); found {
			return username, password, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", u.Logger.NewError(err.Error())
	}
	return "", "", u.Logger.NewError("password log line was not found")
}

func extractGeneratedInitialAdminCredentialsFromLine(line string) (string, string, bool) {
	var fields map[string]any
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return "", "", false
	}

	if jsonStringField(fields, "msg") != generatedInitialAdminPasswordLogMessage {
		return "", "", false
	}

	username := jsonStringField(fields, initialAdminUsernameLogField)
	password := jsonStringField(fields, initialAdminPasswordLogField)
	if username == "" || password == "" {
		return "", "", false
	}
	return username, password, true
}

func jsonStringField(fields map[string]any, key string) string {
	value, ok := fields[key].(string)
	if !ok {
		return ""
	}
	return value
}
