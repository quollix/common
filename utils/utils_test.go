package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/quollix/common/assert"
)

func TestFindDir(t *testing.T) {
	dir, err := FindDir("utils")
	assert.Nil(t, err)
	assert.Equal(t, "utils", filepath.Base(dir))
}

func TestExtractFromLine(t *testing.T) {
	assert.Equal(t, 0, len(extractTagsFromLine("// random comment")))

	tags := extractTagsFromLine("//go:build a && !b || (c && d.e)")
	assert.Equal(t, 4, len(tags))
	assert.Equal(t, "a", tags[0])
	assert.Equal(t, "b", tags[1])
	assert.Equal(t, "c", tags[2])
	assert.Equal(t, "d.e", tags[3])

	tags = extractTagsFromLine("// +build a,!b c")
	assert.Equal(t, 3, len(tags))
	assert.Equal(t, "a", tags[0])
	assert.Equal(t, "b", tags[1])
	assert.Equal(t, "c", tags[2])
}

func TestIsSystemApp(t *testing.T) {
	assert.True(t, IsSystemApp("quollix"))
	assert.False(t, IsSystemApp("someotherapp"))
}

func TestFormatRelativeDuration(t *testing.T) {
	now := time.Date(2026, time.April, 24, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		target   time.Time
		expected string
	}{
		{name: "zero", target: now, expected: "0s"},
		{name: "future seconds", target: now.Add(1 * time.Second), expected: "in 1s"},
		{name: "future days", target: now.Add(6*24*time.Hour + 3*time.Hour), expected: "in 6d"},
		{name: "sub second", target: now.Add(-600 * time.Millisecond), expected: "0s"},
		{name: "seconds", target: now.Add(-(45*time.Second + 900*time.Millisecond)), expected: "45s ago"},
		{name: "minutes", target: now.Add(-(2*time.Minute + 5*time.Second)), expected: "2m ago"},
		{name: "hours", target: now.Add(-(3*time.Hour + 4*time.Minute)), expected: "3h ago"},
		{name: "days", target: now.Add(-(6*24*time.Hour + 3*time.Hour)), expected: "6d ago"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, FormatRelativeDuration(now, testCase.target))
		})
	}
}

func TestPromptUserReadsOneLinePerPrompt(t *testing.T) {
	oldStdin := os.Stdin
	readFile, writeFile, err := os.Pipe()
	assert.Nil(t, err)
	defer func() {
		os.Stdin = oldStdin
		Close(readFile)
		Close(writeFile)
	}()
	os.Stdin = readFile
	_, err = writeFile.WriteString("1\nDELETE\n")
	assert.Nil(t, err)
	assert.Nil(t, writeFile.Close())

	osWrapper := &OsWrapperImpl{}
	firstValue, err := osWrapper.PromptUser("first: ")
	assert.Nil(t, err)
	secondValue, err := osWrapper.PromptUser("second: ")
	assert.Nil(t, err)

	assert.Equal(t, "1", firstValue)
	assert.Equal(t, "DELETE", secondValue)
}

func TestPromptSecretReadsFromNonTerminalStdin(t *testing.T) {
	oldStdin := os.Stdin
	readFile, writeFile, err := os.Pipe()
	assert.Nil(t, err)
	defer func() {
		os.Stdin = oldStdin
		Close(readFile)
		Close(writeFile)
	}()
	os.Stdin = readFile
	_, err = writeFile.WriteString("secret\n")
	assert.Nil(t, err)
	assert.Nil(t, writeFile.Close())

	osWrapper := &OsWrapperImpl{}
	value, err := osWrapper.PromptSecret("secret: ")

	assert.Nil(t, err)
	assert.Equal(t, "secret", value)
}

func TestPromptSecretFailsForEmptyNonTerminalStdin(t *testing.T) {
	oldStdin := os.Stdin
	readFile, writeFile, err := os.Pipe()
	assert.Nil(t, err)
	defer func() {
		os.Stdin = oldStdin
		Close(readFile)
		Close(writeFile)
	}()
	os.Stdin = readFile
	_, err = writeFile.WriteString("\n")
	assert.Nil(t, err)
	assert.Nil(t, writeFile.Close())

	osWrapper := &OsWrapperImpl{}
	value, err := osWrapper.PromptSecret("secret: ")

	assert.Equal(t, "", value)
	assert.NotNil(t, err)
}
