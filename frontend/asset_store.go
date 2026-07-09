package frontend

import (
	"fmt"
	"path"
	"strings"
	"sync"
)

type AssetStore struct {
	mutex      sync.RWMutex
	assetBytes map[string][]byte
	version    string
}

func NewAssetStore(version string) *AssetStore {
	return &AssetStore{
		assetBytes: map[string][]byte{},
		version:    version,
	}
}

func (store *AssetStore) Put(assetPath string, content []byte) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.assetBytes[assetPath] = content
}

func (store *AssetStore) Get(assetPath string) ([]byte, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	content, ok := store.assetBytes[assetPath]
	return content, ok
}

func (store *AssetStore) Has(assetPath string) bool {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	_, ok := store.assetBytes[assetPath]
	return ok
}

func (store *AssetStore) Clear() {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.assetBytes = map[string][]byte{}
}

func (store *AssetStore) VersionedPath(folderPath, fileNameWithoutExtension, fileExtension string) string {
	fileName := fmt.Sprintf("%s.%s.%s", path.Base(fileNameWithoutExtension), store.version, fileExtension)
	filePath := path.Join(path.Dir(fileNameWithoutExtension), fileName)

	folderPath = strings.TrimSuffix(folderPath, "/")
	if folderPath == "" || folderPath == "." {
		return filePath
	}
	return path.Join(folderPath, filePath)
}
