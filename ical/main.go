package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	location := request.QueryStringParameters["timezone"]

	tz, err := time.LoadLocation(location)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("timezone %s is invalid: %w", location, err)
	}

	ee, err := getEvents()
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to get leek duck events: %w", err)
	}

	options := GenerateICalOptions{
		Now:          time.Now(),
		TZ:           tz,
		IncludeTypes: parseList(request.QueryStringParameters["include"]),
		ExcludeTypes: parseList(request.QueryStringParameters["exclude"]),
	}

	log.Printf("Generating ICal from %d events using %+v\n", len(ee), options)

	ics, err := GenerateICal(ee, options)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to generate ics: %w", err)
	}

	log.Printf("Generated ICal with %d events\n", len(ics.Events()))

	return events.APIGatewayProxyResponse{
		Body:       ics.Serialize(),
		Headers:    map[string]string{"Content-Type": "text/calendar"},
		StatusCode: 200,
	}, nil
}

const sourceURL = "https://raw.githubusercontent.com/bigfoott/ScrapedDuck/data/events.min.json"

func getEvents() ([]LeekDuckEvent, error) {
	res, err := http.Get(sourceURL)
	if err != nil {
		return nil, errors.New("failed to download events")
	}
	defer res.Body.Close()

	var ee []LeekDuckEvent
	d := json.NewDecoder(res.Body)
	if err := d.Decode(&ee); err != nil {
		return nil, fmt.Errorf("failed to decode events from leek duck as JSON: %w", err)
	}
	return ee, nil
}

func parseList(raw string) []string {
	input := strings.TrimSpace(raw)

	if input == "" {
		return nil
	}

	return strings.Split(input, ",")
}
