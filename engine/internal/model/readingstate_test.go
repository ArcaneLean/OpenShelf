// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package model

import (
	"testing"
	"time"
)

func TestMostRecentLocation(t *testing.T) {
	rs := ReadingState{
		Locations: map[string]Location{
			"percentage": {
				Value:     10.0,
				UpdatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			"page": {
				Value:     123,
				UpdatedAt: time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC),
			},
		},
	}

	locType, loc, found := rs.MostRecentLocation()

	if locType != "page" {
		t.Errorf("expected page, got %s", locType)
	}

	if loc.Value != 123 {
		t.Errorf("expected 123, got %v", loc.Value)
	}

	if found != true {
		t.Errorf("expected true, got %t", found)
	}
}

func TestSetLocation_InitializesMapAndSetsValues(t *testing.T) {
	rs := &ReadingState{}

	now := time.Date(2026, 1, 25, 10, 15, 0, 0, time.UTC)

	rs.SetLocation("percentage", 42.3, now)

	// 1️⃣ Map should be initialized
	if rs.Locations == nil {
		t.Fatal("Locations map was not initialized")
	}

	// 2️⃣ Location should exist
	loc, ok := rs.Locations["percentage"]
	if !ok {
		t.Fatal("Location 'percentage' not set")
	}

	// 3️⃣ Value should match
	if loc.Value != 42.3 {
		t.Errorf("Expected value 42.3, got %v", loc.Value)
	}

	// 4️⃣ Location timestamp should match
	if !loc.UpdatedAt.Equal(now) {
		t.Errorf("Expected location UpdatedAt %v, got %v", now, loc.UpdatedAt)
	}

	// 5️⃣ Top-level UpdatedAt should match
	if !rs.UpdatedAt.Equal(now) {
		t.Errorf("Expected ReadingState UpdatedAt %v, got %v", now, rs.UpdatedAt)
	}
}

func TestSetLocation_OverwritesExistingLocation(t *testing.T) {
	initialTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	newTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	rs := &ReadingState{
		Locations: map[string]Location{
			"page": {
				Value:     100,
				UpdatedAt: initialTime,
			},
		},
		UpdatedAt: initialTime,
	}

	rs.SetLocation("page", 200, newTime)

	loc := rs.Locations["page"]

	if loc.Value != 200 {
		t.Errorf("Expected value 200, got %v", loc.Value)
	}

	if !loc.UpdatedAt.Equal(newTime) {
		t.Errorf("Expected UpdatedAt %v, got %v", newTime, loc.UpdatedAt)
	}

	if !rs.UpdatedAt.Equal(newTime) {
		t.Errorf("Expected ReadingState UpdatedAt %v, got %v", newTime, rs.UpdatedAt)
	}
}

func TestIsInteroperable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"percentage is interoperable", "percentage", true},
		{"epubcfi is interoperable", "epubcfi", true},
		{"pageNumber is interoperable", "pageNumber", true},
		{"timeSeconds is interoperable", "timeSeconds", true},
		{"unknown type is not interoperable", "customType", false},
		{"empty string is not interoperable", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInteroperable(tt.input)
			if result != tt.expected {
				t.Errorf("IsInteroperable(%q) = %t, expected %t",
					tt.input, result, tt.expected)
			}
		})
	}
}
