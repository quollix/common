package utils

import (
	"gopkg.in/yaml.v3"
)

type SimpleFile struct {
	Name  string
	IsDir bool
}

type FileSystemOperator interface {
	ListFiles(dir string) ([]SimpleFile, error)
	ReadYamlFile(dir string) (map[string]any, error)
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
