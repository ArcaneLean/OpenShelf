package state

import (
	"encoding/json"
	"os"

	"github.com/ArcaneLean/openshelf/internal/model"
)

func LoadReadingState(path string) (*model.ReadingState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state model.ReadingState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func SaveReadingState(path string, state *model.ReadingState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModePerm)
}
