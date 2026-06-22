package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/quollix/taskrunner"
)

type ArtifactDiscovery struct {
	MockeryDirs []string
	WireDirs    []string
}

func BuildArtifacts(tr *taskrunner.TaskRunner, moduleRoot string) {
	err := buildArtifacts(moduleRoot, tr)
	if err != nil {
		Logger.Error(Logger.NewError(err.Error(), DirectoryField, moduleRoot, "result", "failed to generate artifacts"))
		tr.ExitWithError()
	}
}

func buildArtifacts(moduleRoot string, tr *taskrunner.TaskRunner) error {
	RemoveGeneratedArtifacts(moduleRoot)
	discovery, err := DiscoverArtifacts(moduleRoot)
	if err != nil {
		return err
	}
	for _, wireDir := range discovery.WireDirs {
		if hasWireHeaderFile(wireDir) {
			tr.Cmd().Dir(wireDir).Run("go tool wire gen -header_file wire_header.txt")
			continue
		}
		tr.Cmd().Dir(wireDir).Run("go tool wire")
	}
	for _, mockeryDir := range discovery.MockeryDirs {
		tr.Cmd().Dir(mockeryDir).Run("go tool mockery")
	}
	return nil
}

func DiscoverArtifacts(moduleRoot string) (ArtifactDiscovery, error) {
	mockeryDirSet := map[string]struct{}{}
	wireDirSet := map[string]struct{}{}
	wireHeaderDirSet := map[string]struct{}{}

	err := filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		switch d.Name() {
		case ".mockery.yml":
			mockeryDirSet[filepath.Dir(path)] = struct{}{}
		case "wire.go":
			wireDirSet[filepath.Dir(path)] = struct{}{}
		case "wire_header.txt":
			wireHeaderDirSet[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return ArtifactDiscovery{}, Logger.NewError(err.Error(), DirectoryField, moduleRoot)
	}

	for wireHeaderDir := range wireHeaderDirSet {
		if _, exists := wireDirSet[wireHeaderDir]; !exists {
			return ArtifactDiscovery{}, Logger.NewError(
				fmt.Sprintf("found wire_header.txt without wire.go in %s", wireHeaderDir),
				DirectoryField,
				moduleRoot,
			)
		}
	}

	discovery := ArtifactDiscovery{
		MockeryDirs: mapKeys(mockeryDirSet),
		WireDirs:    mapKeys(wireDirSet),
	}
	slices.Sort(discovery.MockeryDirs)
	slices.Sort(discovery.WireDirs)
	return discovery, nil
}

func RemoveGeneratedArtifacts(rootDir string) {
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !IsGeneratedArtifactFile(d.Name()) {
			return nil
		}
		return os.Remove(path) // #nosec G122: generated artifact paths come from the walked repository tree and are deleted by design
	})
	if err != nil {
		Logger.Error(Logger.NewError(err.Error(), DirectoryField, rootDir, "result", "failed to remove generated artifacts"))
		os.Exit(1)
	}
}

func IsGeneratedArtifactFile(fileName string) bool {
	return fileName == "wire_gen.go" || strings.HasSuffix(fileName, "_mock.go")
}

func hasWireHeaderFile(dir string) bool {
	fileInfo, err := os.Stat(filepath.Join(dir, "wire_header.txt"))
	return err == nil && !fileInfo.IsDir()
}

func mapKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
