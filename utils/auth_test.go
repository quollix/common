package utils

import (
	"net/http"
	"testing"
	"time"

	"github.com/quollix/common/assert"
)

var authHelper = &AuthHelperImpl{}

func TestHash(t *testing.T) {
	hashedString := authHelper.GetSHA256Hash("hello")
	assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", hashedString)
}

func TestSaltAndHash(t *testing.T) {
	sampleString := "hello"
	saltedAndHashedString, err := authHelper.SaltAndHash(sampleString)
	assert.Nil(t, err)

	saltedAndHashedString2, err := authHelper.SaltAndHash(sampleString)
	assert.Nil(t, err)

	assert.NotEqual(t, saltedAndHashedString, saltedAndHashedString2)

	assert.True(t, authHelper.DoesMatchSaltedHash(sampleString, saltedAndHashedString))
	assert.True(t, authHelper.DoesMatchSaltedHash(sampleString, saltedAndHashedString2))
	assert.False(t, authHelper.DoesMatchSaltedHash(sampleString+"x", saltedAndHashedString))
	assert.False(t, authHelper.DoesMatchSaltedHash(sampleString+"x", saltedAndHashedString2))
}

func TestGenerateCookie(t *testing.T) {
	cookie, err := authHelper.GenerateCookie()
	assert.Nil(t, err)
	assert.NotNil(t, cookie)
	assert.Equal(t, "auth", cookie.Name)
	assert.True(t, len(cookie.Value) > 0)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.True(t, cookie.Expires.After(time.Now()))
	assert.True(t, cookie.Expires.Before(time.Now().Add(31*24*time.Hour)))
	assert.True(t, cookie.HttpOnly)
}
