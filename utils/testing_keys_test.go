package utils

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
)

func TestLocalTestingKeysCanBeDecoded(t *testing.T) {
	publicKey, err := DecodeAuthorizedEd25519PublicKey(LocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(publicKey))

	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	assert.Equal(t, 64, len(privateKey))

	otherPublicKey, err := DecodeAuthorizedEd25519PublicKey(OtherLocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(otherPublicKey))

	otherPrivateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetOtherLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	assert.Equal(t, 64, len(otherPrivateKey))
	assert.Equal(t, otherPublicKey, otherPrivateKey.Public().(ed25519.PublicKey))
}
