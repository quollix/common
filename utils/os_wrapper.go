package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type OsWrapper interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Remove(path string) error
	RemoveAll(path string) error
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(path string) ([]os.DirEntry, error)
	GetTempDir() (string, error)
	DoesFileExist(path string) (bool, error)
	AllocateLocalhostPort() (string, error)
	Now() time.Time
	PromptUser(prompt string) (string, error)
	Sleep(duration time.Duration)
}

type OsWrapperImpl struct{}

func (o *OsWrapperImpl) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path) // #nosec G304 (CWE-22): Potential file inclusion via variable; no problem since it is called internally
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, Logger.NewError(err.Error(), "path", path)
	}
	return data, nil
}

func (o *OsWrapperImpl) WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return Logger.NewError(err.Error(), "path", path)
	}
	return nil
}

func (o *OsWrapperImpl) Remove(path string) error {
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return Logger.NewError(err.Error(), "path", path)
	}
	return nil
}

func (o *OsWrapperImpl) RemoveAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return Logger.NewError(err.Error(), "path", path)
	}
	return nil
}

func (o *OsWrapperImpl) MkdirAll(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return Logger.NewError(err.Error(), "path", path)
	}
	return nil
}

func (o *OsWrapperImpl) ReadDir(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, Logger.NewError(err.Error(), "path", path)
	}
	return entries, nil
}

func (o *OsWrapperImpl) GetTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "file_system_operator")
	if err != nil {
		return "", Logger.NewError(err.Error())
	}
	return tempDir, nil
}

func (o *OsWrapperImpl) DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, Logger.NewError(err.Error(), "path", path)
	}
	return true, nil
}

func (o *OsWrapperImpl) AllocateLocalhostPort() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", Logger.NewError(err.Error())
	}
	defer Close(listener)
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return "", Logger.NewError("failed to determine allocated localhost port")
	}
	return strconv.Itoa(addr.Port), nil
}

func (o *OsWrapperImpl) Now() time.Time {
	return time.Now().UTC()
}

func (o *OsWrapperImpl) PromptUser(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
		return "", Logger.NewError(err.Error())
	}

	return strings.TrimSpace(value), nil
}

func (o *OsWrapperImpl) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
