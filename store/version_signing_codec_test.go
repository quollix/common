package store

import (
	"testing"
	"time"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

var sampleVersionSigningTimestamp = time.Date(2026, 3, 26, 10, 11, 12, 123456789, time.UTC)

func TestVersionSigningCodec_EncodeVersion(t *testing.T) {
	codec := &VersionSigningCodecImpl{}
	originalVersion := getSampleVersion()

	originalBytes, err := codec.EncodeVersion(originalVersion)
	assert.Nil(t, err)

	sameInputBytes, err := codec.EncodeVersion(getSampleVersion())
	assert.Nil(t, err)
	assert.Equal(t, originalBytes, sameInputBytes)

	testCases := []struct {
		name   string
		adjust func(version *Version)
	}{
		{
			name: "maintainer change changes output",
			adjust: func(version *Version) {
				version.Maintainer = "other-maintainer"
			},
		},
		{
			name: "app name change changes output",
			adjust: func(version *Version) {
				version.AppName = "other-app"
			},
		},
		{
			name: "version name change changes output",
			adjust: func(version *Version) {
				version.VersionName = "v1.2.4"
			},
		},
		{
			name: "content change changes output",
			adjust: func(version *Version) {
				version.Content = []byte("updated compose content")
			},
		},
		{
			name: "creation timestamp change changes output",
			adjust: func(version *Version) {
				version.VersionCreationTimestamp = version.VersionCreationTimestamp.Add(time.Second)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version := getSampleVersion()
			tc.adjust(version)

			actualBytes, err := codec.EncodeVersion(version)
			assert.Nil(t, err)
			assert.NotEqual(t, originalBytes, actualBytes)
		})
	}
}

func TestVersionSigningCodec_EncodeVersion_EmptyFieldReturnsError(t *testing.T) {
	codec := &VersionSigningCodecImpl{}

	testCases := []struct {
		name          string
		adjust        func(version *Version)
		expectedError string
	}{
		{
			name: "empty maintainer returns error",
			adjust: func(version *Version) {
				version.Maintainer = ""
			},
			expectedError: "version maintainer must not be empty",
		},
		{
			name: "empty app name returns error",
			adjust: func(version *Version) {
				version.AppName = ""
			},
			expectedError: "version app name must not be empty",
		},
		{
			name: "empty version name returns error",
			adjust: func(version *Version) {
				version.VersionName = ""
			},
			expectedError: "version name must not be empty",
		},
		{
			name: "empty content returns error",
			adjust: func(version *Version) {
				version.Content = []byte{}
			},
			expectedError: "version content must not be empty",
		},
		{
			name: "zero creation timestamp returns error",
			adjust: func(version *Version) {
				version.VersionCreationTimestamp = time.Time{}
			},
			expectedError: "version creation timestamp must not be zero",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version := getSampleVersion()
			tc.adjust(version)

			_, err := codec.EncodeVersion(version)

			assert.Equal(t, tc.expectedError, u.ExtractError(err))
		})
	}
}

func getSampleVersion() *Version {
	return &Version{
		Maintainer:               "sample-maintainer",
		AppName:                  "sample-app",
		VersionName:              "v1.2.3",
		Content:                  []byte("sample compose content"),
		VersionCreationTimestamp: sampleVersionSigningTimestamp,
	}
}
