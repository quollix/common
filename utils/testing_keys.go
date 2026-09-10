package utils

import (
	"crypto/ed25519"

	"golang.org/x/crypto/ssh"
)

const LocalTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABB0i2GK2a
2Yp5F/MHhsax9jAAAAZAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo
2UszvIrXJqr/pgCpevysNjjwtUigAAAAoNEJgl4UxZrbOp3IDxwGGbR2CZLfDysihPQM/x
TxPiSykOXb7aVcMgxhp2YjWv2LiXk8aKRov+2wt+5VnNys++Adu7bQlt9vFyCbxf37pF7R
RcJSwxz/oGtxR1xSpVoxGu78HX8hUdc8Kfgl+mDHjD/BQVvwPaxhNds8DCsFUwLTrCZDLM
JXgeDrX5FEziouuRnYxpsbVeYtwhG3dIrs+XM=
-----END OPENSSH PRIVATE KEY-----
`

const OtherLocalTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABAtHGz/b/
9rlJJo4ngG1TPUAAAAZAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIKK8Hh2CMrIMbJie
z5Y35lhKGfddc+xOa7Eik0qzK36uAAAAkMPJ78XVGy5Z0WHshkZl3Tp1NMrf/pfr3PX8Kf
THeWEfGqVw38qCWrpheTxJ2uDKKgw+CaqUsc/elqy0e2iI3udqRbXjoaiDwTXl0iSeAIQs
jE6xw9T2HZIJNrZbdS22ktMH5Rt6L6u9hdHVDQgexkA6G3N6ifcl+q+o7GrLKhf+Sn+JKa
or5WpeDeioct72sA==
-----END OPENSSH PRIVATE KEY-----
`

const (
	LocalTestingPublicKeyOpenSSH      = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo2UszvIrXJqr/pgCpevysNjjwtUig"
	OtherLocalTestingPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKK8Hh2CMrIMbJiez5Y35lhKGfddc+xOa7Eik0qzK36u"

	LocalTestingPrivateKeyPassphrase      = "password1"
	OtherLocalTestingPrivateKeyPassphrase = "password2"
)

var (
	LocalTestingPublicKeyOpenSSHBytes      = []byte(LocalTestingPublicKeyOpenSSH)
	OtherLocalTestingPublicKeyOpenSSHBytes = []byte(OtherLocalTestingPublicKeyOpenSSH)
)

func GetLocalTestingPublicKeyRaw() []byte {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(LocalTestingPrivateKeyOpenSSH), []byte(LocalTestingPrivateKeyPassphrase))
	if err != nil {
		panic(err)
	}
	return append([]byte(nil), privateKey.Public().(ed25519.PublicKey)...)
}

func GetOtherLocalTestingPublicKeyRaw() []byte {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH([]byte(OtherLocalTestingPrivateKeyOpenSSH), []byte(OtherLocalTestingPrivateKeyPassphrase))
	if err != nil {
		panic(err)
	}
	return append([]byte(nil), privateKey.Public().(ed25519.PublicKey)...)
}

func GetLocalTestingPublicKeyFingerprintSHA256() string {
	sshPublicKey, err := ssh.NewPublicKey(ed25519.PublicKey(GetLocalTestingPublicKeyRaw()))
	if err != nil {
		panic(err)
	}
	return ssh.FingerprintSHA256(sshPublicKey)
}
