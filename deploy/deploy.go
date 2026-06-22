package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	u "github.com/quollix/common/utils"
	"github.com/quollix/taskrunner"
)

func DeployLocal(tr *taskrunner.TaskRunner, appName string, localServiceYAML string) {
	tempDir, err := os.MkdirTemp("", appName+"-deploy-local-")
	if err != nil {
		tr.Log.Error("failed to create temp dir: %v", err)
		tr.ExitWithError()
	}
	defer u.RemoveDir(tempDir)

	if err := os.WriteFile(filepath.Join(tempDir, "docker-compose.test.yml"), []byte(renderTestCompose(appName, localServiceYAML)), 0o644); err != nil { // #nosec G306: temp compose file for local-only deployment needs normal read permissions for docker tooling
		tr.Log.Error("failed to write test docker compose: %v", err)
		tr.ExitWithError()
	}

	tr.Cmd().Dir(tempDir).Run("docker compose -f docker-compose.test.yml up -d")
	tr.WaitForWebPageToBeReady("http://localhost:8080" + HealthPath)
}

func StartLocalPostgres(tr *taskrunner.TaskRunner, appName string) func() {
	tempDir, err := os.MkdirTemp("", appName+"-postgres-")
	if err != nil {
		tr.Log.Error("failed to create temp dir for postgres compose: %v", err)
		tr.ExitWithError()
	}

	composeFile := filepath.Join(tempDir, "docker-compose.yml")
	content := renderStandalonePostgresCompose(appName)
	if err := os.WriteFile(composeFile, []byte(content), 0o644); err != nil { // #nosec G306: temp compose file for local-only deployment needs normal read permissions for docker tooling
		u.RemoveDir(tempDir)
		tr.Log.Error("failed to write temp postgres compose: %v", err)
		tr.ExitWithError()
	}

	tr.Cmd().Dir(tempDir).Run("docker compose -f %s up -d postgres", composeFile)

	return func() {
		u.RemoveDir(tempDir)
	}
}

func CleanupLocal(tr *taskrunner.TaskRunner, appName string) {
	tr.Cmd().AllowFail().Run("docker rm -f %s %s", appContainerName(appName), postgresContainerName(appName))
	tr.Cmd().AllowFail().Run("docker volume rm -f %s", postgresVolumeName(appName))
	tr.Cmd().AllowFail().Run("docker network rm %s", networkName(appName))
}

func renderTestCompose(appName, appServiceYAML string) string {
	return applyPlaceholders(testComposeTemplate, appName, "", appServiceYAML)
}

func renderStandalonePostgresCompose(appName string) string {
	return applyPlaceholders(standalonePostgresComposeTemplate, appName, "", "")
}

func renderLocalPostgresService(appName string) string {
	return fmt.Sprintf(`postgres:
  image: postgres:17.2
  container_name: %s
  environment:
    - POSTGRES_USER=postgres
    - POSTGRES_HOST_AUTH_METHOD=trust
    - POSTGRES_DB=postgres
  ports:
    - "127.0.0.1:5432:5432"
  tmpfs:
    - /var/lib/postgresql/data
  networks:
    - quollix_%s
	`, postgresContainerName(appName), appName)
}

func applyPlaceholders(template string, appName string, releaseTag string, appServiceYAML string) string {
	renderedAppServiceYAML := strings.NewReplacer(
		"${RELEASE_TAG}", releaseTag,
	).Replace(appServiceYAML)

	return strings.NewReplacer(
		"${APP_NAME}", appName,
		"${RELEASE_TAG}", releaseTag,
		"${APP_SERVICE_BLOCK}", indentServiceBlock(renderedAppServiceYAML),
		"${LOCAL_POSTGRES_SERVICE_BLOCK}", indentServiceBlock(renderLocalPostgresService(appName)),
	).Replace(template)
}

func indentServiceBlock(block string) string {
	trimmed := strings.TrimSpace(block)
	if trimmed == "" {
		return ""
	}

	var builder strings.Builder
	for _, line := range strings.Split(trimmed, "\n") {
		if strings.TrimSpace(line) == "" {
			builder.WriteString("\n")
			continue
		}
		builder.WriteString("  ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	return builder.String()
}

func appContainerName(appName string) string {
	return fmt.Sprintf("%s_%s_%s", u.OfficialMaintainer, appName, appName)
}

func postgresContainerName(appName string) string {
	return fmt.Sprintf("%s_%s_%s", u.OfficialMaintainer, appName, u.OfficialDatabaseAppName)
}

func postgresVolumeName(appName string) string {
	return fmt.Sprintf("%s_%s_data", u.OfficialMaintainer, appName)
}

func networkName(appName string) string {
	return fmt.Sprintf("%s_%s", u.OfficialMaintainer, appName)
}
