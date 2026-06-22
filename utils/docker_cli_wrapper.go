package utils

import (
	"os/exec"
	"strings"
)

type DockerCliWrapper interface {
	IsContainerRunning(containerName string) (bool, error)
	ImageExists(image string) (bool, error)
	DeleteImage(image string) error
}

type DockerCliWrapperImpl struct{}

func (d *DockerCliWrapperImpl) IsContainerRunning(containerName string) (bool, error) {
	output, err := exec.Command( // #nosec G204: fixed docker binary with structured args; containerName is passed as an argument, not shell-expanded
		"docker",
		"ps",
		"--filter",
		"name=^/"+containerName+"$",
		"--format",
		"{{.Names}}",
	).Output()
	if err != nil {
		return false, Logger.NewError(err.Error(), "container_name", containerName)
	}
	return isListedContainer(string(output), containerName), nil
}

func (d *DockerCliWrapperImpl) ImageExists(image string) (bool, error) {
	output, err := exec.Command("docker", "images", "-q", image).Output() // #nosec G204: fixed docker binary with structured args; image is not shell-expanded
	if err != nil {
		return false, Logger.NewError(err.Error(), "image", image)
	}
	return strings.TrimSpace(string(output)) != "", nil
}

func (d *DockerCliWrapperImpl) DeleteImage(image string) error {
	output, err := exec.Command("docker", "image", "rm", "-f", image).CombinedOutput() // #nosec G204: fixed docker binary with structured args; image is not shell-expanded
	if err != nil {
		return Logger.NewError(err.Error(), "image", image, "docker_error_message", strings.TrimSpace(string(output)))
	}
	return nil
}

func isListedContainer(output string, containerName string) bool {
	return strings.TrimSpace(output) == containerName
}
