package api_client

import (
	"encoding/json"

	"github.com/quollix/common/quollix/api"
)

func (c *QuollixClient) uploadBinaryFile(path string, file api.BinaryFile) error {
	_, err := c.Parent.DoRequest(path, file)
	return err
}

func (c *QuollixClient) downloadBinaryFile(path string, payload any) (*api.BinaryFile, error) {
	body, err := c.Parent.DoRequest(path, payload)
	if err != nil {
		return nil, err
	}

	var file api.BinaryFile
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}
