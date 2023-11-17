package main

import (
	"context"
	"net/http"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/labstack/echo/v4"

	"app/internal/service"
)

const dsn = "http://admin:password@localhost:5984"
const dbName = "pokemon-go-ical"
const mimeICS = "text/calendar"

func main() {
	e := echo.New()
	e.Debug = true

	client, err := kivik.New("couch", dsn)
	if err != nil {
		e.Logger.Fatal("failed to initialize couchdb driver", err)
	}

	if err := client.CreateDB(context.TODO(), dbName); err != nil {
		e.Logger.Warnf("failed to create db %q: %w", dbName, err)
	}

	svc := service.Service{
		DB:   client.DB(dbName),
		HTTP: http.DefaultClient,
	}

	e.POST("/calendars", func(c echo.Context) error {
		ctx := c.Request().Context()

		var input service.CreateCalendarInput
		if err := c.Bind(&input); err != nil {
			return err
		}

		output, err := svc.CreateCalendar(ctx, &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, output)
	})

	e.PATCH("/calendars/:id", func(c echo.Context) error {
		ctx := c.Request().Context()

		var input service.UpdateCalendarInput
		if err := c.Bind(&input); err != nil {
			return err
		}

		output, err := svc.UpdateCalendar(ctx, &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, output)
	})

	e.GET("/calendars/:id", func(c echo.Context) error {
		ctx := c.Request().Context()

		var id string
		if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
			return err
		}

		output, err := svc.GetCalendar(ctx, id)
		if err != nil {
			return err
		}

		c.Response().Header().Set(echo.HeaderContentType, mimeICS)
		return c.String(http.StatusOK, output)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
