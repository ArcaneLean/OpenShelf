package canonical

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
)

type Canonical struct {
	Schema        string            `json:"schema"`
	Title         string            `json:"title,omitempty"`
	Authors       []string          `json:"authors,omitempty"`
	Language      string            `json:"language,omitempty"`
	Publisher     string            `json:"publisher,omitempty"`
	PublishedYear int               `json:"publishedYear,omitempty"`
	Identifiers   map[string]string `json:"identifiers,omitempty"`
}

// Normalize ensures canonical identity stability.
func (c *Canonical) Normalize() {
	// Ensure schema is set
	if c.Schema == "" {
		c.Schema = "openshelf.canonical.v1"
	}

	// Trim whitespace
	c.Title = strings.TrimSpace(c.Title)
	c.Publisher = strings.TrimSpace(c.Publisher)

	// Normalize language
	c.Language = strings.ToLower(strings.TrimSpace(c.Language))

	// Normalize authors
	for i := range c.Authors {
		c.Authors[i] = strings.TrimSpace(c.Authors[i])
	}
	sort.Strings(c.Authors)

	// Normalize identifiers
	if len(c.Identifiers) > 0 {
		for k, v := range c.Identifiers {
			nk := strings.ToLower(strings.TrimSpace(k))
			nv := strings.TrimSpace(v)

			delete(c.Identifiers, k)
			if nv != "" {
				c.Identifiers[nk] = nv
			}
		}
	}
}

// deterministicJSON returns canonical JSON with stable key ordering.
func (c *Canonical) deterministicJSON() ([]byte, error) {
	// Ensure normalization first
	c.Normalize()

	// We cannot rely on default map ordering.
	// So we build a deterministic structure manually.

	type ordered struct {
		Schema        string            `json:"schema"`
		Title         string            `json:"title,omitempty"`
		Authors       []string          `json:"authors,omitempty"`
		Language      string            `json:"language,omitempty"`
		Publisher     string            `json:"publisher,omitempty"`
		PublishedYear int               `json:"publishedYear,omitempty"`
		Identifiers   map[string]string `json:"identifiers,omitempty"`
	}

	orderedIdentifiers := make(map[string]string)

	if len(c.Identifiers) > 0 {
		keys := make([]string, 0, len(c.Identifiers))
		for k := range c.Identifiers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			orderedIdentifiers[k] = c.Identifiers[k]
		}
	}

	o := ordered{
		Schema:        c.Schema,
		Title:         c.Title,
		Authors:       c.Authors,
		Language:      c.Language,
		Publisher:     c.Publisher,
		PublishedYear: c.PublishedYear,
		Identifiers:   orderedIdentifiers,
	}

	// Disable HTML escaping for strict determinism
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(o); err != nil {
		return nil, err
	}

	// Remove trailing newline added by Encoder
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// ComputeBookID returns the lowercase hex SHA256 hash
// of the deterministic canonical JSON.
func (c *Canonical) ComputeBookID() (string, error) {
	data, err := c.deterministicJSON()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
