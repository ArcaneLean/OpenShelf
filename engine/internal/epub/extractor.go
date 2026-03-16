// Copyright (C) 2026 Aron Davids
// SPDX-License-Identifier: GPL-3.0-or-later
package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ArcaneLean/openshelf/internal/metadata"
)

// --- container.xml structs ---

type xmlContainer struct {
	Rootfiles []xmlRootfile `xml:"rootfiles>rootfile"`
}

type xmlRootfile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

// --- OPF structs ---

type xmlPackage struct {
	Metadata xmlMetadata `xml:"metadata"`
}

type xmlMetadata struct {
	Titles      []string        `xml:"http://purl.org/dc/elements/1.1/ title"`
	Creators    []string        `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Languages   []string        `xml:"http://purl.org/dc/elements/1.1/ language"`
	Publishers  []string        `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Dates       []string        `xml:"http://purl.org/dc/elements/1.1/ date"`
	Identifiers []xmlIdentifier `xml:"http://purl.org/dc/elements/1.1/ identifier"`
}

type xmlIdentifier struct {
	Value string     `xml:",chardata"`
	Attrs []xml.Attr `xml:",any,attr"`
}

// scheme returns the identifier scheme (e.g. "isbn", "uuid") from attributes,
// falling back to the element's id attribute.
func (id xmlIdentifier) scheme() string {
	for _, a := range id.Attrs {
		if strings.ToLower(a.Name.Local) == "scheme" {
			return strings.ToLower(strings.TrimSpace(a.Value))
		}
	}
	for _, a := range id.Attrs {
		if strings.ToLower(a.Name.Local) == "id" {
			return strings.ToLower(strings.TrimSpace(a.Value))
		}
	}
	return ""
}

// ExtractMetadata opens an EPUB file and returns its mutable metadata.
// FileSha256 and BookID are not set — those are filled in by the identity layer.
func ExtractMetadata(path string) (*metadata.Metadata, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open EPUB: %w", err)
	}
	defer r.Close()

	opfPath, err := findOPFPath(r)
	if err != nil {
		return nil, err
	}

	opfMeta, err := parseOPF(r, opfPath)
	if err != nil {
		return nil, err
	}

	return buildMetadata(opfMeta), nil
}

func findOPFPath(r *zip.ReadCloser) (string, error) {
	f, err := findFile(r, "META-INF/container.xml")
	if err != nil {
		return "", fmt.Errorf("missing META-INF/container.xml: %w", err)
	}

	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	var c xmlContainer
	if err := xml.Unmarshal(data, &c); err != nil {
		return "", fmt.Errorf("cannot parse container.xml: %w", err)
	}

	for _, rf := range c.Rootfiles {
		if rf.MediaType == "application/oebps-package+xml" {
			return rf.FullPath, nil
		}
	}
	if len(c.Rootfiles) > 0 {
		return c.Rootfiles[0].FullPath, nil
	}

	return "", fmt.Errorf("no rootfile found in container.xml")
}

func parseOPF(r *zip.ReadCloser, opfPath string) (*xmlMetadata, error) {
	f, err := findFile(r, opfPath)
	if err != nil {
		return nil, fmt.Errorf("OPF file %q not found in EPUB: %w", opfPath, err)
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var pkg xmlPackage
	if err := xml.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("cannot parse OPF: %w", err)
	}

	return &pkg.Metadata, nil
}

func buildMetadata(opfMeta *xmlMetadata) *metadata.Metadata {
	m := &metadata.Metadata{}

	if len(opfMeta.Titles) > 0 {
		m.Title = strings.TrimSpace(opfMeta.Titles[0])
	}

	for _, cr := range opfMeta.Creators {
		if v := strings.TrimSpace(cr); v != "" {
			m.Authors = append(m.Authors, v)
		}
	}

	if len(opfMeta.Languages) > 0 {
		m.Language = strings.TrimSpace(opfMeta.Languages[0])
	}

	if len(opfMeta.Publishers) > 0 {
		m.Publisher = strings.TrimSpace(opfMeta.Publishers[0])
	}

	if len(opfMeta.Dates) > 0 {
		m.PublishedYear = extractYear(opfMeta.Dates[0])
	}

	if len(opfMeta.Identifiers) > 0 {
		m.Identifiers = make(map[string]string)
		for _, id := range opfMeta.Identifiers {
			v := strings.TrimSpace(id.Value)
			if v == "" {
				continue
			}
			key := id.scheme()
			if key == "" {
				key = "id"
			}
			m.Identifiers[key] = v
		}
	}

	return m
}

// extractYear attempts to parse a year from common EPUB date formats:
// YYYY, YYYY-MM, YYYY-MM-DD, or RFC3339.
func extractYear(s string) int {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return 0
	}
	y, err := strconv.Atoi(s[:4])
	if err != nil {
		return 0
	}
	return y
}

// findFile returns the zip.File with the given name (case-sensitive).
func findFile(r *zip.ReadCloser, name string) (*zip.File, error) {
	for _, f := range r.File {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("file %q not found", name)
}
