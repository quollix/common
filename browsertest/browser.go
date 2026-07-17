package browsertest

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/chromedp/chromedp"
)

func LaunchBrowser() (*Browser, error) {
	browserPath, err := findChromiumPath()
	if err != nil {
		return nil, err
	}

	headless := os.Getenv("HEADFUL") != "true"
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", headless),
		chromedp.IgnoreCertErrors,
	)
	if os.Getenv("CI") == "true" {
		options = append(options, chromedp.NoSandbox)
	}

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithErrorf(chromedpErrorf))
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

func chromedpErrorf(format string, args ...any) {
	if format == "unhandled node event %T" && len(args) == 1 && fmt.Sprintf("%T", args[0]) == "*dom.EventTopLayerElementsUpdated" {
		return
	}
	log.Printf(format, args...)
}

func MustLaunchBrowser() *Browser {
	browser, err := LaunchBrowser()
	if err != nil {
		panic(err)
	}
	return browser
}

func findChromiumPath() (string, error) {
	for _, name := range []string{"chromium", "chromium-browser"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("chromium executable not found in PATH, install chromium or chromium-browser")
}

func (b *Browser) InitialPage() *Page {
	return &Page{browser: b, ctx: b.ctx, cancel: func() {}}
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
