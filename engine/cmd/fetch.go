// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ArcaneLean/openshelf/internal/cli"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch [book-file|book-id]",
	Short: "Fetch the reading state for a book",
	Long: `Fetch retrieves the current reading state for a given book.

You can specify either:
- a path to the book file
- or the book's unique bookId

The command outputs the reading state in JSON format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rs, err := cli.FetchState(args[0])
		if err != nil {
			return err
		}
		data, err := json.Marshal(rs)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
