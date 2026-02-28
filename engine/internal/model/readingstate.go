package model

import (
	"time"
)

var interoperableLocationTypes = map[string]struct{}{
	"percentage":  {},
	"epubcfi":     {},
	"pageNumber":  {},
	"timeSeconds": {},
}

type Location struct {
	Value     any       `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ReadingState struct {
	SpecVersion string              `json:"specVersion"`
	BookID      string              `json:"bookId"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Locations   map[string]Location `json:"location"`
}

func (rs *ReadingState) MostRecentLocation() (string, Location, bool) {
	var latest Location
	var latestType string
	var found bool

	for t, l := range rs.Locations {
		if !found || l.UpdatedAt.After(latest.UpdatedAt) {
			latest = l
			latestType = t
			found = true
		}
	}

	return latestType, latest, found
}

func (rs *ReadingState) SetLocation(locType string, value any, t time.Time) {
	if rs.Locations == nil {
		rs.Locations = make(map[string]Location)
	}
	rs.Locations[locType] = Location{
		Value:     value,
		UpdatedAt: t,
	}
	rs.UpdatedAt = t
}

func IsInteroperable(locType string) bool {
	_, ok := interoperableLocationTypes[locType]
	return ok
}
