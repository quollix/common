package utils

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/quollix/common/assert"
	"github.com/quollix/deepstack"
)

const (
	BrandName               = "quollix"
	OfficialMaintainer      = BrandName
	OfficialBrandAppName    = BrandName
	OfficialDatabaseAppName = "postgres"
)

type Closable interface {
	Close() error
}

func Close(r Closable) {
	RunAndLogIfError(r.Close)
}

func RunAndLogIfError(fn func() error) {
	if err := fn(); err != nil {
		deepstackErr := Logger.NewError(err.Error())
		Logger.Error(deepstackErr)
	}
}

func AssertDeepStackErrorFromRequest(t *testing.T, err error, expectedResponseBodyErrorMessage string) {
	deepstack.AssertDeepStackError(t, err, "request failed", "response_body", expectedResponseBodyErrorMessage, "status_code", http.StatusBadRequest)
}

func AssertInvalidInputError(t *testing.T, err error) {
	deepStackError, ok := err.(*deepstack.DeepStackError)
	assert.True(t, ok)
	assert.Equal(t, "request failed", deepStackError.Message)
	assert.Equal(t, deepStackError.Context["status_code"], http.StatusBadRequest)
	responseBody, ok := deepStackError.Context["response_body"]
	assert.True(t, ok)
	responseBodyString, ok := responseBody.(string)
	assert.True(t, ok)
	assert.True(t, strings.HasPrefix(strings.ToLower(responseBodyString), "invalid input"))
}

var reservedAppNames = map[string]any{
	OfficialBrandAppName:    nil,
	"store":                 nil,
	"quollog":               nil,
	OfficialDatabaseAppName: nil,
}

func IsSystemApp(appName string) bool {
	_, exists := reservedAppNames[appName]
	return exists
}

func WaitMillis(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func Eventually(check func() error) error {
	return EventuallyWithTimeout(3*time.Second, 50*time.Millisecond, check)
}

func EventuallyWithTimeout(timeout time.Duration, interval time.Duration, check func() error) error {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		err := check()
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(interval)
	}

	if lastErr == nil {
		return nil
	}

	return lastErr
}
