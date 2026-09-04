package store

import (
	"bytes"
	"crypto/ed25519"

	u "github.com/quollix/common/utils"
)

func EncodeMaintainerPublicKeyPayload(maintainer string, publicKeyRaw []byte) ([]byte, error) {
	if maintainer == "" {
		return nil, u.Logger.NewError("maintainer must not be empty")
	}
	if len(publicKeyRaw) != ed25519.PublicKeySize {
		return nil, u.Logger.NewError("invalid ed25519 public key")
	}

	buf := &bytes.Buffer{}
	if err := writeVersionSigningField(buf, []byte(maintainer)); err != nil {
		return nil, err
	}
	if err := writeVersionSigningField(buf, publicKeyRaw); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SignMaintainerPublicKey(privateKey ed25519.PrivateKey, maintainer string, publicKeyRaw []byte) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, u.Logger.NewError("invalid ed25519 private key")
	}
	payloadBytes, err := EncodeMaintainerPublicKeyPayload(maintainer, publicKeyRaw)
	if err != nil {
		return nil, err
	}
	return ed25519.Sign(privateKey, payloadBytes), nil
}

func VerifyMaintainerPublicKeySignature(adminPublicKey ed25519.PublicKey, record *MaintainerPublicKeyRecord) (bool, error) {
	if len(adminPublicKey) != ed25519.PublicKeySize {
		return false, u.Logger.NewError("invalid ed25519 public key")
	}
	if record == nil {
		return false, u.Logger.NewError("maintainer public key record must not be nil")
	}
	if len(record.PublicKeySignature) != ed25519.SignatureSize {
		return false, u.Logger.NewError("invalid ed25519 signature")
	}
	payloadBytes, err := EncodeMaintainerPublicKeyPayload(record.Maintainer, record.PublicKeyRaw)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(adminPublicKey, payloadBytes, record.PublicKeySignature), nil
}
