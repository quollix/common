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

const LocalTestingPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo2UszvIrXJqr/pgCpevysNjjwtUig"
const LicenseTokenSigningPublicKeyOpenSSH = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEC2+vXmZ/qX7vwL7Y2CWt65j3765Xe/n84E//XvS5h/" // #nosec G101: public key used by tests, not a credential

var LocalTestingPublicKeyOpenSSHBytes = []byte(LocalTestingPublicKeyOpenSSH)
var LicenseTokenSigningPublicKeyOpenSSHBytes = []byte(LicenseTokenSigningPublicKeyOpenSSH)

func GetLocalTestingPrivateKeyOpenSSH() string {
	return localTestingPrivateKeyOpenSSH
}

func GetLocalTestingPrivateKeyBytes() []byte {
	return []byte(localTestingPrivateKeyOpenSSH)
}

func GetLocalTestingPublicKeyRaw() []byte {
	privateKey, err := DecodeEd25519PrivateKeyOpenSSH(GetLocalTestingPrivateKeyBytes())
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
