package browsertest

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
)

func LaunchBrowser() (*Browser, error) {
	headless := os.Getenv("HEADFUL") != "true"
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.IgnoreCertErrors,
	)
	if os.Getenv("CI") == "true" {
		options = append(options, chromedp.NoSandbox)
	}

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		allocatorCancel()
		return nil, err
	}

	return &Browser{
		ctx: ctx,
		cancel: func() {
			cancel()
			allocatorCancel()
		},
	}, nil
}

func MustLaunchBrowser() *Browser {
	browser, err := LaunchBrowser()
	if err != nil {
		panic(err)
	}
	return browser
}

func (b *Browser) NewIncognito() (*Browser, error) {
	ctx, cancel := chromedp.NewContext(b.ctx, chromedp.WithNewBrowserContext())
	return &Browser{ctx: ctx, cancel: cancel}, nil
}

func (b *Browser) NewPage() (*Page, error) {
	ctx, cancel := chromedp.NewContext(b.ctx)
	page := &Page{browser: b, ctx: ctx, cancel: cancel}
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, err
	}
	return page, nil
}

func (b *Browser) MustPage() *Page {
	page, err := b.NewPage()
	if err != nil {
		panic(err)
	}
	return page
}

func (b *Browser) Close() error {
	b.cancel()
	return nil
}
