package frontend

import (
	"fmt"
	"html/template"
	"strings"
)

type assetTagBuilder struct {
	assetStore         *AssetStore
	publicResourcePath string
	globalAssetNames   []string
}

func (b *assetTagBuilder) buildHeadAssets(pageHtmlPath string) template.HTML {
	basePaths := append([]string{}, b.globalAssetNames...)
	basePath := strings.TrimSuffix(pageHtmlPath, ".html")
	basePaths = append(basePaths, basePath)

	var builder strings.Builder
	for _, file := range basePaths {
		builder.WriteString(b.cssLinkTag(file))
		builder.WriteString(b.jsScriptTag(file))
	}

	return template.HTML(builder.String()) // #nosec G203: Asset paths are generated from configured local resource names.
}

func (b *assetTagBuilder) jsScriptTag(file string) string {
	assetPath := b.assetStore.VersionedPath(b.publicResourcePath, file, "js")
	if b.assetStore.Has(assetPath) {
		return fmt.Sprintf(`<script type="module" src="%s"></script>`+"\n", assetPath)
	}
	return ""
}

func (b *assetTagBuilder) cssLinkTag(file string) string {
	assetPath := b.assetStore.VersionedPath(b.publicResourcePath, file, "css")
	if b.assetStore.Has(assetPath) {
		return fmt.Sprintf(`<link rel="stylesheet" href="%s">`+"\n", assetPath)
	}
	return ""
}
