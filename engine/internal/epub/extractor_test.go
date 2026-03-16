// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package epub

import (
	"testing"

	"github.com/ArcaneLean/openshelf/internal/canonical"
)

const sampleEPUB = "../../../examples/openshelf_test/books/sample.epub"

func TestExtractMetadata(t *testing.T) {
	m, err := ExtractMetadata(sampleEPUB)
	if err != nil {
		t.Fatalf("ExtractMetadata: %v", err)
	}

	if m.Title != "Sample EPUB" {
		t.Errorf("Title = %q, want %q", m.Title, "Sample EPUB")
	}
	if m.Language != "en" {
		t.Errorf("Language = %q, want %q", m.Language, "en")
	}
	// Cache fields must not be set by the extractor
	if m.FileSha256 != "" {
		t.Errorf("FileSha256 should not be set by extractor, got %q", m.FileSha256)
	}
	if m.BookID != "" {
		t.Errorf("BookID should not be set by extractor, got %q", m.BookID)
	}
}

func TestExtractMetadata_ToCanonical(t *testing.T) {
	m, err := ExtractMetadata(sampleEPUB)
	if err != nil {
		t.Fatalf("ExtractMetadata: %v", err)
	}

	c := canonical.FromMetadata(m)
	bookID, err := c.ComputeBookID()
	if err != nil {
		t.Fatalf("ComputeBookID: %v", err)
	}
	if len(bookID) != 64 {
		t.Errorf("bookID len = %d, want 64", len(bookID))
	}

	// Stable: calling again must return the same ID
	bookID2, _ := c.ComputeBookID()
	if bookID != bookID2 {
		t.Errorf("ComputeBookID is not deterministic: %q != %q", bookID, bookID2)
	}
}

func TestExtractYear(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"2021", 2021},
		{"2019-06", 2019},
		{"2018-03-15", 2018},
		{"2020-01-01T00:00:00Z", 2020},
		{"", 0},
		{"bad", 0},
	}
	for _, tc := range cases {
		if got := extractYear(tc.input); got != tc.want {
			t.Errorf("extractYear(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
