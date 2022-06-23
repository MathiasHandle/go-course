package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mathiashandle/go-course/pkg/config"
	"github.com/mathiashandle/go-course/pkg/handlers"
	"github.com/mathiashandle/go-course/pkg/render"
)

var appConfig config.AppConfig
var session *scs.SessionManager

const port = "127.0.0.1:3001"

func main() {
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
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.Newhandlers(repo)
	render.NewTemplates(&appConfig)

	fmt.Println(fmt.Sprintf("Starting app on port %s", port))

	serve := &http.Server{
		Addr:    port,
		Handler: routes(&appConfig),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
