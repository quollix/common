package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"

	u "github.com/quollix/common/utils"
)

type QuollixCertificatesClient struct {
	quollix *QuollixClient
}

func (c *QuollixCertificatesClient) Reset() error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsCertificateReset, nil)
	return err
}

func (c *QuollixCertificatesClient) TryDns01Challenge() (*api.Dns01ChallengeInfo, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsStartDns01CertificateChallenge, nil)
	if err != nil {
		return nil, err
	}
	var info api.Dns01ChallengeInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *QuollixCertificatesClient) Upload(content []byte) error {
	_, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsCertificateUpload, api.BinaryFile{
		FileName: "certificate.pem",
		Content:  content,
	})
	return err
}

func (c *QuollixCertificatesClient) DownloadBundleBytes() ([]byte, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendSettingsCertificateDownload, nil)
	if err != nil {
		return nil, err
	}

	var certificateFile api.BinaryFile
	err = json.Unmarshal(body, &certificateFile)
	if err != nil {
		return nil, err
	}
	if certificateFile.FileName != "certificate.pem" {
		return nil, u.Logger.NewError("unexpected certificate file name", "file_name", certificateFile.FileName)
	}
	return certificateFile.Content, nil
}
