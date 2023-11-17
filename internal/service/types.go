package service

type CalendarConfig struct {
	ID            string   `json:"_id"`
	Rev           string   `json:"_rev,omitempty"`
	Timezone      string   `json:"timezone"`
	IncludeEvents []string `json:"include_events"`
	ExcludeEvents []string `json:"exclude_events"`
}
