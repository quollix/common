package utils

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
	"golang.org/x/crypto/ssh"
)

const unprotectedTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBe9GDTJJdlnnZqaNlLM7yK1yaq/6YAqXr8rDY48LVIoAAAAJiNquafjarm
nwAAAAtzc2gtZWQyNTUxOQAAACBe9GDTJJdlnnZqaNlLM7yK1yaq/6YAqXr8rDY48LVIoA
AAAEBy93cGXFFlR/PyHCdrCjXOjNGi52drotjlb6v7Egg1Tl70YNMkl2Wedmpo2UszvIrX
Jqr/pgCpevysNjjwtUigAAAADnN0b3JlLXRlc3Qta2V5AQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
`

func TestDecodeEd25519PrivateKeyOpenSSH_WrongPassphraseReturnsError(t *testing.T) {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(LocalTestingPrivateKeyOpenSSH), []byte("wrong"))

	assert.NotNil(t, err)
	assert.Nil(t, privateKey)
}

func TestDecodeEd25519PrivateKeyOpenSSH_UnprotectedKeyReturnsError(t *testing.T) {
	parsedKey, err := ssh.ParseRawPrivateKey([]byte(unprotectedTestingPrivateKeyOpenSSH))
	assert.Nil(t, err)
	_, ok := parsedKey.(*ed25519.PrivateKey)
	assert.True(t, ok)

	privateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(unprotectedTestingPrivateKeyOpenSSH), []byte(LocalTestingPrivateKeyPassphrase))

	assert.NotNil(t, err)
	assert.Nil(t, privateKey)
}

func TestIsPrivateKeyPassphraseProtectedOpenSSH_ProtectedKeyReturnsTrue(t *testing.T) {
	protected, err := IsPrivateKeyPassphraseProtectedOpenSSH([]byte(LocalTestingPrivateKeyOpenSSH))

	assert.Nil(t, err)
	assert.True(t, protected)
}

func TestIsPrivateKeyPassphraseProtectedOpenSSH_UnprotectedKeyReturnsFalse(t *testing.T) {
	protected, err := IsPrivateKeyPassphraseProtectedOpenSSH([]byte(unprotectedTestingPrivateKeyOpenSSH))

	assert.Nil(t, err)
	assert.False(t, protected)
}

func TestIsPrivateKeyPassphraseProtectedOpenSSH_InvalidKeyReturnsError(t *testing.T) {
	protected, err := IsPrivateKeyPassphraseProtectedOpenSSH([]byte("invalid"))

	assert.NotNil(t, err)
	assert.False(t, protected)
}
