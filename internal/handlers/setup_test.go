package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/mathiashandle/go-course/internal/config"
	"github.com/mathiashandle/go-course/internal/models"
	"github.com/mathiashandle/go-course/internal/render"
)

var appConfig config.AppConfig
var session *scs.SessionManager

func GetRoutes() http.Handler {
	gob.Register(models.Reservation{})

	appConfig.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	// Setting up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction
	appConfig.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = true

	repo := NewRepo(&appConfig)
	Newhandlers(repo)
	render.NewTemplates(&appConfig)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)

	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.Generals)

	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)

	mux.Post("/search-availability", Repo.PostAvailability)

	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)

	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// Adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// Loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pathtoTemplates = "./../../templates"
var functions = template.FuncMap{}

// Parses all templates including layouts and returns them
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// create empty template cache
	templateCache := map[string]*template.Template{}

	// get all templates
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
