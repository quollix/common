package validation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

const (
	sixtyThreeHexDecimalLetters = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde"
)

type testCaseType struct {
	value       string
	expectError bool
}

func assertCases(t *testing.T, tag string, cases []testCaseType) {
	for _, c := range cases {
		err := Validate("Value", tag, c.value)
		if c.expectError {
			if err == nil {
				fmt.Printf("expected an error, but was nil: %s -> %s\n", tag, c.value)
				t.Fail()
			}
		} else {
			if err != nil {
				fmt.Printf("expected nil, but an errors occurred: %s -> %s\n", tag, c.value)
				t.Fail()
			}
		}
	}
}

func TestValidateDefault(t *testing.T) {
	cases := []testCaseType{
		{"validusername", false},
		{"user123", false},
		{"user.123", true},
		{"user-123", true},
		{"user_123", true},
		{"InvalidUsername", true},
		{"user!@#", true},
		{"us", true},
		{"thisusernameiswaytoolong", true},
	}
	assertCases(t, FieldDefault, cases)
}

func TestValidateUsername(t *testing.T) {
	cases := []testCaseType{
		{"validusername", false},
		{"user123", false},
		{"user-123", false},
		{"user_123", false},
		{"user.123", true},
		{"InvalidUsername", true},
		{"user!@#", true},
		{"us", true},
		{"thisusernameiswaytoolong", true},
	}
	assertCases(t, FieldUsername, cases)
}

func TestValidateVersion(t *testing.T) {
	cases := []testCaseType{
		{"valid.versionname", false},
		{"version123", false},
		{"version.name123", false},
		{"version_name123", true},
		{"invalid.versionname!", true},
		{"ta", true},
		{"this.versionname.is.way.too.long", true},
	}
	assertCases(t, FieldVersionName, cases)
}

func TestValidateFileName(t *testing.T) {
	cases := []testCaseType{
		{"a", false},
		{"file-name_01.txt", false},
		{"name.with.many.parts-and_underscores123", false},
		{"", true},
		{"name with spaces.txt", true},
		{"name/with/slash.txt", true},
		{"name?with?question.txt", true},
		{strings.Repeat("a", 100), false},
		{strings.Repeat("a", 101), true},
	}
	assertCases(t, FieldFileName, cases)
}

func TestValidatePassword(t *testing.T) {
	cases := []testCaseType{
		{"validpassword._-", false},
		{"validpassword!", true},
		{"valid_pass123", false},
		{"InvalidPassword", false}, // uppercase allowed by regex
		{"valid!@#", true},
		{"1234567", true},
		{"12345678", false},
		{"thispasswordiswaytoolong_xxxxx!", true},
	}
	assertCases(t, FieldPassword, cases)
}

func TestValidateCookie(t *testing.T) {
	cases := []testCaseType{
		{sixtyThreeHexDecimalLetters, true},
		{sixtyThreeHexDecimalLetters + "f", false},
		{sixtyThreeHexDecimalLetters + "ff", true},
		{sixtyThreeHexDecimalLetters + "g", true},
		{"", true},
	}
	assertCases(t, FieldSecret, cases)
}

func TestValidateComposeSecretName(t *testing.T) {
	cases := []testCaseType{
		{"SECRET_POSTGRES_PASSWORD", false},
		{"SECRET_1", false},
		{"SECRET_VALUE_123", false},
		{"SECRET_", true},
		{"POSTGRES_PASSWORD", true},
		{"secret_POSTGRES_PASSWORD", true},
		{"SECRET-postgres-password", true},
		{"SECRET_POSTGRES_PASSWORD!", true},
		{strings.Repeat("A", 136), true},
	}
	assertCases(t, FieldComposeSecret, cases)
}

func TestValidateSshPublicKey(t *testing.T) {
	cases := []testCaseType{
		{u.LocalTestingPublicKeyOpenSSH, false},
		{u.OtherLocalTestingPublicKeyOpenSSH, false},
		{u.LocalTestingPublicKeyOpenSSH + " admin@test", false},
		{"ssh-rsa AAAAC3NzaC1lZDI1NTE5AAAAIF70YNMkl2Wedmpo2UszvIrXJqr/pgCpevysNjjwtUig", true},
		{"ssh-ed25519 invalid", true},
		{u.LocalTestingPublicKeyOpenSSH + "\n", true},
		{"", true},
	}
	assertCases(t, FieldSshPublicKey, cases)
}

func TestValidateEmail(t *testing.T) {
	cases := []testCaseType{
		{"", true},
		{"admin@admin.com", false},
		{"@admin.com", true},
		{"admin@.com", true},
		{"admin@admin.", true},
		{"adminadmin.com", true},
		{"admin@admincom", true},
		{strings.Repeat("a", 64) + "@domain.com", false},
		{strings.Repeat("a", 65) + "@domain.com", true},
		{"abc@" + strings.Repeat("b", 253) + ".com", false},
		{"abc@" + strings.Repeat("b", 254) + ".com", true},
	}
	assertCases(t, FieldEmail, cases)
}

func TestValidateEmailOrEmpty(t *testing.T) {
	cases := []testCaseType{
		{"", false},
		{"admin@admin.com", false},
		{"@admin.com", true},
		{"admin@.com", true},
		{"admin@admin.", true},
		{"adminadmin.com", true},
		{"admin@admincom", true},
	}
	assertCases(t, FieldEmailOrEmpty, cases)
}

func TestValidateNumber(t *testing.T) {
	cases := []testCaseType{
		{"0", false},
		{"1", false},
		{"-1", true},
		{"a", true},
		{"A", true},
		{"z", true},
		{"Z", true},
		{"-", true},
		{"_", true},
		{".", true},
		{",", true},
		{"01234567890123456789", false}, // 20 digits
		{"012345678901234567890", true}, // 21 digits
	}
	assertCases(t, FieldNumber, cases)
}

func TestSearchTerm(t *testing.T) {
	cases := []testCaseType{
		{"", false},
		{"a", false},
		{"1", false},
		{"0123456789abcdefghij", false}, // length 20
		{"asdf!", true},
	}
	assertCases(t, FieldSearchTerm, cases)
}

func TestHost(t *testing.T) {
	cases := []testCaseType{
		{"localhost", false},
		{"localhost:8443", false},
		{"localhost123", false},
		{"example.com", false},
		{"my_example-website.com", false},
		{"a.", false},
		{"", false},
		{sixtyThreeHexDecimalLetters + "a", false},
		{sixtyThreeHexDecimalLetters + "ab", true},
	}
	assertCases(t, FieldHost, cases)
}

func TestRemoteHost(t *testing.T) {
	cases := []testCaseType{
		{"localhost", false},
		{"localhost:8443", true},
		{"localhost123", false},
		{"example.com", false},
		{"my_example-website.com", false},
		{"a.", false},
		{"", false},
		{sixtyThreeHexDecimalLetters + "a", false},
		{sixtyThreeHexDecimalLetters + "ab", true},
	}
	assertCases(t, FieldRemoteHost, cases)
}

func TestDomain(t *testing.T) {
	cases := []testCaseType{
		{"localhost", false},
		{"client.example.com", false},
		{"quollix.oidc-client.localhost", false},
		{"https://client.example.com", true},
		{"client.example.com:8443", true},
		{"client.example.com/callback", true},
		{"client.example.com?tenant=main", true},
		{"client.example.com#fragment", true},
		{"client example.com", true},
		{"admin@client.example.com", true},
	}
	assertCases(t, FieldDomain, cases)
}

func TestDomainPath(t *testing.T) {
	cases := []testCaseType{
		{"localhost", false},
		{"auth.example.com", false},
		{"auth.example.com/realms/main", false},
		{"auth.example.com/path_with_symbols-1.2", false},
		{"auth.example.com/path~with~tilde", true},
		{"https://auth.example.com/realms/main", true},
		{"auth.example.com:8443/realms/main", true},
		{"auth.example.com/realms/main?tenant=main", true},
		{"auth.example.com/realms/main#fragment", true},
		{"auth example.com/realms/main", true},
		{"admin@auth.example.com/realms/main", true},
	}
	assertCases(t, FieldDomainPath, cases)
}

var sampleKnownHosts = "# 127.0.0.1:2222 SSH-2.0-OpenSSH_9.9\n# 127.0.0.1:2222 SSH-2.0-OpenSSH_9.9\n[127.0.0.1]:2222 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQClESlJkZf90J0vZZNdAdvl4SpUDt+/VWpiMR8CYbGal8uu09a7UMP9hTeoPacrJtxXRooll7YWv8QRY+/c6UkZHaU4LCOwDJAATHVvKv1ynaGBzGbWK4sGSyTxuzyTYCzcqc1dO+te8qbHh6MI3mC5fF7U+jqU2pJDBfyHb80su4BmyAcSsRc1LgsrHBEYitfsblLWhwzhVRVvD4fRLasfcqpH7ein5peqJPiPOyBsl8+VEpMrH5AzeYsinD5RC84x+0yTOJEQMCdys+EC5i3/Pv3BJ2T/I9VyUoNfF3y9kcxoUIiSj7/kDDhtgAsC87Sv7n5WKrBzkpFpBurLZIaq+ucDUZunE7mbuntc7BI7FIdwxfZl8AgNGAeTAPsbCRORmdYzGNEbgbymMUeNmZYNcrykE8SAsGaaewM+5HnR6x7q7GSHarfIeVSWUDwhMcMCptrsIcSOZlJHEq4hDsb+cILLHQTeOmjuN7O6mLQw5zauIq39YpfzYj9u0PxLBiU=\n# 127.0.0.1:2222 SSH-2.0-OpenSSH_9.9\n[127.0.0.1]:2222 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBLO699LJQo4+GPThGkZ12YP10xfcf6Zn17nLKi85M1b4wBcb9iaBSLeRAMdszf41pWbW1BHlvXBUkfVbSaiqqh0=\n# 127.0.0.1:2222 SSH-2.0-OpenSSH_9.9\n[127.0.0.1]:2222 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIERZ7A/6JHp/4VSE3iKJGPWSV6SnYVfzGGamyHwYDsj4\n# 127.0.0.1:2222 SSH-2.0-OpenSSH_9.9"

func TestKnownHosts(t *testing.T) {
	cases := []testCaseType{
		{sampleKnownHosts, false},
		{sampleKnownHosts + "!", true},
	}
	assertCases(t, FieldKnownHosts, cases)
}

func TestResticBackupIdValidations(t *testing.T) {
	cases := []testCaseType{
		{"06b6458017d1e653195d696653c358e4e6a78772aed17582dd6539287332621f", false},
		{sixtyThreeHexDecimalLetters + "a", false},
		{"06b6458017d1e653195d696653c358e4e6a78772aed17582dd6539287332621fa", true},
		{sixtyThreeHexDecimalLetters + "g", true},
	}
	assertCases(t, FieldResticBackupID, cases)
}

func TestDefaultOrEmpty(t *testing.T) {
	cases := []testCaseType{
		{"", false},
		{"abc", false},
		{"ab", true},
		{"thisisaverylongdefaultvalue", true},
		{"validdefault", false},
		{"invalid_default", true},
	}
	assertCases(t, FieldDefaultOrEmpty, cases)
}

func TestFormatErrorWithContext(t *testing.T) {
	assert.Equal(t, "simple error message", formatErrorWithContext("simple error message", nil))
	assert.Equal(t, "simple error message (id: 42, name: john)", formatErrorWithContext("simple error message", map[string]any{
		"name": "john",
		"id":   42,
	}))
}

func TestIgnore(t *testing.T) {
	cases := []testCaseType{
		{"abcd123,.-#+$%&/()=?", false},
	}
	assertCases(t, FieldIgnore, cases)
}

func TestLooseAllowlist(t *testing.T) {
	cases := []testCaseType{
		{"a", false},
		{strings.Repeat("a", 128), false},
		{strings.Repeat("a", 129), true},
		{"", true},

		{"AZaz09", false},
		{"user.name", false},
		{"user_name", false},
		{"user-name", false},
		{"user+name", false},
		{"user=name", false},
		{"user@name", false},

		{"!@#$%^&*()_-+=.,:;/?[]{}|~<>", false},

		{"abc def", true},
		{"abc\tdef", true},
		{"abc\ndef", true},
		{"abc\rdef", true},

		{`abc"def`, true},
		{`abc'def`, true},
		{`abc\def`, true},
		{"abc`def", true},

		{"ümlaut", true},
	}

	assertCases(t, FieldLoose, cases)
}

func TestOidcSubject(t *testing.T) {
	cases := []testCaseType{
		{"user-123", false},
		{"google-oauth2|123456789", false},
		{"https://issuer.example.com/users/123", false},
		{"subject with spaces", false},
		{"ümlaut-subject", false},
		{strings.Repeat("a", 512), false},

		{"", true},
		{"   ", true},
		{"subject\nwith-newline", true},
		{"subject\twith-tab", true},
		{strings.Repeat("a", 513), true},
		{string([]byte{0xff}), true},
	}

	assertCases(t, FieldOidcSubject, cases)
}

func TestOidcClaim(t *testing.T) {
	cases := []testCaseType{
		{"", false},
		{"Jane Doe", false},
		{"jane.doe@example.com", false},
		{"名字", false},
		{strings.Repeat("a", 1024), false},

		{"Jane\nDoe", true},
		{"Jane\tDoe", true},
		{strings.Repeat("a", 1025), true},
		{string([]byte{0xff}), true},
	}

	assertCases(t, FieldOidcClaim, cases)
}

func TestCredential(t *testing.T) {
	cases := []testCaseType{
		{"client-id", false},
		{"secret with spaces and symbols !@#$%^&*()_+-=[]{}|;:',.<>/?", false},
		{"ümlaut-secret", false},
		{strings.Repeat("a", 1024), false},

		{"", true},
		{"   ", true},
		{"secret\nwith-newline", true},
		{"secret\twith-tab", true},
		{strings.Repeat("a", 1025), true},
		{string([]byte{0xff}), true},
	}

	assertCases(t, FieldCredential, cases)
}

func TestEmailSubject(t *testing.T) {
	cases := []testCaseType{
		{"Privacy policy updated", false},
		{"See https://quollix.org/legal/new-policy.md", false},
		{strings.Repeat("a", 120), false},

		{"", true},
		{"   ", true},
		{"Privacy\npolicy updated", true},
		{"Policy updated – please read", true},
		{strings.Repeat("a", 121), true},
	}

	assertCases(t, FieldEmailSubject, cases)
}

func TestEmailBody(t *testing.T) {
	cases := []testCaseType{
		{"Privacy policy updated.\nSee https://quollix.org/legal/new-policy.md", false},
		{"Line one\r\nLine two", false},
		{"Indented\tline", false},
		{strings.Repeat("a", 5000), false},

		{"", true},
		{"   ", true},
		{"Policy updated – please read", true},
		{strings.Repeat("a", 5001), true},
	}

	assertCases(t, FieldEmailBody, cases)
}
