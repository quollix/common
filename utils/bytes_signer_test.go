package utils

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
)

var crypto = &BytesSignerImpl{}

func TestSignAndVerifyBytes(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.Nil(t, err)

	payloadBytes := []byte(`{"hello":"world"}`)
	signature := crypto.SignBytes(privateKey, payloadBytes)

	isValid := crypto.VerifyBytes(publicKey, payloadBytes, signature)
	assert.True(t, isValid)

	tamperedBytes := []byte(`{"hello":"tampered"}`)
	isValid = crypto.VerifyBytes(publicKey, tamperedBytes, signature)
	assert.False(t, isValid)
}

func TestVerifyBytes_FailsWhenSignatureDoesNotMatch(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(nil)
	assert.Nil(t, err)

	isValid := crypto.VerifyBytes(publicKey, []byte(`{}`), []byte("invalid-signature"))
	assert.False(t, isValid)
}
