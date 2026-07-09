package frontend

import (
	"html/template"
)

type compiledPageTemplate struct {
	PagePath string
	Template *template.Template
}

type RenderRequest struct {
	PageName  string
	PageTitle string
	Globals   any
	Page      any
}

type EngineService interface {
	Render(request RenderRequest) ([]byte, error)
	Reload() error
	Asset(assetPath string) ([]byte, bool)
	RawAsset(assetPath string) ([]byte, bool, error)
}

type frontendData struct {
	PageTitle          string
	HeadAssets         template.HTML
	Version            string
	PublicResourcePath string
}

type templateData struct {
	Frontend frontendData
	Static   any
	Globals  any
	Page     any
}
