package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArcaneLean/openshelf/internal/state"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: openshelf <command> [args]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage: openshelf show <path-to-state-file>")
			os.Exit(1)
		}
		showState(os.Args[2])
	default:
		fmt.Println("Unknown command:", cmd)
	}
}

// showState loads a reading state file and prints summary info
func showState(path string) {
	rs, err := state.LoadReadingState(path)
	if err != nil {
		fmt.Println("Error loading reading state:", err)
		return
	}

	fmt.Println("BookID:", rs.BookID)
	fmt.Println("UpdatedAt:", rs.UpdatedAt.Format(time.RFC3339))

	locType, latest := rs.MostRecentLocation()
	if latest != nil {
		fmt.Println("Most recent location:")
		fmt.Printf("  Type: %s\n", locType)
		fmt.Printf("  Value: %v\n", latest.Value)
		fmt.Printf("  UpdatedAt: %s\n", latest.UpdatedAt.Format(time.RFC3339))
	} else {
		fmt.Println("No locations found in this state.")
	}

	// Optional: print all locations as JSON
	allLocations, _ := json.MarshalIndent(rs.Locations, "  ", "  ")
	fmt.Println("All locations:", string(allLocations))
}
