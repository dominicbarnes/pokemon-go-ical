package service

import (
	"net/http"
	"time"

	"github.com/go-kivik/kivik/v4"
	"github.com/labstack/echo/v4"
)

type Service struct {
	KivikDB *kivik.DB

	// HTTP is a client used for fetching LeekDuck events.
	HTTP *http.Client
}

func validateTZ(input string) (*time.Location, error) {
	tz, err := time.LoadLocation(input)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "timezone is invalid")
	}
	return tz, nil
}
