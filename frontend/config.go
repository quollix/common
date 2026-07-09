package frontend

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/quollix/deepstack"
)

var logger = deepstack.NewDeepStackLogger()

const (
	resourceFolderPath = "resources"
	framePath          = "global/frame.html"
)

type Config struct {
	FrontendFolderPath string
	Version            string
	Static             any
}

type resolvedConfig struct {
	fileSystem         fs.FS
	publicResourcePath string
	framePath          string
	version            string
	static             any
}

func (c Config) resolve() (resolvedConfig, error) {
	frontendFolderPath := strings.TrimSpace(c.FrontendFolderPath)
	if frontendFolderPath == "" {
		return resolvedConfig{}, logger.NewError("frontend folder path must not be empty")
	}
	frontendFolderPath = filepath.Clean(frontendFolderPath)

	resourceFolderPathOnDisk := filepath.Join(frontendFolderPath, resourceFolderPath)
	resourceFolderInfo, err := os.Stat(resourceFolderPathOnDisk)
	if err != nil {
		return resolvedConfig{}, logger.NewError(err.Error(), "resource_folder_path", resourceFolderPathOnDisk)
	}
	if !resourceFolderInfo.IsDir() {
		return resolvedConfig{}, logger.NewError("resource folder path must be a directory", "resource_folder_path", resourceFolderPathOnDisk)
	}

	return resolvedConfig{
		fileSystem:         os.DirFS(resourceFolderPathOnDisk),
		publicResourcePath: "/" + path.Base(filepath.ToSlash(frontendFolderPath)) + "/" + resourceFolderPath,
		framePath:          framePath,
		version:            c.Version,
		static:             c.Static,
	}, nil
}
