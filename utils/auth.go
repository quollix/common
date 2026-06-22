package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHelper interface {
	GenerateCookie() (*http.Cookie, error)
	GenerateSecret() (string, error)
	GetTimeInSevenDays() time.Time
	SaltAndHash(clearText string) (string, error)
	DoesMatchSaltedHash(clearText, saltedHash string) bool
	GetSHA256Hash(plainText string) string
}

type AuthHelperImpl struct{}

func (a *AuthHelperImpl) GenerateCookie() (*http.Cookie, error) {
	secret, err := a.GenerateSecret()
	if err != nil {
		return nil, err
	}
		return &http.Cookie{ // #nosec G124: this shared helper supports localhost HTTP flows; callers that serve browser traffic harden flags before sending
		Name:     "auth",
		Value:    secret,
		Expires:  a.GetTimeInSevenDays(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}, nil
}

func (a *AuthHelperImpl) GenerateSecret() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", Logger.NewError(err.Error())
	}
	return hex.EncodeToString(randomBytes), nil
}

func (a *AuthHelperImpl) GetTimeInSevenDays() time.Time {
	// We need nanosecond precision since some component tests update the cookie expiration time through a request. Since the test takes less than a second, the update is not visible when the precision is set to seconds, which causes the test to fail.
	return time.Now().UTC().Add(7 * 24 * time.Hour)
}

func (a *AuthHelperImpl) SaltAndHash(clearText string) (string, error) {
	hashValue, err := bcrypt.GenerateFromPassword([]byte(clearText), bcrypt.DefaultCost)
	if err != nil {
		return "", Logger.NewError(err.Error())
	}
	return string(hashValue), nil
}

func (a *AuthHelperImpl) DoesMatchSaltedHash(clearText, saltedHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(saltedHash), []byte(clearText)) == nil
}

func (a *AuthHelperImpl) GetSHA256Hash(plainText string) string {
	h := sha256.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}
