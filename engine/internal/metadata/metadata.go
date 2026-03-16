// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package metadata

import (
	"encoding/json"
	"os"
)

type Metadata struct {
	// Cache keys — set by the identity layer, not the extractor.
	FileSha256 string `json:"fileSha256,omitempty"`
	BookID     string `json:"bookId,omitempty"`

	// Mutable book fields — may be corrected by the user at any time.
	Title         string            `json:"title,omitempty"`
	Authors       []string          `json:"authors,omitempty"`
	Language      string            `json:"language,omitempty"`
	Publisher     string            `json:"publisher,omitempty"`
	PublishedYear int               `json:"publishedYear,omitempty"`
	Identifiers   map[string]string `json:"identifiers,omitempty"`
}

func LoadMetadata(path string) (*Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Metadata
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func SaveMetadata(path string, m *Metadata) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
