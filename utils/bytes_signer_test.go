package utils

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
)

var crypto = &BytesSignerImpl{}

func TestSignAndVerifyBytes(t *testing.T) {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	payloadBytes := []byte(`{"hello":"world"}`)
	signature := crypto.SignBytes(privateKey, payloadBytes)

	isValid := crypto.VerifyBytes(publicKey, payloadBytes, signature)
	assert.True(t, isValid)

	tamperedBytes := []byte(`{"hello":"tampered"}`)
	isValid = crypto.VerifyBytes(publicKey, tamperedBytes, signature)
	assert.False(t, isValid)
}

func TestVerifyBytes_FailsWhenSignatureDoesNotMatch(t *testing.T) {
	publicKey := ed25519.PublicKey(GetOtherLocalTestingPublicKeyRaw())

	isValid := crypto.VerifyBytes(publicKey, []byte(`{}`), []byte("invalid-signature"))
	assert.False(t, isValid)
}
