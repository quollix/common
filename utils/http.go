package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type ComponentClient struct {
	Cookie            *http.Cookie
	SetCookieHeader   bool
	RootUrl           string
	Origin            string
	VerifyCertificate bool
}

func SendJsonResponse(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		Logger.Error(err)
		http.Error(w, "marshalling failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		Logger.Error(err.Error())
		return
	}
}

func (c *ComponentClient) DoRequest(path string, payload any) ([]byte, error) {
	resp, err := c.DoRequestWithFullResponse(path, payload)
	if err != nil {
		return nil, err
	}

	respBody, err := assertOkStatusAndExtractBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *ComponentClient) DoRequestWithFullResponse(path string, payload any) (*http.Response, error) {
	url := c.RootUrl + path

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	payloadReader := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", url, payloadReader)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	SetCookieHeaders(req, c)
	if c.Origin != "" {
		req.Header.Set("Origin", c.Origin)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.VerifyCertificate}, // #nosec G402 (CWE-295): TLS InsecureSkipVerify may be true; tolerated by design
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}

	respBody, err := assertOkStatusAndExtractBody(resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Cookies()) == 1 {
		c.Cookie = resp.Cookies()[0]
	}

	// Response body can only be read once. When reading it a second time, an error occurs. So a copy is created.
	newResp := &http.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       io.NopCloser(bytes.NewBuffer(respBody)),
	}
	return newResp, nil
}

func SetCookieHeaders(req *http.Request, c *ComponentClient) {
	if c.SetCookieHeader && c.Cookie != nil {
		req.AddCookie(c.Cookie)
	}
}

func assertOkStatusAndExtractBody(resp *http.Response) ([]byte, error) {
	defer Close(resp.Body)

	var didErrorOccur = false
	if resp.StatusCode != http.StatusOK {
		didErrorOccur = true
	}

	var potentialErrorContext []any
	respBody, err := io.ReadAll(resp.Body)
	if err == nil {
		responseBodyString := string(respBody)
		trimmed := strings.TrimSuffix(responseBodyString, "\n")
		potentialErrorContext = append(potentialErrorContext, "response_body", trimmed)
	} else {
		didErrorOccur = true
		potentialErrorContext = append(potentialErrorContext, "response_body_reading_error", err.Error())
	}

	if didErrorOccur {
		potentialErrorContext = append(potentialErrorContext, "status_code", resp.StatusCode)
		return nil, Logger.NewError("request failed", potentialErrorContext...)
	}

	return respBody, nil
}

func UnpackResponse[T any](object any) (*T, error) {
	respBody, ok := object.([]byte)
	if !ok {
		return nil, Logger.NewError("failed to cast to type []byte")
	}

	var result T
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, Logger.NewError(err.Error())
	}
	return &result, nil
}
