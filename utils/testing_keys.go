package utils

import (
	"crypto/ed25519"

	"golang.org/x/crypto/ssh"
)

const LocalTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABDeZszFO0
2bI7my6e4IRVlpAAAAAQAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo
2UszvIrXJqr/pgCpevysNjjwtUigAAAAoCEG/cUHNKYK56+EQuI0Vuq7QhYcM456zPtisS
uxGjl9iBjZQacg4UIOfWWfpXaxuDUoVYmZUIj03kgaGcI3vQiAu1PCSyWF8vf1QlvlTSJA
LUG4OpWv36w9ZS5d+OOjt+HJbOJUeo5W2zfUVegZQpPxJQQ7Mzhgmb22mmvcWSyGsQNuxi
bRvviIZTH/XcfkuFEWXYZeb6mGoipu4kAnSb8=
-----END OPENSSH PRIVATE KEY-----
`

const OtherLocalTestingPrivateKeyOpenSSH = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABDPIECaF+
on/c3uzTQfQAldAAAAAQAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIKK8Hh2CMrIMbJie
z5Y35lhKGfddc+xOa7Eik0qzK36uAAAAkJcISmnCbK4z7kYjF46Ndor0Q8neC8AF0duLsz
6YjThH9gi2NHYvKMZgjlO9W9Essj2K9CUOxZgKuwRvS2jt5jnZy6jqXQ4GnXEl8YPeXZdl
Af8+zFifRGonfu/YYBianeti39dYFrOnVDlrvJbxzIixwwXobKn2g9ov55tw3xF33hbBOT
EFisxHl/yW/DmU3g==
-----END OPENSSH PRIVATE KEY-----
`

const (
	LocalTestingPublicKeyOpenSSH           = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo2UszvIrXJqr/pgCpevysNjjwtUig"
	LocalTestingPublicKeyFingerprintSHA256 = "SHA256:QL91usdSz5KndtEmrv1z4p4KJTUpMA9Vqhqpzqduhbc"

	OtherLocalTestingPublicKeyOpenSSH           = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKK8Hh2CMrIMbJiez5Y35lhKGfddc+xOa7Eik0qzK36u"
	OtherLocalTestingPublicKeyFingerprintSHA256 = "SHA256:RVqr2+zRJn6XkjcOfl80D6eO3fJygFwqG+jsX4LZPYA"

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
