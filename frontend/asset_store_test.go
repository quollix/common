package frontend

import (
	"testing"

	"github.com/quollix/common/assert"
)

func TestAssetStore_Workflow(t *testing.T) {
	assetStore := NewAssetStore("abc123")
	assetPath := "/frontend/resources/global/frame.abc123.js"
	expectedBytes := []byte("hello")

	assert.False(t, assetStore.Has(assetPath))
	actualBytes, ok := assetStore.Get(assetPath)
	assert.False(t, ok)
	assert.Nil(t, actualBytes)

	assetStore.Put(assetPath, expectedBytes)

	actualBytes, ok = assetStore.Get(assetPath)
	assert.True(t, ok)
	assert.Equal(t, expectedBytes, actualBytes)

	assetStore.Clear()

	assert.False(t, assetStore.Has(assetPath))
}

func TestAssetStore_VersionedPath(t *testing.T) {
	assetStore := NewAssetStore("abc123")

	actualPath := assetStore.VersionedPath("/frontend/resources", "global/frame", "js")
	expectedPath := "/frontend/resources/global/frame.abc123.js"

	assert.Equal(t, expectedPath, actualPath)
}

func TestAssetStore_VersionedPath_WithoutFolderPath(t *testing.T) {
	assetStore := NewAssetStore("abc123")

	actualPath := assetStore.VersionedPath("", "global/frame", "js")

	assert.Equal(t, "global/frame.abc123.js", actualPath)
}

func TestAssetStore_VersionedPath_WithDotFolderPath(t *testing.T) {
	assetStore := NewAssetStore("abc123")

	actualPath := assetStore.VersionedPath(".", "global/frame", "js")

	assert.Equal(t, "global/frame.abc123.js", actualPath)
}

func TestAssetStore_VersionedPath_TrimsFolderPathSuffix(t *testing.T) {
	assetStore := NewAssetStore("abc123")

	actualPath := assetStore.VersionedPath("/frontend/resources/", "global/frame", "js")

	assert.Equal(t, "/frontend/resources/global/frame.abc123.js", actualPath)
}
