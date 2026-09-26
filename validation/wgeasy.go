package validation

import (
	"reflect"
	"regexp"

	u "github.com/quollix/common/utils"
	"gopkg.in/yaml.v3"
)

const (
	wgEasyAppName          = "wgeasy"
	wgEasyImagePlaceholder = "WG_EASY_IMAGE"
	wgEasyComposeYaml      = `
services:
  wgeasy:
    container_name: quollix_wgeasy_wgeasy
    image: WG_EASY_IMAGE
    environment:
      INSECURE: true
    volumes:
      - quollix_wgeasy_wireguard:/etc/wireguard
      - /lib/modules:/lib/modules:ro
    cap_add:
      - NET_ADMIN
      - NET_RAW
      - SYS_MODULE
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
      - net.ipv6.conf.all.disable_ipv6=0
      - net.ipv6.conf.all.forwarding=1
    ports:
      - "51820:51820/udp"
    labels:
      quollix.port: 51821
    networks:
      - quollix_wgeasy
    restart: unless-stopped
volumes:
  quollix_wgeasy_wireguard:
    name: quollix_wgeasy_wireguard
networks:
  quollix_wgeasy:
    external: true
`
)

var wgEasyImageRegex = regexp.MustCompile(`^ghcr\.io/wg-easy/wg-easy:([A-Za-z0-9_][A-Za-z0-9_.-]{0,127})(@sha256:[a-fA-F0-9]{64})?$`)

// IsWgEasyApp reports whether the identity refers to the official wg-easy app.
func IsWgEasyApp(maintainer, appName string) bool {
	return maintainer == u.OfficialMaintainer && appName == wgEasyAppName
}

func validateWgEasyComposeYaml(content []byte) error {
	actual := map[string]any{}
	if err := yaml.Unmarshal(content, &actual); err != nil {
		return u.Logger.NewError("invalid wg-easy app definition")
	}

	services, ok := actual["services"].(map[string]any)
	if !ok {
		return u.Logger.NewError("invalid wg-easy image")
	}
	wgEasyService, ok := services[wgEasyAppName].(map[string]any)
	if !ok {
		return u.Logger.NewError("invalid wg-easy image")
	}
	image, ok := wgEasyService["image"].(string)
	if !ok || !isAllowedWgEasyImage(image) {
		return u.Logger.NewError("invalid wg-easy image")
	}
	wgEasyService["image"] = wgEasyImagePlaceholder

	expected := map[string]any{}
	if err := yaml.Unmarshal([]byte(wgEasyComposeYaml), &expected); err != nil {
		return u.Logger.NewError(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		return u.Logger.NewError("wg-easy app definition does not match the expected configuration")
	}
	return nil
}

func isAllowedWgEasyImage(image string) bool {
	matches := wgEasyImageRegex.FindStringSubmatch(image)
	return matches != nil && matches[1] != "latest"
}
