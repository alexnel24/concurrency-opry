package models

import (
	"fmt"
)

type Performance struct {
	Id          int64
	EventLink   string
	ArtistName  string
	ComboString string
}

func NewPerformance(artistName, eventLink string) *Performance {
	return &Performance{ArtistName: artistName, EventLink: eventLink, ComboString: fmt.Sprintf("%s-%s", artistName, eventLink)}
}

//ToDo: update the time to track which performances are past (potentially add upcoming bool)