package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/mathiashandle/go-course/internal/config"
	"github.com/mathiashandle/go-course/internal/models"
)

var appConfig *config.AppConfig
var pathtoTemplates = "../../templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(templateData *models.TemplateData, req *http.Request) *models.TemplateData {
	templateData.Flash = appConfig.Session.PopString(req.Context(), "flash")
	templateData.Error = appConfig.Session.PopString(req.Context(), "error")
	templateData.Warning = appConfig.Session.PopString(req.Context(), "warning")

	templateData.CSRFToken = nosurf.Token(req)
	return templateData
}

// Renders passed in template
func RenderTemplate(w http.ResponseWriter, req *http.Request, tmpl string, templateData *models.TemplateData) error {
	var tc map[string]*template.Template

	if appConfig.UseCache {
		// get the template cache from app config
		tc = appConfig.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("couldnt get template from template cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, req)
	_ = t.Execute(buf, templateData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

var functions = template.FuncMap{}

// Parses all templates including layouts and returns them
func CreateTemplateCache() (map[string]*template.Template, error) {
	// create empty template cache
	templateCache := map[string]*template.Template{}

	// get all page templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathtoTemplates))
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

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathtoTemplates))
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathtoTemplates))
			if err != nil {
				return templateCache, err
			}
		}
		templateCache[name] = templateSet
	}

	return templateCache, nil
}
