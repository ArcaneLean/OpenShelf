// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"os"

	"github.com/ArcaneLean/openshelf/internal/cli"
	"github.com/spf13/cobra"
)

var (
	bookFlag     bool
	latestFlag   bool
	locationFlag string
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <path-to-state-file>",
	Short: "Display the reading state of a book",
	Long: `Show the reading state of a book from a state file.

You can print only the BookID, the latest location, or a specific location type.
If no flags are provided, all sections are displayed.`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: the path to the state file
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// If no flags are given, treat as all enabled
		if !bookFlag && !latestFlag && locationFlag == "" {
			bookFlag = true
			latestFlag = true
			locationFlag = "*" // print all locations
		}

		if err := cli.ShowState(os.Stdout, path, bookFlag, latestFlag, locationFlag); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Local flags
	showCmd.Flags().BoolVar(&bookFlag, "book", false, "Print only the BookID")
	showCmd.Flags().BoolVar(&latestFlag, "latest", false, "Print only the most recent location")
	showCmd.Flags().StringVar(&locationFlag, "location", "", "Print a specific location type (e.g., page, percentage, epubcfi)")
}
