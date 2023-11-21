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

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"app/internal/service"
)

var authUser, authPass string
var couchDSN, couchDB string
var debug bool
var port int

func init() {
	flag.StringVar(&authUser, "auth-user", "", "The HTTP Basic Auth user to require.")
	flag.StringVar(&authPass, "auth-pass", "", "The HTTP Basic Auth password to require.")
	flag.BoolVar(&debug, "debug", false, "Use to turn on debug logging")
	flag.IntVar(&port, "port", 3000, "The port number to bind to")
	flag.StringVar(&couchDSN, "couch-dsn", "http://admin:password@localhost:5984", "The CouchDB server")
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

	if authUser != "" && authPass != "" {
		e.Use(middleware.BasicAuth(func(user, pass string, c echo.Context) (bool, error) {
			return user == authUser && pass == authPass, nil
		}))
	}

	client, err := initCouchDB(context.TODO())
	if err != nil {
		e.Logger.Fatal(err)
	}

	svc := service.Service{
		DB:   client.DB(couchDB),
		HTTP: http.DefaultClient,
	}

	e.GET("/", healthCheck)
	e.POST("/calendars", svc.CalendarCreate)
	e.PATCH("/calendars/:id", svc.CalendarUpdate)
	e.GET("/calendars/:id", svc.CalendarGet)

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

func initCouchDB(ctx context.Context) (*kivik.Client, error) {
	client, err := kivik.New("couch", couchDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize couchdb driver: %w", err)
	}

	if ok, err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping couchdb: %w", err)
	} else if !ok {
		return nil, errors.New("could not connect to database")
	}

	if ok, err := client.DBExists(ctx, couchDB); err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	} else if !ok {
		return nil, fmt.Errorf("database %q not found", couchDB)
	}

	return client, nil
}

func healthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
