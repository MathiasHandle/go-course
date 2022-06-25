package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mathiashandle/go-course/internal/config"
	"github.com/mathiashandle/go-course/internal/handlers"
	"github.com/mathiashandle/go-course/internal/models"
	"github.com/mathiashandle/go-course/internal/render"
)

var appConfig config.AppConfig
var session *scs.SessionManager

const port = "127.0.0.1:3001"

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting app on port %s", port))

	serve := &http.Server{
		Addr:    port,
		Handler: routes(&appConfig),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	gob.Register(models.Reservation{})

	appConfig.InProduction = false

	// Setting up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction
	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.Newhandlers(repo)
	render.NewTemplates(&appConfig)

	return nil
}
