package ci

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/quollix/common/assert"
)

func TestDiscoverModuleSets(t *testing.T) {
	srcDir := t.TempDir()
	makeModuleDir(t, srcDir, "ci-runner")
	makeModuleDir(t, srcDir, "server")
	makeModuleDir(t, srcDir, "client")
	assert.Nil(t, os.Mkdir(filepath.Join(srcDir, "docs"), 0o755))

	modules, err := discoverModuleSets(srcDir)
	assert.Nil(t, err)
	assert.Equal(t, []string{
		filepath.Join(srcDir, "client"),
		filepath.Join(srcDir, "server"),
	}, modules.artifacts)
}

func TestContainsGoMod(t *testing.T) {
	srcDir := t.TempDir()
	makeModuleDir(t, srcDir, "server")
	assert.Nil(t, os.Mkdir(filepath.Join(srcDir, "docs"), 0o755))

	assert.True(t, containsGoMod(filepath.Join(srcDir, "server")))
	assert.False(t, containsGoMod(filepath.Join(srcDir, "docs")))
}

func makeModuleDir(t *testing.T, parentDir string, name string) {
	t.Helper()
	moduleDir := filepath.Join(parentDir, name)
	assert.Nil(t, os.Mkdir(moduleDir, 0o755))
	assert.Nil(t, os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module sample\n\ngo 1.25.0\n"), 0o644))
}
