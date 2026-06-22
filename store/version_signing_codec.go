package store

import (
	"bytes"
	"encoding/binary"
	"time"

	u "github.com/quollix/common/utils"
)

type VersionSigningCodec interface {
	EncodeVersion(version *Version) ([]byte, error)
}

type VersionSigningCodecImpl struct{}

func (c *VersionSigningCodecImpl) EncodeVersion(version *Version) ([]byte, error) {
	if err := validateVersionSigningVersion(version); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := writeVersionSigningField(buf, []byte(version.Maintainer)); err != nil {
		return nil, err
	}
	if err := writeVersionSigningField(buf, []byte(version.AppName)); err != nil {
		return nil, err
	}
	if err := writeVersionSigningField(buf, []byte(version.VersionName)); err != nil {
		return nil, err
	}
	if err := writeVersionSigningField(buf, version.Content); err != nil {
		return nil, err
	}
	creationTimestamp := version.VersionCreationTimestamp.UTC().Format(time.RFC3339)
	if err := writeVersionSigningField(buf, []byte(creationTimestamp)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func validateVersionSigningVersion(version *Version) error {
	if version == nil {
		return u.Logger.NewError("version must not be nil")
	}
	if version.Maintainer == "" {
		return u.Logger.NewError("version maintainer must not be empty")
	}
	if version.AppName == "" {
		return u.Logger.NewError("version app name must not be empty")
	}
	if version.VersionName == "" {
		return u.Logger.NewError("version name must not be empty")
	}
	if len(version.Content) == 0 {
		return u.Logger.NewError("version content must not be empty")
	}
	if version.VersionCreationTimestamp.IsZero() {
		return u.Logger.NewError("version creation timestamp must not be zero")
	}
	return nil
}

func writeVersionSigningField(buf *bytes.Buffer, field []byte) error {
	if err := binary.Write(buf, binary.BigEndian, uint64(len(field))); err != nil {
		return u.Logger.NewError(err.Error())
	}
	_, err := buf.Write(field)
	if err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}
