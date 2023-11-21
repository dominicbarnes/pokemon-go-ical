package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/internal/service"
)

var couchDSN string
var couchDB string
var debug bool
var port int

func init() {
	flag.BoolVar(&debug, "debug", false, "Use to turn on debug logging")
	flag.IntVar(&port, "port", 3000, "The port number to bind to")
	flag.StringVar(&couchDSN, "couch-dsn", "http://admin:password@localhost:5984", "The CouchDB server")
	flag.StringVar(&couchDB, "couch-database", "pokemon-go-ical", "The CouchDB database to use.")
}

func main() {
	e := echo.New()
	e.Debug = debug
	e.HideBanner = true
	// TODO: clean up logging
	e.Use(middleware.Logger())

	client, err := kivik.New("couch", couchDSN)
	if err != nil {
		e.Logger.Fatal("failed to initialize couchdb driver", err)
	}

	if err := client.CreateDB(context.TODO(), couchDB); err != nil {
		e.Logger.Warnf("failed to create db %q: %w", couchDB, err)
	}

	svc := service.Service{
		DB:   client.DB(couchDB),
		HTTP: http.DefaultClient,
	}

	e.POST("/calendars", svc.CalendarCreate)
	e.PATCH("/calendars/:id", svc.CalendarUpdate)
	e.GET("/calendars/:id", svc.CalendarGet)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
