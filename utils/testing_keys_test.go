package utils

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
	"golang.org/x/crypto/ssh"
)

func TestLocalTestingKeysCanBeDecoded(t *testing.T) {
	publicKey, err := DecodeAuthorizedEd25519PublicKey(LocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(publicKey))

	privateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(LocalTestingPrivateKeyOpenSSH), []byte(LocalTestingPrivateKeyPassphrase))
	assert.Nil(t, err)
	assert.Equal(t, 64, len(privateKey))

	otherPublicKey, err := DecodeAuthorizedEd25519PublicKey(OtherLocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(otherPublicKey))

	otherPrivateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(OtherLocalTestingPrivateKeyOpenSSH), []byte(OtherLocalTestingPrivateKeyPassphrase))
	assert.Nil(t, err)
	assert.Equal(t, 64, len(otherPrivateKey))
	assert.Equal(t, otherPublicKey, otherPrivateKey.Public().(ed25519.PublicKey))
}

func TestLocalTestingPublicKeyFingerprintsMatchFixtures(t *testing.T) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(LocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)
	otherPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(OtherLocalTestingPublicKeyOpenSSHBytes)
	assert.Nil(t, err)

	assert.Equal(t, LocalTestingPublicKeyFingerprintSHA256, ssh.FingerprintSHA256(publicKey))
	assert.Equal(t, OtherLocalTestingPublicKeyFingerprintSHA256, ssh.FingerprintSHA256(otherPublicKey))
}
