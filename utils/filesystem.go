package utils

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/quollix/taskrunner"
)

func FindDir(dirName string) (string, error) {
	currentDir, err := os.Getwd()
	initialDir := currentDir
	if err != nil {
		return "", Logger.NewError(err.Error(), "directory_to_find", dirName)
	}

	for {
		candidatePath := filepath.Join(currentDir, dirName)
		fileInfo, err := os.Stat(candidatePath)
		if err == nil && fileInfo.IsDir() {
			return candidatePath, nil
		}

		parentDir := filepath.Dir(currentDir)
		isRootFolderReached := parentDir == currentDir
		if isRootFolderReached {
			return "", Logger.NewError("folder not found in any parent directory", FolderToFindField, dirName, InitialDirField, initialDir)
		}

		currentDir = parentDir
	}
}

func RemoveDir(path string) {
	RunAndLogIfError(func() error {
		return os.RemoveAll(path)
	})
}

func CollectBuildTags(dir string) ([]string, error) {
	var tags []string
	tagSet := make(map[string]struct{})
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
			f, err := os.Open(path) // #nosec G304 G122: walk input is the intended repository tree and each matched file must be opened
		if err != nil {
			return Logger.NewError(err.Error())
		}
		defer Close(f)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			extractedTags := extractTagsFromLine(line)
			for _, tag := range extractedTags {
				tagSet[tag] = struct{}{}
			}
		}
		err = scanner.Err()
		if err != nil {
			return Logger.NewError(err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags, nil
}

func extractTagsFromLine(line string) []string {
	var tags []string
	if strings.HasPrefix(line, "//go:build") {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			tags = fields[1:]
		}
	}
	if strings.HasPrefix(line, "// +build") {
		fields := strings.Fields(line)
		if len(fields) > 2 {
			tags = fields[2:]
		}
	}

	var filteredTags []string
	for _, tag := range tags {
		if strings.HasPrefix(tag, "!") {
			continue
		}
		filteredTags = append(filteredTags, tag)
	}

	return filteredTags
}

// BuildWholeGoProject this function is meant to detect compile errors in the test code with build tags
func BuildWholeGoProject(Tr *taskrunner.TaskRunner, directory string) {
	buildTags, err := CollectBuildTags(directory)
	if err != nil {
		Tr.Log.Error("Error collecting build tags: %v", err)
		os.Exit(1)
	}
	commaSeparatedBuildTags := strings.Join(buildTags, ",")
	command := fmt.Sprintf("go test -count=1 -tags=%s -run=NO_SUCH_TEST ./...", commaSeparatedBuildTags)
	Tr.Log.Info("checking whether the entire production and test code can be compiled without running tests...")
	Tr.Cmd().Dir(directory).Run("%s", command)
}
