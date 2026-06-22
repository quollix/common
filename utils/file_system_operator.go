package utils

import (
	"context"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"gopkg.in/yaml.v3"
)

const composeConfigConsistencyCheckFailed = "compose config consistency check failed"

type ComposeSyntaxChecker interface {
	CheckDockerComposeSyntax(composeFileBytes []byte) error
}

type SimpleFile struct {
	Name  string
	IsDir bool
}

type FileSystemOperator interface {
	ListFiles(dir string) ([]SimpleFile, error)
	ReadYamlFile(dir string) (map[string]any, error)
	CheckDockerComposeSyntax(composeFileBytes []byte) error
}

type FileSystemOperatorImpl struct {
	OsWrapper OsWrapper
}

func NewFileSystemOperator(osWrapper OsWrapper) FileSystemOperator {
	return &FileSystemOperatorImpl{
		OsWrapper: osWrapper,
	}
}

func (f *FileSystemOperatorImpl) ListFiles(dir string) ([]SimpleFile, error) {
	files, err := f.OsWrapper.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var simpleFiles []SimpleFile
	for _, file := range files {
		simpleFiles = append(simpleFiles, SimpleFile{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		})
	}
	return simpleFiles, nil
}

func (f *FileSystemOperatorImpl) ReadYamlFile(path string) (map[string]any, error) {
	data, err := f.OsWrapper.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var yamlMap map[string]any
	if err := yaml.Unmarshal(data, &yamlMap); err != nil {
		return nil, Logger.NewError(err.Error())
	}

	return yamlMap, nil
}

func (f *FileSystemOperatorImpl) CheckDockerComposeSyntax(composeFileBytes []byte) error {
	configDetails := types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: "docker-compose.yml",
				Content:  composeFileBytes,
			},
		},
	}

	_, err := loader.LoadModelWithContext(context.Background(), configDetails, func(opts *loader.Options) {
		opts.SkipInterpolation = true
		opts.SkipInclude = true
		opts.SkipExtends = true
		opts.SkipResolveEnvironment = true
		opts.ResolvePaths = false
		opts.SetProjectName("compose_check", true)
	})
	if err != nil {
		return Logger.NewError(composeConfigConsistencyCheckFailed, "compose_error_message", err.Error())
	}
	return nil
}
