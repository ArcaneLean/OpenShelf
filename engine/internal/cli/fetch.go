// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArcaneLean/openshelf/internal/identity"
	"github.com/ArcaneLean/openshelf/internal/library"
	"github.com/ArcaneLean/openshelf/internal/model"
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

	lib, err := library.Resolve()
	if err != nil {
		return nil, err
	}

	statePath := lib.StatePath(bookID)
	rs, err := model.LoadReadingState(statePath)
	if os.IsNotExist(err) {
		// First open — create and persist a new empty state.
		rs = model.NewReadingState(bookID)
		if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
			return nil, fmt.Errorf("cannot create state directory: %w", err)
		}
		if err := model.SaveReadingState(statePath, rs); err != nil {
			return nil, fmt.Errorf("cannot save new reading state: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to load reading state for bookID %s: %w", bookID, err)
	}

	return rs, nil
}
