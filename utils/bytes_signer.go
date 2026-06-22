package utils

import (
	"crypto/ed25519"
)

type BytesSigner interface {
	SignBytes(privateKey ed25519.PrivateKey, payloadBytes []byte) []byte
	VerifyBytes(publicKey ed25519.PublicKey, payloadBytes []byte, signature []byte) bool
}

type BytesSignerImpl struct{}

func (c *BytesSignerImpl) SignBytes(privateKey ed25519.PrivateKey, payloadBytes []byte) []byte {
	return ed25519.Sign(privateKey, payloadBytes)
}

func (c *BytesSignerImpl) VerifyBytes(publicKey ed25519.PublicKey, payloadBytes []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, payloadBytes, signature)
}
