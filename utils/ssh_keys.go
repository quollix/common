package utils

import (
	"crypto/ed25519"

	"golang.org/x/crypto/ssh"
)

func DecodeAuthorizedEd25519PublicKey(authorizedKey []byte) (ed25519.PublicKey, error) {
	parsedKey, _, _, _, err := ssh.ParseAuthorizedKey(authorizedKey)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	cryptoPublicKey, ok := parsedKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, Logger.NewError("ssh public key is not convertible to crypto public key")
	}
	publicKey, ok := cryptoPublicKey.CryptoPublicKey().(ed25519.PublicKey)
	if !ok {
		return nil, Logger.NewError("ssh public key is not an ed25519 public key")
	}
	return append(ed25519.PublicKey(nil), publicKey...), nil
}

func DecodeEd25519PrivateKeyOpenSSH(privateKeyOpenSSH []byte) (ed25519.PrivateKey, error) {
	privateKey, err := ssh.ParseRawPrivateKey(privateKeyOpenSSH)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	switch key := privateKey.(type) {
	case ed25519.PrivateKey:
		return key, nil
	case *ed25519.PrivateKey:
		return *key, nil
	default:
		return nil, Logger.NewError("test private key is not an ed25519 private key")
	}
}
