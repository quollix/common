package api_client

import (
	"encoding/json"
	api "github.com/quollix/common/quollix/api"
	"io"
	"net/http"
	"net/url"
	"strings"

	u "github.com/quollix/common/utils"
)

type AppAccessClient struct {
	quollix *QuollixClient
}

var noRedirectHttpClient = &http.Client{
	Transport:     httpClient.Transport,
	CheckRedirect: func(request *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
}

func (c *AppAccessClient) GetSecret(appName string) (string, error) {
	body, err := c.quollix.Parent.DoRequest(api.Paths.BackendSecret, api.AppAccessSecretRequest{
		AppName: appName,
	})
	if err != nil {
		return "", err
	}

	var secret string
	err = json.Unmarshal(body, &secret)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func (c *AppAccessClient) DoAppRequest(appAccessCookie *http.Cookie, httpMethod string, appUrl string, httpBody io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, appUrl, httpBody)
	if err != nil {
		return nil, err
	}
	if appAccessCookie != nil {
		req.AddCookie(appAccessCookie)
	}
	return httpClient.Do(req)
}

func (c *AppAccessClient) ExchangeSecretForAppAccessCookie(secret string, appUrl string) (*http.Cookie, error) {
	urlWithSecret, err := getAppUrlWithSecret(appUrl, secret)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urlWithSecret, nil)
	if err != nil {
		return nil, err
	}
	if c.quollix.Parent.Cookie != nil {
		req.AddCookie(c.quollix.Parent.Cookie)
	}
	resp, err := noRedirectHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer u.Close(resp.Body)
	if resp.StatusCode != http.StatusFound {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, u.Logger.NewError("request failed", "status_code", resp.StatusCode, "response_body", strings.TrimSpace(string(body)))
	}
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		return nil, u.Logger.NewError("expected app cookie", "cookie_count", len(cookies))
	}
	return cookies[0], nil
}

func getAppUrlWithSecret(appUrl string, secret string) (string, error) {
	parsedUrl, err := url.Parse(appUrl)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	query.Set(api.BrandAppQuerySecretName, secret)
	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String(), nil
}
