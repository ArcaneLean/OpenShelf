package metadata

import (
	"encoding/json"
	"os"
)

type Metadata struct {
	FileSha256 string `json:"fileSha256"`
	BookID     string `json:"bookId"`
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
