package browser

import (
	"net/http"
	"time"

	"github.com/quollix/common/browsertest"
	"github.com/quollix/common/quollix/api"
	"github.com/quollix/common/quollix/api_client"
	u "github.com/quollix/common/utils"
)

type Browser struct {
	BaseURL string
	Client  *api_client.QuollixClient
	Browser *browsertest.Browser
	Page    *browsertest.Page

	InstalledApps *InstalledAppsPageHelpers
}

func NewBrowser(baseURL string) (*Browser, error) {
	browser, err := browsertest.LaunchBrowser()
	if err != nil {
		return nil, u.Logger.NewError(err.Error())
	}

	page, err := browser.NewPage()
	if err != nil {
		_ = browser.Close()
		return nil, u.Logger.NewError(err.Error())
	}

	quollixBrowser := &Browser{
		BaseURL: baseURL,
		Client:  api_client.NewQuollixClientForRootUrl(baseURL),
		Browser: browser,
		Page:    page,
	}
	quollixBrowser.InstalledApps = &InstalledAppsPageHelpers{Page: page}
	return quollixBrowser, nil
}

func (b *Browser) Close() error {
	var closeErr error
	if b.Page != nil {
		closeErr = b.Page.Close()
	}
	if b.Browser != nil {
		if err := b.Browser.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	if closeErr != nil {
		return u.Logger.NewError(closeErr.Error())
	}
	return nil
}

func (b *Browser) Visit(path string) error {
	if err := b.Page.Navigate(b.BaseURL + path); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (b *Browser) SyncedLoginWithClient(username, password string) error {
	if err := b.Visit(api.Paths.FrontendSignIn); err != nil {
		return err
	}
	if err := b.Client.Auth.SignIn(username, password); err != nil {
		return err
	}
	if err := b.syncBrowserCookieFromClient(); err != nil {
		return err
	}
	return b.Visit(api.Paths.FrontendInstalledApps)
}

func (b *Browser) SyncedLoginWithBrowser(username, password string) error {
	if err := b.Visit(api.Paths.FrontendSignIn); err != nil {
		return err
	}
	if err := b.Page.MustElement("#sign-in-tab").Click(); err != nil {
		return u.Logger.NewError(err.Error())
	}
	if err := b.Page.MustElement("#username-input").Input(username); err != nil {
		return u.Logger.NewError(err.Error())
	}
	if err := b.Page.MustElement("#password-input").Input(password); err != nil {
		return u.Logger.NewError(err.Error())
	}
	if err := b.Page.DoAndWaitLoad(func() {
		b.Page.MustElement("#sign-in-button").MustClick()
	}); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return b.waitUntilClientCookieSyncedFromBrowser()
}

func (b *Browser) UseAuthCookie(cookie *http.Cookie) error {
	b.Client.Parent.Cookie = cookie
	return b.syncBrowserCookieFromClient()
}

func (b *Browser) ClearSession() error {
	b.Client.Parent.Cookie = nil
	return b.clearBrowserCookies()
}

func (b *Browser) syncBrowserCookieFromClient() error {
	if err := b.clearBrowserCookies(); err != nil {
		return err
	}
	if b.Client.Parent.Cookie == nil {
		return nil
	}
	if err := b.Page.SetCookie(b.Client.Parent.Cookie, b.BaseURL); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (b *Browser) clearBrowserCookies() error {
	if err := b.Page.ClearCookies(); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (b *Browser) syncClientCookieFromBrowser() error {
	cookie, err := b.getAuthCookieFromBrowser()
	if err != nil {
		return err
	}
	b.Client.Parent.Cookie = cookie
	return nil
}

func (b *Browser) waitUntilClientCookieSyncedFromBrowser() error {
	var lastErr error
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if err := b.syncClientCookieFromBrowser(); err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return nil
	}
	return lastErr
}

func (b *Browser) getAuthCookieFromBrowser() (*http.Cookie, error) {
	cookies, err := b.Page.Cookies(b.BaseURL)
	if err != nil {
		return nil, u.Logger.NewError(err.Error())
	}
	for _, cookie := range cookies {
		if cookie.Name != api.BrandAppAuthCookieName {
			continue
		}
		return &http.Cookie{ // #nosec G124: tests reconstruct the browser cookie only for local test replay
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
		}, nil
	}
	return nil, u.Logger.NewError("no auth cookie found")
}

func (b *Browser) OpenInstalledAppsPage() *InstalledAppsPageHelpers {
	if err := b.Visit(api.Paths.FrontendInstalledApps); err != nil {
		panic(err)
	}
	return b.InstalledApps
}
