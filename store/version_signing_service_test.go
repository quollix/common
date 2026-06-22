package store

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

type versionSigningServiceTestDependencies struct {
	service     VersionSigningService
	serviceImpl *VersionSigningServiceImpl
	codec       *VersionSigningCodecMock
	bytesSigner *u.BytesSignerMock
}

func setupVersionSigningServiceTestDependencies(t *testing.T) *versionSigningServiceTestDependencies {
	codec := NewVersionSigningCodecMock(t)
	bytesSigner := u.NewBytesSignerMock(t)
	serviceImpl := &VersionSigningServiceImpl{
		Codec:       codec,
		BytesSigner: bytesSigner,
	}

	return &versionSigningServiceTestDependencies{
		service:     serviceImpl,
		serviceImpl: serviceImpl,
		codec:       codec,
		bytesSigner: bytesSigner,
	}
}

func TestSignVersion_InvalidPrivateKeyReturnsError(t *testing.T) {
	deps := setupVersionSigningServiceTestDependencies(t)

	signature, err := deps.service.SignVersion(ed25519.PrivateKey("short"), getSampleVersion())
	assert.Nil(t, signature)
	assert.Equal(t, "invalid ed25519 private key", u.ExtractError(err))
}

func TestSignVersion_HappyPath(t *testing.T) {
	deps := setupVersionSigningServiceTestDependencies(t)
	version := getSampleVersion()
	privateKey := ed25519.PrivateKey(make([]byte, ed25519.PrivateKeySize))
	payloadBytes := []byte("payload-bytes")
	expectedSignature := []byte("signature")

	deps.codec.EXPECT().EncodeVersion(version).Return(payloadBytes, nil)
	deps.bytesSigner.EXPECT().SignBytes(privateKey, payloadBytes).Return(expectedSignature)

	signature, err := deps.service.SignVersion(privateKey, version)
	assert.Nil(t, err)
	assert.Equal(t, expectedSignature, signature)
}

func TestVerifyVersionSignature_InvalidPublicKeyReturnsError(t *testing.T) {
	deps := setupVersionSigningServiceTestDependencies(t)

	ok, err := deps.service.VerifyVersionSignature(ed25519.PublicKey("short"), getSampleVersion())
	assert.False(t, ok)
	assert.Equal(t, "invalid ed25519 public key", u.ExtractError(err))
}

func TestVerifyVersionSignature_HappyPath(t *testing.T) {
	deps := setupVersionSigningServiceTestDependencies(t)
	version := getSampleVersion()
	version.Signature = []byte("signature")
	publicKey := ed25519.PublicKey(make([]byte, ed25519.PublicKeySize))
	payloadBytes := []byte("payload-bytes")

	deps.codec.EXPECT().EncodeVersion(version).Return(payloadBytes, nil)
	deps.bytesSigner.EXPECT().VerifyBytes(publicKey, payloadBytes, version.Signature).Return(true)

	ok, err := deps.service.VerifyVersionSignature(publicKey, version)
	assert.Nil(t, err)
	assert.True(t, ok)
}
