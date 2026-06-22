package store

import (
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

func TestAppStoreOfficialMaintainerPublicKeyOpenSSHCanBeDecoded(t *testing.T) {
	publicKey, err := u.DecodeAuthorizedEd25519PublicKey([]byte(AppStoreOfficialMaintainerPublicKeyOpenSSH))
	assert.Nil(t, err)
	assert.Equal(t, 32, len(publicKey))
}
