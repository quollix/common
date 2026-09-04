package utils

import (
	"crypto/ed25519"

	"golang.org/x/crypto/ssh"
)

const localTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBe9GDTJJdlnnZqaNlLM7yK1yaq/6YAqXr8rDY48LVIoAAAAJiNquafjarm
nwAAAAtzc2gtZWQyNTUxOQAAACBe9GDTJJdlnnZqaNlLM7yK1yaq/6YAqXr8rDY48LVIoA
AAAEBy93cGXFFlR/PyHCdrCjXOjNGi52drotjlb6v7Egg1Tl70YNMkl2Wedmpo2UszvIrX
Jqr/pgCpevysNjjwtUigAAAADnN0b3JlLXRlc3Qta2V5AQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
`

const otherLocalTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCivB4dgjKyDGyYns+WN+ZYShn3XXPsTmuxIpNKsyt+rgAAAJBv1sTyb9bE
8gAAAAtzc2gtZWQyNTUxOQAAACCivB4dgjKyDGyYns+WN+ZYShn3XXPsTmuxIpNKsyt+rg
AAAEAZ0+KIT35e+H8CrvPJ30F238nVZ30ma2hAEQxjIdnHl6K8Hh2CMrIMbJiez5Y35lhK
Gfddc+xOa7Eik0qzK36uAAAAC2JhaWVyQHRpbmt5AQI=
-----END OPENSSH PRIVATE KEY-----
`

const (
	LocalTestingPublicKeyOpenSSH      = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo2UszvIrXJqr/pgCpevysNjjwtUig"
	OtherLocalTestingPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKK8Hh2CMrIMbJiez5Y35lhKGfddc+xOa7Eik0qzK36u"
)

var (
	LocalTestingPublicKeyOpenSSHBytes      = []byte(LocalTestingPublicKeyOpenSSH)
	OtherLocalTestingPublicKeyOpenSSHBytes = []byte(OtherLocalTestingPublicKeyOpenSSH)
)

func GetLocalTestingPrivateKeyOpenSSH() string {
	return localTestingPrivateKeyOpenSSH
}

func GetOtherLocalTestingPrivateKeyOpenSSH() string {
	return otherLocalTestingPrivateKeyOpenSSH
}

func GetLocalTestingPrivateKeyBytes() []byte {
	return []byte(localTestingPrivateKeyOpenSSH)
}

func GetOtherLocalTestingPrivateKeyBytes() []byte {
	return []byte(otherLocalTestingPrivateKeyOpenSSH)
}

func GetLocalTestingPublicKeyRaw() []byte {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetLocalTestingPrivateKeyBytes())
	if err != nil {
		panic(err)
	}
	return append([]byte(nil), privateKey.Public().(ed25519.PublicKey)...)
}

func GetOtherLocalTestingPublicKeyRaw() []byte {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetOtherLocalTestingPrivateKeyBytes())
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
