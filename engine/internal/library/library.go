package library

import (
	"fmt"
	"os"
	"path/filepath"
)

type Library struct {
	Root string
}

func Resolve() (*Library, error) {
	root := os.Getenv("OPENSHELF_LIBRARY")
	if root == "" {
		return nil, fmt.Errorf("OPENSHELF_LIBRARY is not set")
	}

	return &Library{Root: root}, nil
}

func (l *Library) StatePath(bookId string) string {
	return filepath.Join(l.Root, ".state", bookId+".json")
}

func (l *Library) MetadataPath(hash string) string {
	return filepath.Join(l.Root, ".metadata", hash+".json")
}

func (l *Library) CanonicalPath(bookId string) string {
	return filepath.Join(l.Root, ".canonical", bookId+".json")
}
