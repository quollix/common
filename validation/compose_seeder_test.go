package validation

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	"gopkg.in/yaml.v3"
)

const (
	inputCompose = `
services:
  gitea:
    image: gitea/gitea:1.20.2
    container_name: samplemaintainer_gitea_gitea
    environment:
      WELL_KNOWN_ENDPOINT: https://${BASE_DOMAIN}/.well-known/openid-configuration
      CLIENT_ID: ${CLIENT_ID}
      CLIENT_SECRET: ${CLIENT_SECRET}
      IANA_TIMEZONE: ${IANA_TIMEZONE}
      OPTIONAL: ${SECRET_OPTIONAL}
      LOGS_RETENTION: ${LOGS_RETENTION_DAYS}d
      LOGS_STORAGE: ${LOGS_MAX_DISK_SPACE_GB}GB
      METRICS_RETENTION: ${METRICS_RETENTION_DAYS}d
    volumes:
      - samplemaintainer_gitea_data:/data

  giteadb:
    image: mariadb:10.5
    container_name: samplemaintainer_gitea_giteadb

volumes:
  samplemaintainer_gitea_data:
    name: samplemaintainer_gitea_data
`
	expectedCompletedCompose = `
services:
    gitea:
        image: gitea/gitea:1.20.2
        container_name: samplemaintainer_gitea_gitea
        environment:
            WELL_KNOWN_ENDPOINT: https://my-domain.com/.well-known/openid-configuration
            CLIENT_ID: abcdef
            CLIENT_SECRET: ghijkl
            IANA_TIMEZONE: Etc/UTC
            OPTIONAL: ${SECRET_OPTIONAL}
            LOGS_RETENTION: 30d
            LOGS_STORAGE: 5GB
            METRICS_RETENTION: 30d
        restart: unless-stopped
        networks:
            - samplemaintainer_gitea
        cap_drop:
            - ALL
        cap_add:
            - CAP_NET_BIND_SERVICE
            - CAP_CHOWN
            - CAP_FOWNER
            - CAP_SETGID
            - CAP_SETUID
            - CAP_DAC_OVERRIDE
        volumes:
            - samplemaintainer_gitea_data:/data

    giteadb:
        image: mariadb:10.5
        container_name: samplemaintainer_gitea_giteadb
        restart: unless-stopped
        networks:
            - samplemaintainer_gitea
        cap_drop:
            - ALL
        cap_add:
            - CAP_NET_BIND_SERVICE
            - CAP_CHOWN
            - CAP_FOWNER
            - CAP_SETGID
            - CAP_SETUID
            - CAP_DAC_OVERRIDE

networks:
    samplemaintainer_gitea:
        external: true

volumes:
    samplemaintainer_gitea_data:
        name: samplemaintainer_gitea_data
`
)

func TestCompleteDockerComposeYaml(t *testing.T) {
	dataMap := map[string]string{
		"IANA_TIMEZONE":          "Etc/UTC",
		"BASE_DOMAIN":            "my-domain.com",
		"CLIENT_ID":              "abcdef",
		"CLIENT_SECRET":          "ghijkl",
		"LOGS_RETENTION_DAYS":    "30",
		"LOGS_MAX_DISK_SPACE_GB": "5",
		"METRICS_RETENTION_DAYS": "30",
	}
	actualCompletedCompose, err := CompleteDockerComposeYaml("samplemaintainer", "gitea", []byte(inputCompose), dataMap)
	assert.Nil(t, err)
	assert.False(t, strings.ContainsRune(string(actualCompletedCompose), '\t'))
	assertYamlEquality(t, []byte(expectedCompletedCompose), actualCompletedCompose)
}

func assertYamlEquality(t *testing.T, a, b []byte) {
	var m1, m2 map[string]any
	if err := yaml.Unmarshal(a, &m1); err != nil {
		t.Fail()
	}
	if err := yaml.Unmarshal(b, &m2); err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(m1, m2) {
		fmt.Printf("EXPECTED YAML: \n\n%v\n", string(a))
		fmt.Printf("\n\nACTUAL YAML: \n\n%v\n", string(b))
		t.Fail()
	}
}
