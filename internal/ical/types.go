package ical

import "strings"

type LeekDuckEvent struct {
	ID      string `json:"eventID,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"eventType,omitempty"`
	Heading string `json:"heading,omitempty"`
	Link    string `json:"link,omitempty"`
	Image   string `json:"image,omitempty"`
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`

	ExtraData *LeekDuckEventExtraData `json:"extraData,omitempty"`
}

type LeekDuckEventPokemon struct {
	Name       string `json:"name,omitempty"`
	CanBeShiny bool   `json:"canBeShiny,omitempty"`
	Image      string `json:"image,omitempty"`
}

type LeekDuckEventBonus struct {
	Text  string `json:"text,omitempty"`
	Image string `json:"image,omitempty"`
}

type LeekDuckEventExtraData struct {
	Spotlight   *LeekDuckEventExtraDataSpotlight   `json:"spotlight,omitempty"`
	RaidBattles *LeekDuckEventExtraDataRaidBattles `json:"raidbattles,omitempty"`
	// TODO: communityday
}

type LeekDuckEventExtraDataSpotlight struct {
	LeekDuckEventPokemon
	Bonus string                 `json:"bonus,omitempty"`
	List  []LeekDuckEventPokemon `json:"list,omitempty"`
}

type LeekDuckEventExtraDataRaidBattles struct {
	Bosses  []LeekDuckEventPokemon `json:"bosses,omitempty"`
	Shinies []LeekDuckEventPokemon `json:"shinies,omitempty"`
}

type LeekDuckEventExtraDataCommunityDay struct {
	Spawns           []LeekDuckEventPokemon `json:"spawns,omitempty"`
	Bonuses          []LeekDuckEventBonus   `json:"bonuses,omitempty"`
	BonusDisclaimers []string               `json:"bonusDisclaimers,omitempty"`
	Shinies          []LeekDuckEventPokemon `json:"shinies,omitempty"`
	// TODO: specialresearch
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

	if e.Type == "pokemon-spotlight-hour" {
		if e.ExtraData != nil && e.ExtraData.Spotlight != nil && e.ExtraData.Spotlight.Bonus != "" {
			sb.WriteRune(' ')
			sb.WriteRune('(')
			sb.WriteString(e.ExtraData.Spotlight.Bonus)
			sb.WriteRune(')')
		}
	}

	return sb.String()
}

func (e LeekDuckEvent) Description() string {
	sb := new(strings.Builder)

	sb.WriteString("See more details at LeekDuck.com:\n\n")
	sb.WriteString(e.Link)

	return sb.String()
}
