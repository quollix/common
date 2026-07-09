package frontend

import (
	"bytes"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"
)

type templateEngine struct{}

func (e *templateEngine) compileHTML(frameTemplateText, pageTemplateText string) (*htmltemplate.Template, error) {
	tmpl, err := htmltemplate.New("").Parse(frameTemplateText)
	if err != nil {
		return nil, logger.NewError(err.Error())
	}
	tmpl, err = tmpl.Parse(pageTemplateText)
	if err != nil {
		return nil, logger.NewError(err.Error())
	}
	return tmpl, nil
}

func (e *templateEngine) compileText(templateName, templateText string, data templateData) ([]byte, error) {
	tmpl, err := texttemplate.New(templateName).Option("missingkey=error").Parse(templateText)
	if err != nil {
		return nil, logger.NewError(err.Error(), "template_name", templateName)
	}
	var outputBuffer bytes.Buffer
	if err := tmpl.Execute(&outputBuffer, data); err != nil {
		return nil, logger.NewError(err.Error(), "template_name", templateName)
	}
	return outputBuffer.Bytes(), nil
}

func (e *templateEngine) preprocessPageTemplate(content string) string {
	var builder strings.Builder
	builder.WriteString(`{{define "content"}}` + "\n")
	builder.WriteString(content + "\n")
	builder.WriteString(`{{end}}` + "\n")
	return builder.String()
}
