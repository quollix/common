package utils

import (
	"testing"

	"github.com/quollix/common/assert"
)

func TestLocalTestingKeysCanBeDecoded(t *testing.T) {
	publicKey, err := DecodeAuthorizedEd25519PublicKey(LocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(publicKey))

	productionPublicKey, err := DecodeAuthorizedEd25519PublicKey(LicenseTokenSigningPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(productionPublicKey))

	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	assert.Equal(t, 64, len(privateKey))
}
