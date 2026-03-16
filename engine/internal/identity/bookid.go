// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/ArcaneLean/openshelf/internal/library"
	"github.com/ArcaneLean/openshelf/internal/metadata"
	"github.com/ArcaneLean/openshelf/pkg/utils"
)

func GetBookID(path string) (string, error) {
	lib, err := library.Resolve()
	if err != nil {
		return "", err
	}

	hash, err := utils.HashFile(path)
	if err != nil {
		return "", err
	}

	metaPath := lib.MetadataPath(hash)
	meta, err := metadata.LoadMetadata(metaPath)
	if err != nil {
		// Cache miss — EPUB extraction not yet implemented.
		// User must register the book manually.
		return "", fmt.Errorf("no metadata found for this file; run 'openshelf register' to add it")
	}

	return meta.BookID, nil
}

func ComputeBookID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
