package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kivik/v4/couchdb"
	"github.com/labstack/echo/v4"
)

type Service struct {
	KivikDSN      string
	KivikDatabase string

	// HTTP is a client used for fetching LeekDuck events.
	HTTP *http.Client
}

func (s *Service) AuthCouchDB(user, pass string, c echo.Context) (bool, error) {
	ctx := c.Request().Context()
	db, err := initCouchDB(ctx, s.KivikDSN, s.KivikDatabase, user, pass)
	if err != nil {
		return false, err
	}
	c.Set("KivikDB", db)
	return true, nil
}

func validateTZ(input string) (*time.Location, error) {
	tz, err := time.LoadLocation(input)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "timezone is invalid")
	}
	return tz, nil
}

func initCouchDB(ctx context.Context, dsn, db, user, pass string) (*kivik.DB, error) {
	client, err := kivik.New("couch", dsn, couchdb.BasicAuth(user, pass))
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
