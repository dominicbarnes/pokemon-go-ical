package service

import (
	"net/http"

	"github.com/go-kivik/kivik/v4"
)

type Service struct {
	// DB is the client for the configuration store.
	DB *kivik.DB

	// HTTP is a client used for fetching LeekDuck events.
	HTTP *http.Client
}
