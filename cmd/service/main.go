package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/internal/service"
)

var couchDSN, couchDB string
var debug bool
var port int

func init() {
	flag.BoolVar(&debug, "debug", false, "Use to turn on debug logging")
	flag.IntVar(&port, "port", 3000, "The port number to bind to")
	flag.StringVar(&couchDSN, "couch-dsn", "http://localhost:5984", "The CouchDB server")
	flag.StringVar(&couchDB, "couch-database", "pokemon-go-ical", "The CouchDB database to use.")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf(":%d", port)

	e := echo.New()
	e.Debug = debug
	e.HideBanner = true
	// TODO: clean up logging
	e.Use(middleware.Logger())

	svc := service.Service{
		KivikDSN:      couchDSN,
		KivikDatabase: couchDB,
		HTTP:          http.DefaultClient,
	}

	e.GET("/", healthCheck)

	g := e.Group("/calendars")
	g.Use(middleware.BasicAuth(svc.AuthCouchDB))
	g.POST("", svc.CalendarCreate)
	g.PATCH("/:id", svc.CalendarUpdate)
	g.GET("/:id", svc.CalendarGet)

	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			e.Logger.Fatalf("shutting down the server: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
