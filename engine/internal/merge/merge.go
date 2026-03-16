// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package merge

import (
	"github.com/ArcaneLean/openshelf/internal/model"
)

func MergeReadingStates(a, b model.ReadingState) model.ReadingState {
	if a.BookID != b.BookID {
		panic("cannot merge reading states with different BookIDs")
	}

	if a.UpdatedAt.After(b.UpdatedAt) {
		return a
	}
	return b
}
