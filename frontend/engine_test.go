package frontend

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	"github.com/quollix/deepstack"
)

type sampleStatic struct {
	Brand    string
	IconPath string
}

type samplePageData struct {
	Message    string
	SignInPath string
}

type sampleGlobals struct {
	Username string
}

func newTestEngine(t *testing.T) *Engine {
	engine, err := NewEngine(Config{
		FrontendFolderPath: "testdata/frontend",
		Version:            "abc123",
		Static: sampleStatic{
			Brand:    "Sample",
			IconPath: "/frontend/resources/images/logo.svg",
		},
	})
	assert.Nil(t, err)

	return engine
}

func TestNewEngine_RequiresFrontendFolderPath(t *testing.T) {
	_, err := NewEngine(Config{})

	deepstack.AssertDeepStackError(t, err, "frontend folder path must not be empty")
}

func TestEngine_Render_UsesFixtureAndInjectsAssets(t *testing.T) {
	engine := newTestEngine(t)

	outputBytes, err := engine.Render(RenderRequest{
		PageName:  "sample",
		PageTitle: "Custom title",
		Globals:   sampleGlobals{Username: "admin"},
		Page:      samplePageData{Message: "Hello from page data", SignInPath: "/sign-in"},
	})
	assert.Nil(t, err)

	output := string(outputBytes)
	assertContains(t, output, "<title>Custom title - Sample</title>")
	assertContains(t, output, `<link rel="icon" href="/frontend/resources/images/logo.svg" type="image/svg">`)
	assertContains(t, output, `/frontend/resources/global/frame.abc123.css`)
	assertContains(t, output, `/frontend/resources/global/frame.abc123.js`)
	assertContains(t, output, `/frontend/resources/global/global.abc123.css`)
	assertContains(t, output, `/frontend/resources/global/global.abc123.js`)
	assertContains(t, output, `/frontend/resources/pages/sample/sample.abc123.css`)
	assertContains(t, output, `/frontend/resources/pages/sample/sample.abc123.js`)
	assertContains(t, output, `<p id="sample-message">Hello from page data</p>`)
	assertContains(t, output, `<p id="sample-username">admin</p>`)
	assertContains(t, output, `<a id="sample-sign-in-link" href="/sign-in">Sign in</a>`)
}

func TestEngine_Render_UsesPageDataWhenProvided(t *testing.T) {
	engine := newTestEngine(t)

	outputBytes, err := engine.Render(RenderRequest{
		PageName: "sample",
		Globals:  sampleGlobals{Username: "admin"},
		Page: samplePageData{
			Message:    "Hello from page data",
			SignInPath: "/tenant/sign-in",
		},
	})
	assert.Nil(t, err)

	output := string(outputBytes)
	assertContains(t, output, "<title>Sample - Sample</title>")
	assertContains(t, output, `<a id="sample-sign-in-link" href="/tenant/sign-in">Sign in</a>`)
}

func TestEngine_Render_DerivesPageTitle(t *testing.T) {
	engine := newTestEngine(t)

	outputBytes, err := engine.Render(RenderRequest{
		PageName: "render-contract",
		Page:     samplePageData{Message: "contract"},
	})
	assert.Nil(t, err)

	output := string(outputBytes)
	assertContains(t, output, "<title>Render contract - Sample</title>")
	assertContains(t, output, `<p id="render-contract-message">contract</p>`)
	assertContains(t, output, `<p id="render-contract-brand">Sample</p>`)
	assertContains(t, output, `<p id="render-contract-version">abc123</p>`)
	assertContains(t, output, `<p id="render-contract-public-path">/frontend/resources</p>`)
}

func TestEngine_Render_PageWithoutPageAssetsOnlyIncludesGlobalAssets(t *testing.T) {
	engine := newTestEngine(t)

	outputBytes, err := engine.Render(RenderRequest{
		PageName: "no-assets",
		Page:     samplePageData{Message: "no page assets"},
	})
	assert.Nil(t, err)

	output := string(outputBytes)
	assertContains(t, output, `/frontend/resources/global/frame.abc123.css`)
	assertContains(t, output, `/frontend/resources/global/global.abc123.js`)
	assertNotContains(t, output, `/frontend/resources/pages/no-assets/no-assets.abc123.css`)
	assertNotContains(t, output, `/frontend/resources/pages/no-assets/no-assets.abc123.js`)
}

func TestEngine_Render_MissingPageReturnsDeepStackError(t *testing.T) {
	engine := newTestEngine(t)

	_, err := engine.Render(RenderRequest{PageName: "missing"})

	deepstack.AssertDeepStackError(t, err, "template not found", "page_name", "missing")
}

func TestEngine_Render_BlankPageNameReturnsError(t *testing.T) {
	engine := newTestEngine(t)

	_, err := engine.Render(RenderRequest{PageName: " "})

	deepstack.AssertDeepStackError(t, err, "page name must not be empty")
}

func TestEngine_Render_InvalidPageDataReturnsError(t *testing.T) {
	engine := newTestEngine(t)

	_, err := engine.Render(RenderRequest{
		PageName: "sample",
		Page:     struct{}{},
	})

	deepstack.AssertDeepStackError(t, err, "template: :3:34: executing \"content\" at <.Page.Message>: can't evaluate field Message in type interface {}", "page_name", "sample")
}

func TestEngine_AssetsAreRenderedAndVendorAssetsAreNotRendered(t *testing.T) {
	engine := newTestEngine(t)

	sampleJS, ok := engine.Asset("/frontend/resources/pages/sample/sample.abc123.js")
	assert.True(t, ok)
	assertContains(t, string(sampleJS), `brand: "Sample"`)
	assertContains(t, string(sampleJS), `lessThan: "<"`)

	vendorCSS, ok := engine.Asset("/frontend/resources/vendor/css/vendor.abc123.css")
	assert.True(t, ok)
	assertContains(t, string(vendorCSS), `{{.Missing.Value}}`)

	fontBytes, ok := engine.Asset("/frontend/resources/global/font.abc123.woff2")
	assert.True(t, ok)
	assert.Equal(t, []byte("sample font\n"), fontBytes)
}

func TestEngine_RawAsset_ReadsUnversionedResource(t *testing.T) {
	engine := newTestEngine(t)

	contentBytes, ok, err := engine.RawAsset("/frontend/resources/vendor/css/vendor.css")
	assert.Nil(t, err)
	assert.True(t, ok)
	assertContains(t, string(contentBytes), `{{.Missing.Value}}`)
}

func TestEngine_RawAsset_MissingResourceReturnsNotFound(t *testing.T) {
	engine := newTestEngine(t)

	contentBytes, ok, err := engine.RawAsset("/frontend/resources/missing.css")

	assert.Nil(t, err)
	assert.False(t, ok)
	assert.Nil(t, contentBytes)
}

func TestEngine_RawAsset_OutsidePublicResourcePathReturnsError(t *testing.T) {
	engine := newTestEngine(t)

	contentBytes, ok, err := engine.RawAsset("/other/resources/vendor/css/vendor.css")

	deepstack.AssertDeepStackError(t, err, "asset path must start with public resource path", "path", "/other/resources/vendor/css/vendor.css", "public_resource_path", "/frontend/resources")
	assert.False(t, ok)
	assert.Nil(t, contentBytes)
}

func TestDeriveTitle_BlankPageNameReturnsEmptyTitle(t *testing.T) {
	actualTitle := deriveTitle(" .html ")

	assert.Equal(t, "", actualTitle)
}

func assertContains(t *testing.T, actual, expectedSubstring string) {
	assert.True(t, strings.Contains(actual, expectedSubstring))
}

func assertNotContains(t *testing.T, actual, unexpectedSubstring string) {
	assert.False(t, strings.Contains(actual, unexpectedSubstring))
}
