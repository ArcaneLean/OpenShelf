// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeBookID(t *testing.T) {
	// Create temporary directory
	dir := t.TempDir()

	// Create temporary file
	filePath := filepath.Join(dir, "testbook.txt")
	content := []byte("OpenShelf test content")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Compute expected SHA256 manually
	expectedHash := sha256.Sum256(content)
	expected := hex.EncodeToString(expectedHash[:])

	// Call function under test
	result, err := ComputeBookID(filePath)
	if err != nil {
		t.Fatalf("ComputeBookID returned error: %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestComputeBookID_FileNotFound(t *testing.T) {
	_, err := ComputeBookID("nonexistent_file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
