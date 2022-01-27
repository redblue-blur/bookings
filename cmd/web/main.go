package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/redblue-blur/bookings/internal/config"
	"github.com/redblue-blur/bookings/internal/handlers"
	"github.com/redblue-blur/bookings/internal/helpers"
	"github.com/redblue-blur/bookings/internal/models"
	"github.com/redblue-blur/bookings/internal/render"
)

const portno = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("starting application on port %s", portno)

	srv := &http.Server{
		Addr:    portno,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
func run() error {
	//what am i going to put in the session
	gob.Register(models.Reservation{})
	//true if in Production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return nil
}
