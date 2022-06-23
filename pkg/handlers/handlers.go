package handlers

import (
	"net/http"

	"github.com/mathiashandle/go-course/pkg/config"
	"github.com/mathiashandle/go-course/pkg/models"
	"github.com/mathiashandle/go-course/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func Newhandlers(repo *Repository) {
	Repo = repo
}

// Home page handler
func (rep *Repository) Home(w http.ResponseWriter, req *http.Request) {
	remoteIp := req.RemoteAddr
	rep.App.Session.Put(req.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := rep.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
