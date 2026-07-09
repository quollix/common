package frontend

import (
	"bytes"
	"errors"
	"io/fs"
	"path"
	"strings"
)

const frameTemplateName = "frame"

type Engine struct {
	config         resolvedConfig
	assetStore     *AssetStore
	assetTags      *assetTagBuilder
	templateEngine *templateEngine
	compiledPages  map[string]*compiledPageTemplate
}

var _ EngineService = (*Engine)(nil)

func NewEngine(config Config) (*Engine, error) {
	resolvedConfig, err := config.resolve()
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		config:         resolvedConfig,
		assetStore:     NewAssetStore(resolvedConfig.version),
		templateEngine: &templateEngine{},
		compiledPages:  map[string]*compiledPageTemplate{},
	}
	engine.assetTags = &assetTagBuilder{
		assetStore:         engine.assetStore,
		publicResourcePath: resolvedConfig.publicResourcePath,
		globalAssetNames:   []string{"global/frame", "global/global"},
	}

	if err := engine.Reload(); err != nil {
		return nil, err
	}

	return engine, nil
}

func (e *Engine) Reload() error {
	e.assetStore.Clear()
	if err := e.walkResourcesAndAddAssets("js"); err != nil {
		return err
	}
	if err := e.walkResourcesAndAddAssets("css"); err != nil {
		return err
	}
	if err := e.walkResourcesAndAddAssets("woff2"); err != nil {
		return err
	}

	framePath := e.resourcePath(e.config.framePath)
	frameBytes, err := fs.ReadFile(e.config.fileSystem, framePath)
	if err != nil {
		return logger.NewError(err.Error(), "path", framePath)
	}

	e.compiledPages, err = e.walkHTMLTemplatesAndCompile(string(frameBytes))
	return err
}

func (e *Engine) Render(request RenderRequest) ([]byte, error) {
	request.PageName = strings.TrimSpace(request.PageName)
	if request.PageName == "" {
		return nil, logger.NewError("page name must not be empty")
	}

	compiledPage, ok := e.compiledPages[request.PageName]
	if !ok {
		return nil, logger.NewError("template not found", "page_name", request.PageName)
	}

	pageTitle := request.PageTitle
	if pageTitle == "" {
		pageTitle = deriveTitle(request.PageName)
	}

	data := templateData{
		Frontend: frontendData{
			PageTitle:          pageTitle,
			HeadAssets:         e.assetTags.buildHeadAssets(compiledPage.PagePath),
			Version:            e.config.version,
			PublicResourcePath: e.config.publicResourcePath,
		},
		Static:  e.config.static,
		Globals: request.Globals,
		Page:    request.Page,
	}

	var outputBuffer bytes.Buffer
	if err := compiledPage.Template.ExecuteTemplate(&outputBuffer, frameTemplateName, data); err != nil {
		return nil, logger.NewError(err.Error(), "page_name", request.PageName)
	}

	return outputBuffer.Bytes(), nil
}

func (e *Engine) Asset(assetPath string) ([]byte, bool) {
	return e.assetStore.Get(assetPath)
}

func (e *Engine) RawAsset(assetPath string) ([]byte, bool, error) {
	relativePath, ok := strings.CutPrefix(assetPath, e.config.publicResourcePath+"/")
	if !ok {
		return nil, false, logger.NewError("asset path must start with public resource path", "path", assetPath, "public_resource_path", e.config.publicResourcePath)
	}

	contentBytes, err := fs.ReadFile(e.config.fileSystem, e.resourcePath(relativePath))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, logger.NewError(err.Error(), "path", assetPath)
	}
	return contentBytes, true, nil
}

func (e *Engine) resourceRoot() string {
	return "."
}

func (e *Engine) resourcePath(relativePath string) string {
	return strings.Trim(relativePath, "/")
}

func (e *Engine) walkResourcesAndAddAssets(fileExtension string) error {
	return fs.WalkDir(e.config.fileSystem, e.resourceRoot(), func(currentPath string, dirEntry fs.DirEntry, walkError error) error {
		return e.addAsset(currentPath, dirEntry, walkError, fileExtension)
	})
}

func (e *Engine) walkHTMLTemplatesAndCompile(frameText string) (map[string]*compiledPageTemplate, error) {
	compiledPages := map[string]*compiledPageTemplate{}
	err := fs.WalkDir(e.config.fileSystem, e.resourceRoot(), func(currentPath string, dirEntry fs.DirEntry, walkError error) error {
		return e.compileHTMLTemplate(currentPath, dirEntry, walkError, frameText, compiledPages)
	})
	return compiledPages, err
}

func (e *Engine) compileHTMLTemplate(currentPath string, dirEntry fs.DirEntry, walkError error, frameText string, compiledPages map[string]*compiledPageTemplate) error {
	if walkError != nil {
		return logger.NewError(walkError.Error(), "current_path", currentPath)
	}
	if dirEntry.IsDir() || !strings.HasSuffix(currentPath, ".html") {
		return nil
	}

	relativePath := strings.TrimPrefix(currentPath, e.resourceRoot()+"/")
	if relativePath == e.config.framePath {
		return nil
	}

	pageBytes, err := fs.ReadFile(e.config.fileSystem, currentPath)
	if err != nil {
		return logger.NewError(err.Error(), "path", currentPath)
	}

	pageName := strings.ToLower(strings.TrimSuffix(path.Base(currentPath), ".html"))
	pageText := e.templateEngine.preprocessPageTemplate(string(pageBytes))

	compiledTemplate, err := e.templateEngine.compileHTML(frameText, pageText)
	if err != nil {
		return logger.AddContext(err, "page_name", pageName, "path", currentPath)
	}

	compiledPages[pageName] = &compiledPageTemplate{
		Template: compiledTemplate,
		PagePath: relativePath,
	}

	logger.Debug("compiled template", "name", pageName, "path", currentPath)
	return nil
}

func (e *Engine) addAsset(currentPath string, dirEntry fs.DirEntry, walkError error, fileExtension string) error {
	if walkError != nil {
		return logger.NewError(walkError.Error(), "current_path", currentPath)
	}
	if dirEntry.IsDir() || !strings.HasSuffix(currentPath, "."+fileExtension) {
		return nil
	}

	fileBytes, err := fs.ReadFile(e.config.fileSystem, currentPath)
	if err != nil {
		return logger.NewError(err.Error(), "path", currentPath)
	}

	relativePath := strings.TrimPrefix(currentPath, e.resourceRoot()+"/")
	fileNameWithoutExtension := strings.TrimSuffix(path.Base(relativePath), "."+fileExtension)
	fileFolder := path.Dir(relativePath)
	injectedPath := e.assetStore.VersionedPath(e.config.publicResourcePath, path.Join(fileFolder, fileNameWithoutExtension), fileExtension)

	renderedBytes := fileBytes
	if !isVendorAsset(relativePath) {
		renderedBytes, err = e.templateEngine.compileText(currentPath, string(fileBytes), e.staticTemplateData())
		if err != nil {
			return err
		}
	}

	e.assetStore.Put(injectedPath, renderedBytes)
	logger.Debug("generated injected web resource", "path", currentPath)
	return nil
}

func (e *Engine) staticTemplateData() templateData {
	return templateData{
		Frontend: frontendData{
			Version:            e.config.version,
			PublicResourcePath: e.config.publicResourcePath,
		},
		Static: e.config.static,
	}
}

func deriveTitle(pageName string) string {
	s := strings.TrimSpace(pageName)
	s = strings.TrimSuffix(s, ".html")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		logger.Error("derive title of empty page", "page_name", pageName)
		return ""
	}

	lowerTitle := strings.ToLower(s)
	return strings.ToUpper(lowerTitle[:1]) + lowerTitle[1:]
}

func isVendorAsset(relativePath string) bool {
	return strings.HasPrefix(relativePath, "vendor/")
}
