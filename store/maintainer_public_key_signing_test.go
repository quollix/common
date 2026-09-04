package store

import (
	"crypto/ed25519"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

func TestMaintainerPublicKeySigning_HappyPath(t *testing.T) {
	privateKey, err := u.DecodeEd25519PrivateKeyOpenSSH(u.GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	publicKey := u.GetLocalTestingPublicKeyRaw()

	signature, err := SignMaintainerPublicKey(privateKey, "appuser", publicKey)
	assert.Nil(t, err)

	ok, err := VerifyMaintainerPublicKeySignature(privateKey.Public().(ed25519.PublicKey), &MaintainerPublicKeyRecord{
		Maintainer:         "appuser",
		PublicKeyRaw:       publicKey,
		PublicKeySignature: signature,
	})
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestMaintainerPublicKeySigning_DifferentMaintainerFailsVerification(t *testing.T) {
	privateKey, err := u.DecodeEd25519PrivateKeyOpenSSH(u.GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	publicKey := u.GetLocalTestingPublicKeyRaw()

	signature, err := SignMaintainerPublicKey(privateKey, "appuser", publicKey)
	assert.Nil(t, err)

	ok, err := VerifyMaintainerPublicKeySignature(privateKey.Public().(ed25519.PublicKey), &MaintainerPublicKeyRecord{
		Maintainer:         "other",
		PublicKeyRaw:       publicKey,
		PublicKeySignature: signature,
	})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestMaintainerPublicKeySigning_DifferentPublicKeyFailsVerification(t *testing.T) {
	privateKey, err := u.DecodeEd25519PrivateKeyOpenSSH(u.GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	publicKey := u.GetLocalTestingPublicKeyRaw()
	differentPublicKey := u.GetOtherLocalTestingPublicKeyRaw()

	signature, err := SignMaintainerPublicKey(privateKey, "appuser", publicKey)
	assert.Nil(t, err)

	ok, err := VerifyMaintainerPublicKeySignature(privateKey.Public().(ed25519.PublicKey), &MaintainerPublicKeyRecord{
		Maintainer:         "appuser",
		PublicKeyRaw:       differentPublicKey,
		PublicKeySignature: signature,
	})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestMaintainerPublicKeySigning_DifferentAdminPublicKeyFailsVerification(t *testing.T) {
	privateKey, err := u.DecodeEd25519PrivateKeyOpenSSH(u.GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)
	publicKey := u.GetLocalTestingPublicKeyRaw()
	differentAdminPublicKey := ed25519.PublicKey(u.GetOtherLocalTestingPublicKeyRaw())

	signature, err := SignMaintainerPublicKey(privateKey, "appuser", publicKey)
	assert.Nil(t, err)

	ok, err := VerifyMaintainerPublicKeySignature(differentAdminPublicKey, &MaintainerPublicKeyRecord{
		Maintainer:         "appuser",
		PublicKeyRaw:       publicKey,
		PublicKeySignature: signature,
	})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestMaintainerPublicKeySigning_RejectsInvalidSignatureLength(t *testing.T) {
	privateKey, err := u.DecodeEd25519PrivateKeyOpenSSH(u.GetLocalTestingPrivateKeyBytes())
	assert.Nil(t, err)

	ok, err := VerifyMaintainerPublicKeySignature(privateKey.Public().(ed25519.PublicKey), &MaintainerPublicKeyRecord{
		Maintainer:         "appuser",
		PublicKeyRaw:       u.GetLocalTestingPublicKeyRaw(),
		PublicKeySignature: []byte("short"),
	})
	assert.False(t, ok)
	assert.Equal(t, "invalid ed25519 signature", u.ExtractError(err))
}
