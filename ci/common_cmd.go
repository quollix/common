package ci

import (
	"os"
	"path/filepath"
	"sort"

	u "github.com/quollix/common/utils"
	"github.com/quollix/taskrunner"
	"github.com/spf13/cobra"
)

type moduleSets struct {
	artifacts []string
}

func NewCommonCmd(tr *taskrunner.TaskRunner, srcDir string, _ string) *cobra.Command {
	commonCmd := &cobra.Command{
		Use:   "common",
		Short: "shared repository maintenance commands",
	}

	commonCmd.AddCommand(
		newArtifactsCmd(tr, srcDir),
		u.NewDockerRuntimeCleanupCommand(tr),
	)
	return commonCmd
}

func newArtifactsCmd(tr *taskrunner.TaskRunner, srcDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "artifacts",
		Short: "generate artifacts for all non-ci-runner Go modules in src",
		RunE: func(cmd *cobra.Command, args []string) error {
			modules, err := discoverModuleSets(srcDir)
			if err != nil {
				return err
			}
			for _, moduleDir := range modules.artifacts {
				u.BuildArtifacts(tr, moduleDir)
			}
			return nil
		},
	}
}

func discoverModuleSets(srcDir string) (moduleSets, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return moduleSets{}, u.Logger.NewError(err.Error(), u.DirectoryField, srcDir)
	}
	modules := moduleSets{
		artifacts: make([]string, 0, len(entries)),
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		moduleDir := filepath.Join(srcDir, entry.Name())
		if !containsGoMod(moduleDir) {
			continue
		}
		if entry.Name() != "ci-runner" {
			modules.artifacts = append(modules.artifacts, moduleDir)
		}
	}
	sort.Strings(modules.artifacts)
	return modules, nil
}

func containsGoMod(dir string) bool {
	fileInfo, err := os.Stat(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}
