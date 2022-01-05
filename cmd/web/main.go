package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/redblue-blur/bookings/internal/config"
	"github.com/redblue-blur/bookings/internal/handlers"
	"github.com/redblue-blur/bookings/internal/render"
)

const portno = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	//true if in Production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("starting application on port %s", portno)

	srv := &http.Server{
		Addr:    portno,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
