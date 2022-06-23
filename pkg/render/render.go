package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/mathiashandle/go-course/pkg/config"
	"github.com/mathiashandle/go-course/pkg/models"
)

var appConfig *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Renders passed in template
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	// get the template cache from app config
	tc := appConfig.TemplateCache

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Couldnt get template from template cache")
	}

	buf := new(bytes.Buffer)
	templateData = AddDefaultData(templateData)
	_ = t.Execute(buf, templateData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

var functions = template.FuncMap{}

// Parses all templates including layouts and returns them
func CreateTemplateCache() (map[string]*template.Template, error) {
	// create empty template cache
	templateCache := map[string]*template.Template{}

	// get all templates
	pages, err := filepath.Glob("../../templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		// path to current template
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}
		templateCache[name] = templateSet
	}

	return templateCache, nil
}
