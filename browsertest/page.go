package browsertest

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	if err := chromedp.Run(p.ctx, chromedp.WaitReady("body", chromedp.ByQuery)); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (p *Page) Info() (*PageInfo, error) {
	var currentURL string
	if err := chromedp.Run(p.ctx, chromedp.Location(&currentURL)); err != nil {
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
		if err := chromedp.Run(ctx); err != nil {
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
	if err := chromedp.Run(p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetCookie(cookie.Name, cookie.Value).
			WithURL(baseURL).
			WithPath(cookiePath(cookie)).
			WithSecure(cookie.Secure).
			WithHTTPOnly(cookie.HttpOnly).
			Do(ctx)
	})); err != nil {
		return u.Logger.NewError(err.Error(), "url", baseURL)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := chromedp.Run(p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.ClearBrowserCookies().Do(ctx)
	})); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (p *Page) Cookies(url string) ([]*network.Cookie, error) {
	var cookies []*network.Cookie
	err := chromedp.Run(p.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().WithURLs([]string{url}).Do(ctx)
		return err
	}))
	if err != nil {
		return nil, u.Logger.NewError(err.Error(), "url", url)
	}
	return cookies, nil
}

func (p *Page) count(expr string) (int, error) {
	var count int
	err := chromedp.Run(p.ctx, chromedp.EvaluateAsDevTools(fmt.Sprintf(`(() => %s.length)()`, expr), &count))
	if err != nil {
		return 0, u.Logger.NewError(err.Error())
	}
	return count, nil
}

func (p *Page) runAndWaitReady(actions ...chromedp.Action) error {
	if err := chromedp.Run(p.ctx, append(actions, chromedp.WaitReady("body", chromedp.ByQuery))...); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}
