package store

import (
	"crypto/ed25519"

	u "github.com/quollix/common/utils"
)

type VersionSigningService interface {
	SignVersion(privateKey ed25519.PrivateKey, version *Version) ([]byte, error)
	VerifyVersionSignature(publicKey ed25519.PublicKey, version *Version) (bool, error)
}

type VersionSigningServiceImpl struct {
	Codec       VersionSigningCodec
	BytesSigner u.BytesSigner
}

func (s *VersionSigningServiceImpl) SignVersion(privateKey ed25519.PrivateKey, version *Version) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, u.Logger.NewError("invalid ed25519 private key")
	}
	payloadBytes, err := s.Codec.EncodeVersion(version)
	if err != nil {
		return nil, err
	}
	return s.BytesSigner.SignBytes(privateKey, payloadBytes), nil
}

func (s *VersionSigningServiceImpl) VerifyVersionSignature(publicKey ed25519.PublicKey, version *Version) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, u.Logger.NewError("invalid ed25519 public key")
	}
	payloadBytes, err := s.Codec.EncodeVersion(version)
	if err != nil {
		return false, err
	}
	return s.BytesSigner.VerifyBytes(publicKey, payloadBytes, version.Signature), nil
}
