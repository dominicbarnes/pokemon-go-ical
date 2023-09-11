package main

import "strings"

type LeekDuckEvent struct {
	ID        string `json:"eventID,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"eventType,omitempty"`
	Heading   string `json:"heading,omitempty"`
	Link      string `json:"link,omitempty"`
	Image     string `json:"image,omitempty"`
	Start     string `json:"start,omitempty"`
	End       string `json:"end,omitempty"`
	ExtraData any    `json:"extraData,omitempty"`
}

func (e LeekDuckEvent) Title() string {
	sb := new(strings.Builder)

	if e.Heading != "" {
		sb.WriteRune('[')
		sb.WriteString(e.Heading)
		sb.WriteRune(']')
		sb.WriteRune(' ')
	}

	sb.WriteString(e.Name)

	return sb.String()
}
