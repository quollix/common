package browsertest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	u "github.com/quollix/common/utils"
)

func (p *Page) Close() error {
	p.cancel()
	return nil
}

func (p *Page) CancelTimeout() {
	p.cancel()
}

func (p *Page) Timeout(timeout time.Duration) *Page {
	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	return &Page{browser: p.browser, ctx: ctx, cancel: cancel}
}

func (p *Page) Navigate(url string) error {
	return p.runAndWaitReady(chromedp.Navigate(url))
}

func (p *Page) MustNavigate(url string) *Page {
	if err := p.Navigate(url); err != nil {
		panic(err)
	}
	return p
}

func (p *Page) MustReload() {
	if err := p.Reload(); err != nil {
		panic(err)
	}
}

func (p *Page) Reload() error {
	return p.runAndWaitReady(chromedp.Reload())
}

func (p *Page) WaitLoad() error {
	if err := chromedp.Do(p.ctx, chromedp.WaitReady(chromedp.CSS("body"))); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (p *Page) Info() (*PageInfo, error) {
	currentURL, err := chromedp.Run(p.ctx, chromedp.Location())
	if err != nil {
		return nil, u.Logger.NewError(err.Error())
	}
	return &PageInfo{URL: currentURL}, nil
}

func (p *Page) DoAndWaitLoad(action func() error) error {
	if err := action(); err != nil {
		return err
	}
	return p.WaitLoad()
}

func (p *Page) WaitOpen() func() (*Page, error) {
	targetIDChannel := chromedp.WaitNewTarget(p.ctx, func(info *target.Info) bool {
		return info.Type == "page"
	})
	return func() (*Page, error) {
		targetID := <-targetIDChannel
		ctx, cancel := chromedp.NewContext(p.browser.ctx, chromedp.WithTargetID(targetID))
		page := &Page{browser: p.browser, ctx: ctx, cancel: cancel}
		if err := chromedp.Do(ctx); err != nil {
			cancel()
			return nil, u.Logger.NewError(err.Error())
		}
		return page, nil
	}
}

func (p *Page) Element(selector string) (*Element, error) {
	return findElement(p, "document", selector)
}

func (p *Page) MustElement(selector string) *Element {
	element, err := p.Element(selector)
	if err != nil {
		panic(err)
	}
	return element
}

func (p *Page) ElementMatchingText(selector string, textPattern string) (*Element, error) {
	return findElementMatchingText(p, selector, textPattern)
}

func (p *Page) MustElementMatchingText(selector string, textPattern string) *Element {
	element, err := p.ElementMatchingText(selector, textPattern)
	if err != nil {
		panic(err)
	}
	return element
}

func (p *Page) Elements(selector string) ([]*Element, error) {
	return findElements(p, "document", selector)
}

func (p *Page) MustElements(selector string) []*Element {
	elements, err := p.Elements(selector)
	if err != nil {
		panic(err)
	}
	return elements
}

func (p *Page) Has(selector string) (bool, *Element, error) {
	return hasElement(p, "document", selector)
}

func (p *Page) SetCookie(cookie *http.Cookie, baseURL string) error {
	if _, err := chromedp.Call(p.ctx, network.SetCookie, network.SetCookieParams{
		Name:     cookie.Name,
		Value:    cookie.Value,
		URL:      baseURL,
		Path:     cookiePath(cookie),
		Secure:   new(cookie.Secure),
		HTTPOnly: new(cookie.HttpOnly),
	}); err != nil {
		return u.Logger.NewError(err.Error(), "url", baseURL)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if _, err := chromedp.Call(p.ctx, network.ClearBrowserCookies, cdp.Empty{}); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (p *Page) Cookies(url string) ([]*network.Cookie, error) {
	result, err := chromedp.Call(p.ctx, network.GetCookies, network.GetCookiesParams{URLs: []string{url}})
	if err != nil {
		return nil, u.Logger.NewError(err.Error(), "url", url)
	}
	return result.Cookies, nil
}

func (p *Page) count(expr string) (int, error) {
	count, err := chromedp.Run(p.ctx, chromedp.EvaluateAsDevTools[int](fmt.Sprintf(`(() => %s.length)()`, expr)))
	if err != nil {
		return 0, u.Logger.NewError(err.Error())
	}
	return count, nil
}

func (p *Page) runAndWaitReady(actions ...chromedp.Action[chromedp.Void]) error {
	if err := chromedp.Do(p.ctx, append(actions, chromedp.WaitReady(chromedp.CSS("body")))...); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}
