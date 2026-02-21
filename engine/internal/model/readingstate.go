package model

import "time"

var InteroperableLocationTypes = []string{
	"percentage",
	"epubcfi",
	"pageNumber",
	"timeSeconds",
}

type Location struct {
	Value     interface{} `json:"value"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type ReadingState struct {
	SpecVersion string              `json:"specVersion"`
	BookID      string              `json:"bookId"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Locations   map[string]Location `json:"location"`
}

func (rs *ReadingState) MostRecentLocation() (locType string, loc *Location) {
	var latest *Location
	var latestType string
	for t, l := range rs.Locations {
		if latest == nil || l.UpdatedAt.After(latest.UpdatedAt) {
			copyLoc := l // avoid pointer to iteration variable
			latest = &copyLoc
			latestType = t
		}
	}
	return latestType, latest
}

func (rs *ReadingState) SetLocation(locType string, value interface{}, t time.Time) {
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
	for _, t := range InteroperableLocationTypes {
		if t == locType {
			return true
		}
	}
	return false
}
