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

func (c *QuollixCertificatesClient) UploadCertificateBundle(content []byte) error {
	return c.quollix.uploadBinaryFile(api.Paths.BackendSettingsCertificateUpload, api.BinaryFile{
		FileName: "certificate.pem",
		Content:  content,
	})
}

func (c *QuollixCertificatesClient) DownloadCertificateBundleBytes() ([]byte, error) {
	certificateFile, err := c.quollix.downloadBinaryFile(api.Paths.BackendSettingsCertificateDownload, nil)
	if err != nil {
		return nil, err
	}

	if certificateFile.FileName != "certificate.pem" {
		return nil, u.Logger.NewError("unexpected certificate file name", "file_name", certificateFile.FileName)
	}
	return certificateFile.Content, nil
}

func (c *QuollixCertificatesClient) UploadAcmeAccountPrivateKey(content []byte) error {
	return c.quollix.uploadBinaryFile(api.Paths.BackendSettingsAcmeAccountPrivateKeyUpload, api.BinaryFile{
		FileName: "acme_account_private_key.pem",
		Content:  content,
	})
}

func (c *QuollixCertificatesClient) DownloadAcmeAccountPrivateKey() (*api.BinaryFile, error) {
	return c.quollix.downloadBinaryFile(api.Paths.BackendSettingsAcmeAccountPrivateKeyDownload, nil)
}
