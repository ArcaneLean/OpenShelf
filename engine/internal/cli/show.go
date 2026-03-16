// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ArcaneLean/openshelf/internal/state"
)

// ShowState prints reading state info with optional flags.
// If no flags are given, all sections are printed.
func ShowState(w io.Writer, path string, bookOnly, latestOnly bool, locationType string) error {
	rs, err := state.LoadReadingState(path)
	if err != nil {
		return fmt.Errorf("error loading reading state: %w", err)
	}

	// If no flags are given, act as if all are enabled
	if !bookOnly && !latestOnly && locationType == "" {
		bookOnly = true
		latestOnly = true
		locationType = "*" // special marker to print all locations
	}

	// --book flag
	if bookOnly {
		fmt.Fprintln(w, "BookID:", rs.BookID)
	}

	// --latest flag
	if latestOnly {
		locType, latest, found := rs.MostRecentLocation()
		if found {
			fmt.Fprintln(w, "Most recent location:")
			fmt.Fprintf(w, "  Type: %s\n", locType)
			fmt.Fprintf(w, "  Value: %v\n", latest.Value)
			fmt.Fprintf(w, "  UpdatedAt: %s\n", latest.UpdatedAt.Format(time.RFC3339))
		} else {
			fmt.Fprintln(w, "No locations found in this state.")
		}
	}

	// --location <type> flag or "*" (all)
	if locationType != "" {
		if locationType == "*" {
			allLocations, _ := json.MarshalIndent(rs.Locations, "  ", "  ")
			fmt.Fprintln(w, "All locations:", string(allLocations))
		} else if loc, ok := rs.Locations[locationType]; ok {
			fmt.Fprintf(w, "%s location:\n", locationType)
			fmt.Fprintf(w, "  Value: %v\n", loc.Value)
			fmt.Fprintf(w, "  UpdatedAt: %s\n", loc.UpdatedAt.Format(time.RFC3339))
		} else {
			fmt.Fprintf(w, "No location of type '%s' found.\n", locationType)
		}
	}

	return nil
}
