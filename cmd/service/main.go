package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/segmentio/events/v2"

	"app/internal/service"
)

var couchDSN, couchDB string
var writeUser, writePass string
var debug bool
var port int

func init() {
	flag.BoolVar(&debug, "debug", false, "Use to turn on debug logging")
	flag.IntVar(&port, "port", 3000, "The port number to bind to")
	flag.StringVar(&couchDSN, "couch-dsn", "http://localhost:5984", "The CouchDB server")
	flag.StringVar(&couchDB, "couch-database", "pokemon-go-ical", "The CouchDB database to use")
	flag.StringVar(&writeUser, "write-user", "allow-writes", "Required for write endpoints")
	flag.StringVar(&writePass, "write-pass", "", "Required for write endpoints")
	flag.Parse()
}

func main() {
	ctx, cancel := events.WithSignals(context.Background(), os.Interrupt)
	defer cancel()

	addr := fmt.Sprintf(":%d", port)

	e := echo.New()
	e.Debug = debug
	e.HideBanner = true
	// TODO: clean up logging
	e.Use(middleware.Logger())

	db, err := initCouchDB(ctx, couchDSN, couchDB)
	if err != nil {
		e.Logger.Fatalf("failed to initialize couchdb client: %w", err)
	}

	svc := service.Service{
		KivikDB: db,
		HTTP:    http.DefaultClient,
	}

	// read operations (anonymous)
	e.GET("/", healthCheck)
	e.GET("/calendars/:id", svc.CalendarGet)

	// write operations (require write auth)
	g := e.Group("/calendars")
	g.Use(middleware.BasicAuth(basicAuth))
	g.POST("", svc.CalendarCreate)
	g.PATCH("/:id", svc.CalendarUpdate)

	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			e.Logger.Fatalf("shutting down the server: %w", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func initCouchDB(ctx context.Context, dsn, db string) (*kivik.DB, error) {
	client, err := kivik.New("couch", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize couchdb driver: %w", err)
	}

	if ok, err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping couchdb: %w", err)
	} else if !ok {
		return nil, errors.New("could not connect to database")
	}

	if ok, err := client.DBExists(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	} else if !ok {
		return nil, fmt.Errorf("database %q does not exist", db)
	}

	return client.DB(db), nil
}

func basicAuth(user, pass string, c echo.Context) (bool, error) {
	return user == writeUser && pass == writePass, nil
}
