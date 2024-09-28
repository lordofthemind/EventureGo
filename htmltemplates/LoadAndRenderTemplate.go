package htmltemplates

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed templates/*.html
var emailTemplates embed.FS

// LoadAndRenderTemplate loads and renders the specified template with the provided data
func LoadAndRenderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFS(emailTemplates, "templates/"+templateName)
	if err != nil {
		return "", err
	}

	var renderedTemplate strings.Builder
	err = tmpl.Execute(&renderedTemplate, data)
	if err != nil {
		return "", err
	}

	return renderedTemplate.String(), nil
}
