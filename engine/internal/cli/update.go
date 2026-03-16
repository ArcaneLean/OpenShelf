// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ArcaneLean/openshelf/internal/identity"
	"github.com/ArcaneLean/openshelf/internal/library"
	"github.com/ArcaneLean/openshelf/internal/model"
)

func UpdateState(w io.Writer, bookOrID string, locType string, value any, t time.Time) error {
	var bookID string
	info, err := os.Stat(bookOrID)
	if err == nil && !info.IsDir() {
		// It's a file path — resolve to bookID.
		bookID, err = identity.GetBookID(bookOrID)
		if err != nil {
			return fmt.Errorf("failed to compute bookID: %w", err)
		}
	} else {
		// Treat as a bookID directly.
		bookID = bookOrID
	}

	lib, err := library.Resolve()
	if err != nil {
		return err
	}

	statePath := lib.StatePath(bookID)
	rs, err := model.LoadReadingState(statePath)
	if err != nil {
		return fmt.Errorf("failed to load reading state for bookID %s: %w", bookID, err)
	}

	rs.SetLocation(locType, value, t)

	if err := model.SaveReadingState(statePath, rs); err != nil {
		return fmt.Errorf("failed to save reading state: %w", err)
	}

	return nil
}
