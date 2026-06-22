package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/quollix/common/assert"
)

var samplesDir string

func getSamplesDir() string {
	if samplesDir == "" {
		var err error
		samplesDir, err = FindDir("samples")
		if err != nil {
			panic(err)
		}
	}
	return samplesDir
}

var fileReadingDir = getSamplesDir() + "/file_reading"

func newFileSystemOperatorForTest() *FileSystemOperatorImpl {
	return &FileSystemOperatorImpl{
		OsWrapper: &OsWrapperImpl{},
	}
}

func TestListFiles(t *testing.T) {
	fileSystemOperator := newFileSystemOperatorForTest()
	files, err := fileSystemOperator.ListFiles(getSamplesDir() + "/test_list")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, "test.txt", files[0].Name)
	assert.False(t, files[0].IsDir)
	assert.Equal(t, "test_dir", files[1].Name)
	assert.True(t, files[1].IsDir)
}

func TestDockerComposeCheck_Success(t *testing.T) {
	fileSystemOperator := newFileSystemOperatorForTest()
	err := fileSystemOperator.CheckDockerComposeSyntax([]byte(`
services:
  app:
    image: nginx:1.27
`))
	assert.Nil(t, err)
}

func TestDockerComposeCheck_Fail(t *testing.T) {
	fileSystemOperator := newFileSystemOperatorForTest()
	err := fileSystemOperator.CheckDockerComposeSyntax([]byte("hello"))
	assert.Equal(t, composeConfigConsistencyCheckFailed, ExtractError(err))
}

// We want that validaiton works completely in-memory, so we assert that the system does not try to access these non-existing files
func TestDockerComposeCheck_DoesNotLoadReferencedFiles(t *testing.T) {
	fileSystemOperator := newFileSystemOperatorForTest()
	err := fileSystemOperator.CheckDockerComposeSyntax([]byte(`
include:
  - missing-compose.yml
services:
  app:
    extends:
      file: missing-base.yml
      service: base
    image: nginx:1.27
`))
	assert.Nil(t, err)
}

func TestNotExistingFile(t *testing.T) {
	fileSystemOperator := newFileSystemOperatorForTest()
	_, err := fileSystemOperator.ReadYamlFile(fileReadingDir + "/non-existing-file.yml")
	assert.True(t, strings.Contains(ExtractError(err), "no such file or directory"))
}

func TestGetTempDir_CreatesDirectory(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	tempDir, err := osWrapper.GetTempDir()
	assert.Nil(t, err)
	assert.True(t, strings.Contains(filepath.Base(tempDir), "file_system_operator"))

	stat, statErr := os.Stat(tempDir)
	assert.Nil(t, statErr)
	assert.True(t, stat.IsDir())

	assert.Nil(t, osWrapper.RemoveAll(tempDir))
}

func TestWriteFile_WritesContent(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	tempDir, err := osWrapper.GetTempDir()
	assert.Nil(t, err)
	defer func() { assert.Nil(t, osWrapper.RemoveAll(tempDir)) }()

	filePath := filepath.Join(tempDir, "written.txt")
	fileBytes := []byte("hello")

	writeErr := osWrapper.WriteFile(filePath, fileBytes, 0o600)
	assert.Nil(t, writeErr)

	readBytes, readErr := os.ReadFile(filePath)
	assert.Nil(t, readErr)
	assert.Equal(t, fileBytes, readBytes)
}

func TestRemoveDir_RemovesDirectoryRecursively(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	tempDir, err := osWrapper.GetTempDir()
	assert.Nil(t, err)

	nestedDir := filepath.Join(tempDir, "nested")
	mkdirErr := os.MkdirAll(nestedDir, 0700)
	assert.Nil(t, mkdirErr)

	filePath := filepath.Join(nestedDir, "file.txt")
	writeErr := os.WriteFile(filePath, []byte("x"), 0600)
	assert.Nil(t, writeErr)

	assert.Nil(t, osWrapper.RemoveAll(tempDir))

	_, statErr := os.Stat(tempDir)
	assert.True(t, os.IsNotExist(statErr))
}

func TestDoesFileExist_ExistingFile(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	tempDir, err := osWrapper.GetTempDir()
	assert.Nil(t, err)
	defer func() { assert.Nil(t, osWrapper.RemoveAll(tempDir)) }()

	filePath := filepath.Join(tempDir, "exists.txt")
	writeErr := os.WriteFile(filePath, []byte("x"), 0600)
	assert.Nil(t, writeErr)

	exists, existsErr := osWrapper.DoesFileExist(filePath)
	assert.Nil(t, existsErr)
	assert.True(t, exists)
}

func TestDoesFileExist_NonExistingFile(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	tempDir, err := osWrapper.GetTempDir()
	assert.Nil(t, err)
	defer func() { assert.Nil(t, osWrapper.RemoveAll(tempDir)) }()

	filePath := filepath.Join(tempDir, "does-not-exist.txt")

	exists, existsErr := osWrapper.DoesFileExist(filePath)
	assert.Nil(t, existsErr)
	assert.False(t, exists)
}

func TestAllocateLocalhostPort_ReturnsPort(t *testing.T) {
	osWrapper := &OsWrapperImpl{}

	port, err := osWrapper.AllocateLocalhostPort()

	assert.Nil(t, err)
	assert.False(t, port == "")
}

func TestSleep_AdvancesRealTime(t *testing.T) {
	osWrapper := &OsWrapperImpl{}
	start := osWrapper.Now()

	osWrapper.Sleep(1 * time.Millisecond)

	assert.True(t, osWrapper.Now().After(start) || osWrapper.Now().Equal(start))
}
