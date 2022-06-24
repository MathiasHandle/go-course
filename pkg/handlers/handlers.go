package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

	render.RenderTemplate(w, req, "home.page.tmpl", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := rep.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, req, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Startdate %s Enddate %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON renders the search availability page
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, req *http.Request) {
	res := jsonResponse{
		Ok:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(res, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, req, "contact.page.tmpl", &models.TemplateData{})
}
