// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openshelf",
	Short: "Manage your reading state files",
	Long: `OpenShelf CLI is a tool to view and update reading states of books.

You can use it to display the current state of a book or update
locations like page, percentage, or EPUB CFI positions.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {}
