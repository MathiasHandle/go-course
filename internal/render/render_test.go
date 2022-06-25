package render

import (
	"net/http"
	"testing"

	"github.com/mathiashandle/go-course/internal/models"
)

func getSession() (*http.Request, error) {
	req, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := req.Context()
	ctx, _ = session.Load(ctx, req.Header.Get("X-Session"))

	req = req.WithContext(ctx)

	return req, nil
}

func TestAddDefaultData(t *testing.T) {
	var templateData models.TemplateData

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(req.Context(), "flash", "123")

	result := AddDefaultData(&templateData, req)
	if result.Flash != "123" {
		t.Error("flash value of 123 not ofund in session")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathtoTemplates = "./../../templates"

	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	appConfig.TemplateCache = templateCache

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var w myWriter

	err = RenderTemplate(&w, req, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&w, req, "non-existent.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(appConfig)
}

func TestCreateTemplateCache(t *testing.T) {
	pathtoTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
