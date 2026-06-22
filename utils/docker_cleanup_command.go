package utils

import (
	"strings"

	"github.com/quollix/taskrunner"
	"github.com/spf13/cobra"
)

const protectedDockerContainerPrefix = "work-codex-"

func NewDockerRuntimeCleanupCommand(tr *taskrunner.TaskRunner) *cobra.Command {
	return &cobra.Command{
		Use:   "docker-cleanup",
		Short: "removes all Docker containers, volumes, and networks but keeps images",
		Long:  "removes all Docker containers, then prunes Docker volumes and networks while keeping images intact",
		Run: func(cmd *cobra.Command, args []string) {
			for _, container := range getCleanupContainers(tr) {
				tr.Cmd().AllowFail().Run("docker rm -f %s", container)
			}
			tr.Cmd().AllowFail().Run("docker volume prune -af")
			tr.Cmd().AllowFail().Run("docker network prune -f")
			tr.Log.Info("Docker runtime cleanup finished.")
		},
	}
}

func getCleanupContainers(tr *taskrunner.TaskRunner) []string {
	output := tr.Cmd().AllowFail().Run("docker ps -a --format '{{.Names}}'").Output()
	containers := strings.Fields(output)
	filtered := make([]string, 0, len(containers))
	for _, container := range containers {
		if strings.HasPrefix(container, protectedDockerContainerPrefix) {
			continue
		}
		filtered = append(filtered, container)
	}
	return filtered
}
