// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ArcaneLean/openshelf/internal/cli"
	"github.com/spf13/cobra"
)

var timeFlag string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <path-to-state-file> <location-type> <value>",
	Short: "Update a location in a reading state file",
	Long: `Update the reading state of a book.

You can set the value of a location type (e.g., page, percentage, epubcfi)
and optionally provide a timestamp. If no timestamp is given, the current time is used.`,
	Args: cobra.ExactArgs(3), // requires path, location type, and value
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		locType := args[1]
		valueArg := args[2]

		// Parse timestamp if provided
		t := time.Now()
		if timeFlag != "" {
			parsedTime, err := time.Parse(time.RFC3339, timeFlag)
			if err != nil {
				return fmt.Errorf("invalid time format, must be RFC3339: %w", err)
			}
			t = parsedTime
		}

		// Convert value to int or float if possible, else leave as string
		var value any
		if intVal, err := strconv.Atoi(valueArg); err == nil {
			value = intVal
		} else if floatVal, err := strconv.ParseFloat(valueArg, 64); err == nil {
			value = floatVal
		} else {
			value = valueArg
		}

		if err := cli.UpdateState(os.Stdout, path, locType, value, t); err != nil {
			return err
		}

		fmt.Println("State updated successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Local flag
	updateCmd.Flags().StringVar(&timeFlag, "time", "", "Optional timestamp in RFC3339 format (default: now)")
}
