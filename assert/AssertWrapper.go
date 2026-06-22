package assert

import (
	"log/slog"
	"testing"

	"github.com/quollix/deepstack"
	"github.com/stretchr/testify/require"
)

var logger = deepstack.NewDeepStackLogger(deepstack.NewRawConsoleHandler(slog.LevelInfo))

func True(t *testing.T, condition bool) {
	require.True(t, condition)
}

func False(t *testing.T, condition bool) {
	require.False(t, condition)
}

func Equal(t *testing.T, expected any, actual any) {
	require.Equal(t, expected, actual)
}

func Nil(t *testing.T, object any) {
	err, ok := object.(*deepstack.DeepStackError)
	if ok {
		logger.Error(err)
		require.Fail(t, "expected nil, but got a deepstack error instead, see error log")
		return
	}
	require.Nil(t, object)
}

func NotNil(t *testing.T, object any) {
	require.NotNil(t, object)
}

func NotEqual(t *testing.T, expected any, actual any) {
	require.NotEqual(t, expected, actual)
}
