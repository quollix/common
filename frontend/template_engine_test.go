package frontend

import (
	"bytes"
	"testing"

	"github.com/quollix/common/assert"
)

func TestTemplateEngine_CompileHTML_EscapesInScriptContext(t *testing.T) {
	engine := &templateEngine{}
	frameTemplateText := `{{define "frame"}}<h1>hello</h1>{{template "content" .}}{{end}}`
	pageTemplateText := `{{define "content"}}<script>const lessThan = "{{ "<" }}";</script>{{end}}`

	compiledTemplate, err := engine.compileHTML(frameTemplateText, pageTemplateText)
	assert.Nil(t, err)

	var outputBuffer bytes.Buffer
	err = compiledTemplate.ExecuteTemplate(&outputBuffer, "frame", nil)
	assert.Nil(t, err)

	expected := `<h1>hello</h1><script>const lessThan = "\u003c";</script>`
	assert.Equal(t, expected, outputBuffer.String())
}

func TestTemplateEngine_CompileHTML_InvalidFrameTemplateReturnsError(t *testing.T) {
	engine := &templateEngine{}

	_, err := engine.compileHTML(`{{define "frame"`, "")

	assert.NotNil(t, err)
}

func TestTemplateEngine_CompileHTML_InvalidPageTemplateReturnsError(t *testing.T) {
	engine := &templateEngine{}

	_, err := engine.compileHTML(`{{define "frame"}}{{end}}`, `{{define "content"`)

	assert.NotNil(t, err)
}

func TestTemplateEngine_CompileText_DoesNotEscapeInScriptContext(t *testing.T) {
	engine := &templateEngine{}
	data := templateData{
		Page: samplePageData{SignInPath: "/sign-in"},
	}
	templateText := `<script>const loginPath = "{{.Page.SignInPath}}"; const lessThan = "<";</script>`

	outputBytes, err := engine.compileText("sample", templateText, data)
	assert.Nil(t, err)

	expected := `<script>const loginPath = "/sign-in"; const lessThan = "<";</script>`
	assert.Equal(t, expected, string(outputBytes))
}

func TestTemplateEngine_CompileText_InvalidTemplateReturnsError(t *testing.T) {
	engine := &templateEngine{}

	_, err := engine.compileText("sample", "{{.Page.SignInPath", templateData{})

	assert.NotNil(t, err)
}

func TestTemplateEngine_CompileText_MissingKeyReturnsError(t *testing.T) {
	engine := &templateEngine{}

	_, err := engine.compileText("sample", "{{.Page.SignInPath}}", templateData{Page: map[string]string{}})

	assert.NotNil(t, err)
}

func TestTemplateEngine_PreprocessPageTemplate(t *testing.T) {
	engine := &templateEngine{}
	actual := engine.preprocessPageTemplate("<h1>Installed</h1>")

	expected := `{{define "content"}}
<h1>Installed</h1>
{{end}}
`

	assert.Equal(t, expected, actual)
}
