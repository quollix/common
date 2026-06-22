package utils

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/quollix/deepstack"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DirectoryField        = "directory"
	SizeInBytesField      = "size_in_bytes"
	FolderToFindField     = "folder_to_find"
	InitialDirField       = "initial_directory"
	HostField             = "host"
	CurrentAttemptField   = "current_attempt"
	MaximumAttemptsFields = "maximum_attempts"
	PortField             = "port"
)

var (
	Logger               = deepstack.NewDeepStackLogger(deepstack.NewRawConsoleHandler(slog.LevelInfo))
	OperationFailedError = "Operation failed"
)

func NewFileHandler(level slog.Level) *slog.JSONHandler {
	opts := &slog.HandlerOptions{AddSource: true, Level: level}
	logDir := "data/logs"
	_ = os.MkdirAll(logDir, 0700)
	logFile := &lumberjack.Logger{
		Filename: logDir + "/app.log",
		MaxSize:  100,
		MaxAge:   30,
		Compress: true,
	}
	return slog.NewJSONHandler(logFile, opts)
}

func WriteResponseErrorAlways(w http.ResponseWriter, err error, context ...any) {
	writeResponseErrorInternal(w, nil, true, err, context...)
}

// WriteResponseError Here we can define a set of expected error messages that can safely be exposed to clients. If the error matches one of these strings, it is returned unchanged. Otherwise, a default message is used to ensure that no sensitive error messages are returned.
func WriteResponseError(w http.ResponseWriter, expectedErrors map[string]any, err error, context ...any) {
	writeResponseErrorInternal(w, expectedErrors, false, err, context...)
}

func writeResponseErrorInternal(
	w http.ResponseWriter,
	expectedErrors map[string]any,
	alwaysWriteActualError bool,
	err error,
	context ...any,
) {
	var actualErrorMessage string
	deepStackError, isDeepStackError := err.(*deepstack.DeepStackError)
	if isDeepStackError {
		actualErrorMessage = deepStackError.Message
	} else {
		Logger.Warn("expected a DeepStack error, but got a regular error")
		actualErrorMessage = err.Error()
	}

	isExpectedError := false
	if expectedErrors != nil {
		_, isExpectedError = expectedErrors[actualErrorMessage]
	}
	if !isExpectedError && strings.HasPrefix(strings.ToLower(actualErrorMessage), "invalid input") {
		isExpectedError = true
	}

	shouldWriteActualError := alwaysWriteActualError || isExpectedError

	if shouldWriteActualError {
		if isDeepStackError {
			Logger.Info(deepStackError, context...)
		} else {
			Logger.Info(actualErrorMessage, context...)
		}
		http.Error(w, actualErrorMessage, http.StatusBadRequest)
		return
	}

	if isDeepStackError {
		Logger.Error(deepStackError, context...)
	} else {
		context = append(context, "stack_trace", string(debug.Stack()))
		Logger.Error(actualErrorMessage, context...)
	}
	http.Error(w, OperationFailedError, http.StatusBadRequest)
}

func MapOf(elements ...string) map[string]any {
	var newMap = map[string]any{}
	for _, element := range elements {
		newMap[element] = nil
	}
	return newMap
}

func ExtractError(err error) string {
	if err == nil {
		return ""
	}

	deepStackError, ok := err.(*deepstack.DeepStackError)
	if ok {
		return deepStackError.Message
	} else {
		return err.Error()
	}
}
