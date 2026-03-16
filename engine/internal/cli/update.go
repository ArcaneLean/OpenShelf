// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/ArcaneLean/openshelf/internal/model"
)

func UpdateState(w io.Writer, path string, locType string, value any, t time.Time) error {
	rs, err := model.LoadReadingState(path)
	if err != nil {
		return fmt.Errorf("error loading reading state: %w", err)
	}

	rs.SetLocation(locType, value, t)

	err = model.SaveReadingState(path, rs)
	if err != nil {
		return fmt.Errorf("error saving reading state: %w", err)
	}

	return nil
}
