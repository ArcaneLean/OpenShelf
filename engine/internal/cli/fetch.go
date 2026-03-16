// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cli

import (
	"fmt"
	"os"

	"github.com/ArcaneLean/openshelf/internal/identity"
	"github.com/ArcaneLean/openshelf/internal/model"
	"github.com/ArcaneLean/openshelf/internal/state"
)

func FetchState(bookOrFile string) (*model.ReadingState, error) {
	var bookID string
	info, err := os.Stat(bookOrFile)
	if err == nil && !info.IsDir() {
		// It's a file
		bookID, err = identity.GetBookID(bookOrFile)
		if err != nil {
			return nil, fmt.Errorf("failed to compute bookID: %w", err)
		}
	} else {
		// treat as bookID
		bookID = bookOrFile
	}

	rsPath, err := state.GetStatePath(bookID)
	// handle error

	rs, err := state.LoadReadingState(rsPath)
	if err != nil {
		return nil, fmt.Errorf("reading state not found for bookID %s: %w", bookID, err)
	}
	return rs, nil
}
